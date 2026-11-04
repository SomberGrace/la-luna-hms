<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { authStore } from '#lib/stores/authStore';
  import { goto } from '$app/navigation';

  let rooms = $state<Array<{
    room_number: string;
    room_type: string;
    bed_type: string;
    room_rate: number;
    status: string;
    description: string;
    image_url: string;
  }>>([]);
  let loading = $state(true);
  let error = $state('');
  let selectedRoom = $state<typeof rooms[0] | null>(null);
  let showBookingModal = $state(false);

  // Booking form fields
  let firstName = $state('');
  let lastName = $state('');
  let phone = $state('');
  let email = $state('');
  let nationality = $state('Malawian');
  let city = $state('');
  let passportNo = $state('');
  let postalCode = $state('');
  let address = $state('');
  let arrivalDate = $state(new Date().toISOString().split('T')[0]);
  let departureDate = $state(new Date(Date.now() + 86400000).toISOString().split('T')[0]);
  let paymentMode = $state('Cash');
  let isBooking = $state(false);
  let bookingResult = $state<{ success: boolean; message: string; cardNo?: string } | null>(null);
  let profileLoaded = $state(false);
  let showEditDetails = $state(false);
  let existingBooking = $state<{ room_number: string; card_no: string } | null>(null);

  let auth = $state<{ isLoggedIn: boolean; token: string; role: string; fullName?: string }>({ isLoggedIn: false, token: '', role: '' });
  authStore.subscribe(v => { auth = v; });

  let nights = $derived(Math.max(1, Math.ceil((new Date(departureDate).getTime() - new Date(arrivalDate).getTime()) / 86400000)));
  let totalCost = $derived(selectedRoom ? selectedRoom.room_rate * nights : 0);

  function fmtMwk(n: number) {
    return 'MWK ' + n.toLocaleString();
  }

  // Room type visual config
  function getRoomStyle(type: string) {
    switch (type) {
      case 'Standard': return { gradient: 'from-blue-500/20 to-slate-900', border: 'border-blue-500/30', icon: 'ph-bed', color: 'text-blue-400' };
      case 'Executive': return { gradient: 'from-amber-500/20 to-slate-900', border: 'border-amber-500/30', icon: 'ph-crown', color: 'text-amber-400' };
      case 'Apartment': return { gradient: 'from-emerald-500/20 to-slate-900', border: 'border-emerald-500/30', icon: 'ph-buildings', color: 'text-emerald-400' };
      case 'Sector Suite': return { gradient: 'from-purple-500/20 to-slate-900', border: 'border-purple-500/30', icon: 'ph-star', color: 'text-purple-400' };
      default: return { gradient: 'from-slate-800 to-slate-900', border: 'border-slate-700', icon: 'ph-bed', color: 'text-slate-400' };
    }
  }

  function getRoomImage(room: { room_type: string; image_url?: string }) {
    if (room.image_url && room.image_url.startsWith('http')) return room.image_url;
    switch (room.room_type) {
      case 'Executive':
        return 'https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?auto=format&fit=crop&w=800&q=80';
      case 'Apartment':
        return 'https://images.unsplash.com/photo-1591088398332-8a7791972843?auto=format&fit=crop&w=800&q=80';
      case 'Sector Suite':
        return 'https://images.unsplash.com/photo-1618773928121-c32242e63f39?auto=format&fit=crop&w=800&q=80';
      default:
        return 'https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=800&q=80';
    }
  }

  $effect(() => {
    fetchRooms();
    if (auth.isLoggedIn && auth.role === 'guest' && !profileLoaded) {
      fetchProfile();
    }
  });

  async function fetchProfile() {
    try {
      const res = await apiFetch('/api/v1/guest/profile', {
        headers: { 'Authorization': `Bearer ${auth.token}` }
      });
      if (res.ok) {
        const p = await res.json();
        if (p.first_name) firstName = p.first_name;
        if (p.last_name) lastName = p.last_name;
        if (p.contact_no && p.contact_no !== 'N/A') phone = p.contact_no;
        if (p.email) email = p.email;
        if (p.nationality) nationality = p.nationality;
        if (p.city) city = p.city;
        if (p.passport_no) passportNo = p.passport_no;
        if (p.postal_code) postalCode = p.postal_code;
        if (p.address) address = p.address;
        profileLoaded = true;
      }

      const bRes = await apiFetch('/api/v1/guest/bookings', {
        headers: { 'Authorization': `Bearer ${auth.token}` }
      });
      if (bRes.ok) {
        const bData = await bRes.json();
        if (Array.isArray(bData)) {
          const active = bData.find((b: any) => b.folio_status !== 'Closed');
          existingBooking = active ? { room_number: active.room_number, card_no: active.card_no } : null;
        }
      }
    } catch {
      // Profile fetch is best-effort
    }
  }

  async function fetchRooms() {
    loading = true;
    error = '';
    try {
      const res = await apiFetch('/api/v1/guest/rooms');
      if (!res.ok) throw new Error('Failed to load rooms');
      rooms = await res.json();
    } catch (e: any) {
      error = e.message || 'Could not load rooms';
      rooms = [];
    } finally {
      loading = false;
    }
  }

  function openBooking(room: typeof rooms[0]) {
    selectedRoom = room;
    showBookingModal = true;
    bookingResult = null;
  }

  async function submitBooking(e: Event) {
    e.preventDefault();
    if (!selectedRoom) return;

    if (!auth.isLoggedIn || auth.role !== 'guest') {
      bookingResult = { success: false, message: 'Please log in or create a guest account first.' };
      return;
    }

    isBooking = true;
    bookingResult = null;
    try {
      const res = await apiFetch('/api/v1/guest/book', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${auth.token}`
        },
        body: JSON.stringify({
          room_number: selectedRoom.room_number,
          arrival_date: arrivalDate,
          departure_date: departureDate,
          payment_mode: paymentMode,
          first_name: firstName,
          last_name: lastName,
          contact_no: phone,
          email: email,
          nationality: nationality,
          city: city,
          passport_no: passportNo,
          postal_code: postalCode,
          address: address
        })
      });

      const data = await res.json();
      if (res.ok) {
        bookingResult = { success: true, message: `Booking confirmed! Card No: ${data.card_no}. ${data.nights} night(s) totalling ${fmtMwk(data.total)}.`, cardNo: data.card_no };
        fetchRooms(); // refresh availability
      } else {
        bookingResult = { success: false, message: data.error || 'Booking failed' };
      }
    } catch {
      bookingResult = { success: false, message: 'Network error. Please try again.' };
    } finally {
      isBooking = false;
    }
  }
</script>

<div class="space-y-6 pb-8">
  <!-- Hero Section -->
  <div class="text-center space-y-2 py-4">
    <h1 class="text-2xl sm:text-3xl font-bold text-white font-['Oxanium',sans-serif] tracking-wide">Book Your Stay</h1>
    <p class="text-sm text-slate-400 max-w-md mx-auto">Browse our rooms and reserve your perfect stay at La Luna Hotel, Lilongwe</p>

    {#if !auth.isLoggedIn || auth.role !== 'guest'}
      <div class="mt-3 p-3 rounded-2xl bg-amber-500/10 border border-amber-500/30 text-amber-300 text-xs max-w-md mx-auto">
        <i class="ph ph-info text-sm"></i>
        <a href="/book/login" class="underline font-bold hover:text-amber-200">Log in or create an account</a> to book a room
      </div>
    {:else}
      <p class="text-xs text-emerald-400 font-semibold"><i class="ph ph-check-circle text-sm"></i> Logged in as {auth.fullName}</p>
      {#if existingBooking}
        <div class="mt-3 p-3.5 rounded-2xl bg-amber-500/10 border border-amber-500/30 text-amber-300 text-xs max-w-md mx-auto space-y-1">
          <p class="font-bold flex items-center justify-center space-x-1.5 text-amber-400">
            <i class="ph ph-bell-simple-ringing text-sm"></i>
            <span>Active Reservation: Room {existingBooking.room_number}</span>
          </p>
          <p class="text-[11px] text-slate-400">
            You already have an active reservation (Card: {existingBooking.card_no}). Multiple concurrent bookings are disabled.
          </p>
          <div class="pt-1">
            <a href="/guest" class="inline-flex items-center space-x-1 px-3 py-1 rounded-lg bg-amber-500/20 text-amber-300 border border-amber-500/40 text-[11px] font-bold hover:bg-amber-500/30 transition">
              <span>Go to Guest Portal</span>
              <i class="ph ph-arrow-right text-xs"></i>
            </a>
          </div>
        </div>
      {/if}
    {/if}
  </div>

  <!-- Room Cards Grid -->
  {#if loading}
    <div class="flex justify-center py-12">
      <div class="w-8 h-8 border-2 border-amber-500 border-t-transparent rounded-full animate-spin"></div>
    </div>
  {:else if error}
    <div role="alert" class="mx-auto max-w-md rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-center">
      <p class="text-sm font-semibold text-rose-300">Rooms could not be loaded</p>
      <p class="mt-1 text-xs text-rose-200/80">{error}</p>
      <button
        onclick={fetchRooms}
        class="mt-3 rounded-xl border border-rose-400/30 px-4 py-2 text-xs font-semibold text-rose-200 transition hover:bg-rose-500/10"
      >
        Try again
      </button>
    </div>
  {:else if rooms.length === 0}
    <div class="text-center py-12">
      <i class="ph ph-bed text-4xl text-slate-600"></i>
      <p class="text-slate-400 mt-2 text-sm">No rooms available at the moment</p>
    </div>
  {:else}
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      {#each rooms as room}
        {@const style = getRoomStyle(room.room_type)}
        <div class="bg-gradient-to-br {style.gradient} border {style.border} rounded-3xl overflow-hidden shadow-xl hover:shadow-2xl transition-shadow group">
          <!-- Room Image -->
          <div class="h-44 bg-slate-800 relative overflow-hidden">
            <img
              src={getRoomImage(room)}
              alt="Room {room.room_number} - {room.room_type}"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              loading="lazy"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-slate-950 via-slate-950/30 to-transparent z-10"></div>
            <div class="absolute bottom-3 left-4 z-20">
              <span class="text-[10px] px-2 py-0.5 rounded-full bg-emerald-500/80 text-slate-950 font-bold uppercase tracking-wider shadow">Available</span>
            </div>
            <div class="absolute bottom-3 right-4 z-20">
              <span class="text-xs font-bold text-white font-['Oxanium',sans-serif] bg-slate-900/80 px-2 py-0.5 rounded-lg border border-slate-700">Room {room.room_number}</span>
            </div>
          </div>

          <!-- Room Details -->
          <div class="p-4 space-y-3">
            <div>
              <h3 class="text-base font-bold text-white group-hover:text-amber-300 transition-colors">{room.room_type}</h3>
              <p class="text-[11px] text-slate-400 mt-1 leading-relaxed line-clamp-2">{room.description}</p>
            </div>

            <div class="flex items-center justify-between text-xs">
              <div class="flex items-center space-x-3 text-slate-400">
                <span class="flex items-center space-x-1">
                  <i class="ph ph-bed text-sm"></i>
                  <span>{room.bed_type}</span>
                </span>
              </div>
              <div class="text-right">
                <span class="text-lg font-bold text-amber-400 font-['Oxanium',sans-serif]">{fmtMwk(room.room_rate)}</span>
                <span class="text-[10px] text-slate-500 block">per night</span>
              </div>
            </div>

            {#if existingBooking}
              <button
                disabled
                class="w-full py-3 rounded-2xl bg-slate-800 text-slate-500 font-bold text-xs uppercase tracking-wider cursor-not-allowed border border-slate-700 flex items-center justify-center space-x-1.5 opacity-60"
                title="You already have an active reservation (Room {existingBooking.room_number})"
              >
                <i class="ph ph-lock text-sm"></i>
                <span>Room {existingBooking.room_number} Already Reserved</span>
              </button>
            {:else}
              <button
                onclick={() => openBooking(room)}
                class="w-full py-3 rounded-2xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/20 active:scale-95 transition flex items-center justify-center space-x-1.5"
              >
                <i class="ph ph-calendar-check text-sm"></i>
                <span>Book This Room</span>
              </button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<!-- Booking Modal -->
{#if showBookingModal && selectedRoom}
  <div class="fixed inset-0 z-50 flex items-start justify-center bg-slate-950/85 backdrop-blur-sm p-4 overflow-y-auto" role="dialog" aria-modal="true">
    <div class="w-full max-w-lg bg-slate-900 border border-slate-700 rounded-3xl p-5 sm:p-6 shadow-2xl my-8 space-y-5">
      <!-- Modal Header -->
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-lg font-bold text-white">Book Room {selectedRoom.room_number}</h2>
          <p class="text-xs text-slate-400">{selectedRoom.room_type} · {selectedRoom.bed_type} Bed · {fmtMwk(selectedRoom.room_rate)}/night</p>
        </div>
        <button onclick={() => { showBookingModal = false; bookingResult = null; }} class="p-2 text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition" aria-label="Close">
          <i class="ph ph-x text-lg"></i>
        </button>
      </div>

      {#if bookingResult?.success}
        <!-- Success State -->
        <div class="text-center space-y-4 py-4">
          <div class="w-16 h-16 rounded-full bg-emerald-500/20 text-emerald-400 mx-auto flex items-center justify-center">
            <i class="ph ph-check-circle text-4xl"></i>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">Booking Confirmed!</h3>
            <p class="text-xs text-slate-400 mt-1">{bookingResult.message}</p>
          </div>
          <div class="flex space-x-3">
            <button
              onclick={() => { showBookingModal = false; bookingResult = null; }}
              class="flex-1 py-3 rounded-xl bg-slate-800 hover:bg-slate-700 text-white text-xs font-bold transition"
            >Browse More Rooms</button>
            <button
              onclick={() => goto('/guest')}
              class="flex-1 py-3 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 text-xs font-bold transition"
            >Go to Guest Portal</button>
          </div>
        </div>
      {:else}
        <!-- Booking Form -->
        <form onsubmit={submitBooking} class="space-y-4">
          <!-- Date Selection -->
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label for="arr-date" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1">Check-in Date</label>
              <input id="arr-date" type="date" bind:value={arrivalDate} min={new Date().toISOString().split('T')[0]}
                class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500" />
            </div>
            <div>
              <label for="dep-date" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1">Check-out Date</label>
              <input id="dep-date" type="date" bind:value={departureDate} min={arrivalDate}
                class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500" />
            </div>
          </div>

          <!-- Cost Summary -->
          <div class="bg-slate-950/80 rounded-2xl p-3 border border-slate-800 flex items-center justify-between text-xs">
            <span class="text-slate-400">{nights} night{nights !== 1 ? 's' : ''} × {fmtMwk(selectedRoom.room_rate)}</span>
            <span class="text-lg font-bold text-amber-400 font-['Oxanium',sans-serif]">{fmtMwk(totalCost)}</span>
          </div>

          <!-- Guest Details -->
          <div class="space-y-3">
            {#if profileLoaded && firstName && lastName && !showEditDetails}
              <!-- Compact verified identity card -->
              <div class="bg-emerald-500/10 border border-emerald-500/30 rounded-2xl p-3.5 space-y-1.5">
                <div class="flex items-center justify-between">
                  <p class="text-[10px] text-emerald-400 font-bold uppercase tracking-wider flex items-center space-x-1">
                    <i class="ph ph-check-circle text-xs"></i>
                    <span>Guest Profile Verified</span>
                  </p>
                  <button type="button" onclick={() => { showEditDetails = true; }}
                    class="text-[10px] text-slate-400 hover:text-white underline">Edit Details</button>
                </div>
                <p class="text-sm font-bold text-white">{firstName} {lastName}</p>
                <div class="flex flex-wrap gap-x-4 gap-y-0.5 text-[11px] text-slate-400">
                  {#if passportNo}<span><i class="ph ph-identification-card text-xs text-amber-400"></i> {passportNo}</span>{/if}
                  {#if phone}<span><i class="ph ph-phone text-xs text-amber-400"></i> {phone}</span>{/if}
                  {#if email}<span><i class="ph ph-envelope text-xs text-amber-400"></i> {email}</span>{/if}
                  {#if nationality}<span>{nationality}</span>{/if}
                  {#if city}<span>{city}</span>{/if}
                </div>
              </div>
            {:else}
              <div class="flex items-center justify-between">
                <p class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">Guest Details</p>
                {#if profileLoaded && showEditDetails}
                  <button type="button" onclick={() => { showEditDetails = false; }}
                    class="text-[10px] text-amber-400 hover:text-amber-300 underline">Use Saved Profile</button>
                {/if}
              </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="b-fname" class="sr-only">First Name</label>
                <input id="b-fname" type="text" bind:value={firstName} placeholder="First Name *" required
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
              <div>
                <label for="b-lname" class="sr-only">Last Name</label>
                <input id="b-lname" type="text" bind:value={lastName} placeholder="Last Name *" required
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="b-phone" class="sr-only">Phone</label>
                <input id="b-phone" type="tel" bind:value={phone} placeholder="Phone Number *" required
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
              <div>
                <label for="b-email" class="sr-only">Email</label>
                <input id="b-email" type="email" bind:value={email} placeholder="Email Address"
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="b-nat" class="sr-only">Nationality</label>
                <input id="b-nat" type="text" bind:value={nationality} placeholder="Nationality"
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
              <div>
                <label for="b-city" class="sr-only">City</label>
                <input id="b-city" type="text" bind:value={city} placeholder="City"
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="b-pass" class="sr-only">National ID / Passport</label>
                <input id="b-pass" type="text" bind:value={passportNo} placeholder="National ID / Passport *" required
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
              <div>
                <label for="b-post" class="sr-only">Postal Code</label>
                <input id="b-post" type="text" bind:value={postalCode} placeholder="Postal Code"
                  class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
              </div>
            </div>
            <div>
              <label for="b-addr" class="sr-only">Address</label>
              <input id="b-addr" type="text" bind:value={address} placeholder="Physical Address"
                class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600" />
            </div>
            {/if}
            <div>
              <label for="b-pay" class="block text-slate-400 text-[10px] font-bold uppercase tracking-wider mb-1">Payment Method</label>
              <select id="b-pay" bind:value={paymentMode}
                class="w-full px-3 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white text-xs focus:outline-none focus:border-amber-500">
                <option>Cash</option>
                <option>Credit Card</option>
                <option>Cheque</option>
              </select>
            </div>
          </div>

          {#if bookingResult && !bookingResult.success}
            <div class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs">
              <i class="ph ph-warning text-sm"></i> {bookingResult.message}
            </div>
          {/if}

          <button
            type="submit"
            disabled={isBooking || !firstName || !lastName || !phone}
            class="w-full py-3.5 rounded-2xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-95 transition disabled:opacity-50 flex items-center justify-center space-x-1.5"
          >
            {#if isBooking}
              <div class="w-4 h-4 border-2 border-slate-950 border-t-transparent rounded-full animate-spin"></div>
              <span>Processing...</span>
            {:else}
              <i class="ph ph-calendar-check text-sm"></i>
              <span>Confirm Booking · {fmtMwk(totalCost)}</span>
            {/if}
          </button>
        </form>
      {/if}
    </div>
  </div>
{/if}
