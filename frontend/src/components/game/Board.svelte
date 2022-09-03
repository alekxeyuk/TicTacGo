<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Game } from './game';

    const dispatch = createEventDispatcher<{
        gameEvent: { action: string; index: number };
    }>();

    export let gameInstance: Game;

    function handleInput(index: number) {
        dispatch("gameEvent", { action: "click", index });
    }
</script>

<div id="game-board">
    {#each { length: 3 } as _, i}
        <div class="table clear-both">
            {#each { length: 3 } as _, j}
                <button
                    class="bg-white border-2 border-solid border-slate-800 hover:border-rose-400 float-left text-2xl font-bold h-10 w-10 m-px p-0 text-center"
                    on:click={() => handleInput(i * 3 + j)}
                >
                    { $gameInstance.board[i * 3 + j] }
                </button>
            {/each}
        </div>
    {/each}
</div>
