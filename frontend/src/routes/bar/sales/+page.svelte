<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { getBarItemImage } from '#lib/productImages';
  import { onMount } from 'svelte';
  type SaleEntry = { id: string; item: string; qty: number; unitPrice: number; total: number; time: string };
  type BarItem = { id: string; name: string; unitPrice: number; costPrice: number; stock: number };

  let barItems = $state<BarItem[]>([]);
  let todaySales = $state<SaleEntry[]>([]);
  let selectedItemId = $state('');
  let qty = $state(1);
  let isPosting = $state(false);
  let isLoading = $state(true);
  let errorMessage = $state('');
  let toast = $state('');

  let selectedItem = $derived(barItems.find(i => i.id === selectedItemId));
  let saleTotal = $derived((selectedItem?.unitPrice ?? 0) * qty);
  let grandTotal = $derived(todaySales.reduce((s, e) => s + e.total, 0));

  onMount(() => {
    void loadItems();
  });

  async function loadItems() {
    isLoading = true;
    errorMessage = '';
    try {
      const response = await apiFetch('/api/v1/bar/items');
      if (!response.ok) throw new Error('Unable to load bar inventory.');
      barItems = await response.json();
      if (!barItems.some(item => item.id === selectedItemId)) {
        selectedItemId = barItems[0]?.id ?? '';
      }
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : 'Unable to load bar inventory.';
    } finally {
      isLoading = false;
    }
  }

  async function recordSale() {
    if (!selectedItem || qty < 1 || qty > selectedItem.stock) return;
    isPosting = true;
    try {
      const response = await apiFetch('/api/v1/bar/sales', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ item: selectedItem.name, qty, unitPrice: selectedItem.unitPrice, total: saleTotal })
      });
      if (!response.ok) throw new Error('Unable to record sale. Check the item and stock, then retry.');
      const now = new Date().toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' });
      const entry: SaleEntry = { id: Date.now().toString(), item: selectedItem.name, qty, unitPrice: selectedItem.unitPrice, total: saleTotal, time: now };
      todaySales = [entry, ...todaySales];
      await loadItems();
      qty = 1;
      toast = `Recorded: ${entry.item} ×${entry.qty}`;
      setTimeout(() => toast = '', 3000);
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : 'Unable to record sale.';
    } finally {
      isPosting = false;
    }
  }

  function fmtMwk(n: number) { return 'MWK ' + n.toLocaleString(); }
</script>

<div class="space-y-5 pb-4">
  <h1 class="text-xl font-bold text-white pt-1">Record Sales</h1>

  {#if toast}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center">{toast}</div>
  {/if}
  {#if errorMessage}
    <div role="alert" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">
      {errorMessage}
      <button onclick={loadItems} class="ml-2 underline">Retry</button>
    </div>
  {/if}

  <!-- Sale Entry Form -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-4">
    <h2 class="font-semibold text-white text-sm">New Sale Entry</h2>

    {#if isLoading}
      <p class="text-sm text-slate-400">Loading live bar inventory…</p>
    {:else if barItems.length === 0}
      <p class="text-sm text-slate-400">No bar items are available in inventory.</p>
    {:else}
    <div>
      <label for="sale-item" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Item</label>
      {#if selectedItem}
        <div class="mb-2 flex items-center gap-3 rounded-xl border border-slate-700 bg-slate-950/70 p-2">
          <img src={getBarItemImage(selectedItem.name)} alt="" class="h-14 w-14 rounded-lg object-cover" />
          <div>
            <p class="text-sm font-semibold text-white">{selectedItem.name}</p>
            <p class="text-xs text-slate-400">{fmtMwk(selectedItem.unitPrice)} per item</p>
          </div>
        </div>
      {/if}
      <select id="sale-item" bind:value={selectedItemId} class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-purple-500">
        {#each barItems as item}
          <option value={item.id} disabled={item.stock <= 0}>{item.name} — {fmtMwk(item.unitPrice)} ({item.stock} in stock)</option>
        {/each}
      </select>
    </div>

    <div>
      <label for="sale-qty" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Quantity</label>
      <div class="flex items-center space-x-3">
        <button onclick={() => qty = Math.max(1, qty - 1)} class="w-12 h-12 rounded-xl bg-slate-800 text-slate-200 text-xl font-bold border border-slate-700 active:scale-95" aria-label="Decrease qty">−</button>
        <input id="sale-qty" type="number" bind:value={qty} min="1" class="flex-1 text-center px-4 py-3 bg-slate-950 border border-slate-700 rounded-xl text-white text-lg font-bold focus:outline-none focus:border-purple-500" />
        <button onclick={() => qty += 1} class="w-12 h-12 rounded-xl bg-slate-800 text-slate-200 text-xl font-bold border border-slate-700 active:scale-95" aria-label="Increase qty">+</button>
      </div>
    </div>

    <div class="bg-slate-800 rounded-xl p-3 flex justify-between items-center">
      <span class="text-sm text-slate-400">Sale Total</span>
      <span class="text-xl font-bold text-purple-300 font-['Oxanium',sans-serif]">{fmtMwk(saleTotal)}</span>
    </div>

    <button
      onclick={recordSale}
      disabled={isPosting || qty < 1 || !selectedItem || qty > (selectedItem?.stock ?? 0)}
      class="w-full py-3.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-bold text-sm uppercase tracking-wider shadow-lg active:scale-98 transition disabled:opacity-50"
    >
      {isPosting ? 'Recording...' : 'Record Sale'}
    </button>
    {/if}
  </div>

  <!-- Today's sales log -->
  <div>
    <div class="flex justify-between items-center mb-3">
      <h2 class="font-bold text-white text-sm">Sales Recorded This Session</h2>
      <span class="text-xs text-emerald-400 font-semibold">{fmtMwk(grandTotal)}</span>
    </div>
    <div class="space-y-2">
      {#each todaySales as s}
        <div class="bg-slate-900 border border-slate-800 rounded-xl px-4 py-3 flex items-center justify-between">
          <div class="flex items-center gap-3">
            <img src={getBarItemImage(s.item)} alt="" class="h-11 w-11 rounded-lg object-cover" loading="lazy" />
            <div>
            <p class="text-sm text-white font-semibold">{s.item}</p>
            <p class="text-xs text-slate-400">{s.qty} unit{s.qty !== 1 ? 's' : ''} · {s.time}</p>
            </div>
          </div>
          <span class="text-sm font-bold text-white">{fmtMwk(s.total)}</span>
        </div>
      {/each}
      {#if todaySales.length === 0}
        <div class="text-center py-8 text-slate-500">
          <i class="ph ph-wine text-3xl mb-2 block"></i>
          <p class="text-xs">No sales recorded in this session.</p>
        </div>
      {/if}
    </div>
  </div>
</div>
