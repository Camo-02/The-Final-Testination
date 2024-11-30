<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import HintButton from './HintButton.svelte';

	export let storyContent: string;
	export let cheatsheetContent: string;
	export let textualHintUsedInThisSession: boolean;
	export let freezeTimeUsedInThisSession: boolean;
	export let fillBlockUsedInThisSession: boolean;
	export let hintsContent: {
		playerCoins: number;
		timeFreezePrice: number;
		timeFreezeUsed: number;
		solutionHintPrice: number;
		solutionHintUsed: number;
		textualHintPrice: number;
		textualHintUsed: number;
	};

	const emit = createEventDispatcher();

	let showing: 'Cheatsheet' | 'Story' | 'Hints' = 'Story';
	function showCheatsheet() {
		showing = 'Cheatsheet';
	}

	function showStory() {
		showing = 'Story';
	}

	function showHints() {
		showing = 'Hints';
	}

</script>

<div
	class="absolute top-[3%] left-[1%] text-white bg-game-blue border-4 border-game-cyan w-[30%] h-[90%] z-[45] font-game bg-opacity-80 backdrop-blur-[3px]"
>
	<ul class="list-style-none flex flex-row justify-evenly">
		<li class="w-full">
			<button
				id="story"
				on:click={showStory}
				class={`${
					showing === 'Story' ? '' : 'bg-game-cyan'
				} w-full py-3 flex items-center justify-center text-[120%] font-bold`}>Story</button
			>
		</li>
		<li class="w-full">
			<button
				id="cheatsheet"
				class={`${
					showing === 'Cheatsheet' ? '' : 'bg-game-cyan'
				} w-full py-3 flex items-center justify-center text-[120%] font-bold bg-opacity-1`}
				on:click={showCheatsheet}>Info</button
			>
		</li>
		<li class="w-full">
			<button
				id="hints"
				class={`${
					showing === 'Hints' ? '' : 'bg-game-cyan'
				} w-full py-3 flex items-center justify-center text-[120%] font-bold bg-opacity-1`}
				on:click={showHints}>Hints</button
			>
		</li>
	</ul>

	<div
		class="scrollbarcustom overflow-y-scroll scrollbar scrollbar-thumb-game-cyan scrollbar-track-transparent w-full h-[90%] relative p-[3%] py-[5%] overflow-x-auto"
	>
		<p id="sidebar" class="text-lg whitespace-pre-wrap break-words overflow-auto">
			{#if showing === 'Cheatsheet'}
				<p class="text-[80%] whitespace-pre-wrap break-words overflow-auto">{cheatsheetContent}</p>
			{:else if showing === 'Story'}
				<p class="text-[80%] whitespace-pre-wrap break-words overflow-auto">{storyContent}</p>
			{:else if showing === 'Hints'}
				<div class="flex flex-col w-full items-center">
					<div class="flex flex-row gap-4 items-center justify-center w-11/12 text-xl mx-auto mb-6">
						<img src="/coins.svg" alt="coin" class="w-[15%] h-[15%]" />
						<h1>Coins {hintsContent.playerCoins}</h1>
					</div>

					<HintButton
						on:click={() => {
							emit('fillBlock');
						}}
						used={hintsContent.solutionHintUsed > 0 || fillBlockUsedInThisSession}
						price={hintsContent.solutionHintPrice}
						name={'Fill Block'}
						icon={'/fill_a_block_icon.svg'}
						totalCoins={hintsContent.playerCoins}
					/>

					<HintButton
						on:click={() => {
							emit('freezeTime');
						}}
						used={hintsContent.timeFreezeUsed > 0 || freezeTimeUsedInThisSession}
						price={hintsContent.timeFreezePrice}
						name={'Freeze time'}
						icon={'/time_freeze_icon.svg'}
						totalCoins={hintsContent.playerCoins}
					/>
					<HintButton
						on:click={() => {
							emit('textualHint');
						}}
						used={textualHintUsedInThisSession ? false : hintsContent.textualHintUsed > 0}
						price={textualHintUsedInThisSession ? 0 : hintsContent.textualHintPrice}
						name={'Textual hint'}
						icon={'/textual_hint_icon.svg'}
						totalCoins={hintsContent.playerCoins}
					/>
				</div>
			{/if}
		</p>
	</div>
</div>
