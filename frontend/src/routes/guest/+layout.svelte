<script lang="ts">
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { authStore } from '#lib/stores/authStore';
  import { guestStore } from '#lib/stores/guestStore';

  let { children } = $props();

  const stay = guestStore.stay;
  let showWifiModal = $state(false);
  let copyFeedback = $state(false);

  let auth = $state<{ isLoggedIn: boolean; role: string; fullName: string }>({ isLoggedIn: false, role: '', fullName: '' });
  authStore.subscribe(v => { auth = v; });

  let isGuest = $derived(auth.isLoggedIn && auth.role === 'guest');

  function isActive(path: string): boolean {
    if (path === '/guest') return page.url.pathname === '/guest';
    return page.url.pathname.startsWith(path);
  }

  function copyWifi() {
    if (typeof navigator !== 'undefined' && navigator.clipboard) {
      navigator.clipboard.writeText('LaLunaGuest2025');
      copyFeedback = true;
      setTimeout(() => { copyFeedback = false; }, 2000);
    }
  }

  function handleLogout() {
    authStore.logout();
    goto('/book');
  }
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-['Inter',sans-serif] pb-24 select-none">
  <!-- Top Mobile Sticky Header -->
  <header class="sticky top-0 z-40 bg-slate-900/90 backdrop-blur-md border-b border-slate-800/80 px-4 py-3 shadow-lg">
    <div class="max-w-md mx-auto flex items-center justify-between">
      <!-- Hotel Brand & Guest Identification -->
      <a href="/guest" class="flex items-center space-x-2.5 active:scale-95 transition-transform">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-amber-400 via-amber-500 to-amber-600 flex items-center justify-center shadow-lg shadow-amber-500/20 text-slate-950 font-bold">
          <i class="ph ph-moon text-xl"></i>
        </div>
        <div>
          <div class="flex items-center space-x-1.5">
            <span class="font-['Oxanium',sans-serif] font-bold tracking-wider text-base text-white">LA LUNA</span>
            <span class="text-[10px] uppercase font-bold tracking-widest bg-amber-500/20 text-amber-300 px-1.5 py-0.5 rounded border border-amber-500/30">
              Guest
            </span>
          </div>
          <p class="text-[11px] text-slate-400 font-medium truncate max-w-[140px]">
            {auth.fullName || $stay.guestName || 'Resident'}
          </p>
        </div>
      </a>

      <!-- Quick Actions: WiFi & Booking/Exit -->
      <div class="flex items-center space-x-2">
        <button
          onclick={() => (showWifiModal = true)}
          class="flex items-center space-x-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-amber-300 border border-slate-700 active:scale-95 transition"
          title="Hotel WiFi Details"
        >
          <i class="ph ph-wifi-high text-sm"></i>
          <span class="text-[11px]">WiFi</span>
        </button>

        {#if auth.isLoggedIn}
          <button
            onclick={handleLogout}
            class="p-2 rounded-xl bg-slate-800 text-slate-400 hover:text-white border border-slate-700/60 active:scale-95 transition"
            title="Sign Out"
          >
            <i class="ph ph-sign-out text-base"></i>
          </button>
        {:else}
          <a
            href="/book/login"
            class="px-2.5 py-1 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 text-xs font-bold active:scale-95 transition"
          >
            Login
          </a>
        {/if}
      </div>
    </div>
  </header>

  <!-- Auth Gate Check -->
  {#if !auth.isLoggedIn}
    <main class="flex-1 w-full max-w-md mx-auto px-4 py-12 text-center space-y-4">
      <div class="w-16 h-16 rounded-2xl bg-amber-500/10 border-2 border-amber-500/30 text-amber-400 mx-auto flex items-center justify-center">
        <i class="ph ph-lock text-3xl"></i>
      </div>
      <div>
        <h2 class="text-xl font-bold text-white font-['Oxanium',sans-serif]">Guest Portal Access</h2>
        <p class="text-xs text-slate-400 mt-1 max-w-xs mx-auto">
          Please log in or create an account to view your stay details, dining, and room bill.
        </p>
      </div>
      <div class="flex flex-col space-y-2 max-w-xs mx-auto pt-2">
        <a
          href="/book/login"
          class="py-3 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider text-center transition active:scale-95 shadow-lg shadow-amber-500/20"
        >
          Sign In / Create Account
        </a>
        <a
          href="/book"
          class="py-3 rounded-xl bg-slate-800 hover:bg-slate-700 text-white font-bold text-xs uppercase tracking-wider text-center transition border border-slate-700 active:scale-95"
        >
          Browse & Book Rooms
        </a>
      </div>
    </main>
  {:else}
    <!-- Main Content Body -->
    <main class="flex-1 w-full max-w-md mx-auto px-4 py-4 overflow-y-auto">
      {@render children()}
    </main>
  {/if}

  <!-- Fixed Bottom 4-Tile Navigation (Simplified: Keycard Removed) -->
  {#if auth.isLoggedIn}
    <nav class="fixed bottom-0 inset-x-0 z-40 bg-slate-900/95 backdrop-blur-lg border-t border-slate-800 shadow-2xl">
      <div class="max-w-md mx-auto grid grid-cols-4 h-16">
        <!-- 1. Stay Hub -->
        <a
          href="/guest"
          class="flex flex-col items-center justify-center text-center transition-colors {isActive('/guest') && page.url.pathname === '/guest'
            ? 'text-amber-400 font-semibold'
            : 'text-slate-400 hover:text-slate-200'}"
        >
          <i class="ph ph-house-line text-xl mb-0.5"></i>
          <span class="text-[10px] tracking-tight">Stay</span>
        </a>

        <!-- 2. Dining / Room Service -->
        <a
          href="/guest/dining"
          class="flex flex-col items-center justify-center text-center transition-colors {isActive('/guest/dining')
            ? 'text-amber-400 font-semibold'
            : 'text-slate-400 hover:text-slate-200'}"
        >
          <i class="ph ph-fork-knife text-xl mb-0.5"></i>
          <span class="text-[10px] tracking-tight">Dining</span>
        </a>

        <!-- 3. Folio & Bill -->
        <a
          href="/guest/folio"
          class="flex flex-col items-center justify-center text-center transition-colors {isActive('/guest/folio')
            ? 'text-amber-400 font-semibold'
            : 'text-slate-400 hover:text-slate-200'}"
        >
          <i class="ph ph-receipt text-xl mb-0.5"></i>
          <span class="text-[10px] tracking-tight">My Bill</span>
        </a>

        <!-- 4. Services -->
        <a
          href="/guest/services"
          class="flex flex-col items-center justify-center text-center transition-colors {isActive('/guest/services')
            ? 'text-amber-400 font-semibold'
            : 'text-slate-400 hover:text-slate-200'}"
        >
          <i class="ph ph-bell-ringing text-xl mb-0.5"></i>
          <span class="text-[10px] tracking-tight">Services</span>
        </a>
      </div>
    </nav>
  {/if}

  <!-- Hotel WiFi Modal -->
  {#if showWifiModal}
    <div
      class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/80 backdrop-blur-sm p-4"
      role="dialog"
      aria-modal="true"
    >
      <div class="w-full max-w-sm bg-slate-900 border border-slate-700 rounded-3xl p-6 space-y-4 shadow-2xl text-center">
        <div class="w-12 h-12 rounded-2xl bg-amber-500/20 text-amber-400 mx-auto flex items-center justify-center">
          <i class="ph ph-wifi-high text-2xl"></i>
        </div>

        <div>
          <h2 class="text-lg font-bold text-white">Hotel WiFi</h2>
          <p class="text-xs text-slate-400 mt-0.5">Complimentary High-Speed Internet</p>
        </div>

        <div class="bg-slate-950/80 rounded-2xl p-4 border border-slate-800 space-y-3 text-left">
          <div>
            <span class="text-[10px] text-slate-400 uppercase font-semibold block">Network SSID</span>
            <span class="text-sm font-bold text-white font-mono">LaLuna-Guest-WiFi</span>
          </div>

          <div>
            <span class="text-[10px] text-slate-400 uppercase font-semibold block">Password</span>
            <div class="flex items-center justify-between">
              <span class="text-sm font-bold text-amber-300 font-mono">LaLunaGuest2025</span>
              <button
                onclick={copyWifi}
                class="px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-xs text-slate-200 font-medium active:scale-95 transition"
              >
                {copyFeedback ? '✓ Copied' : 'Copy'}
              </button>
            </div>
          </div>
        </div>

        <button
          onclick={() => (showWifiModal = false)}
          class="w-full py-3 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider transition active:scale-95"
        >
          Done
        </button>
      </div>
    </div>
  {/if}
</div>
