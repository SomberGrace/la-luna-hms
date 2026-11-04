// frontend/src/lib/stores/roomStore.ts
import { writable, derived, get } from 'svelte/store';
import type { Room, Folio, OfflineAction, RoomStatus } from '#lib/types/room';
import { apiFetch } from '#lib/api';

// Rooms are loaded from the live database; no hardcoded defaults.
const INITIAL_ROOMS: Room[] = [];

const INITIAL_FOLIOS: Record<string, Folio> = {};

const API_BASE = '/api/v1/frontdesk';

// Storage Keys
const ROOMS_CACHE_KEY = 'la_luna_rooms_cache';
const FOLIOS_CACHE_KEY = 'la_luna_folios_cache';
const OFFLINE_QUEUE_KEY = 'la_luna_offline_queue';

// Online status store
export const isOnline = writable<boolean>(typeof navigator !== 'undefined' ? navigator.onLine : true);

if (typeof window !== 'undefined') {
  window.addEventListener('online', () => {
    isOnline.set(true);
    roomStore.syncOfflineQueue();
  });
  window.addEventListener('offline', () => {
    isOnline.set(false);
  });
}

function createRoomStore() {
  // Initialize state with cache or initial values
  let cachedRooms: Room[] = INITIAL_ROOMS;
  let cachedFolios: Record<string, Folio> = INITIAL_FOLIOS;
  let cachedQueue: OfflineAction[] = [];

  if (typeof window !== 'undefined') {
    try {
      const savedRooms = localStorage.getItem(ROOMS_CACHE_KEY);
      if (savedRooms) cachedRooms = JSON.parse(savedRooms);
      const savedFolios = localStorage.getItem(FOLIOS_CACHE_KEY);
      if (savedFolios) cachedFolios = JSON.parse(savedFolios);
      const savedQueue = localStorage.getItem(OFFLINE_QUEUE_KEY);
      if (savedQueue) cachedQueue = JSON.parse(savedQueue);
    } catch (e) {
      console.warn('Failed to parse cached local data:', e);
    }
  }

  const rooms = writable<Room[]>(cachedRooms);
  const folios = writable<Record<string, Folio>>(cachedFolios);
  const offlineQueue = writable<OfflineAction[]>(cachedQueue);
  const isLoading = writable<boolean>(false);
  const syncStatus = writable<'idle' | 'syncing' | 'synced' | 'error'>('idle');

  function saveRoomsLocal(data: Room[]) {
    rooms.set(data);
    if (typeof window !== 'undefined') {
      localStorage.setItem(ROOMS_CACHE_KEY, JSON.stringify(data));
    }
  }

  function saveFoliosLocal(data: Record<string, Folio>) {
    folios.set(data);
    if (typeof window !== 'undefined') {
      localStorage.setItem(FOLIOS_CACHE_KEY, JSON.stringify(data));
    }
  }

  function saveQueueLocal(data: OfflineAction[]) {
    offlineQueue.set(data);
    if (typeof window !== 'undefined') {
      localStorage.setItem(OFFLINE_QUEUE_KEY, JSON.stringify(data));
    }
  }

  // Network Fetch with fallback
  async function fetchRooms() {
    isLoading.set(true);
    try {
      const res = await apiFetch(`${API_BASE}/rooms`, { method: 'GET' });
      if (!res.ok) throw new Error(`HTTP error ${res.status}`);
      const data: Room[] = await res.json();
      saveRoomsLocal(data);
      isOnline.set(true);
    } catch (err) {
      console.warn('Backend unavailable, using offline/cached rooms:', err);
      // Keep cached data
    } finally {
      isLoading.set(false);
    }
  }

  async function fetchFolio(roomNumber: string): Promise<Folio | null> {
    try {
      const res = await apiFetch(`${API_BASE}/folio?room_number=${roomNumber}`);
      if (res.ok) {
        const data: Folio = await res.json();
        folios.update((all) => {
          const next = { ...all, [roomNumber]: data };
          saveFoliosLocal(next);
          return next;
        });
        return data;
      }
    } catch (err) {
      console.warn('Backend unavailable, using cached folio:', err);
    }
    // Fallback to local cache
    const current = get(folios);
    return current[roomNumber] || null;
  }

  async function checkIn(payload: {
    guest_name: string;
    guest_phone?: string;
    guest_id_doc?: string;
    room_number: string;
    nights: number;
    deposit_paid: number;
    reservation_code?: string;
  }): Promise<{ success: boolean; folio?: Folio; message: string }> {
    const currentRooms = get(rooms);
    const room = currentRooms.find((r) => r.room_number === payload.room_number);
    if (!room) {
      return { success: false, message: 'Room not found' };
    }
    // Only new walk-ins require room to be currently 'Available'.
    // Online reservations already hold the room (status is 'Reserved' or 'Available').
    if (!payload.reservation_code && room.status !== 'Available') {
      return { success: false, message: 'Room is not available' };
    }

    try {
      const res = await apiFetch(`${API_BASE}/checkin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        const data = await res.json();
        // Update local room status
        const updatedRooms = currentRooms.map((r) =>
          r.room_number === payload.room_number ? { ...r, status: 'Occupied' as RoomStatus } : r
        );
        saveRoomsLocal(updatedRooms);
        // Update local folio
        if (data.folio) {
          folios.update((f) => {
            const next = { ...f, [payload.room_number]: data.folio };
            saveFoliosLocal(next);
            return next;
          });
        }
        return { success: true, folio: data.folio, message: 'Check-in successful' };
      }
    } catch (err) {
      console.warn('Network failed during check-in, queueing offline operation:', err);
    }

    // Offline optimistic execution
    const totalCharge = room.price_per_night * payload.nights;
    const balance = totalCharge - (payload.deposit_paid || 0);
    const nowIso = new Date().toISOString();
    const today = nowIso.split('T')[0];
    const checkoutDate = new Date(Date.now() + payload.nights * 24 * 3600 * 1000).toISOString().split('T')[0];

    const localFolio: Folio = {
      id: `fol-${payload.room_number}`,
      guest_id: `g-offline-${Date.now()}`,
      guest_name: payload.guest_name,
      guest_phone: payload.guest_phone,
      guest_id_doc: payload.guest_id_doc,
      room_number: payload.room_number,
      check_in_date: today,
      check_out_date: checkoutDate,
      total_charges: totalCharge,
      total_paid: payload.deposit_paid || 0,
      balance: balance,
      items: [
        {
          id: `item-${Date.now()}`,
          description: `${payload.nights} Night(s) Accommodation - ${room.type}`,
          amount: totalCharge,
          category: 'Accommodation',
          timestamp: nowIso
        }
      ]
    };

    if (payload.deposit_paid > 0) {
      localFolio.items.push({
        id: `pay-${Date.now()}`,
        description: 'Initial Deposit Received',
        amount: -payload.deposit_paid,
        category: 'Payment',
        timestamp: nowIso
      });
    }

    // Update local state
    const updatedRooms = currentRooms.map((r) =>
      r.room_number === payload.room_number ? { ...r, status: 'Occupied' as RoomStatus } : r
    );
    saveRoomsLocal(updatedRooms);

    folios.update((f) => {
      const next = { ...f, [payload.room_number]: localFolio };
      saveFoliosLocal(next);
      return next;
    });

    // Queue action
    const action: OfflineAction = {
      id: `act-${Date.now()}`,
      action_type: 'checkin',
      payload,
      timestamp: nowIso,
      status: 'pending'
    };
    saveQueueLocal([...get(offlineQueue), action]);

    return {
      success: true,
      folio: localFolio,
      message: 'Check-in processed locally (queued for server sync)'
    };
  }

  async function recordPayment(payload: {
    room_number: string;
    amount: number;
    method: string;
    reference?: string;
  }): Promise<{ success: boolean; folio?: Folio; message: string }> {
    try {
      const res = await apiFetch(`${API_BASE}/payment`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        const data = await res.json();
        folios.update((f) => {
          const next = { ...f, [payload.room_number]: data.folio };
          saveFoliosLocal(next);
          return next;
        });
        return { success: true, folio: data.folio, message: 'Payment recorded successfully' };
      }
    } catch (err) {
      console.warn('Network failed during payment, queuing offline operation:', err);
    }

    // Offline optimistic update
    const currentFolios = get(folios);
    const folio = currentFolios[payload.room_number];
    if (!folio) {
      return { success: false, message: 'No active folio found for this room' };
    }

    const updatedFolio: Folio = {
      ...folio,
      total_paid: folio.total_paid + payload.amount,
      balance: folio.balance - payload.amount,
      items: [
        ...folio.items,
        {
          id: `pay-${Date.now()}`,
          description: `Payment Received (${payload.method})${payload.reference ? ` - Ref: ${payload.reference}` : ''}`,
          amount: -payload.amount,
          category: 'Payment',
          timestamp: new Date().toISOString()
        }
      ]
    };

    saveFoliosLocal({ ...currentFolios, [payload.room_number]: updatedFolio });

    const action: OfflineAction = {
      id: `act-${Date.now()}`,
      action_type: 'payment',
      payload,
      timestamp: new Date().toISOString(),
      status: 'pending'
    };
    saveQueueLocal([...get(offlineQueue), action]);

    return {
      success: true,
      folio: updatedFolio,
      message: 'Payment recorded locally (queued for server sync)'
    };
  }

  async function checkOut(payload: { room_number: string }): Promise<{ success: boolean; message: string }> {
    const currentFolios = get(folios);
    const folio = currentFolios[payload.room_number];

    if (folio && folio.balance > 0) {
      return {
        success: false,
        message: `Cannot check out: Outstanding balance of MWK ${folio.balance.toLocaleString()} remaining`
      };
    }

    try {
      const res = await apiFetch(`${API_BASE}/checkout`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        // Room status to Available (released for booking)
        const updatedRooms = get(rooms).map((r) =>
          r.room_number === payload.room_number ? { ...r, status: 'Available' as RoomStatus } : r
        );
        saveRoomsLocal(updatedRooms);
        const { [payload.room_number]: _, ...rest } = currentFolios;
        saveFoliosLocal(rest);
        return { success: true, message: 'Check-out completed successfully. Room is now Available.' };
      }
      const errData = await res.json();
      return { success: false, message: errData.error || errData.message || 'Check-out rejected' };
    } catch (err) {
      console.warn('Network failed during checkout, queuing offline operation:', err);
    }

    // Offline optimistic execution
    const updatedRooms = get(rooms).map((r) =>
      r.room_number === payload.room_number ? { ...r, status: 'Available' as RoomStatus } : r
    );
    saveRoomsLocal(updatedRooms);
    const { [payload.room_number]: _, ...rest } = currentFolios;
    saveFoliosLocal(rest);

    const action: OfflineAction = {
      id: `act-${Date.now()}`,
      action_type: 'checkout',
      payload,
      timestamp: new Date().toISOString(),
      status: 'pending'
    };
    saveQueueLocal([...get(offlineQueue), action]);

    return {
      success: true,
      message: 'Check-out processed locally. Room marked Cleaning (queued for server sync).'
    };
  }

  async function extendStay(payload: {
    room_number: string;
    additional_nights: number;
    pay_now?: boolean;
    payment_method?: string;
  }): Promise<{ success: boolean; folio?: Folio; message: string }> {
    try {
      const res = await apiFetch(`${API_BASE}/extension`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        const data = await res.json();
        folios.update((f) => {
          const next = { ...f, [payload.room_number]: data.folio };
          saveFoliosLocal(next);
          return next;
        });
        return { success: true, folio: data.folio, message: 'Stay extended successfully' };
      }
    } catch (err) {
      console.warn('Network failed during extension, queuing offline operation:', err);
    }

    // Offline optimistic execution
    const currentFolios = get(folios);
    const folio = currentFolios[payload.room_number];
    const room = get(rooms).find((r) => r.room_number === payload.room_number);
    if (!folio || !room) {
      return { success: false, message: 'Stay or room not found' };
    }

    const addCharge = room.price_per_night * payload.additional_nights;
    const currentCO = new Date(folio.check_out_date);
    const newCO = new Date(currentCO.getTime() + payload.additional_nights * 24 * 3600 * 1000).toISOString().split('T')[0];

    const updatedFolio: Folio = {
      ...folio,
      check_out_date: newCO,
      total_charges: folio.total_charges + addCharge,
      balance: folio.balance + (payload.pay_now ? 0 : addCharge),
      total_paid: folio.total_paid + (payload.pay_now ? addCharge : 0),
      items: [
        ...folio.items,
        {
          id: `ext-${Date.now()}`,
          description: `Stay Extension: +${payload.additional_nights} Night(s) (${room.type})`,
          amount: addCharge,
          category: 'Accommodation',
          timestamp: new Date().toISOString()
        }
      ]
    };

    if (payload.pay_now) {
      updatedFolio.items.push({
        id: `pay-ext-${Date.now()}`,
        description: `Extension Payment Received (${payload.payment_method || 'Cash'})`,
        amount: -addCharge,
        category: 'Payment',
        timestamp: new Date().toISOString()
      });
    }

    saveFoliosLocal({ ...currentFolios, [payload.room_number]: updatedFolio });

    const action: OfflineAction = {
      id: `act-${Date.now()}`,
      action_type: 'extension',
      payload,
      timestamp: new Date().toISOString(),
      status: 'pending'
    };
    saveQueueLocal([...get(offlineQueue), action]);

    return {
      success: true,
      folio: updatedFolio,
      message: 'Stay extended locally (queued for server sync)'
    };
  }

  async function updateRoomStatus(
    roomNumber: string,
    status: RoomStatus,
    housekeeper?: string,
    priority?: 'Normal' | 'High' | 'Urgent'
  ): Promise<{ success: boolean; message: string }> {
    try {
      const res = await apiFetch(`${API_BASE}/rooms/status`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ room_number: roomNumber, status, housekeeper, priority })
      });
      if (res.ok) {
        const updated = get(rooms).map((r) =>
          r.room_number === roomNumber ? { ...r, status, housekeeper, priority } : r
        );
        saveRoomsLocal(updated);
        return { success: true, message: `Room ${roomNumber} updated to ${status}` };
      }
    } catch (err) {
      console.warn('Network failed during status update:', err);
    }

    // Local optimistic update
    const updated = get(rooms).map((r) =>
      r.room_number === roomNumber ? { ...r, status, housekeeper, priority } : r
    );
    saveRoomsLocal(updated);

    const action: OfflineAction = {
      id: `act-${Date.now()}`,
      action_type: 'room_status',
      payload: { room_number: roomNumber, status, housekeeper, priority },
      timestamp: new Date().toISOString(),
      status: 'pending'
    };
    saveQueueLocal([...get(offlineQueue), action]);

    return { success: true, message: `Room ${roomNumber} updated locally to ${status}` };
  }

  async function syncOfflineQueue(): Promise<{ synced: number; failed: number }> {
    const queue = get(offlineQueue);
    if (queue.length === 0) return { synced: 0, failed: 0 };

    syncStatus.set('syncing');
    let synced = 0;
    let failed = 0;
    const remainingQueue: OfflineAction[] = [];

    for (const item of queue) {
      try {
        let endpoint = '';
        if (item.action_type === 'checkin') endpoint = `${API_BASE}/checkin`;
        else if (item.action_type === 'payment') endpoint = `${API_BASE}/payment`;
        else if (item.action_type === 'checkout') endpoint = `${API_BASE}/checkout`;
        else if (item.action_type === 'extension') endpoint = `${API_BASE}/extension`;
        else if (item.action_type === 'room_status') endpoint = `${API_BASE}/rooms/status`;

        const res = await apiFetch(endpoint, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(item.payload)
        });

        if (res.ok) {
          synced++;
        } else {
          failed++;
          remainingQueue.push({ ...item, status: 'failed', error_message: `HTTP ${res.status}` });
        }
      } catch (err: any) {
        failed++;
        remainingQueue.push({ ...item, status: 'failed', error_message: err.message });
      }
    }

    saveQueueLocal(remainingQueue);
    syncStatus.set(failed === 0 ? 'synced' : 'error');
    // Refresh rooms and folios from server if any synced
    if (synced > 0) {
      await fetchRooms();
    }

    return { synced, failed };
  }

  function clearQueue() {
    saveQueueLocal([]);
  }

  return {
    rooms: { subscribe: rooms.subscribe },
    folios: { subscribe: folios.subscribe },
    offlineQueue: { subscribe: offlineQueue.subscribe },
    isLoading: { subscribe: isLoading.subscribe },
    syncStatus: { subscribe: syncStatus.subscribe },
    fetchRooms,
    fetchFolio,
    checkIn,
    recordPayment,
    checkOut,
    extendStay,
    updateRoomStatus,
    syncOfflineQueue,
    clearQueue
  };
}

export const roomStore = createRoomStore();

// Derived Stores for Quick Mobile Metrics
export const availableRoomsCount = derived(roomStore.rooms, ($rooms) =>
  $rooms.filter((r) => r.status === 'Available').length
);

export const occupiedRoomsCount = derived(roomStore.rooms, ($rooms) =>
  $rooms.filter((r) => r.status === 'Occupied').length
);

export const cleaningRoomsCount = derived(roomStore.rooms, ($rooms) =>
  $rooms.filter((r) => r.status === 'Cleaning').length
);

export const outOfServiceRoomsCount = derived(roomStore.rooms, ($rooms) =>
  $rooms.filter((r) => r.status === 'OutOfService').length
);

export const pendingOfflineCount = derived(roomStore.offlineQueue, ($queue) => $queue.length);
