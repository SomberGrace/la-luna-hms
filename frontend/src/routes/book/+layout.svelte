<script lang="ts">
  let { children } = $props();
</script>

<div class="min-h-screen bg-slate-950 text-slate-100 font-['Inter',sans-serif] select-none">
  <!-- Top Header -->
  <header class="sticky top-0 z-40 bg-slate-900/90 backdrop-blur-md border-b border-slate-800/80 px-4 py-3 shadow-lg">
    <div class="max-w-4xl mx-auto flex items-center justify-between">
      <a href="/book" class="flex items-center space-x-2.5 active:scale-95 transition-transform">
        <div class="w-9 h-9 rounded-xl bg-gradient-to-br from-amber-400 via-amber-500 to-amber-600 flex items-center justify-center shadow-lg shadow-amber-500/20 text-slate-950 font-bold">
          <i class="ph ph-moon text-xl"></i>
        </div>
        <div>
          <span class="font-['Oxanium',sans-serif] font-bold tracking-wider text-base text-white">LA LUNA</span>
          <p class="text-[10px] text-slate-400 font-medium">Book Your Stay</p>
        </div>
      </a>

      <div class="flex items-center space-x-2">
        <a
          href="/book/login"
          class="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl text-xs font-bold bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 active:scale-95 transition"
        >
          <i class="ph ph-sign-in text-sm"></i>
          <span>Login</span>
        </a>
        <a
          href="/"
          class="p-2 rounded-xl bg-slate-800 text-slate-400 hover:text-white border border-slate-700/60 active:scale-95 transition"
          title="Back to Staff Portal"
        >
          <i class="ph ph-arrow-left text-base"></i>
        </a>
      </div>
    </div>
  </header>

  <main class="w-full max-w-4xl mx-auto px-4 py-6">
    {@render children()}
  </main>

  <!-- Footer -->
  <footer class="border-t border-slate-800 py-6 text-center">
    <p class="text-xs text-slate-500">La Luna Hotel · Lilongwe, Malawi · Tel: +265 1 234 567</p>
    <p class="text-[10px] text-slate-600 mt-1">© 2025 La Luna Hotel Management System</p>
  </footer>
</div>
