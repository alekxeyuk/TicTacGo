import { writable } from "svelte/store";

const defaultState = {
    board: Array(9).fill(''),
    isXNext: true,
};

function createStore() {
    const { subscribe, set, update } = writable(defaultState);

    return {
        subscribe,
        reset: () => set(defaultState),
        updateBoard: (index: number) => update((state) => {
            // const newBoard = [...state.board];
            state.board[index] = state.isXNext ? 'X' : 'O';
            return {
                board: state.board,
                isXNext: !state.isXNext,
            };
        }),
    };
}

async function getRandomRoom(): Promise<{uuid:string}> {
    const response = await fetch("http://localhost:80/room/random", {credentials: 'include'});
    return response.json();
}

// Game class definition
export class Game {
    public store = createStore();
    private mutex = false;
    private channelId = '';

    async init(): Promise<string> {
        const channelId = await getRandomRoom();
        this.channelId = channelId.uuid;
        return this.channelId;
    }

    move(index: number) {
        if (this.mutex) {
            return;
        }
        this.mutex = true;

        fetch(`http://localhost:80/room/${this.channelId}/move?cell=${index}`, {
            method: 'PUT',
            credentials: 'include',
        }).then(() => {
                // this.mutex = false;
        }).finally(() => {
            this.mutex = false;
        });

        this.store.updateBoard(index);
    }

    reset() {
        this.store.reset();
    }

    updateBoard(index: number) {
        this.store.updateBoard(index);
    }

    subscribe(callback: (state: typeof defaultState) => void) {
        return this.store.subscribe(callback);
    }
}
