<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { onMount } from 'svelte';
  type MenuItem = { id: string; name: string; price: number; category: string; available: boolean };
  type DocketItem = { menuItemId: string; name: string; price: number; qty: number };
  type Docket = {
    id: string;
    tableNo: string;
    roomNo: string;
    items: DocketItem[];
    status: 'open' | 'room_charge' | 'settled';
    kitchenStatus: 'Pending' | 'Cooking' | 'Ready' | 'Served';
    time: string;
    notes: string;
  };
  type MealGuest = {
    room_number: string;
    card_no: string;
    guest_name: string;
    room_type: string;
    included_covers: number;
    reconciled_today?: boolean;
    reconciled_breakfast?: boolean;
    reconciled_dinner?: boolean;
  };

  let activeTab = $state<'dockets' | 'reconciliation'>('dockets');
  let activeMeal = $state<'breakfast' | 'dinner'>('breakfast');

  let dockets = $state<Docket[]>([]);
  let menuItems = $state<MenuItem[]>([]);
  let occupiedRooms = $state<string[]>([]);
  let mealGuests = $state<MealGuest[]>([]);
  let isLoading = $state(true);
  let loadingGuests = $state(true);
  let loadError = $state('');

  // New docket form
  let showNewForm = $state(false);
  let mode = $state<'table' | 'room'>('table');
  let tableNo = $state('');
  let roomNo = $state('');
  let docketItems = $state<DocketItem[]>([]);
  let notes = $state('');
  let isSubmitting = $state(false);
  let isSettling = $state(false);
  let toast = $state('');

  let docketTotal = $derived(docketItems.reduce((s, i) => s + i.price * i.qty, 0));

  let breakfastServedCount = $derived(mealGuests.filter(g => g.reconciled_breakfast || g.reconciled_today).length);
  let dinnerServedCount = $derived(mealGuests.filter(g => g.reconciled_dinner).length);

  onMount(() => {
    void loadDockets();
    void loadMenu();
    void loadOccupiedRooms();
    void fetchGuestsList();
  });

  async function loadDockets() {
    isLoading = true;
    try {
      const res = await apiFetch('/api/v1/restaurant/dockets');
      if (!res.ok) throw new Error('Unable to load restaurant dockets.');
      dockets = await res.json();
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load restaurant dockets.';
    } finally {
      isLoading = false;
    }
  }

  async function loadMenu() {
    try {
      const res = await apiFetch('/api/v1/restaurant/menu');
      if (!res.ok) throw new Error('Unable to load menu items.');
      menuItems = await res.json();
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load menu items.';
    }
  }

  async function loadOccupiedRooms() {
    try {
      const res = await apiFetch('/api/v1/restaurant/occupied-rooms');
      if (!res.ok) throw new Error('Unable to load occupied rooms.');
      occupiedRooms = await res.json();
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load occupied rooms.';
    }
  }

  async function fetchGuestsList() {
    loadingGuests = true;
    try {
      const res = await apiFetch('/api/v1/restaurant/breakfast');
      if (!res.ok) throw new Error('Unable to load occupied guest rooms.');
      const data: MealGuest[] = await res.json();
      mealGuests = data.map((guest) => ({
        ...guest,
        reconciled_breakfast: guest.reconciled_breakfast ?? guest.reconciled_today ?? false,
        reconciled_dinner: guest.reconciled_dinner ?? false
      }));
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load occupied guest rooms.';
    } finally {
      loadingGuests = false;
    }
  }

  async function reconcileMeal(mg: MealGuest, meal: 'breakfast' | 'dinner') {
    const mealLabel = meal === 'breakfast' ? 'Breakfast' : 'Dinner';
    try {
      const res = await apiFetch(`/api/v1/restaurant/${meal}/reconcile`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          room_number: mg.room_number,
          covers: mg.included_covers,
          notes: `Room-assigned ${mealLabel}`
        })
      });
      if (!res.ok) throw new Error(`Unable to reconcile ${mealLabel.toLowerCase()}.`);
      toast = `${mealLabel} reconciled for Room ${mg.room_number} (${mg.included_covers} cover(s)).`;
      mealGuests = mealGuests.map(g => {
        if (g.room_number !== mg.room_number) return g;
        return meal === 'breakfast'
          ? { ...g, reconciled_breakfast: true, reconciled_today: true }
          : { ...g, reconciled_dinner: true };
      });
    } catch (error) {
      loadError = error instanceof Error ? error.message : `Unable to reconcile ${mealLabel.toLowerCase()}.`;
    }
    setTimeout(() => { toast = ''; }, 3000);
  }

  function addMenuItem(m: MenuItem) {
    const existing = docketItems.find(i => i.menuItemId === m.id);
    if (existing) {
      docketItems = docketItems.map(i => i.menuItemId === m.id ? { ...i, qty: i.qty + 1 } : i);
    } else {
      docketItems = [...docketItems, { menuItemId: m.id, name: m.name, price: m.price, qty: 1 }];
    }
  }

  function removeItem(id: string) {
    docketItems = docketItems.filter(i => i.menuItemId !== id);
  }

  function changeQty(id: string, delta: number) {
    docketItems = docketItems.map(i => {
      if (i.menuItemId !== id) return i;
      const q = i.qty + delta;
      return q <= 0 ? null : { ...i, qty: q };
    }).filter(Boolean) as DocketItem[];
  }

  async function submitDocket(chargeType: 'table' | 'room_charge') {
    if (docketItems.length === 0) return;
    isSubmitting = true;
    try {
      const res = await apiFetch('/api/v1/restaurant/dockets', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          table_no: chargeType === 'table' ? tableNo : '',
          room_no: chargeType === 'room_charge' ? roomNo : '',
          items: docketItems.map(({ menuItemId, qty }) => ({ menuItemId, qty })),
          notes,
          charge_type: chargeType
        })
      });
      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error(body?.error || 'Unable to create docket.');
      }
      tableNo = ''; roomNo = ''; docketItems = []; notes = '';
      showNewForm = false;
      toast = chargeType === 'room_charge' ? 'Charged to room folio and sent to kitchen.' : 'Docket submitted to kitchen.';
      await loadDockets();
      setTimeout(() => toast = '', 3000);
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to create docket.';
    } finally {
      isSubmitting = false;
    }
  }

  async function settleDocket(id: string) {
    isSettling = true;
    try {
      const res = await apiFetch('/api/v1/restaurant/dockets/settle', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ docket_id: id })
      });
      if (!res.ok) {
        const body = await res.json().catch(() => null);
        throw new Error(body?.error || 'Unable to settle docket.');
      }
      toast = `${id} settled.`;
      await loadDockets();
      setTimeout(() => toast = '', 2500);
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to settle docket.';
    } finally {
      isSettling = false;
    }
  }

  function statusBadge(s: string) {
    if (s === 'open') return 'bg-amber-500/20 text-amber-400 border-amber-500/30';
    if (s === 'room_charge') return 'bg-blue-500/20 text-blue-400 border-blue-500/30';
    return 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30';
  }
  function statusLabel(s: string) {
    return s === 'open' ? 'Open' : s === 'room_charge' ? 'Room Charge' : 'Settled';
  }
  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }
