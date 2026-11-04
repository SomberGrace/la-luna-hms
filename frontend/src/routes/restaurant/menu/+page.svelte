<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { getRestaurantItemImage } from '#lib/productImages';
  import { onMount } from 'svelte';

  type MenuItem = {
    id: string;
    name: string;
    category: 'Breakfast' | 'Lunch' | 'Dinner' | 'Snack' | 'Beverage';
    price: number;
    costPrice: number;
    available: boolean;
    isAvailable: boolean;
    image_url?: string;
  };
  type MenuItemResponse = {
    id: string;
    name: string;
    category: MenuItem['category'];
    price: number;
    costPrice: number;
    available: boolean;
    is_available: boolean;
    image_url: string;
  };

  const categories = ['All', 'Breakfast', 'Lunch', 'Dinner', 'Snack', 'Beverage'] as const;
  let activeCategory = $state<string>('All');
  let showAddModal = $state(false);
  let isSaving = $state(false);
  let saveMsg = $state('');
  let loadingItems = $state(true);
  let updatingItem = $state('');
  let loadError = $state('');

  let items = $state<MenuItem[]>([]);

  // New item form state
  let newName = $state('');
  let newCategory = $state<MenuItem['category']>('Lunch');
  let newPrice = $state('');
  let newCostPrice = $state('');
  let newImageUrl = $state('');
  let newAvailable = $state(true);

  onMount(() => {
    void loadItems();
  });

  async function loadItems() {
    loadingItems = true;
    try {
      const res = await apiFetch('/api/v1/restaurant/menu');
      if (!res.ok) throw new Error('Unable to load menu items.');
      const data: MenuItemResponse[] = await res.json();
      items = data.map((item) => ({
        id: item.id,
        name: item.name,
        category: item.category,
        price: item.price,
        costPrice: item.costPrice,
        available: item.available,
        isAvailable: item.is_available,
        image_url: item.image_url
      }));
      loadError = '';
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to load menu items.';
    } finally {
      loadingItems = false;
    }
  }

  let filtered = $derived(
    activeCategory === 'All' ? items : items.filter(i => i.category === activeCategory)
  );

  function catBadge(cat: string) {
    const m: Record<string, string> = {
      Breakfast: 'bg-amber-500/20 text-amber-400',
      Lunch:     'bg-blue-500/20 text-blue-400',
      Dinner:    'bg-purple-500/20 text-purple-400',
      Snack:     'bg-green-500/20 text-green-400',
      Beverage:  'bg-cyan-500/20 text-cyan-400',
    };
    return m[cat] ?? 'bg-slate-700 text-slate-300';
  }

  async function toggleAvailable(item: MenuItem) {
    updatingItem = item.id;
    try {
      const res = await apiFetch('/api/v1/restaurant/menu/availability', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: item.id, available: !item.isAvailable })
      });
      if (!res.ok) throw new Error('Unable to update menu availability.');
      await loadItems();
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to update menu availability.';
    } finally {
      updatingItem = '';
    }
  }

  async function addItem() {
    if (!newName.trim() || !newPrice) return;
    isSaving = true;
    saveMsg = '';
    try {
      const response = await apiFetch('/api/v1/restaurant/menu', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: newName.trim(),
          category: newCategory,
          price: Number(newPrice),
          costPrice: Number(newCostPrice) || 0,
          available: newAvailable,
          image_url: newImageUrl.trim()
        })
      });
      if (!response.ok) throw new Error('Unable to add menu item.');
      await loadItems();
      newName = ''; newPrice = ''; newCostPrice = ''; newImageUrl = ''; newAvailable = true;
      showAddModal = false;
      saveMsg = 'Item added.';
      setTimeout(() => saveMsg = '', 2000);
    } catch (error) {
      loadError = error instanceof Error ? error.message : 'Unable to add menu item.';
    } finally {
      isSaving = false;
    }
  }

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }
</script>

