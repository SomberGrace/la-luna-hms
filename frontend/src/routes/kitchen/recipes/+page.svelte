<script lang="ts">
  type IngredientReq = {
    name: string;
    qty: number;
    unit: string;
    approxCostMWK: number;
  };

  type DishRecipe = {
    dishId: string;
    dishName: string;
    category: 'Breakfast' | 'Lunch' | 'Dinner' | 'Snack';
    sellingPrice: number;
    ingredients: IngredientReq[];
  };

  let recipes = $state<DishRecipe[]>([
    {
      dishId: 'd-1',
      dishName: 'Nsima with Stewed Chicken',
      category: 'Lunch',
      sellingPrice: 3500,
      ingredients: [
        { name: 'Maize Flour (Gramil)', qty: 0.35, unit: 'kg', approxCostMWK: 350 },
        { name: 'Whole Chicken (Portion)', qty: 0.25, unit: 'bird', approxCostMWK: 1200 },
        { name: 'Tomatoes & Onions', qty: 0.15, unit: 'kg', approxCostMWK: 150 },
        { name: 'Cooking Oil', qty: 0.05, unit: 'litres', approxCostMWK: 150 }
      ]
    },
    {
      dishId: 'd-2',
      dishName: 'Grilled Lake Tilapia',
      category: 'Dinner',
      sellingPrice: 4500,
      ingredients: [
        { name: 'Fresh Tilapia', qty: 1, unit: 'piece', approxCostMWK: 1800 },
        { name: 'Cooking Oil & Spices', qty: 0.08, unit: 'litres', approxCostMWK: 250 },
        { name: 'Salad / Lemon', qty: 0.1, unit: 'kg', approxCostMWK: 150 }
      ]
    },
    {
      dishId: 'd-3',
      dishName: 'Chips & Sauce',
      category: 'Snack',
      sellingPrice: 1200,
      ingredients: [
        { name: 'Irish Potatoes', qty: 0.45, unit: 'kg', approxCostMWK: 450 },
        { name: 'Cooking Oil', qty: 0.1, unit: 'litres', approxCostMWK: 250 },
        { name: 'Tomato Sauce', qty: 0.05, unit: 'litres', approxCostMWK: 80 }
      ]
    },
    {
      dishId: 'd-4',
      dishName: 'Beef Steak with Mash',
      category: 'Dinner',
      sellingPrice: 8500,
      ingredients: [
        { name: 'Beef Tenderloin', qty: 0.35, unit: 'kg', approxCostMWK: 2800 },
        { name: 'Irish Potatoes', qty: 0.3, unit: 'kg', approxCostMWK: 300 },
        { name: 'Butter & Herbs', qty: 0.05, unit: 'kg', approxCostMWK: 400 }
      ]
    }
  ]);

  let selectedCategory = $state<'All' | 'Breakfast' | 'Lunch' | 'Dinner' | 'Snack'>('All');

  let filteredRecipes = $derived(
    selectedCategory === 'All'
      ? recipes
      : recipes.filter(r => r.category === selectedCategory)
  );

  function getBOMCost(r: DishRecipe) {
    return r.ingredients.reduce((sum, item) => sum + item.approxCostMWK, 0);
  }

  function getMargin(r: DishRecipe) {
    const cost = getBOMCost(r);
    return r.sellingPrice - cost;
  }

  function getMarginPercent(r: DishRecipe) {
    if (r.sellingPrice === 0) return 0;
    return Math.round((getMargin(r) / r.sellingPrice) * 100);
  }

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }
</script>

<div class="space-y-4 pb-4">
  <div class="pt-1">
    <h1 class="text-xl font-bold text-white tracking-tight">Recipe BOM Costing</h1>
    <p class="text-xs text-slate-400">Bill of Materials: Linking Dishes to Stock Ingredients</p>
  </div>

  <!-- Category Filter -->
  <div class="flex space-x-2 overflow-x-auto pb-1 scrollbar-none">
    {#each ['All', 'Breakfast', 'Lunch', 'Dinner', 'Snack'] as cat}
      <button
        onclick={() => (selectedCategory = cat as any)}
        class="flex-shrink-0 px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all {selectedCategory === cat ? 'bg-emerald-500 text-slate-950 border-emerald-500 shadow-lg shadow-emerald-500/25' : 'bg-slate-900 text-slate-400 border-slate-700 hover:border-slate-500'}"
      >
        {cat}
      </button>
    {/each}
  </div>

  <!-- Recipes List -->
  <div class="space-y-3">
    {#each filteredRecipes as recipe}
      {@const cost = getBOMCost(recipe)}
      {@const margin = getMargin(recipe)}
      {@const marginPct = getMarginPercent(recipe)}

      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-md space-y-3">
        <!-- Dish Title & Badges -->
        <div class="flex items-start justify-between">
          <div>
            <h2 class="text-base font-bold text-white">{recipe.dishName}</h2>
            <span class="text-[10px] px-2 py-0.5 rounded-full bg-slate-800 text-slate-300 border border-slate-700 font-semibold uppercase">
              {recipe.category}
            </span>
          </div>
          <div class="text-right">
            <span class="text-xs text-slate-400 block">Menu Price</span>
            <span class="text-base font-bold text-emerald-400 font-['Oxanium',sans-serif]">
              {fmtMwk(recipe.sellingPrice)}
            </span>
          </div>
        </div>

        <!-- BOM Ingredients Table -->
        <div class="bg-slate-950/60 rounded-xl p-3 border border-slate-800/80 space-y-2">
          <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">
            Recipe Ingredients (Per Portion)
          </p>
          <div class="space-y-1.5">
            {#each recipe.ingredients as ing}
              <div class="flex items-center justify-between text-xs">
                <span class="text-slate-300">
                  {ing.name} <span class="text-slate-500">({ing.qty} {ing.unit})</span>
                </span>
                <span class="text-slate-400 font-mono">{fmtMwk(ing.approxCostMWK)}</span>
              </div>
            {/each}
          </div>
        </div>

        <!-- Financial Summary for Kitchen & Management -->
        <div class="grid grid-cols-3 gap-2 text-center pt-1 border-t border-slate-800">
          <div class="bg-slate-800/60 rounded-xl p-2">
            <span class="text-[10px] text-slate-400 uppercase font-semibold block">Food Cost</span>
            <span class="text-xs font-bold text-amber-300 font-['Oxanium',sans-serif]">{fmtMwk(cost)}</span>
          </div>
          <div class="bg-slate-800/60 rounded-xl p-2">
            <span class="text-[10px] text-slate-400 uppercase font-semibold block">Gross Margin</span>
            <span class="text-xs font-bold text-emerald-400 font-['Oxanium',sans-serif]">{fmtMwk(margin)}</span>
          </div>
          <div class="bg-slate-800/60 rounded-xl p-2">
            <span class="text-[10px] text-slate-400 uppercase font-semibold block">Margin %</span>
            <span class="text-xs font-bold text-cyan-300 font-['Oxanium',sans-serif]">{marginPct}%</span>
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>

