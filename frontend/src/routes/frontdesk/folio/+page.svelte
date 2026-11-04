<!-- frontend/src/routes/frontdesk/folio/+page.svelte -->
<script lang="ts">
  import { page } from '$app/state';
  import { roomStore } from '#lib/stores/roomStore';
  import type { Folio, FolioItem, Room } from '#lib/types/room';
  import { onMount } from 'svelte';

  const rooms = roomStore.rooms;
  const folios = roomStore.folios;

  let selectedRoomNumber = $state('102');

  // Rooms that have active folios
  let availableFolioRooms = $derived(
    Object.keys($folios).sort()
  );

  let currentFolio = $derived<Folio | null>(
    selectedRoomNumber ? $folios[selectedRoomNumber] || null : null
  );

  let currentRoom = $derived(
    $rooms.find((r: Room) => r.room_number === selectedRoomNumber) || null
  );

  onMount(() => {
    const queryRoom = page.url.searchParams.get('room');
    if (queryRoom && $folios[queryRoom]) {
      selectedRoomNumber = queryRoom;
    } else if (availableFolioRooms.length > 0) {
      selectedRoomNumber = availableFolioRooms[0];
    }
  });

  function getCategoryColor(category: string): string {
    switch (category) {
      case 'Accommodation':
        return 'bg-amber-500/10 text-amber-400 border-amber-500/30';
      case 'Payment':
        return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30';
      case 'Restaurant':
        return 'bg-orange-500/10 text-orange-400 border-orange-500/30';
      case 'Bar':
        return 'bg-purple-500/10 text-purple-400 border-purple-500/30';
      default:
        return 'bg-slate-700/40 text-slate-300 border-slate-600/30';
    }
  }

  function handlePrintReceipt() {
    if (typeof window !== 'undefined') {
      window.print();
    }
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">GUEST BILLING FOLIO</h1>
      <p class="text-xs text-slate-400">Itemized Charges & Account Settlement</p>
    </div>
    <button
      onclick={handlePrintReceipt}
      class="px-2.5 py-1.5 rounded-xl bg-slate-800 text-slate-300 hover:text-white border border-slate-700 text-xs flex items-center space-x-1"
      title="Print Folio Statement"
    >
      <i class="ph ph-printer text-base"></i>
      <span>Print</span>
    </button>
  </div>

  <!-- Room Folio Selector -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5 space-y-2">
    <label for="room-folio-select" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400">
      Active Room Account
    </label>
    {#if availableFolioRooms.length === 0}
      <p class="text-xs text-slate-500">No active guest folios currently open.</p>
    {:else}
      <select
        id="room-folio-select"
        bind:value={selectedRoomNumber}
        class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2.5 text-xs text-white focus:outline-none focus:border-amber-500 font-medium"
      >
        {#each availableFolioRooms as rNo}
          {@const f = $folios[rNo]}
          <option value={rNo}>
            Room {rNo} - {f.guest_name} (Bal: MWK {f.balance.toLocaleString()})
          </option>
        {/each}
      </select>
    {/if}
  </div>

  {#if currentFolio}
    <!-- Folio Header & Guest Identity -->
    <div class="bg-gradient-to-br from-slate-900 to-slate-950 border border-slate-800 rounded-2xl p-4 shadow-xl space-y-3">
      <div class="flex items-start justify-between">
        <div>
          <span class="text-[10px] font-mono uppercase bg-amber-500/10 text-amber-400 px-2 py-0.5 rounded border border-amber-500/20">
            {currentFolio.id}
          </span>
          <h2 class="text-base font-bold text-white mt-1.5">{currentFolio.guest_name}</h2>
          {#if currentFolio.guest_phone || currentFolio.guest_id_doc}
            <p class="text-[11px] text-slate-400 mt-0.5 font-mono">
              {currentFolio.guest_phone || ''} {currentFolio.guest_id_doc ? `&bull; ${currentFolio.guest_id_doc}` : ''}
            </p>
          {/if}
        </div>
        <div class="text-right">
          <span class="font-['Oxanium',sans-serif] text-2xl font-bold text-amber-400">
            {currentFolio.room_number}
          </span>
          <p class="text-[11px] text-slate-400">{currentRoom?.type || 'Room'}</p>
        </div>
      </div>

      <!-- Stay Period -->
      <div class="grid grid-cols-2 gap-2 pt-2 border-t border-slate-800/80 text-xs">
        <div>
          <span class="text-[10px] text-slate-500 uppercase block font-semibold">Check-In</span>
          <span class="font-mono text-slate-300">{currentFolio.check_in_date}</span>
        </div>
        <div class="text-right">
          <span class="text-[10px] text-slate-500 uppercase block font-semibold">Check-Out</span>
          <span class="font-mono text-slate-300">{currentFolio.check_out_date} (10:00 AM)</span>
        </div>
      </div>

      <!-- Financial Totals Triple Banner -->
      <div class="grid grid-cols-3 gap-2 pt-2 border-t border-slate-800/80 text-center">
        <div class="bg-slate-950/80 rounded-xl p-2 border border-slate-800">
          <span class="text-[9px] uppercase tracking-wider text-slate-400 block">Total Charges</span>
          <span class="text-xs font-bold font-mono text-white">MWK {currentFolio.total_charges.toLocaleString()}</span>
        </div>
        <div class="bg-slate-950/80 rounded-xl p-2 border border-slate-800">
          <span class="text-[9px] uppercase tracking-wider text-slate-400 block">Total Paid</span>
          <span class="text-xs font-bold font-mono text-emerald-400">MWK {currentFolio.total_paid.toLocaleString()}</span>
        </div>
        <div class="bg-slate-950/80 rounded-xl p-2 border {currentFolio.balance > 0 ? 'border-amber-500/40 bg-amber-500/5' : 'border-emerald-500/40 bg-emerald-500/5'}">
          <span class="text-[9px] uppercase tracking-wider text-slate-400 block">Balance Due</span>
          <span class="text-xs font-bold font-mono {currentFolio.balance > 0 ? 'text-amber-400' : 'text-emerald-400'}">
            MWK {currentFolio.balance.toLocaleString()}
          </span>
        </div>
      </div>
    </div>

    <!-- Itemized Ledger Line Items -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-xl space-y-3">
      <div class="flex items-center justify-between pb-1 border-b border-slate-800">
        <h3 class="text-xs font-bold uppercase tracking-wider text-slate-300">Itemized Charges & Receipts</h3>
        <span class="text-[11px] text-slate-500 font-mono">{currentFolio.items.length} items</span>
      </div>

      <div class="space-y-2.5 max-h-72 overflow-y-auto pr-1">
        {#each currentFolio.items as item (item.id)}
          <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-2.5 flex items-center justify-between text-xs">
            <div class="flex-1 pr-2">
              <div class="flex items-center space-x-1.5 mb-1">
                <span class="text-[9px] font-semibold uppercase px-1.5 py-0.2 rounded border {getCategoryColor(item.category)}">
                  {item.category}
                </span>
                <span class="text-[10px] text-slate-500 font-mono">
                  {item.timestamp ? item.timestamp.split('T')[0] : ''}
                </span>
              </div>
              <p class="font-medium text-slate-200">{item.description}</p>
            </div>
            <div class="text-right">
              <span class="font-mono font-bold {item.amount < 0 ? 'text-emerald-400' : 'text-white'}">
                {item.amount < 0 ? '-' : ''}MWK {Math.abs(item.amount).toLocaleString()}
              </span>
            </div>
          </div>
        {/each}
      </div>
    </div>

    <!-- Quick Action Bar -->
    <div class="grid grid-cols-3 gap-2">
      <a
        href="/frontdesk/payments?room={currentFolio.room_number}"
        class="py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs text-center shadow-md active:scale-95 transition flex items-center justify-center space-x-1"
      >
        <i class="ph ph-plus-circle text-sm"></i>
        <span>Add Pay</span>
      </a>

      <a
        href="/frontdesk/extension?room={currentFolio.room_number}"
        class="py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold text-xs text-center border border-slate-700 active:scale-95 transition flex items-center justify-center space-x-1"
      >
        <i class="ph ph-calendar-plus text-sm"></i>
        <span>Extend</span>
      </a>

      <a
        href="/frontdesk/checkout?room={currentFolio.room_number}"
        class="py-2.5 rounded-xl bg-rose-600 hover:bg-rose-500 text-white font-bold text-xs text-center shadow-md active:scale-95 transition flex items-center justify-center space-x-1"
      >
        <i class="ph ph-sign-out text-sm"></i>
        <span>Checkout</span>
      </a>
    </div>
  {/if}
</div>
