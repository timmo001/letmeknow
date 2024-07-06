<script lang="ts">
  import { onMount } from "svelte";
  import { getCurrent, LogicalSize } from "@tauri-apps/api/window";

  import type { Register } from "../types/register";
  import type { Notification } from "../types/notification";

  const NOTIFICATION_CONNECTING: Notification = {
    title: "LetMeKnow Client Connecting",
    subtitle: "Connecting to server...",
    timeout: 5000,
  };

  const NOTIFICATION_CONNECTED: Notification = {
    title: "LetMeKnow Client Connected",
    subtitle: "Listening for notifications...",
    timeout: 5000,
  };

  let height: number = 60.0;
  let notification: Notification = NOTIFICATION_CONNECTING;

  let timeoutReconnect: NodeJS.Timeout | null = null;
  let timeoutHide: NodeJS.Timeout | null = null;

  async function hideWindow(): Promise<void> {
    console.log("Hiding window");
    await getCurrent().hide();
  }

  async function showWindow(): Promise<void> {
    console.log("Showing window");
    await getCurrent().show();
  }

  async function resize(
    height: number = 66,
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
      counter++;

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

      // 500ms * 20 = 10s
      if (counter > 20) clearInterval(interval);
    }, 500);
  }

  async function updateData(n: Notification): Promise<void> {
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

    // Set the default timeout if not provided
    if (!n.timeout) n.timeout = 10000;

    // Hide the window after the timeout
    if (timeoutHide) clearTimeout(timeoutHide);
    timeoutHide = setTimeout(() => {
      hideWindow();
    }, n.timeout);
  }

  async function reconnectToServer(): Promise<void> {
    // updateData({
    //   title: "LetMeKnow Client Disconnected",
    //   subtitle: "Reconnecting to server in 5 seconds...",
    //   timeout: 5000,
    // });

    // Reconnect to the server
    console.log("Reconnecting to the server in 5s");
    if (timeoutReconnect) clearTimeout(timeoutReconnect);
    timeoutReconnect = setTimeout(() => {
      setupServerConnection();
    }, 5000);
  }

  async function setupServerConnection(): Promise<void> {
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
      };

      socket.onopen = () => {
        console.log("Connection established");
        updateData(NOTIFICATION_CONNECTED);

        // Register the client with the server
        const register: Register = {
          type: "register",
          userID: "client", // TODO: Use a real user ID
        };
        socket.send(JSON.stringify(register));
      };

      socket.onmessage = (event) => {
        console.log("Message received:", event.data);

        let data;
        try {
          data = JSON.parse(event.data);
        } catch (error) {
          console.error("Error parsing JSON:", error);
        }

        console.log("Data type:", data.type);

        // Test if data is a successful registration
        if (data.type === "register" && data.success) {
          console.log("Registration successful");
          return;
        }

        // Test if data is a valid notification
        if (data.type === "notification") {
          const notification: Notification = data;
          updateData(notification);
        }
      };
    } catch (error) {
      console.error("Error:", error);
      reconnectToServer();
    }
  }

  onMount(() => {
    updateData(NOTIFICATION_CONNECTING);
    setupServerConnection()
      .then(() => {
        console.log("Server setup completed");
      })
      .catch((error) => {
        console.error("Caught error:", error);
        reconnectToServer();
      });
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
    @apply text-2xl font-medium mt-2;
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
