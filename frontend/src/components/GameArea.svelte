<script lang="ts">
	import BlockButton from './BlockButton.svelte';
	import VacantBlock from './VacantBlock.svelte';
	import HintedBlock from './HintedBlock.svelte';
	import SkeletonBlock from './SkeletonBlock.svelte';
	import FilledBlock from './FilledBlock.svelte';

	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { EXPOSED_BASE_API_URL } from '$src/constants';
	import { createEventDispatcher } from 'svelte';

	export let skeleton: Record<string, string>;
	export let solution_length: number;
	export let blocks: string[];
	export let selecting_block_to_fill: boolean;

	const SOURCE_BLOCK_LIST = 'block-list';
	const SOURCE_SOLUTION = 'solution';

	const dispatch = createEventDispatcher();
	let hinted_content = 'fill';
	let selected_block_index: number | null = -1;
	let answer_correctly: boolean | null = null;
	let answer: (string | null)[] = [];
	let matches: boolean[] | null = null;

	let possibilities = JSON.parse(JSON.stringify(blocks)) as string[];

	onMount(() => {
		answer = Array.from({ length: solution_length }).map((_, index) => {
			if (skeleton[index.toString()]) return skeleton[index.toString()];
			return null;
		});
	});

	async function inject() {
		if (answer.includes(null)) {
			alert('Please fill all the boxes');
			return;
		}

		const url = `${EXPOSED_BASE_API_URL}/blocks/${$page.params.id}/check-answer`;

		const res = await fetch(url, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				blocks: answer
			})
		});
		if (res.status === 200) {
			const resultBody = await res.json();
			matches = answer.map((_) => true);
			answer_correctly = true;
			dispatch('submittedAnswer', {
				answer_correctly,
				next_level_id: resultBody.next_level_id,
				score: resultBody.score,
				multiplier: resultBody.multiplier
			});
		} else if (res.status === 400) {
			matches = (await res.json()).matches;
			answer_correctly = false;
			dispatch('submittedAnswer', { answer_correctly });
		} else {
			alert('Something went wrong');
		}
	}

	function onRemove(removed_index: number, text: string) {
		possibilities.push(text);
		possibilities = possibilities;
		answer[removed_index] = null;
	}

	function removeOnHint(removed_index: number, text: string) {
		possibilities = possibilities.filter((item) => item !== text);
	}

	async function fill_block_selected(index: number) {
		let res = fetch(`${EXPOSED_BASE_API_URL}/game/${$page.params.id}/hint`, {
			method: 'POST',
			credentials: 'include',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				hint_type: 'fill',
				order: index
			})
		});

		let res_json = await res;
		if (res_json.status !== 200) {
			alert('Something went wrong');
			return null;
		}
		let res_body = await res_json.json();
		removeOnHint(index, res_body.hintContent);

		return res_body.hintContent;
	}

	function removeDuplicates(correct_index: number | null, text: string) {
		if (correct_index === null) return;
		for (let i = 0; i < answer.length; i++) {
			if (answer[i] === text && i !== correct_index) {
				answer[i] = null;
			}
		}
	}

	let highlightedCellIndex: number | null = null;
</script>

<main
	class={`${
		answer_correctly !== null ? 'disabled pb-6' : 'pb-0'
	} bg-black text-sm flex flex-col items-center justify-between border-8 border-black pb-0 rounded-lg font-game absolute z-40 top-[7.5%] right-[11%]  w-[54%] h-[55%]`}
