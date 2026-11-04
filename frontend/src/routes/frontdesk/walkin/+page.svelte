<!-- frontend/src/routes/frontdesk/walkin/+page.svelte -->
<script lang="ts">
  import { page } from '$app/state';
  import { roomStore, isOnline } from '#lib/stores/roomStore';
  import type { Room, Folio } from '#lib/types/room';
  import { onMount } from 'svelte';

  const rooms = roomStore.rooms;

  // Form State
  let guestName = $state('');
  let guestPhone = $state('');
  let guestIdDoc = $state('');
  let selectedRoomNumber = $state('');
  let nights = $state(1);
  let depositAmount = $state(0);
  let paymentMethod = $state('Cash');
  let depositMode = $state<'full' | 'firstNight' | 'custom'>('full');

  let isSubmitting = $state(false);
  let errorMessage = $state('');
  let completedFolio = $state<Folio | null>(null);

  // Available rooms
  let availableRooms = $derived($rooms.filter((r: Room) => r.status === 'Available'));

  let selectedRoom = $derived(
    $rooms.find((r: Room) => r.room_number === selectedRoomNumber) || null
  );

  let totalCharges = $derived(
    selectedRoom ? selectedRoom.price_per_night * nights : 0
  );

  // Auto calculate deposit depending on depositMode
  $effect(() => {
    if (!selectedRoom) return;
    if (depositMode === 'full') {
      depositAmount = totalCharges;
    } else if (depositMode === 'firstNight') {
      depositAmount = selectedRoom.price_per_night;
    }
  });

  onMount(() => {
    const queryRoom = page.url.searchParams.get('room');
    if (queryRoom) {
      selectedRoomNumber = queryRoom;
    } else if (availableRooms.length > 0) {
      selectedRoomNumber = availableRooms[0].room_number;
    }
  });

  async function handleWalkInSubmit(e: Event) {
    e.preventDefault();
    if (!guestName.trim()) {
      errorMessage = 'Please enter guest full name.';
      return;
    }
    if (!guestIdDoc.trim()) {
      errorMessage = 'Guest ID or Passport number is required by policy.';
      return;
    }
    if (!selectedRoomNumber) {
      errorMessage = 'Please select an available room.';
      return;
    }

    errorMessage = '';
    isSubmitting = true;

    try {
      const res = await roomStore.checkIn({
        guest_name: guestName.trim(),
        guest_phone: guestPhone.trim(),
        guest_id_doc: guestIdDoc.trim(),
        room_number: selectedRoomNumber,
        nights: Number(nights),
        deposit_paid: Number(depositAmount)
      });

      if (res.success && res.folio) {
        completedFolio = res.folio;
      } else {
        errorMessage = res.message || 'Walk-in check-in failed';
      }
    } catch (err: any) {
      errorMessage = err.message || 'An unexpected error occurred.';
    } finally {
      isSubmitting = false;
    }
  }

  function handleReset() {
    guestName = '';
    guestPhone = '';
    guestIdDoc = '';
    nights = 1;
    completedFolio = null;
    depositMode = 'full';
    if (availableRooms.length > 0) {
      selectedRoomNumber = availableRooms[0].room_number;
    }
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">RAPID WALK-IN</h1>
      <p class="text-xs text-slate-400">Immediate Front Desk Registration</p>
    </div>
    {#if !$isOnline}
      <span class="text-[10px] font-bold bg-amber-500/20 text-amber-300 border border-amber-500/30 px-2 py-0.5 rounded-full flex items-center">
        <i class="ph ph-wifi-slash mr-1"></i> Offline Mode Active
      </span>
    {/if}
  </div>

  {#if errorMessage}
    <div class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs flex items-center justify-between">
      <span>{errorMessage}</span>
      <button type="button" aria-label="Close error" onclick={() => (errorMessage = '')}><i class="ph ph-x text-sm"></i></button>
    </div>
  {/if}

  {#if completedFolio}
    <!-- Success Confirmation Receipt View -->
    <div class="bg-slate-900 border border-emerald-500/40 rounded-2xl p-5 shadow-2xl space-y-4 animate-scale-up">
      <div class="text-center space-y-1">
        <div class="w-12 h-12 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center mx-auto text-2xl">
          <i class="ph ph-check-fat"></i>
        </div>
        <h2 class="text-base font-bold text-white">Registration Complete!</h2>
        <p class="text-xs text-slate-400">Room {completedFolio.room_number} is locked & assigned</p>
      </div>

      <div class="bg-slate-950 border border-slate-800 rounded-xl p-3.5 space-y-2 text-xs">
        <div class="flex justify-between">
          <span class="text-slate-400">Guest:</span>
          <span class="font-bold text-white">{completedFolio.guest_name}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Room Assigned:</span>
          <span class="font-bold text-amber-400 font-['Oxanium',sans-serif]">Room {completedFolio.room_number}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Stay Period:</span>
          <span class="font-mono text-slate-300">{completedFolio.check_in_date} &rarr; {completedFolio.check_out_date}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Total Charges:</span>
          <span class="font-mono text-white">MWK {completedFolio.total_charges.toLocaleString()}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Deposit Paid ({paymentMethod}):</span>
          <span class="font-mono text-emerald-400 font-semibold">MWK {completedFolio.total_paid.toLocaleString()}</span>
        </div>
        <div class="flex justify-between pt-1 border-t border-slate-800 font-bold">
          <span class="text-slate-300">Folio Balance:</span>
          <span class="font-mono {completedFolio.balance > 0 ? 'text-amber-400' : 'text-emerald-400'}">
            MWK {completedFolio.balance.toLocaleString()}
          </span>
        </div>
      </div>

      <!-- Resident Portal Access Box for Walk-in Guests -->
      <div class="p-3.5 bg-slate-950 border border-amber-500/30 rounded-xl space-y-1.5 text-left">
        <div class="flex items-center space-x-1.5 text-amber-300 font-bold text-xs">
          <i class="ph ph-device-mobile text-amber-400 text-base"></i>
          <span>Guest Mobile Portal Access</span>
        </div>
        <p class="text-[11px] text-slate-300">
          Provide these credentials so the guest can sign in to the Guest Portal on their phone (In-Room Dining, Folio, Services):
        </p>
        <div class="flex items-center space-x-4 bg-slate-900 px-3 py-2 rounded-lg border border-slate-800 text-xs">
          <div>
            <span class="text-[10px] text-slate-500 uppercase block font-semibold">Room</span>
            <span class="font-bold text-white font-mono">{completedFolio.room_number}</span>
          </div>
          <div>
            <span class="text-[10px] text-slate-500 uppercase block font-semibold">Card Reference</span>
            <span class="font-bold text-amber-400 font-mono">{completedFolio.card_no}</span>
          </div>
          <div class="flex-1 text-right">
            <span class="text-[10px] text-slate-400 block">Tab: <strong class="text-white">Resident Sign-In</strong></span>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-2 pt-1">
        <a
          href="/frontdesk/folio?room={completedFolio.room_number}"
          class="py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs text-center shadow-md active:scale-95 transition"
        >
          View Full Folio &rarr;
        </a>
        <button
          onclick={handleReset}
          class="py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-xs active:scale-95 transition"
        >
          New Registration
        </button>
      </div>
    </div>
  {:else}
    <!-- Form -->
    <form onsubmit={handleWalkInSubmit} class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-4 shadow-xl">
      <!-- Section 1: Guest Information -->
      <div class="space-y-2.5">
        <div class="flex items-center space-x-1.5 text-xs font-bold uppercase tracking-wider text-amber-400">
          <i class="ph ph-user text-sm"></i>
          <span>1. Guest Identification</span>
        </div>

        <div>
          <label for="walkin-guest-name" class="block text-[11px] font-semibold text-slate-400 mb-1">Full Legal Name *</label>
          <input
            id="walkin-guest-name"
            type="text"
            placeholder="e.g. Johnathan Phiri"
            bind:value={guestName}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2.5 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500"
          />
        </div>

        <div class="grid grid-cols-2 gap-2">
          <div>
            <label for="walkin-guest-phone" class="block text-[11px] font-semibold text-slate-400 mb-1">Phone Number</label>
            <input
              id="walkin-guest-phone"
              type="tel"
              placeholder="+265 999 000 000"
              bind:value={guestPhone}
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500"
            />
          </div>
          <div>
            <label for="walkin-guest-id" class="block text-[11px] font-semibold text-slate-400 mb-1">National ID / Passport *</label>
            <input
              id="walkin-guest-id"
              type="text"
              placeholder="e.g. MW-ID-7819"
              bind:value={guestIdDoc}
              required
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500 font-mono"
            />
          </div>
        </div>
      </div>

      <!-- Section 2: Room Selection -->
      <div class="space-y-2.5 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between text-xs font-bold uppercase tracking-wider text-amber-400">
          <span class="flex items-center space-x-1.5">
            <i class="ph ph-bed text-sm"></i>
            <span>2. Room Selection</span>
          </span>
          <span class="text-[10px] lowercase text-slate-400 font-normal">
            {availableRooms.length} available
          </span>
        </div>

        {#if availableRooms.length === 0}
          <div class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-400 text-xs text-center">
            No rooms available for immediate check-in.
          </div>
        {:else}
          <div class="grid grid-cols-3 gap-2">
            {#each availableRooms as r}
              <button
                type="button"
                onclick={() => (selectedRoomNumber = r.room_number)}
                class="p-2.5 rounded-xl border text-center transition flex flex-col items-center {selectedRoomNumber === r.room_number
                  ? 'bg-amber-500 text-slate-950 font-bold border-amber-500 shadow-md shadow-amber-500/20'
                  : 'bg-slate-950 border-slate-800 text-slate-300 hover:border-slate-700'}"
              >
                <span class="font-['Oxanium',sans-serif] text-sm font-bold">{r.room_number}</span>
                <span class="text-[10px] opacity-80">{r.type}</span>
                <span class="text-[9px] font-mono mt-0.5 opacity-90">MWK {(r.price_per_night / 1000).toFixed(0)}k</span>
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Section 3: Duration of Stay -->
      <div class="space-y-2.5 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between text-xs font-bold uppercase tracking-wider text-amber-400">
          <span class="flex items-center space-x-1.5">
            <i class="ph ph-calendar text-sm"></i>
            <span>3. Duration of Stay</span>
          </span>
          <span class="text-[11px] font-mono text-white">{nights} Night(s)</span>
        </div>

        <!-- Quick Nights Pills -->
        <div class="flex items-center space-x-2">
          {#each [1, 2, 3, 5, 7] as n}
            <button
              type="button"
              onclick={() => (nights = n)}
              class="flex-1 py-1.5 rounded-xl border text-xs font-bold transition {nights === n
                ? 'bg-amber-500 text-slate-950 border-amber-500'
                : 'bg-slate-950 border-slate-800 text-slate-300 hover:text-white'}"
            >
              {n}d
            </button>
          {/each}
        </div>
      </div>

      <!-- Section 4: Payment & Deposit -->
      <div class="space-y-2.5 pt-2 border-t border-slate-800/80">
        <div class="flex items-center justify-between text-xs font-bold uppercase tracking-wider text-amber-400">
          <span class="flex items-center space-x-1.5">
            <i class="ph ph-wallet text-sm"></i>
            <span>4. Initial Payment</span>
          </span>
          <span class="text-[11px] font-mono text-white font-bold">Total: MWK {totalCharges.toLocaleString()}</span>
        </div>

        <!-- Deposit Mode Selector -->
        <div class="grid grid-cols-3 gap-2 text-[11px]">
          <button
            type="button"
            onclick={() => (depositMode = 'full')}
            class="py-1.5 rounded-xl border font-semibold transition {depositMode === 'full'
              ? 'bg-amber-500 text-slate-950 font-bold border-amber-500'
              : 'bg-slate-950 border-slate-800 text-slate-400'}"
          >
            Full (100%)
          </button>
          <button
            type="button"
            onclick={() => (depositMode = 'firstNight')}
            class="py-1.5 rounded-xl border font-semibold transition {depositMode === 'firstNight'
              ? 'bg-amber-500 text-slate-950 font-bold border-amber-500'
              : 'bg-slate-950 border-slate-800 text-slate-400'}"
          >
            1st Night
          </button>
          <button
            type="button"
            onclick={() => (depositMode = 'custom')}
            class="py-1.5 rounded-xl border font-semibold transition {depositMode === 'custom'
              ? 'bg-amber-500 text-slate-950 font-bold border-amber-500'
              : 'bg-slate-950 border-slate-800 text-slate-400'}"
          >
            Custom
          </button>
        </div>

        {#if depositMode === 'custom'}
          <div>
            <label for="walkin-deposit-amount" class="block text-[11px] font-semibold text-slate-400 mb-1">Custom Deposit (MWK)</label>
            <input
              id="walkin-deposit-amount"
              type="number"
              min="0"
              step="1000"
              bind:value={depositAmount}
              class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white font-mono focus:border-amber-500"
            />
          </div>
        {/if}

        <!-- Payment Method Options -->
        <div>
          <span id="walkin-payment-method-label" class="block text-[11px] font-semibold text-slate-400 mb-1">Payment Method</span>
          <div role="group" aria-labelledby="walkin-payment-method-label" class="grid grid-cols-4 gap-1.5">
            {#each ['Cash', 'Visa', 'Airtel', 'Mpamba'] as method}
              <button
                type="button"
                onclick={() => (paymentMethod = method)}
                class="py-1.5 px-1 rounded-xl border text-[11px] font-semibold text-center transition {paymentMethod === method
                  ? 'bg-emerald-500 text-slate-950 font-bold border-emerald-500'
                  : 'bg-slate-950 border-slate-800 text-slate-400'}"
              >
                {method}
              </button>
            {/each}
          </div>
        </div>
      </div>

      <!-- Action Button -->
      <button
        type="submit"
        disabled={isSubmitting || availableRooms.length === 0}
        class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-98 transition disabled:opacity-50"
      >
        {isSubmitting ? 'Registering Guest...' : 'Complete Walk-In Check-In'}
      </button>
    </form>
  {/if}
</div>