<div class="space-y-4 pb-4">
  <div class="flex items-center justify-between pt-1">
    <h1 class="text-xl font-bold text-white">Menu Items</h1>
    <span class="text-xs text-slate-400">{items.length} items</span>
  </div>

  {#if saveMsg}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center">{saveMsg}</div>
  {/if}
  {#if loadError}
    <div role="alert" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs text-center">{loadError}</div>
  {/if}

  <!-- Category Filter Chips -->
  <div class="flex space-x-2 overflow-x-auto pb-1 scrollbar-none">
    {#each categories as cat}
      <button
        onclick={() => activeCategory = cat}
        class="flex-shrink-0 px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all {activeCategory === cat ? 'bg-orange-500 text-slate-950 border-orange-500 shadow-lg shadow-orange-500/25' : 'bg-slate-900 text-slate-400 border-slate-700 hover:border-slate-500'}"
      >
        {cat}
      </button>
    {/each}
  </div>

  <!-- Items List -->
  <div class="space-y-2.5">
    {#if loadingItems}
      <div class="text-center py-12 text-slate-500 text-sm">Loading menu…</div>
    {/if}
    {#each filtered as item}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden {!item.available ? 'opacity-50' : ''}">
        <div class="flex">
          <div class="w-20 h-20 flex-shrink-0 bg-slate-800">
            <img src={getRestaurantItemImage(item.name, item.category, item.image_url)} alt={item.name} class="w-full h-full object-cover" loading="lazy" />
          </div>
          <div class="flex-1 p-3.5 flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center space-x-2 mb-1">
                <span class="font-semibold text-white text-sm">{item.name}</span>
                <span class="text-[10px] px-1.5 py-0.5 rounded-full font-semibold {catBadge(item.category)}">{item.category}</span>
              </div>
              <p class="text-xs text-slate-400">Cost: {fmtMwk(item.costPrice)} · Margin: {fmtMwk(item.price - item.costPrice)}</p>
            </div>
            <div class="flex flex-col items-end space-y-2">
              <span class="text-orange-400 font-bold text-sm font-['Oxanium',sans-serif]">{fmtMwk(item.price)}</span>
              <button
                onclick={() => toggleAvailable(item)}
                disabled={updatingItem === item.id}
                aria-label={`Toggle ${item.name} menu availability`}
                class="text-[10px] px-2.5 py-1 rounded-full font-semibold border transition-all disabled:opacity-50 {item.isAvailable ? 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30' : 'bg-slate-700 text-slate-400 border-slate-600'}"
              >
                {item.isAvailable ? (item.available ? '● Available' : '● Low Stock') : '○ Off Menu'}
              </button>
            </div>
          </div>
        </div>
      </div>
    {/each}
    {#if !loadingItems && filtered.length === 0 && !loadError}
      <div class="text-center py-12 text-slate-500">
        <i class="ph ph-fork-knife text-4xl mb-3 block"></i>
        <p class="text-sm">No items in this category</p>
      </div>
    {/if}
  </div>
</div>

<!-- Add Item FAB -->
<button
  onclick={() => showAddModal = true}
  class="fixed bottom-20 right-4 w-14 h-14 rounded-full bg-orange-500 hover:bg-orange-400 text-slate-950 shadow-xl shadow-orange-500/40 flex items-center justify-center active:scale-95 transition-transform z-30"
  aria-label="Add menu item"
>
  <i class="ph ph-plus text-2xl font-bold"></i>
</button>

<!-- Add Item Modal -->
{#if showAddModal}
  <div class="fixed inset-0 z-50 flex items-end justify-center bg-slate-950/80 backdrop-blur-sm" role="dialog" aria-modal="true">
    <div class="w-full max-w-md bg-slate-900 border border-slate-700 rounded-t-3xl p-6 space-y-4 shadow-2xl">
      <div class="flex items-center justify-between">
        <h2 class="font-bold text-white text-lg">Add Menu Item</h2>
        <button onclick={() => showAddModal = false} class="p-2 rounded-xl bg-slate-800 text-slate-400 hover:text-white" aria-label="Close modal">
          <i class="ph ph-x text-lg"></i>
        </button>
      </div>

      <div class="space-y-3">
        <div>
          <label for="item-name" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Item Name</label>
          <input id="item-name" type="text" bind:value={newName} placeholder="e.g. Grilled Chicken" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500" />
        </div>
        <div>
          <label for="item-cat" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Category</label>
          <select id="item-cat" bind:value={newCategory} class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500">
            {#each ['Breakfast','Lunch','Dinner','Snack','Beverage'] as c}
              <option value={c}>{c}</option>
            {/each}
          </select>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="item-price" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Selling Price (MWK)</label>
            <input id="item-price" type="number" bind:value={newPrice} placeholder="0" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500" />
          </div>
          <div>
            <label for="item-cost" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Cost Price (MWK)</label>
            <input id="item-cost" type="number" bind:value={newCostPrice} placeholder="0" class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500" />
          </div>
        </div>
        <div>
          <label for="item-img" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Image URL (Optional)</label>
          <input id="item-img" type="url" bind:value={newImageUrl} placeholder="https://images.unsplash.com/..." class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-orange-500" />
        </div>
        <div class="flex items-center space-x-3">
          <button
            type="button"
            onclick={() => newAvailable = !newAvailable}
            class="w-12 h-6 rounded-full transition-colors {newAvailable ? 'bg-emerald-500' : 'bg-slate-700'} relative"
            role="switch"
            aria-checked={newAvailable}
            aria-label="Toggle availability"
          >
            <span class="absolute top-0.5 w-5 h-5 bg-white rounded-full shadow transition-all {newAvailable ? 'right-0.5' : 'left-0.5'}"></span>
          </button>
          <span class="text-sm text-slate-300">{newAvailable ? 'Available on menu' : 'Off menu'}</span>
        </div>
      </div>

      <button
        onclick={addItem}
        disabled={isSaving || !newName.trim() || !newPrice}
        class="w-full py-3.5 rounded-xl bg-orange-500 hover:bg-orange-400 text-slate-950 font-bold text-sm uppercase tracking-wider shadow-lg active:scale-98 transition disabled:opacity-50"
      >
        {isSaving ? 'Saving...' : 'Add to Menu'}
      </button>
    </div>
  </div>
{/if}
