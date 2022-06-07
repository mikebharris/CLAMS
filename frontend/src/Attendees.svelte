<script lang="ts">
	import { beforeUpdate } from "svelte";
	import Summary from "./Summary.svelte";

	let attendees: JSON[];

	$: fetch(`${process.env.API_GATEWAY_URL}/attendees`)
		.then(r => r.json())
		.then(data => {
			attendees = data.Attendees;
			window.scrollTo(0, 0);
		});
</script>

{#if attendees}
	{#each attendees as attendee, i}
		<Summary {attendee} {i}/>
	{/each}
{:else}
	<p class="loading">loading...</p>
{/if}

<style>
	.loading {
		opacity: 0;
		animation: 0.4s 0.8s forwards fade-in;
	}

	@keyframes fade-in {
		from { opacity: 0; }
		to { opacity: 1; }
	}
</style>