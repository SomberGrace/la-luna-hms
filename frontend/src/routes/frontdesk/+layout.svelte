<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { isOnline, pendingOfflineCount, roomStore } from '#lib/stores/roomStore';
  import { authStore } from '#lib/stores/authStore';
  import { onMount } from 'svelte';

  let { children } = $props();

  function isAuthorizedStaff(role: string, loggedIn: boolean): boolean {
    return loggedIn && role !== 'guest' && (role === 'frontdesk' || role === 'manager');
  }

  onMount(() => {
    if (!isAuthorizedStaff($authStore.role, $authStore.isLoggedIn)) {
      goto('/');
      return;
    }
    roomStore.fetchRooms();
  });

  $effect(() => {
    if (typeof window !== 'undefined' && !isAuthorizedStaff($authStore.role, $authStore.isLoggedIn)) {
      goto('/');
    }
  });

  function handleLogout() {
    authStore.logout();
    goto('/');
  }

  // Active path checking
  function isActive(path: string): boolean {
    if (path === '/frontdesk') {
      return page.url.pathname === '/frontdesk';
    }
    return page.url.pathname.startsWith(path);
  }
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-['Inter',sans-serif] pb-24 select-none">
  <!-- Top Mobile Sticky Header -->
  <header class="sticky top-0 z-40 bg-slate-900/90 backdrop-blur-md border-b border-slate-800 px-4 py-3 shadow-md">
    <div class="max-w-md mx-auto flex items-center justify-between">
      <!-- Brand & Hub Link -->
      <a href="/frontdesk" class="flex items-center space-x-2.5 active:scale-95 transition-transform">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-amber-500 to-amber-600 flex items-center justify-center shadow-lg shadow-amber-500/20 text-slate-950 font-bold">
          <i class="ph ph-moon-stars text-xl"></i>
        </div>
        <div>
          <div class="flex items-center space-x-1.5">
            <span class="font-['Oxanium',sans-serif] font-bold tracking-wider text-base text-white">LA LUNA</span>
            <span class="text-[10px] uppercase font-bold tracking-widest bg-amber-500/20 text-amber-400 px-1.5 py-0.5 rounded">Mobile FD</span>
          </div>
          <p class="text-[11px] text-slate-400 font-medium">Front Desk Reception</p>
        </div>
      </a>

      <!-- Quick Status Badges: Sync / Offline & User Role -->
      <div class="flex items-center space-x-2">
        <!-- Offline / Sync Status Pill -->
        <a
          href="/frontdesk/offline"
          class="flex items-center space-x-1.5 px-2.5 py-1 rounded-full text-xs font-semibold border transition-colors {$isOnline
            ? 'bg-emerald-950/60 border-emerald-700/50 text-emerald-400'
            : 'bg-rose-950/70 border-rose-600 text-rose-300 animate-pulse'}"
          title={$isOnline ? 'Online - Live backend synced' : 'Offline - Queued locally'}
        >
          <span class="w-2 h-2 rounded-full {$isOnline ? 'bg-emerald-400 animate-ping' : 'bg-rose-400'}"></span>
          <span class="text-[11px]">{$isOnline ? 'Live' : 'Offline'}</span>
          {#if $pendingOfflineCount > 0}
            <span class="ml-1 bg-amber-500 text-slate-950 font-bold px-1.5 py-0.2 rounded-full text-[10px]">
              {$pendingOfflineCount}
            </span>
          {/if}
        </a>

        <!-- Quick Menu link to Reports / Settings -->
        <a
          href="/frontdesk/reports"
          class="p-2 rounded-lg bg-slate-800 text-slate-300 hover:text-amber-400 active:bg-slate-700 border border-slate-700/60 transition-colors"
          title="Daily Summary Reports"
        >
          <i class="ph ph-chart-pie-slice text-lg"></i>
        </a>

        <!-- Sign out button -->
        <button
          type="button"
          onclick={handleLogout}
          class="p-2 rounded-lg bg-slate-800 text-slate-300 hover:text-rose-400 active:bg-slate-700 border border-slate-700/60 transition-colors"
          title="Sign Out of Station"
        >
          <i class="ph ph-sign-out text-lg"></i>
        </button>
      </div>
    </div>
  </header>

  <!-- Main Mobile Content Area -->
  <main class="flex-1 w-full max-w-md mx-auto px-4 py-4 overflow-y-auto">
    {@render children()}
  </main>

  <!-- Fixed Bottom 5-Tile Touch Navigation (Thumb-friendly mobile layout) -->
  <nav class="fixed bottom-0 inset-x-0 z-40 bg-slate-900/95 backdrop-blur-lg border-t border-slate-800 shadow-2xl">
    <div class="max-w-md mx-auto grid grid-cols-5 h-16">
      <!-- 1. Hub -->
      <a
        href="/frontdesk"
        class="flex flex-col items-center justify-center text-center transition-colors {isActive('/frontdesk') && page.url.pathname === '/frontdesk'
          ? 'text-amber-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        <i class="ph ph-squares-four text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold tracking-tight">Hub</span>
      </a>

      <!-- 2. Rooms -->
      <a
        href="/frontdesk/rooms"
        class="flex flex-col items-center justify-center text-center transition-colors {isActive('/frontdesk/rooms')
          ? 'text-amber-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        <i class="ph ph-bed text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold tracking-tight">Rooms</span>
      </a>

      <!-- 3. Walk-In (Center Highlight Action) -->
      <a
        href="/frontdesk/walkin"
        class="flex flex-col items-center justify-center text-center group"
      >
        <div class="w-11 h-11 -mt-5 rounded-full bg-gradient-to-tr from-amber-600 to-amber-400 text-slate-950 flex items-center justify-center shadow-lg shadow-amber-500/40 border-2 border-slate-950 active:scale-95 transition-transform">
          <i class="ph ph-user-plus text-xl font-bold"></i>
        </div>
        <span class="text-[10px] font-semibold mt-1 {isActive('/frontdesk/walkin') ? 'text-amber-400' : 'text-slate-300'}">Walk-In</span>
      </a>

      <!-- 4. Folio & Billing -->
      <a
        href="/frontdesk/folio"
        class="flex flex-col items-center justify-center text-center transition-colors {isActive('/frontdesk/folio') || isActive('/frontdesk/payments')
          ? 'text-amber-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        <i class="ph ph-receipt text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold tracking-tight">Folio</span>
      </a>

      <!-- 5. Check-Out Terminal -->
      <a
        href="/frontdesk/checkout"
        class="flex flex-col items-center justify-center text-center transition-colors {isActive('/frontdesk/checkout')
          ? 'text-amber-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        <i class="ph ph-sign-out text-xl mb-0.5"></i>
        <span class="text-[10px] font-semibold tracking-tight">Checkout</span>
      </a>
    </div>
  </nav>
</div>
