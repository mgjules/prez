<script lang="ts">
  import { inertia, page } from "@inertiajs/svelte";

  // Navigation links
  let navLinks = [
    { href: "/", label: "Home", active: false },
    { href: "/users", label: "Users", active: false },
  ];

  $: {
    navLinks.forEach((link) => {
      link.active = link.href === $page.url;
    });
    navLinks = navLinks;
  }
</script>

<svelte:head>
  <title>Gonertia Demo</title>
  <meta name="description" content="Gonertia demo application with Svelte" />
</svelte:head>

<div class="min-h-screen bg-gray-100">
  <!-- Navigation Header -->
  <nav class="z-50 bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex">
          <!-- Logo -->
          <div class="flex-shrink-0 flex items-center">
            <a
              use:inertia
              href="/"
              class="text-xl font-bold text-indigo-600 hover:text-indigo-700"
            >
              Gonertia
            </a>
          </div>

          <!-- Navigation Links -->
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            {#each navLinks as link (link.href)}
              <a
                use:inertia
                href={link.href}
                class="inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium {link.active
                  ? 'border-indigo-500 text-gray-900'
                  : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'}"
              >
                {link.label}
              </a>
            {/each}
          </div>
        </div>

        <!-- Right side (could add user menu, etc.) -->
        <div class="flex items-center">
          <div class="text-sm text-gray-500">Demo Application</div>
        </div>
      </div>
    </div>
  </nav>

  <!-- Page Content -->
  <main class="py-6">
    <slot />
  </main>

  <!-- Footer -->
  <footer class="bg-white border-t border-gray-200">
    <div class="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8">
      <div class="text-center text-sm text-gray-500">
        Built with Go + Gonertia + Svelte
      </div>
    </div>
  </footer>
</div>

<style>
  /* Custom styles for layout */
  .border-b-2 {
    border-bottom-width: 2px;
  }
</style>

