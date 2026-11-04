<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { authStore } from '#lib/stores/authStore';
  import { guestStore } from '#lib/stores/guestStore';
  import { roomStore } from '#lib/stores/roomStore';
  import { onMount } from 'svelte';
  import type { Folio } from '#lib/types/room';

  const stay = guestStore.stay;
  const folios = roomStore.folios;

  let auth = $state<{ token: string; isLoggedIn: boolean }>({ token: '', isLoggedIn: false });
  authStore.subscribe(v => { auth = v; });

  let isCheckedIn = $state(false);
  let currentFolio = $state<Folio | null>(null);
  let requestSuccessMsg = $state('');

  onMount(async () => {
    if (auth.isLoggedIn && auth.token) {
      try {
        const bRes = await apiFetch('/api/v1/guest/bookings', {
          headers: { 'Authorization': `Bearer ${auth.token}` }
        });
        if (bRes.ok) {
          const bData = await bRes.json();
          if (Array.isArray(bData) && bData.length > 0) {
            const active = bData.find((b: any) => b.folio_status !== 'Closed');
            isCheckedIn = active ? (active.checked_in === true && active.room_status === 'Occupied') : false;
          } else {
            isCheckedIn = false;
          }
        }

        const fRes = await apiFetch('/api/v1/guest/folio', {
          headers: { 'Authorization': `Bearer ${auth.token}` }
        });
        if (fRes.ok) {
          const fData = await fRes.json();
          if (fData.has_folio && fData.payment_status !== 'Closed') {
            currentFolio = {
              id: fData.folio_id,
              guest_id: fData.card_no,
              guest_name: 'Resident',
              guest_phone: '',
              room_number: fData.room_number,
              check_in_date: fData.arrival_date,
              check_out_date: fData.departure_date,
              total_charges: fData.room_charges + fData.incidental_charges,
              total_paid: 0,
              balance: fData.total_balance,
              items: [
                {
                  id: 'it-rm',
                  description: `Room Accommodation (Room ${fData.room_number})`,
                  amount: fData.room_charges,
                  category: 'Accommodation',
                  timestamp: fData.arrival_date
                }
              ]
            };
          } else {
            currentFolio = null;
          }
        }
      } catch {
        isCheckedIn = false;
        currentFolio = null;
      }
    } else {
      isCheckedIn = false;
      currentFolio = null;
    }
  });

  function requestSettlement() {
    requestSuccessMsg = '✓ Settlement request sent to Front Desk. Clerk will prepare your receipt.';
    setTimeout(() => {
      requestSuccessMsg = '';
    }, 4000);
  }

  function handlePrint() {
    window.print();
  }

  function fmtMwk(n: number) {
    return 'MWK ' + Math.abs(n).toLocaleString();
  }

  function getCategoryStyle(cat: string) {
    switch (cat) {
      case 'Accommodation':
        return 'bg-blue-500/20 text-blue-300 border-blue-500/30';
      case 'Restaurant':
        return 'bg-orange-500/20 text-orange-300 border-orange-500/30';
      case 'Bar':
        return 'bg-purple-500/20 text-purple-300 border-purple-500/30';
      case 'Payment':
        return 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30';
      default:
        return 'bg-slate-800 text-slate-300 border-slate-700';
    }
  }
</script>

