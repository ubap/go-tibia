<script>
  import { bot } from '$lib/botStore.svelte';
</script>

<div class="space-y-6">
  <!-- Page Header -->
  <header>
    <h1 class="text-2xl font-bold text-slate-100">Character Stats</h1>
    <p class="text-slate-400 text-sm">Real-time data from game session.</p>
  </header>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">

    <!-- Character Name Card -->
    <div class="bg-slate-900 border border-slate-800 p-5 rounded-xl shadow-sm">
      <span class="text-xs font-bold text-slate-500 uppercase tracking-widest">Character Name</span>
      <div class="mt-1 text-xl font-semibold text-orange-500">
        {bot.name}
      </div>
    </div>

    <!-- Position Card -->
    <div class="bg-slate-900 border border-slate-800 p-5 rounded-xl shadow-sm">
      <span class="text-xs font-bold text-slate-500 uppercase tracking-widest text-block mb-3">World Position</span>

      <div class="flex items-center gap-4 mt-2">
        <!-- X Coordinate -->
        <div class="flex flex-col">
          <span class="text-[10px] text-slate-500 font-mono">X</span>
          <span class="text-lg font-mono text-slate-200 bg-slate-800 px-2 py-1 rounded">{bot.x}</span>
        </div>

        <!-- Y Coordinate -->
        <div class="flex flex-col">
          <span class="text-[10px] text-slate-500 font-mono">Y</span>
          <span class="text-lg font-mono text-slate-200 bg-slate-800 px-2 py-1 rounded">{bot.y}</span>
        </div>

        <!-- Z Coordinate (Floor) -->
        <div class="flex flex-col">
          <span class="text-[10px] text-slate-500 font-mono">Z</span>
          <span class="text-lg font-mono text-orange-400 bg-slate-800 px-2 py-1 rounded">{bot.z}</span>
        </div>
      </div>
    </div>

  </div>
</div>