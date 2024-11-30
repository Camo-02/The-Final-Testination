<script lang="ts">
	import NavBar from '$components/NavBar.svelte';
	import { BASE_API_URL } from '$src/constants';
	import SplashScreen from '$components/SplashScreen.svelte';
	import { onMount } from 'svelte';
	import { split } from 'postcss/lib/list';

	type ProfileData = {
		profileImage: string;
		username: string;
		email: string;
		levels:
			| {
					title: string;
					score: number | null;
					time_freeze_points_used: number | null;
					textual_hint_points_used: number | null;
					hint_solution_points_used: number | null;
					start_time: string;
					end_time: string | null;
			  }[]
			| null;
	};

	type AvailableIcon = {
		id: string;
		svg: string;
	};
	let checkStatus = true;
	let total_coins: number = 0;
	let data: ProfileData;
	let loading = true;
	let show_change_icon = false;
	let selected_icon = '1';
	let availabe_icons: AvailableIcon[] = [];
	let icons_svg_elments: HTMLDivElement[] = [];
	let profile_icon_element: HTMLImageElement;
	async function selectIcon() {
		const res = await fetch(`${BASE_API_URL}/player/availableIcons`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.status === 401) {
			window.location.replace('/');
			return;
		}

		icons_svg_elments = [];

		availabe_icons = await res.json();
		availabe_icons.forEach(async (icon) => {
			let icon_svg = document.createElement('div');
			icon_svg.style.backgroundImage = `url("data:image/svg+xml,${encodeURIComponent(icon.svg)}")`;
			icon_svg.style.width = '10rem';
			icon_svg.style.height = '10rem';
			icon_svg.style.backgroundSize = 'contain';
			icon_svg.style.backgroundRepeat = 'no-repeat';
			icon_svg.style.cursor = 'pointer';
			icon_svg.style.margin = '0.5rem';
			icon_svg.id = icon.id;
			icons_svg_elments = [...icons_svg_elments, icon_svg];
		});
	}

	function getNewProfileImg(svg_string: string) {
		let blob = new Blob([svg_string], { type: 'image/svg+xml' });
		let url = URL.createObjectURL(blob);
		let imgElem = document.createElement('img');
		imgElem.src = url;
		imgElem.style.objectFit = 'contain';
		imgElem.id = 'profileIcon';
		return imgElem;
	}

	async function changeIcon() {
		const res = await fetch(`${BASE_API_URL}/player/changeIcon`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ icon: selected_icon })
		});

		if (res.status === 401) {
			window.location.replace('/');
			return;
		}
		let new_icon = availabe_icons.find((icon) => icon.id === selected_icon);
		profile_icon_element = getNewProfileImg(new_icon?.svg ?? '');
		show_change_icon = false;
	}

	async function getProfile() {
		const res = await fetch(`${BASE_API_URL}/player/profile`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (res.status === 401) {
			window.location.replace('/');
			checkStatus = false;
			return;
		}
		data = await res.json();
		total_coins = (data.levels ?? []).reduce((acc, level) => {
			let points_spent_on_hints = 0;
			points_spent_on_hints += level.time_freeze_points_used ?? 0;
			points_spent_on_hints += level.textual_hint_points_used ?? 0;
			points_spent_on_hints += level.hint_solution_points_used ?? 0;
			return acc + (level.score ?? 0) - points_spent_on_hints;
		}, 0);
	}

	onMount(async () => {
		loading = true;
		await getProfile();
		if (checkStatus) {
			profile_icon_element = getNewProfileImg(data.profileImage);
			loading = false;
		}
	});
	async function logout() {
		const res = await fetch(`${BASE_API_URL}/player/logout`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		});
		const data = await res.json();
		if (res.ok) {
			window.location.reload();
		} else {
			alert(data.error);
		}
	}
</script>

