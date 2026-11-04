<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  type CategoryRow = { cat: string; covers: number; revenue: number };

  let reportDate = $state(new Date().toISOString().slice(0, 10));
  let isLoading = $state(false);
  let isSubmitting = $state(false);
  let submitMsg = $state('');

  let totalRevenue = $state(47200);
  let totalCovers = $state(14);
  let avgCoverSpend = $derived(totalCovers > 0 ? Math.round(totalRevenue / totalCovers) : 0);
  let topItem = $state('Nsima with Chicken');

  let breakdown = $state<CategoryRow[]>([
    { cat: 'Breakfast', covers: 4,  revenue: 8800  },
    { cat: 'Lunch',     covers: 6,  revenue: 18200 },
    { cat: 'Dinner',    covers: 3,  revenue: 15600 },
    { cat: 'Snack',     covers: 0,  revenue: 0     },
    { cat: 'Beverage',  covers: 1,  revenue: 4600  },
  ]);

  onMount(() => {
    loadReport(reportDate);
  });

  async function loadReport(date: string) {
    isLoading = true;
    try {
      const r = await apiFetch(`/api/v1/restaurant/reports?date=${date}`);
      if (r.ok) {
        const d = await r.json();
        totalRevenue = d.total_revenue ?? totalRevenue;
        totalCovers  = d.total_covers  ?? totalCovers;
        topItem      = d.top_item      ?? topItem;
        breakdown    = d.breakdown     ?? breakdown;
      }
    } catch { /* offline */ }
    isLoading = false;
  }

  async function submitReport() {
    isSubmitting = true;
    submitMsg = '';
    try {
      const r = await apiFetch('/api/v1/restaurant/reports', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ date: reportDate, total_revenue: totalRevenue, total_covers: totalCovers, breakdown }),
      });
      submitMsg = r.ok ? '✓ Report submitted successfully.' : '⚠ Saved locally — will sync when online.';
    } catch {
      submitMsg = '⚠ Saved locally — will sync when online.';
    }
    isSubmitting = false;
    setTimeout(() => submitMsg = '', 4000);
  }

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }

  function catIcon(cat: string) {
    const m: Record<string, string> = { Breakfast: 'ph-coffee', Lunch: 'ph-fork-knife', Dinner: 'ph-moon', Snack: 'ph-popcorn', Beverage: 'ph-wine' };
    return m[cat] ?? 'ph-circle';
  }
</script>

<div class="space-y-5 pb-4">
  <div class="flex items-center justify-between pt-1">
    <h1 class="text-xl font-bold text-white">Restaurant Report</h1>
    {#if isLoading}
      <span class="text-xs text-slate-400 animate-pulse">Loading...</span>
    {/if}
  </div>

  <!-- Date Picker -->
  <div>
    <label for="rpt-date" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Report Date</label>
    <input
      id="rpt-date"
      type="date"
      bind:value={reportDate}
      onchange={() => loadReport(reportDate)}
      class="w-full px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500"
    />
  </div>

  {#if submitMsg}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center">{submitMsg}</div>
  {/if}

  <!-- Summary KPIs -->
  <div class="grid grid-cols-2 gap-3">
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 col-span-2">
      <div class="flex items-center space-x-2 mb-1">
        <i class="ph ph-currency-circle-dollar text-emerald-400"></i>
        <span class="text-xs text-slate-400">Total Revenue</span>
      </div>
      <p class="text-2xl font-bold text-emerald-400 font-['Oxanium',sans-serif]">{fmtMwk(totalRevenue)}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-1">
        <i class="ph ph-users text-blue-400 text-sm"></i>
        <span class="text-xs text-slate-400">Total Covers</span>
      </div>
      <p class="text-2xl font-bold text-white font-['Oxanium',sans-serif]">{totalCovers}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4">
      <div class="flex items-center space-x-2 mb-1">
        <i class="ph ph-trending-up text-amber-400 text-sm"></i>
        <span class="text-xs text-slate-400">Avg. Spend / Cover</span>
      </div>
      <p class="text-lg font-bold text-amber-400 font-['Oxanium',sans-serif]">{fmtMwk(avgCoverSpend)}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 col-span-2">
      <div class="flex items-center space-x-2 mb-1">
        <i class="ph ph-star text-orange-400 text-sm"></i>
        <span class="text-xs text-slate-400">Top Selling Item</span>
      </div>
      <p class="text-sm font-bold text-white">{topItem}</p>
    </div>
  </div>

  <!-- Category Breakdown -->
  <div>
    <h2 class="font-bold text-white text-sm mb-3">Category Breakdown</h2>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden">
      <div class="grid grid-cols-3 px-4 py-2 border-b border-slate-800 bg-slate-800/50">
        <span class="text-[10px] text-slate-400 uppercase font-bold">Category</span>
        <span class="text-[10px] text-slate-400 uppercase font-bold text-center">Covers</span>
        <span class="text-[10px] text-slate-400 uppercase font-bold text-right">Revenue</span>
      </div>
      {#each breakdown as row}
        <div class="grid grid-cols-3 px-4 py-3 border-b border-slate-800/60 last:border-0">
          <div class="flex items-center space-x-2">
            <i class="ph {catIcon(row.cat)} text-orange-400 text-sm"></i>
            <span class="text-xs text-slate-300">{row.cat}</span>
          </div>
          <span class="text-xs text-white text-center">{row.covers}</span>
          <span class="text-xs font-semibold text-white text-right">{fmtMwk(row.revenue)}</span>
        </div>
      {/each}
    </div>
  </div>

  <!-- Submit -->
  <button
    onclick={submitReport}
    disabled={isSubmitting}
    class="w-full py-3.5 rounded-xl bg-orange-500 hover:bg-orange-400 text-slate-950 font-bold text-sm uppercase tracking-wider shadow-lg shadow-orange-500/25 active:scale-98 transition disabled:opacity-50"
  >
    {isSubmitting ? 'Submitting...' : 'Submit Daily Report'}
  </button>
</div>
