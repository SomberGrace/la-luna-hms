<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  type DocketSummary = {
    id: string;
    tableNo: string;
    roomNo: string;
    items: { name: string; qty: number }[];
    total: number;
    status: 'open' | 'room_charge' | 'settled';
    time: string;
  };

  let greeting = $state('Good morning');
  let shiftDate = $state('');
  let activeTables = $state(0);
  let todayDockets = $state(0);
  let todayRevenue = $state(0);
  let pendingDockets = $state(0);
  let recentDockets = $state<DocketSummary[]>([]);
  let loadError = $state('');

  onMount(() => {
    const h = new Date().getHours();
    greeting = h < 12 ? 'Good morning' : h < 17 ? 'Good afternoon' : 'Good evening';
    shiftDate = new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });

    void loadDashboard();
  });

  async function loadDashboard() {
    try {
      const [summaryResponse, docketsResponse] = await Promise.all([
        apiFetch('/api/v1/restaurant/summary'),
        apiFetch('/api/v1/restaurant/dockets')
      ]);
      if (!summaryResponse.ok || !docketsResponse.ok) {
        throw new Error('Unable to load restaurant dashboard data.');
      }
      const summary = await summaryResponse.json();
      const dockets: DocketSummary[] = await docketsResponse.json();
      activeTables = summary.active_tables;
      todayDockets = summary.today_dockets;
      todayRevenue = summary.today_revenue;
      pendingDockets = summary.pending_dockets;
      recentDockets = dockets.slice(0, 3);
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load restaurant dashboard data.';
    }
  }

  function statusBadge(status: string) {
    if (status === 'open') return 'bg-amber-500/20 text-amber-400 border-amber-500/30';
    if (status === 'room_charge') return 'bg-blue-500/20 text-blue-400 border-blue-500/30';
    return 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30';
  }

  function statusLabel(status: string) {
    if (status === 'open') return 'Open';
    if (status === 'room_charge') return 'Room Charge';
    return 'Settled';
  }

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }
</script>

<div class="space-y-5 pb-4">
  <!-- Greeting -->
  <div class="pt-1">
    <p class="text-slate-400 text-xs">{shiftDate}</p>
    <h1 class="text-2xl font-bold text-white mt-0.5">{greeting} 👋</h1>
    <p class="text-slate-400 text-sm">Restaurant Attendant Station</p>
  </div>
  {#if loadError}
    <div role="alert" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">{loadError}</div>
  {/if}

  <!-- KPI Cards -->
  <div class="grid grid-cols-2 gap-3">
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-2">
        <div class="w-7 h-7 rounded-lg bg-orange-500/20 flex items-center justify-center">
          <i class="ph ph-table text-orange-400 text-sm"></i>
        </div>
        <span class="text-slate-400 text-xs">Active Tables</span>
      </div>
      <p class="text-2xl font-bold text-white font-['Oxanium',sans-serif]">{activeTables}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-2">
        <div class="w-7 h-7 rounded-lg bg-blue-500/20 flex items-center justify-center">
          <i class="ph ph-users text-blue-400 text-sm"></i>
        </div>
        <span class="text-slate-400 text-xs">Dockets Today</span>
      </div>
      <p class="text-2xl font-bold text-white font-['Oxanium',sans-serif]">{todayDockets}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-2">
        <div class="w-7 h-7 rounded-lg bg-emerald-500/20 flex items-center justify-center">
          <i class="ph ph-currency-circle-dollar text-emerald-400 text-sm"></i>
        </div>
        <span class="text-slate-400 text-xs">Today's Revenue</span>
      </div>
      <p class="text-lg font-bold text-emerald-400 font-['Oxanium',sans-serif]">{fmtMwk(todayRevenue)}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-2">
        <div class="w-7 h-7 rounded-lg bg-amber-500/20 flex items-center justify-center">
          <i class="ph ph-clock text-amber-400 text-sm"></i>
        </div>
        <span class="text-slate-400 text-xs">Pending Dockets</span>
      </div>
      <p class="text-2xl font-bold text-amber-400 font-['Oxanium',sans-serif]">{pendingDockets}</p>
    </div>
  </div>

  <!-- Quick Actions -->
  <div class="grid grid-cols-3 gap-2.5">
    <a href="/restaurant/dockets" class="flex flex-col items-center justify-center bg-orange-500 hover:bg-orange-400 active:scale-95 text-slate-950 rounded-2xl p-4 font-bold text-xs shadow-lg shadow-orange-500/25 transition">
      <i class="ph ph-plus-circle text-2xl mb-1.5"></i>
      New Docket
    </a>
    <a href="/restaurant/menu" class="flex flex-col items-center justify-center bg-slate-800 hover:bg-slate-700 active:scale-95 text-slate-200 rounded-2xl p-4 font-semibold text-xs border border-slate-700/60 transition">
      <i class="ph ph-list text-2xl mb-1.5"></i>
      View Menu
    </a>
    <a href="/restaurant/reports" class="flex flex-col items-center justify-center bg-slate-800 hover:bg-slate-700 active:scale-95 text-slate-200 rounded-2xl p-4 font-semibold text-xs border border-slate-700/60 transition">
      <i class="ph ph-chart-pie-slice text-2xl mb-1.5"></i>
      Daily Report
    </a>
  </div>

  <!-- Recent Dockets -->
  <div>
    <div class="flex items-center justify-between mb-3">
      <h2 class="font-bold text-white text-sm">Recent Dockets</h2>
      <a href="/restaurant/dockets" class="text-orange-400 text-xs font-semibold">View all →</a>
    </div>
    <div class="space-y-2.5">
      {#each recentDockets as d}
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 flex items-center justify-between">
          <div>
            <div class="flex items-center space-x-2 mb-1">
              <span class="font-['Oxanium',sans-serif] font-bold text-sm text-white">{d.id}</span>
              <span class="text-[10px] px-1.5 py-0.5 rounded-full border font-semibold {statusBadge(d.status)}">{statusLabel(d.status)}</span>
            </div>
            <p class="text-xs text-slate-400">
              {d.tableNo ? 'Table ' + d.tableNo : 'Room ' + d.roomNo} · {d.items.length} item{d.items.length !== 1 ? 's' : ''} · {new Date(d.time).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' })}
            </p>
          </div>
          <p class="text-sm font-bold text-white">{fmtMwk(d.total)}</p>
        </div>
      {/each}
      {#if recentDockets.length === 0 && !loadError}
        <p class="text-center text-xs text-slate-500 py-5">No restaurant dockets found.</p>
      {/if}
    </div>
  </div>
</div>