{#if loading}
	<SplashScreen />
{:else}
	<main class="min-h-screen bg-platform-white font-platform relative">
		<NavBar />
		{#if show_change_icon}
			<div
				class="absolute w-full h-full bg-platform-black bg-opacity-75 top-0 z-[900] flex flex-col items-center justify-center"
			>
				<div
					class="bg-platform-white flex flex-col items-center p-5 rounded-[70px] border-8 border-platform-black min-h-[80vh] max-h-[80vh] w-5/12 relative"
				>
					<button
						on:click={() => {
							show_change_icon = false;
							icons_svg_elments = [];
						}}
						class="absolute top-5 right-5 p-2 m-2"
					>
						<img src="/cross.svg" class="w-6 h-6" alt="close" />
					</button>
					<h1 class="text-3xl font-bold text-platform-black mt-14 text-primary">
						Change Profile Icon
					</h1>
					{#if availabe_icons.length === 0}
						<div class="w-11/12 h-[75%] mx-auto flex items-center justify-center">
							<img src="/logoanimated.svg" class="w-40 h-40" alt="loading" />
						</div>
					{:else}
						<div
							id="availableIcons"
							class="w-11/12 h-[75%] mx-auto mt-14 overflow-y-scroll grid grid-cols-2 p-3 place-items-center"
						>
							{#each icons_svg_elments as icon}
								<button
									id="icon{icon.id}"
									on:click={() => {
										icon.style.borderRadius = '100%';
										icon.style.border = '10px solid #E41BF6';
										selected_icon = icon.id;
									}}
									on:focusout={() => {
										icon.style.border = 'none';
									}}
								>
									{@html icon.outerHTML}
								</button>
							{/each}
						</div>
						<button
							id="saveIconBtn"
							on:click={() => {
								changeIcon();
							}}
							class="w-8/12 py-4 bg-primary text-platform-white font-bold text-2xl rounded-full mt-14 mb-10"
							>Save</button
						>
					{/if}
				</div>
			</div>
		{/if}
		<div class="flex flex-row pt-20">
			<div class="flex flex-col items-center justify-center w-3/12 min-h-[70vh]">
				<div class="relative">
					<div id="icon" class="w-52 h-52">
						{@html profile_icon_element.outerHTML}
					</div>
					<button
						on:click={async () => {
							show_change_icon = true;
							await selectIcon();
						}}
						id="selectIconButton"
						class="z-[1000] top-0 absolute w-full bg-platform-white p-auto h-full rounded-full opacity-0 hover:opacity-100 flex justify-center items-center bg-opacity-0 hover:bg-opacity-50 transform ease-in-out duration-300"
					>
						<img src="/modify_icon.svg" class="w-14 h-14" alt="edit profile logo" />
					</button>
				</div>
				<h1 class="text-5xl font-bold mt-5 text-platform-black" id="username">{data.username}</h1>
				<h1 class="text-3xl font-bold mt-5 text-platform-black" id="coins">
					{total_coins} Coins
				</h1>
				<!--LOG OUT button -->
				<button
					on:click={() => {
						//delete the cookie testination_login
						window.location.replace('/login');
						logout();
						//logout after redirecting to login page to avoid the error on firefox browser
					}}
					class="w-8/12 py-4 bg-primary text-platform-white font-bold text-2xl rounded-full mt-14"
					>Log Out</button
				>
			</div>
			<div
				class="items-center w-8/12 border-8 border-platform-black rounded-[50px] h-[70vh] flex flex-col shadow-2xl bg-platform-white overflow-y-scroll"
				class:justify-center={data.levels === null || data.levels.length === 0}
				class:justify-between={data.levels !== null && data.levels.length > 0}
				id="levels-container"
			>
				{#if data.levels !== null && data.levels.length > 0}
					{#each data.levels as level, idx}
						<div
							class="w-full p-5 flex flex-row my-14 hover:bg-gray-400 hover:bg-opacity-25 transform ease-in-out duration-300"
						>
							<div class="flex items-center justify-center w-5/12 flex-col">
								<h1 class="text-3xl font-bold text-platform-black mb-5" id="level-title-{idx}">
									{level.title}
								</h1>
								<h1 class="text-3xl font-bold text-platform-black" id="level-score-{idx}">
									{level.end_time === null || (level.score ?? 0) === 0 ? 'Incomplete' : level.score}
								</h1>
							</div>
							<div class="flex items-center justify-center w-7/12">
								<div class="flex flex-col items-center justify-center">
									<div class="flex flex-row mb-5 items-center justify-between">
										<h1 class="text-xl font-bold text-platform-black" id="level-started-{idx}">
											{#if level.end_time === null}
												Started at
											{/if}

											{new Date(level.start_time).toLocaleString()}
										</h1>
										{#if level.end_time !== null}
											<img src="/horizontal_arrow_profile.svg" class="w-10 h-10 mx-2" alt="arrow" />
											<h1 class="text-xl font-bold text-platform-black" id="level-end-{idx}">
												{new Date(level.end_time).toLocaleString()}
											</h1>
										{/if}
									</div>
									<div class="flex flex-row w-full justify-between">
										{#if level.time_freeze_points_used !== null}
											<div class="flex flex-row items-center justify-between">
												<img src="/time_freeze_icon.svg" class="w-10 h-10 mr-2" alt="time freeze" />
												<h1 class="text-xl font-bold text-platform-black" id="level-freeze-{idx}">
													{level.time_freeze_points_used}
												</h1>
											</div>
										{/if}
										{#if level.textual_hint_points_used !== null}
											<div class="flex flex-row items-center justify-between">
												<img src="/text_hint_icon.svg" class="w-10 h-10 mr-2" alt="textual hint" />
												<h1 class="text-xl font-bold text-platform-black" id="level-textual-{idx}">
													{level.textual_hint_points_used}
												</h1>
											</div>
										{/if}
										{#if level.hint_solution_points_used !== null}
											<div class="flex flex-row items-center justify-between">
												<img
													src="/fill_block_icon.svg"
													class="w-10 h-10 mr-2"
													alt="hint solution"
												/>
												<h1 class="text-xl font-bold text-platform-black" id="level-solution-{idx}">
													{level.hint_solution_points_used}
												</h1>
											</div>
										{/if}
									</div>
								</div>
							</div>
						</div>
					{/each}
				{:else}
					<h2 class="text-3xl font-bold text-platform-black" id="no-levels-message">
						No levels completed yet
					</h2>
				{/if}
			</div>
		</div>
	</main>
{/if}