>
	<div id="input-area" class="flex flex-row items-center justify-start flex-wrap">
		{#each answer as _block, index}
			{#if skeleton[index.toString()]}
				<!-- If the block is also in the skeleton -->
				<!-- Same as skeleton[index.toString()]==={answer} -->
				<SkeletonBlock text={skeleton[index.toString()]} />
			{:else if answer[index]}
				<!-- key forces re-render of children elements when the value changes -->
				{#key matches && matches[index]}
					<!-- else, if it is only in the answer, it means that the user filled the box -->
					<FilledBlock
						on:dragstart={(event) => {
							event.dataTransfer?.setData(
								'text/plain',
								JSON.stringify({ index, source: SOURCE_SOLUTION })
							);
						}}
						on:remove={() => onRemove(index, answer[index])}
						text={answer[index]}
						isCorrect={matches && matches[index]}
					/>
				{/key}
			{:else}
				<!-- svelte-ignore a11y-no-static-element-interactions -->
				<div
					class="w-[20%] m-3"
					on:dragenter={() => {
						highlightedCellIndex = index;
					}}
					on:dragleave={() => {
						highlightedCellIndex = null;
					}}
					on:drop|preventDefault={(event) => {
						highlightedCellIndex = null;
						if (!event.dataTransfer) return;
						const json = event.dataTransfer.getData('text/plain');
						if (!json) return;
						const data = JSON.parse(json); // { source: string, index: number }

						if (data.source === SOURCE_SOLUTION) {
							// If dragging from another solution block, then move the block to the current index
							const source_index = data.index;

							const text = answer[source_index];

							answer[source_index] = null;
							answer[index] = text;
						} else if (data.source === SOURCE_BLOCK_LIST) {
							// Else we're dragging from the block list, then move the block to the current index from the possibilities
							const possibility_index = data.index;

							const text = possibilities[possibility_index];

							possibilities = possibilities.toSpliced(possibility_index, 1);
							answer[index] = text;
						}
					}}
					ondragover="return false"
				>
					<!-- else, it is empty -->
					{#if selecting_block_to_fill}
						<VacantBlock
							on:block_selected={async (e) => {
								hinted_content = await fill_block_selected(e.detail.block_index);
								if (hinted_content === null) return;
								selected_block_index = e.detail.block_index;
								dispatch('blockFilled');
								selecting_block_to_fill = false;
								answer[e.detail.block_index] = hinted_content;
								removeDuplicates(selected_block_index, hinted_content);
							}}
							hinted_content={selected_block_index === index ? hinted_content : 'fill'}
							self_index={index}
							{selecting_block_to_fill}
							highlighted={highlightedCellIndex === index}
						/>
					{:else if selected_block_index === index}
						<HintedBlock text_content={hinted_content} isCorrect={matches && matches[index]} />
					{:else}
						<VacantBlock
							hinted_content="fill"
							self_index={index}
							{selecting_block_to_fill}
							highlighted={highlightedCellIndex === index}
						/>
					{/if}
				</div>
			{/if}
		{/each}
	</div>

	<div
		id="block-list"
		class="h-[35%] overscroll-contain overflow-y-scroll flex flex-row items-center justify-start flex-wrap p-3 border-2 border-game-green w-[95%] justify-self-end gap-3 mb-3"
		on:drop|preventDefault={(event) => {
			if (!event.dataTransfer) return;
			const json = event.dataTransfer.getData('text/plain');
			if (!json) return;
			const data = JSON.parse(json); // { text: string, index: number }

			if (data.source === SOURCE_BLOCK_LIST) return;
			const answer_index = data.index;

			const text = answer[answer_index];

			possibilities.push(text);
			possibilities = possibilities;
			answer[answer_index] = null;
		}}
		ondragover="return false"
	>
		{#each possibilities as block, index}
			<BlockButton
				on:dragstart={(event) => {
					event.dataTransfer.setData(
						'text/plain',
						JSON.stringify({ index, source: SOURCE_BLOCK_LIST })
					);
				}}
				text={block}
				disabled={answer_correctly !== null}
			/>
		{/each}
	</div>

	{#if answer_correctly === null}
		<button
			id="inject-button"
			on:click={inject}
			class="px-[8%] py-[0.8%] bg-game-green text-black text-[90%] font-semibold"
		>
			Inject
		</button>
	{/if}
</main>
