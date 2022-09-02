<script lang="ts">
    import { onMount } from "svelte";
    import { createChannelStore } from "../../stores/indexChannelStore";
    import { Game } from "./game";
    import Board from "./Board.svelte";

    let messages: string[] = [];
    let gameInstance = new Game();
    let roomId = "";
    $: status = `Next player: ${$gameInstance.isXNext ? 'X' : 'O'}`;

    function handleBoardMessage(
        a: CustomEvent<{ action: string; index: number }>
    ) {
        console.log("handleBoardMessage", a.detail);
        messages = [...messages, a.detail.action + ": " + a.detail.index];

        switch (a.detail.action) {
            case "move":
                gameInstance.move(a.detail.index);
                break;
            // case "reset":
            //     gameInstance.reset();
            //     break;
            // case "start":
            //     gameInstance.start();
            //     break;
            // case "stop":
            //     gameInstance.stop();
            //     break;
        }
    }

    onMount(async () => {
        roomId = await gameInstance.init();

        const serverEvents = createChannelStore(roomId, "time", true);

        serverEvents.subscribe((incomingEvent) => {
            messages = [...messages, incomingEvent];
            console.log(incomingEvent);
        });

        return serverEvents.close;
    });
</script>


<div class="game">
    <div class="game-board">
        <Board {gameInstance} on:gameEvent={handleBoardMessage} />
    </div>
    <div class="game-info">
        <h2>{roomId}</h2>
        <div>{status}</div>
        <ul>
            {#each messages as message}
                <li>{message}</li>
            {/each}
        </ul>
    </div>
</div>