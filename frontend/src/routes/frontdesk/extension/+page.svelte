<!-- frontend/src/routes/frontdesk/extension/+page.svelte -->
<script lang="ts">
  import { page } from '$app/state';
  import { roomStore } from '#lib/stores/roomStore';
  import type { Room, Folio } from '#lib/types/room';
  import { onMount } from 'svelte';

  const rooms = roomStore.rooms;
  const folios = roomStore.folios;

  let selectedRoomNumber = $state('');
  let additionalNights = $state(1);
  let payNow = $state(false);
  let paymentMethod = $state('Cash');

  let isSubmitting = $state(false);
  let feedbackMessage = $state<{ text: string; isError: boolean } | null>(null);

  // List of currently occupied rooms
  let occupiedRooms = $derived($rooms.filter((r: Room) => r.status === 'Occupied'));

  let activeRoom = $derived(
    occupiedRooms.find((r: Room) => r.room_number === selectedRoomNumber) || null
  );

  let activeFolio = $derived<Folio | null>(
    selectedRoomNumber ? $folios[selectedRoomNumber] || null : null
  );

  let additionalCharge = $derived(
    activeRoom ? activeRoom.price_per_night * additionalNights : 0
  );

  let projectedCheckoutDate = $derived.by(() => {
    if (!activeFolio?.check_out_date) return '';
    try {
      const co = new Date(activeFolio.check_out_date);
      const newCo = new Date(co.getTime() + additionalNights * 24 * 3600 * 1000);
      return newCo.toISOString().split('T')[0];
    } catch (e) {
      return '';
    }
  });

  onMount(() => {
    const queryRoom = page.url.searchParams.get('room');
    if (queryRoom && occupiedRooms.some((r) => r.room_number === queryRoom)) {
      selectedRoomNumber = queryRoom;
    } else if (occupiedRooms.length > 0) {
      selectedRoomNumber = occupiedRooms[0].room_number;
    }
  });

  async function handleExtensionSubmit(e: Event) {
    e.preventDefault();
    if (!selectedRoomNumber) {
      feedbackMessage = { text: 'Please select an occupied room.', isError: true };
      return;
    }

    isSubmitting = true;
    feedbackMessage = null;

    try {
      const res = await roomStore.extendStay({
        room_number: selectedRoomNumber,
        additional_nights: Number(additionalNights),
        pay_now: payNow,
        payment_method: payNow ? paymentMethod : undefined
      });

      if (res.success) {
        feedbackMessage = {
          text: `Stay extended successfully for Room ${selectedRoomNumber}! New checkout date: ${projectedCheckoutDate}.`,
          isError: false
        };
      } else {
        feedbackMessage = { text: res.message || 'Stay extension failed.', isError: true };
      }
    } catch (err: any) {
      feedbackMessage = { text: err.message || 'Failed to process extension.', isError: true };
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">STAY EXTENSION</h1>
      <p class="text-xs text-slate-400">Extend In-House Guest Stay Duration</p>
    </div>
    <span class="text-xs font-mono bg-amber-500/10 text-amber-400 border border-amber-500/20 px-2.5 py-1 rounded-full font-bold">
      {occupiedRooms.length} Occupied
    </span>
  </div>

  {#if feedbackMessage}
    <div
      class="p-3.5 rounded-2xl border text-xs flex items-start space-x-2 {feedbackMessage.isError
        ? 'bg-rose-500/10 border-rose-500/30 text-rose-300'
        : 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300'}"
    >
      <i class="ph {feedbackMessage.isError ? 'ph-warning-circle text-rose-400' : 'ph-check-circle text-emerald-400'} text-base mt-0.5"></i>
      <div class="flex-1">
        <p class="font-semibold">{feedbackMessage.text}</p>
        {#if !feedbackMessage.isError}
          <div class="mt-2 flex space-x-2">
            <a
              href="/frontdesk/folio?room={selectedRoomNumber}"
              class="px-3 py-1.5 rounded-xl bg-amber-500 text-slate-950 font-bold text-xs shadow-md"
            >
              View Updated Folio &rarr;
            </a>
          </div>
        {/if}
      </div>
    </div>
  {/if}

  {#if occupiedRooms.length === 0}
    <div class="bg-slate-900 border border-slate-800 rounded-2xl p-8 text-center text-slate-500 text-xs">
      <i class="ph ph-bed text-3xl mb-2 text-slate-600 block"></i>
      No guests are currently checked in to extend.
    </div>
  {:else}
    <form onsubmit={handleExtensionSubmit} class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-4 shadow-xl">
      <!-- Select Occupied Room -->
      <div>
        <label for="room-extension-select" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
          Select In-House Room *
        </label>
        <select
          id="room-extension-select"
          bind:value={selectedRoomNumber}
          class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2.5 text-xs text-white focus:outline-none focus:border-amber-500 font-medium"
        >
          {#each occupiedRooms as r}
            {@const f = $folios[r.room_number]}
            <option value={r.room_number}>
              Room {r.room_number} &bull; {f?.guest_name || 'Guest'} &bull; C/O: {f?.check_out_date || 'N/A'}
            </option>
          {/each}
        </select>
      </div>

      <!-- Current Stay Snapshot -->
      {#if activeFolio && activeRoom}
        <div class="bg-slate-950 border border-slate-800/80 rounded-xl p-3.5 space-y-2 text-xs">
          <div class="flex items-center justify-between">
            <span class="text-slate-400">Current Guest:</span>
            <span class="font-bold text-white">{activeFolio.guest_name}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-slate-400">Room Category:</span>
            <span class="text-amber-400 font-semibold">{activeRoom.type} (MWK {activeRoom.price_per_night.toLocaleString()}/nt)</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-slate-400">Current Checkout Date:</span>
            <span class="font-mono text-slate-200 font-bold">{activeFolio.check_out_date}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-slate-400">Current Balance:</span>
            <span class="font-mono {activeFolio.balance > 0 ? 'text-amber-400' : 'text-emerald-400'} font-bold">
              MWK {activeFolio.balance.toLocaleString()}
            </span>
          </div>
        </div>
      {/if}

      <!-- Additional Nights Selection -->
      <div>
        <div class="flex items-center justify-between mb-1.5">
          <span id="extend-duration-label" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400">
            Extend Duration
          </span>
          <span class="text-xs font-mono font-bold text-amber-400">+{additionalNights} Night(s)</span>
        </div>

        <div role="group" aria-labelledby="extend-duration-label" class="flex items-center space-x-2">
          {#each [1, 2, 3, 5, 7] as n}
            <button
              type="button"
              onclick={() => (additionalNights = n)}
              class="flex-1 py-2 rounded-xl border text-xs font-bold transition {additionalNights === n
                ? 'bg-amber-500 text-slate-950 border-amber-500 shadow-md shadow-amber-500/20'
                : 'bg-slate-950 border-slate-800 text-slate-300 hover:text-white'}"
            >
              +{n}d
            </button>
          {/each}
        </div>
      </div>

      <!-- Projection Card -->
      <div class="bg-amber-500/10 border border-amber-500/20 rounded-xl p-3.5 space-y-1.5 text-xs">
        <div class="flex justify-between text-amber-200">
          <span>New Checkout Date:</span>
          <span class="font-mono font-bold text-amber-400">{projectedCheckoutDate} (10:00 AM)</span>
        </div>
        <div class="flex justify-between text-amber-200">
          <span>Additional Accommodation Charge:</span>
          <span class="font-mono font-bold text-white">MWK {additionalCharge.toLocaleString()}</span>
        </div>
      </div>

      <!-- Immediate Payment Toggle -->
      <div class="space-y-2.5 pt-1">
        <label class="flex items-center justify-between p-3 rounded-xl bg-slate-950 border border-slate-800 cursor-pointer">
          <span class="text-xs font-semibold text-slate-200">Collect Extension Payment Now</span>
          <input
            type="checkbox"
            bind:checked={payNow}
            class="w-4 h-4 rounded text-amber-500 focus:ring-amber-500 bg-slate-800 border-slate-700"
          />
        </label>

        {#if payNow}
          <div class="grid grid-cols-4 gap-1.5 animate-fade-in">
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
        {/if}
      </div>

      <!-- Action Button -->
      <button
        type="submit"
        disabled={isSubmitting || !activeFolio}
        class="w-full py-3.5 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-98 transition disabled:opacity-50"
      >
        {isSubmitting ? 'Updating Stay...' : 'Confirm Stay Extension'}
      </button>
    </form>
  {/if}
</div>
