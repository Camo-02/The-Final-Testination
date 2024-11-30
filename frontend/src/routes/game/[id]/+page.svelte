<script lang="ts">
	import GameArea from '$components/GameArea.svelte';
	import Cheatsheet from '$components/Cheatsheet.svelte';
	import GameDialog from '$components/GameDialog.svelte';
	import { BASE_API_URL } from '$src/constants';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import SplashScreen from '$components/SplashScreen.svelte';

	type GameData = {
		cheatsheet: string;
		story: string;
		skeleton: Record<number, string>;
		blocks: string[];
		solution_length: number;
		winning_message: string;
		background: string;
		title: string;
		player_coins: number;
		freeze_time_duration: number;
		freeze_time_price: number;
		freeze_time_used: number;
		textual_hint_price: number;
		textual_hint_used: number;
		solution_hint_price: number;
		solution_hint_used: number;
	};

	let answer: boolean | null = null;
	let next_level_id: string | null = null;
	let showSplash: boolean = true;
	let data: GameData = {
		title: '',
		background: '',
		winning_message: '',
		cheatsheet: '',
		story: '',
		skeleton: {},
		blocks: [],
		player_coins: 0,
		freeze_time_duration: 0,
		solution_length: 0,
		freeze_time_price: 0,
		freeze_time_used: 0,
		textual_hint_price: 0,
		textual_hint_used: 0,
		solution_hint_price: 0,
		solution_hint_used: 0
	};

	let selecting_block_to_fill = false;
	let time_is_frozen = false;
	let minutes = 0;
	let seconds = 0;
	let timer: string | number | NodeJS.Timeout | undefined;
	let fetching = true;
	let score = 0;
	let time_slot = 0;
	let time_text = 'Crap';
	let great = 80;
	let medium = 60;
	let notsogood = 40;
	let crap = 20;
	let perfect_class =
		'bg-gradient-to-r from-green-500 to-green-700 flex items-center justify-center rounded-r-full py-3';
	let great_class =
		'bg-gradient-to-r from-yellow-500 to-green-500 flex items-center justify-center py-3';
	let medium_class =
		'bg-gradient-to-r from-orange-500 to-yellow-500 flex items-center justify-center py-3';
	let notsogood_class =
		'bg-gradient-to-r from-red-500 to-orange-500 flex items-center justify-center py-3';
	let crap_class =
		'bg-gradient-to-r from-red-700 to-red-500 flex items-center justify-center rounded-l-full py-3';
	let not_reached_class = 'bg-black bg-opacity-50 flex items-center justify-center py-3';
	let not_reached_beginning_class =
		'bg-black bg-opacity-50 flex items-center justify-center rounded-l-full py-3';
	let not_reached_end_class =
		'bg-black bg-opacity-50 flex items-center justify-center rounded-r-full py-3';

	async function getGameData() {
		const res = await fetch(`${BASE_API_URL}/game/${$page.params.id}`, {
			credentials: 'include'
		});

		if (!res.ok) {
			if (res.status === 401) {
				window.location.replace('/login');
			} else {
				window.location.replace('/');
			}
			return;
		}

		data = await res.json();
		fetching = false;
	}

	onMount(async () => {
		await new Promise((r) => setTimeout(r, 3000));
		showSplash = false;
		await getGameData();
		if (imgContainer) {
			imgContainer.style.backgroundImage = `url("data:image/svg+xml,${encodeURIComponent(
				data.background
			)}")`;
		}

		timer = setInterval(() => {
			if (seconds === 60) {
				minutes += 1;
				seconds = 0;
			}
			if (!time_is_frozen) {
				seconds += 1;
			}
		}, 1000);
	});

	function gameResult(
		event: CustomEvent<{
			answer_correctly: boolean;
			next_level_id: string | null;
			score: number;
			multiplier: number;
		}>
	) {
		answer = event.detail.answer_correctly;
		if (!answer) {
			return;
		}
		clearInterval(timer);
		score = event.detail.score;
		next_level_id = event.detail.next_level_id;
		time_slot = event.detail.multiplier * 100;
		if (next_level_id == '') {
			next_level_id = null;
		}

		if (time_slot > great) {
			time_text = "Incredible! You've got it immediately!";
		} else if (time_slot > medium) {
			time_text = 'Great job! You are pretty fast, but there is still room for improvement.';
		} else if (time_slot > notsogood) {
			time_text = 'Nice job, although perhaps you should read a bit more about this subject.';
		} else if (time_slot > crap) {
			time_text = 'You really took your time, eh? Maybe go back to study again about this topic.';
		} else {
			time_text = 'Did you even read the assignment?';
		}
	}

	let menu_visible = false;
	let manuallyClosed = false;
	let tutorial_visible = false;
	let imgContainer: HTMLDivElement;

	async function useHint(type: 'freeze' | 'textual' | 'fill') {
		const res = await fetch(`${BASE_API_URL}/game/${$page.params.id}/hint`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ hint_type: type })
		});
		if (res.status !== 200) {
			alert('Something went wrong, please try again');
		}
		return res;
	}

	let textualHint = '';
	let textualHintVisible = false;
	let freezeHintUsed = false;
	let fillBlockUsed = false;
