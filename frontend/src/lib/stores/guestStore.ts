// frontend/src/lib/stores/guestStore.ts
import { writable, get } from 'svelte/store';

export interface GuestStay {
  cardNo: string;
  guestName: string;
  phone: string;
  email: string;
  roomNumber: string;
  roomType: 'Standard' | 'Executive' | 'Apartment';
  bedType: 'Single' | 'Double';
  ratePerNight: number;
  arrivalDate: string;
  departureDate: string;
  nights: number;
  wifiSsid: string;
  wifiPass: string;
  dndActive: boolean;
  makeupRoomActive: boolean;
  digitalKeyUnlocked: boolean;
}

export interface DiningOrderItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  specialInstructions?: string;
}

export interface DiningOrder {
  id: string;
  orderNumber: string;
  roomNumber: string;
  items: DiningOrderItem[];
  totalAmount: number;
  status: 'Received' | 'Cooking' | 'Dispatched' | 'Delivered';
  createdAt: string;
  deliveryNotes?: string;
}

export interface ServiceRequest {
  id: string;
  type: string;
  icon: string;
  details: string;
  roomNumber: string;
  status: 'Pending' | 'Attending' | 'Completed';
  createdAt: string;
}

const DEFAULT_STAY: GuestStay = {
  cardNo: '',
  guestName: '',
  phone: '',
  email: '',
  roomNumber: '',
  roomType: 'Standard',
  bedType: 'Single',
  ratePerNight: 0,
  arrivalDate: '',
  departureDate: '',
  nights: 0,
  wifiSsid: 'LaLuna-Resident-5G',
  wifiPass: 'LunaResort2025!',
  dndActive: false,
  makeupRoomActive: false,
  digitalKeyUnlocked: false
};

const DEFAULT_ORDERS: DiningOrder[] = [];

const DEFAULT_SERVICES: ServiceRequest[] = [];

const GUEST_STORAGE_KEY = 'la_luna_guest_stay_v1';
const GUEST_ORDERS_KEY = 'la_luna_guest_orders_v1';
const GUEST_SERVICES_KEY = 'la_luna_guest_services_v1';

function createGuestStore() {
  let initialStay = DEFAULT_STAY;
  let initialOrders = DEFAULT_ORDERS;
  let initialServices = DEFAULT_SERVICES;

  if (typeof window !== 'undefined') {
    try {
      const savedStay = localStorage.getItem(GUEST_STORAGE_KEY);
      if (savedStay) initialStay = JSON.parse(savedStay);

      const savedOrders = localStorage.getItem(GUEST_ORDERS_KEY);
      if (savedOrders) initialOrders = JSON.parse(savedOrders);

      const savedServices = localStorage.getItem(GUEST_SERVICES_KEY);
      if (savedServices) initialServices = JSON.parse(savedServices);
    } catch (e) {
      console.warn('Failed to parse local guest data:', e);
    }
  }

  const stay = writable<GuestStay>(initialStay);
  const orders = writable<DiningOrder[]>(initialOrders);
  const services = writable<ServiceRequest[]>(initialServices);

  function saveStay(updated: GuestStay) {
    stay.set(updated);
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem(GUEST_STORAGE_KEY, JSON.stringify(updated));
      } catch (e) {
        console.warn('Failed to save stay:', e);
      }
    }
  }

  function saveOrders(updated: DiningOrder[]) {
    orders.set(updated);
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem(GUEST_ORDERS_KEY, JSON.stringify(updated));
      } catch (e) {
        console.warn('Failed to save orders:', e);
      }
    }
  }

  function saveServices(updated: ServiceRequest[]) {
    services.set(updated);
    if (typeof window !== 'undefined') {
      try {
        localStorage.setItem(GUEST_SERVICES_KEY, JSON.stringify(updated));
      } catch (e) {
        console.warn('Failed to save services:', e);
      }
    }
  }

  return {
    stay: {
      subscribe: stay.subscribe,
      update: stay.update,
      set: saveStay,
      toggleDnd: () => {
        stay.update(s => {
          const updated = { ...s, dndActive: !s.dndActive, makeupRoomActive: s.dndActive ? s.makeupRoomActive : false };
          saveStay(updated);
          return updated;
        });
      },
      toggleMakeupRoom: () => {
        stay.update(s => {
          const updated = { ...s, makeupRoomActive: !s.makeupRoomActive, dndActive: s.makeupRoomActive ? s.dndActive : false };
          saveStay(updated);
          return updated;
        });
      },
      unlockDigitalKey: () => {
        stay.update(s => {
          const updated = { ...s, digitalKeyUnlocked: true };
          saveStay(updated);
          return updated;
        });
        setTimeout(() => {
          stay.update(s => {
            const relocked = { ...s, digitalKeyUnlocked: false };
            saveStay(relocked);
            return relocked;
          });
        }, 5000);
      },
      switchStay: (roomNumber: string, guestName: string) => {
        stay.update(s => {
          const updated: GuestStay = {
            ...s,
            roomNumber,
            guestName,
            cardNo: `CARD-${roomNumber}`
          };
          saveStay(updated);
          return updated;
        });
      }
    },
    orders: {
      subscribe: orders.subscribe,
      addOrder: (newOrder: DiningOrder) => {
        orders.update(prev => {
          const updated = [newOrder, ...prev];
          saveOrders(updated);
          return updated;
        });
      },
      updateStatus: (orderId: string, status: DiningOrder['status']) => {
        orders.update(prev => {
          const updated = prev.map(o => o.id === orderId ? { ...o, status } : o);
          saveOrders(updated);
          return updated;
        });
      }
    },
    services: {
      subscribe: services.subscribe,
      addRequest: (type: string, details: string, icon: string = 'ph-bell') => {
        const curStay = get(stay);
        const req: ServiceRequest = {
          id: 'srv-' + Date.now(),
          type,
          icon,
          details,
          roomNumber: curStay.roomNumber,
          status: 'Pending',
          createdAt: new Date().toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' }) + ' Today'
        };
        services.update(prev => {
          const updated = [req, ...prev];
          saveServices(updated);
          return updated;
        });
      },
      updateStatus: (id: string, status: ServiceRequest['status']) => {
        services.update(prev => {
          const updated = prev.map(r => r.id === id ? { ...r, status } : r);
          saveServices(updated);
          return updated;
        });
      }
    }
  };
}

export const guestStore = createGuestStore();

