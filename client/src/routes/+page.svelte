<script lang="ts">
  import { onMount } from "svelte";
  import { getCurrent, LogicalSize } from "@tauri-apps/api/window";

  let height = 60;

  async function resize(height: number = 60, width: number = 370) {
    console.log("Resizing window:", { width, height });

    // Resize the window
    await getCurrent().setSize(new LogicalSize(width, height));
  }

  function updateWindowSize() {
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

  onMount(() => {
    updateWindowSize();
  });
</script>

<div class="container">
  <img src="https://via.placeholder.com/480x270" alt="Placeholder" />
  <h1>Title</h1>
  <h4>Subtitle</h4>
  <p>
    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam nec
    fermentum nunc.
  </p>
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
