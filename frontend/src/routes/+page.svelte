<script lang="ts">
  import { goto } from '$app/navigation';
  import { APP_ROLE, ROLE_ROUTE_MAP, isRoleAllowedForApp, normalizeUserRole, type UserRole } from '#lib/mobile';
  import { apiFetch } from '#lib/api';
  import { getGuestLandingPath } from '#lib/guestBookings';
  import { authStore } from '#lib/stores/authStore';
  import { onMount } from 'svelte';

  // Login credentials and options (Svelte 5 runes)
  let username = $state('');
  let password = $state('');
  let selectedRole = $state<UserRole>(APP_ROLE === 'shared' ? 'frontdesk' : APP_ROLE);
  let showPassword = $state(false);

  // UI state management
  let isLoading = $state(false);
  let errorMessage = $state('');

  onMount(() => {
    if (APP_ROLE === 'guest' && $authStore.isLoggedIn && $authStore.role === 'guest') {
      void getGuestLandingPath()
        .then((destination) => goto(destination))
        .catch((error) => {
          errorMessage = error instanceof Error ? error.message : 'Unable to verify your current reservation.';
        });
      return;
    }

    // If a guest account session exists in storage, reset so staff station sign-in is required
    if (APP_ROLE !== 'guest' && $authStore.isLoggedIn && $authStore.role === 'guest') {
      authStore.logout();
    }
  });

  // Role metadata definition for RBAC context
  const roles = [
    { id: 'frontdesk', label: 'Front Desk Officer', path: '/frontdesk', icon: 'ph-device-mobile', color: 'from-amber-500 to-amber-600', badge: 'text-amber-400 bg-amber-500/20' },
    { id: 'restaurant', label: 'Restaurant Attendant', path: '/restaurant', icon: 'ph-fork-knife', color: 'from-orange-500 to-orange-600', badge: 'text-orange-400 bg-orange-500/20' },
    { id: 'bar', label: 'Bar Attendant', path: '/bar', icon: 'ph-wine', color: 'from-purple-500 to-purple-600', badge: 'text-purple-400 bg-purple-500/20' },
    { id: 'kitchen', label: 'Head Chef / Kitchen', path: '/kitchen', icon: 'ph-cooking-pot', color: 'from-emerald-500 to-teal-600', badge: 'text-emerald-400 bg-emerald-500/20' },
    { id: 'manager', label: 'Hotel Manager', path: '/manager', icon: 'ph-briefcase', color: 'from-indigo-500 to-indigo-700', badge: 'text-indigo-400 bg-indigo-500/20' },
    { id: 'guest', label: 'Guest Booking & Portal', path: '/book', icon: 'ph-user', color: 'from-amber-400 to-amber-600', badge: 'text-amber-300 bg-amber-500/20' }
  ];

  function getDestinationPath(role: string): string {
    const allowedRole = normalizeUserRole(role);
    if (!allowedRole || !isRoleAllowedForApp(allowedRole)) {
      return ROLE_ROUTE_MAP[APP_ROLE];
    }
    const found = roles.find(r => r.id === allowedRole);
    return found ? found.path : ROLE_ROUTE_MAP[APP_ROLE];
  }

  function getStaffFullName(user: string, role: string): string {
    return user || 'Staff Member';
  }

  /** Submits credentials to the backend and trusts only the verified role in its response. */
  async function handleLogin(e: Event) {
    e.preventDefault();
    if (!username.trim() || !password.trim()) {
      errorMessage = 'Please enter both username and password.';
      return;
    }

    errorMessage = '';
    isLoading = true;
    try {
      const response = await apiFetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: username,
          password: password,
          role: selectedRole
        })
      });

      if (!response.ok) {
        throw new Error('Invalid credentials or role access denied.');
      }

      const data = await response.json();
      if (!data.token || !data.user?.role) {
        throw new Error('The server returned an invalid login response.');
      }

      const verifiedRole = normalizeUserRole(data.user.role);
      if (!verifiedRole || !isRoleAllowedForApp(verifiedRole)) {
        throw new Error('This account does not have access to this app.');
      }

      authStore.login({
        token: data.token,
        username: data.user?.username || username,
        fullName: data.user?.full_name || getStaffFullName(username, verifiedRole),
        role: verifiedRole
      });

      goto(verifiedRole === 'guest' ? await getGuestLandingPath() : getDestinationPath(verifiedRole));
    } catch (error) {
      errorMessage = error instanceof Error ? error.message : 'Unable to sign in. Check your connection and try again.';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-slate-950 flex flex-col items-center justify-center p-4">
  <!-- Brand Header -->
  <div class="flex flex-col items-center justify-center mb-6 text-center">
    <div class="inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-amber-500/10 border-2 border-amber-400 text-amber-400 mb-3 shadow-lg shadow-amber-500/20">
      <i class="ph ph-moon text-4xl"></i>
    </div>
    <h1 class="font-['Oxanium',sans-serif] text-3xl font-bold text-white tracking-wider">LA LUNA HMS</h1>
    <p class="text-xs text-slate-400 font-medium mt-1">Integrated Hotel Management System</p>
  </div>

  <div class="max-w-md w-full bg-slate-900 border border-slate-800 rounded-3xl shadow-2xl overflow-hidden p-6 sm:p-8 space-y-6">
    <!-- Header -->
    <div class="text-center">
      <h2 class="font-['Inter',sans-serif] font-bold text-2xl text-white tracking-tight">
        {APP_ROLE === 'shared' ? 'STATION SIGN IN' : `${APP_ROLE.toUpperCase()} SIGN IN`}
      </h2>
      <p class="text-xs text-slate-400 mt-1">
        {APP_ROLE === 'shared' ? 'Select your workstation department to enter' : 'Sign in to your role-specific app'}
      </p>
    </div>

    {#if APP_ROLE === 'shared'}
    <!-- Quick Access Department Hubs Carousel -->
    <div class="space-y-2">
      <p class="text-[10px] text-slate-400 uppercase font-bold tracking-wider">Direct Station Access</p>
      <div class="grid grid-cols-2 gap-2">
        <a
          href="/frontdesk"
          class="flex items-center space-x-2 p-2.5 rounded-xl bg-amber-500/10 border border-amber-500/30 hover:bg-amber-500/20 transition active:scale-95"
        >
          <div class="w-7 h-7 rounded-lg bg-amber-500 text-slate-950 flex items-center justify-center font-bold">
            <i class="ph ph-device-mobile text-sm"></i>
          </div>
          <div class="text-left truncate">
            <span class="text-xs font-bold text-amber-300 block leading-tight">Front Desk</span>
            <span class="text-[10px] text-slate-400">11 Views</span>
          </div>
        </a>

        <a
          href="/restaurant"
          class="flex items-center space-x-2 p-2.5 rounded-xl bg-orange-500/10 border border-orange-500/30 hover:bg-orange-500/20 transition active:scale-95"
        >
          <div class="w-7 h-7 rounded-lg bg-orange-500 text-slate-950 flex items-center justify-center font-bold">
            <i class="ph ph-fork-knife text-sm"></i>
          </div>
          <div class="text-left truncate">
            <span class="text-xs font-bold text-orange-300 block leading-tight">Restaurant</span>
            <span class="text-[10px] text-slate-400">Dockets & Menu</span>
          </div>
        </a>

        <a
          href="/bar"
          class="flex items-center space-x-2 p-2.5 rounded-xl bg-purple-500/10 border border-purple-500/30 hover:bg-purple-500/20 transition active:scale-95"
        >
          <div class="w-7 h-7 rounded-lg bg-purple-600 text-white flex items-center justify-center font-bold">
            <i class="ph ph-wine text-sm"></i>
          </div>
          <div class="text-left truncate">
            <span class="text-xs font-bold text-purple-300 block leading-tight">Bar Station</span>
            <span class="text-[10px] text-slate-400">Inventory & COGS</span>
          </div>
        </a>

        <a
          href="/kitchen"
          class="flex items-center space-x-2 p-2.5 rounded-xl bg-emerald-500/10 border border-emerald-500/30 hover:bg-emerald-500/20 transition active:scale-95"
        >
          <div class="w-7 h-7 rounded-lg bg-emerald-500 text-slate-950 flex items-center justify-center font-bold">
            <i class="ph ph-cooking-pot text-sm"></i>
          </div>
          <div class="text-left truncate">
            <span class="text-xs font-bold text-emerald-300 block leading-tight">Kitchen KDS</span>
            <span class="text-[10px] text-slate-400">Orders & Stock</span>
          </div>
        </a>

        <a
          href="/manager"
          class="flex items-center justify-between p-2.5 rounded-xl bg-indigo-500/10 border border-indigo-500/30 hover:bg-indigo-500/20 transition active:scale-95"
        >
          <div class="flex items-center space-x-2">
            <div class="w-7 h-7 rounded-lg bg-indigo-600 text-white flex items-center justify-center font-bold">
              <i class="ph ph-briefcase text-sm"></i>
            </div>
            <div class="text-left">
              <span class="text-xs font-bold text-indigo-300 block leading-tight">General Manager</span>
              <span class="text-[10px] text-slate-400">Audits & Reports</span>
            </div>
          </div>
          <span class="text-xs text-indigo-400 font-bold">&rarr;</span>
        </a>

        <a
          href="/book"
          class="flex items-center justify-between p-2.5 rounded-xl bg-gradient-to-r from-amber-500/20 to-amber-600/10 border border-amber-500/40 hover:bg-amber-500/30 transition active:scale-95 shadow-md shadow-amber-500/10"
        >
          <div class="flex items-center space-x-2">
            <div class="w-7 h-7 rounded-lg bg-amber-500 text-slate-950 flex items-center justify-center font-bold">
              <i class="ph ph-calendar-check text-sm"></i>
            </div>
            <div class="text-left">
              <span class="text-xs font-bold text-amber-300 block leading-tight">Guest Booking & Portal</span>
              <span class="text-[10px] text-slate-300">Book Room · Self-Service</span>
            </div>
          </div>
          <span class="text-xs text-amber-400 font-bold">&rarr;</span>
        </a>
      </div>
    </div>
    {/if}

    {#if errorMessage}
      <div class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs">
        {errorMessage}
      </div>
    {/if}

    <!-- Login Form -->
    <form onsubmit={handleLogin} class="space-y-4">
      {#if APP_ROLE === 'shared'}
      <div>
        <label for="login-role" class="block text-slate-400 text-xs font-semibold uppercase tracking-wider mb-1.5">Workstation Role</label>
        <select
          id="login-role"
          bind:value={selectedRole}
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500"
        >
          {#each roles as r}
            <option value={r.id}>{r.label}</option>
          {/each}
        </select>
      </div>
      {:else}
      <p class="text-xs text-slate-400">Sign in to the {roles.find((r) => r.id === APP_ROLE)?.label} app.</p>
      {/if}

      <div>
        <label for="login-user" class="block text-slate-400 text-xs font-semibold uppercase tracking-wider mb-1.5">Username</label>
        <input
          id="login-user"
          type="text"
          placeholder="e.g. reception1 or chef1"
          bind:value={username}
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600"
        />
      </div>

      <div>
        <label for="login-pass" class="block text-slate-400 text-xs font-semibold uppercase tracking-wider mb-1.5">Password</label>
        <div class="relative">
          <input
            id="login-pass"
            type={showPassword ? 'text' : 'password'}
            bind:value={password}
            placeholder="••••••••"
            class="w-full px-4 py-2.5 bg-slate-950 border border-slate-800 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600 pr-12"
          />
          <button
            type="button"
            onclick={() => (showPassword = !showPassword)}
            class="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-slate-400 hover:text-slate-200"
          >
            {showPassword ? 'Hide' : 'Show'}
          </button>
        </div>
      </div>

      <button
        type="submit"
        disabled={isLoading}
        class="w-full py-3.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-98 transition disabled:opacity-50 mt-2"
      >
        {isLoading ? 'Authenticating Station...' : 'Sign In to Station'}
      </button>
    </form>
  </div>
</div>
