<script lang="ts">
	import { onMount } from 'svelte';
	import { BASE_API_URL } from '$src/constants';
	import NavBar from '$components/NavBar.svelte';

	interface Level {
		game_id: string;
		title: string;
		game_order: number;
		max_score: number;
		score_achieved: number;
		description: string;
	}

	let available_levels: Level[] = [];
	let fetching_levels: boolean = true;
	let levels: number[];
	let level_rows: number[][];
	let show_level_modal: boolean = false;
	let current_level_selected: number = 0;
	let available_level_class =
		'z-[1000] bg-theme rounded-full max-h-20 min-h-20 w-20 h-20 max-w-20 min-w-20 flex items-center justify-center';
	let unavailable_level_class =
		'z-[1000] bg-gray-500 rounded-full max-h-20 min-h-20 w-20 h-20 max-w-20 min-w-20 flex items-center justify-center';
	let left_available_level_class =
		'ml-auto z-[1000] bg-theme rounded-full max-h-20 min-h-20 w-20 h-20 max-w-20 min-w-20 flex items-center justify-center';
	let right_available_level_class =
		'mr-auto z-[1000] bg-theme rounded-full max-h-20 min-h-20 w-20 h-20 max-w-20 min-w-20 flex items-center justify-center';
	let left_unavailable_level_class =
		'ml-auto z-[1000] bg-gray-500 rounded-full max-h-20 min-h-20 w-20 h-20 max-w-20 min-w-20 flex items-center justify-center';
	let right_unavailable_level_class =
		'mr-auto z-[1000] bg-gray-500 rounded-full max-h-20 min-h-20 w-20 h-20 max-w-20 min-w-20 flex items-center justify-center';

	onMount(async () => {
		const res = await fetch(`${BASE_API_URL}/player/availableLevels`, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!res.ok) {
			window.location.replace('/login');
		}

		available_levels = await res.json();

		levels = Array.from({ length: available_levels.length }, (_, i) => i + 1);

		level_rows = Array.from({ length: Math.ceil(levels.length / 5) }).map((_, i) =>
			levels.slice(i * 5, i * 5 + 5)
		);

		level_rows.forEach((row, index) => {
			if (index % 2 != 0) {
				row = row.reverse();
			}
		});

		fetching_levels = false;
	});
</script>

