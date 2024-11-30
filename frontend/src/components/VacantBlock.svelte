<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	const emit = createEventDispatcher();
	export let highlighted: boolean;
	export let selecting_block_to_fill: boolean;
	export let self_index: number;
	export let hinted_content: string;
	const dispatch = createEventDispatcher();
	let normal_class = 'text-black bg-game-gray rounded-xl text-center text-md py-3 px-5 w-full';
	let selecting_class =
		'text-black bg-game-green text-black rounded-xl text-center text-md py-3 px-5 w-full animate-pulse';
</script>

<button
	on:click={() => {
		if (selecting_block_to_fill) {
			let block_index = self_index;
			emit('block_selected', { block_index });
		}
		dispatch('emptySelect');
	}}
	class:opacity-75={highlighted}
	class={selecting_block_to_fill ? selecting_class : normal_class}
>
	{hinted_content}
</button>
