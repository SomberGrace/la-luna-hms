<script lang="ts">
  import { guestStore } from '#lib/stores/guestStore';

  const stay = guestStore.stay;

  const schedules = [
    {
      title: 'Luna Main Dining Room',
      hours: '06:30 - 22:00',
      details: 'Buffet breakfast (06:30-10:00), A la carte lunch & dinner',
      icon: 'ph-fork-knife',
      badge: 'Level 1'
    },
    {
      title: 'Moonlit Bar & Lounge',
      hours: '10:00 - 23:00',
      details: 'Craft cocktails, chilled Malawian beers, sunset wines',
      icon: 'ph-wine',
      badge: 'Poolside'
    },
    {
      title: 'Infinity Swimming Pool',
      hours: '07:00 - 19:00',
      details: 'Complimentary pool towels available at towel station',
      icon: 'ph-swimming-pool',
      badge: 'Beachfront'
    },
    {
      title: 'Fitness & Wellness Studio',
      hours: '06:00 - 21:00',
      details: 'Cardio equipment, free weights, yoga mats',
      icon: 'ph-barbell',
      badge: 'Level 2'
    }
  ];

  const extensions = [
    { name: 'Front Desk Reception', ext: 'Ext. 0', icon: 'ph-phone' },
    { name: 'In-Room Dining Kitchen', ext: 'Ext. 2', icon: 'ph-cooking-pot' },
    { name: 'Housekeeping & Linen', ext: 'Ext. 3', icon: 'ph-sparkle' },
    { name: 'General Manager On Duty', ext: 'Ext. 4', icon: 'ph-briefcase' },
    { name: 'Emergency & First Aid', ext: 'Ext. 9', icon: 'ph-first-aid' }
  ];
</script>

<div class="space-y-4 pb-4">
  <!-- Title -->
  <div class="pt-1">
    <h1 class="text-xl font-bold text-white tracking-tight">Resort Guide & Directory</h1>
    <p class="text-xs text-slate-400">Everything you need to know about La Luna Resort</p>
  </div>

  <!-- Essential Policies Warning Card -->
  <div class="bg-gradient-to-r from-amber-500/10 via-amber-600/15 to-transparent border border-amber-500/30 rounded-3xl p-4 shadow-lg space-y-2">
    <div class="flex items-center space-x-2 text-amber-300">
      <i class="ph ph-clock-countdown text-xl"></i>
      <h2 class="text-xs font-bold uppercase tracking-wider">Key Hotel Stay Policies</h2>
    </div>
    <ul class="text-xs text-slate-300 space-y-1.5 list-disc list-inside">
      <li><strong>Check-out Time:</strong> Strict <strong>10:00 AM</strong> departure policy.</li>
      <li><strong>Zero-Balance Settlement:</strong> All folio balances must be settled at Front Desk before departure.</li>
      <li><strong>Quiet Relaxation Hours:</strong> Please observe quiet hours from 22:30 PM.</li>
      <li><strong>Non-Smoking Rooms:</strong> All suites are 100% non-smoking indoors.</li>
    </ul>
  </div>

  <!-- Facility Timings -->
  <div class="space-y-2.5">
    <h2 class="text-xs font-bold text-slate-400 uppercase tracking-wider">Dining & Facility Hours</h2>
    <div class="space-y-2">
      {#each schedules as item}
        <div class="bg-slate-900 border border-slate-800 rounded-2xl p-3.5 flex items-start space-x-3 shadow-md">
          <div class="w-9 h-9 rounded-xl bg-slate-800 text-amber-400 flex items-center justify-center flex-shrink-0 mt-0.5">
            <i class="ph {item.icon} text-lg"></i>
          </div>
          <div class="flex-1">
            <div class="flex items-center justify-between mb-0.5">
              <h3 class="font-bold text-white text-xs">{item.title}</h3>
              <span class="text-[9px] px-2 py-0.5 rounded-full bg-slate-800 text-slate-300 border border-slate-700 font-semibold">
                {item.badge}
              </span>
            </div>
            <p class="text-xs text-amber-300 font-semibold font-mono">{item.hours}</p>
            <p class="text-[11px] text-slate-400 mt-0.5 leading-relaxed">{item.details}</p>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <!-- Internal Telephone Extensions -->
  <div class="bg-slate-900 border border-slate-800 rounded-3xl p-4 shadow-lg space-y-3">
    <h2 class="text-xs font-bold text-white uppercase tracking-wider">Room Phone Speed Dial</h2>
    <div class="grid grid-cols-1 gap-2">
      {#each extensions as item}
        <div class="bg-slate-950/70 border border-slate-800/80 rounded-xl p-2.5 flex items-center justify-between text-xs">
          <div class="flex items-center space-x-2">
            <i class="ph {item.icon} text-amber-400 text-sm"></i>
            <span class="text-slate-300 font-medium">{item.name}</span>
          </div>
          <span class="font-bold text-white font-['Oxanium',sans-serif] px-2 py-0.5 rounded bg-slate-800 text-[11px]">
            {item.ext}
          </span>
        </div>
      {/each}
    </div>
  </div>

  <!-- Lake Malawi Local Experiences -->
  <div class="bg-slate-900 border border-slate-800 rounded-3xl p-4 shadow-lg space-y-2.5">
    <div class="flex items-center space-x-2 text-indigo-300">
      <i class="ph ph-sun-horizon text-lg"></i>
      <h2 class="text-xs font-bold uppercase tracking-wider">Lake Malawi Excursions</h2>
    </div>
    <p class="text-xs text-slate-400 leading-relaxed">
      Inquire with the Front Desk Concierge to arrange any of our private guided Lake Malawi experiences:
    </p>
    <div class="space-y-1.5 text-xs text-slate-300">
      <div class="flex items-center space-x-2">
        <i class="ph ph-sailboat text-amber-400"></i>
        <span>Sunset Catamaran Cruise with sparkling wine</span>
      </div>
      <div class="flex items-center space-x-2">
        <i class="ph ph-fish text-cyan-400"></i>
        <span>Snorkeling & Diving at Otter Point (Cichlid fish reserve)</span>
      </div>
      <div class="flex items-center space-x-2">
        <i class="ph ph-basket text-orange-400"></i>
        <span>Mangochi Cultural Village & Artisan Market Tour</span>
      </div>
    </div>
  </div>
</div>

