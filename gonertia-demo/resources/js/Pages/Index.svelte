<script lang="ts">
  import { inertia } from "@inertiajs/svelte";

  export let text: string;
</script>

<svelte:head>
  <title>Gonertia - Home</title>
</svelte:head>

<div class="relative isolate px-6 pt-14 lg:px-8">
  <div class="mx-auto max-w-4xl py-32 sm:py-48 lg:py-56">
    <div class="text-center">
      <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
        {text}
      </h1>
      <p class="mt-6 text-lg leading-8 text-gray-600">
        When you want that sweet SPA experience without the API brainrot.
      </p>
      <div class="mt-10 flex items-center justify-center gap-x-6">
        <a
          use:inertia
          href="/users"
          class="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >Get started</a
        >
      </div>
    </div>
  </div>
</div>