{#if fetching_levels}
	<main class="min-h-screen bg-theme"></main>
{:else}
	<main class="min-h-screen bg-platform-white relative">
		{#if show_level_modal}
			<div
				class="absolute min-h-full min-w-full flex items-center justify-center bg-platform-black bg-opacity-50 z-[3000]"
			>
				<div
					id="modal"
					class="bg-platform-white border-8 border-platform-black rounded-[50px] p-10 flex items-center justify-center flex-col min-h-[70vh] w-5/12 shadow-[0_35px_60px_-15px_rgba(0,0,0,0.3)] shadow-platform-black relative"
				>
					<button
						id="close-modal-btn"
						class="absolute top-10 right-10 cursor-pointer w-5 h-5"
						on:click={() => {
							show_level_modal = false;
						}}
					>
						<img src="/cross.svg" alt="close modal" />
					</button>
					{#if current_level_selected - 1 != levels.length - 1}
						<h1
							class:bg-gray-500={available_levels[current_level_selected - 1].max_score < 0}
							class:bg-theme={available_levels[current_level_selected - 1].max_score >= 0}
							class=" z-[1000] rounded-full w-28 h-28 flex items-center justify-center text-platform-white font-platform font-black text-5xl"
						>
							{current_level_selected}
						</h1>
					{:else}
						<div
							class:bg-gray-500={available_levels[current_level_selected - 1].max_score < 0}
							class:bg-theme={available_levels[current_level_selected - 1].max_score >= 0}
							class="rounded-full w-28 h-28 flex items-center justify-center"
						>
							<img src="/skull.svg" alt="lock" class="w-10 h-10" />
						</div>
					{/if}

					<!-- Do not show specifics of level if it is locked -->
					<h1 id="level-title" class="font-platform font-black text-3xl text-platform-black my-5">
						{#if available_levels[current_level_selected - 1].max_score >= 0}
							{available_levels[current_level_selected - 1].title}
						{:else}
							Level {current_level_selected}
						{/if}
					</h1>
					<p
						id="level-description"
						class="font-platform font-medium text-xl text-platform-black mb-5 text-center"
					>
						{#if available_levels[current_level_selected - 1].max_score >= 0}
							{available_levels[current_level_selected - 1].description}
							<br />
							Score Achieved: {available_levels[current_level_selected - 1].score_achieved} out of {available_levels[
								current_level_selected - 1
							].max_score}
						{:else}
							Locked
						{/if}
					</p>

					{#if available_levels[current_level_selected - 1].score_achieved > 0}
						<p
							id="already-completed-text"
							class="underline font-platform font-medium text-xl text-platform-black mb-5 text-center"
						>
							You have already completed this level, playing it again will yield no points.
						</p>
					{/if}
					{#if available_levels[current_level_selected - 1].max_score >= 0}
						<a
							id="play-game-btn"
							href={'/game/' + available_levels[current_level_selected - 1].game_id}
							class="flex justify-center items-center z-[3000] bg-primary hover:bg-theme text-platform-white font-platform font-bold py-2 px-4 rounded-full w-4/12 mb-[1%]"
						>
							Play
						</a>
					{/if}
				</div>
			</div>
		{/if}

		<NavBar />
		<div class="flex items-center justify-around">
			<div
				class="w-3/12 border-8 border-platform-black rounded-[50px] min-h-[70vh] flex items-center justify-center flex-col shadow-[0_35px_60px_-15px_rgba(0,0,0,0.3)] shadow-platform-black"
			>
				<h1
					class="px-3 text-theme font-platform font-black text-3xl bg-theme bg-opacity-20 transform ease-in-out delay-100 duration-300 w-full py-5 text-center mb-3"
				>
					James Bond - GoldTest
				</h1>
				<h1
					class="px-3 text-theme font-platform font-regular text-3xl hover:bg-theme hover:bg-opacity-20 transform ease-in-out delay-100 duration-300 w-full py-5 text-center mb-3"
				>
					Indiana Jones - the Crystal Test
				</h1>
				<h1
					class="px-3 text-theme font-platform font-regular text-3xl hover:bg-theme hover:bg-opacity-20 transform ease-in-out delay-100 duration-300 w-full py-5 text-center"
				>
					Test Attack
				</h1>
			</div>

			<div
				class="w-8/12 border-8 border-platform-black rounded-[50px] flex flex-col py-10 items-center min-h-[70vh] justify-center shadow-[0_35px_60px_-15px_rgba(0,0,0,0.3)] shadow-platform-black"
				id="gameMap"
			>
				{#each level_rows as row, row_index}
					<div class="relative flex flex-row justify-between items-center mx-auto w-10/12">
						{#if row.length >= 2}
							<img
								src="/dotted_line_horizontal.svg"
								alt="dotted line"
								class="absolute top-1/2 h-auto translate-y-[-50%]"
							/>
						{/if}

						{#each row as level, level_index_in_row}
							{#if row.length > 1}
								{#if available_levels.some((avilable_level) => avilable_level.game_order === level)}
									<div class="relative flex">
										<button
											class={available_levels[level - 1].max_score >= 0
												? available_level_class
												: unavailable_level_class}
											on:click={() => {
												current_level_selected = level;
												show_level_modal = true;
											}}
										>
											{#if available_levels[level - 1].max_score < 0}
												<img src="/Lock.svg" alt="lock" class="w-8 h-8" />
											{:else if row_index % 2 != 0 && level_index_in_row == 0 && row_index == level_rows.length - 1}
												<img src="/skull.svg" alt="lock" class="w-8 h-8" />
											{:else if row_index % 2 == 0 && level_index_in_row == row.length - 1 && row_index == level_rows.length - 1}
												<img src="/skull.svg" alt="lock" class="w-8 h-8" />
											{:else}
												<h1
													class="text-center text-platform-white font-platform font-black text-2xl"
												>
													{level}
												</h1>
											{/if}
										</button>
										{#if available_levels[level - 1].score_achieved > 0}
											<img
												src="/Check.svg"
												alt="check"
												class="w-8 h-8 absolute top-0 right-0 z-[2000]"
											/>
										{/if}
									</div>
								{:else}
									<button
										disabled
										class={available_levels[level - 1].max_score >= 0
											? available_level_class
											: unavailable_level_class}
									>
										<img src="/Lock.svg" alt="lock" class="w-8 h-8" /> as previous case-->
										{#if row_index % 2 != 0 && level_index_in_row == 0 && row_index == level_rows.length - 1}
											<img src="/skull.svg" alt="lock" class="w-8 h-8" />
										{:else if row_index % 2 == 0 && level_index_in_row == row.length - 1 && row_index == level_rows.length - 1}
											<img src="/skull.svg" alt="lock" class="w-8 h-8" />
										{:else}
											<h1 class="text-center text-platform-white font-platform font-black text-2xl">
												{level}
											</h1>
										{/if}
									</button>
								{/if}
							{/if}

							{#if row.length == 1 && row_index % 2 != 0}
								<button
									class={available_levels[level - 1].max_score != -1
										? left_available_level_class
										: left_unavailable_level_class}
									on:click={() => {
										current_level_selected = level;
										show_level_modal = true;
									}}
								>
									{#if available_levels[level - 1].max_score != -1}
										<img src="/skull.svg" alt="lock" class="w-8 h-8" />
									{:else}
										<img src="/Lock.svg" alt="lock" class="w-8 h-8" />
									{/if}
								</button>
							{/if}

							{#if row.length == 1 && row_index % 2 == 0}
								<button
									class={available_levels[level - 1].max_score != -1
										? right_available_level_class
										: right_unavailable_level_class}
									on:click={() => {
										current_level_selected = level;
										show_level_modal = true;
									}}
								>
									{#if available_levels[level - 1].max_score != -1}
										<img src="/skull.svg" alt="lock" class="w-8 h-8" />
									{:else}
										<img src="/Lock.svg" alt="lock" class="w-8 h-8" />
									{/if}
								</button>
							{/if}
						{/each}
					</div>
					{#if row_index % 2 == 0 && row_index != level_rows.length - 1}
						<div class="relative w-10/12 mx-auto flex items-center justify-end h-[50px]">
							<img
								src="/dotted_line_vertical.svg"
								alt="dotted line"
								class="h-[50px] absolute right-8 translate-x-1/2"
							/>
						</div>
					{/if}
					{#if row_index % 2 != 0 && row_index != level_rows.length - 1}
						<div class="relative w-10/12 mx-auto flex items-center justify-start h-[50px]">
							<img
								src="/dotted_line_vertical.svg"
								alt="dotted line"
								class="h-[50px] absolute left-8 translate-x-1/2"
							/>
						</div>
					{/if}
				{/each}
			</div>
		</div>
	</main>
{/if}
