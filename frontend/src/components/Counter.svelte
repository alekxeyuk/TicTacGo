<script type="ts">
    import { onMount } from "svelte";
    import { createChannelStore } from "../stores/indexChannelStore";


	let count = 0;

    // get the rooms count from the server
    async function getRoomsCount() {
        const response = await fetch('http://localhost:80/room/count');
        const data = await response.json();
        count = data.count;
    }

    onMount(() => {
        getRoomsCount();

        const store = createChannelStore('global');

        store.subscribe(incomingMessages => {
            console.log(incomingMessages);
        });

        return store.close;
    });
</script>

<span>
	There is {count} open {count === 1 ? 'room' : 'rooms'}
</span>