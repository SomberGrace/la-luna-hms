<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { authStore } from '#lib/stores/authStore';
  import { onMount } from 'svelte';

  let { children } = $props();

  function isAuthorized(role: string, loggedIn: boolean): boolean {
    return loggedIn && role !== 'guest' && (role === 'bar' || role === 'manager');
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
    if (path === '/bar') return page.url.pathname === '/bar';
    return page.url.pathname.startsWith(path);
  }
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-['Inter',sans-serif] pb-24 select-none">
  <header class="sticky top-0 z-40 bg-slate-900/90 backdrop-blur-md border-b border-slate-800 px-4 py-3 shadow-md">
    <div class="max-w-md mx-auto flex items-center justify-between">
      <a href="/bar" class="flex items-center space-x-2.5 active:scale-95 transition-transform">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-purple-500 to-purple-700 flex items-center justify-center shadow-lg shadow-purple-500/20 text-white font-bold">
          <i class="ph ph-wine text-xl"></i>
        </div>
        <div>
          <div class="flex items-center space-x-1.5">
            <span class="font-['Oxanium',sans-serif] font-bold tracking-wider text-base text-white">LA LUNA</span>
            <span class="text-[10px] uppercase font-bold tracking-widest bg-purple-500/20 text-purple-300 px-1.5 py-0.5 rounded">Bar</span>
          </div>
          <p class="text-[11px] text-slate-400 font-medium">Bar Inventory & Audit</p>
        </div>
      </a>
      <button type="button" onclick={handleLogout} class="p-2 rounded-lg bg-slate-800 text-slate-300 hover:text-purple-400 border border-slate-700/60 transition-colors" title="Sign Out">
        <i class="ph ph-sign-out text-lg"></i>
      </button>
    </div>
  </header>

  <main class="flex-1 w-full max-w-md mx-auto px-4 py-4">
    {@render children()}
  </main>

  <nav class="fixed bottom-0 inset-x-0 z-40 bg-slate-900/95 backdrop-blur-lg border-t border-slate-800 shadow-2xl">
    <div class="max-w-md mx-auto grid grid-cols-4 h-16">
      <a href="/bar" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/bar') && page.url.pathname === '/bar' ? 'text-purple-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-house text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Hub</span>
      </a>
      <a href="/bar/stock" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/bar/stock') ? 'text-purple-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-package text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Stock</span>
      </a>
      <a href="/bar/sales" class="flex flex-col items-center justify-center text-center group">
        <div class="w-11 h-11 -mt-5 rounded-full bg-gradient-to-tr from-purple-700 to-purple-400 text-white flex items-center justify-center shadow-lg shadow-purple-500/40 border-2 border-slate-950 active:scale-95 transition-transform">
          <i class="ph ph-plus text-xl font-bold"></i>
        </div>
        <span class="text-[10px] font-semibold mt-1 {isActive('/bar/sales') ? 'text-purple-400' : 'text-slate-300'}">Sales</span>
      </a>
      <a href="/bar/reports" class="flex flex-col items-center justify-center text-center transition-colors {isActive('/bar/reports') ? 'text-purple-400' : 'text-slate-400 hover:text-slate-200'}">
        <i class="ph ph-chart-bar text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold">Reports</span>
      </a>
    </div>
  </nav>
</div>

