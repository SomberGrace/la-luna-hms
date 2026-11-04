<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  let todayStr = $state(new Date().toISOString().slice(0, 10));
  let formattedDate = $state('');

  // Executive figures matching executive_summaries table
  let roomCash = $state(90000);
  let roomVisa = $state(120000);
  let roomCheque = $state(45000);
  let roomCredit = $state(0);
  let restaurantSales = $state(47200);
  let barSales = $state(38400);
  let pettyCash = $state(15000);

  // Derived financials
  let totalRoomSales = $derived(roomCash + roomVisa + roomCheque + roomCredit);
  let grossRevenue = $derived(totalRoomSales + restaurantSales + barSales);
  let netBankingDeposit = $derived(grossRevenue - pettyCash);

  // Room Inventory metrics
  let totalRooms = $state(8);
  let occupiedRooms = $state(6);
  let occupancyRate = $derived(Math.round((occupiedRooms / totalRooms) * 100));

  // Department Sign-off Statuses
  let pendingAuditsCount = $state(3);

  onMount(() => {
    formattedDate = new Date().toLocaleDateString('en-GB', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });

    apiFetch(`/api/v1/manager/summary?date=${todayStr}`)
      .then(res => res.json())
      .then(data => {
        if (data.room_cash !== undefined) roomCash = data.room_cash;
        if (data.restaurant_sales !== undefined) restaurantSales = data.restaurant_sales;
        if (data.bar_sales !== undefined) barSales = data.bar_sales;
      })
      .catch(() => {
        // Fallback demo data active
      });
  });

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }
</script>

