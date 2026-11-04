<!-- frontend/src/routes/frontdesk/rooms/+page.svelte -->
<script lang="ts">
  import { page } from '$app/state';
  import { roomStore, availableRoomsCount, occupiedRoomsCount, cleaningRoomsCount, outOfServiceRoomsCount } from '#lib/stores/roomStore';
  import type { Room, RoomStatus, RoomType } from '#lib/types/room';

  const rooms = roomStore.rooms;
  const folios = roomStore.folios;

  // Initial filter from query params if present
  let statusFilter = $state<string>('All');
  let typeFilter = $state<string>('All');
  let searchRoom = $state('');

  // Status modal state
  let selectedRoomForStatus = $state<Room | null>(null);
  let newStatus = $state<RoomStatus>('Available');
  let housekeeperName = $state('');
  let statusFeedback = $state('');

  $effect(() => {
    const urlStatus = page.url.searchParams.get('status');
    if (urlStatus) {
      statusFilter = urlStatus;
    }
  });

  let filteredRooms = $derived(
    $rooms.filter((room: Room) => {
      // Status filter
      if (statusFilter !== 'All' && room.status !== statusFilter) return false;
      // Type filter
      if (typeFilter !== 'All' && room.type !== typeFilter) return false;
      // Search
      if (searchRoom.trim()) {
        const q = searchRoom.toLowerCase();
        const folio = $folios[room.room_number];
        const guest = folio?.guest_name?.toLowerCase() || '';
        return room.room_number.toLowerCase().includes(q) || guest.includes(q);
      }
      return true;
    })
  );

  function getStatusStyle(status: RoomStatus): { badge: string; border: string; bg: string; dot: string } {
    switch (status) {
      case 'Available':
        return {
          badge: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30',
          border: 'border-emerald-900/40 hover:border-emerald-600/60',
          bg: 'bg-slate-900/90',
          dot: 'bg-emerald-400'
        };
      case 'Occupied':
        return {
          badge: 'bg-amber-500/10 text-amber-400 border-amber-500/30',
          border: 'border-amber-900/40 hover:border-amber-600/60',
          bg: 'bg-slate-900/90',
          dot: 'bg-amber-400'
        };
      case 'Cleaning':
        return {
          badge: 'bg-sky-500/10 text-sky-400 border-sky-500/30',
          border: 'border-sky-900/40 hover:border-sky-600/60',
          bg: 'bg-slate-900/90',
          dot: 'bg-sky-400'
        };
      case 'OutOfService':
      default:
        return {
          badge: 'bg-rose-500/10 text-rose-400 border-rose-500/30',
          border: 'border-rose-900/40 hover:border-rose-600/60',
          bg: 'bg-slate-900/90',
          dot: 'bg-rose-400'
        };
    }
  }

  function openStatusModal(room: Room) {
    selectedRoomForStatus = room;
    newStatus = room.status;
    housekeeperName = room.housekeeper || '';
    statusFeedback = '';
  }

  async function handleQuickStatusChange(room: Room, targetStatus: RoomStatus) {
    const res = await roomStore.updateRoomStatus(room.room_number, targetStatus);
    statusFeedback = res.message;
    setTimeout(() => {
      statusFeedback = '';
    }, 2500);
  }

  async function handleModalStatusSave() {
    if (!selectedRoomForStatus) return;
    const res = await roomStore.updateRoomStatus(
      selectedRoomForStatus.room_number,
      newStatus,
      housekeeperName
    );
    statusFeedback = res.message;
    setTimeout(() => {
      selectedRoomForStatus = null;
      statusFeedback = '';
    }, 1200);
  }
</script>