<div class="space-y-4 pb-4">
  <!-- Title & Print Header -->
  <div class="pt-1 flex items-center justify-between no-print">
    <div>
      <h1 class="text-xl font-bold text-white tracking-tight">My Room Statement</h1>
      <p class="text-xs text-slate-400">Live Folio for Room {$stay.roomNumber}</p>
    </div>
    <button
      onclick={handlePrint}
      class="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 text-xs font-semibold active:scale-95 transition"
    >
      <i class="ph ph-printer text-sm"></i>
      <span>Print / PDF</span>
    </button>
  </div>

  {#if requestSuccessMsg}
    <div class="p-3 rounded-2xl bg-emerald-500/20 border border-emerald-500/40 text-emerald-300 text-xs text-center font-bold">
      {requestSuccessMsg}
    </div>
  {/if}

  {#if !isCheckedIn}
    <div class="p-3.5 rounded-2xl bg-amber-500/15 border border-amber-500/30 text-amber-300 text-xs text-center space-y-1 no-print">
      <div class="flex items-center justify-center space-x-1.5 font-bold">
        <i class="ph ph-lock text-sm"></i>
        <span>Pending Front Desk Check-In</span>
      </div>
      <p class="text-[11px] text-slate-300">
        This is an estimated booking statement. Active in-house folio charges will be confirmed upon physical check-in at the Front Desk Reception.
      </p>
    </div>
  {/if}

  <!-- Balance Banner Card -->
  {#if isCheckedIn && currentFolio}
    <div class="bg-gradient-to-br from-indigo-950 via-slate-900 to-slate-900 border border-indigo-700/50 rounded-3xl p-5 shadow-2xl relative overflow-hidden">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs text-indigo-300 uppercase font-bold tracking-wider">Current Balance Due</span>
        <span class="text-[10px] px-2 py-0.5 rounded-full border font-bold uppercase tracking-wider {currentFolio.balance <= 0 ? 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30' : 'bg-amber-500/20 text-amber-300 border-amber-500/40'}">
          {currentFolio.balance <= 0 ? 'Settled in Full' : 'Open Balance'}
        </span>
      </div>

      <p class="text-3xl font-extrabold text-white font-['Oxanium',sans-serif] tracking-tight">
        MWK {currentFolio.balance.toLocaleString()}
      </p>

      <!-- Mini Breakdown -->
      <div class="grid grid-cols-2 gap-3 mt-4 pt-3 border-t border-slate-800 text-xs">
        <div>
          <span class="text-[10px] text-slate-400 block">Total Incurred Charges</span>
          <span class="font-bold text-white font-['Oxanium',sans-serif]">MWK {currentFolio.total_charges.toLocaleString()}</span>
        </div>
        <div>
          <span class="text-[10px] text-slate-400 block">Total Payments / Deposits</span>
          <span class="font-bold text-emerald-400 font-['Oxanium',sans-serif]">MWK {currentFolio.total_paid.toLocaleString()}</span>
        </div>
      </div>
    </div>

    <!-- Official Itemized Ledger Slip -->
    <div class="bg-slate-900 border border-slate-800 rounded-3xl p-4 shadow-lg space-y-3 print:bg-white print:text-black print:border-black">
      <div class="flex items-center justify-between border-b border-slate-800 pb-2">
        <div>
          <h2 class="text-xs font-bold text-white uppercase tracking-wider">Itemized Line Items</h2>
          <p class="text-[11px] text-slate-400">Card Ref: {$stay.cardNo}</p>
        </div>
        <span class="text-xs text-amber-400 font-mono font-semibold">{currentFolio.items.length} Transactions</span>
      </div>

      <!-- Line items list -->
      <div class="space-y-2">
        {#each currentFolio.items as item}
          <div class="p-3 rounded-2xl bg-slate-950/60 border border-slate-800/80 space-y-1 text-xs">
            <div class="flex items-start justify-between">
              <span class="font-semibold text-white leading-tight">{item.description}</span>
              <span class="font-bold font-['Oxanium',sans-serif] {item.amount < 0 ? 'text-emerald-400' : 'text-slate-200'}">
                {item.amount < 0 ? '-' : ''}{fmtMwk(item.amount)}
              </span>
            </div>

            <div class="flex items-center justify-between pt-1 text-[10px]">
              <span class="px-2 py-0.5 rounded-full border font-bold uppercase {getCategoryStyle(item.category)}">
                {item.category}
              </span>
              <span class="text-slate-400">{item.timestamp ? item.timestamp.split('T')[0] : 'During Stay'}</span>
            </div>
          </div>
        {/each}
      </div>

      <!-- Final Totals Footer -->
      <div class="pt-3 border-t border-slate-800 space-y-1.5 text-xs">
        <div class="flex justify-between text-slate-400">
          <span>Accommodation Charges</span>
          <span class="font-mono text-white">MWK {currentFolio.total_charges.toLocaleString()}</span>
        </div>
        <div class="flex justify-between text-slate-400">
          <span>Advance Deposit Paid</span>
          <span class="font-mono text-emerald-400">-MWK {currentFolio.total_paid.toLocaleString()}</span>
        </div>
        <div class="flex justify-between text-sm font-bold text-white pt-1.5 border-t border-slate-800">
          <span>Net Balance Due</span>
          <span class="text-amber-400 font-['Oxanium',sans-serif]">MWK {currentFolio.balance.toLocaleString()}</span>
        </div>
      </div>
    </div>

    <!-- Settlement Actions (Hidden on Print) -->
    <div class="space-y-2 no-print">
      <button
        onclick={requestSettlement}
        class="w-full py-3.5 rounded-2xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/25 active:scale-98 transition flex items-center justify-center space-x-1.5"
      >
        <i class="ph ph-credit-card text-base"></i>
        <span>Request Settlement at Front Desk</span>
      </button>

      <p class="text-[11px] text-center text-slate-400">
        Questions regarding your bill? Dial <a href="tel:0" class="text-amber-400 underline font-semibold">Ext. 0</a> for 24/7 Front Desk assistance.
      </p>
    </div>
  {/if}
</div>
