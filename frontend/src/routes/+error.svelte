<script lang="ts">
	import { page } from '$app/stores';

	let errorMsg = 'There was an error handling your request. Please try again later.';

	if ($page.error?.message) {
		const body = JSON.parse($page.error!.message);
		if (body.error) {
			errorMsg = body.error;
		} else if (body.queryError) {
			errorMsg = body.queryError;
		}
	}

	JSON.parse($page.error?.message ?? '').mes;
</script>

<main class="h-full w-full flex flex-col gap-6 items-center mt-12">
	<h1 class="text-3xl">Error</h1>
	<h2>{errorMsg}</h2>
	<a
		href="/"
		class="bg-primary hover:bg-theme text-platform-white font-platform font-bold py-2 px-4 rounded-full"
	>
		Go home
	</a>
</main>
