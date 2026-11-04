<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { authStore } from '#lib/stores/authStore';
  import { guestStore } from '#lib/stores/guestStore';
  import { goto } from '$app/navigation';

  let auth = $state<{ fullName: string; token: string; isLoggedIn: boolean }>({ fullName: '', token: '', isLoggedIn: false });
  authStore.subscribe(v => { auth = v; });

  const orders = guestStore.orders;
  const services = guestStore.services;

  // Active items count
  let activeOrders = $derived(
    $orders.filter(o => o.status !== 'Delivered')
  );

  let activeRequests = $derived(
    $services.filter(s => s.status !== 'Completed')
  );

  // Booking data from API
  let bookings = $state<Array<{
    card_no: string;
    room_number: string;
    room_type: string;
    bed_type: string;
    room_rate: number;
    arrival_date: string;
    departure_date: string;
    nights: number;
    payment_mode: string;
    room_charges: number;
    incidental_charges: number;
    total_balance: number;
    folio_status: string;
    room_status: string;
    checked_in: boolean;
  }>>([]);

  let loadingBookings = $state(true);
  let isCheckingIn = $state(false);
  let checkInBanner = $state('');
  let dndActive = $state(false);
  let makeupActive = $state(false);

  let activeBooking = $derived(
    bookings.find(b => b.folio_status !== 'Closed') || null
  );
  let hasActiveBooking = $derived(activeBooking !== null);
  let isCheckedIn = $derived(hasActiveBooking && activeBooking?.checked_in === true);

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }

  $effect(() => {
    if (auth.isLoggedIn) {
      fetchBookings();
    }
  });

  async function fetchBookings() {
    loadingBookings = true;
    try {
      const res = await apiFetch('/api/v1/guest/bookings', {
        headers: {
          'Authorization': `Bearer ${auth.token}`
        }
      });
      if (res.ok) {
        const data = await res.json();
        bookings = data;
        if (data.length > 0) {
          // Sync with guestStore
          const b = data[0];
          guestStore.stay.switchStay(b.room_number, auth.fullName || 'Resident');
        }
      }
    } catch {
      // offline fallback
    } finally {
      loadingBookings = false;
    }
  }

  async function handleCheckIn() {
    if (!activeBooking) return;
    isCheckingIn = true;
    checkInBanner = '';
    try {
      const res = await apiFetch('/api/v1/guest/checkin', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${auth.token}`,
          'Content-Type': 'application/json'
        }
      });
      const data = await res.json();
      if (res.ok) {
        checkInBanner = data.message || 'Checked in successfully!';
        // Update local status
        bookings = bookings.map(b => b.card_no === activeBooking?.card_no ? { ...b, room_status: 'Occupied' } : b);
      } else {
        checkInBanner = data.error || 'Check-in must be completed at the Front Desk Reception.';
      }
    } catch {
      checkInBanner = 'Online self-check-in is disabled. Please check in with the Front Desk Receptionist.';
    } finally {
      isCheckingIn = false;
      setTimeout(() => { checkInBanner = ''; }, 5000);
    }
  }
</script>

<div class="space-y-4 pb-4">
  <!-- Personalized Welcome Header -->
  <div class="pt-1 flex items-start justify-between">
    <div>
      <p class="text-xs text-amber-400 font-semibold tracking-wide uppercase">Welcome to La Luna Hotel</p>
      <h1 class="text-2xl font-bold text-white tracking-tight">{auth.fullName || 'Guest'}</h1>
      {#if activeBooking}
        <p class="text-xs text-slate-400">
          {activeBooking.room_type} · Room {activeBooking.room_number}
        </p>
      {:else}
        <p class="text-xs text-slate-400">Guest Resident Account</p>
      {/if}
    </div>

    {#if hasActiveBooking}
      <button
        type="button"
        disabled
        class="text-[11px] px-3 py-1.5 rounded-xl bg-slate-800 text-slate-500 font-bold flex items-center space-x-1 cursor-not-allowed opacity-60 border border-slate-700"
        title="You already have an active room reservation (Room {activeBooking?.room_number})"
      >
        <i class="ph ph-lock text-xs"></i>
        <span>Room Booked</span>
      </button>
    {:else}
      <a
        href="/book"
        class="text-[11px] px-3 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold active:scale-95 transition flex items-center space-x-1 shadow-md shadow-amber-500/20"
      >
        <i class="ph ph-plus text-xs"></i>
        <span>Book Room</span>
      </a>
    {/if}
  </div>

  <!-- Check-in notification feedback -->
  {#if checkInBanner}
    <div class="p-3.5 rounded-2xl bg-emerald-500/20 border border-emerald-500/40 text-emerald-300 text-xs text-center font-bold flex items-center justify-center space-x-2 animate-bounce">
      <i class="ph ph-check-circle text-base"></i>
      <span>{checkInBanner}</span>
    </div>
  {/if}

  {#if loadingBookings}
    <div class="flex justify-center py-12">
      <div class="w-8 h-8 border-2 border-amber-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

  <!-- NO ACTIVE BOOKING STATE: Feature lock-down -->
  {:else if !hasActiveBooking}
    <div class="bg-gradient-to-br from-slate-900 to-slate-950 border border-amber-500/30 rounded-3xl p-6 text-center space-y-4 shadow-xl">
      <div class="w-16 h-16 rounded-2xl bg-amber-500/10 border-2 border-amber-500/30 text-amber-400 mx-auto flex items-center justify-center">
        <i class="ph ph-calendar-blank text-3xl"></i>
      </div>
      <div>
        <h2 class="text-lg font-bold text-white font-['Oxanium',sans-serif]">No Active Room Booking</h2>
        <p class="text-xs text-slate-400 mt-1.5 max-w-xs mx-auto leading-relaxed">
          Your account is not currently tied to any active reservation. Most hotel services are locked until a room is booked.
        </p>
      </div>

      <a
        href="/book"
        class="inline-flex items-center space-x-2 px-6 py-3.5 rounded-2xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-95 transition"
      >
        <i class="ph ph-calendar-check text-base"></i>
        <span>Browse Available Rooms & Book</span>
      </a>
    </div>

    <!-- Disabled/Locked Services Preview -->
    <div>
      <div class="flex items-center justify-between mb-2.5">
        <h2 class="text-xs font-bold text-slate-500 uppercase tracking-wider">Hotel Services (Requires Booking)</h2>
        <span class="text-[10px] text-slate-500 flex items-center space-x-1">
          <i class="ph ph-lock text-xs"></i>
          <span>Locked</span>
        </span>
      </div>

      <div class="grid grid-cols-2 gap-3 opacity-50 pointer-events-none select-none">
        <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between h-24">
          <div class="w-8 h-8 rounded-lg bg-orange-500/20 text-orange-400 flex items-center justify-center">
            <i class="ph ph-fork-knife text-base"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white">In-Room Dining</p>
            <p class="text-[10px] text-slate-400">Locked</p>
          </div>
        </div>

        <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between h-24">
          <div class="w-8 h-8 rounded-lg bg-blue-500/20 text-blue-400 flex items-center justify-center">
            <i class="ph ph-receipt text-base"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white">Room Bill & Folio</p>
            <p class="text-[10px] text-slate-400">Locked</p>
          </div>
        </div>

        <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between h-24">
          <div class="w-8 h-8 rounded-lg bg-emerald-500/20 text-emerald-400 flex items-center justify-center">
            <i class="ph ph-bell-ringing text-base"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white">Room Services</p>
            <p class="text-[10px] text-slate-400">Locked</p>
          </div>
        </div>

        <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between h-24">
          <div class="w-8 h-8 rounded-lg bg-purple-500/20 text-purple-400 flex items-center justify-center">
            <i class="ph ph-sparkle text-base"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white">Housekeeping</p>
            <p class="text-[10px] text-slate-400">Locked</p>
          </div>
        </div>
      </div>
    </div>

  <!-- ACTIVE BOOKING STATE -->
  {:else if activeBooking}
    <!-- Active Room Stay Card -->
    <div class="bg-gradient-to-br from-amber-600/15 via-slate-900 to-slate-900 border border-amber-500/30 rounded-3xl p-5 shadow-xl relative overflow-hidden">
      <div class="absolute -right-6 -bottom-6 w-32 h-32 rounded-full bg-amber-500/5 blur-xl pointer-events-none"></div>

      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center space-x-2">
          <div class="w-8 h-8 rounded-xl bg-amber-500 text-slate-950 flex items-center justify-center font-bold">
            <i class="ph ph-bed text-lg"></i>
          </div>
          <div>
            <span class="text-xs font-bold text-white uppercase tracking-wider block">Room {activeBooking.room_number}</span>
            <span class="text-[10px] text-slate-400">{activeBooking.room_type} · {activeBooking.bed_type} Bed</span>
          </div>
        </div>

        <span class="text-[10px] px-2.5 py-1 rounded-full border font-bold uppercase tracking-wider {isCheckedIn ? 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40' : 'bg-amber-500/20 text-amber-300 border-amber-500/40'}">
          {isCheckedIn ? 'Checked In' : 'Pending Arrival'}
        </span>
      </div>

      <div class="flex items-center justify-between pt-2 text-xs">
        <div>
          <p class="text-slate-400">Card Ref</p>
          <p class="text-sm font-bold text-white font-['Oxanium',sans-serif]">{activeBooking.card_no}</p>
        </div>
        <div class="text-right">
          <p class="text-slate-400">Nightly Rate</p>
          <p class="text-sm font-bold text-amber-400 font-['Oxanium',sans-serif]">{fmtMwk(activeBooking.room_rate)}</p>
        </div>
      </div>

      <!-- Front Desk Check-in Required Notice -->
      {#if !isCheckedIn}
        <div class="mt-4 pt-3 border-t border-slate-800">
          <div class="bg-amber-500/10 border border-amber-500/30 rounded-xl p-3 text-center space-y-1">
            <div class="flex items-center justify-center space-x-1.5 text-amber-400 text-xs font-bold">
              <i class="ph ph-clock text-base"></i>
              <span>Pending Front Desk Check-In</span>
            </div>
            <p class="text-[10px] text-slate-400 leading-relaxed">
              Your booking is confirmed. Please visit the Front Desk Reception with your ID to complete check-in. Hotel services will be activated once the receptionist processes your arrival.
            </p>
          </div>
        </div>
      {/if}
    </div>

    {#if isCheckedIn}
    <!-- Room Privacy & Housekeeping Controls (Only available when checked in) -->
    <div class="grid grid-cols-2 gap-3">
      <!-- DND Toggle -->
      <button
        onclick={() => (dndActive = !dndActive)}
        class="bg-slate-900 border rounded-2xl p-4 text-left transition active:scale-95 shadow-md {dndActive ? 'border-rose-500/60 bg-rose-950/20' : 'border-slate-800 hover:border-slate-700'}"
      >
        <div class="flex items-center justify-between mb-2">
          <i class="ph ph-prohibit text-xl {dndActive ? 'text-rose-400' : 'text-slate-400'}"></i>
          <span class="w-2.5 h-2.5 rounded-full {dndActive ? 'bg-rose-400 animate-ping' : 'bg-slate-700'}"></span>
        </div>
        <p class="text-xs font-bold text-white">Do Not Disturb</p>
        <p class="text-[10px] text-slate-400 mt-0.5">
          {dndActive ? 'Active — Privacy on' : 'Off — Staff may assist'}
        </p>
      </button>

      <!-- Make Up Room Toggle -->
      <button
        onclick={() => (makeupActive = !makeupActive)}
        class="bg-slate-900 border rounded-2xl p-4 text-left transition active:scale-95 shadow-md {makeupActive ? 'border-emerald-500/60 bg-emerald-950/20' : 'border-slate-800 hover:border-slate-700'}"
      >
        <div class="flex items-center justify-between mb-2">
          <i class="ph ph-sparkle text-xl {makeupActive ? 'text-emerald-400' : 'text-slate-400'}"></i>
          <span class="w-2.5 h-2.5 rounded-full {makeupActive ? 'bg-emerald-400 animate-ping' : 'bg-slate-700'}"></span>
        </div>
        <p class="text-xs font-bold text-white">Make Up Room</p>
        <p class="text-[10px] text-slate-400 mt-0.5">
          {makeupActive ? 'Queued for cleaning' : 'Tap to request service'}
        </p>
      </button>
    </div>

    <!-- Active Orders / Requests Tracker -->
    {#if activeOrders.length > 0 || activeRequests.length > 0}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-1.5">
            <span class="w-2 h-2 rounded-full bg-amber-400 animate-ping"></span>
            <h2 class="text-xs font-bold text-white uppercase tracking-wider">Active Requests</h2>
          </div>
          <span class="text-[10px] text-slate-400">{activeOrders.length + activeRequests.length} Active</span>
        </div>

        {#each activeOrders as ord}
          <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-3 flex items-center justify-between text-xs">
            <div>
              <div class="flex items-center space-x-1.5">
                <i class="ph ph-cooking-pot text-amber-400 text-sm"></i>
                <span class="font-bold text-white">Dining ({ord.orderNumber})</span>
              </div>
              <p class="text-slate-400 text-[11px] mt-0.5">
                {ord.items.length} item{ord.items.length !== 1 ? 's' : ''} · {fmtMwk(ord.totalAmount)}
              </p>
            </div>
            <span class="text-[10px] px-2 py-0.5 rounded-full bg-amber-500/20 text-amber-300 border border-amber-500/40 font-bold uppercase">
              {ord.status}
            </span>
          </div>
        {/each}

        {#each activeRequests as srv}
          <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-3 flex items-center justify-between text-xs">
            <div>
              <div class="flex items-center space-x-1.5">
                <i class="ph {srv.icon} text-emerald-400 text-sm"></i>
                <span class="font-bold text-white">{srv.type}</span>
              </div>
              <p class="text-slate-400 text-[11px] mt-0.5">{srv.details}</p>
            </div>
            <span class="text-[10px] px-2 py-0.5 rounded-full bg-emerald-500/20 text-emerald-300 border border-emerald-500/40 font-bold uppercase">
              {srv.status}
            </span>
          </div>
        {/each}
      </div>
    {/if}

    <!-- Quick Action Grid -->
    <div>
      <h2 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2.5">Resident Services</h2>
      <div class="grid grid-cols-2 gap-3">
        <!-- In-Room Dining -->
        <a
          href="/guest/dining"
          class="bg-slate-900 border border-slate-800 hover:border-amber-500/40 rounded-2xl p-4 text-left transition active:scale-95 shadow-md group flex flex-col justify-between h-28"
        >
          <div class="w-10 h-10 rounded-xl bg-orange-500/20 text-orange-400 flex items-center justify-center">
            <i class="ph ph-fork-knife text-xl"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white group-hover:text-amber-300 transition-colors">Dining & Breakfast</p>
            <p class="text-[10px] text-slate-400">Order meals to Room {activeBooking.room_number}</p>
          </div>
        </a>

        <!-- My Bill & Folio -->
        <a
          href="/guest/folio"
          class="bg-slate-900 border border-slate-800 hover:border-amber-500/40 rounded-2xl p-4 text-left transition active:scale-95 shadow-md group flex flex-col justify-between h-28"
        >
          <div class="w-10 h-10 rounded-xl bg-blue-500/20 text-blue-400 flex items-center justify-center">
            <i class="ph ph-receipt text-xl"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white group-hover:text-amber-300 transition-colors">Room Folio</p>
            <p class="text-[10px] text-slate-400">View bill & charges</p>
          </div>
        </a>

        <!-- Housekeeping & Concierge -->
        <a
          href="/guest/services"
          class="bg-slate-900 border border-slate-800 hover:border-amber-500/40 rounded-2xl p-4 text-left transition active:scale-95 shadow-md group flex flex-col justify-between h-28"
        >
          <div class="w-10 h-10 rounded-xl bg-emerald-500/20 text-emerald-400 flex items-center justify-center">
            <i class="ph ph-bell-ringing text-xl"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white group-hover:text-amber-300 transition-colors">Amenities & Towels</p>
            <p class="text-[10px] text-slate-400">Request guest essentials</p>
          </div>
        </a>

        <!-- Hotel Guide -->
        <a
          href="/guest/guide"
          class="bg-slate-900 border border-slate-800 hover:border-amber-500/40 rounded-2xl p-4 text-left transition active:scale-95 shadow-md group flex flex-col justify-between h-28"
        >
          <div class="w-10 h-10 rounded-xl bg-purple-500/20 text-purple-400 flex items-center justify-center">
            <i class="ph ph-compass text-xl"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white group-hover:text-amber-300 transition-colors">Hotel Guide</p>
            <p class="text-[10px] text-slate-400">Hours, policies & contacts</p>
          </div>
        </a>
      </div>
    </div>

    {:else}
      <!-- Not checked in: locked services -->
      <div class="bg-slate-900/60 border border-amber-500/20 rounded-2xl p-5 text-center space-y-3">
        <i class="ph ph-lock-simple text-3xl text-amber-400 block"></i>
        <h3 class="text-sm font-bold text-white">Services Locked — Check-In Required</h3>
        <p class="text-[11px] text-slate-400 max-w-xs mx-auto leading-relaxed">
          Room services, dining, folio, and housekeeping features will become available once the Front Desk Receptionist completes your check-in.
        </p>
      </div>

      <!-- Stay Information still visible even before check-in -->
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-2.5 text-xs">
        <div class="flex items-center justify-between">
          <span class="text-slate-400">Stay Duration</span>
          <span class="font-bold text-white">{activeBooking.nights} Night{activeBooking.nights !== 1 ? 's' : ''} ({activeBooking.arrival_date} – {activeBooking.departure_date})</span>
        </div>
        <div class="flex items-center justify-between border-t border-slate-800/80 pt-2">
          <span class="text-slate-400">Total Charges</span>
          <span class="font-bold text-amber-300">{fmtMwk(activeBooking.total_balance)}</span>
        </div>
        <div class="flex items-center justify-between border-t border-slate-800/80 pt-2">
          <span class="text-slate-400">Checkout Time</span>
          <span class="font-bold text-amber-300">10:00 AM Departure</span>
        </div>
        <div class="flex items-center justify-between border-t border-slate-800/80 pt-2">
          <span class="text-slate-400">Reception Desk</span>
          <a href="tel:0" class="font-bold text-indigo-400 hover:underline">Ext. 0 (24/7)</a>
        </div>
      </div>
    {/if}
  {/if}
</div>
