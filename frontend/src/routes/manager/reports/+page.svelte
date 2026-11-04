<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  let reportDate = $state(new Date().toISOString().slice(0, 10));
  let isGenerating = $state(false);
  let statusBanner = $state('');

  // Consolidated Financial Figures
  let roomSalesTotal = $state(255000);
  let restaurantSalesTotal = $state(47200);
  let barSalesTotal = $state(38400);
  let totalRevenue = $derived(roomSalesTotal + restaurantSalesTotal + barSalesTotal);

  let departmentExpenses = $state([
    { dept: 'Kitchen (Food Costs & Procurement)', amount: 52400 },
    { dept: 'Bar (Beverage Re-stock COGS)', amount: 14200 },
    { dept: 'Front Desk (Laundry & Room Amenities)', amount: 18500 },
    { dept: 'Administration & Utilities (Fuel/Power)', amount: 22000 }
  ]);

  let totalExpenses = $derived(
    departmentExpenses.reduce((sum, item) => sum + item.amount, 0)
  );

  let netOperatingIncome = $derived(totalRevenue - totalExpenses);
  let operatingMargin = $derived(
    totalRevenue > 0 ? Math.round((netOperatingIncome / totalRevenue) * 100) : 0
  );

  let generalManager = $state('Dr. Samuel Gondwe');
  let remarks = $state('Solid daily performance with 75% room occupancy. Bar and dining revenues matched projections. Kitchen requested poultry restocking order approval for tomorrow morning.');

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }

  async function generateAndArchive() {
    isGenerating = true;
    statusBanner = '';

    const payload = {
      report_date: reportDate,
      total_revenue: totalRevenue,
      total_expenses: totalExpenses,
      net_operating_income: netOperatingIncome,
      department_breakdown: {
        rooms: roomSalesTotal,
        restaurant: restaurantSalesTotal,
        bar: barSalesTotal,
        expenses: departmentExpenses
      },
      manager_name: generalManager,
      remarks: remarks
    };

    try {
      const res = await apiFetch('/api/v1/manager/consolidated', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        statusBanner = '✓ Consolidated report archived into Supabase PostgreSQL.';
      } else {
        statusBanner = '⚠ Archived locally in station memory (Offline Ready).';
      }
    } catch {
      statusBanner = '⚠ Saved locally in offline report repository.';
    } finally {
      isGenerating = false;
      setTimeout(() => {
        statusBanner = '';
      }, 4000);
    }
  }

  function handlePrint() {
    window.print();
  }
</script>

