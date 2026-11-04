<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { getBarItemImage } from '#lib/productImages';
  import { onMount } from 'svelte';
  type BarItem = {
    id: string;
    name: string;
    unitPrice: number;
    costPrice: number;
    stock: number;
    openingStock: number;
    addStock: number;
    soldQty: number;
    closingStock: number;
    salesRevenue: number;
  };

  let auditDate = $state(new Date().toISOString().slice(0, 10));
  let isSaving = $state(false);
  let isLoading = $state(true);
  let saveMsg = $state('');
  let errorMessage = $state('');
  let items = $state<BarItem[]>([]);

  let showAddModal = $state(false);
  let newName = $state('');
  let newUnitPrice = $state('');
  let newCostPrice = $state('');
  let newOpeningStock = $state('');

  onMount(() => {
    void loadLedger();
  });

  async function loadLedger() {
    isLoading = true;
    errorMessage = '';
    try {
      const response = await apiFetch(`/api/v1/bar/ledger?date=${auditDate}`);
      if (!response.ok) throw new Error('Unable to load bar stock ledger.');
      items = await response.json();
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : 'Unable to load bar stock ledger.';
    } finally {
      isLoading = false;
    }
  }

  function isLow(item: BarItem) { return item.closingStock < 5; }

  function updateSold(id: string, val: string) {
    const sold = parseFloat(val) || 0;
    items = items.map(i => {
      if (i.id !== id) return i;
      const closing = Math.max(0, (i.openingStock + i.addStock) - sold);
      return { ...i, soldQty: sold, closingStock: closing, salesRevenue: sold * i.unitPrice };
    });
  }

  function updateAddStock(id: string, val: string) {
    const add = parseFloat(val) || 0;
    items = items.map(i => {
      if (i.id !== id) return i;
      const closing = Math.max(0, (i.openingStock + add) - i.soldQty);
      return { ...i, addStock: add, closingStock: closing };
    });
  }

  async function saveAudit() {
    isSaving = true;
    saveMsg = '';
    errorMessage = '';
    try {
      const r = await apiFetch('/api/v1/bar/ledger', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ audit_date: auditDate, items }),
      });
      if (!r.ok) throw new Error('Unable to save the bar stock audit.');
      saveMsg = 'Audit saved to the database.';
      await loadLedger();
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : 'Unable to save the bar stock audit.';
    }
    isSaving = false;
    setTimeout(() => saveMsg = '', 4000);
  }

  function addBarItem() {
    if (!newName.trim()) return;
    const opening = parseFloat(newOpeningStock) || 0;
    items = [...items, {
      id: Date.now().toString(),
      name: newName.trim(),
      unitPrice: parseFloat(newUnitPrice) || 0,
      costPrice: parseFloat(newCostPrice) || 0,
      stock: opening,
      openingStock: opening,
      addStock: 0,
      soldQty: 0,
      closingStock: opening,
      salesRevenue: 0,
    }];
    newName = ''; newUnitPrice = ''; newCostPrice = ''; newOpeningStock = '';
    showAddModal = false;
  }

  function fmtMwk(n: number) { return 'MWK ' + n.toLocaleString(); }
</script>

