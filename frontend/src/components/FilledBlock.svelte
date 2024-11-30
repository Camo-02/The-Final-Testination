<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Cross from '$components/Cross.svelte';

	export let text: string;
	export let isCorrect: boolean | null;

	const correctClasses = 'bg-game-green text-black border-game-green';
	const incorrectClasses = 'bg-game-red text-black border-game-red';
	const defaultClasses = 'bg-black text-game-green border-game-green';

	let classes = defaultClasses;
	if (isCorrect === true) {
		classes = correctClasses;
	} else if (isCorrect === false) {
		classes = incorrectClasses;
	}

	const dispatch = createEventDispatcher();
</script>

<div
	class="w-[19%] m-2 px-5 py-3 border-2 border-game-green rounded-xl overflow-hidden relative group {classes}"
	on:dragstart
	draggable={isCorrect === null}
>
	<div class="flex-grow w-full text-[90%] text-center">{text}</div>
	{#if isCorrect === null}
		<button
			on:click={() => {
				dispatch('remove');
			}}
			class="flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all bg-game-green ease-in-out duration-200 px-[1%] py-[1%] w-[20%] h-[101%] absolute top-0 right-0"
		>
			<img src="/trashcan.svg" alt="close" class="w-[60%] h-[60%]" />
		</button>
	{:else if isCorrect}
		<div
			class="flex items-center justify-center px-[1%] py-[1%] w-[15%] h-[101%] absolute top-0 right-0"
		>
			<img src="/tick.svg" alt="tick" class="w-[60%] h-[60%]" />
		</div>
	{:else}
		<div
			class="flex items-center justify-center px-[1%] py-[1%] w-[15%] h-[101%] absolute top-0 right-0"
		>
			<Cross width="w-[50%]" height="h-[50%]" />
		</div>
	{/if}
</div>
