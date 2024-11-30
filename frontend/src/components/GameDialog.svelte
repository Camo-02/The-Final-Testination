<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import Cross from '$components/Cross.svelte';

	export let text: string;
	export let title: string;

	const dispatch = createEventDispatcher();
</script>

<div
	id="game-dialog"
	class="bg-game-blue text-game-cyan border-4 border-game-cyan flex flex-col items-center justify-center p-5 max-w-[40vw] h-fit z-50 relative"
>
	<h2 class="font-bold text-5xl mb-10">{title}</h2>

	<button
		id="close-button"
		class="w-4 h-4 absolute top-9 right-10 cursor-pointer"
		on:click={() => {
			dispatch('dialogClose');
		}}
	>
		<Cross fill="fill-game-cyan" />
	</button>
	<p class="w-10/12 text-left mb-10 text-2xl" style="white-space: pre-line">{text}</p>

	<slot />
</div>
