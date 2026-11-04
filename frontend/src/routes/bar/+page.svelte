<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  let greeting = $state('Good morning');
  let shiftDate = $state('');
  let totalSales = $state(0);
  let totalCogs = $state(0);
  let grossMargin = $derived(totalSales - totalCogs);
  let itemsLow = $state(0);
  let topSeller = $state('—');
  let loadError = $state('');

  onMount(() => {
    const h = new Date().getHours();
    greeting = h < 12 ? 'Good morning' : h < 17 ? 'Good afternoon' : 'Good evening';
    shiftDate = new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
    void loadSummary();
  });

  async function loadSummary() {
    try {
      const date = new Date().toISOString().slice(0, 10);
      const [reportResponse, itemsResponse] = await Promise.all([
        apiFetch(`/api/v1/bar/reports?date=${date}`),
        apiFetch('/api/v1/bar/items')
      ]);
      if (!reportResponse.ok || !itemsResponse.ok) throw new Error('Unable to load bar dashboard data.');
      const report = await reportResponse.json();
      const items = await itemsResponse.json();
      totalSales = report.total_bar_sales;
      totalCogs = report.total_cogs;
      itemsLow = items.filter((item: { stock: number }) => item.stock < 5).length;
      const topItem = [...report.ledger].sort((a: { revenue: number }, b: { revenue: number }) => b.revenue - a.revenue)[0];
      topSeller = topItem?.revenue > 0 ? topItem.name : '—';
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load bar dashboard data.';
    }
  }

  function fmtMwk(n: number) { return 'MWK ' + n.toLocaleString(); }
</script>

<div class="space-y-5 pb-4">
  <div class="pt-1">
    <p class="text-slate-400 text-xs">{shiftDate}</p>
    <h1 class="text-2xl font-bold text-white mt-0.5">{greeting} 🍷</h1>
    <p class="text-slate-400 text-sm">Bar Attendant Station</p>
  </div>
  {#if loadError}
    <div role="alert" class="rounded-xl border border-rose-500/20 bg-rose-500/10 p-3 text-xs text-rose-300">{loadError}</div>
  {/if}

  <div class="grid grid-cols-2 gap-3">
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 col-span-2">
      <div class="flex items-center space-x-2 mb-1">
        <i class="ph ph-currency-circle-dollar text-emerald-400"></i>
        <span class="text-xs text-slate-400">Today's Bar Sales</span>
      </div>
      <p class="text-2xl font-bold text-emerald-400 font-['Oxanium',sans-serif]">{fmtMwk(totalSales)}</p>
      <p class="text-xs text-slate-500 mt-1">COGS: {fmtMwk(totalCogs)} · Gross Margin: {fmtMwk(grossMargin)}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-1">
        <i class="ph ph-warning text-rose-400 text-sm"></i>
        <span class="text-xs text-slate-400">Low Stock Items</span>
      </div>
      <p class="text-2xl font-bold text-rose-400 font-['Oxanium',sans-serif]">{itemsLow}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-1">
        <i class="ph ph-star text-purple-400 text-sm"></i>
        <span class="text-xs text-slate-400">Top Seller</span>
      </div>
      <p class="text-sm font-bold text-white">{topSeller}</p>
    </div>
  </div>

  <div class="grid grid-cols-3 gap-2.5">
    <a href="/bar/sales" class="flex flex-col items-center justify-center bg-purple-600 hover:bg-purple-500 active:scale-95 text-white rounded-2xl p-4 font-bold text-xs shadow-lg shadow-purple-500/25 transition">
      <i class="ph ph-plus-circle text-2xl mb-1.5"></i>
      Record Sale
    </a>
    <a href="/bar/stock" class="flex flex-col items-center justify-center bg-slate-800 hover:bg-slate-700 active:scale-95 text-slate-200 rounded-2xl p-4 font-semibold text-xs border border-slate-700/60 transition">
      <i class="ph ph-package text-2xl mb-1.5"></i>
      Stock Audit
    </a>
    <a href="/bar/reports" class="flex flex-col items-center justify-center bg-slate-800 hover:bg-slate-700 active:scale-95 text-slate-200 rounded-2xl p-4 font-semibold text-xs border border-slate-700/60 transition">
      <i class="ph ph-chart-pie-slice text-2xl mb-1.5"></i>
      Daily Report
    </a>
  </div>

  <!-- Low stock alert -->
  {#if itemsLow > 0}
    <div class="bg-rose-950/50 border border-rose-700/40 rounded-2xl p-4 flex items-start space-x-3">
      <i class="ph ph-warning-circle text-rose-400 text-xl mt-0.5 flex-shrink-0"></i>
      <div>
        <p class="text-sm font-bold text-rose-300">Low Stock Alert</p>
        <p class="text-xs text-rose-400/80 mt-0.5">{itemsLow} item{itemsLow !== 1 ? 's are' : ' is'} running low. Update stock before service.</p>
        <a href="/bar/stock" class="text-xs text-rose-300 font-semibold underline mt-1 inline-block">View Stock →</a>
      </div>
    </div>
  {/if}
</div>
