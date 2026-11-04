<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { authStore } from '#lib/stores/authStore';
  import { guestStore } from '#lib/stores/guestStore';
  import { onMount } from 'svelte';

  const stay = guestStore.stay;
  const orders = guestStore.orders;

  type MenuItem = {
    id: string;
    name: string;
    category: 'Breakfast' | 'Mains' | 'Snacks' | 'Drinks';
    price: number;
    description: string;
    image_url?: string;
    badge?: string;
    popular?: boolean;
  };

  const menuItems: MenuItem[] = [
    {
      id: 'm1',
      name: 'Nsima with Stewed Chicken',
      category: 'Mains',
      price: 3500,
      description: 'Traditional Malawian maize meal served with rich tomato chicken stew and fresh seasonal greens.',
      image_url: 'https://images.unsplash.com/photo-1547496502-affa22d38842?auto=format&fit=crop&w=600&q=80',
      badge: 'Local Specialty',
      popular: true
    },
    {
      id: 'm2',
      name: 'Grilled Lake Tilapia (Chambo)',
      category: 'Mains',
      price: 4500,
      description: 'Fresh whole Lake Malawi Chambo fish grilled to perfection with lemon butter and hand-cut chips.',
      image_url: 'https://images.unsplash.com/photo-1519708227418-c8fd9a32b7a2?auto=format&fit=crop&w=600&q=80',
      badge: 'Lake Catch',
      popular: true
    },
    {
      id: 'm3',
      name: 'Prime Beef Tenderloin Steak',
      category: 'Mains',
      price: 8500,
      description: 'Pan-seared tenderloin served with garlic mashed potatoes, roasted vegetables, and pepper sauce.',
      image_url: 'https://images.unsplash.com/photo-1544025162-d76694265947?auto=format&fit=crop&w=600&q=80',
      badge: "Chef's Cut"
    },
    {
      id: 'm4',
      name: 'Full English Resort Breakfast',
      category: 'Breakfast',
      price: 5500,
      description: 'Two eggs, bacon rashers, pork sausages, baked beans, grilled tomato, and toasted sourdough.',
      image_url: 'https://images.unsplash.com/photo-1533089860892-a7c6f0a88666?auto=format&fit=crop&w=600&q=80',
      badge: 'Breakfast Special'
    },
    {
      id: 'm5',
      name: 'Warm Mandasi & Local Honey',
      category: 'Breakfast',
      price: 800,
      description: 'Freshly fried fluffy African donuts served with butter, local Zomba honey, and berry jam.',
      image_url: 'https://images.unsplash.com/photo-1589367920969-ab8e050bbb04?auto=format&fit=crop&w=600&q=80',
      badge: 'Traditional'
    },
    {
      id: 'm6',
      name: 'Crispy Chips & Peri-Peri Sauce',
      category: 'Snacks',
      price: 1200,
      description: 'Golden potatoes fried crisp, tossed in house seasoned salt with tangy dipping sauce.',
      image_url: 'https://images.unsplash.com/photo-1573080496219-bb080dd4f877?auto=format&fit=crop&w=600&q=80',
      badge: 'Quick Bite'
    },
    {
      id: 'm7',
      name: 'Fresh Mango & Passion Juice',
      category: 'Drinks',
      price: 1500,
      description: 'Freshly pressed seasonal tropical fruits, chilled and served over ice.',
      image_url: 'https://images.unsplash.com/photo-1613478223719-2ab802602423?auto=format&fit=crop&w=600&q=80',
      badge: '100% Fresh'
    },
    {
      id: 'm8',
      name: 'Carlsberg Green Lager (500ml)',
      category: 'Drinks',
      price: 1200,
      description: 'Malawi brewed premium lager beer, ice cold.',
      image_url: 'https://images.unsplash.com/photo-1608270119358-1c4b92b6a50f?auto=format&fit=crop&w=600&q=80',
      badge: 'Chilled'
    },
    {
      id: 'm9',
      name: 'Malawi Gin & Tonic',
      category: 'Drinks',
      price: 2500,
      description: 'Signature MGT with legendary Malawi Gin, crisp tonic water, and fresh lime wedge.',
      image_url: 'https://images.unsplash.com/photo-1514362545857-3bc16c4c7d1b?auto=format&fit=crop&w=600&q=80',
      badge: 'Resort Signature'
    }
  ];

  let auth = $state<{ token: string; isLoggedIn: boolean }>({ token: '', isLoggedIn: false });
  authStore.subscribe(v => { auth = v; });

  let isCheckedIn = $state(false); // strictly false until verified occupied & open
  let hasCheckedStatus = $state(false);

  onMount(async () => {
    if (auth.isLoggedIn && auth.token) {
      try {
        const res = await apiFetch('/api/v1/guest/bookings', {
          headers: { 'Authorization': `Bearer ${auth.token}` }
        });
        if (res.ok) {
          const data = await res.json();
          if (Array.isArray(data) && data.length > 0) {
            const active = data.find((b: any) => b.folio_status !== 'Closed');
            isCheckedIn = active ? (active.checked_in === true && active.room_status === 'Occupied') : false;
          } else {
            isCheckedIn = false;
          }
        }
      } catch {
        isCheckedIn = false;
      }
    } else {
      isCheckedIn = false;
    }
    hasCheckedStatus = true;
  });

  let selectedCategory = $state<'All' | 'Breakfast' | 'Mains' | 'Snacks' | 'Drinks'>('All');
  let cart = $state<{ item: MenuItem; quantity: number }[]>([]);
  let deliveryNotes = $state('');
  let showCartSheet = $state(false);
  let toastMsg = $state('');

  let filteredItems = $derived(
    selectedCategory === 'All'
      ? menuItems
      : menuItems.filter(i => i.category === selectedCategory)
  );

  let cartCount = $derived(
    cart.reduce((sum, entry) => sum + entry.quantity, 0)
  );

  let cartTotal = $derived(
    cart.reduce((sum, entry) => sum + entry.item.price * entry.quantity, 0)
  );

  function addToCart(item: MenuItem) {
    const existing = cart.find(c => c.item.id === item.id);
    if (existing) {
      cart = cart.map(c => c.item.id === item.id ? { ...c, quantity: c.quantity + 1 } : c);
    } else {
      cart = [...cart, { item, quantity: 1 }];
    }
    toastMsg = `Added ${item.name} to tray`;
    setTimeout(() => {
      toastMsg = '';
    }, 1500);
  }

  function updateQty(itemId: string, delta: number) {
    cart = cart
      .map(c => {
        if (c.item.id !== itemId) return c;
        const q = c.quantity + delta;
        return q <= 0 ? null : { ...c, quantity: q };
      })
      .filter(Boolean) as { item: MenuItem; quantity: number }[];
  }

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }
</script>

