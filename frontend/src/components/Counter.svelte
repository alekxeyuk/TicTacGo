<script>
    import { onMount } from "svelte";

    
	let count = 0;

    // get the rooms count from the server
    async function getRoomsCount() {
        const response = await fetch('/api/rooms');
        const data = await response.json();
        return data.count;
    }

    onMount(async () => {
        count = await getRoomsCount();
    });
</script>

<span>
	There is {count} open {count === 1 ? 'room' : 'rooms'}
</span>