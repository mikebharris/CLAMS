<script lang="ts">
    let stats: JSON
    $: fetch(`${process.env.API_GATEWAY_URL}/report`)
        .then(r => r.json())
        .then(data => {
            stats = data;
            window.scrollTo(0, 0);
        });

    function toSquids(number: number): string {
        return Number(number).toLocaleString('en', { style: 'currency', currency: 'GBP' })
    }
</script>

{#if stats}
    <dl>
        <dt>Total Attendees</dt>
        <dd>{stats.TotalAttendees}</dd>
        <dt>Total Kids</dt>
        <dd>{stats.TotalKids}</dd>
        <dt>TotalNightsCamped</dt>
        <dd>{stats.TotalNightsCamped}</dd>
        <dt>TotalCampingCharge</dt>
        <dd>{toSquids(stats.TotalCampingCharge)}</dd>
        <dt>TotalPaid</dt>
        <dd>{toSquids(stats.TotalPaid)}</dd>
        <dt>TotalToPay</dt>
        <dd>{toSquids(stats.TotalToPay)}</dd>
        <dt>TotalIncome</dt>
        <dd>{toSquids(stats.TotalIncome)}</dd>
        <dt>AveragePaidByAttendee</dt>
        <dd>{toSquids(stats.AveragePaidByAttendee)}</dd>
    </dl>
{/if}

<style>
    dd {
        margin-bottom: 10px;
        text-align: right;
        margin-right:250px;
    }
</style>