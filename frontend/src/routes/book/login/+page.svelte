<script lang="ts">
  import { goto } from '$app/navigation';
  import { apiFetch } from '#lib/api';
  import { getGuestLandingPath } from '#lib/guestBookings';
  import { authStore } from '#lib/stores/authStore';

  let mode = $state<'login' | 'register' | 'resident'>('login');
  let isLoading = $state(false);
  let errorMsg = $state('');
  let successMsg = $state('');
  let showPassword = $state(false);

  // Login fields
  let loginUsername = $state('');
  let loginPassword = $state('');

  // Resident walk-in fields
  let resRoomNumber = $state('');
  let resCardNo = $state('');
  let resLastName = $state('');

  // Registration fields
  let regUsername = $state('');
  let regEmail = $state('');
  let regPassword = $state('');
  let regConfirm = $state('');
  let regFirstName = $state('');
  let regLastName = $state('');
  let regPhone = $state('');
  let regNationality = $state('Malawian');
  let regCity = $state('Lilongwe');
  let regPassportNo = $state('');
  let regPostalCode = $state('');
  let regAddress = $state('');

  async function handleLogin(e: Event) {
    e.preventDefault();
    if (!loginUsername.trim() || !loginPassword.trim()) {
      errorMsg = 'Please enter username and password.';
      return;
    }
    errorMsg = '';
    isLoading = true;
    try {
      const res = await apiFetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: loginUsername, password: loginPassword, role: 'guest' })
      });
      const data = await res.json();
      if (!res.ok) {
        errorMsg = data.error || 'Login failed';
        return;
      }
      authStore.login({
        userId: data.user?.user_id || '',
        token: data.token,
        username: data.user?.username || loginUsername,
        fullName: data.user?.full_name || loginUsername,
        role: 'guest'
      });
      goto(await getGuestLandingPath());
    } catch {
      errorMsg = 'Connection error. Please try again.';
    } finally {
      isLoading = false;
    }
  }

  async function handleResidentLogin(e: Event) {
    e.preventDefault();
    if (!resRoomNumber.trim()) {
      errorMsg = 'Please enter your Room Number.';
      return;
    }
    errorMsg = '';
    isLoading = true;
    try {
      const res = await apiFetch('/api/v1/guest/resident-login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          room_number: resRoomNumber.trim(),
          card_no: resCardNo.trim(),
          last_name: resLastName.trim()
        })
      });
      const data = await res.json();
      if (!res.ok) {
        errorMsg = data.error || 'Failed to authenticate resident access.';
        return;
      }
      authStore.login({
        userId: data.user?.user_id || 'resident',
        token: data.token,
        username: data.user?.username || `resident-${resRoomNumber}`,
        fullName: data.user?.full_name || `Resident Room ${resRoomNumber}`,
        role: 'guest'
      });
      successMsg = `Welcome, ${data.user?.full_name || 'Room ' + resRoomNumber}! Redirecting to Guest Portal...`;
      setTimeout(() => goto('/guest'), 800);
    } catch {
      errorMsg = 'Connection error. Please try again.';
    } finally {
      isLoading = false;
    }
  }

  async function handleRegister(e: Event) {
    e.preventDefault();
    errorMsg = '';

    if (!regUsername.trim() || !regEmail.trim() || !regPassword || !regFirstName.trim() || !regLastName.trim()) {
      errorMsg = 'Please fill in all required fields.';
      return;
    }
    if (regPassword !== regConfirm) {
      errorMsg = 'Passwords do not match.';
      return;
    }
    if (regPassword.length < 6) {
      errorMsg = 'Password must be at least 6 characters.';
      return;
    }

    isLoading = true;
    try {
      const res = await apiFetch('/api/v1/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: regUsername,
          email: regEmail,
          password: regPassword,
          first_name: regFirstName,
          last_name: regLastName,
          phone: regPhone,
          nationality: regNationality,
          city: regCity,
          passport_no: regPassportNo,
          postal_code: regPostalCode,
          address: regAddress
        })
      });
      const data = await res.json();
      if (!res.ok) {
        errorMsg = data.error || 'Registration failed';
        return;
      }
      // Auto-login after registration
      authStore.login({
        userId: data.user?.user_id || '',
        token: data.token,
        username: data.user?.username || regUsername,
        fullName: data.user?.full_name || `${regFirstName} ${regLastName}`,
        role: 'guest'
      });
      successMsg = 'Account created! Redirecting to booking portal...';
      setTimeout(() => goto('/book'), 1500);
    } catch {
      errorMsg = 'Connection error. Please try again.';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="max-w-md mx-auto space-y-6 pb-8">
  <!-- Header -->
  <div class="text-center space-y-2 pt-4">
    <div class="w-14 h-14 rounded-2xl bg-amber-500/10 border-2 border-amber-500/40 text-amber-400 mx-auto flex items-center justify-center">
      <i class="ph ph-user-circle text-3xl"></i>
    </div>
    <h1 class="text-xl font-bold text-white">Guest Account</h1>
    <p class="text-xs text-slate-400">Sign in or create an account to book your stay</p>
  </div>

  <!-- Tab Switcher -->
  <div class="flex bg-slate-900 rounded-2xl p-1 border border-slate-800">
    <button
      onclick={() => { mode = 'login'; errorMsg = ''; successMsg = ''; }}
      class="flex-1 py-2.5 rounded-xl text-xs font-bold transition {mode === 'login' ? 'bg-amber-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-white'}"
    >Sign In</button>
    <button
      onclick={() => { mode = 'register'; errorMsg = ''; successMsg = ''; }}
      class="flex-1 py-2.5 rounded-xl text-xs font-bold transition {mode === 'register' ? 'bg-amber-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-white'}"
    >Create Account</button>
    <button
      onclick={() => { mode = 'resident'; errorMsg = ''; successMsg = ''; }}
      class="flex-1 py-2.5 rounded-xl text-xs font-bold transition {mode === 'resident' ? 'bg-amber-500 text-slate-950 shadow-lg' : 'text-slate-400 hover:text-white'}"
    >Resident Sign-In</button>
  </div>

  {#if successMsg}
    <div class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-xs text-center font-bold">
      <i class="ph ph-check-circle text-sm"></i> {successMsg}
    </div>
  {/if}

  {#if errorMsg}
    <div class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs">
      <i class="ph ph-warning text-sm"></i> {errorMsg}
    </div>
  {/if}

  <!-- Login Form -->
  {#if mode === 'login'}
    <form onsubmit={handleLogin} class="bg-slate-900 border border-slate-800 rounded-3xl p-5 space-y-4 shadow-xl">
      <div>
        <label for="g-user" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Username</label>
        <input id="g-user" type="text" bind:value={loginUsername} placeholder="Your username"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>
      <div>
        <label for="g-pass" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Password</label>
        <div class="relative">
          <input id="g-pass" type={showPassword ? 'text' : 'password'} bind:value={loginPassword} placeholder="••••••••"
            class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600 pr-12" />
          <button type="button" onclick={() => (showPassword = !showPassword)} class="absolute right-3 top-1/2 -translate-y-1/2 text-xs text-slate-400 hover:text-slate-200">
            {showPassword ? 'Hide' : 'Show'}
          </button>
        </div>
      </div>
      <button type="submit" disabled={isLoading}
        class="w-full py-3.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-95 transition disabled:opacity-50">
        {isLoading ? 'Signing in...' : 'Sign In'}
      </button>
    </form>

  <!-- Resident Walk-in Form -->
  {:else if mode === 'resident'}
    <form onsubmit={handleResidentLogin} class="bg-slate-900 border border-slate-800 rounded-3xl p-5 space-y-4 shadow-xl">
      <div class="p-3 bg-amber-500/10 border border-amber-500/20 rounded-2xl text-[11px] text-amber-300">
        <i class="ph ph-info mr-1 text-sm"></i> Checked in at the Front Desk? Sign in directly using your Room Number and Registration Card Reference.
      </div>
      <div>
        <label for="res-room" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Room Number *</label>
        <input id="res-room" type="text" bind:value={resRoomNumber} placeholder="e.g. 101" required
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600 font-mono" />
      </div>
      <div>
        <label for="res-card" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Card Reference</label>
        <input id="res-card" type="text" bind:value={resCardNo} placeholder="e.g. CARD-101-1234 (from your bill/card)"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600 font-mono" />
      </div>
      <div>
        <label for="res-last" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Last Name (Optional Fallback)</label>
        <input id="res-last" type="text" bind:value={resLastName} placeholder="e.g. Banda"
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>
      <button type="submit" disabled={isLoading}
        class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-95 transition disabled:opacity-50">
        {isLoading ? 'Verifying Resident...' : 'Enter Resident Portal'}
      </button>
    </form>

  <!-- Registration Form -->
  {:else}
    <form onsubmit={handleRegister} class="bg-slate-900 border border-slate-800 rounded-3xl p-5 space-y-4 shadow-xl">
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="r-fname" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">First Name *</label>
          <input id="r-fname" type="text" bind:value={regFirstName} placeholder="John" required
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
        </div>
        <div>
          <label for="r-lname" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Last Name *</label>
          <input id="r-lname" type="text" bind:value={regLastName} placeholder="Doe" required
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
        </div>
      </div>
      <div>
        <label for="r-user" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Username *</label>
        <input id="r-user" type="text" bind:value={regUsername} placeholder="Choose a username" required
          class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>
      <div>
        <label for="r-email" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Email Address *</label>
        <input id="r-email" type="email" bind:value={regEmail} placeholder="john@example.com" required
          class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>
      <div>
        <label for="r-phone" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Phone Number *</label>
        <input id="r-phone" type="tel" bind:value={regPhone} placeholder="+265 999 000 000" required
          class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="r-nat" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Nationality</label>
          <input id="r-nat" type="text" bind:value={regNationality} placeholder="Malawian"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
        </div>
        <div>
          <label for="r-city" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">City</label>
          <input id="r-city" type="text" bind:value={regCity} placeholder="Lilongwe"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
        </div>
      </div>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label for="r-pass" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">National ID / Passport</label>
          <input id="r-pass" type="text" bind:value={regPassportNo} placeholder="e.g. MW-PP-12345"
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
        </div>
        <div>
          <label for="r-post" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Postal Code</label>
          <input id="r-post" type="text" bind:value={regPostalCode} placeholder="P.O. Box ..."
            class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
        </div>
      </div>
      <div>
        <label for="r-addr" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Physical Address</label>
        <input id="r-addr" type="text" bind:value={regAddress} placeholder="Area 10, Lilongwe"
          class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>
      <div>
        <label for="r-pass" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Password *</label>
        <input id="r-pass" type="password" bind:value={regPassword} placeholder="Min 6 characters" required
          class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>
      <div>
        <label for="r-confirm" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1.5">Confirm Password *</label>
        <input id="r-confirm" type="password" bind:value={regConfirm} placeholder="Repeat password" required
          class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
      </div>

      <button type="submit" disabled={isLoading}
        class="w-full py-3.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-95 transition disabled:opacity-50">
        {isLoading ? 'Creating Account...' : 'Create Guest Account'}
      </button>
    </form>
  {/if}

  <!-- Link back to booking -->
  <div class="text-center">
    <a href="/book" class="text-xs text-amber-400 hover:text-amber-300 font-semibold">
      <i class="ph ph-arrow-left text-sm"></i> Back to Room Browser
    </a>
  </div>
</div>