<div class="space-y-4 pb-4">
  <div class="flex items-center justify-between pt-1">
    <h1 class="text-xl font-bold text-white">Stock Audit</h1>
    <button
      onclick={() => showAddModal = true}
      class="flex items-center space-x-1.5 px-3 py-1.5 bg-purple-600 text-white text-xs font-bold rounded-xl active:scale-95 transition"
    >
      <i class="ph ph-plus text-sm"></i> Add Item
    </button>
  </div>

  <!-- Date -->
  <div>
    <label for="audit-date" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Audit Date</label>
    <input id="audit-date" type="date" bind:value={auditDate} class="w-full px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-purple-500" />
  </div>

  {#if saveMsg}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center">{saveMsg}</div>
  {/if}
  {#if errorMessage}
    <div role="alert" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">
      {errorMessage}
      <button onclick={loadLedger} class="ml-2 underline">Retry</button>
    </div>
  {/if}

  <!-- Items -->
  <div class="space-y-3">
    {#if isLoading}
      <p class="py-8 text-center text-sm text-slate-400">Loading live bar inventory…</p>
    {/if}
    {#each items as item}
      <div class="bg-slate-900 border {isLow(item) ? 'border-rose-700/60' : 'border-slate-800'} rounded-2xl p-4 space-y-3">
        <div class="flex items-center justify-between gap-3">
          <div class="flex min-w-0 items-center gap-3">
            <img src={getBarItemImage(item.name)} alt="" class="h-14 w-14 flex-shrink-0 rounded-xl object-cover" loading="lazy" />
            <div class="min-w-0">
              <p class="font-semibold text-white text-sm">{item.name}</p>
              <p class="text-xs text-slate-400">Unit: {fmtMwk(item.unitPrice)} · Cost: {fmtMwk(item.costPrice)}</p>
            </div>
          </div>
          {#if isLow(item)}
            <span class="text-[10px] bg-rose-500/20 text-rose-400 border border-rose-500/30 px-2 py-0.5 rounded-full font-semibold">Low Stock</span>
          {/if}
        </div>
        <div class="grid grid-cols-3 gap-2 text-center">
          <div class="bg-slate-800 rounded-xl p-2">
            <p class="text-[10px] text-slate-400 mb-0.5">Opening</p>
            <p class="text-sm font-bold text-white">{item.openingStock}</p>
          </div>
          <div class="bg-slate-800 rounded-xl p-2">
            <p class="text-[10px] text-slate-400 mb-0.5">+ Added</p>
            <input
              type="number"
              value={item.addStock}
              oninput={(e) => updateAddStock(item.id, (e.target as HTMLInputElement).value)}
              class="w-full text-sm font-bold text-purple-300 bg-transparent text-center focus:outline-none"
            />
          </div>
          <div class="bg-slate-800 rounded-xl p-2">
            <p class="text-[10px] text-slate-400 mb-0.5">Total</p>
            <p class="text-sm font-bold text-white">{item.openingStock + item.addStock}</p>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-2 text-center">
          <div class="bg-slate-800 rounded-xl p-2">
            <p class="text-[10px] text-slate-400 mb-0.5">Sold</p>
            <input
              type="number"
              value={item.soldQty}
              oninput={(e) => updateSold(item.id, (e.target as HTMLInputElement).value)}
              class="w-full text-sm font-bold text-amber-400 bg-transparent text-center focus:outline-none"
            />
          </div>
          <div class="bg-slate-800 rounded-xl p-2">
            <p class="text-[10px] text-slate-400 mb-0.5">Closing</p>
            <p class="text-sm font-bold {isLow(item) ? 'text-rose-400' : 'text-emerald-400'}">{item.closingStock}</p>
          </div>
        </div>
        <div class="flex justify-between items-center pt-1 border-t border-slate-800">
          <span class="text-xs text-slate-400">Sales Revenue</span>
          <span class="text-sm font-bold text-emerald-400">{fmtMwk(item.salesRevenue)}</span>
        </div>
      </div>
    {/each}
    {#if !isLoading && !errorMessage && items.length === 0}
      <p class="py-8 text-center text-sm text-slate-400">No bar inventory records found.</p>
    {/if}
  </div>

  <button
    onclick={saveAudit}
    disabled={isSaving}
    class="w-full py-3.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-bold text-sm uppercase tracking-wider shadow-lg active:scale-98 transition disabled:opacity-50"
  >
    {isSaving ? 'Saving...' : 'Save Audit'}
  </button>
</div>

<!-- Add Item Modal -->
{#if showAddModal}
  <div class="fixed inset-0 z-50 flex items-end justify-center bg-slate-950/80 backdrop-blur-sm" role="dialog" aria-modal="true">
    <div class="w-full max-w-md bg-slate-900 border border-slate-700 rounded-t-3xl p-5 space-y-4 shadow-2xl">
      <div class="flex items-center justify-between">
        <h2 class="font-bold text-white text-lg">Add Bar Item</h2>
        <button onclick={() => showAddModal = false} class="p-2 rounded-xl bg-slate-800 text-slate-400" aria-label="Close"><i class="ph ph-x text-lg"></i></button>
      </div>
      <div class="space-y-3">
        <div>
          <label for="bi-name" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Item Name</label>
          <input id="bi-name" type="text" bind:value={newName} placeholder="e.g. Savanna Dry" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-purple-500" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="bi-unit" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Unit Price (MWK)</label>
            <input id="bi-unit" type="number" bind:value={newUnitPrice} placeholder="0" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-purple-500" />
          </div>
          <div>
            <label for="bi-cost" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Cost Price (MWK)</label>
            <input id="bi-cost" type="number" bind:value={newCostPrice} placeholder="0" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-purple-500" />
          </div>
        </div>
        <div>
          <label for="bi-open" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Opening Stock (units)</label>
          <input id="bi-open" type="number" bind:value={newOpeningStock} placeholder="0" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-purple-500" />
        </div>
      </div>
      <button onclick={addBarItem} disabled={!newName.trim()} class="w-full py-3.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-bold text-sm uppercase tracking-wider disabled:opacity-50 transition">
        Add Item
      </button>
    </div>
  </div>
{/if}
