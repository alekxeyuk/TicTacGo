<script lang="ts">
    import { onMount } from "svelte";
    import { createChannelStore } from "../../stores/indexChannelStore";
    import { Game } from "./game";
    import Board from "./Board.svelte";

    // let messages: string[] = [];
    let gameInstance = new Game();
    let roomId = "";
    $: status = `Next player: ${$gameInstance.isXNext ? "O" : "X"}`;

    function handleBoardMessage(
        a: CustomEvent<{ action: string; index: number }>
    ) {
        console.log("handleBoardMessage", a.detail);

        switch (a.detail.action) {
            case "click":
                gameInstance.sendClick(a.detail.index);
                break;
        }
    }

    function handleServerMessage(a: string) {
        console.log("handleGameMessage", a);

        const data = JSON.parse(a);

        switch (data.action) {
            case "move":
                gameInstance.placeSignAtIndex(data.index, data.sign);
                break;
            case "win":
                gameInstance.setWinner(data.sign);
                break;
        }
    }

    onMount(async () => {
        roomId = await gameInstance.init();

        const serverEvents = createChannelStore(roomId, "game", true);

        serverEvents.subscribe((incomingEvent) => {
            if (incomingEvent !== "") {
                handleServerMessage(incomingEvent);
            }
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
        <div>{status}</div>
        <div>Your sign is: {$gameInstance.mySign}</div>
        {#if $gameInstance.winner !== ""}
            <div>
                Winner: {$gameInstance.winner === $gameInstance.mySign
                    ? "You"
                    : "Opponent"}
            </div>
        {/if}
    </div>
</div>

<style>
    .game {
        display: flex;
        flex-direction: row;
    }

    .game-board {
        margin-right: 20px;
    }

    .game-info {
        margin-top: 20px;
    }
</style>