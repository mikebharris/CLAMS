<script lang="ts">
    export let attendee: JSON;
    export let i: Number;
</script>

<article>
    <h2>#{i}: <a href="#/attendee/{attendee.AuthCode}">{attendee.AuthCode}</a></h2>
    <p class="meta"><a href="#/attendee/{attendee.AuthCode}">{attendee.Name}</a> is coming on {attendee.ArrivalDay}</p>
</article>

<style>
    article {
        position: relative;
        padding: 0 0 0 2em;
        border-bottom: 1px solid #eee;
    }

    h2 {
        font-size: 1em;
        margin: 0.5em 0;
    }

    a {
        color: #333;
    }
</style>