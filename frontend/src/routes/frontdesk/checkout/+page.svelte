<!-- frontend/src/routes/frontdesk/checkout/+page.svelte -->
<script lang="ts">
  import { page } from '$app/state';
  import { roomStore } from '#lib/stores/roomStore';
  import { authStore } from '#lib/stores/authStore';
  import type { Folio, Room } from '#lib/types/room';
  import { onMount } from 'svelte';

  const rooms = roomStore.rooms;
  const folios = roomStore.folios;

  let selectedRoomNumber = $state('');
  let keysReturned = $state(true);
  let roomConditionConfirmed = $state(true);

  let isSubmitting = $state(false);
  let errorMessage = $state('');
  let successReceipt = $state<{
    room: string;
    guest: string;
    checkoutTime: string;
    officer: string;
  } | null>(null);

  // Available occupied rooms
  let occupiedRooms = $derived($rooms.filter((r: Room) => r.status === 'Occupied'));

  $effect(() => {
    if (selectedRoomNumber && !$folios[selectedRoomNumber]) {
      roomStore.fetchFolio(selectedRoomNumber);
    }
  });

  let activeFolio = $derived<Folio | null>(
    selectedRoomNumber ? $folios[selectedRoomNumber] || null : null
  );

  let activeRoom = $derived<Room | null>(
    $rooms.find((r: Room) => r.room_number === selectedRoomNumber) || null
  );

  let isZeroBalance = $derived(
    activeFolio ? activeFolio.balance <= 0 : false
  );

  onMount(() => {
    const queryRoom = page.url.searchParams.get('room');
    if (queryRoom && occupiedRooms.some((r) => r.room_number === queryRoom)) {
      selectedRoomNumber = queryRoom;
    } else if (occupiedRooms.length > 0) {
      selectedRoomNumber = occupiedRooms[0].room_number;
    }
  });

  async function handleCheckOut(e: Event) {
    e.preventDefault();
    if (!selectedRoomNumber) {
      errorMessage = 'Please select a room to check out.';
      return;
    }

    if (!isZeroBalance) {
      errorMessage = `Zero-Balance Lock Active: Outstanding balance of MWK ${activeFolio?.balance.toLocaleString()} remaining. Please settle balance before checking out.`;
      return;
    }

    if (!keysReturned) {
      errorMessage = 'Please confirm room keys have been returned by guest.';
      return;
    }

    errorMessage = '';
    isSubmitting = true;

    try {
      const res = await roomStore.checkOut({ room_number: selectedRoomNumber });
      if (res.success) {
        successReceipt = {
          room: selectedRoomNumber,
          guest: activeFolio?.guest_name || 'Guest',
          checkoutTime: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          officer: $authStore.fullName
        };
      } else {
        errorMessage = res.message || 'Check-out failed.';
      }
    } catch (err: any) {
      errorMessage = err.message || 'Error executing checkout.';
    } finally {
      isSubmitting = false;
    }
  }

  function handleReset() {
    successReceipt = null;
    errorMessage = '';
    if (occupiedRooms.length > 0) {
      selectedRoomNumber = occupiedRooms[0].room_number;
    }
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">CHECK-OUT TERMINAL</h1>
      <p class="text-xs text-slate-400">Room Release & Housekeeping Handover</p>
    </div>
    <!-- 10:00 AM Strict Policy Tag -->
    <span class="text-[10px] font-mono bg-amber-500/10 text-amber-400 border border-amber-500/20 px-2 py-0.5 rounded-full font-bold">
      10:00 AM Rule
    </span>
  </div>

  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs flex items-start space-x-2">
      <i class="ph ph-prohibit text-base text-rose-400 mt-0.5"></i>
      <div class="flex-1">
        <p class="font-bold">Check-Out Blocked</p>
        <p class="mt-0.5">{errorMessage}</p>
        {#if activeFolio && activeFolio.balance > 0}
          <div class="mt-2.5">
            <a
              href="/frontdesk/payments?room={selectedRoomNumber}"
              class="inline-flex items-center px-3 py-1.5 rounded-xl bg-amber-500 text-slate-950 font-bold text-xs shadow-md active:scale-95 transition"
            >
              <i class="ph ph-currency-circle-dollar mr-1"></i> Settle MWK {activeFolio.balance.toLocaleString()} Now &rarr;
            </a>
          </div>
        {/if}
      </div>
    </div>
  {/if}

  {#if successReceipt}
    <!-- Check-Out Clearance Slip / Gate Pass -->
    <div class="bg-slate-900 border border-emerald-500/40 rounded-3xl p-5 shadow-2xl space-y-4 animate-scale-up">
      <div class="text-center space-y-1">
        <div class="w-12 h-12 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center mx-auto text-2xl">
          <i class="ph ph-check-fat"></i>
        </div>
        <h2 class="text-base font-bold text-white">Check-Out Successful</h2>
        <p class="text-xs text-slate-400">Room {successReceipt.room} handed over to Housekeeping</p>
      </div>

      <div class="bg-slate-950 border border-slate-800 rounded-2xl p-4 space-y-2.5 text-xs">
        <div class="flex justify-between">
          <span class="text-slate-400">Guest:</span>
          <span class="font-bold text-white">{successReceipt.guest}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Room Released:</span>
          <span class="font-['Oxanium',sans-serif] font-bold text-amber-400">Room {successReceipt.room}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">New Room Status:</span>
          <span class="text-sky-400 font-bold flex items-center">
            <i class="ph ph-broom mr-1"></i> Cleaning (High Priority)
          </span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Departure Timestamp:</span>
          <span class="font-mono text-slate-200">{successReceipt.checkoutTime}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Cleared By Officer:</span>
          <span class="text-slate-300">{successReceipt.officer}</span>
        </div>
        <div class="flex justify-between pt-2 border-t border-slate-800 text-emerald-400 font-semibold">
          <span>Final Balance:</span>
          <span class="font-mono font-bold">MWK 0.00 (Zero-Balance Verified)</span>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-2 pt-1">
        <a
          href="/frontdesk/housekeeping"
          class="py-2.5 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs text-center shadow-md active:scale-95 transition"
        >
          Track Housekeeping &rarr;
        </a>
        <button
          onclick={handleReset}
          class="py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-xs active:scale-95 transition"
        >
          Checkout Another
        </button>
      </div>
    </div>
  {:else}
    {#if occupiedRooms.length === 0}
      <div class="bg-slate-900 border border-slate-800 rounded-2xl p-8 text-center text-slate-500 text-xs">
        <i class="ph ph-bed text-3xl mb-2 text-slate-600 block"></i>
        No occupied rooms currently in-house for departure.
      </div>
    {:else}
      <form onsubmit={handleCheckOut} class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-4 shadow-xl">
        <!-- Room Selector -->
        <div>
          <label for="checkout-room-select" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
            Select Occupied Room *
          </label>
          <select
            id="checkout-room-select"
            bind:value={selectedRoomNumber}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2.5 text-xs text-white focus:outline-none focus:border-amber-500 font-medium"
          >
            {#each occupiedRooms as r}
              {@const f = $folios[r.room_number]}
              <option value={r.room_number}>
                Room {r.room_number} &bull; {f?.guest_name || 'Guest'} &bull; Bal: MWK {f?.balance.toLocaleString() || 0}
              </option>
            {/each}
          </select>
        </div>

        <!-- Zero-Balance Compliance Inspector -->
        {#if activeFolio}
          <div class="p-3.5 rounded-2xl border {isZeroBalance ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300' : 'bg-rose-500/10 border-rose-500/40 text-rose-300'} text-xs space-y-2">
            <div class="flex items-center justify-between">
              <span class="font-bold flex items-center">
                {#if isZeroBalance}
                  <i class="ph ph-check-circle text-base mr-1.5 text-emerald-400"></i>
                  Zero-Balance Lock Satisfied
                {:else}
                  <i class="ph ph-lock text-base mr-1.5 text-rose-400"></i>
                  Zero-Balance Lock Enforced
                {/if}
              </span>
              <span class="font-mono font-bold text-sm">
                MWK {activeFolio.balance.toLocaleString()}
              </span>
            </div>

            {#if !isZeroBalance}
              <p class="text-[11px] text-rose-400">
                Hotel policy strictly forbids room release until all outstanding charges are settled.
              </p>
              <div class="pt-1">
                <a
                  href="/frontdesk/payments?room={selectedRoomNumber}"
                  class="w-full py-2 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs flex items-center justify-center space-x-1 shadow-md transition"
                >
                  <i class="ph ph-currency-circle-dollar text-sm"></i>
                  <span>Collect MWK {activeFolio.balance.toLocaleString()} Payment</span>
                </a>
              </div>
            {:else}
              <p class="text-[11px] text-emerald-400">
                Folio is fully paid. Guest is eligible for immediate check-out and key return.
              </p>
            {/if}
          </div>
        {/if}

        <!-- Check-Out Protocols & Checkboxes -->
        <div class="space-y-2 bg-slate-950 border border-slate-800 rounded-xl p-3 text-xs">
          <label class="flex items-center space-x-2.5 cursor-pointer">
            <input
              type="checkbox"
              bind:checked={keysReturned}
              class="w-4 h-4 rounded text-amber-500 focus:ring-amber-500 bg-slate-800 border-slate-700"
            />
            <span class="text-slate-300 font-medium">Physical room keycards collected</span>
          </label>

          <label class="flex items-center space-x-2.5 cursor-pointer">
            <input
              type="checkbox"
              bind:checked={roomConditionConfirmed}
              class="w-4 h-4 rounded text-amber-500 focus:ring-amber-500 bg-slate-800 border-slate-700"
            />
            <span class="text-slate-300 font-medium">No missing room inventory / damages reported</span>
          </label>
        </div>

        <!-- 10:00 AM Strict Checkout Enforcement Notice -->
        <div class="text-[11px] text-slate-400 bg-slate-950/60 p-2.5 rounded-xl border border-slate-800/80 flex items-center space-x-2">
          <i class="ph ph-info text-amber-400 text-base"></i>
          <span>Releasing the room automatically sets status to <strong>Cleaning</strong> for the Housekeeping team.</span>
        </div>

        <!-- Action Button -->
        <button
          type="submit"
          disabled={isSubmitting || !isZeroBalance || !keysReturned}
          class="w-full py-3.5 rounded-xl bg-gradient-to-r {isZeroBalance && keysReturned ? 'from-rose-500 to-rose-600 hover:from-rose-400 hover:to-rose-500 text-white cursor-pointer shadow-rose-500/25' : 'from-slate-700 to-slate-800 text-slate-500 cursor-not-allowed'} font-bold text-xs uppercase tracking-wider shadow-lg active:scale-98 transition"
        >
          {#if !isZeroBalance}
            <i class="ph ph-lock mr-1"></i> Balance Must Be Zero (MWK {activeFolio?.balance.toLocaleString()})
          {:else if isSubmitting}
            Processing Release...
          {:else}
            Release Room & Handover to Housekeeping
          {/if}
        </button>
      </form>
    {/if}
  {/if}
</div>
