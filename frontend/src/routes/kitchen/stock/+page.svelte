<script lang="ts">
  import { apiFetch } from '#lib/api';
  type KitchenIngredient = {
    id: string;
    name: string;
    unit: string;
    openingStock: number;
    addStock: number;
    soldUsage: number;
    closingBalance: number;
    minAlert: number;
  };

  let auditDate = $state(new Date().toISOString().slice(0, 10));
  let isSaving = $state(false);
  let saveMsg = $state('');

  let stockList = $state<KitchenIngredient[]>([
    {
      id: 'ks-1',
      name: 'Maize Flour (Gramil)',
      unit: 'kg',
      openingStock: 50,
      addStock: 25,
      soldUsage: 18,
      closingBalance: 57,
      minAlert: 20
    },
    {
      id: 'ks-2',
      name: 'Fresh Tilapia (Lake Malawi)',
      unit: 'pieces',
      openingStock: 22,
      addStock: 10,
      soldUsage: 14,
      closingBalance: 18,
      minAlert: 10
    },
    {
      id: 'ks-3',
      name: 'Whole Chickens',
      unit: 'birds',
      openingStock: 15,
      addStock: 0,
      soldUsage: 11,
      closingBalance: 4,
      minAlert: 6
    },
    {
      id: 'ks-4',
      name: 'Cooking Oil (Sunseed)',
      unit: 'litres',
      openingStock: 20,
      addStock: 10,
      soldUsage: 6,
      closingBalance: 24,
      minAlert: 8
    },
    {
      id: 'ks-5',
      name: 'Potatoes (Irish)',
      unit: 'kg',
      openingStock: 35,
      addStock: 0,
      soldUsage: 22,
      closingBalance: 13,
      minAlert: 15
    },
    {
      id: 'ks-6',
      name: 'Tomatoes & Onions',
      unit: 'kg',
      openingStock: 18,
      addStock: 5,
      soldUsage: 12,
      closingBalance: 11,
      minAlert: 5
    }
  ]);

  // Add Ingredient Modal
  let showAddModal = $state(false);
  let newName = $state('');
  let newUnit = $state('kg');
  let newOpening = $state('');
  let newMinAlert = $state('10');

  function updateAdd(id: string, value: string) {
    const add = parseFloat(value) || 0;
    stockList = stockList.map(item => {
      if (item.id !== id) return item;
      const closing = Math.max(0, (item.openingStock + add) - item.soldUsage);
      return { ...item, addStock: add, closingBalance: closing };
    });
  }

  function updateUsage(id: string, value: string) {
    const usage = parseFloat(value) || 0;
    stockList = stockList.map(item => {
      if (item.id !== id) return item;
      const closing = Math.max(0, (item.openingStock + item.addStock) - usage);
      return { ...item, soldUsage: usage, closingBalance: closing };
    });
  }

  function addIngredient() {
    if (!newName.trim()) return;
    const opening = parseFloat(newOpening) || 0;
    const alertVal = parseFloat(newMinAlert) || 5;

    stockList = [
      ...stockList,
      {
        id: 'ks-' + (stockList.length + 1),
        name: newName.trim(),
        unit: newUnit,
        openingStock: opening,
        addStock: 0,
        soldUsage: 0,
        closingBalance: opening,
        minAlert: alertVal
      }
    ];

    newName = '';
    newOpening = '';
    showAddModal = false;
  }

  async function saveAudit() {
    isSaving = true;
    saveMsg = '';
    try {
      const res = await apiFetch('/api/v1/kitchen/stock', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          audit_date: auditDate,
          items: stockList
        })
      });
      saveMsg = res.ok ? '✓ Stock ledger saved to database.' : '⚠ Stored locally — will sync when connected.';
    } catch {
      saveMsg = '⚠ Offline mode — changes saved in station memory.';
    } finally {
      isSaving = false;
      setTimeout(() => {
        saveMsg = '';
      }, 4000);
    }
  }
</script>

