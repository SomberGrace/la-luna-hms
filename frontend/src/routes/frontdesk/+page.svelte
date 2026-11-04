<!-- frontend/src/routes/frontdesk/+page.svelte -->
<script lang="ts">
  import {
    roomStore,
    availableRoomsCount,
    occupiedRoomsCount,
    cleaningRoomsCount,
    outOfServiceRoomsCount,
    pendingOfflineCount,
    isOnline
  } from '#lib/stores/roomStore';
  import { authStore } from '#lib/stores/authStore';
  import { apiFetch } from '#lib/api';
  import { getPendingArrivals } from '#lib/frontdeskReservations';
  import { onMount } from 'svelte';
  import type { ReservationItem } from '#lib/types/guest';

  const rooms = roomStore.rooms;
  const folios = roomStore.folios;

  let searchQuery = $state('');
  let pendingArrivals = $state<ReservationItem[]>([]);
  let loadingArrivals = $state(false);
  let checkInStatusMsg = $state<{ text: string; isError: boolean } | null>(null);

  onMount(async () => {
    await roomStore.fetchRooms();
    await loadPendingArrivals();
  });

  async function loadPendingArrivals() {
    loadingArrivals = true;
    try {
      const res = await apiFetch('/api/v1/frontdesk/reservations');
      if (res.ok) {
        const data: unknown = await res.json();
        const arrivals = getPendingArrivals(data);
        if (arrivals) pendingArrivals = arrivals;
      }
    } catch (e) {
      console.warn('Failed to load pending arrivals:', e);
    } finally {
      loadingArrivals = false;
    }
  }

  async function quickClearArrival(r: ReservationItem) {
    checkInStatusMsg = null;
    try {
      const res = await roomStore.checkIn({
        guest_name: r.guest_name,
        guest_phone: r.guest_phone || 'N/A',
        guest_id_doc: r.guest_id_doc || 'Online Verified ID',
        room_number: r.room_number,
        nights: Number(r.nights || 1),
        deposit_paid: Number(r.deposit_paid || 0),
        reservation_code: r.res_code
      });

      if (res.success) {
        checkInStatusMsg = { text: `Guest ${r.guest_name} cleared! Checked in to Room ${r.room_number}.`, isError: false };
        pendingArrivals = pendingArrivals.filter((p) => p.res_code !== r.res_code);
        await roomStore.fetchRooms();
        await roomStore.fetchFolio(r.room_number);
      } else {
        checkInStatusMsg = { text: res.message || 'Check-in failed.', isError: true };
      }
    } catch (err: any) {
      checkInStatusMsg = { text: err.message || 'Check-in failed.', isError: true };
    }
    setTimeout(() => { checkInStatusMsg = null; }, 5000);
  }

  // Rooms that are still pending arrival
  let pendingRoomNumbers = $derived(new Set(pendingArrivals.map((p) => p.room_number)));

  // Filter occupied rooms for quick stay monitoring — strictly exclude any room that is still in pending arrivals
  let activeStays = $derived(
    $rooms
      .filter((r) => r.status === 'Occupied' && !pendingRoomNumbers.has(r.room_number))
      .map((r) => {
        const folio = $folios[r.room_number];
        return {
          room: r,
          folio: folio || null
        };
      })
      .filter((stay) => {
        if (!searchQuery.trim()) return true;
        const q = searchQuery.toLowerCase();
        return (
          stay.room.room_number.toLowerCase().includes(q) ||
          (stay.folio?.guest_name && stay.folio.guest_name.toLowerCase().includes(q))
        );
      })
  );

  // Reactive fetch live folios for all occupied rooms
  $effect(() => {
    const occupied = $rooms.filter((r) => r.status === 'Occupied');
    for (const r of occupied) {
      if (!$folios[r.room_number]) {
        roomStore.fetchFolio(r.room_number);
      }
    }
  });

  function formatMWK(amount: number): string {
    return 'MWK ' + (amount || 0).toLocaleString();
  }
