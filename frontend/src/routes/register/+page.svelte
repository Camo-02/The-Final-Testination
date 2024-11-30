<script lang="ts">
	import { BASE_API_URL } from '$src/constants';
	import { onMount } from 'svelte';

	let showPassword = false;
	let username: string;
	let email: string;
	let password: string;

	async function register() {
		const res = await fetch(`${BASE_API_URL}/player/register`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ username, password, email })
		});
		const data = await res.json();
		if (res.ok) {
			window.location.replace('/');
			console.log('====================================');
			console.log(username, password, email);
			console.log('====================================');
		} else {
			alert(data.error);
		}
	}

	onMount(async () => {
		const res = await fetch(`${BASE_API_URL}/player/availableLevels`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.ok) {
			window.location.replace('/');
		}
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
		<input
			id="user"
			type="email"
			bind:value={username}
			placeholder="username"
			class="z-[1000] w-1/2 mx-auto rounded-full outline-none px-5 py-3 bg-platform-white placeholder:text-theme font-platform placeholder:text-opacity-50 text-theme font-bold border-4 border-platform-black mb-[1%] shadow-[0_35px_40px_-15px_rgba(0,0,0,0.5)] shadow-platform-black"
		/>
		<input
			id="mail"
			type="email"
			bind:value={email}
			placeholder="email"
			class="z-[1000] w-1/2 mx-auto rounded-full outline-none px-5 py-3 bg-platform-white placeholder:text-theme font-platform placeholder:text-opacity-50 text-theme font-bold border-4 border-platform-black mb-[1%] shadow-[0_35px_40px_-15px_rgba(0,0,0,0.5)] shadow-platform-black"
		/>
		<div class="flex relative w-1/2 mx-auto mb-[5%]">
			<!-- <input type={showPassword ? "text" : "password"} bind:value={password} placeholder="password" class="z-[1000] w-full rounded-full outline-none px-5 py-3 bg-platform-white placeholder:text-theme font-platform placeholder:text-opacity-50 font-bold border-4 border-platform-black  shadow-lg text-theme shadow-[0_35px_60px_-15px_rgba(0,0,0,0.3)] shadow-platform-black"/>-->
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
					placeholder="password"
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
			id="registerBtn"
			class=" z-[1000] bg-primary hover:bg-theme text-platform-white font-platform font-bold py-2 px-4 rounded-full w-4/12 mb-[1%]"
			on:click={register}>Register</button
		>
		<h3 class=" z-[1000] text-platform-black font-platform font-regular">
			or <span class="font-bold underline"
				><a href="/login" class="text-platform-purple hover:text-platform-purple-dark"
					>Login</a
				>
	</div>
</main>
