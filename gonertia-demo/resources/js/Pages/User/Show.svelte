<script lang="ts">
  import { inertia } from "@inertiajs/svelte";

  type User = {
    id: string;
    name: string;
    email: string;
    created_at: Date;
    updated_at: Date;
  };

  export let user: User;

  function formatDate(dateValue: Date | string) {
    const date = new Date(dateValue);
    return date.toLocaleDateString() + " " + date.toLocaleTimeString();
  }
</script>

<svelte:head>
  <title>Gonertia - User Profile</title>
</svelte:head>

<div class="relative isolate px-6 py-8 lg:px-8">
  <div class="mx-auto max-w-6xl">
    <div class="px-4 sm:px-6 lg:px-8">
      <!-- User Profile Header -->
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900">
          User Profile
        </h1>
      </div>

      <!-- User Information Card -->
      <div class="bg-white shadow-lg rounded-lg overflow-hidden">
        <div class="px-6 py-8">
          <div class="flex flex-col gap-6 md:flex-row">
            <!-- Avatar Section -->
            <div class="flex flex-2 justify-center items-center">
              <div
                class="w-32 h-32 rounded-full bg-indigo-100 flex items-center justify-center"
              >
                <div
                  class="w-24 h-24 rounded-full bg-indigo-600 flex items-center justify-center text-white text-3xl font-bold"
                >
                  {user.name.charAt(0).toUpperCase()}
                </div>
              </div>
            </div>

            <!-- User Details -->
            <div class="flex-3 space-y-6">
              <div>
                <h3 class="text-lg font-medium text-gray-900">
                  Contact Information
                </h3>
                <dl class="mt-2 space-y-1">
                  <div
                    class="flex justify-between py-3 border-b border-gray-200"
                  >
                    <dt class="text-sm font-medium text-gray-500">Full Name</dt>
                    <dd class="ml-4 text-sm text-gray-900">{user.name}</dd>
                  </div>
                  <div
                    class="flex justify-between py-3 border-b border-gray-200"
                  >
                    <dt class="text-sm font-medium text-gray-500">
                      Email Address
                    </dt>
                    <dd class="ml-4 text-sm text-gray-900">{user.email}</dd>
                  </div>
                </dl>
              </div>

              <div>
                <h3 class="text-lg font-medium text-gray-900">
                  Account Information
                </h3>
                <dl class="mt-2 space-y-1">
                  <div
                    class="flex justify-between py-3 border-b border-gray-200"
                  >
                    <dt class="text-sm font-medium text-gray-500">User ID</dt>
                    <dd class="ml-4 text-sm text-gray-900 font-mono">
                      {user.id}
                    </dd>
                  </div>
                  <div
                    class="flex justify-between py-3 border-b border-gray-200"
                  >
                    <dt class="text-sm font-medium text-gray-500">
                      Member Since
                    </dt>
                    <dd class="ml-4 text-sm text-gray-900">
                      {formatDate(user.created_at)}
                    </dd>
                  </div>
                  <div class="flex justify-between py-3">
                    <dt class="text-sm font-medium text-gray-500">
                      Last Updated
                    </dt>
                    <dd class="ml-4 text-sm text-gray-900">
                      {formatDate(user.updated_at)}
                    </dd>
                  </div>
                </dl>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Back Button -->
      <div class="mt-8 flex items-center justify-center">
        <a
          use:inertia
          href="/users"
          class="text-sm font-semibold leading-6 text-gray-900 hover:text-gray-700"
        >
          <span aria-hidden="true">←</span> Back to Users
        </a>
      </div>
    </div>
  </div>
</div>

