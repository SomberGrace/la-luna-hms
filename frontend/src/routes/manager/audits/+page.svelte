<script lang="ts">
  type AuditItem = {
    id: string;
    department: 'Front Desk' | 'Bar' | 'Kitchen' | 'Restaurant';
    clerk: string;
    timestamp: string;
    keyMetricLabel: string;
    keyMetricValue: string;
    details: string;
    status: 'Submitted' | 'Approved' | 'Disputed';
  };

  let activeTab = $state<'All' | 'Front Desk' | 'Bar' | 'Kitchen' | 'Restaurant'>('All');
  let toastMsg = $state('');

  let audits = $state<AuditItem[]>([
    {
      id: 'AUD-FD-01',
      department: 'Front Desk',
      clerk: 'Faith Mwanza (Reception)',
      timestamp: '14:00 Today',
      keyMetricLabel: 'Room Revenue',
      keyMetricValue: 'MWK 255,000',
      details: '6 Occupied rooms, 2 Check-outs with Zero Balance verified. Handover Float: MWK 50,000.',
      status: 'Submitted'
    },
    {
      id: 'AUD-BAR-01',
      department: 'Bar',
      clerk: 'Peter Banda (Barman)',
      timestamp: '14:15 Today',
      keyMetricLabel: 'Bar Sales & COGS',
      keyMetricValue: 'MWK 38,400',
      details: 'Stock audited: 12 Castle Lagers, 17 Carlsberg Greens sold. Physical closing matches ledger balance.',
      status: 'Submitted'
    },
    {
      id: 'AUD-KT-01',
      department: 'Kitchen',
      clerk: 'James Phiri (Head Chef)',
      timestamp: '14:30 Today',
      keyMetricLabel: 'Portions & Wastage',
      keyMetricValue: '38 Portions',
      details: 'Food cost MWK 52,400. Wastage recorded MWK 2,800. Chicken stock critically low (4 birds remaining).',
      status: 'Submitted'
    },
    {
      id: 'AUD-REST-01',
      department: 'Restaurant',
      clerk: 'Grace Tembo (Waiter/Captain)',
      timestamp: '14:45 Today',
      keyMetricLabel: 'Dining Revenue',
      keyMetricValue: 'MWK 47,200',
      details: '14 Covers served. All dockets settled or charged to room folios.',
      status: 'Approved'
    }
  ]);

  let filteredAudits = $derived(
    activeTab === 'All'
      ? audits
      : audits.filter(a => a.department === activeTab)
  );

  function approveAudit(id: string) {
    audits = audits.map(a => (a.id === id ? { ...a, status: 'Approved' } : a));
    toastMsg = `Audit ${id} officially Approved and signed off by GM.`;
    setTimeout(() => {
      toastMsg = '';
    }, 3000);
  }

  function disputeAudit(id: string) {
    const reason = prompt('Enter discrepancy reason to return for correction:', 'Physical cash discrepancy detected');
    if (reason) {
      audits = audits.map(a => (a.id === id ? { ...a, status: 'Disputed' } : a));
      toastMsg = `Audit ${id} flagged as Disputed (${reason}).`;
      setTimeout(() => {
        toastMsg = '';
      }, 3500);
    }
  }

  function getStatusStyle(status: AuditItem['status']) {
    switch (status) {
      case 'Submitted':
        return 'bg-amber-500/20 text-amber-300 border-amber-500/30';
      case 'Approved':
        return 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30';
      case 'Disputed':
        return 'bg-rose-500/20 text-rose-300 border-rose-500/30 animate-pulse';
    }
  }

  function getDeptBadge(dept: AuditItem['department']) {
    switch (dept) {
      case 'Front Desk':
        return 'bg-amber-500/20 text-amber-300';
      case 'Bar':
        return 'bg-purple-500/20 text-purple-300';
      case 'Kitchen':
        return 'bg-emerald-500/20 text-emerald-300';
      case 'Restaurant':
        return 'bg-orange-500/20 text-orange-300';
    }
  }
</script>