</script>

<div class="space-y-5">
  <!-- Officer Duty Card -->
  <div class="bg-gradient-to-r from-slate-900 to-slate-800 border border-slate-800 rounded-2xl p-4 shadow-lg flex items-center justify-between">
    <div>
      <p class="text-xs text-amber-400 font-semibold tracking-wide uppercase">On-Duty Front Desk</p>
      <h2 class="text-lg font-bold text-white tracking-tight">{$authStore.fullName}</h2>
      <p class="text-[12px] text-slate-400 flex items-center mt-0.5">
        <i class="ph ph-clock mr-1 text-amber-400"></i> {$authStore.shift}
      </p>
    </div>
    <div class="text-right">
      <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-semibold {$isOnline ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">
        <span class="w-1.5 h-1.5 rounded-full mr-1.5 {$isOnline ? 'bg-emerald-400' : 'bg-rose-400'}"></span>
        {$isOnline ? 'Online' : 'Offline Mode'}
      </span>
      <p class="text-[11px] text-slate-500 mt-1 font-mono">Standard C/O: 10:00 AM</p>
    </div>
  </div>

  <!-- Operational Status KPI 4-Grid -->
  <div class="grid grid-cols-2 gap-3">
    <!-- Available -->
    <a
      href="/frontdesk/rooms?status=Available"
      class="bg-slate-900 border border-emerald-800/40 hover:border-emerald-500/50 p-3.5 rounded-xl shadow-sm transition active:scale-95 group"
    >
      <div class="flex items-center justify-between">
        <span class="text-xs font-semibold text-emerald-400">Available</span>
        <div class="w-7 h-7 rounded-lg bg-emerald-500/10 flex items-center justify-center text-emerald-400">
          <i class="ph ph-door-open text-base"></i>
        </div>
      </div>
      <div class="mt-2 flex items-baseline justify-between">
        <span class="text-2xl font-bold font-['Oxanium',sans-serif] text-white">{$availableRoomsCount}</span>
        <span class="text-[10px] text-slate-400 group-hover:text-emerald-400 transition-colors">Assign &rarr;</span>
      </div>
    </a>

    <!-- Occupied -->
    <a
      href="/frontdesk/rooms?status=Occupied"
      class="bg-slate-900 border border-amber-800/40 hover:border-amber-500/50 p-3.5 rounded-xl shadow-sm transition active:scale-95 group"
    >
      <div class="flex items-center justify-between">
        <span class="text-xs font-semibold text-amber-400">Occupied</span>
        <div class="w-7 h-7 rounded-lg bg-amber-500/10 flex items-center justify-center text-amber-400">
          <i class="ph ph-key text-base"></i>
        </div>
      </div>
      <div class="mt-2 flex items-baseline justify-between">
        <span class="text-2xl font-bold font-['Oxanium',sans-serif] text-white">{$occupiedRoomsCount}</span>
        <span class="text-[10px] text-slate-400 group-hover:text-amber-400 transition-colors">Stays &rarr;</span>
      </div>
    </a>

    <!-- Cleaning / Housekeeping -->
    <a
      href="/frontdesk/housekeeping"
      class="bg-slate-900 border border-sky-800/40 hover:border-sky-500/50 p-3.5 rounded-xl shadow-sm transition active:scale-95 group"
    >
      <div class="flex items-center justify-between">
        <span class="text-xs font-semibold text-sky-400">Housekeeping</span>
        <div class="w-7 h-7 rounded-lg bg-sky-500/10 flex items-center justify-center text-sky-400">
          <i class="ph ph-broom text-base"></i>
        </div>
      </div>
      <div class="mt-2 flex items-baseline justify-between">
        <span class="text-2xl font-bold font-['Oxanium',sans-serif] text-white">{$cleaningRoomsCount}</span>
        <span class="text-[10px] text-slate-400 group-hover:text-sky-400 transition-colors">Track &rarr;</span>
      </div>
    </a>

    <!-- Out Of Service -->
    <a
      href="/frontdesk/rooms?status=OutOfService"
      class="bg-slate-900 border border-rose-800/40 hover:border-rose-500/50 p-3.5 rounded-xl shadow-sm transition active:scale-95 group"
    >
      <div class="flex items-center justify-between">
        <span class="text-xs font-semibold text-rose-400">Out of Service</span>
        <div class="w-7 h-7 rounded-lg bg-rose-500/10 flex items-center justify-center text-rose-400">
          <i class="ph ph-wrench text-base"></i>
        </div>
      </div>
      <div class="mt-2 flex items-baseline justify-between">
        <span class="text-2xl font-bold font-['Oxanium',sans-serif] text-white">{$outOfServiceRoomsCount}</span>
        <span class="text-[10px] text-slate-400 group-hover:text-rose-400 transition-colors">View &rarr;</span>
      </div>
    </a>
  </div>

  <!-- Command Actions 6-Tile Fast Access -->
  <div>
    <h3 class="text-xs font-semibold uppercase tracking-wider text-slate-400 mb-2.5 flex items-center justify-between">
      <span>Command Actions</span>
      <span class="text-[11px] text-amber-500 lowercase">touch optimized</span>
    </h3>
    <div class="grid grid-cols-3 gap-2.5">
      <a
        href="/frontdesk/walkin"
        class="bg-slate-900 border border-amber-600/30 hover:border-amber-500 p-3 rounded-xl flex flex-col items-center text-center shadow-md active:scale-95 transition"
      >
        <div class="w-10 h-10 rounded-full bg-amber-500/15 text-amber-400 flex items-center justify-center mb-1.5">
          <i class="ph ph-user-plus text-xl"></i>
        </div>
        <span class="text-xs font-bold text-white">Walk-In</span>
        <span class="text-[10px] text-slate-400">Rapid Register</span>
      </a>

      <a
        href="/frontdesk/checkin"
        class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-3 rounded-xl flex flex-col items-center text-center shadow-md active:scale-95 transition"
      >
        <div class="w-10 h-10 rounded-full bg-emerald-500/15 text-emerald-400 flex items-center justify-center mb-1.5">
          <i class="ph ph-identification-card text-xl"></i>
        </div>
        <span class="text-xs font-bold text-white">Check-In</span>
        <span class="text-[10px] text-slate-400">Reservations</span>
      </a>

      <a
        href="/frontdesk/payments"
        class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-3 rounded-xl flex flex-col items-center text-center shadow-md active:scale-95 transition"
      >
        <div class="w-10 h-10 rounded-full bg-teal-500/15 text-teal-400 flex items-center justify-center mb-1.5">
          <i class="ph ph-currency-circle-dollar text-xl"></i>
        </div>
        <span class="text-xs font-bold text-white">Payments</span>
        <span class="text-[10px] text-slate-400">Settle Balance</span>
      </a>

      <a
        href="/frontdesk/extension"
        class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-3 rounded-xl flex flex-col items-center text-center shadow-md active:scale-95 transition"
      >
        <div class="w-10 h-10 rounded-full bg-indigo-500/15 text-indigo-400 flex items-center justify-center mb-1.5">
          <i class="ph ph-calendar-plus text-xl"></i>
        </div>
        <span class="text-xs font-bold text-white">Extension</span>
        <span class="text-[10px] text-slate-400">Extend Stay</span>
      </a>

      <a
        href="/frontdesk/checkout"
        class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-3 rounded-xl flex flex-col items-center text-center shadow-md active:scale-95 transition"
      >
        <div class="w-10 h-10 rounded-full bg-rose-500/15 text-rose-400 flex items-center justify-center mb-1.5">
          <i class="ph ph-sign-out text-xl"></i>
        </div>
        <span class="text-xs font-bold text-white">Checkout</span>
        <span class="text-[10px] text-slate-400">Zero-Balance</span>
      </a>

      <a
        href="/frontdesk/reports"
        class="bg-slate-900 border border-slate-800 hover:border-slate-700 p-3 rounded-xl flex flex-col items-center text-center shadow-md active:scale-95 transition"
      >
        <div class="w-10 h-10 rounded-full bg-amber-500/15 text-amber-400 flex items-center justify-center mb-1.5">
          <i class="ph ph-chart-bar text-xl"></i>
        </div>
        <span class="text-xs font-bold text-white">Reports</span>
        <span class="text-[10px] text-slate-400">Reception Shift</span>
      </a>
    </div>
  </div>

  <!-- Feedback banner -->
  {#if checkInStatusMsg}
    <div class="p-3.5 rounded-2xl border text-xs flex items-center justify-between {checkInStatusMsg.isError ? 'bg-rose-500/10 border-rose-500/30 text-rose-300' : 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300'}">
      <div class="flex items-center space-x-2">
        <i class="ph {checkInStatusMsg.isError ? 'ph-warning-circle' : 'ph-check-circle'} text-base"></i>
        <span class="font-bold">{checkInStatusMsg.text}</span>
      </div>
      <button onclick={() => (checkInStatusMsg = null)} class="text-slate-400 hover:text-white text-xs">&times;</button>
    </div>
  {/if}

  <!-- Pending Arrivals Clearance Card -->
  <div class="bg-slate-900 border border-amber-500/30 rounded-2xl p-4 shadow-lg space-y-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-2">
        <div class="w-2.5 h-2.5 rounded-full bg-amber-400 animate-pulse"></div>
        <h3 class="font-bold text-sm text-white">Pending Expected Arrivals</h3>
      </div>
      <div class="flex items-center space-x-2">
        <span class="text-xs font-mono font-bold bg-amber-500/10 text-amber-300 border border-amber-500/20 px-2 py-0.5 rounded-full">
          {pendingArrivals.length} pending
        </span>
        <button
          onclick={loadPendingArrivals}
          disabled={loadingArrivals}
          class="text-[11px] text-slate-400 hover:text-amber-400 p-1 rounded hover:bg-slate-800 transition"
          title="Refresh"
        >
          <i class="ph ph-arrows-clockwise {loadingArrivals ? 'animate-spin' : ''}"></i>
        </button>
      </div>
    </div>

    {#if loadingArrivals && pendingArrivals.length === 0}
      <p class="text-xs text-slate-500 py-3 text-center">Loading pending arrivals...</p>
    {:else if pendingArrivals.length === 0}
      <div class="py-3 px-3.5 bg-slate-950/60 rounded-xl border border-slate-800/80 text-center text-xs text-slate-400">
        No pending arrivals waiting. Online bookings appear here until checked in at Front Desk.
      </div>
    {:else}
      <div class="space-y-2 max-h-56 overflow-y-auto pr-1">
        {#each pendingArrivals as r}
          <div class="bg-slate-950/80 border border-slate-800/80 hover:border-amber-500/40 rounded-xl p-3 flex items-center justify-between transition">
            <div>
              <div class="flex items-center space-x-2">
                <span class="text-xs font-bold text-white">{r.guest_name}</span>
                <span class="text-[10px] font-mono bg-slate-800 text-slate-400 px-1.5 py-0.5 rounded">{r.res_code}</span>
                {#if r.room_number}
                  <span class="text-[10px] font-mono bg-amber-500/10 text-amber-300 border border-amber-500/20 px-1.5 py-0.5 rounded">Room {r.room_number}</span>
                {/if}
              </div>
              <p class="text-[11px] text-slate-400 mt-0.5">
                {r.room_type} Room &bull; {r.nights} Night(s) &bull; Check-in: {r.check_in_date}
              </p>
            </div>
            <div class="flex items-center space-x-1.5">
              <a
                href="/frontdesk/checkin?room={r.room_number}"
                class="text-[11px] font-semibold text-slate-300 hover:text-white bg-slate-800 hover:bg-slate-700 px-2.5 py-1.5 rounded-lg border border-slate-700 transition"
              >
                Review
              </a>
              <button
                type="button"
                onclick={() => quickClearArrival(r)}
                class="text-[11px] font-bold text-slate-950 bg-amber-500 hover:bg-amber-400 px-3 py-1.5 rounded-lg shadow-md shadow-amber-500/20 transition active:scale-95 flex items-center space-x-1"
              >
                <i class="ph ph-check-bold text-xs"></i>
                <span>Confirm Check-In</span>
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Live Occupied Stays Quick Monitor -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-lg space-y-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-2">
        <i class="ph ph-bed text-amber-400 text-lg"></i>
        <h3 class="font-bold text-sm text-white">Active Guest Stays</h3>
      </div>
      <span class="text-xs text-slate-400">{activeStays.length} active</span>
    </div>

    <!-- Quick search input -->
    <div class="relative">
      <input
        type="text"
        placeholder="Filter by guest name or room..."
        bind:value={searchQuery}
        class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-amber-500 pl-8"
      />
      <i class="ph ph-magnifying-glass absolute left-2.5 top-2.5 text-slate-500 text-xs"></i>
    </div>

    <!-- Active list -->
    <div class="space-y-2.5 mt-2 max-h-72 overflow-y-auto pr-1">
      {#if activeStays.length === 0}
        <div class="text-center py-6 text-slate-500 text-xs">
          <i class="ph ph-house text-2xl mb-1 text-slate-600 block"></i>
          No occupied rooms found.
        </div>
      {:else}
        {#each activeStays as stay (stay.room.room_number)}
          <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-3 flex flex-col space-y-2.5">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-2">
                <span class="w-8 h-8 rounded-lg bg-amber-500/20 text-amber-400 font-['Oxanium',sans-serif] font-bold text-sm flex items-center justify-center">
                  {stay.room.room_number}
                </span>
                <div>
                  <h4 class="text-xs font-bold text-white">
                    {stay.folio ? stay.folio.guest_name : 'Guest In-House'}
                  </h4>
                  <p class="text-[11px] text-slate-400">{stay.room.type} Room</p>
                </div>
              </div>
              <div class="text-right">
                {#if stay.folio}
                  <span class="text-xs font-bold {stay.folio.balance > 0 ? 'text-amber-400' : 'text-emerald-400'}">
                    {formatMWK(stay.folio.balance)}
                  </span>
                  <p class="text-[10px] text-slate-500">
                    {stay.folio.balance > 0 ? 'Balance due' : 'Paid in full'}
                  </p>
                {:else}
                  <span class="text-[11px] text-slate-400">Loading folio...</span>
                {/if}
              </div>
            </div>

            <!-- Quick Action Buttons for this Stay -->
            <div class="grid grid-cols-4 gap-1.5 pt-1 border-t border-slate-800/60 text-center">
              <a
                href="/frontdesk/folio?room={stay.room.room_number}"
                class="py-1 px-1 rounded-lg bg-slate-800 text-[10px] font-semibold text-slate-300 hover:text-amber-400 transition"
              >
                Folio
              </a>
              <a
                href="/frontdesk/payments?room={stay.room.room_number}"
                class="py-1 px-1 rounded-lg bg-slate-800 text-[10px] font-semibold text-emerald-400 hover:bg-slate-700 transition"
              >
                Pay
              </a>
              <a
                href="/frontdesk/extension?room={stay.room.room_number}"
                class="py-1 px-1 rounded-lg bg-slate-800 text-[10px] font-semibold text-sky-400 hover:bg-slate-700 transition"
              >
                Extend
              </a>
              <a
                href="/frontdesk/checkout?room={stay.room.room_number}"
                class="py-1 px-1 rounded-lg bg-slate-800 text-[10px] font-semibold text-rose-400 hover:bg-slate-700 transition"
              >
                C/Out
              </a>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>
