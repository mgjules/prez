<script lang="ts">
  import { inertia, Deferred } from "@inertiajs/svelte";
  import { router } from "@inertiajs/svelte";

  type User = {
    id: string;
    name: string;
    email: string;
    created_at: Date;
    updated_at: Date;
  };

  export let users: User[];

  // Modal state
  let isModalOpen = false;
  let isDeleteModalOpen = false;
  let editingUser: User | null = null;
  let userToDelete: User | null = null;

  // Form data
  let formData = {
    name: "",
    email: "",
  };

  // Date formatting
  function formatDate(dateValue: Date | string) {
    const date = new Date(dateValue);
    return date.toLocaleDateString() + " " + date.toLocaleTimeString();
  }

  // Modal functions
  function openCreateModal() {
    editingUser = null;
    formData = { name: "", email: "" };
    isModalOpen = true;
  }

  function openEditModal(user: User) {
    editingUser = user;
    formData = { name: user.name, email: user.email };
    isModalOpen = true;
  }

  function closeModal() {
    isModalOpen = false;
    editingUser = null;
    formData = { name: "", email: "" };
  }

  function openDeleteModal(user: User) {
    userToDelete = user;
    isDeleteModalOpen = true;
  }

  function closeDeleteModal() {
    isDeleteModalOpen = false;
    userToDelete = null;
  }

  // Form submission
  async function handleSubmit() {
    try {
      if (editingUser) {
        // Update existing user - send both name and email (email is readonly but needed for backend)
        router.patch("/users/", {
          name: formData.name,
          email: formData.email,
          id: editingUser.id,
        });
      } else {
        // Create new user - both name and email
        router.post("/users/", formData);
      }
      closeModal();
    } catch (error) {
      console.error("Error saving user:", error);
    }
  }

  // Delete confirmation
  async function confirmDelete() {
    if (userToDelete) {
      try {
        router.delete("/users/", {
          data: { id: userToDelete.id },
        });
        closeDeleteModal();
      } catch (error) {
        console.error("Error deleting user:", error);
      }
    }
  }
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
  <div class="mx-auto max-w-6xl py-8 sm:py-12 lg:py-16">
    <div class="px-4 sm:px-6 lg:px-8">
      <!-- Create User Button Section -->
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900">Users</h1>
        <button
          on:click={openCreateModal}
          class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
        >
          + Create User
        </button>
      </div>

      <!-- Users Table -->
      <div class="bg-white shadow-sm rounded-lg overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Name
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Email
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Created
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Updated
              </th>
              <th
                class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <Deferred data="users">
              <svelte:fragment slot="fallback">
                <tr>
                  <td colspan="5" class="px-6 py-4 text-center text-gray-500">
                    Loading...
                  </td>
                </tr>
              </svelte:fragment>

              {#each users ?? [] as user}
                <tr class="hover:bg-gray-50 transition-colors">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">
                      {user.name}
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">{user.email}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">
                      {formatDate(user.created_at)}
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-gray-500">
                      {formatDate(user.updated_at)}
                    </div>
                  </td>
                  <td
                    class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                  >
                    <button
                      on:click={() => openEditModal(user)}
                      class="text-indigo-600 hover:text-indigo-900 mr-3 cursor-pointer"
                    >
                      Edit
                    </button>
                    <button
                      on:click={() => openDeleteModal(user)}
                      class="text-red-600 hover:text-red-900 cursor-pointer"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              {/each}
            </Deferred>
          </tbody>
        </table>
      </div>

      <!-- Back Button -->
      <div class="mt-8 flex items-center justify-center">
        <a
          use:inertia
          href="/"
          class="text-sm font-semibold leading-6 text-gray-900 hover:text-gray-700"
          ><span aria-hidden="true">←</span> Back</a
        >
      </div>
    </div>
  </div>

  <!-- Create/Edit User Modal -->
  {#if isModalOpen}
    <div
      class="fixed inset-0 bg-gray-600 opacity-95 overflow-y-auto h-full w-full z-50"
    >
      <div
        class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white"
      >
        <div class="mt-3">
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            {editingUser ? "Edit User" : "Create User"}
          </h3>
          <form on:submit|preventDefault={handleSubmit}>
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700">Name</label
              >
              <input
                type="text"
                bind:value={formData.name}
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                required
              />
            </div>
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700"
                >Email</label
              >
              <input
                type="email"
                bind:value={formData.email}
                class="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 {editingUser !=
                null
                  ? 'bg-gray-100'
                  : ''}"
                required={editingUser == null}
                readonly={editingUser != null}
              />
              {#if editingUser != null}
                <p class="mt-1 text-xs text-gray-500">Email cannot be edited</p>
              {/if}
            </div>
            <div class="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                on:click={closeModal}
                class="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400 transition-colors cursor-pointer"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors cursor-pointer"
              >
                {editingUser ? "Update" : "Create"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  {/if}

  <!-- Delete Confirmation Modal -->
  {#if isDeleteModalOpen}
    <div
      class="fixed inset-0 bg-gray-600 opacity-95 overflow-y-auto h-full w-full z-50"
    >
      <div
        class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white"
      >
        <div class="mt-3 text-center">
          <div
            class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100"
          >
            <svg
              class="h-6 w-6 text-red-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z"
              />
            </svg>
          </div>
          <h3 class="text-lg leading-6 font-medium text-gray-900 mt-2">
            Delete User
          </h3>
          <div class="mt-2 px-7 py-3">
            <p class="text-sm text-gray-500">
              Are you sure you want to delete {userToDelete?.name}? This action
              cannot be undone.
            </p>
          </div>
          <div class="flex justify-center space-x-3 mt-4">
            <button
              on:click={closeDeleteModal}
              class="px-4 py-2 bg-gray-300 text-gray-700 rounded-md hover:bg-gray-400 transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              on:click={confirmDelete}
              class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors cursor-pointer"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  {/if}

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