</script>

<div class="space-y-4 pb-4">
  <!-- Header with Tab Navigation -->
  <div class="flex items-center justify-between pt-1">
    <div>
      <h1 class="text-xl font-bold text-white font-['Oxanium',sans-serif]">Restaurant Operations</h1>
      <p class="text-xs text-slate-400">Order tickets & morning breakfast reconciliation</p>
    </div>

    {#if activeTab === 'dockets'}
      <button
        onclick={() => showNewForm = true}
        class="flex items-center space-x-1.5 px-4 py-2 bg-orange-500 hover:bg-orange-400 text-slate-950 font-bold text-xs rounded-xl active:scale-95 transition shadow-lg shadow-orange-500/20"
      >
        <i class="ph ph-plus text-sm"></i>
        <span>New Docket</span>
      </button>
    {/if}
  </div>

  <!-- Tab Selector -->
  <div class="flex bg-slate-900 border border-slate-800 rounded-2xl p-1">
    <button
      onclick={() => activeTab = 'dockets'}
      class="flex-1 py-2.5 rounded-xl text-xs font-bold transition flex items-center justify-center space-x-1.5 {activeTab === 'dockets' ? 'bg-orange-500 text-slate-950 shadow-md' : 'text-slate-400 hover:text-white'}"
    >
      <i class="ph ph-receipt text-sm"></i>
      <span>Active Dockets ({dockets.filter(d => d.status !== 'settled').length})</span>
    </button>
    <button
      onclick={() => activeTab = 'reconciliation'}
      class="flex-1 py-2.5 rounded-xl text-xs font-bold transition flex items-center justify-center space-x-1.5 {activeTab === 'reconciliation' ? 'bg-orange-500 text-slate-950 shadow-md' : 'text-slate-400 hover:text-white'}"
    >
      <i class="ph ph-fork-knife text-sm"></i>
      <span>Meal Reconciliation</span>
    </button>
  </div>

  {#if toast}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center font-bold">{toast}</div>
  {/if}
  {#if loadError}
    <div role="alert" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">{loadError}</div>
  {/if}

  <!-- TAB 1: DOCKETS LIST -->
  {#if activeTab === 'dockets'}
    <div class="space-y-3">
      {#if isLoading}
        <div class="text-center py-12 text-slate-500 text-sm">Loading dockets…</div>
      {:else}
      {#each dockets as d}
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-md">
          <div class="flex items-start justify-between mb-3">
            <div>
              <div class="flex items-center space-x-2 mb-0.5">
                <span class="font-['Oxanium',sans-serif] font-bold text-white">{d.id}</span>
                <span class="text-[10px] px-1.5 py-0.5 rounded-full border font-semibold {statusBadge(d.status)}">{statusLabel(d.status)}</span>
              </div>
              <p class="text-xs text-slate-400">
                {d.tableNo ? 'Table ' + d.tableNo : 'Room ' + d.roomNo} · {new Date(d.time).toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' })}
              </p>
            </div>
            <span class="font-bold text-white font-['Oxanium',sans-serif]">{fmtMwk(d.items.reduce((s, i) => s + i.price * i.qty, 0))}</span>
          </div>
          <div class="space-y-1 mb-3">
            {#each d.items as item}
              <div class="flex justify-between text-xs text-slate-400">
                <span>{item.name} × {item.qty}</span>
                <span>{fmtMwk(item.price * item.qty)}</span>
              </div>
            {/each}
            {#if d.notes}
              <p class="text-xs text-slate-500 italic">Note: {d.notes}</p>
            {/if}
            {#if d.status !== 'settled'}
              <p class="text-[11px] text-slate-500">Kitchen: {d.kitchenStatus}</p>
            {/if}
          </div>
          {#if d.status !== 'settled' && d.kitchenStatus === 'Served'}
            <button
              onclick={() => settleDocket(d.id)}
              disabled={isSettling}
              class="w-full py-2 rounded-xl bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-400 font-semibold text-xs border border-emerald-500/30 transition disabled:opacity-50"
            >
              Mark as Settled
            </button>
          {:else if d.status !== 'settled'}
            <p class="w-full py-2 text-center rounded-xl bg-slate-800 text-slate-500 text-xs border border-slate-700">Settlement available after kitchen marks served</p>
          {/if}
        </div>
      {/each}
      {#if dockets.length === 0 && !loadError}
        <div class="text-center py-12 text-slate-500">
          <i class="ph ph-receipt text-4xl mb-3 block"></i>
          <p class="text-sm">No dockets found</p>
        </div>
      {/if}
      {/if}
    </div>

  <!-- TAB 2: ROOM MEAL RECONCILIATION (BREAKFAST & DINNER) -->
  {:else}
    <div class="space-y-3">
      <!-- Meal Type Selector Toggle -->
      <div class="grid grid-cols-2 gap-2 bg-slate-900 border border-slate-800 rounded-2xl p-1.5">
        <button
          type="button"
          onclick={() => activeMeal = 'breakfast'}
          class="py-2.5 rounded-xl text-xs font-bold transition flex items-center justify-center space-x-1.5 {activeMeal === 'breakfast' ? 'bg-amber-500 text-slate-950 shadow-md' : 'text-slate-400 hover:text-white'}"
        >
          <i class="ph ph-coffee text-sm"></i>
          <span>Breakfast ({breakfastServedCount}/{mealGuests.length})</span>
        </button>
        <button
          type="button"
          onclick={() => activeMeal = 'dinner'}
          class="py-2.5 rounded-xl text-xs font-bold transition flex items-center justify-center space-x-1.5 {activeMeal === 'dinner' ? 'bg-purple-500 text-white shadow-md' : 'text-slate-400 hover:text-white'}"
        >
          <i class="ph ph-moon text-sm"></i>
          <span>Dinner ({dinnerServedCount}/{mealGuests.length})</span>
        </button>
      </div>

      {#if activeMeal === 'breakfast'}
        <div class="bg-amber-500/10 border border-amber-500/30 rounded-2xl p-3.5 text-xs text-amber-300">
          <div class="flex items-center space-x-2 font-bold mb-1">
            <i class="ph ph-coffee text-sm"></i>
            <span>Breakfast Reconciliation Policy</span>
          </div>
          <p class="text-[11px] text-slate-300 leading-relaxed">
            In-house guests are entitled to complimentary breakfast based on room bed type. Reconcile room covers upon serving. Extra covers can be charged as room dockets.
          </p>
        </div>
      {:else}
        <div class="bg-purple-500/10 border border-purple-500/30 rounded-2xl p-3.5 text-xs text-purple-300">
          <div class="flex items-center space-x-2 font-bold mb-1">
            <i class="ph ph-moon text-sm"></i>
            <span>Dinner Reconciliation Policy</span>
          </div>
          <p class="text-[11px] text-slate-300 leading-relaxed">
            Reconcile room-assigned dinner covers for resident guests dining at the restaurant. Automatically increments today's restaurant dinner cover counts.
          </p>
        </div>
      {/if}

      <div class="space-y-2.5">
        {#each mealGuests as mg}
          {@const isReconciled = activeMeal === 'breakfast' ? (mg.reconciled_breakfast || mg.reconciled_today) : mg.reconciled_dinner}
          <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-md flex items-center justify-between">
            <div class="space-y-1">
              <div class="flex items-center space-x-2">
                <span class="text-sm font-bold text-white font-['Oxanium',sans-serif]">Room {mg.room_number}</span>
                <span class="text-[10px] px-2 py-0.5 rounded-full border font-semibold {isReconciled ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40' : (activeMeal === 'breakfast' ? 'bg-amber-500/20 text-amber-300 border-amber-500/40' : 'bg-purple-500/20 text-purple-300 border-purple-500/40')}">
                  {isReconciled ? `${activeMeal === 'breakfast' ? 'Breakfast' : 'Dinner'} Served` : `Pending ${activeMeal === 'breakfast' ? 'Breakfast' : 'Dinner'}`}
                </span>
              </div>
              <p class="text-xs text-slate-300 font-medium">{mg.guest_name}</p>
              <p class="text-[11px] text-slate-500">{mg.room_type} · {mg.included_covers} cover{mg.included_covers !== 1 ? 's' : ''} ({mg.card_no})</p>
            </div>

            <div>
              {#if isReconciled}
                <div class="flex items-center space-x-1 text-emerald-400 text-xs font-bold px-3 py-1.5 rounded-xl bg-emerald-500/10 border border-emerald-500/20">
                  <i class="ph ph-check-circle text-base"></i>
                  <span>Served</span>
                </div>
              {:else}
                <button
                  onclick={() => reconcileMeal(mg, activeMeal)}
                  class="px-3.5 py-2 rounded-xl {activeMeal === 'breakfast' ? 'bg-amber-500 hover:bg-amber-400 text-slate-950' : 'bg-purple-600 hover:bg-purple-500 text-white'} font-bold text-xs uppercase tracking-wider active:scale-95 transition shadow-md flex items-center space-x-1"
                >
                  <i class="ph ph-check text-sm"></i>
                  <span>Reconcile</span>
                </button>
              {/if}
            </div>
          </div>
        {/each}

        {#if mealGuests.length === 0}
          <div class="text-center py-12 text-slate-500">
            <i class="ph ph-bed text-4xl mb-3 block"></i>
            <p class="text-sm">No in-house resident guests found for reconciliation</p>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

<!-- New Docket Sheet -->
{#if showNewForm}
  <div class="fixed inset-0 z-50 flex items-end justify-center bg-slate-950/80 backdrop-blur-sm" role="dialog" aria-modal="true">
    <div class="w-full max-w-md bg-slate-900 border border-slate-700 rounded-t-3xl p-5 space-y-4 shadow-2xl max-h-[90vh] overflow-y-auto">
      <div class="flex items-center justify-between">
        <h2 class="font-bold text-white text-lg">New Docket</h2>
        <button onclick={() => showNewForm = false} class="p-2 rounded-xl bg-slate-800 text-slate-400" aria-label="Close">
          <i class="ph ph-x text-lg"></i>
        </button>
      </div>

      <!-- Table / Room toggle -->
      <div class="flex rounded-xl bg-slate-800 p-1">
        <button onclick={() => mode = 'table'} class="flex-1 py-2 rounded-lg text-xs font-bold transition {mode === 'table' ? 'bg-orange-500 text-slate-950' : 'text-slate-400'}">Table</button>
        <button onclick={() => mode = 'room'} class="flex-1 py-2 rounded-lg text-xs font-bold transition {mode === 'room' ? 'bg-orange-500 text-slate-950' : 'text-slate-400'}">Room Charge</button>
      </div>

      {#if mode === 'table'}
        <div>
          <label for="tbl-no" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Table Number</label>
          <input id="tbl-no" type="text" bind:value={tableNo} placeholder="e.g. 3" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500" />
        </div>
      {:else}
        <div>
          <label for="rm-no" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Room Number</label>
          <select id="rm-no" bind:value={roomNo} disabled={occupiedRooms.length === 0} class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500 disabled:opacity-50">
            <option value="">Select an occupied room</option>
            {#each occupiedRooms as room}
              <option value={room}>{room}</option>
            {/each}
          </select>
          {#if occupiedRooms.length === 0}
            <p class="mt-1 text-[11px] text-slate-500">No occupied rooms are currently available for charging.</p>
          {/if}
        </div>
      {/if}

      <!-- Menu item selection -->
      <div>
        <p class="text-xs text-slate-400 font-semibold uppercase tracking-wider mb-2">Add Items</p>
        <div class="grid grid-cols-2 gap-2">
          {#each menuItems as m}
            <button onclick={() => addMenuItem(m)} disabled={!m.available} class="flex flex-col items-start bg-slate-800 hover:bg-slate-700 active:scale-95 border border-slate-700 rounded-xl p-3 text-left transition disabled:opacity-40 disabled:cursor-not-allowed">
              <span class="text-xs font-semibold text-white leading-tight">{m.name}</span>
              <span class="text-orange-400 text-xs mt-1">{m.available ? `MWK ${m.price.toLocaleString()}` : 'Unavailable'}</span>
            </button>
          {/each}
        </div>
        {#if menuItems.length === 0}
          <p class="text-xs text-slate-500">No menu items are available.</p>
        {/if}
      </div>

      <!-- Selected items -->
      {#if docketItems.length > 0}
        <div class="bg-slate-800 rounded-xl p-3 space-y-2">
          {#each docketItems as item}
            <div class="flex items-center justify-between">
              <span class="text-xs text-white flex-1 truncate">{item.name}</span>
              <div class="flex items-center space-x-2">
                <button onclick={() => changeQty(item.menuItemId, -1)} class="w-6 h-6 rounded-full bg-slate-700 text-slate-200 text-xs flex items-center justify-center" aria-label="Decrease">−</button>
                <span class="text-xs text-white w-4 text-center">{item.qty}</span>
                <button onclick={() => changeQty(item.menuItemId, 1)} class="w-6 h-6 rounded-full bg-slate-700 text-slate-200 text-xs flex items-center justify-center" aria-label="Increase">+</button>
                <button onclick={() => removeItem(item.menuItemId)} class="w-6 h-6 rounded-full bg-rose-500/20 text-rose-400 text-xs flex items-center justify-center" aria-label="Remove"><i class="ph ph-x"></i></button>
              </div>
            </div>
          {/each}
          <div class="border-t border-slate-700 pt-2 flex justify-between">
            <span class="text-xs font-bold text-white">Total</span>
            <span class="text-xs font-bold text-orange-400">{fmtMwk(docketTotal)}</span>
          </div>
        </div>
      {/if}

      <div>
        <label for="dkt-notes" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Notes (optional)</label>
        <input id="dkt-notes" type="text" bind:value={notes} placeholder="e.g. No spice, extra sauce" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500" />
      </div>

      <div class="grid grid-cols-2 gap-3">
        {#if mode === 'room'}
          <button
            onclick={() => submitDocket('room_charge')}
            disabled={isSubmitting || docketItems.length === 0 || !roomNo || occupiedRooms.length === 0}
            class="py-3 rounded-xl bg-blue-500/80 hover:bg-blue-500 text-white font-bold text-xs uppercase tracking-wider disabled:opacity-50 transition active:scale-95"
          >
            <i class="ph ph-bed mr-1"></i>Charge to Room & Send
          </button>
        {:else}
          <button
            onclick={() => submitDocket('table')}
            disabled={isSubmitting || docketItems.length === 0 || !tableNo}
            class="py-3 rounded-xl bg-orange-500 hover:bg-orange-400 text-slate-950 font-bold text-xs uppercase tracking-wider disabled:opacity-50 transition active:scale-95"
          >
            <i class="ph ph-cooking-pot mr-1"></i>Submit to Kitchen
          </button>
          <button
            onclick={() => showNewForm = false}
            class="py-3 rounded-xl bg-slate-800 text-slate-300 font-semibold text-xs border border-slate-700 transition"
          >
            Cancel
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}
