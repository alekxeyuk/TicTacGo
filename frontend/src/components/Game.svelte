<script lang="ts">
    import { onMount } from "svelte";
    import { createChannelStore } from "../stores/indexChannelStore";

    let roomId = "";
    let messages: string[] = [];

    async function getRandomRoom() {
        const response = await fetch("http://localhost:80/room/random", {credentials: 'include'});
        const data = await response.json();
        return data.uuid;
    }

    onMount(async () => {
        roomId = await getRandomRoom();

        const serverEvents = createChannelStore(roomId, "time", true);

        serverEvents.subscribe((incomingEvent) => {
            messages = [...messages, incomingEvent];
            console.log(incomingEvent);
        });

        return serverEvents.close;
    });
</script>

<h2>{roomId}</h2>
<ul>
    {#each messages as message}
        <li>{message}</li>
    {/each}
</ul>
