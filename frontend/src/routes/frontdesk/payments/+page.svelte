<!-- frontend/src/routes/frontdesk/payments/+page.svelte -->
<script lang="ts">
  import { page } from '$app/state';
  import { roomStore, isOnline } from '#lib/stores/roomStore';
  import { authStore } from '#lib/stores/authStore';
  import type { Folio } from '#lib/types/room';
  import { onMount } from 'svelte';

  const folios = roomStore.folios;

  let selectedRoomNumber = $state('102');
  let amount = $state(0);
  let method = $state('Cash');
  let reference = $state('');

  let isSubmitting = $state(false);
  let feedbackMessage = $state<{ text: string; isError: boolean } | null>(null);
  let receiptData = $state<{
    room: string;
    guest: string;
    amount: number;
    method: string;
    ref: string;
    time: string;
    remaining: number;
  } | null>(null);

  // Active rooms with folios
  let availableFolioRooms = $derived(Object.keys($folios).sort());

  let activeFolio = $derived<Folio | null>(
    selectedRoomNumber ? $folios[selectedRoomNumber] || null : null
  );

  let outstandingBalance = $derived(activeFolio ? activeFolio.balance : 0);

  onMount(() => {
    const queryRoom = page.url.searchParams.get('room');
    if (queryRoom && $folios[queryRoom]) {
      selectedRoomNumber = queryRoom;
    } else if (availableFolioRooms.length > 0) {
      selectedRoomNumber = availableFolioRooms[0];
    }
    // Set initial amount to current balance
    if (activeFolio && activeFolio.balance > 0) {
      amount = activeFolio.balance;
    }
  });

  function setPresetAmount(val: number) {
    amount = val;
  }

  function settleFull() {
    if (activeFolio) {
      amount = Math.max(0, activeFolio.balance);
    }
  }

  async function handlePaymentSubmit(e: Event) {
    e.preventDefault();
    if (!selectedRoomNumber) {
      feedbackMessage = { text: 'Please select an active room.', isError: true };
      return;
    }
    if (amount <= 0) {
      feedbackMessage = { text: 'Payment amount must be greater than zero.', isError: true };
      return;
    }

    isSubmitting = true;
    feedbackMessage = null;

    try {
      const res = await roomStore.recordPayment({
        room_number: selectedRoomNumber,
        amount: Number(amount),
        method,
        reference: reference.trim() || undefined
      });

      if (res.success && res.folio) {
        receiptData = {
          room: selectedRoomNumber,
          guest: res.folio.guest_name,
          amount: Number(amount),
          method,
          ref: reference.trim() || `TXN-${Date.now().toString().slice(-6)}`,
          time: new Date().toLocaleTimeString(),
          remaining: res.folio.balance
        };
        feedbackMessage = { text: 'Payment verified and recorded successfully!', isError: false };
      } else {
        feedbackMessage = { text: res.message || 'Payment submission failed.', isError: true };
      }
    } catch (err: any) {
      feedbackMessage = { text: err.message || 'Payment processing error.', isError: true };
    } finally {
      isSubmitting = false;
    }
  }

  function handleReset() {
    receiptData = null;
    feedbackMessage = null;
    reference = '';
    settleFull();
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">PAYMENT TERMINAL</h1>
      <p class="text-xs text-slate-400">Cash, Card & Mobile Money Settlement</p>
    </div>
    <span class="text-xs font-mono bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 px-2.5 py-1 rounded-full font-bold">
      Terminal Active
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
      </div>
    </div>
  {/if}

  {#if receiptData}
    <!-- Payment Receipt Slip Card -->
    <div class="bg-slate-900 border border-emerald-500/40 rounded-3xl p-5 shadow-2xl space-y-4 animate-scale-up">
      <div class="text-center space-y-1">
        <div class="w-12 h-12 rounded-full bg-emerald-500/20 text-emerald-400 flex items-center justify-center mx-auto text-2xl">
          <i class="ph ph-receipt"></i>
        </div>
        <h2 class="text-base font-bold text-white">Payment Receipt</h2>
        <p class="text-xs text-slate-400 font-mono">Ref: {receiptData.ref}</p>
      </div>

      <div class="bg-slate-950 border border-slate-800 rounded-2xl p-4 space-y-2.5 text-xs">
        <div class="flex justify-between">
          <span class="text-slate-400">Room Number:</span>
          <span class="font-bold font-['Oxanium',sans-serif] text-amber-400">Room {receiptData.room}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Guest Name:</span>
          <span class="font-bold text-white">{receiptData.guest}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Method:</span>
          <span class="font-semibold text-slate-200">{receiptData.method}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-slate-400">Cashier:</span>
          <span class="text-slate-300">{$authStore.fullName}</span>
        </div>
        <div class="flex justify-between pt-2 border-t border-slate-800 text-sm font-bold">
          <span class="text-slate-300">Amount Paid:</span>
          <span class="font-mono text-emerald-400">MWK {receiptData.amount.toLocaleString()}</span>
        </div>
        <div class="flex justify-between text-xs font-semibold">
          <span class="text-slate-400">Outstanding Balance:</span>
          <span class="font-mono {receiptData.remaining > 0 ? 'text-amber-400' : 'text-emerald-400'}">
            MWK {receiptData.remaining.toLocaleString()}
          </span>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-2 pt-1">
        {#if receiptData.remaining <= 0}
          <a
            href="/frontdesk/checkout?room={receiptData.room}"
            class="py-2.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs text-center shadow-md active:scale-95 transition"
          >
            Check Out Room &rarr;
          </a>
        {:else}
          <a
            href="/frontdesk/folio?room={receiptData.room}"
            class="py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-bold text-xs text-center border border-slate-700 active:scale-95 transition"
          >
            View Folio
          </a>
        {/if}
        <button
          onclick={handleReset}
          class="py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 font-semibold text-xs active:scale-95 transition"
        >
          New Payment
        </button>
      </div>
    </div>
  {:else}
    <form onsubmit={handlePaymentSubmit} class="bg-slate-900 border border-slate-800 rounded-2xl p-4 space-y-4 shadow-xl">
      <!-- Select Room with Folio -->
      <div>
        <label for="payment-room-select" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
          Select Room & Guest *
        </label>
        <select
          id="payment-room-select"
          bind:value={selectedRoomNumber}
          onchange={settleFull}
          class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2.5 text-xs text-white focus:outline-none focus:border-amber-500 font-medium"
        >
          {#each availableFolioRooms as rNo}
            {@const f = $folios[rNo]}
            <option value={rNo}>
              Room {rNo} &bull; {f.guest_name} &bull; Bal: MWK {f.balance.toLocaleString()}
            </option>
          {/each}
        </select>
      </div>

      <!-- Current Balance Banner -->
      {#if activeFolio}
        <div class="bg-slate-950 border {outstandingBalance > 0 ? 'border-amber-500/40 bg-amber-500/5' : 'border-emerald-500/40 bg-emerald-500/5'} rounded-2xl p-4 flex items-center justify-between">
          <div>
            <span class="text-[10px] uppercase tracking-wider text-slate-400 font-bold block">Outstanding Balance</span>
            <span class="text-2xl font-bold font-['Oxanium',sans-serif] {outstandingBalance > 0 ? 'text-amber-400' : 'text-emerald-400'}">
              MWK {outstandingBalance.toLocaleString()}
            </span>
          </div>
          {#if outstandingBalance > 0}
            <button
              type="button"
              onclick={settleFull}
              class="px-3 py-1.5 rounded-xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs shadow-md active:scale-95 transition"
            >
              Pay Full
            </button>
          {:else}
            <span class="text-xs text-emerald-400 font-semibold flex items-center">
              <i class="ph ph-check-circle mr-1 text-base"></i> Fully Settled
            </span>
          {/if}
        </div>
      {/if}

      <!-- Payment Amount -->
      <div>
        <div class="flex items-center justify-between mb-1">
          <label for="payment-amount-input" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400">
            Payment Amount (MWK) *
          </label>
          <span class="text-[11px] text-slate-500 font-mono">MWK {amount.toLocaleString()}</span>
        </div>
        <input
          id="payment-amount-input"
          type="number"
          min="1"
          step="500"
          bind:value={amount}
          required
          class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2.5 text-base font-bold text-white font-mono focus:outline-none focus:border-amber-500 shadow-inner"
        />

        <!-- Quick Presets -->
        <div class="flex items-center space-x-2 mt-2">
          {#each [10000, 20000, 45000, 75000] as preset}
            <button
              type="button"
              onclick={() => setPresetAmount(preset)}
              class="flex-1 py-1 rounded-lg bg-slate-950 border border-slate-800 text-[10px] text-slate-300 hover:text-white font-mono active:scale-95 transition"
            >
              {(preset / 1000).toFixed(0)}k
            </button>
          {/each}
        </div>
      </div>

      <!-- Payment Method -->
      <div>
        <label class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1.5">
          Payment Method *
        </label>
        <div class="grid grid-cols-3 gap-2">
          {#each ['Cash', 'Visa / POS', 'Airtel Money', 'TNM Mpamba', 'Cheque', 'Bank Transfer'] as m}
            <button
              type="button"
              onclick={() => (method = m)}
              class="py-2 px-1 rounded-xl border text-xs font-semibold text-center transition {method === m
                ? 'bg-amber-500 text-slate-950 font-bold border-amber-500 shadow-md shadow-amber-500/20'
                : 'bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200'}"
            >
              {m}
            </button>
          {/each}
        </div>
      </div>

      <!-- Reference / Authorization Slip Number -->
      <div>
        <label for="payment-ref" class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">
          Reference / Slip / Transaction ID
        </label>
        <input
          id="payment-ref"
          type="text"
          placeholder="e.g. POS-882194 / MP-99201"
          bind:value={reference}
          class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3 py-2 text-xs text-white placeholder-slate-600 focus:outline-none focus:border-amber-500 font-mono"
        />
      </div>

      <!-- Action Button -->
      <button
        type="submit"
        disabled={isSubmitting || amount <= 0}
        class="w-full py-3.5 rounded-xl bg-gradient-to-r from-emerald-500 to-emerald-600 hover:from-emerald-400 hover:to-emerald-500 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-emerald-500/25 active:scale-98 transition disabled:opacity-50"
      >
        {isSubmitting ? 'Recording Payment...' : `Record Payment of MWK ${amount.toLocaleString()}`}
      </button>
    </form>
  {/if}
</div>
