<script lang="ts">
  import { onMount } from 'svelte';

  type Theme = 'light' | 'dark' | 'colorpop';
  const themes: { id: Theme; label: string; description: string }[] = [
    { id: 'light', label: 'Luna Light', description: 'Warm, quiet neutrals' },
    { id: 'dark', label: 'Midnight', description: 'Low-light navy surfaces' },
    { id: 'colorpop', label: 'Color Pop', description: 'Cream with playful pastels' }
  ];

  let theme = $state<Theme>('light');
  let settingsOpen = $state(false);

  onMount(() => {
    try {
      const savedTheme = localStorage.getItem('la-luna-theme');
      if (savedTheme === 'dark' || savedTheme === 'light' || savedTheme === 'colorpop') {
        theme = savedTheme;
        document.documentElement.dataset.theme = savedTheme;
      }
    } catch {
      document.documentElement.dataset.theme = theme;
    }
  });

  function chooseTheme(selectedTheme: Theme) {
    theme = selectedTheme;
    document.documentElement.dataset.theme = selectedTheme;
    settingsOpen = false;
    try {
      localStorage.setItem('la-luna-theme', selectedTheme);
    } catch {
      // The selected theme remains active until the page is reloaded.
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') settingsOpen = false;
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="theme-settings">
  {#if settingsOpen}
    <section class="theme-menu" aria-label="Appearance settings">
      <div class="theme-menu-heading">
        <div>
          <p class="theme-menu-eyebrow">PERSONALIZE</p>
          <h2>Choose your look</h2>
        </div>
        <button type="button" class="theme-menu-close" onclick={() => settingsOpen = false} aria-label="Close appearance settings">
          <i class="ph ph-x" aria-hidden="true"></i>
        </button>
      </div>

      <div class="theme-options" role="group" aria-label="App theme">
        {#each themes as option}
          <button
            type="button"
            class="theme-option"
            class:theme-option-selected={theme === option.id}
            aria-pressed={theme === option.id}
            onclick={() => chooseTheme(option.id)}
          >
            <span class={`theme-preview theme-preview-${option.id}`} aria-hidden="true">
              <span></span><span></span><span></span>
            </span>
            <span class="theme-option-copy">
              <strong>{option.label}</strong>
              <small>{option.description}</small>
            </span>
            {#if theme === option.id}
              <i class="ph ph-check-circle theme-option-check" aria-hidden="true"></i>
            {/if}
          </button>
        {/each}
      </div>
    </section>
  {/if}

  <button
    type="button"
    class="theme-toggle"
    onclick={() => settingsOpen = !settingsOpen}
    aria-expanded={settingsOpen}
    aria-haspopup="dialog"
    aria-label="Appearance settings"
    title="Appearance settings"
  >
    <i class="ph ph-gear-six" aria-hidden="true"></i>
    <span>Appearance</span>
  </button>
</div>
