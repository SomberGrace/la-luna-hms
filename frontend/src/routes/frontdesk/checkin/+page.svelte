<!-- frontend/src/routes/frontdesk/checkin/+page.svelte -->
<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { getPendingArrivals } from '#lib/frontdeskReservations';
  import { page } from '$app/state';
  import { roomStore } from '#lib/stores/roomStore';
  import type { ReservationItem } from '#lib/types/guest';
  import type { Room } from '#lib/types/room';
  import { onMount } from 'svelte';

  const rooms = roomStore.rooms;

  // Form inputs
  let resCode = $state('');
  let guestName = $state('');
  let guestPhone = $state('');
  let guestIdDoc = $state('');
  let selectedRoomNumber = $state('');
  let nights = $state(1);
  let depositPaid = $state(0);
  let paymentMethod = $state('Cash');

  // Pending arrivals list (loaded live from database)
  let pendingReservations = $state<ReservationItem[]>([]);
  let isLoadingReservations = $state(false);

  let isSubmitting = $state(false);
  let feedbackMessage = $state<{ text: string; isError: boolean } | null>(null);
  let checkedInFolioRoom = $state<string | null>(null);

  // Available rooms for assignment
  let availableRooms = $derived($rooms.filter((r) => r.status === 'Available'));
  let selectableRooms = $derived(
    $rooms.filter((r) => r.status === 'Available' || r.room_number === selectedRoomNumber)
  );

  let selectedRoom = $derived(
    $rooms.find((r) => r.room_number === selectedRoomNumber) || null
  );

  let totalEstimatedCharge = $derived(
    selectedRoom ? selectedRoom.price_per_night * nights : 0
  );

  let remainingDue = $derived(Math.max(0, totalEstimatedCharge - depositPaid));

  onMount(async () => {
    await roomStore.fetchRooms();
    await loadPendingArrivals();

    // If room specified in URL query
    const urlRoom = page.url.searchParams.get('room');
    if (urlRoom) {
      selectedRoomNumber = urlRoom;
    } else if (availableRooms.length > 0) {
      selectedRoomNumber = availableRooms[0].room_number;
    }
  });

  async function loadPendingArrivals() {
    isLoadingReservations = true;
    try {
      const res = await apiFetch('/api/v1/frontdesk/reservations');
      if (res.ok) {
        const data: unknown = await res.json();
        const arrivals = getPendingArrivals(data);
        if (arrivals) pendingReservations = arrivals;
      }
    } catch (e) {
      console.warn('Failed to load live reservations:', e);
    } finally {
      isLoadingReservations = false;
    }
  }

  function selectReservation(res: ReservationItem) {
    resCode = res.res_code;
    guestName = res.guest_name;
    guestPhone = res.guest_phone;
    guestIdDoc = res.guest_id_doc || 'Online Verified ID';
    nights = res.nights;
    depositPaid = res.deposit_paid;
    selectedRoomNumber = res.room_number;
  }

  async function quickClearCheckIn(r: ReservationItem) {
    isSubmitting = true;
    feedbackMessage = null;
    try {
      const res = await roomStore.checkIn({
        guest_name: r.guest_name,
        guest_phone: r.guest_phone || 'N/A',
        guest_id_doc: r.guest_id_doc || 'Online Verified ID',
        room_number: r.room_number,
        nights: Number(r.nights || 1),
        deposit_paid: Number(r.deposit_paid || 0),
        reservation_code: r.res_code
      });

      if (res.success) {
        feedbackMessage = { text: `Guest ${r.guest_name} cleared and checked in to Room ${r.room_number}!`, isError: false };
        checkedInFolioRoom = r.room_number;
        pendingReservations = pendingReservations.filter((p) => p.res_code !== r.res_code);
        await roomStore.fetchRooms();
        await roomStore.fetchFolio(r.room_number);
      } else {
        feedbackMessage = { text: res.message || 'Check-in failed.', isError: true };
      }
    } catch (err: any) {
      feedbackMessage = { text: err.message || 'An error occurred during check-in.', isError: true };
    } finally {
      isSubmitting = false;
    }
  }

  async function handleSubmitCheckIn(e: Event) {
    e.preventDefault();
    if (!guestName.trim() || !selectedRoomNumber) {
      feedbackMessage = { text: 'Please fill in guest name and choose an available room.', isError: true };
      return;
    }

    if (!guestIdDoc.trim()) {
      feedbackMessage = { text: 'Guest National ID or Passport document is required for verification.', isError: true };
      return;
    }

    isSubmitting = true;
    feedbackMessage = null;

    try {
      const res = await roomStore.checkIn({
        guest_name: guestName.trim(),
        guest_phone: guestPhone.trim(),
        guest_id_doc: guestIdDoc.trim(),
        room_number: selectedRoomNumber,
        nights: Number(nights),
        deposit_paid: Number(depositPaid),
        reservation_code: resCode.trim() || undefined
      });

      if (res.success) {
        feedbackMessage = { text: `Guest successfully checked in to Room ${selectedRoomNumber}!`, isError: false };
        checkedInFolioRoom = selectedRoomNumber;

        // Remove from pending arrivals list if matched
        if (resCode) {
          pendingReservations = pendingReservations.filter((r) => r.res_code !== resCode);
        }
      } else {
        feedbackMessage = { text: res.message || 'Check-in failed.', isError: true };
      }
    } catch (err: any) {
      feedbackMessage = { text: err.message || 'An error occurred during check-in.', isError: true };
    } finally {
      isSubmitting = false;
    }
  }

  function resetForm() {
    resCode = '';
    guestName = '';
    guestPhone = '';
    guestIdDoc = '';
    depositPaid = 0;
    nights = 1;
    checkedInFolioRoom = null;
    feedbackMessage = null;
    if (availableRooms.length > 0) {
      selectedRoomNumber = availableRooms[0].room_number;
    }
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">RESERVATION CHECK-IN</h1>
      <p class="text-xs text-slate-400">Process Expected Arrivals & ID Verification</p>
    </div>
    <span class="text-xs font-mono bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 px-2.5 py-1 rounded-full font-bold">
      {availableRooms.length} Available
    </span>
  </div>

  <!-- Pending Arrivals Quick Select Carousel/Cards -->
  {#if !checkedInFolioRoom}
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5 space-y-2">
      <div class="flex items-center justify-between">
        <span class="text-xs font-bold text-amber-400 flex items-center">
          <i class="ph ph-clock-countdown mr-1 text-sm"></i> Pending Arrivals Today ({pendingReservations.length})
        </span>
        <button
          type="button"
          onclick={loadPendingArrivals}
          disabled={isLoadingReservations}
          class="text-[11px] text-slate-400 hover:text-amber-400 flex items-center space-x-1 px-2.5 py-1 rounded-lg bg-slate-800 border border-slate-700 hover:border-amber-500/50 transition-colors disabled:opacity-50"
        >
          <i class="ph ph-arrows-clockwise {isLoadingReservations ? 'animate-spin' : ''}"></i>
          <span>{isLoadingReservations ? 'Checking...' : 'Refresh'}</span>
        </button>
      </div>

      {#if isLoadingReservations && pendingReservations.length === 0}
        <p class="text-xs text-slate-500 py-2 text-center">Loading pending arrivals...</p>
      {:else if pendingReservations.length === 0}
        <div class="py-2.5 px-3 bg-slate-950/60 rounded-xl border border-slate-800/80 text-center text-xs text-slate-400">
          No pending arrivals found. Online guest bookings will appear here instantly.
        </div>
      {:else}
        <div class="space-y-2">
          {#each pendingReservations as r}
            <div
              class="w-full bg-slate-950/80 border border-slate-800 rounded-xl p-3 flex items-center justify-between transition hover:border-amber-500/40"
            >
              <div>
                <div class="flex items-center space-x-2">
                  <span class="text-xs font-bold text-white">{r.guest_name}</span>
                  <span class="text-[10px] font-mono bg-slate-800 text-slate-400 px-1.5 py-0.5 rounded">{r.res_code}</span>
                  {#if r.room_number}
                    <span class="text-[10px] font-mono bg-amber-500/10 text-amber-300 border border-amber-500/20 px-1.5 py-0.5 rounded">Room {r.room_number}</span>
                  {/if}
                </div>
                <p class="text-[11px] text-slate-400 mt-1">
                  {r.room_type} Room &bull; {r.nights} Night(s) &bull; Dep: MWK {r.deposit_paid.toLocaleString()}
                </p>
              </div>
              <div class="flex items-center space-x-2">
                <button
                  type="button"
                  onclick={() => selectReservation(r)}
                  class="text-[11px] font-semibold text-slate-300 hover:text-white bg-slate-800 hover:bg-slate-700 px-2.5 py-1.5 rounded-lg border border-slate-700 transition"
                >
                  Edit / Review
                </button>
                <button
                  type="button"
                  onclick={() => quickClearCheckIn(r)}
                  disabled={isSubmitting}
                  class="text-[11px] font-bold text-slate-950 bg-amber-500 hover:bg-amber-400 px-3 py-1.5 rounded-lg shadow-md shadow-amber-500/20 transition active:scale-95 disabled:opacity-50 flex items-center space-x-1"
                >
                  <i class="ph ph-check-bold text-xs"></i>
                  <span>Clear & Check-In</span>
                </button>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {/if}

  {#if feedbackMessage}
    <div
      class="p-3.5 rounded-2xl border text-xs flex items-start space-x-2 {feedbackMessage.isError
        ? 'bg-rose-500/10 border-rose-500/30 text-rose-300'
        : 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300'}"
    >
      <i class="ph {feedbackMessage.isError ? 'ph-warning-circle text-rose-400' : 'ph-check-circle text-emerald-400'} text-base mt-0.5"></i>
      <div class="flex-1">
        <p class="font-semibold">{feedbackMessage.text}</p>
        {#if checkedInFolioRoom}
          <div class="mt-2.5 p-3 rounded-xl bg-slate-900 border border-amber-500/30 space-y-1 text-left">
            <span class="text-xs font-bold text-amber-300 flex items-center">
              <i class="ph ph-device-mobile mr-1.5 text-amber-400"></i> Resident Guest Portal Credentials:
            </span>
            <p class="text-xs text-slate-300">
              Room: <span class="font-bold text-white font-mono">{checkedInFolioRoom}</span> &bull; 
              Card Ref: <span class="font-bold text-amber-400 font-mono">{resCode || 'CARD-' + checkedInFolioRoom}</span>
            </p>
            <p class="text-[10px] text-slate-400">
              The guest can sign into the Guest Portal on their mobile device using their Room Number & Card Reference.
            </p>
          </div>

          <div class="flex items-center space-x-2 mt-3">
            <a
              href="/frontdesk/folio?room={checkedInFolioRoom}"
              class="px-3 py-1.5 rounded-xl bg-amber-500 text-slate-950 font-bold text-xs shadow-md"
            >
              Open Guest Folio &rarr;
            </a>
            <button
              onclick={resetForm}
              class="px-3 py-1.5 rounded-xl bg-slate-800 text-slate-300 font-semibold text-xs"
            >
              Check-In Another
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/if}

  <!-- Check-In Form -->
  {#if !checkedInFolioRoom}
    <form onsubmit={handleSubmitCheckIn} class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-3.5 shadow-xl">
      <!-- Reservation Code (Optional) -->
      <div>
        <label for="res-code" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
          Reservation Code (Optional)
        </label>
        <input
          id="res-code"
          type="text"
          placeholder="e.g. RES-2026-001"
          bind:value={resCode}
          class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500 font-mono"
        />
      </div>

      <!-- Guest Full Name -->
      <div>
        <label for="guest-name" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
          Guest Full Name *
        </label>
        <input
          id="guest-name"
          type="text"
          placeholder="e.g. David Mwale"
          bind:value={guestName}
          required
          class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500"
        />
      </div>

      <!-- Phone & National ID / Passport -->
      <div class="grid grid-cols-2 gap-2.5">
        <div>
          <label for="guest-phone" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
            Mobile Number
          </label>
          <input
            id="guest-phone"
            type="tel"
            placeholder="+265 ..."
            bind:value={guestPhone}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500"
          />
        </div>

        <div>
          <label for="guest-id" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
            National ID / Passport *
          </label>
          <input
            id="guest-id"
            type="text"
            placeholder="e.g. MW-ID-..."
            bind:value={guestIdDoc}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500 font-mono"
          />
        </div>
      </div>

      <!-- Room Selection -->
      <div>
        <label for="room-select" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
          Assign Room *
        </label>
        {#if selectableRooms.length === 0}
          <div class="p-2.5 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-400 text-xs">
            No rooms currently available or assigned to this reservation. Check the Housekeeping tab to expedite cleaning.
          </div>
        {:else}
          <select
            id="room-select"
            bind:value={selectedRoomNumber}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2.5 text-xs text-white focus:outline-none focus:border-amber-500 font-medium"
          >
            {#each selectableRooms as r}
              <option value={r.room_number}>
                Room {r.room_number} &bull; {r.type} &bull; MWK {r.price_per_night.toLocaleString()}/night ({r.status})
              </option>
            {/each}
          </select>
        {/if}
      </div>

      <!-- Stay Duration & Rate Breakdown -->
      <div class="grid grid-cols-2 gap-2.5">
        <div>
          <label for="nights-count" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
            Number of Nights
          </label>
          <input
            id="nights-count"
            type="number"
            min="1"
            max="60"
            bind:value={nights}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-amber-500"
          />
        </div>

        <div>
          <label for="deposit-amount" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
            Deposit Paid (MWK)
          </label>
          <input
            id="deposit-amount"
            type="number"
            min="0"
            step="1000"
            bind:value={depositPaid}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white focus:outline-none focus:border-amber-500 font-mono"
          />
        </div>
      </div>

      <!-- Financial Calculation Pill -->
      <div class="bg-slate-950 border border-slate-800/80 rounded-xl p-3 space-y-1.5 text-xs">
        <div class="flex justify-between text-slate-400">
          <span>Estimated Total:</span>
          <span class="text-white font-mono font-semibold">MWK {totalEstimatedCharge.toLocaleString()}</span>
        </div>
        <div class="flex justify-between text-slate-400">
          <span>Deposit Collected:</span>
          <span class="text-emerald-400 font-mono font-semibold">- MWK {(depositPaid || 0).toLocaleString()}</span>
        </div>
        <div class="flex justify-between pt-1 border-t border-slate-800 text-white font-bold">
          <span>Estimated Balance Due:</span>
          <span class="text-amber-400 font-mono">MWK {remainingDue.toLocaleString()}</span>
        </div>
      </div>

      <!-- Action Button -->
      <button
        type="submit"
        disabled={isSubmitting || selectableRooms.length === 0}
        class="w-full py-3 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-98 transition disabled:opacity-50"
      >
        {isSubmitting ? 'Processing Check-In...' : 'Confirm & Check-In Guest'}
      </button>
    </form>
  {/if}
</div>
