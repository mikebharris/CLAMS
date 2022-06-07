<script lang="ts">
    import {onMount} from 'svelte';
    import Attendees from './Attendees.svelte';
    import Attendee from './Attendee.svelte';
    import Report from './Report.svelte';

    let attendee: JSON;

    console.log(process.env.API_GATEWAY_URL)

    async function hashchange() {
        // the poor man's router!
        const path = window.location.hash.slice(1);

        if (path.startsWith('/attendee/')) {
            const authCode = path.slice(10);
            await fetch(`${process.env.API_GATEWAY_URL}/attendee/${authCode}`)
                .then(r => r.json())
                .then(data => {
                    attendee = data.Attendees[0];
                    window.scrollTo(0, 0);
                });
            console.log(attendee)
            window.scrollTo(0, 0);
        } else if (path.startsWith('/attendees')) {
            attendee = null;
        } else {
            window.location.hash = '/report';
        }
    }

    onMount(hashchange);
</script>

<svelte:window on:hashchange={hashchange}/>

<main>
    <h1>Welcome to CLAMS</h1>
    <div class="attendees">
        <h2>Attendees</h2>
        {#if attendee}
            <Attendee {attendee} returnTo="#/attendees"/>
        {:else}
            <Attendees/>
        {/if}
    </div>
    <div class="stats">
        <h2>Some Figures</h2>
        <Report/>
    </div>
</main>

<style>
    main {
        position: relative;
        max-width: 800px;
        margin: 0 auto;
        min-height: 101vh;
        padding: 1em;
    }

    div.attendees {
        width: 50%;
        float: left
    }

    div.stats {
        width: 50%;
        float: right
    }

    main :global(.meta) {
        color: #999;
        font-size: 12px;
        margin: 0 0 1em 0;
    }

    main :global(a) {
        color: rgb(0, 0, 150);
    }
</style>