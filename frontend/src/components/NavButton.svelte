<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	let { link, buttonText, color } = $props();
</script>

<!-- Safelist: text-logo-blue shadow-logo-red shadow-logo-blue -->
<button
	onclick={() =>
		goto('/login', {
			state: {
				...page.state,
				isLoginSelected: link == '/login'
			}
		})}
	class="text-{color} bg-neutral rounded-xl font-sans font-bold shadow-md hover:shadow-lg shadow-{color} m-2 p-3"
>
	{buttonText}
</button>
