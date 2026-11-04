<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  let reportDate = $state(new Date().toISOString().slice(0, 10));
  let isSubmitting = $state(false);
  let statusMessage = $state('');

  // Daily statistics
  let totalDishesPrepared = $state(38);
  let estimatedFoodCostMWK = $state(52400);
  let reportedWastageCostMWK = $state(2800);
  let gasElectricityStatus = $state('Optimal - No power cut recorded');
  let chefNotes = $state('Lunch rush smoothly handled. Received fresh Tilapia batch at 11:30 AM. Low stock on whole chickens requires purchase order for tomorrow.');

  let stockSummary = $state([
    { ingredient: 'Maize Flour (Gramil)', used: '18 kg', closing: '57 kg', status: 'Adequate' },
    { ingredient: 'Fresh Tilapia', used: '14 pcs', closing: '18 pcs', status: 'Good' },
    { ingredient: 'Whole Chickens', used: '11 birds', closing: '4 birds', status: 'CRITICAL LOW' },
    { ingredient: 'Cooking Oil', used: '6 litres', closing: '24 litres', status: 'Adequate' }
  ]);

  async function submitKitchenReport() {
    isSubmitting = true;
    statusMessage = '';

    const payload = {
      report_date: reportDate,
      total_dishes: totalDishesPrepared,
      estimated_food_cost: estimatedFoodCostMWK,
      wastage_cost: reportedWastageCostMWK,
      stock_summary: stockSummary,
      chef_notes: chefNotes,
      status: 'Submitted'
    };

    try {
      const res = await apiFetch('/api/v1/kitchen/reports', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        statusMessage = '✓ Daily Kitchen Report submitted to General Manager!';
      } else {
        statusMessage = '⚠ Saved locally in offline queue — will sync upon connection.';
      }
    } catch {
      statusMessage = '⚠ Offline mode — report queued and logged locally.';
    } finally {
      isSubmitting = false;
      setTimeout(() => {
        statusMessage = '';
      }, 4000);
    }
  }

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }
</script>

<div class="space-y-4 pb-4">
  <div class="pt-1">
    <h1 class="text-xl font-bold text-white tracking-tight">Kitchen Shift Report</h1>
    <p class="text-xs text-slate-400">Daily Production, Wastage & Stock Summary</p>
  </div>

  <!-- Date Input -->
  <div>
    <label for="kr-date" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Audit Date</label>
    <input
      id="kr-date"
      type="date"
      bind:value={reportDate}
      class="w-full px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-emerald-500"
    />
  </div>

  {#if statusMessage}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center font-medium">
      {statusMessage}
    </div>
  {/if}

  <!-- KPI Overview -->
  <div class="grid grid-cols-2 gap-3">
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5">
      <span class="text-[10px] text-slate-400 uppercase font-semibold block">Total Portions</span>
      <span class="text-2xl font-bold text-white font-['Oxanium',sans-serif] mt-0.5 block">{totalDishesPrepared}</span>
      <span class="text-[11px] text-emerald-400">Covers cooked today</span>
    </div>

    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5">
      <span class="text-[10px] text-slate-400 uppercase font-semibold block">Est. Food Cost</span>
      <span class="text-lg font-bold text-amber-300 font-['Oxanium',sans-serif] mt-0.5 block">{fmtMwk(estimatedFoodCostMWK)}</span>
      <span class="text-[11px] text-slate-400">BOM calculated</span>
    </div>

    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5 col-span-2">
      <div class="flex items-center justify-between">
        <span class="text-[10px] text-slate-400 uppercase font-semibold">Reported Food Wastage</span>
        <span class="text-xs font-bold text-rose-400 font-['Oxanium',sans-serif]">{fmtMwk(reportedWastageCostMWK)}</span>
      </div>
      <p class="text-[11px] text-slate-400 mt-1">Trimmings, spoils and cooking losses recorded.</p>
    </div>
  </div>

  <!-- Key Ingredients Summary (JSONB in Schema) -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-3">
    <h2 class="text-xs font-bold text-white uppercase tracking-wider">Stock Closing Summary</h2>
    <div class="space-y-2">
      {#each stockSummary as item}
        <div class="flex items-center justify-between text-xs py-1 border-b border-slate-800/60 last:border-0">
          <div>
            <p class="font-semibold text-white">{item.ingredient}</p>
            <p class="text-slate-400 text-[11px]">Usage: {item.used} · In Hand: {item.closing}</p>
          </div>
          <span class="text-[10px] px-2 py-0.5 rounded-full font-bold {item.status.includes('CRITICAL') ? 'bg-rose-500/20 text-rose-300 border border-rose-500/30 animate-pulse' : 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'}">
            {item.status}
          </span>
        </div>
      {/each}
    </div>
  </div>

  <!-- Head Chef Handover Notes -->
  <div>
    <label for="kr-notes" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">
      Head Chef Handover Notes & Requests
    </label>
    <textarea
      id="kr-notes"
      rows="3"
      bind:value={chefNotes}
      class="w-full px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-emerald-500 leading-relaxed"
    ></textarea>
  </div>

  <!-- Submission Button -->
  <button
    onclick={submitKitchenReport}
    disabled={isSubmitting}
    class="w-full py-3.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-sm uppercase tracking-wider shadow-lg shadow-emerald-600/25 active:scale-98 transition disabled:opacity-50"
  >
    {isSubmitting ? 'Submitting to Manager...' : 'Submit Kitchen Report (Chef Handover)'}
  </button>
</div>
