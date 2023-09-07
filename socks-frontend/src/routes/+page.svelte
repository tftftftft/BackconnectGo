<script lang="ts">
  import {Toast ,getToastStore } from '@skeletonlabs/skeleton';
  import type { ToastSettings, ToastStore } from '@skeletonlabs/skeleton';
  import { goto } from '$app/navigation';

  let username: string = "";
  let password: string = "";
  const toast: ToastStore = getToastStore();

  async function checkCredentials() {
  try {
    console.log('Checking credentials...');
    const response = await fetch('http://127.0.0.1:30000/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify({ username, password })
    });

    const data = await response.json();

// Print out data for debugging
  console.log("Received data from /api/login:", data);

    if (data.status === "success") {
      // Store session data (You may also want to encrypt the data)
      localStorage.setItem('session', JSON.stringify(data.session));

      // Redirect to dashboard
      goto('/dashboard');
    } else {
      const toastSettings: ToastSettings = {
        message: "Credentials are invalid.",
        background: 'bg-error-600 text-white',
        classes: ''
      };
      toast.trigger(toastSettings);
    }
  } catch (error) {
    console.error('There was an error!', error);
  }
}
</script>

<div class="flex flex-col items-center justify-center min-h-screen mt-[-100px]">
  <form class="w-full max-w-sm">
    <div class="mb-4">
      <label for="username" class="block text-sm font-bold mb-2">Username:</label>
      <input bind:value={username} class="input" title="Input (text)" type="text" placeholder="Username"/>
    </div>
    <div class="mb-4">
      <label for="password" class="block text-sm font-bold mb-2">Password:</label>
      <input bind:value={password} class="input" title="Input (text)" type="password" placeholder="Password"/>
    </div>
    <div class="flex items-center justify-center">
      <button type="button" class="btn variant-filled" on:click={checkCredentials}>Login</button>
    </div>
  </form>
  <!-- Toast container -->
  <Toast />
</div>