<div class="space-y-4 pb-4">
  <div class="pt-1">
    <h1 class="text-xl font-bold text-white tracking-tight">Shift Audits & Approvals</h1>
    <p class="text-xs text-slate-400">Department Daily Closures Sign-Off Console</p>
  </div>

  {#if toastMsg}
    <div class="p-3 rounded-xl bg-indigo-500/20 border border-indigo-500/40 text-indigo-200 text-xs text-center font-medium">
      {toastMsg}
    </div>
  {/if}

  <!-- Department Tabs -->
  <div class="flex space-x-2 overflow-x-auto pb-1 scrollbar-none">
    {#each ['All', 'Front Desk', 'Bar', 'Kitchen', 'Restaurant'] as tab}
      <button
        onclick={() => (activeTab = tab as any)}
        class="flex-shrink-0 px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all {activeTab === tab ? 'bg-indigo-600 text-white border-indigo-600 shadow-lg shadow-indigo-600/30' : 'bg-slate-900 text-slate-400 border-slate-700 hover:border-slate-500'}"
      >
        {tab}
      </button>
    {/each}
  </div>

  <!-- Audit Items -->
  <div class="space-y-3">
    {#each filteredAudits as item}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-lg space-y-3">
        <!-- Card Header -->
        <div class="flex items-start justify-between">
          <div>
            <div class="flex items-center space-x-2 mb-1">
              <span class="font-['Oxanium',sans-serif] font-bold text-white text-sm">{item.id}</span>
              <span class="text-[10px] px-2 py-0.5 rounded-full font-bold uppercase {getDeptBadge(item.department)}">
                {item.department}
              </span>
            </div>
            <p class="text-xs text-slate-400">{item.clerk} · {item.timestamp}</p>
          </div>
          <span class="text-[10px] px-2.5 py-1 rounded-full border font-bold uppercase tracking-wider {getStatusStyle(item.status)}">
            {item.status}
          </span>
        </div>

        <!-- Metric Highlight -->
        <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-3 flex items-center justify-between">
          <span class="text-xs text-slate-400">{item.keyMetricLabel}</span>
          <span class="text-base font-bold text-white font-['Oxanium',sans-serif]">{item.keyMetricValue}</span>
        </div>

        <!-- Details -->
        <p class="text-xs text-slate-300 leading-relaxed bg-slate-800/40 p-2.5 rounded-xl border border-slate-800/60">
          {item.details}
        </p>

        <!-- Actions -->
        {#if item.status === 'Submitted'}
          <div class="grid grid-cols-2 gap-2 pt-1">
            <button
              onclick={() => approveAudit(item.id)}
              class="py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs uppercase tracking-wider transition active:scale-95 shadow-md flex items-center justify-center space-x-1"
            >
              <i class="ph ph-check text-sm"></i>
              <span>Approve Audit</span>
            </button>
            <button
              onclick={() => disputeAudit(item.id)}
              class="py-2.5 rounded-xl bg-rose-600/20 hover:bg-rose-600/30 border border-rose-500/40 text-rose-300 font-bold text-xs uppercase tracking-wider transition active:scale-95 flex items-center justify-center space-x-1"
            >
              <i class="ph ph-warning text-sm"></i>
              <span>Flag Discrepancy</span>
            </button>
          </div>
        {:else if item.status === 'Approved'}
          <div class="flex items-center justify-center space-x-1.5 py-1.5 text-xs text-emerald-400 font-semibold bg-emerald-950/30 rounded-xl border border-emerald-800/30">
            <i class="ph ph-seal-check text-base"></i>
            <span>Signed off & Added to Consolidated Ledger</span>
          </div>
        {:else}
          <div class="flex items-center justify-center space-x-1.5 py-1.5 text-xs text-rose-300 font-semibold bg-rose-950/30 rounded-xl border border-rose-800/30">
            <i class="ph ph-x-circle text-base"></i>
            <span>Disputed — Waiting on Department Correction</span>
          </div>
        {/if}
      </div>
    {/each}
  </div>
</div>

