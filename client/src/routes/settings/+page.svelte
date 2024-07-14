<script lang="ts">
  import { onMount } from "svelte";
  import { getCurrentWindow } from "@tauri-apps/api/window";
  import { invoke } from "@tauri-apps/api/core";
  import { listen } from "@tauri-apps/api/event";
  import LoaderCircle from "lucide-svelte/icons/loader-circle";

  import type { Settings } from "$lib/types/settings";
  import { DEFAULT_SETTINGS } from "$lib/constants";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";
  import { Switch } from "$lib/components/ui/switch";

  let isSaving = false;

  let settings: Settings = DEFAULT_SETTINGS;

  async function getCurrentWindowSettings(): Promise<void> {
    settings = (await invoke("get_settings")) as Settings;
    console.log("Current settings", settings);
  }

  async function hideWindow(): Promise<void> {
    console.log("Hiding window");
    await getCurrentWindow().hide();
  }

  async function showWindow(): Promise<void> {
    console.log("Showing window");
    await getCurrentWindow().show();
  }

  async function saveSettings(): Promise<void> {
    console.log("Saving settings", settings);
    isSaving = true;

    // Add a small delay to show the loader animation
    await new Promise((resolve) => setTimeout(resolve, 100));

    const result = await invoke<Settings>("update_settings", { settings });
    if (result) {
      console.log("Settings saved");
      settings = result;
      isSaving = false;
      await hideWindow();
    }
  }

  onMount(() => {
    // Handle the window being shown
    listen("show", async () => {
      await getCurrentWindowSettings();
      await showWindow();
    });
  });
</script>

<div class="container grid grid-cols-1 gap-2">
  <h1>Settings</h1>
  <section class="grid grid-cols-1 gap-2 mt-2 overflow-y-auto ps-1 pe-2">
    <div class="grid items-center grid-cols-2 gap-2">
      <Label for="autostart">Autostart</Label>
      <Switch id="autostart" bind:checked={settings.autostart} />
    </div>

    <div class="grid items-center grid-cols-2 gap-2">
      <Label for="log-level">Log Level</Label>
      <Input id="log-level" type="text" bind:value={settings.log_level} />
    </div>

    <div class="grid items-center grid-cols-2 gap-2">
      <Label for="server-host">Host</Label>
      <Input id="server-host" type="text" bind:value={settings.server.host} />
    </div>
    <div class="grid items-center grid-cols-2 gap-2">
      <Label for="server-port">Port</Label>
      <Input id="server-port" type="number" bind:value={settings.server.port} />
    </div>

    <div class="grid items-center grid-cols-2 gap-2">
      <Button class="bg-primary" on:click={hideWindow}>Cancel</Button>
      <Button disabled={isSaving} on:click={saveSettings}>
        {#if isSaving}
          <LoaderCircle class="w-4 h-4 mr-2 animate-spin" />
          Please wait...
        {:else}
          Save
        {/if}
      </Button>
    </div>
  </section>
</div>

<style lang="postcss">
  .container {
    height: 98vh;
    margin: 0;
    padding: 0.8rem;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    background-color: rgba(15, 23, 42, 0.7);
    border-radius: 1rem;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  }

  h1 {
    @apply text-2xl font-medium mt-2;
  }
</style>
