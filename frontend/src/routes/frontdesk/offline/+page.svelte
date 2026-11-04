<!-- frontend/src/routes/frontdesk/offline/+page.svelte -->
<script lang="ts">
  import { roomStore, isOnline, pendingOfflineCount } from '#lib/stores/roomStore';
  import type { OfflineAction } from '#lib/types/room';

  const offlineQueue = roomStore.offlineQueue;
  const syncStatus = roomStore.syncStatus;

  let isSyncing = $state(false);
  let syncResult = $state<{ synced: number; failed: number } | null>(null);

  async function handleTriggerSync() {
    isSyncing = true;
    syncResult = null;
    try {
      const res = await roomStore.syncOfflineQueue();
      syncResult = res;
    } catch (e) {
      console.error(e);
    } finally {
      isSyncing = false;
    }
  }

  function handleSimulateOfflineToggle() {
    isOnline.update((val) => !val);
  }

  function handleClearQueue() {
    if (confirm('Are you sure you want to clear the offline sync queue? Any unsynced changes will be lost.')) {
      roomStore.clearQueue();
      syncResult = null;
    }
  }

  function getActionBadge(type: string): { label: string; color: string } {
    switch (type) {
      case 'checkin':
        return { label: 'Check-In', color: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30' };
      case 'payment':
        return { label: 'Payment', color: 'bg-teal-500/20 text-teal-400 border-teal-500/30' };
      case 'checkout':
        return { label: 'Check-Out', color: 'bg-rose-500/20 text-rose-400 border-rose-500/30' };
      case 'extension':
        return { label: 'Extension', color: 'bg-indigo-500/20 text-indigo-400 border-indigo-500/30' };
      case 'room_status':
        return { label: 'Status Update', color: 'bg-sky-500/20 text-sky-400 border-sky-500/30' };
      default:
        return { label: type, color: 'bg-slate-700 text-slate-300 border-slate-600' };
    }
  }
</script>

<div class="space-y-4">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-xl font-bold font-['Oxanium',sans-serif] text-white">OFFLINE SYNC & QUEUE</h1>
      <p class="text-xs text-slate-400">Network Resiliency & Local Power-Grid Fallback</p>
    </div>
    <button
      onclick={handleSimulateOfflineToggle}
      class="px-2.5 py-1 rounded-xl text-[10px] font-bold border transition-colors {$isOnline
        ? 'bg-slate-800 text-slate-400 border-slate-700 hover:text-white'
        : 'bg-amber-500/20 text-amber-300 border-amber-500/40'}"
      title="Test offline behavior"
    >
      {$isOnline ? 'Test Offline Mode' : 'Restore Online Mode'}
    </button>
  </div>

  <!-- Network Connection Status Banner -->
  <div class="bg-gradient-to-r from-slate-900 to-slate-950 border {$isOnline ? 'border-emerald-600/40' : 'border-rose-600/60'} rounded-3xl p-5 shadow-xl space-y-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <div class="w-12 h-12 rounded-2xl {$isOnline ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/30 animate-pulse'} flex items-center justify-center text-2xl">
          <i class="ph {$isOnline ? 'ph-wifi-high' : 'ph-wifi-slash'}"></i>
        </div>
        <div>
          <h2 class="text-base font-bold text-white flex items-center space-x-2">
            <span>{$isOnline ? 'Online - Live Central Server' : 'Offline - Local Fallback Active'}</span>
          </h2>
          <p class="text-xs text-slate-400 mt-0.5">
            {$isOnline
              ? 'Transactions communicate directly with the local Go microservice.'
              : 'Local cache active. Transactions are locked and queued for sync.'}
          </p>
        </div>
      </div>
    </div>

    <!-- Sync Action Button -->
    <div class="pt-2 border-t border-slate-800/80 flex items-center justify-between">
      <div class="text-xs font-mono">
        <span class="text-slate-400">Queued:</span>
        <strong class="text-amber-400 font-bold ml-1">{$pendingOfflineCount} operations</strong>
      </div>

      <button
        type="button"
        onclick={handleTriggerSync}
        disabled={isSyncing || $pendingOfflineCount === 0 || !$isOnline}
        class="px-4 py-2 rounded-xl bg-gradient-to-r from-amber-500 to-amber-600 hover:from-amber-400 hover:to-amber-500 text-slate-950 font-bold text-xs uppercase tracking-wider shadow-lg shadow-amber-500/20 active:scale-95 transition disabled:opacity-50 flex items-center space-x-1.5"
      >
        <i class="ph ph-arrows-clockwise text-sm {isSyncing ? 'animate-spin' : ''}"></i>
        <span>{isSyncing ? 'Syncing...' : 'Sync Queue Now'}</span>
      </button>
    </div>
  </div>

  {#if syncResult}
    <div class="p-3.5 rounded-2xl border text-xs {syncResult.failed === 0 ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300' : 'bg-amber-500/10 border-amber-500/30 text-amber-300'} flex items-center justify-between">
      <span>Successfully synchronized {syncResult.synced} transaction(s). {syncResult.failed > 0 ? `(${syncResult.failed} failed/server offline)` : ''}</span>
      <i class="ph {syncResult.failed === 0 ? 'ph-check-circle text-emerald-400' : 'ph-warning text-amber-400'} text-base"></i>
    </div>
  {/if}

  <!-- Queue Items List -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-xl space-y-3">
    <div class="flex items-center justify-between pb-2 border-b border-slate-800">
      <h3 class="text-xs font-bold uppercase tracking-wider text-slate-300">Pending Sync Queue</h3>
      {#if $offlineQueue.length > 0}
        <button
          onclick={handleClearQueue}
          class="text-[11px] text-rose-400 hover:text-rose-300 font-semibold"
        >
          Clear Queue
        </button>
      {/if}
    </div>

    <div class="space-y-2.5 max-h-64 overflow-y-auto pr-1">
      {#if $offlineQueue.length === 0}
        <div class="text-center py-6 text-slate-500 text-xs">
          <i class="ph ph-tray text-3xl mb-1 text-slate-600 block"></i>
          Queue is empty. All front desk records are up to date!
        </div>
      {:else}
        {#each $offlineQueue as item (item.id)}
          {@const badge = getActionBadge(item.action_type)}
          <div class="bg-slate-950/70 border border-slate-800 rounded-xl p-3 flex items-start justify-between text-xs">
            <div class="space-y-1">
              <div class="flex items-center space-x-2">
                <span class="text-[10px] font-bold uppercase px-2 py-0.5 rounded border {badge.color}">
                  {badge.label}
                </span>
                <span class="font-bold text-white font-['Oxanium',sans-serif]">
                  Room {item.payload?.room_number || 'N/A'}
                </span>
              </div>
              <p class="text-[11px] text-slate-300">
                {#if item.action_type === 'checkin'}
                  Guest: {item.payload?.guest_name} &bull; {item.payload?.nights} Night(s) &bull; Dep: MWK {item.payload?.deposit_paid?.toLocaleString() || 0}
                {:else if item.action_type === 'payment'}
                  Amount: MWK {item.payload?.amount?.toLocaleString()} ({item.payload?.method})
                {:else if item.action_type === 'checkout'}
                  Release Room & Transfer to Housekeeping
                {:else if item.action_type === 'extension'}
                  Extension: +{item.payload?.additional_nights} Night(s)
                {:else if item.action_type === 'room_status'}
                  New Status: {item.payload?.status}
                {/if}
              </p>
              <span class="text-[10px] text-slate-500 font-mono block">
                {new Date(item.timestamp).toLocaleTimeString()} &bull; ID: {item.id}
              </span>
              {#if item.error_message}
                <span class="text-[10px] text-rose-400 block font-semibold">
                  Error: {item.error_message}
                </span>
              {/if}
            </div>

            <span class="text-[10px] font-mono px-2 py-0.5 rounded {item.status === 'failed' ? 'bg-rose-500/10 text-rose-400 border border-rose-500/20' : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'}">
              {item.status}
            </span>
          </div>
        {/each}
      {/if}
    </div>
  </div>

  <!-- Reception SOP Guide for Power / Internet Outages -->
  <div class="bg-slate-900 border border-slate-800 rounded-2xl p-4 shadow-xl space-y-2 text-xs">
    <div class="flex items-center space-x-2 text-amber-400 font-bold">
      <i class="ph ph-shield-check text-base"></i>
      <h3 class="uppercase tracking-wider">Outage Standard Operating Procedure (SOP)</h3>
    </div>
    <ul class="space-y-1.5 text-slate-400 text-[11px] list-disc list-inside">
      <li><strong>Double-Booking Prevention:</strong> Rooms assigned during an outage are immediately locked in local memory. Do not reassign the room to another guest.</li>
      <li><strong>Cash / Mobile Receipts:</strong> Collect cash or verify mobile transaction SMS on front desk phone before confirming check-in.</li>
      <li><strong>Automatic Recovery:</strong> The moment power or local network restores, click <em>Sync Queue Now</em> to commit transactions to the database.</li>
    </ul>
  </div>
</div>
