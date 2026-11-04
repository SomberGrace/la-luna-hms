<!-- frontend/src/routes/frontdesk/housekeeping/+page.svelte -->
<script lang="ts">
  import { roomStore, cleaningRoomsCount, outOfServiceRoomsCount, availableRoomsCount } from '#lib/stores/roomStore';
  import type { Room, RoomStatus } from '#lib/types/room';

  const rooms = roomStore.rooms;

  let filterTab = $state<'turnover' | 'all'>('turnover');
  let selectedFloor = $state<string>('All');
  let notification = $state('');

  // Turnover rooms are Cleaning or OutOfService
  let displayedRooms = $derived(
    $rooms.filter((r: Room) => {
      if (filterTab === 'turnover') {
        if (r.status !== 'Cleaning' && r.status !== 'OutOfService') return false;
      }
      if (selectedFloor !== 'All' && String(r.floor || 1) !== selectedFloor) return false;
      return true;
    })
  );

  async function handleMarkClean(room: Room) {
    const res = await roomStore.updateRoomStatus(room.room_number, 'Available', undefined, 'Normal');
    showNotification(`Room ${room.room_number} is now Cleaned & Available for check-in!`);
  }

  async function handleSetPriority(room: Room, priority: 'Normal' | 'High' | 'Urgent') {
    await roomStore.updateRoomStatus(room.room_number, room.status, room.housekeeper, priority);
    showNotification(`Room ${room.room_number} priority set to ${priority}`);
  }

  async function handleSendToMaintenance(room: Room) {
    await roomStore.updateRoomStatus(room.room_number, 'OutOfService', 'Maintenance Team', 'High');
    showNotification(`Room ${room.room_number} marked Out of Service for Maintenance.`);
  }

  function showNotification(msg: string) {
    notification = msg;
    setTimeout(() => {
      notification = '';
    }, 2800);
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">HOUSEKEEPING</h1>
      <p class="text-xs text-slate-400">Room Turnover & Cleaning Status Tracker</p>
    </div>
    <div class="flex items-center space-x-1.5 text-xs font-mono">
      <span class="bg-sky-500/10 text-sky-400 border border-sky-500/20 px-2 py-0.5 rounded-full font-bold">
        {$cleaningRoomsCount} Cleaning
      </span>
      <span class="bg-rose-500/10 text-rose-400 border border-rose-500/20 px-2 py-0.5 rounded-full font-bold">
        {$outOfServiceRoomsCount} OOO
      </span>
    </div>
  </div>

  {#if notification}
    <div class="p-3 rounded-2xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-xs flex items-center justify-between animate-fade-in">
      <span>{notification}</span>
      <i class="ph ph-check-circle text-base text-emerald-400"></i>
    </div>
  {/if}

  <!-- View Mode Switcher -->
  <div class="grid grid-cols-2 gap-2 bg-slate-900 p-1 rounded-2xl border border-slate-800 text-xs">
    <button
      onclick={() => (filterTab = 'turnover')}
      class="py-2 rounded-xl font-bold transition {filterTab === 'turnover'
        ? 'bg-amber-500 text-slate-950 shadow-md shadow-amber-500/20'
        : 'text-slate-400 hover:text-white'}"
    >
      Turnover Queue ({$cleaningRoomsCount + $outOfServiceRoomsCount})
    </button>
    <button
      onclick={() => (filterTab = 'all')}
      class="py-2 rounded-xl font-bold transition {filterTab === 'all'
        ? 'bg-amber-500 text-slate-950 shadow-md shadow-amber-500/20'
        : 'text-slate-400 hover:text-white'}"
    >
      All Rooms ({$rooms.length})
    </button>
  </div>

  <!-- Floor Selector Filter -->
  <div class="flex items-center space-x-2 text-xs">
    <span class="text-slate-400 font-medium">Floor:</span>
    {#each ['All', '1', '2', '3'] as fl}
      <button
        onclick={() => (selectedFloor = fl)}
        class="px-2.5 py-1 rounded-lg border transition-colors {selectedFloor === fl
          ? 'bg-slate-800 border-amber-500 text-amber-400 font-bold'
          : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200'}"
      >
        {fl === 'All' ? 'All Floors' : `Floor ${fl}`}
      </button>
    {/each}
  </div>

  <!-- Housekeeping Cards List -->
  <div class="space-y-3 pt-1">
    {#if displayedRooms.length === 0}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-8 text-center text-slate-500 text-xs">
        <i class="ph ph-sparkle text-3xl mb-2 text-emerald-400 block"></i>
        <p class="font-semibold text-white">All Rooms Cleaned & Inspected!</p>
        <p class="mt-1">No rooms currently awaiting housekeeping or maintenance.</p>
      </div>
    {:else}
      {#each displayedRooms as room (room.room_number)}
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-lg space-y-3">
          <div class="flex items-start justify-between">
            <div class="flex items-center space-x-3">
              <div class="w-12 h-12 rounded-xl bg-slate-950 border border-slate-800 flex flex-col items-center justify-center font-['Oxanium',sans-serif]">
                <span class="text-base font-bold text-white">{room.room_number}</span>
                <span class="text-[9px] uppercase tracking-wider text-slate-400">Fl {room.floor || 1}</span>
              </div>
              <div>
                <div class="flex items-center space-x-2">
                  <h3 class="font-bold text-sm text-white">{room.type}</h3>
                  {#if room.priority === 'Urgent'}
                    <span class="text-[9px] font-bold uppercase bg-rose-500/20 text-rose-400 px-1.5 py-0.5 rounded border border-rose-500/30 animate-pulse">
                      Urgent
                    </span>
                  {:else if room.priority === 'High'}
                    <span class="text-[9px] font-bold uppercase bg-amber-500/20 text-amber-400 px-1.5 py-0.5 rounded border border-amber-500/30">
                      High Priority
                    </span>
                  {/if}
                </div>
                <p class="text-xs text-slate-400 mt-0.5">
                  Assigned: <strong class="text-slate-200">{room.housekeeper || 'Unassigned'}</strong>
                </p>
              </div>
            </div>

            <!-- Operational Status Chip -->
            <span class="px-2.5 py-1 rounded-full text-[10px] font-bold border {room.status === 'Cleaning' ? 'bg-sky-500/10 text-sky-400 border-sky-500/30' : room.status === 'Available' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30' : room.status === 'Occupied' ? 'bg-amber-500/10 text-amber-400 border-amber-500/30' : 'bg-rose-500/10 text-rose-400 border-rose-500/30'}">
              {room.status}
            </span>
          </div>

          <!-- Actions Bar -->
          <div class="pt-2 border-t border-slate-800/80 flex items-center justify-between">
            {#if room.status === 'Cleaning'}
              <div class="flex items-center space-x-1.5">
                <button
                  type="button"
                  onclick={() => handleSetPriority(room, 'Urgent')}
                  class="px-2 py-1 rounded-lg bg-slate-800 text-[10px] text-amber-400 hover:bg-slate-700 font-semibold"
                >
                  Expedite
                </button>
                <button
                  type="button"
                  onclick={() => handleSendToMaintenance(room)}
                  class="px-2 py-1 rounded-lg bg-slate-800 text-[10px] text-rose-400 hover:bg-slate-700 font-semibold"
                >
                  Report OOO
                </button>
              </div>

              <button
                type="button"
                onclick={() => handleMarkClean(room)}
                class="px-3 py-1.5 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-slate-950 font-bold text-xs flex items-center space-x-1 shadow-md active:scale-95 transition"
              >
                <i class="ph ph-check-fat text-xs"></i>
                <span>Mark Cleaned</span>
              </button>
            {:else if room.status === 'OutOfService'}
              <span class="text-[11px] text-rose-400 flex items-center">
                <i class="ph ph-wrench mr-1"></i> Maintenance Active
              </span>
              <button
                type="button"
                onclick={() => roomStore.updateRoomStatus(room.room_number, 'Cleaning', 'Grace Phiri', 'Normal')}
                class="px-3 py-1.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs active:scale-95 transition"
              >
                Send to Clean
              </button>
            {:else}
              <span class="text-[11px] text-slate-500">In Service</span>
              <button
                type="button"
                onclick={() => roomStore.updateRoomStatus(room.room_number, 'Cleaning', 'Staff Assigned', 'Normal')}
                class="px-2.5 py-1 rounded-lg bg-slate-800 text-slate-300 hover:text-white text-xs"
              >
                Request Touch-Up
              </button>
            {/if}
          </div>
        </div>
      {/each}
    {/if}
  </div>
</div>