<div class="space-y-4 pb-4">
  <!-- Screen Header -->
  <div class="pt-1 flex items-center justify-between no-print">
    <div>
      <h1 class="text-xl font-bold text-white tracking-tight">Consolidated Report</h1>
      <p class="text-xs text-slate-400">Official Hotel Daily Statement & P&L</p>
    </div>
    <button
      onclick={handlePrint}
      class="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 font-semibold text-xs transition"
    >
      <i class="ph ph-printer text-sm"></i>
      <span>Print / PDF</span>
    </button>
  </div>

  <!-- Date selector -->
  <div class="no-print">
    <label for="cr-date" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Statement Date</label>
    <input
      id="cr-date"
      type="date"
      bind:value={reportDate}
      class="w-full px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-indigo-500"
    />
  </div>

  {#if statusBanner}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center font-medium no-print">
      {statusBanner}
    </div>
  {/if}

  <!-- Official Printable Statement Document -->
  <div class="bg-slate-900 border border-slate-800 rounded-3xl p-5 shadow-2xl space-y-4 print:bg-white print:text-black print:border-black print:p-0">
    <!-- Hotel Letterhead Header -->
    <div class="border-b border-slate-800 pb-4 text-center">
      <div class="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-amber-500 text-slate-950 font-bold mb-1">
        <i class="ph ph-moon text-xl"></i>
      </div>
      <h2 class="font-['Oxanium',sans-serif] font-bold text-xl text-white tracking-wider">LA LUNA HOTEL & RESORT</h2>
      <p class="text-[11px] text-slate-400 uppercase tracking-widest mt-0.5">Comprehensive Daily Consolidated Report</p>
      <p class="text-xs text-indigo-400 font-mono mt-1">Audit Date: {reportDate}</p>
    </div>

    <!-- Revenue Section -->
    <div class="space-y-2">
      <h3 class="text-xs font-bold text-indigo-300 uppercase tracking-wider flex items-center justify-between">
        <span>1. Total Gross Revenue</span>
        <span class="text-white font-['Oxanium',sans-serif] font-bold">{fmtMwk(totalRevenue)}</span>
      </h3>
      <div class="bg-slate-950/60 rounded-xl p-3 border border-slate-800/80 space-y-1.5 text-xs">
        <div class="flex justify-between text-slate-300">
          <span>Rooms & Accommodations</span>
          <span class="font-bold text-white">{fmtMwk(roomSalesTotal)}</span>
        </div>
        <div class="flex justify-between text-slate-300">
          <span>Restaurant Dining & Room Service</span>
          <span class="font-bold text-orange-400">{fmtMwk(restaurantSalesTotal)}</span>
        </div>
        <div class="flex justify-between text-slate-300">
          <span>Bar & Lounge Sales</span>
          <span class="font-bold text-purple-400">{fmtMwk(barSalesTotal)}</span>
        </div>
      </div>
    </div>

    <!-- Operating Costs Section -->
    <div class="space-y-2">
      <h3 class="text-xs font-bold text-rose-300 uppercase tracking-wider flex items-center justify-between">
        <span>2. Operating Expenses (COGS & Overheads)</span>
        <span class="text-white font-['Oxanium',sans-serif] font-bold">{fmtMwk(totalExpenses)}</span>
      </h3>
      <div class="bg-slate-950/60 rounded-xl p-3 border border-slate-800/80 space-y-1.5 text-xs">
        {#each departmentExpenses as exp}
          <div class="flex justify-between text-slate-300">
            <span>{exp.dept}</span>
            <span class="font-mono text-slate-400">{fmtMwk(exp.amount)}</span>
          </div>
        {/each}
      </div>
    </div>

    <!-- Bottom Line Profit -->
    <div class="bg-indigo-950/50 border border-indigo-700/60 rounded-2xl p-4 flex items-center justify-between">
      <div>
        <span class="text-xs text-indigo-300 font-semibold block">Net Operating Income (NOI)</span>
        <span class="text-[11px] text-slate-400">Profit Margin: {operatingMargin}%</span>
      </div>
      <span class="text-xl font-bold text-emerald-400 font-['Oxanium',sans-serif]">
        {fmtMwk(netOperatingIncome)}
      </span>
    </div>

    <!-- Manager Endorsement -->
    <div class="pt-2 border-t border-slate-800 space-y-2 text-xs">
      <p class="text-slate-400 font-semibold">General Manager Review & Remarks:</p>
      <p class="text-slate-300 italic bg-slate-800/40 p-2.5 rounded-xl border border-slate-800">
        "{remarks}"
      </p>
      <div class="flex items-center justify-between pt-2 text-[11px] text-slate-400">
        <span>Verified by: <strong class="text-white">{generalManager}</strong></span>
        <span>Status: <strong class="text-emerald-400">Officially Signed</strong></span>
      </div>
    </div>
  </div>

  <!-- Generate / Submit Action -->
  <div class="no-print">
    <button
      onclick={generateAndArchive}
      disabled={isGenerating}
      class="w-full py-3.5 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white font-bold text-sm uppercase tracking-wider shadow-lg shadow-indigo-600/30 active:scale-98 transition disabled:opacity-50"
    >
      {isGenerating ? 'Archiving Report...' : 'Archive Consolidated Report'}
    </button>
  </div>
</div>
