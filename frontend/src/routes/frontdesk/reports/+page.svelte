<!-- frontend/src/routes/frontdesk/reports/+page.svelte -->
<script lang="ts">
  import {
    roomStore,
    availableRoomsCount,
    occupiedRoomsCount,
    cleaningRoomsCount,
    outOfServiceRoomsCount
  } from '#lib/stores/roomStore';
  import { authStore } from '#lib/stores/authStore';
  import { onMount } from 'svelte';

  const rooms = roomStore.rooms;
  const folios = roomStore.folios;

  let handoverNotes = $state('');
  let notesSaved = $state(false);

  // Derived financial & operational stats
  let totalRooms = $derived($rooms.length);
  let occupancyPercent = $derived(
    totalRooms > 0 ? Math.round(($occupiedRoomsCount / totalRooms) * 100) : 0
  );

  let totalRevenueCollected = $derived(
    Object.values($folios).reduce((sum, f) => sum + (f.total_paid || 0), 0)
  );

  let totalOutstanding = $derived(
    Object.values($folios).reduce((sum, f) => sum + (f.balance || 0), 0)
  );

  let totalGrossCharges = $derived(
    Object.values($folios).reduce((sum, f) => sum + (f.total_charges || 0), 0)
  );

  // Payment method breakdowns from all items
  let paymentsBreakdown = $derived.by(() => {
    const methods: Record<string, number> = {
      Cash: 0,
      Visa: 0,
      Airtel: 0,
      Mpamba: 0
    };

    for (const f of Object.values($folios)) {
      for (const item of f.items) {
        if (item.category === 'Payment' && item.amount < 0) {
          const paid = Math.abs(item.amount);
          const desc = item.description.toLowerCase();
          if (desc.includes('cash') || desc.includes('deposit')) {
            methods.Cash += paid;
          } else if (desc.includes('visa') || desc.includes('pos')) {
            methods.Visa += paid;
          } else if (desc.includes('airtel')) {
            methods.Airtel += paid;
          } else if (desc.includes('mpamba')) {
            methods.Mpamba += paid;
          } else {
            methods.Cash += paid;
          }
        }
      }
    }
    return methods;
  });

  onMount(() => {
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('la_luna_handover_notes');
      if (saved) handoverNotes = saved;
    }
  });

  function saveNotes() {
    if (typeof window !== 'undefined') {
      localStorage.setItem('la_luna_handover_notes', handoverNotes);
      notesSaved = true;
      setTimeout(() => (notesSaved = false), 2000);
    }
  }

  function handlePrintSummary() {
    if (typeof window !== 'undefined') {
      window.print();
    }
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">RECEPTION REPORTS</h1>
      <p class="text-xs text-slate-400">Daily Reception & Shift Handover Summary</p>
    </div>
    <button
      onclick={handlePrintSummary}
      class="px-2.5 py-1.5 rounded-xl bg-slate-800 text-slate-300 hover:text-white border border-slate-700 text-xs flex items-center space-x-1"
    >
      <i class="ph ph-printer text-base"></i>
      <span>Print</span>
    </button>
  </div>

  <!-- Executive Occupancy Gauge Card -->
  <div class="bg-gradient-to-r from-slate-900 to-slate-950 border border-slate-800 rounded-3xl p-5 shadow-xl space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <span class="text-[10px] font-bold uppercase tracking-wider text-amber-400">Occupancy Metric</span>
        <h2 class="text-3xl font-bold font-['Oxanium',sans-serif] text-white">{occupancyPercent}%</h2>
        <p class="text-xs text-slate-400 mt-0.5">
          {$occupiedRoomsCount} of {totalRooms} rooms occupied
        </p>
      </div>

      <!-- Radial or Progress Badge -->
      <div class="w-16 h-16 rounded-2xl bg-amber-500/10 border border-amber-500/30 flex flex-col items-center justify-center text-amber-400 font-bold">
        <i class="ph ph-chart-pie-slice text-2xl"></i>
        <span class="text-[10px] font-mono mt-0.5">{$occupiedRoomsCount}/{totalRooms}</span>
      </div>
    </div>

    <!-- Inventory Distribution Bar -->
    <div class="space-y-1.5">
      <div class="h-3 w-full bg-slate-800 rounded-full overflow-hidden flex">
        <div style="width: {($occupiedRoomsCount / totalRooms) * 100}%" class="bg-amber-500" title="Occupied"></div>
        <div style="width: {($availableRoomsCount / totalRooms) * 100}%" class="bg-emerald-500" title="Available"></div>
        <div style="width: {($cleaningRoomsCount / totalRooms) * 100}%" class="bg-sky-500" title="Cleaning"></div>
        <div style="width: {($outOfServiceRoomsCount / totalRooms) * 100}%" class="bg-rose-500" title="Out of Service"></div>
      </div>
      <div class="flex justify-between text-[10px] text-slate-400 font-mono">
        <span class="text-amber-400">{$occupiedRoomsCount} Occ</span>
        <span class="text-emerald-400">{$availableRoomsCount} Avail</span>
        <span class="text-sky-400">{$cleaningRoomsCount} Clean</span>
        <span class="text-rose-400">{$outOfServiceRoomsCount} OOO</span>
      </div>
    </div>
  </div>

  <!-- Financial Collection Overview -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-xl space-y-3">
    <h3 class="text-xs font-bold uppercase tracking-wider text-slate-300 flex items-center justify-between">
      <span>Financial Collections</span>
      <span class="text-[10px] font-mono text-emerald-400 font-normal">Active In-House Folios</span>
    </h3>

    <div class="grid grid-cols-2 gap-2.5">
      <div class="bg-slate-950 p-3 rounded-xl border border-slate-800/80">
        <span class="text-[10px] text-slate-400 uppercase block font-semibold">Total Revenue Collected</span>
        <span class="text-base font-bold font-mono text-emerald-400 mt-1 block">
          MWK {totalRevenueCollected.toLocaleString()}
        </span>
      </div>

      <div class="bg-slate-950 p-3 rounded-xl border border-slate-800/80">
        <span class="text-[10px] text-slate-400 uppercase block font-semibold">Outstanding Balance Due</span>
        <span class="text-base font-bold font-mono text-amber-400 mt-1 block">
          MWK {totalOutstanding.toLocaleString()}
        </span>
      </div>
    </div>

    <!-- Payment Channel Breakdown -->
    <div class="pt-2 border-t border-slate-800/80 space-y-2">
      <span class="text-[11px] font-semibold text-slate-400 block">Collections by Channel:</span>
      <div class="grid grid-cols-2 gap-2 text-xs">
        <div class="flex justify-between p-2 rounded-lg bg-slate-950/60 border border-slate-800/60">
          <span class="text-slate-400">Cash:</span>
          <span class="font-mono font-bold text-white">MWK {paymentsBreakdown.Cash.toLocaleString()}</span>
        </div>
        <div class="flex justify-between p-2 rounded-lg bg-slate-950/60 border border-slate-800/60">
          <span class="text-slate-400">Visa POS:</span>
          <span class="font-mono font-bold text-white">MWK {paymentsBreakdown.Visa.toLocaleString()}</span>
        </div>
        <div class="flex justify-between p-2 rounded-lg bg-slate-950/60 border border-slate-800/60">
          <span class="text-slate-400">Airtel Money:</span>
          <span class="font-mono font-bold text-white">MWK {paymentsBreakdown.Airtel.toLocaleString()}</span>
        </div>
        <div class="flex justify-between p-2 rounded-lg bg-slate-950/60 border border-slate-800/60">
          <span class="text-slate-400">TNM Mpamba:</span>
          <span class="font-mono font-bold text-white">MWK {paymentsBreakdown.Mpamba.toLocaleString()}</span>
        </div>
      </div>
    </div>
  </div>

  <!-- Shift Handover Notes & Log -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-xl space-y-3">
    <div class="flex items-center justify-between">
      <h3 class="text-xs font-bold uppercase tracking-wider text-slate-300">Shift Handover Log</h3>
      <span class="text-[11px] text-amber-400 font-mono">{$authStore.shift}</span>
    </div>

    <div class="text-xs text-slate-400">
      Logged Officer: <strong class="text-white">{$authStore.fullName}</strong>
    </div>

    <textarea
      rows="3"
      placeholder="Type shift notes for incoming reception team (e.g., Room 202 requested late checkout at 12:00, Room 102 deposit confirmed by manager)..."
      bind:value={handoverNotes}
      class="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500"
    ></textarea>

    <div class="flex items-center justify-between pt-1">
      <span class="text-[11px] {notesSaved ? 'text-emerald-400 font-semibold' : 'text-slate-500'}">
        {notesSaved ? 'Handover notes saved!' : 'Auto-cached locally'}
      </span>
      <button
        type="button"
        onclick={saveNotes}
        class="px-4 py-2 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-md active:scale-95 transition"
      >
        Save Handover Notes
      </button>
    </div>
  </div>
</div>
