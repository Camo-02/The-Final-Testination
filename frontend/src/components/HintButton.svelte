<script lang="ts">
	export let used: boolean;
	export let price: number;
	export let name: string;
	export let icon: string;
	export let totalCoins: number;
</script>

<button
	disabled={used || totalCoins < price}
	class="flex flex-col gap-2 items-center justify-start
            disabled:cursor-not-allowed
            hover:scale-105 disabled:hover:scale-100 disabled:hover:bg-inherit hover:bg-game-cyan hover:bg-opacity-50
            transition ease-in-out duration-300 py-4 px-2 w-fit"
	on:click
>
	<img
		src={icon}
		alt="{name} icon"
		class="w-[25%] h-[25%]"
		class:opacity-50={used || totalCoins < price}
	/>

	<div class="flex flex-col" class:opacity-50={used || totalCoins < price}>
		<h1 class="text-xl font-semibold">{name}</h1>
		<h2>
			{used ? 'Used' : `Price: ${price}`}
		</h2>
	</div>
</button>
