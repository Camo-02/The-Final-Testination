<script lang="ts">
	import { BASE_API_URL } from '$src/constants';
	import { onMount } from 'svelte';
  
	let showPassword = false;
	let credential: string;
	let password: string;
  
	async function login() {
	  const res = await fetch(`${BASE_API_URL}/player/login`, {
		method: 'POST',
		credentials: 'include',
		headers: {
		  'Content-Type': 'application/json'
		},
		body: JSON.stringify({ credential, password })
	  });
	  const data = await res.json();
	  if (res.ok) {
		window.location.replace('/');
	  } else {
		alert(data.error);
	  }
	}
  
	async function handleCredentialResponse(response: any) {
	console.log(response);

	try {
		const res = await fetch(`${BASE_API_URL}/player/googleLogin`, {
		method: 'POST',
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ token: response.credential }) // Send the Google token to the backend
		});

		const data = await res.json();
		if (res.ok) {
		// Successful Google login, redirect the user
		window.location.replace('/');
		} else {
		// Error during login
		alert(data.error || 'Google login failed');
		}
	} catch (error) {
		console.error('Error during Google login:', error);
		alert('An error occurred during Google login.');
	}
}

  
	onMount(() => {
	// Initialize the Google Sign-In client with the client ID and callback function
	  google.accounts.id.initialize({
		client_id: 'Insert_Your_Client_ID', // Google Client ID
		callback: handleCredentialResponse
	  });
	   // Render the Google Sign-In button in the div with id 'google-signin-button'
	  google.accounts.id.renderButton(
		document.getElementById('google-signin-button'),
		{ theme: 'filled_blue', 
		  size: 'large',
		  shape: 'pill' 
		  } // Customize the button
	  );
	  google.accounts.id.prompt(); // Show the One Tap prompt if necessary
	});
  </script>
  
  <main
	class="relative min-h-screen bg-platform-white flex flex-col items-center justify-center w-full overflow-hidden"
>
	<img
		src="/purple_block.svg"
		alt="purple block"
		class="absolute bottom-0 left-0 w-1/2 h-auto rotate-180"
	/>
	<img src="/purple_block.svg" alt="purple block" class="absolute top-0 right-0 w-1/2 h-auto" />
	<div class="w-1/2 mx-auto flex flex-col items-center">
		<img src="/logo.svg" alt="The Final Testination logo" class="w-[25%] mb-[5%]" />
		<div id="google-signin-button" class="mb-4"></div> <!-- Google Sign-In button -->
		<input
			id="credential"
			type="email"
			bind:value={credential}
			placeholder="Email or Username"
			class="z-[1000] w-1/2 mx-auto rounded-full outline-none px-5 py-3 bg-platform-white placeholder:text-theme font-platform placeholder:text-opacity-50 text-theme font-bold border-4 border-platform-black mb-[1%] shadow-[0_35px_40px_-15px_rgba(0,0,0,0.5)] shadow-platform-black"
		/>
		<div class="flex relative w-1/2 mx-auto mb-[5%]">
			{#if showPassword}
				<input
					id="password"
					type="text"
					bind:value={password}
					placeholder="password"
					class="z-[1000] w-full rounded-full outline-none px-5 py-3 bg-platform-white placeholder:text-theme font-platform placeholder:text-opacity-50 font-bold border-4 border-platform-black text-theme shadow-[0_35px_60px_-15px_rgba(0,0,0,0.3)] shadow-platform-black"
				/>
				<button
					on:click={() => {
						showPassword = !showPassword;
					}}
				>
					<img
						src="/hide_password_button.svg"
						alt="eye"
						class=" z-[1001] absolute right-4 top-4 w-6 h-6"
					/>
				</button>
			{:else}
				<input
					id="password"
					type="password"
					bind:value={password}
					placeholder="Password"
					class="z-[1000] w-full rounded-full outline-none px-5 py-3 bg-platform-white placeholder:text-theme font-platform placeholder:text-opacity-50 font-bold border-4 border-platform-black text-theme shadow-[0_35px_60px_-15px_rgba(0,0,0,0.3)] shadow-platform-black"
				/>

				<button
					on:click={() => {
						showPassword = !showPassword;
					}}
				>
					<img
						src="/show_password_button.svg"
						alt="eye"
						class=" z-[1001] absolute right-4 top-4 w-6 h-6"
					/>
				</button>
			{/if}
		</div>
		<button
			id="loginBtn"
			class=" z-[1000] bg-primary hover:bg-theme text-platform-white font-platform font-bold py-2 px-4 rounded-full w-4/12 mb-[1%]"
			on:click={login}>Login</button
		>
		<h3 class=" z-[1000] text-platform-black font-platform font-regular">
			or <span class="font-bold underline"
				><a href="/register" class="text-platform-purple hover:text-platform-purple-dark"
					>Register</a
				></span
			>
		</h3>
	</div>
</main>
