<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { authStore } from '#lib/stores/authStore';
  import { onMount } from 'svelte';

  let { children } = $props();

  function isAuthorized(role: string, loggedIn: boolean): boolean {
    return loggedIn && role !== 'guest' && (role === 'kitchen' || role === 'manager');
  }

  onMount(() => {
    if (!isAuthorized($authStore.role, $authStore.isLoggedIn)) {
      goto('/');
    }
  });

  $effect(() => {
    if (typeof window !== 'undefined' && !isAuthorized($authStore.role, $authStore.isLoggedIn)) {
      goto('/');
    }
  });

  function handleLogout() {
    authStore.logout();
    goto('/');
  }

  function isActive(path: string): boolean {
    if (path === '/kitchen') return page.url.pathname === '/kitchen';
    return page.url.pathname.startsWith(path);
  }
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-['Inter',sans-serif] pb-24 select-none">
  <!-- Sticky Header -->
  <header class="sticky top-0 z-40 bg-slate-900/90 backdrop-blur-md border-b border-slate-800 px-4 py-3 shadow-md">
    <div class="max-w-md mx-auto flex items-center justify-between">
      <a href="/kitchen" class="flex items-center space-x-2.5 active:scale-95 transition-transform">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center shadow-lg shadow-emerald-500/20 text-slate-950 font-bold">
          <i class="ph ph-cooking-pot text-xl"></i>
        </div>
        <div>
          <div class="flex items-center space-x-1.5">
            <span class="font-['Oxanium',sans-serif] font-bold tracking-wider text-base text-white">LA LUNA</span>
            <span class="text-[10px] uppercase font-bold tracking-widest bg-emerald-500/20 text-emerald-300 px-1.5 py-0.5 rounded">Kitchen KDS</span>
          </div>
          <p class="text-[11px] text-slate-400 font-medium">Head Chef & Stock Station</p>
        </div>
      </a>
      <button type="button" onclick={handleLogout} class="p-2 rounded-lg bg-slate-800 text-slate-300 hover:text-emerald-400 border border-slate-700/60 transition-colors" title="Sign Out">
        <i class="ph ph-sign-out text-lg"></i>
      </button>
    </div>
  </header>

  <!-- Content -->
  <main class="flex-1 w-full max-w-md mx-auto px-4 py-4">
    {@render children()}
  </main>

  <!-- Bottom Navigation -->
  <nav class="fixed bottom-0 inset-x-0 z-40 bg-slate-900/95 backdrop-blur-lg border-t border-slate-800 shadow-2xl">
    <div class="max-w-md mx-auto grid grid-cols-4 h-16">
      <a href="/kitchen" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/kitchen') && page.url.pathname === '/kitchen' ? 'text-emerald-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-fire text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Orders</span>
      </a>
      <a href="/kitchen/stock" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/kitchen/stock') ? 'text-emerald-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-scales text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Stock</span>
      </a>
      <a href="/kitchen/recipes" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/kitchen/recipes') ? 'text-emerald-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-notebook text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">BOM Recipes</span>
      </a>
      <a href="/kitchen/reports" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/kitchen/reports') ? 'text-emerald-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-chart-pie-slice text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Reports</span>
      </a>
    </div>
  </nav>
</div>

