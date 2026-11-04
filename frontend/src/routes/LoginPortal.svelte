<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { apiFetch } from '#lib/api';

  // Login credentials and options
  let username = '';
  let password = '';
  let selectedRole = 'frontdesk';
  let rememberMe = false;
  let showPassword = false;

  // UI state management
  let isLoading = false;
  let errorMessage = '';

  const dispatch = createEventDispatcher();

  // Role metadata definition for RBAC context
  const roles = [
    { id: 'frontdesk', label: 'Front Desk Officer' },
    { id: 'restaurant', label: 'Restaurant Attendant' },
    { id: 'bar', label: 'Bar Attendant' },
    { id: 'kitchen', label: 'Head Chef / Kitchen' },
    { id: 'manager', label: 'Hotel Manager' }
  ];

  /**
   * Submits credentials to the Golang backend auth endpoint.
   */
  async function handleLogin() {
    if (!username.trim() || !password.trim()) {
      errorMessage = 'Please enter both username and password.';
      return;
    }

    errorMessage = '';
    isLoading = true;

    try {
      // POST request to Golang backend route
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
        const errorData = await response.json();
        throw new Error(errorData.message || 'Invalid credentials or unauthorized role access.');
      }

      const data = await response.json();

      // Notify parent component of successful login
      dispatch('loginSuccess', {
        token: data.token,
        user: data.user,
        role: selectedRole
      });

    } catch (err: any) {
      errorMessage = err.message || 'Authentication failed. Please check connection or credentials.';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-slate-900 text-slate-100 flex flex-col justify-center items-center p-4 sm:p-6 lg:p-8">
  <!-- Brand Header -->
  <div class="sm:mx-auto sm:w-full sm:max-w-md text-center mb-6">
    <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-amber-500/10 border border-amber-500/30 text-amber-400 mb-3 shadow-lg shadow-amber-500/5">
      <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
      </svg>
    </div>
    <h1 class="text-3xl font-extrabold tracking-tight text-white">La Luna Hotel</h1>
    <p class="mt-1 text-sm text-slate-400">Management System & Portal Access</p>
  </div>

  <!-- Login Card Container -->
  <div class="w-full sm:max-w-md bg-slate-800/80 backdrop-blur-md border border-slate-700/60 rounded-2xl shadow-2xl p-6 sm:p-8">
    <div class="mb-6">
      <h2 class="text-xl font-bold text-slate-100">Portal Authentication</h2>
      <p class="text-xs text-slate-400 mt-1">Select your operational role and credentials to proceed.</p>
    </div>

    <!-- Error Alert Box -->
    {#if errorMessage}
      <div class="mb-5 p-3.5 rounded-lg bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs flex items-start space-x-2">
        <svg class="w-4 h-4 text-rose-400 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{errorMessage}</span>
      </div>
    {/if}

    <form on:submit|preventDefault={handleLogin} class="space-y-5">
      <!-- Role Selector -->
      <div>
        <label for="role" class="block text-xs font-semibold uppercase tracking-wider text-slate-300 mb-2">
          Operational Role
        </label>
        <div class="relative">
          <select
            id="role"
            bind:value={selectedRole}
            class="w-full bg-slate-900/90 border border-slate-700 rounded-xl px-4 py-2.5 text-sm text-slate-200 focus:outline-none focus:ring-2 focus:ring-amber-500/50 focus:border-amber-500 transition-all appearance-none cursor-pointer"
          >
            {#each roles as role}
              <option value={role.id} class="bg-slate-900 text-slate-200">{role.label}</option>
            {/each}
          </select>
          <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3 text-slate-400">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Username / Staff ID Input -->
      <div>
        <label for="username" class="block text-xs font-semibold uppercase tracking-wider text-slate-300 mb-2">
          Username / Staff ID
        </label>
        <div class="relative">
          <input
            id="username"
            type="text"
            bind:value={username}
            placeholder="e.g., frontdesk"
            required
            class="w-full bg-slate-900/90 border border-slate-700 rounded-xl px-4 py-2.5 pl-10 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-amber-500/50 focus:border-amber-500 transition-all"
          />
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Password Input -->
      <div>
        <label for="password" class="block text-xs font-semibold uppercase tracking-wider text-slate-300 mb-2">
          Password
        </label>
        <div class="relative">
          <input
            id="password"
            type={showPassword ? 'text' : 'password'}
            bind:value={password}
            placeholder="••••••••"
            required
            class="w-full bg-slate-900/90 border border-slate-700 rounded-xl px-4 py-2.5 pl-10 pr-10 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-amber-500/50 focus:border-amber-500 transition-all"
          />
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-500">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </div>
          <button
            type="button"
            on:click={() => (showPassword = !showPassword)}
            class="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-500 hover:text-slate-300 transition-colors"
          >
            {#if showPassword}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858-5.908a8.88 8.88 0 012.122-.381c4.478 0 8.268 2.943 9.542 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21M3 3l18 18" />
              </svg>
            {:else}
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
            {/if}
          </button>
        </div>
      </div>

      <!-- Remember Me Option -->
      <div class="flex items-center justify-between text-xs">
        <label class="flex items-center space-x-2 cursor-pointer text-slate-400 hover:text-slate-200 transition-colors">
          <input
            type="checkbox"
            bind:checked={rememberMe}
            class="w-4 h-4 rounded border-slate-700 bg-slate-900 text-amber-500 focus:ring-amber-500/40 focus:ring-offset-slate-800"
          />
          <span>Keep session active</span>
        </label>
      </div>

      <!-- Submit Button -->
      <button
        type="submit"
        disabled={isLoading}
        class="w-full bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold py-2.5 px-4 rounded-xl shadow-lg shadow-amber-500/20 hover:shadow-amber-500/30 transition-all flex items-center justify-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
      >
        {#if isLoading}
          <svg class="animate-spin h-4 w-4 text-slate-950" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span>Authenticating...</span>
        {:else}
          <span>Sign In to Dashboard</span>
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
          </svg>
        {/if}
      </button>
    </form>
  </div>
</div>