<div class="space-y-4 pb-4">
  <!-- Top Greeting & Date -->
  <div class="pt-1 flex items-start justify-between">
    <div>
      <p class="text-xs text-slate-400">{formattedDate}</p>
      <h1 class="text-2xl font-bold text-white tracking-tight">Executive Desk</h1>
      <p class="text-xs text-indigo-400 font-semibold">General Manager Station</p>
    </div>
    <a
      href="/manager/audits"
      class="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl bg-amber-500/20 border border-amber-500/30 text-amber-300 text-xs font-semibold active:scale-95 transition"
    >
      <i class="ph ph-bell-ringing"></i>
      <span>{pendingAuditsCount} Audits</span>
    </a>
  </div>

  <!-- Grand Total Revenue Banner -->
  <div class="bg-gradient-to-br from-indigo-950 via-slate-900 to-slate-900 border border-indigo-700/50 rounded-3xl p-5 shadow-2xl relative overflow-hidden">
    <div class="absolute -right-6 -bottom-6 w-32 h-32 rounded-full bg-indigo-500/10 blur-xl pointer-events-none"></div>

    <div class="flex items-center justify-between mb-2">
      <span class="text-xs text-indigo-300 uppercase font-bold tracking-wider">Gross Hotel Revenue</span>
      <span class="text-[11px] px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-400 font-bold border border-emerald-500/30">
        Today's Close
      </span>
    </div>

    <p class="text-3xl font-extrabold text-white font-['Oxanium',sans-serif] tracking-tight">
      {fmtMwk(grossRevenue)}
    </p>

    <!-- Sub Breakdown -->
    <div class="grid grid-cols-3 gap-2 mt-4 pt-3 border-t border-slate-800 text-xs">
      <div>
        <span class="text-[10px] text-slate-400 block">Rooms ({Math.round((totalRoomSales / grossRevenue) * 100)}%)</span>
        <span class="font-bold text-white font-['Oxanium',sans-serif]">{fmtMwk(totalRoomSales)}</span>
      </div>
      <div>
        <span class="text-[10px] text-slate-400 block">Dining ({Math.round((restaurantSales / grossRevenue) * 100)}%)</span>
        <span class="font-bold text-orange-400 font-['Oxanium',sans-serif]">{fmtMwk(restaurantSales)}</span>
      </div>
      <div>
        <span class="text-[10px] text-slate-400 block">Bar ({Math.round((barSales / grossRevenue) * 100)}%)</span>
        <span class="font-bold text-purple-400 font-['Oxanium',sans-serif]">{fmtMwk(barSales)}</span>
      </div>
    </div>
  </div>

  <!-- Key Executive Metrics -->
  <div class="grid grid-cols-2 gap-3">
    <!-- Occupancy Rate -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-2">
      <div class="flex items-center justify-between">
        <span class="text-xs text-slate-400 uppercase font-semibold">Occupancy Rate</span>
        <i class="ph ph-bed text-indigo-400 text-base"></i>
      </div>
      <div class="flex items-baseline space-x-2">
        <span class="text-2xl font-bold text-white font-['Oxanium',sans-serif]">{occupancyRate}%</span>
        <span class="text-xs text-slate-400">({occupiedRooms}/{totalRooms} rooms)</span>
      </div>
      <!-- Progress Bar -->
      <div class="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
        <div class="bg-indigo-500 h-full rounded-full transition-all" style="width: {occupancyRate}%"></div>
      </div>
    </div>

    <!-- Banking Deposit -->
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-1">
      <div class="flex items-center justify-between">
        <span class="text-xs text-slate-400 uppercase font-semibold">Net Banking</span>
        <i class="ph ph-bank text-emerald-400 text-base"></i>
      </div>
      <p class="text-lg font-bold text-emerald-400 font-['Oxanium',sans-serif] mt-1">
        {fmtMwk(netBankingDeposit)}
      </p>
      <p class="text-[10px] text-slate-500">
        Less {fmtMwk(pettyCash)} petty cash
      </p>
    </div>
  </div>

  <!-- Room Payment Modes Split (Audit requirement) -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-3">
    <div class="flex items-center justify-between">
      <h2 class="text-xs font-bold text-white uppercase tracking-wider">Room Revenue By Settlement Method</h2>
      <span class="text-xs text-slate-400 font-mono">{fmtMwk(totalRoomSales)}</span>
    </div>

    <div class="grid grid-cols-2 gap-2 text-xs">
      <div class="bg-slate-800/80 rounded-xl p-2.5 flex items-center justify-between">
        <div class="flex items-center space-x-1.5">
          <i class="ph ph-money text-emerald-400 text-base"></i>
          <span class="text-slate-300">Cash</span>
        </div>
        <span class="font-bold text-white font-['Oxanium',sans-serif]">{fmtMwk(roomCash)}</span>
      </div>

      <div class="bg-slate-800/80 rounded-xl p-2.5 flex items-center justify-between">
        <div class="flex items-center space-x-1.5">
          <i class="ph ph-credit-card text-blue-400 text-base"></i>
          <span class="text-slate-300">Visa / POS</span>
        </div>
        <span class="font-bold text-white font-['Oxanium',sans-serif]">{fmtMwk(roomVisa)}</span>
      </div>

      <div class="bg-slate-800/80 rounded-xl p-2.5 flex items-center justify-between">
        <div class="flex items-center space-x-1.5">
          <i class="ph ph-note text-amber-400 text-base"></i>
          <span class="text-slate-300">Cheque</span>
        </div>
        <span class="font-bold text-white font-['Oxanium',sans-serif]">{fmtMwk(roomCheque)}</span>
      </div>

      <div class="bg-slate-800/80 rounded-xl p-2.5 flex items-center justify-between">
        <div class="flex items-center space-x-1.5">
          <i class="ph ph-user-check text-purple-400 text-base"></i>
          <span class="text-slate-300">Credit / Acc</span>
        </div>
        <span class="font-bold text-white font-['Oxanium',sans-serif]">{fmtMwk(roomCredit)}</span>
      </div>
    </div>
  </div>

  <!-- Quick Action Tiles -->
  <div class="grid grid-cols-2 gap-3">
    <a
      href="/manager/audits"
      class="flex flex-col items-center justify-center bg-indigo-600 hover:bg-indigo-500 active:scale-95 text-white rounded-2xl p-4 font-bold text-xs shadow-lg shadow-indigo-600/25 transition space-y-1.5"
    >
      <i class="ph ph-stamp text-2xl"></i>
      <span>Review Shift Audits</span>
    </a>

    <a
      href="/manager/reports"
      class="flex flex-col items-center justify-center bg-slate-800 hover:bg-slate-700 active:scale-95 text-slate-200 rounded-2xl p-4 font-semibold text-xs border border-slate-700/60 transition space-y-1.5"
    >
      <i class="ph ph-printer text-2xl"></i>
      <span>Consolidated PDF</span>
    </a>
  </div>
</div>