<div class="space-y-4">
  <!-- Page Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">ROOM INVENTORY</h1>
      <p class="text-xs text-slate-400">Live Availability & Occupancy Grid</p>
    </div>
    <a
      href="/frontdesk/walkin"
      class="px-3 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs flex items-center space-x-1 shadow-md shadow-amber-500/20 active:scale-95 transition"
    >
      <i class="ph ph-plus-circle text-base"></i>
      <span>Assign</span>
    </a>
  </div>

  {#if statusFeedback}
    <div class="p-2.5 rounded-xl bg-amber-500/10 border border-amber-500/30 text-amber-300 text-xs flex items-center justify-between animate-fade-in">
      <span>{statusFeedback}</span>
      <i class="ph ph-check-circle text-base text-amber-400"></i>
    </div>
  {/if}

  <!-- Search Input -->
  <div class="relative">
    <input
      type="text"
      placeholder="Search room number or in-house guest..."
      bind:value={searchRoom}
      class="w-full bg-slate-900 border border-slate-800 rounded-xl px-3.5 py-2.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-amber-500 pl-9 shadow-inner"
    />
    <i class="ph ph-magnifying-glass absolute left-3 top-3 text-slate-500 text-sm"></i>
  </div>

  <!-- Operational State Filter Chips (Horizontal Scroll) -->
  <div class="flex items-center space-x-2 overflow-x-auto pb-1 scrollbar-none text-xs">
    <button
      onclick={() => (statusFilter = 'All')}
      class="px-3 py-1.5 rounded-xl font-semibold whitespace-nowrap transition-colors {statusFilter === 'All'
        ? 'bg-amber-500 text-slate-950 font-bold shadow-md shadow-amber-500/20'
        : 'bg-slate-900 text-slate-400 hover:text-white border border-slate-800'}"
    >
      All ({$rooms.length})
    </button>
    <button
      onclick={() => (statusFilter = 'Available')}
      class="px-3 py-1.5 rounded-xl font-semibold whitespace-nowrap transition-colors {statusFilter === 'Available'
        ? 'bg-emerald-500 text-slate-950 font-bold shadow-md shadow-emerald-500/20'
        : 'bg-slate-900 text-emerald-400 border border-slate-800'}"
    >
      Available ({$availableRoomsCount})
    </button>
    <button
      onclick={() => (statusFilter = 'Occupied')}
      class="px-3 py-1.5 rounded-xl font-semibold whitespace-nowrap transition-colors {statusFilter === 'Occupied'
        ? 'bg-amber-500 text-slate-950 font-bold shadow-md shadow-amber-500/20'
        : 'bg-slate-900 text-amber-400 border border-slate-800'}"
    >
      Occupied ({$occupiedRoomsCount})
    </button>
    <button
      onclick={() => (statusFilter = 'Cleaning')}
      class="px-3 py-1.5 rounded-xl font-semibold whitespace-nowrap transition-colors {statusFilter === 'Cleaning'
        ? 'bg-sky-500 text-slate-950 font-bold shadow-md shadow-sky-500/20'
        : 'bg-slate-900 text-sky-400 border border-slate-800'}"
    >
      Cleaning ({$cleaningRoomsCount})
    </button>
    <button
      onclick={() => (statusFilter = 'OutOfService')}
      class="px-3 py-1.5 rounded-xl font-semibold whitespace-nowrap transition-colors {statusFilter === 'OutOfService'
        ? 'bg-rose-500 text-slate-950 font-bold shadow-md shadow-rose-500/20'
        : 'bg-slate-900 text-rose-400 border border-slate-800'}"
    >
      Out of Service ({$outOfServiceRoomsCount})
    </button>
  </div>

  <!-- Room Type Filter Chips -->
  <div class="flex items-center space-x-2 text-[11px]">
    <span class="text-slate-500 font-medium">Type:</span>
    {#each ['All', 'Standard', 'Executive', 'Apartment'] as t}
      <button
        onclick={() => (typeFilter = t)}
        class="px-2 py-1 rounded-lg border transition-colors {typeFilter === t
          ? 'bg-slate-800 border-amber-500/60 text-amber-400 font-bold'
          : 'bg-slate-950 border-slate-800/80 text-slate-400 hover:text-slate-200'}"
      >
        {t}
      </button>
    {/each}
  </div>

  <!-- Rooms Cards Grid -->
  <div class="space-y-3 pt-1">
    {#if filteredRooms.length === 0}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-8 text-center text-slate-500 text-xs">
        <i class="ph ph-bed text-3xl mb-2 text-slate-600 block"></i>
        No rooms match the selected criteria.
      </div>
    {:else}
      {#each filteredRooms as room (room.room_number)}
        {@const style = getStatusStyle(room.status)}
        {@const folio = $folios[room.room_number]}
        <div class="{style.bg} border {style.border} rounded-2xl p-3.5 shadow-md transition-all flex flex-col space-y-3">
          <!-- Room Header & Status Badge -->
          <div class="flex items-start justify-between">
            <div class="flex items-center space-x-3">
              <div class="w-12 h-12 rounded-xl bg-slate-950 border border-slate-800 flex flex-col items-center justify-center font-['Oxanium',sans-serif] shadow-inner">
                <span class="text-base font-bold text-white tracking-wide">{room.room_number}</span>
                <span class="text-[9px] uppercase tracking-widest text-slate-400 font-mono">Fl {room.floor || 1}</span>
              </div>
              <div>
                <div class="flex items-center space-x-2">
                  <h3 class="font-bold text-sm text-white">{room.type}</h3>
                  <span class="text-[11px] font-mono text-amber-400 font-semibold">
                    MWK {room.price_per_night.toLocaleString()}<span class="text-slate-500 text-[10px]">/nt</span>
                  </span>
                </div>
                {#if room.status === 'Occupied' && folio}
                  <p class="text-xs text-slate-300 font-medium flex items-center mt-0.5">
                    <i class="ph ph-user mr-1 text-amber-400 text-xs"></i>
                    {folio.guest_name}
                  </p>
                {:else if room.status === 'Cleaning' && room.housekeeper}
                  <p class="text-xs text-sky-400 flex items-center mt-0.5">
                    <i class="ph ph-broom mr-1"></i> {room.housekeeper}
                  </p>
                {:else}
                  <p class="text-[11px] text-slate-500 mt-0.5">Ready for guest assignment</p>
                {/if}
              </div>
            </div>

            <!-- Status Pill -->
            <button
              onclick={() => openStatusModal(room)}
              class="px-2.5 py-1 rounded-full text-[11px] font-bold border {style.badge} flex items-center space-x-1.5 active:scale-95 transition"
              title="Click to change status"
            >
              <span class="w-2 h-2 rounded-full {style.dot}"></span>
              <span>{room.status}</span>
              <i class="ph ph-caret-down text-[10px] ml-0.5 opacity-60"></i>
            </button>
          </div>

          <!-- Bottom Context / Action Bar -->
          <div class="pt-2 border-t border-slate-800/60 flex items-center justify-between text-xs">
            {#if room.status === 'Available'}
              <span class="text-[11px] text-emerald-400 flex items-center">
                <i class="ph ph-check text-xs mr-1"></i> Clean & Inspected
              </span>
              <div class="flex items-center space-x-1.5">
                <button
                  onclick={() => handleQuickStatusChange(room, 'OutOfService')}
                  class="px-2 py-1 rounded-lg bg-slate-800 text-slate-400 hover:text-rose-400 text-[10px] transition"
                >
                  Mark OOO
                </button>
                <a
                  href="/frontdesk/walkin?room={room.room_number}"
                  class="px-2.5 py-1 rounded-lg bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-[11px] flex items-center space-x-1 transition shadow-sm"
                >
                  <i class="ph ph-user-plus text-xs"></i>
                  <span>Walk-In</span>
                </a>
              </div>
            {:else if room.status === 'Occupied'}
              <div class="text-[11px]">
                {#if folio}
                  <span class="{folio.balance > 0 ? 'text-amber-400' : 'text-emerald-400'} font-semibold font-mono">
                    Bal: MWK {folio.balance.toLocaleString()}
                  </span>
                {/if}
              </div>
              <div class="flex items-center space-x-1.5">
                <a
                  href="/frontdesk/folio?room={room.room_number}"
                  class="px-2 py-1 rounded-lg bg-slate-800 text-slate-300 hover:text-white text-[10px] font-semibold transition"
                >
                  Folio
                </a>
                <a
                  href="/frontdesk/checkout?room={room.room_number}"
                  class="px-2 py-1 rounded-lg bg-rose-500/20 border border-rose-500/30 text-rose-300 hover:bg-rose-500 hover:text-white text-[10px] font-semibold transition"
                >
                  Checkout
                </a>
              </div>
            {:else if room.status === 'Cleaning'}
              <span class="text-[11px] text-sky-400 flex items-center">
                <i class="ph ph-clock mr-1 text-xs"></i> Assigned to Housekeeping
              </span>
              <button
                onclick={() => handleQuickStatusChange(room, 'Available')}
                class="px-2.5 py-1 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-[11px] flex items-center space-x-1 transition shadow-sm"
              >
                <i class="ph ph-check text-xs"></i>
                <span>Mark Clean</span>
              </button>
            {:else}
              <span class="text-[11px] text-rose-400">Maintenance Required</span>
              <button
                onclick={() => handleQuickStatusChange(room, 'Cleaning')}
                class="px-2.5 py-1 rounded-lg bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-[11px] flex items-center space-x-1 transition shadow-sm"
              >
                <span>Send to Clean</span>
              </button>
            {/if}
          </div>
        </div>
      {/each}
    {/if}
  </div>

  <!-- Change Status Dialog Modal -->
  {#if selectedRoomForStatus}
    <div class="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-end sm:items-center justify-center p-3">
      <div class="bg-slate-900 border border-slate-800 rounded-3xl w-full max-w-sm p-5 space-y-4 shadow-2xl animate-slide-up">
        <div class="flex items-center justify-between pb-2 border-b border-slate-800">
          <div>
            <h3 class="font-bold text-base text-white">Room {selectedRoomForStatus.room_number} Status</h3>
            <p class="text-xs text-slate-400">{selectedRoomForStatus.type} - Floor {selectedRoomForStatus.floor || 1}</p>
          </div>
          <button
            onclick={() => (selectedRoomForStatus = null)}
            class="w-8 h-8 rounded-full bg-slate-800 text-slate-400 hover:text-white flex items-center justify-center"
          >
            <i class="ph ph-x text-base"></i>
          </button>
        </div>

        <div class="space-y-3 text-xs">
          <div>
            <label class="block text-slate-400 font-semibold mb-1">Target Operational Status</label>
            <div class="grid grid-cols-2 gap-2">
              {#each ['Available', 'Occupied', 'Cleaning', 'OutOfService'] as st}
                <button
                  type="button"
                  onclick={() => (newStatus = st as RoomStatus)}
                  class="py-2 px-3 rounded-xl border text-center font-semibold transition {newStatus === st
                    ? 'bg-amber-500 text-slate-950 font-bold border-amber-500 shadow-md shadow-amber-500/20'
                    : 'bg-slate-950 border-slate-800 text-slate-300'}"
                >
                  {st}
                </button>
              {/each}
            </div>
          </div>

          <div>
            <label class="block text-slate-400 font-semibold mb-1">Housekeeper / Attendant Name</label>
            <input
              type="text"
              placeholder="e.g. Grace Phiri"
              bind:value={housekeeperName}
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500"
            />
          </div>
        </div>

        <div class="flex items-center space-x-2 pt-2">
          <button
            type="button"
            onclick={() => (selectedRoomForStatus = null)}
            class="flex-1 py-2.5 rounded-xl bg-slate-800 text-slate-300 font-semibold text-xs hover:bg-slate-700"
          >
            Cancel
          </button>
          <button
            type="button"
            onclick={handleModalStatusSave}
            class="flex-1 py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-lg shadow-amber-500/20"
          >
            Apply Status
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>
