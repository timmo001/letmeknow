<script lang="ts">
  import { onMount } from "svelte";
  import { getCurrent } from "@tauri-apps/api/window";
  import { invoke } from "@tauri-apps/api/core";

  async function hideWindow(): Promise<void> {
    console.log("Hiding window");
    await getCurrent().hide();
  }

  async function showWindow(): Promise<void> {
    console.log("Showing window");
    await getCurrent().show();
    await invoke("set_window", {});
  }

  onMount(() => {
    // TODO: Hide the window on mount
    // hideWindow();
  });
</script>

<div class="container">
  <h1>Settings</h1>
  <h4>Manage settings</h4>

  <input
    {type}
    className={cn(
      "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
      className
    )}
    {ref}
    {...props}
  />
</div>

<style lang="postcss">
  .container {
    height: 98vh;
    margin: 0;
    padding: 0.8rem;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    background-color: rgba(15, 23, 42, 0.5);
    border-radius: 1rem;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  }

  h1 {
    @apply text-2xl font-medium mt-2;
  }

  h4 {
    @apply text-sm font-thin italic ms-2;
  }

  p {
    @apply text-base mt-1;
  }
</style>
