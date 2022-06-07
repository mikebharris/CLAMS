<script lang="ts">
    export let attendee: JSON;
    export let returnTo: String;
</script>

<a href={returnTo}>&laquo; back</a>
<article>
    <h1>{attendee.Name}</h1>
    <dl>
        <dt>AuthCode</dt>
        <dd>{attendee.AuthCode}</dd>
        <dt>Email</dt>
        <dd><a href="mailto:{attendee.Email}">{attendee.Email}</a></dd>
        <dt>Telephone</dt>
        <dd>{attendee.Telephone}</dd>
        <dt>Kids Count</dt>
        <dd>{attendee.NumberOfKids}</dd>
        <dt>Diet</dt>
        <dd>{attendee.Diet}</dd>
        <dt>Money Owed</dt>
        <dd>
            {#if attendee.Financials.AmountDue > 0}
                <div class="indebt">£{attendee.Financials.AmountDue}.00p</div>
            {:else}
                <div class="paidup">None :-)</div>
            {/if}
        </dd>
    </dl>
</article>

<style>
    article {
        margin: 0 0 1em 0;
    }

    a {
        display: block;
        margin: 0 0 1em 0;
    }

    h1 {
        font-size: 1.4em;
        margin: 0;
    }

    dd {
        margin-bottom: 10px
    }

    div.indebt {
        color: red;
        font-weight: bold;
        text-blink: yes;
    }

    div.paidup {
        color: green;
        font-weight: bold;
        text-blink: yes;
    }
</style>