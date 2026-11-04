<script lang="ts">
  import { apiFetch } from '#lib/api';
  import { authStore } from '#lib/stores/authStore';
  import { guestStore } from '#lib/stores/guestStore';
  import { onMount } from 'svelte';

  const stay = guestStore.stay;
  const services = guestStore.services;

  let auth = $state<{ token: string; isLoggedIn: boolean }>({ token: '', isLoggedIn: false });
  authStore.subscribe(v => { auth = v; });

  let isCheckedIn = $state(false);
  let customRequestText = $state('');
  let isSubmitting = $state(false);
  let toastMsg = $state('');

  onMount(async () => {
    if (auth.isLoggedIn && auth.token) {
      try {
        const res = await apiFetch('/api/v1/guest/bookings', {
          headers: { 'Authorization': `Bearer ${auth.token}` }
        });
        if (res.ok) {
          const data = await res.json();
          if (Array.isArray(data) && data.length > 0) {
            const active = data.find((b: any) => b.folio_status !== 'Closed');
            isCheckedIn = active ? (active.checked_in === true && active.room_status === 'Occupied') : false;
          } else {
            isCheckedIn = false;
          }
        }
      } catch {
        isCheckedIn = false;
      }
    } else {
      isCheckedIn = false;
    }
  });

  // Service presets
  const serviceOptions = [
    {
      title: 'Fresh Bath Towels',
      description: '2 extra plush bath sheets & hand towels',
      icon: 'ph-sparkle',
      color: 'bg-emerald-500/20 text-emerald-400'
    },
    {
      title: 'Pillow Menu & Linen Change',
      description: 'Extra hypoallergenic pillows & fresh sheets',
      icon: 'ph-bed',
      color: 'bg-blue-500/20 text-blue-400'
    },
    {
      title: 'Morning Wake-Up Call',
      description: 'Scheduled phone call from Front Desk at 07:00 AM',
      icon: 'ph-alarm',
      color: 'bg-amber-500/20 text-amber-400'
    },
    {
      title: 'Luggage Assistance / Porter',
      description: 'Bellhop to assist with bags and checkout departure',
      icon: 'ph-suitcase-rolling',
      color: 'bg-purple-500/20 text-purple-400'
    },
    {
      title: 'Late Check-out Request',
      description: 'Request extension until 12:00 PM (Subject to availability)',
      icon: 'ph-clock-countdown',
      color: 'bg-rose-500/20 text-rose-400'
    }
  ];

  function sendQuickRequest(title: string, desc: string, icon: string) {
    guestStore.services.addRequest(title, desc, icon);
    toastMsg = `✓ Request for "${title}" dispatched to Front Desk team!`;
    setTimeout(() => {
      toastMsg = '';
    }, 3500);
  }

  function submitCustomRequest(e: Event) {
    e.preventDefault();
    if (!customRequestText.trim()) return;
    isSubmitting = true;

    guestStore.services.addRequest('Concierge Assistance', customRequestText.trim(), 'ph-chat-circle-dots');
    customRequestText = '';
    isSubmitting = false;
    toastMsg = '✓ Your request has been sent to the duty concierge!';
    setTimeout(() => {
      toastMsg = '';
    }, 3500);
  }

  function getStatusStyle(status: string) {
    switch (status) {
      case 'Pending':
        return 'bg-amber-500/20 text-amber-300 border-amber-500/40 animate-pulse';
      case 'Attending':
        return 'bg-blue-500/20 text-blue-300 border-blue-500/40';
      case 'Completed':
        return 'bg-emerald-500/20 text-emerald-300 border-emerald-500/30';
      default:
        return 'bg-slate-800 text-slate-300 border-slate-700';
    }
  }
</script>

