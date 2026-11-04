<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  let reportDate = $state(new Date().toISOString().slice(0, 10));
  let isSubmitting = $state(false);
  let submitMsg = $state('');
  let loadError = $state('');

  let totalSales = $state(0);
  let totalCogs  = $state(0);
  let grossMargin = $derived(totalSales - totalCogs);
  let marginPct   = $derived(totalSales > 0 ? Math.round((grossMargin / totalSales) * 100) : 0);

  type LedgerRow = { name: string; sold: number; revenue: number; cogs: number };
  let ledger = $state<LedgerRow[]>([]);

  onMount(() => {
    void loadReport();
  });

  async function loadReport() {
    loadError = '';
    try {
      const r = await apiFetch(`/api/v1/bar/reports?date=${reportDate}`);
      if (!r.ok) throw new Error('Unable to load the bar report.');
      const d = await r.json();
      totalSales = d.total_bar_sales;
      totalCogs = d.total_cogs;
      ledger = d.ledger;
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load the bar report.';
    }
  }

  async function submitReport() {
    isSubmitting = true;
    submitMsg = '';
    try {
      const r = await apiFetch('/api/v1/bar/reports', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ report_date: reportDate, total_bar_sales: totalSales, total_cogs: totalCogs, ledger }),
      });
      if (!r.ok) throw new Error('Unable to submit the bar report.');
      submitMsg = 'Bar report submitted to Manager.';
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to submit the bar report.';
    }
    isSubmitting = false;
    setTimeout(() => submitMsg = '', 4000);
  }

  function fmtMwk(n: number) { return 'MWK ' + n.toLocaleString(); }
</script>

<div class="space-y-5 pb-4">
  <h1 class="text-xl font-bold text-white pt-1">Bar Daily Report</h1>

  <div>
    <label for="rpt-date" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Report Date</label>
    <input id="rpt-date" type="date" bind:value={reportDate} onchange={loadReport} class="w-full px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-purple-500" />
  </div>

  {#if loadError}
    <div role="alert" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">
      {loadError}
      <button onclick={loadReport} class="ml-2 underline">Retry</button>
    </div>
  {/if}
  {#if submitMsg}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center">{submitMsg}</div>
  {/if}

  <!-- Summary -->
  <div class="grid grid-cols-1 gap-3">
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <p class="text-xs text-slate-400 mb-1">Total Bar Sales</p>
      <p class="text-2xl font-bold text-emerald-400 font-['Oxanium',sans-serif]">{fmtMwk(totalSales)}</p>
    </div>
    <div class="grid grid-cols-3 gap-3">
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3">
        <p class="text-[10px] text-slate-400 mb-1">COGS</p>
        <p class="text-sm font-bold text-rose-400">{fmtMwk(totalCogs)}</p>
      </div>
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3">
        <p class="text-[10px] text-slate-400 mb-1">Gross Margin</p>
        <p class="text-sm font-bold text-purple-300">{fmtMwk(grossMargin)}</p>
      </div>
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3">
        <p class="text-[10px] text-slate-400 mb-1">Margin %</p>
        <p class="text-sm font-bold text-amber-400">{marginPct}%</p>
      </div>
    </div>
  </div>

  <!-- Item Ledger -->
  <div>
    <h2 class="font-bold text-white text-sm mb-3">Item Ledger</h2>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div class="grid grid-cols-4 px-4 py-2 border-b border-slate-800 bg-slate-800/50">
        <span class="text-[10px] text-slate-400 uppercase font-bold col-span-2">Item</span>
        <span class="text-[10px] text-slate-400 uppercase font-bold text-center">Sold</span>
        <span class="text-[10px] text-slate-400 uppercase font-bold text-right">Revenue</span>
      </div>
      {#each ledger as row}
        <div class="grid grid-cols-4 px-4 py-3 border-b border-slate-800/60 last:border-0 items-center">
          <span class="text-xs text-slate-300 col-span-2 truncate">{row.name}</span>
          <span class="text-xs text-white text-center">{row.sold}</span>
          <span class="text-xs font-semibold text-white text-right">{fmtMwk(row.revenue)}</span>
        </div>
      {/each}
    </div>
  </div>

  <button
    onclick={submitReport}
    disabled={isSubmitting}
    class="w-full py-3.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-bold text-sm uppercase tracking-wider shadow-lg active:scale-98 transition disabled:opacity-50"
  >
    {isSubmitting ? 'Submitting...' : 'Submit Bar Report'}
  </button>
</div>
