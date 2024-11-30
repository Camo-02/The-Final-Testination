<script lang="ts">
	import Cross from '$components/Cross.svelte';

	const correctClasses = 'bg-game-green text-black border-game-green';
	const incorrectClasses = 'bg-game-red text-black border-game-red';
	const defaultClasses = 'bg-game-green text-black';

	let classes = defaultClasses;
	export let isCorrect: boolean | null;
	export let text_content: string;

	if (isCorrect === true) {
		classes = correctClasses;
	} else if (isCorrect === false) {
		classes = incorrectClasses;
	}
</script>

<div class="w-auto m-2 px-5 py-3 rounded-xl overflow-hidden relative group {classes}">
	{#if isCorrect === null}
		<h1 class="flex-grow w-full text-[90%] text-center">{text_content}</h1>
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
