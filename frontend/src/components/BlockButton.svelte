<script lang="ts">
	export let text: string;
	export let disabled: boolean;
</script>

<button on:dragstart draggable={true} {disabled} class="w-[19%] m-3">
	<div class="bg-black text-game-green border-game-green border-2 rounded-xl py-4 px-6">
		<h1 class="p-auto text-md">{text}</h1>
	</div>
</button>