</script>

{#if fetching || showSplash}
	<SplashScreen />
{:else}
	<main class="min-h-screen">
		<div bind:this={imgContainer} class="w-full min-h-screen bg-cover relative">
			{#if menu_visible}
				<div
					class="z-[999] flex flex-col items-center justify-center w-8/12 h-[60%] absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-platform-white rounded-2xl border-4 border-platform-black"
				>
					<div class="w-full h-full relative flex flex-col items-center justify-center">
						<button
							class="absolute top-5 right-5"
							on:click={() => {
								menu_visible = false;
							}}
						>
							<img src="/cross.svg" alt="snowflake" class="w-5 h-5" />
						</button>

						<h1 class="text-7xl text-theme font-black font-platform mb-10 pt-10">Menu</h1>
						<div class="flex flex-row items-center justify-around w-3/12 mb-10 relative">
							{#if time_is_frozen}
								<img id="snowflake_side_l" src="/snowflake.svg" alt="snowflake" class="w-10 h-10" />
							{/if}
							<h1 id="timer" class="text-5xl text-theme font-bold font-platform">
								{minutes < 10 ? '0' + minutes : minutes}:{seconds < 10 ? '0' + seconds : seconds}
							</h1>
							{#if time_is_frozen}
								<img id="snowflake_side_r" src="/snowflake.svg" alt="snowflake" class="w-10 h-10" />
							{/if}
						</div>
						<a
							id="quitBtn"
							href="/"
							class="cursor-pointer text-theme hover:font-bold w-full text-center py-5 text-[150%] hover:bg-primary hover:bg-opacity-20 transition ease-in-out duration-300"
						>
							Quit the game
						</a>

						<button
							id="tutorialBtn"
							class="cursor-pointer text-theme hover:text-theme hover:font-bold w-full py-5 text-center text-[150%] hover:bg-primary hover:bg-opacity-20 transition ease-in-out duration-300"
							on:click={() => {
								tutorial_visible = true;
								menu_visible = false;
							}}
						>
							How to play
						</button>
					</div>
				</div>
			{/if}

			<div class="absolute top-[3%] right-[1%] flex flex-col gap-6 items-center">
				<button
					id="overlayMenuBtn"
					class="w-14 h-14"
					on:click={() => {
						menu_visible = !menu_visible;
					}}
				>
					<img
						alt="hamburger menu"
						src="/logobutton.svg"
						class="drop-shadow-[(0_35px_35px_rgba(0,0,0,1)] z-[999] opacity-50 hover:opacity-100 transition duration-300 ease-in-out w-full h-full"
					/>
				</button>
				{#if time_is_frozen}
					<img
						id="timeFreezeAd"
						src="/white_snowflake.svg"
						alt="time frozen icon"
						class="w-10 h-10 opacity-50 transform ease-in-out duration-300 hover:opacity-100"
					/>
				{/if}
			</div>
			{#if textualHintVisible}
				<div id="textualHintAd" class="flex items-center justify-center h-screen">
					<GameDialog
						on:dialogClose={() => {
							textualHintVisible = false;
						}}
						title="Textual Hint"
						text={textualHint}
					></GameDialog>
				</div>
			{/if}
			{#if tutorial_visible}
				<div class="flex items-center justify-center h-screen">
					<GameDialog
						on:dialogClose={() => {
							tutorial_visible = false;
						}}
						title="How to play"
						text="Your objective is to complete the input by dragging the desired blocks to the empty space.
						You can remove blocks by either using the trash icon or by dragging them back to their starting position.
						When you think you've completed the input correctly, click the Inject button.
						
						Your score is based on time, hints, and number of attempts.
						
						You have three hints at your disposal:
						- Textual Hint: Will give you a textual clue about the solution.
						- Time Freeze: Will freeze the timer for a certain amount of time.
						- Fill Block: Will fill a empty space of your choice.
						
						To use hints you have to spend your coins, which you earn by compleating each level.
						
						REMEMBER: Hints can be used only once per level, so use them wisely."
					></GameDialog>
				</div>
			{/if}
			{#if manuallyClosed}
				<button
					id="open-game-dialog-btn"
					class="z-[999] absolute top-[12%] right-[1.25%] bg-game-blue border-game-cyan border-4 text-game-cyan px-5 py-3"
					on:click={() => {
						manuallyClosed = false;
					}}
				>
					see result
				</button>
			{/if}
			<Cheatsheet
				cheatsheetContent={data.cheatsheet}
				storyContent={data.story}
				hintsContent={{
					playerCoins: data.player_coins,
					timeFreezePrice: data.freeze_time_price,
					timeFreezeUsed: data.freeze_time_used,
					textualHintPrice: data.textual_hint_price,
					textualHintUsed: data.textual_hint_used,
					solutionHintPrice: data.solution_hint_price,
					solutionHintUsed: data.solution_hint_used
				}}
				textualHintUsedInThisSession={textualHint !== ''}
				freezeTimeUsedInThisSession={freezeHintUsed}
				fillBlockUsedInThisSession={fillBlockUsed}
				on:fillBlock={() => {
					selecting_block_to_fill = true;
				}}
				on:freezeTime={() => {
					useHint('freeze').then((res) => {
						if (res.status !== 200) {
							return;
						}
						data.freeze_time_used = data.freeze_time_price;
						data.player_coins -= data.freeze_time_price;
						freezeHintUsed = true;
					});
					time_is_frozen = true;
					setTimeout(() => {
						time_is_frozen = false;
					}, data.freeze_time_duration * 1000);
				}}
				on:textualHint={() => {
					if (textualHint != '') {
						textualHintVisible = true;
						return;
					}
					useHint('textual').then(async (res) => {
						if (res.status !== 200) {
							return;
						}
						const json = await res.json();
						data.textual_hint_used = data.textual_hint_price;
						data.player_coins -= data.textual_hint_price;
						textualHint = json.hintContent;
						textualHintVisible = true;
					});
				}}
			/>
			<GameArea
				{selecting_block_to_fill}
				skeleton={data.skeleton}
				blocks={data.blocks}
				solution_length={data.solution_length}
				on:submittedAnswer={gameResult}
				on:blockFilled={() => {
					data.solution_hint_used = data.solution_hint_price;
					data.player_coins -= data.solution_hint_price;
					fillBlockUsed = true;
				}}
			/>
		</div>
	</main>

	{#if answer !== null && !manuallyClosed}
		<div
			class=" z-[50] absolute top-0 left-0 w-full h-full bg-black bg-opacity-70 flex flex-col items-center justify-center"
		>
			{#if answer === true}
				<GameDialog
					title="Success!"
					text={data.winning_message}
					on:dialogClose={() => {
						manuallyClosed = true;
					}}
				>
					{#if score !== 0}
						<!-- comment about time result -->
						<p class="mb-5 text-center">{time_text}</p>

						<!-- render time_slot in a progress line -->
						<div
							class="w-11/12 mx-auto h-5 bg-blue-800 bg-opacity-[50%] rounded-full grid grid-cols-5 mb-10"
						>
							<div class={time_slot >= 0 ? crap_class : not_reached_beginning_class}>
								{#if time_slot <= crap}
									<img src="/whiteskull.svg" alt="time skull" class="w-6 h-6" />
								{/if}
							</div>
							<div class={time_slot > crap ? notsogood_class : not_reached_class}>
								{#if time_slot <= notsogood && time_slot > crap}
									<img src="/whiteskull.svg" alt="time skull" class="w-6 h-6" />
								{/if}
							</div>
							<div class={time_slot > notsogood ? medium_class : not_reached_class}>
								{#if time_slot <= medium && time_slot > notsogood}
									<img src="/whiteskull.svg" alt="time skull" class="w-6 h-6" />
								{/if}
							</div>
							<div class={time_slot > medium ? great_class : not_reached_class}>
								{#if time_slot <= great && time_slot > medium}
									<img src="/whiteskull.svg" alt="time skull" class="w-6 h-6" />
								{/if}
							</div>
							<div class={time_slot > great ? perfect_class : not_reached_end_class}>
								{#if time_slot > great}
									<img src="/whiteskull.svg" alt="time skull" class="w-6 h-6" />
								{/if}
							</div>
						</div>
					{/if}

					<div class="flex flex-row gap-6 my-10">
						<a class="bg-game-cyan text-white px-5 py-3" href="/">exit the game</a>
						{#if next_level_id != null}
							<a
								class="bg-game-cyan text-white px-5 py-3"
								href="/game/{next_level_id}"
								on:click={() => {
									manuallyClosed = true;
									window.location.replace(`/game/${next_level_id}`);
								}}
								>next level
							</a>
						{/if}
					</div>
				</GameDialog>
			{:else if answer === false}
				<GameDialog
					on:dialogClose={() => {
						manuallyClosed = true;
					}}
					title="You lose!"
					text="You still have to practice on this stuff, maybe you're better on the field..."
				>
					<div class="flex flex-row gap-6">
						<a class="bg-game-cyan text-white px-5 py-3" href="/">exit the game</a>
						<button
							class="bg-game-cyan text-white px-5 py-3"
							on:click={() => {
								window.location.reload();
							}}>try again</button
						>
					</div>
				</GameDialog>
			{/if}
		</div>
	{/if}
{/if}