<div class="space-y-4 pb-20">
  <!-- Title & Active Room Banner -->
  <div class="pt-1 flex items-start justify-between">
    <div>
      <h1 class="text-xl font-bold text-white tracking-tight">In-Room Dining</h1>
      <p class="text-xs text-slate-400">Delivered directly to Room {$stay.roomNumber}</p>
    </div>
    <div class="flex items-center space-x-1 px-2.5 py-1 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-[11px] font-semibold">
      <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-ping"></span>
      <span>Kitchen Open</span>
    </div>
  </div>

  {#if toastMsg}
    <div class="p-3 rounded-2xl bg-amber-500/20 border border-amber-500/40 text-amber-200 text-xs text-center font-bold shadow-lg">
      {toastMsg}
    </div>
  {/if}

  <div class="p-3.5 rounded-2xl bg-amber-500/15 border border-amber-500/30 text-amber-200 text-xs text-center">
    Online ordering is not available yet. No order will be placed or charged to your folio. Please contact Front Desk to order.
  </div>

  {#if !isCheckedIn}
    <div class="p-3.5 rounded-2xl bg-amber-500/15 border border-amber-500/30 text-amber-300 text-xs text-center space-y-1">
      <div class="flex items-center justify-center space-x-1.5 font-bold">
        <i class="ph ph-lock text-sm"></i>
        <span>Front Desk Check-In Required</span>
      </div>
      <p class="text-[11px] text-slate-300">
        In-room dining ordering activates once the Front Desk completes your check-in.
      </p>
    </div>
  {/if}

  <!-- Category Pills -->
  <div class="flex space-x-2 overflow-x-auto pb-1 scrollbar-none">
    {#each ['All', 'Breakfast', 'Mains', 'Snacks', 'Drinks'] as cat}
      <button
        onclick={() => (selectedCategory = cat as any)}
        class="flex-shrink-0 px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all {selectedCategory === cat ? 'bg-amber-500 text-slate-950 border-amber-500 shadow-md shadow-amber-500/30' : 'bg-slate-900 text-slate-400 border-slate-700 hover:border-slate-500'}"
      >
        {cat}
      </button>
    {/each}
  </div>

  <!-- Menu Dish Cards -->
  <div class="space-y-3">
    {#each filteredItems as item}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl overflow-hidden shadow-md hover:border-slate-700 transition flex flex-col">
        {#if item.image_url}
          <div class="relative h-36 w-full overflow-hidden bg-slate-950">
            <img src={item.image_url} alt={item.name} class="w-full h-full object-cover" loading="lazy" />
            <div class="absolute inset-0 bg-gradient-to-t from-slate-900 via-transparent to-transparent"></div>
            {#if item.badge}
              <span class="absolute top-2.5 left-2.5 text-[9px] px-2 py-0.5 rounded-full font-bold uppercase tracking-wider bg-amber-500 text-slate-950 shadow-md">
                {item.badge}
              </span>
            {/if}
            <span class="absolute bottom-2.5 right-2.5 text-xs font-bold text-amber-300 font-['Oxanium',sans-serif] bg-slate-950/80 px-2 py-0.5 rounded-lg border border-slate-800">
              {fmtMwk(item.price)}
            </span>
          </div>
        {/if}

        <div class="p-3.5 space-y-2 flex-1 flex flex-col justify-between">
          <div>
            <div class="flex items-center justify-between mb-0.5">
              <h2 class="font-bold text-white text-sm">{item.name}</h2>
              {#if !item.image_url}
                <span class="text-sm font-bold text-amber-400 font-['Oxanium',sans-serif]">
                  {fmtMwk(item.price)}
                </span>
              {/if}
            </div>
            <p class="text-xs text-slate-400 leading-relaxed">{item.description}</p>
          </div>

          <div class="flex items-center justify-between pt-1 border-t border-slate-800/80">
            <span class="text-[11px] text-slate-400">Prep time: ~15-20 mins</span>
            <button
              onclick={() => addToCart(item)}
              disabled={!isCheckedIn}
              class="px-3.5 py-1.5 rounded-xl bg-slate-800 hover:bg-amber-500 hover:text-slate-950 text-slate-200 font-bold text-xs border border-slate-700 transition active:scale-95 flex items-center space-x-1 disabled:opacity-40 disabled:cursor-not-allowed"
            >
              <i class="ph {isCheckedIn ? 'ph-plus' : 'ph-lock'} text-sm"></i>
              <span>{isCheckedIn ? 'Add to Tray' : 'Locked'}</span>
            </button>
          </div>
        </div>
      </div>
    {/each}
  </div>

  <!-- Recent / Active Orders History Accordion -->
  {#if $orders.length > 0}
    <div class="pt-2">
      <h2 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2">My In-Room Dining Orders</h2>
      <div class="space-y-2.5">
        {#each $orders as ord}
          <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5 space-y-2">
            <div class="flex items-center justify-between">
              <div>
                <span class="font-bold text-white text-xs font-['Oxanium',sans-serif]">{ord.orderNumber}</span>
                <span class="text-[10px] text-slate-400 block">{ord.createdAt}</span>
              </div>
              <span class="text-[10px] px-2 py-0.5 rounded-full font-bold uppercase tracking-wider {ord.status === 'Delivered' ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30' : 'bg-amber-500/20 text-amber-300 border border-amber-500/40 animate-pulse'}">
                {ord.status}
              </span>
            </div>

            <div class="space-y-1 text-xs text-slate-300">
              {#each ord.items as line}
                <div class="flex justify-between">
                  <span>{line.name} × {line.quantity}</span>
                  <span class="font-mono text-slate-400">{fmtMwk(line.price * line.quantity)}</span>
                </div>
              {/each}
            </div>

            <div class="flex items-center justify-between pt-1.5 border-t border-slate-800 text-xs">
              <span class="text-slate-400">Charged to Room Folio:</span>
              <span class="font-bold text-amber-300 font-['Oxanium',sans-serif]">{fmtMwk(ord.totalAmount)}</span>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>

<!-- Floating Delivery Tray Bar (When Cart Has Items) -->
{#if cartCount > 0}
  <div class="fixed bottom-20 inset-x-0 z-40 max-w-md mx-auto px-4">
    <button
      onclick={() => (showCartSheet = true)}
      class="w-full bg-gradient-to-r from-amber-500 to-amber-600 text-slate-950 rounded-2xl p-3.5 shadow-2xl flex items-center justify-between active:scale-98 transition font-bold"
    >
      <div class="flex items-center space-x-2.5">
        <div class="w-8 h-8 rounded-xl bg-slate-950 text-amber-400 flex items-center justify-center font-bold text-sm">
          {cartCount}
        </div>
        <div class="text-left">
          <span class="text-xs block leading-tight">View Room {$stay.roomNumber} Tray</span>
          <span class="text-[10px] text-slate-900/80">Tap to review & charge to stay</span>
        </div>
      </div>
      <span class="text-base font-['Oxanium',sans-serif] font-extrabold">{fmtMwk(cartTotal)}</span>
    </button>
  </div>
{/if}

<!-- Cart Review & Room Charge Bottom Sheet -->
{#if showCartSheet}
  <div class="fixed inset-0 z-50 flex items-end justify-center bg-slate-950/80 backdrop-blur-sm" role="dialog" aria-modal="true">
    <div class="w-full max-w-md bg-slate-900 border border-slate-700 rounded-t-3xl p-5 space-y-4 shadow-2xl max-h-[85vh] overflow-y-auto">
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <div>
          <h2 class="font-bold text-white text-base">In-Room Dining Tray</h2>
          <p class="text-xs text-amber-400">Delivering to Room {$stay.roomNumber} · {$stay.guestName}</p>
        </div>
        <button onclick={() => (showCartSheet = false)} class="p-2 text-slate-400 hover:text-white" aria-label="Close tray">
          <i class="ph ph-x text-lg"></i>
        </button>
      </div>

      <!-- Itemized Cart -->
      <div class="space-y-2.5">
        {#each cart as entry}
          <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-3 flex items-center justify-between text-xs">
            <div class="flex-1 pr-2 truncate">
              <p class="font-bold text-white">{entry.item.name}</p>
              <p class="text-slate-400 text-[11px]">{fmtMwk(entry.item.price)} each</p>
            </div>
            <div class="flex items-center space-x-2">
              <button onclick={() => updateQty(entry.item.id, -1)} class="w-7 h-7 rounded-lg bg-slate-800 text-slate-200 font-bold active:scale-95" aria-label="Decrease">−</button>
              <span class="w-4 text-center font-bold text-white">{entry.quantity}</span>
              <button onclick={() => updateQty(entry.item.id, 1)} class="w-7 h-7 rounded-lg bg-slate-800 text-slate-200 font-bold active:scale-95" aria-label="Increase">+</button>
            </div>
          </div>
        {/each}
      </div>

      <!-- Delivery Instructions -->
      <div>
        <label for="tray-notes" class="block text-xs text-slate-400 font-semibold uppercase tracking-wider mb-1">
          Special Dietary or Delivery Notes
        </label>
        <input
          id="tray-notes"
          type="text"
          bind:value={deliveryNotes}
          placeholder="e.g. Extra napkins, no spicy pepper, cutlery for 2"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500"
        />
      </div>

      <!-- Cost Breakdown -->
      <div class="bg-slate-950/80 rounded-2xl p-4 border border-slate-800 space-y-1.5 text-xs">
        <div class="flex justify-between text-slate-400">
          <span>Items Total</span>
          <span class="font-mono text-white">{fmtMwk(cartTotal)}</span>
        </div>
        <div class="flex justify-between text-slate-400">
          <span>Room Delivery Service</span>
          <span class="text-emerald-400 font-semibold">COMPLIMENTARY</span>
        </div>
        <div class="flex justify-between text-base font-bold text-white pt-2 border-t border-slate-800">
          <span>Total To Charge</span>
          <span class="text-amber-400 font-['Oxanium',sans-serif]">{fmtMwk(cartTotal)}</span>
        </div>
      </div>

      <!-- Submit Order Button -->
      <button
        disabled={true}
        class="w-full py-3.5 rounded-2xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-xl shadow-amber-500/30 active:scale-98 transition disabled:opacity-50 flex items-center justify-center space-x-1.5"
      >
        <i class="ph ph-receipt text-base"></i>
        <span>Online ordering unavailable</span>
      </button>
    </div>
  </div>
{/if}
