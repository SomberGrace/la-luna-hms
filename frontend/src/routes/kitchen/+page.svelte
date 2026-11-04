<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';

  type TicketStatus = 'Pending' | 'Cooking' | 'Ready' | 'Served';
  type OrderItem = { name: string; qty: number };
  type OrderTicket = {
    id: string;
    target: string;
    items: OrderItem[];
    status: TicketStatus;
    minutesAgo: number;
    specialInstructions?: string;
  };

  let activeTickets = $state<OrderTicket[]>([]);
  let isLoading = $state(true);
  let updatingTicket = $state('');
  let loadError = $state('');
  let toastMsg = $state('');

  let pendingCount = $derived(activeTickets.filter(ticket => ticket.status === 'Pending').length);
  let cookingCount = $derived(activeTickets.filter(ticket => ticket.status === 'Cooking').length);
  let readyCount = $derived(activeTickets.filter(ticket => ticket.status === 'Ready').length);

  async function loadTickets() {
    try {
      const response = await apiFetch('/api/v1/kitchen/dockets');
      if (!response.ok) throw new Error('Unable to load kitchen tickets.');
      const records: {
        id: string;
        roomNo: string;
        tableNo: string;
        createdAt: string;
        kitchenStatus: TicketStatus;
        notes: string;
        items: OrderItem[];
      }[] = await response.json();
      activeTickets = records.map((record) => ({
        id: record.id,
        target: record.tableNo ? `Table ${record.tableNo}` : `Room ${record.roomNo}`,
        items: record.items,
        status: record.kitchenStatus,
        minutesAgo: Math.max(0, Math.floor((Date.now() - Date.parse(record.createdAt)) / 60000)),
        specialInstructions: record.notes || undefined
      }));
      loadError = '';
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load kitchen tickets.';
    } finally {
      isLoading = false;
    }
  }

  async function nextStatus(ticket: OrderTicket) {
    const next: Partial<Record<TicketStatus, TicketStatus>> = {
      Pending: 'Cooking',
      Cooking: 'Ready',
      Ready: 'Served'
    };
    const status = next[ticket.status];
    if (!status) return;

    updatingTicket = ticket.id;
    try {
      const response = await apiFetch('/api/v1/kitchen/dockets/status', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ docket_id: ticket.id, status })
      });
      if (!response.ok) {
        const body = await response.json().catch(() => null);
        throw new Error(body?.error || 'Unable to update kitchen ticket.');
      }
      toastMsg = `${ticket.id} marked ${status.toLowerCase()}.`;
      await loadTickets();
      setTimeout(() => { toastMsg = ''; }, 2500);
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to update kitchen ticket.';
    } finally {
      updatingTicket = '';
    }
  }

  function getStatusBadge(status: OrderTicket['status']) {
    switch (status) {
      case 'Pending':
        return 'bg-amber-500/20 text-amber-300 border-amber-500/40';
      case 'Cooking':
        return 'bg-blue-500/20 text-blue-300 border-blue-500/40 animate-pulse';
      case 'Ready':
        return 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40';
      case 'Served':
        return 'bg-slate-800 text-slate-400 border-slate-700';
    }
  }

  function getNextBtnText(status: OrderTicket['status']) {
    switch (status) {
      case 'Pending':
        return 'Start Cooking 🔥';
      case 'Cooking':
        return 'Mark Ready 🍽️';
      case 'Ready':
        return 'Complete & Serve ✨';
      case 'Served':
        return 'Completed';
    }
  }

  onMount(() => {
    void loadTickets();
    const refresh = setInterval(() => { void loadTickets(); }, 15000);
    return () => clearInterval(refresh);
  });
</script>