<div class="space-y-4 pb-4">
  <!-- Title -->
  <div class="pt-1 flex items-start justify-between">
    <div>
      <h1 class="text-xl font-bold text-white tracking-tight">Concierge & Requests</h1>
      <p class="text-xs text-slate-400">Personalized resident service for Room {$stay.roomNumber}</p>
    </div>
    <a
      href="tel:0"
      class="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl bg-indigo-600 hover:bg-indigo-500 text-white text-xs font-bold active:scale-95 transition shadow-md shadow-indigo-600/30"
    >
      <i class="ph ph-phone-call text-sm"></i>
      <span>Dial Ext. 0</span>
    </a>
  </div>

  {#if toastMsg}
    <div class="p-3 rounded-2xl bg-emerald-500/20 border border-emerald-500/40 text-emerald-300 text-xs text-center font-bold shadow-lg">
      {toastMsg}
    </div>
  {/if}

  {#if !isCheckedIn}
    <div class="p-3.5 rounded-2xl bg-amber-500/15 border border-amber-500/30 text-amber-300 text-xs text-center space-y-1">
      <div class="flex items-center justify-center space-x-1.5 font-bold">
        <i class="ph ph-lock text-sm"></i>
        <span>Front Desk Check-In Required</span>
      </div>
      <p class="text-[11px] text-slate-300">
        Guest amenities and concierge requests activate once the Front Desk completes your room check-in.
      </p>
    </div>
  {/if}

  <!-- 1-Tap Service Presets Grid -->
  <div>
    <h2 class="text-xs font-bold text-slate-400 uppercase tracking-wider mb-2.5">Fast Service Requests</h2>
    <div class="grid grid-cols-2 gap-3">
      {#each serviceOptions as srv}
        <button
          onclick={() => sendQuickRequest(srv.title, srv.description, srv.icon)}
          disabled={!isCheckedIn}
          class="bg-slate-900 border border-slate-800 hover:border-amber-500/40 rounded-2xl p-3.5 text-left transition active:scale-95 shadow-md flex flex-col justify-between h-32 group disabled:opacity-40 disabled:cursor-not-allowed"
        >
          <div class="w-8 h-8 rounded-xl {srv.color} flex items-center justify-center">
            <i class="ph {srv.icon} text-lg"></i>
          </div>
          <div>
            <p class="text-xs font-bold text-white group-hover:text-amber-300 transition-colors leading-tight">{srv.title}</p>
            <p class="text-[10px] text-slate-400 mt-1 line-clamp-2 leading-tight">{srv.description}</p>
          </div>
        </button>
      {/each}
    </div>
  </div>

  <!-- Custom Request Form -->
  <div class="bg-slate-900 border border-slate-800 rounded-3xl p-4 shadow-lg space-y-3">
    <div class="flex items-center space-x-2">
      <i class="ph ph-chat-circle-dots text-amber-400 text-lg"></i>
      <h2 class="text-xs font-bold text-white uppercase tracking-wider">Custom Concierge Inquiry</h2>
    </div>

    <form onsubmit={submitCustomRequest} class="space-y-3">
      <div>
        <label for="custom-request-text" class="sr-only">Special concierge request</label>
        <textarea
          id="custom-request-text"
          rows="3"
          bind:value={customRequestText}
          placeholder="E.g., Please arrange an iron and ironing board, or ice bucket to Room 201..."
          class="w-full px-4 py-2.5 bg-slate-950 border border-slate-700 rounded-2xl text-white text-xs focus:outline-none focus:border-amber-500 placeholder-slate-600 leading-relaxed"
        ></textarea>
      </div>

      <button
        type="submit"
        disabled={isSubmitting || !customRequestText.trim() || !isCheckedIn}
        class="w-full py-3 rounded-2xl bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold text-xs uppercase tracking-wider transition active:scale-98 disabled:opacity-50 flex items-center justify-center space-x-1"
      >
        <i class="ph ph-paper-plane-tilt text-sm"></i>
        <span>{isSubmitting ? 'Sending...' : 'Send Message to Front Desk'}</span>
      </button>
    </form>
  </div>

  <!-- Submitted Requests Log -->
  {#if $services.length > 0}
    <div class="space-y-2">
      <h2 class="text-xs font-bold text-slate-400 uppercase tracking-wider">My Requests Log</h2>
      <div class="space-y-2">
        {#each $services as req}
          <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5 flex items-center justify-between text-xs shadow-sm">
            <div class="flex items-start space-x-2.5">
              <div class="w-8 h-8 rounded-xl bg-slate-800 text-amber-400 flex items-center justify-center flex-shrink-0 mt-0.5">
                <i class="ph {req.icon} text-base"></i>
              </div>
              <div>
                <p class="font-bold text-white">{req.type}</p>
                <p class="text-slate-400 text-[11px] leading-tight mt-0.5">{req.details}</p>
                <span class="text-[10px] text-slate-400 block mt-1">{req.createdAt}</span>
              </div>
            </div>

            <span class="text-[10px] px-2 py-0.5 rounded-full border font-bold uppercase tracking-wider flex-shrink-0 {getStatusStyle(req.status)}">
              {req.status}
            </span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
