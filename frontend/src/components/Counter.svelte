<script lang="ts">
    import { onMount } from "svelte";
    import { createChannelStore } from "../stores/indexChannelStore";


	let count = 0;

    // get the rooms count from the server
    async function getRoomsCount() {
        const response = await fetch('http://localhost:80/room/count', {credentials: 'include'});
        const data = await response.json();
        count = data.room_count;
    }

    onMount(() => {
        getRoomsCount();

        const serverEvents = createChannelStore('global', 'count', true);

        serverEvents.subscribe(incomingCount => {
            count = parseInt(incomingCount);
        });

        return serverEvents.close;
    });
</script>

<span>
	There is {count} open {count === 1 ? 'room' : 'rooms'}
</span>