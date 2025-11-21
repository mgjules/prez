<script lang="ts">
  import { inertia, Deferred } from "@inertiajs/svelte";

  type User = {
    id: string;
    name: string;
    email: string;
    created_at: Date;
    updated_at: Date;
  };

  export let users: User[];
</script>

<svelte:head>
  <title>Gonertia - Users</title>
</svelte:head>

<div class="relative isolate px-6 pt-14 lg:px-8">
  <div
    class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
    aria-hidden="true"
  >
    <div
      class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"
      style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
    />
  </div>
  <div class="mx-auto max-w-4xl py-32 sm:py-48 lg:py-56">
    <div class="text-center">
      <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
        Users
      </h1>
      <div class="mt-10 space-y-4">
        <Deferred data="users">
          <svelte:fragment slot="fallback">
            <div>Loading...</div>
          </svelte:fragment>

          {#each users as user}
            <div>
              <pre>ID: {user.id}</pre>
              <pre>Name: {user.name}</pre>
              <pre>Email: {user.email}</pre>
              <pre>Created at: {user.created_at}</pre>
              <pre>Updated at: {user.updated_at}</pre>
            </div>
          {/each}
        </Deferred>
      </div>
      <div class="mt-10 flex items-center justify-center gap-x-6">
        <a
          use:inertia
          href="/"
          class="text-sm font-semibold leading-6 text-gray-900"
          ><span aria-hidden="true">←</span> Back</a
        >
      </div>
    </div>
  </div>
  <div
    class="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
    aria-hidden="true"
  >
    <div
      class="relative left-[calc(50%+3rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%+36rem)] sm:w-[72.1875rem]"
      style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
    />
  </div>
</div>
