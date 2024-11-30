<script lang="ts">
	import { BASE_API_URL, PAGE_SIZE } from '$src/constants';
	import NavBar from '$components/NavBar.svelte';
	import { onMount } from 'svelte';
	import SplashScreen from '$components/SplashScreen.svelte';

	type LeaderboardEntry = {
		username: string;
		score: number;
	};

	let data: {
		currentPage: number;
		entries: LeaderboardEntry[];
		pages: number;
	} = { currentPage: 1, entries: [], pages: 1 };
	let fetching = true;
	let currentPage = data.currentPage;
	let entries = data.entries;
	let maxPages = data.pages;

	async function getLeaderboardPage(page: number) {
		const res = await fetch(`${BASE_API_URL}/leaderboard/${page}`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!res.ok) {
			console.error(
				'There was an error fetching the leaderboard data:' + res.status + ' ' + res.statusText
			);
			return;
		}

		const data = await res.json();
		entries = data.entries;
		currentPage = data.currentPage;

		// This rarely comes up, but if a player enters the leaderboard and the leaderboard gets another page,
		// after the first load, with this page we will be able to retrieve it
		maxPages = data.pages;
		fetching = false;
	}

	$: startIndex = (currentPage - 1) * PAGE_SIZE;

	onMount(() => {
		getLeaderboardPage(1);
	});
</script>

{#if fetching}
	<SplashScreen />
{:else}
	<main class="min-h-screen bg-platform-white">
		<NavBar />
		<div class="flex flex-col items-center justify-content w-full">
			<h1 class="text-theme font-platform text-5xl text-center font-black mb-10">Leaderboard</h1>
			<div class="flex items-center">
				<button
					id="go_back_button"
					disabled={currentPage === 1}
					on:click={() => getLeaderboardPage(currentPage - 1)}
				>
					<img
						src="/leftarrow.svg"
						class="w-8 h-8 cursor-pointer transform hover:-translate-x-2 ease-in-out duration-200"
						alt="Go back button"
					/>
				</button>
				<p class="text-theme font-platform mx-8 mt-5 mb-10 text-4xl font-black translate-y-[20%]">
					{currentPage}
				</p>
				<button
					id="go_forward_button"
					disabled={currentPage === maxPages}
					on:click={() => getLeaderboardPage(currentPage + 1)}
				>
					<img
						src="/rightarrow.svg"
						class="w-8 h-8 cursor-pointer transform hover:translate-x-2 ease-in-out duration-200"
						alt="Go forward button"
					/>
				</button>
			</div>
			<div class="w-full" id="contestants_list">
				{#each entries as entry, index}
					<div
						class="grid grid-flow-row grid-cols-3 gap-4 hover:bg-theme hover:bg-opacity-20 transform ease-in-out duration-300 my-5 py-3 hover:font-bold"
					>
						<div
							class="flex flex-col items-center justify-content text-theme text-3xl font-platform w-full"
						>
							{index + startIndex + 1}
						</div>
						<div
							class="flex flex-col items-center justify-content text-theme text-3xl font-platform w-full"
						>
							{entry.username}
						</div>
						<div
							class="flex flex-col items-center justify-content text-theme text-3xl font-platform w-full"
						>
							{entry.score}
						</div>
					</div>
				{/each}
			</div>
		</div>
	</main>
{/if}
