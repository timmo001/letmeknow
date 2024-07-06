<script lang="ts">
  import { onMount } from "svelte";
  import { getCurrent, LogicalSize } from "@tauri-apps/api/window";

  import type { Notification } from "../types/notificaiton";

  let height: number = 60.0;
  let notification: Notification = {
    title: "Title",
    subtitle: "Subtitle",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec fermentum nunc.",
    image: { url: "https://via.placeholder.com/480x270" },
  };

  async function hideWindow(): Promise<void> {
    console.log("Hiding window");
    await getCurrent().hide();
  }

  async function showWindow(): Promise<void> {
    console.log("Showing window");
    await getCurrent().show();
  }

  async function resize(
    height: number = 60,
    width: number = 370
  ): Promise<void> {
    console.log("Resizing window:", { width, height });

    // Resize the window
    await getCurrent().setSize(new LogicalSize(width, height));
  }

  function updateWindowSize(): void {
    // Update the window size up to 10 times. Allows the window to resize properly with images loading.
    let counter = 0;
    const interval = setInterval(() => {
      // Get .container element height
      const container = document.querySelector(".container");
      if (!container) return;
      const heightComputed = window.getComputedStyle(container).height;
      // console.log("Container height:", heightComputed);

      // Remove "px" from the height value
      const heightValue = parseFloat(heightComputed.replace("px", ""));
      // console.log("Container height value:", heightValue);

      if (heightValue === height) return;

      // Update the height value
      height = heightValue;

      // Resize the window
      resize(heightValue);
      counter++;
      // 500ms * 20 = 10s
      if (counter > 20) clearInterval(interval);
    }, 500);
  }

  function updateData(n: Notification): void {
    // Update data
    console.log("Updating data:", n);

    notification = n;

    // Update the window title
    if (n.title) {
      document.title = n.title;
    }

    // Show the window
    showWindow();

    // Update the window size
    updateWindowSize();
  }

  function reconnectToServer(): void {
    // Reconnect to the server
    console.log("Reconnecting to the server in 5s");
    setTimeout(() => {
      setupServerConnection();
    }, 5000);
  }

  function setupServerConnection(): void {
    // Setup server connection
    console.log("Setting up server connection");

    try {
      const socket = new WebSocket("ws://localhost:8080/websocket");

      socket.onclose = () => {
        console.log("Connection closed");
        reconnectToServer();
      };

      socket.onerror = (error) => {
        console.error("Error:", error);
        reconnectToServer();
      };

      socket.onopen = () => {
        console.log("Connection established");
      };

      socket.onmessage = (event) => {
        console.log("Message received:", event.data);
      };
    } catch (error) {
      console.error("Error:", error);
      reconnectToServer();
    }
  }

  onMount(() => {
    updateData(notification);
    // hideWindow();
    // showWindow();
    setupServerConnection();
  });
</script>

<div class="container">
  {#if notification.image}
    <img src={notification.image.url} alt={notification.title} />
  {/if}
  {#if notification.title}
    <h1>{notification.title}</h1>
  {/if}
  {#if notification.subtitle}
    <h4>{notification.subtitle}</h4>
  {/if}
  {#if notification.content}
    <p>{notification.content}</p>
  {/if}
</div>

<style lang="postcss">
  .container {
    margin: 0;
    padding: 0.8rem;
    display: flex;
    flex-direction: column;
    justify-content: center;
    background-color: rgba(15, 23, 42, 0.6);
    border-radius: 1rem;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  }

  h1 {
    @apply text-3xl font-medium mt-2;
  }

  h4 {
    @apply text-sm font-thin italic ms-2;
  }

  p {
    @apply text-base mt-1;
  }

  img {
    width: 100%;
    height: auto;
    border-radius: 0.5rem;
  }
</style>