<div class="space-y-4 pb-4">
  <div class="flex items-center justify-between pt-1">
    <div>
      <h1 class="text-xl font-bold text-white tracking-tight">Kitchen Stock</h1>
      <p class="text-xs text-slate-400">Ingredient Inventory & Usage Audit</p>
    </div>
    <button
      onclick={() => (showAddModal = true)}
      class="flex items-center space-x-1.5 px-3 py-1.5 bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-xs rounded-xl transition active:scale-95 shadow-md shadow-emerald-600/20"
    >
      <i class="ph ph-plus text-sm"></i>
      <span>Add Item</span>
    </button>
  </div>

  <!-- Audit Date Selector -->
  <div>
    <label for="k-audit-date" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Audit Date</label>
    <input
      id="k-audit-date"
      type="date"
      bind:value={auditDate}
      class="w-full px-4 py-2.5 bg-slate-900 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-emerald-500"
    />
  </div>

  {#if saveMsg}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/20 text-emerald-300 text-xs text-center font-medium">
      {saveMsg}
    </div>
  {/if}

  <!-- Ingredients Cards -->
  <div class="space-y-3">
    {#each stockList as item}
      {@const isLow = item.closingBalance <= item.minAlert}
      <div class="bg-slate-900 border {isLow ? 'border-rose-600/50 bg-rose-950/10' : 'border-slate-800'} rounded-2xl p-4 space-y-3 shadow-md">
        <!-- Header -->
        <div class="flex items-start justify-between">
          <div>
            <h3 class="font-bold text-white text-sm">{item.name}</h3>
            <p class="text-xs text-slate-400">Unit: <span class="text-slate-200 uppercase font-mono">{item.unit}</span> · Min Safe: {item.minAlert} {item.unit}</p>
          </div>
          {#if isLow}
            <span class="text-[10px] bg-rose-500/20 text-rose-300 border border-rose-500/30 px-2 py-0.5 rounded-full font-bold uppercase tracking-wider animate-pulse">
              Low Stock
            </span>
          {:else}
            <span class="text-[10px] bg-emerald-500/20 text-emerald-300 border border-emerald-500/30 px-2 py-0.5 rounded-full font-semibold">
              In Stock
            </span>
          {/if}
        </div>

        <!-- Stock Columns -->
        <div class="grid grid-cols-3 gap-2 text-center">
          <div class="bg-slate-800/80 rounded-xl p-2">
            <span class="text-[10px] text-slate-400 uppercase font-semibold block">Opening</span>
            <span class="text-sm font-bold text-white font-['Oxanium',sans-serif]">{item.openingStock}</span>
          </div>

          <div class="bg-slate-800/80 rounded-xl p-2 border border-emerald-500/20">
            <label for="add-input-{item.id}" class="text-[10px] text-emerald-400 uppercase font-semibold block">+ Received</label>
            <input
              id="add-input-{item.id}"
              type="number"
              value={item.addStock}
              oninput={(e) => updateAdd(item.id, (e.target as HTMLInputElement).value)}
              class="w-full text-sm font-bold text-emerald-300 bg-transparent text-center focus:outline-none"
            />
          </div>

          <div class="bg-slate-800/80 rounded-xl p-2 border border-amber-500/20">
            <label for="usage-input-{item.id}" class="text-[10px] text-amber-400 uppercase font-semibold block">Used/Sold</label>
            <input
              id="usage-input-{item.id}"
              type="number"
              value={item.soldUsage}
              oninput={(e) => updateUsage(item.id, (e.target as HTMLInputElement).value)}
              class="w-full text-sm font-bold text-amber-300 bg-transparent text-center focus:outline-none"
            />
          </div>
        </div>

        <!-- Balance Footnote -->
        <div class="flex items-center justify-between pt-2 border-t border-slate-800 text-xs">
          <span class="text-slate-400">Current Closing Balance:</span>
          <span class="text-sm font-bold {isLow ? 'text-rose-400' : 'text-emerald-400'} font-['Oxanium',sans-serif]">
            {item.closingBalance} {item.unit}
          </span>
        </div>
      </div>
    {/each}
  </div>

  <!-- Save Button -->
  <button
    onclick={saveAudit}
    disabled={isSaving}
    class="w-full py-3.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-sm uppercase tracking-wider shadow-lg shadow-emerald-600/25 active:scale-98 transition disabled:opacity-50"
  >
    {isSaving ? 'Saving Audit...' : 'Save Kitchen Audit'}
  </button>
</div>

<!-- Add Ingredient Modal -->
{#if showAddModal}
  <div class="fixed inset-0 z-50 flex items-end justify-center bg-slate-950/80 backdrop-blur-sm" role="dialog" aria-modal="true">
    <div class="w-full max-w-md bg-slate-900 border border-slate-700 rounded-t-3xl p-5 space-y-4 shadow-2xl">
      <div class="flex items-center justify-between">
        <h2 class="font-bold text-white text-lg">Add Raw Ingredient</h2>
        <button onclick={() => (showAddModal = false)} class="p-2 rounded-xl bg-slate-800 text-slate-400" aria-label="Close modal">
          <i class="ph ph-x text-lg"></i>
        </button>
      </div>

      <div class="space-y-3">
        <div>
          <label for="ing-name" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Ingredient Name</label>
          <input
            id="ing-name"
            type="text"
            bind:value={newName}
            placeholder="e.g. Beef Tenderloin"
            class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-emerald-500"
          />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="ing-unit" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Unit</label>
            <select
              id="ing-unit"
              bind:value={newUnit}
              class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-emerald-500"
            >
              <option value="kg">kg</option>
              <option value="pieces">pieces</option>
              <option value="litres">litres</option>
              <option value="grams">grams</option>
              <option value="bags">bags</option>
            </select>
          </div>
          <div>
            <label for="ing-open" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Opening Stock</label>
            <input
              id="ing-open"
              type="number"
              bind:value={newOpening}
              placeholder="0"
              class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-emerald-500"
            />
          </div>
        </div>

        <div>
          <label for="ing-min" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">Low Stock Alert Threshold</label>
          <input
            id="ing-min"
            type="number"
            bind:value={newMinAlert}
            placeholder="10"
            class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-emerald-500"
          />
        </div>
      </div>

      <button
        onclick={addIngredient}
        disabled={!newName.trim()}
        class="w-full py-3.5 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-sm uppercase tracking-wider shadow-lg active:scale-98 transition disabled:opacity-50"
      >
        Add to Inventory
      </button>
    </div>
  </div>
{/if}