<div class="space-y-4 pb-4">
  <!-- Station Header -->
  <div class="pt-1 flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold text-white tracking-tight">Kitchen Orders</h1>
      <p class="text-xs text-slate-400">Live Kitchen Display System (KDS)</p>
    </div>
    <div class="flex items-center space-x-1.5 px-2.5 py-1 rounded-full bg-emerald-950/60 border border-emerald-700/50 text-emerald-400 text-xs">
      <span class="w-2 h-2 rounded-full bg-emerald-400 animate-ping"></span>
      <span class="text-[11px] font-semibold">KDS Live</span>
    </div>
  </div>

  {#if toastMsg}
    <div class="p-3 rounded-xl bg-emerald-500/20 border border-emerald-500/30 text-emerald-300 text-xs text-center font-medium">
      {toastMsg}
    </div>
  {/if}
  {#if loadError}
    <div role="alert" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">
      {loadError}
    </div>
  {/if}

  <!-- Kitchen KPIs -->
  <div class="grid grid-cols-3 gap-2.5">
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3 text-center">
      <p class="text-[10px] text-slate-400 uppercase font-semibold">Pending</p>
      <p class="text-xl font-bold text-amber-400 font-['Oxanium',sans-serif] mt-0.5">{pendingCount}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3 text-center">
      <p class="text-[10px] text-slate-400 uppercase font-semibold">Cooking</p>
      <p class="text-xl font-bold text-blue-400 font-['Oxanium',sans-serif] mt-0.5">{cookingCount}</p>
    </div>
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3 text-center">
      <p class="text-[10px] text-slate-400 uppercase font-semibold">Ready</p>
      <p class="text-xl font-bold text-emerald-400 font-['Oxanium',sans-serif] mt-0.5">{readyCount}</p>
    </div>
  </div>

  <!-- Active Order Tickets -->
  <div class="space-y-3">
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-bold text-white tracking-wide uppercase">Current Tickets</h2>
      <a href="/kitchen/stock" class="text-xs text-emerald-400 hover:underline">Check Ingredients &rarr;</a>
    </div>

    {#if isLoading}
      <div class="text-center py-12 text-slate-500 text-sm">Loading kitchen tickets…</div>
    {/if}
    {#each activeTickets as ticket}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-lg space-y-3 transition-all {ticket.status === 'Served' ? 'opacity-40' : ''}">
        <!-- Ticket Header -->
        <div class="flex items-center justify-between border-b border-slate-800/80 pb-2.5">
          <div class="flex items-center space-x-2">
            <span class="font-['Oxanium',sans-serif] font-bold text-white text-base">{ticket.id}</span>
            <span class="text-xs font-semibold px-2 py-0.5 rounded-lg bg-slate-800 text-slate-200 border border-slate-700">
              {ticket.target}
            </span>
          </div>
          <div class="flex items-center space-x-2">
            <span class="text-[11px] text-slate-400 flex items-center">
              <i class="ph ph-clock mr-1"></i> {ticket.minutesAgo}m ago
            </span>
            <span class="text-[10px] px-2 py-0.5 rounded-full border font-bold uppercase tracking-wider {getStatusBadge(ticket.status)}">
              {ticket.status}
            </span>
          </div>
        </div>

        <!-- Ticket Items -->
        <div class="space-y-1.5">
          {#each ticket.items as item}
            <div class="flex items-start justify-between">
              <div class="flex items-center space-x-2">
                <span class="w-5 h-5 rounded-md bg-emerald-500/20 text-emerald-400 text-xs font-bold flex items-center justify-center">
                  {item.qty}
                </span>
                <span class="text-sm font-semibold text-white">{item.name}</span>
              </div>
            </div>
          {/each}
        </div>

        {#if ticket.specialInstructions}
          <div class="bg-amber-950/40 border border-amber-600/30 rounded-xl p-2 text-[11px] text-amber-200 flex items-center space-x-2">
            <i class="ph ph-warning-circle text-amber-400 text-sm flex-shrink-0"></i>
            <span>{ticket.specialInstructions}</span>
          </div>
        {/if}

        <!-- Status Action -->
        {#if ticket.status !== 'Served'}
          <button
            onclick={() => nextStatus(ticket)}
            disabled={updatingTicket === ticket.id}
            class="w-full py-2.5 rounded-xl font-bold text-xs uppercase tracking-wider transition active:scale-98 shadow-md flex items-center justify-center space-x-1.5 {ticket.status === 'Pending' ? 'bg-amber-500 hover:bg-amber-400 text-slate-950' : ticket.status === 'Cooking' ? 'bg-blue-600 hover:bg-blue-500 text-white' : 'bg-emerald-600 hover:bg-emerald-500 text-white'}"
          >
            <span>{getNextBtnText(ticket.status)}</span>
          </button>
        {/if}
      </div>
    {/each}

    {#if !isLoading && activeTickets.length === 0 && !loadError}
      <div class="text-center py-12 text-slate-500">
        <i class="ph ph-check-circle text-4xl mb-3 block text-emerald-500"></i>
        <p class="text-sm font-semibold text-slate-300">All orders served!</p>
        <p class="text-xs text-slate-500 mt-1">Waiting for new restaurant dockets.</p>
      </div>
    {/if}
  </div>
</div>
