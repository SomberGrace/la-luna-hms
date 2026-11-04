<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { authStore } from '#lib/stores/authStore';
  import { onMount } from 'svelte';

  let { children } = $props();

  function isAuthorized(role: string, loggedIn: boolean): boolean {
    return loggedIn && role !== 'guest' && (role === 'restaurant' || role === 'manager');
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
    if (path === '/restaurant') return page.url.pathname === '/restaurant';
    return page.url.pathname.startsWith(path);
  }
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-['Inter',sans-serif] pb-24 select-none">
  <!-- Sticky Header -->
  <header class="sticky top-0 z-40 bg-slate-900/90 backdrop-blur-md border-b border-slate-800 px-4 py-3 shadow-md">
    <div class="max-w-md mx-auto flex items-center justify-between">
      <a href="/restaurant" class="flex items-center space-x-2.5 active:scale-95 transition-transform">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center shadow-lg shadow-orange-500/20 text-slate-950 font-bold">
          <i class="ph ph-fork-knife text-xl"></i>
        </div>
        <div>
          <div class="flex items-center space-x-1.5">
            <span class="font-['Oxanium',sans-serif] font-bold tracking-wider text-base text-white">LA LUNA</span>
            <span class="text-[10px] uppercase font-bold tracking-widest bg-orange-500/20 text-orange-400 px-1.5 py-0.5 rounded">Restaurant</span>
          </div>
          <p class="text-[11px] text-slate-400 font-medium">Dining & Dockets</p>
        </div>
      </a>
      <button type="button" onclick={handleLogout} class="p-2 rounded-lg bg-slate-800 text-slate-300 hover:text-orange-400 border border-slate-700/60 transition-colors" title="Sign Out">
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
      <a href="/restaurant" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/restaurant') && page.url.pathname === '/restaurant' ? 'text-orange-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-house text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Hub</span>
      </a>
      <a href="/restaurant/menu" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/restaurant/menu') ? 'text-orange-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-fork-knife text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Menu</span>
      </a>
      <a href="/restaurant/dockets" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/restaurant/dockets') ? 'text-orange-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-receipt text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Dockets</span>
      </a>
      <a href="/restaurant/reports" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/restaurant/reports') ? 'text-orange-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-chart-bar text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Reports</span>
      </a>
    </div>
  </nav>
</div>
