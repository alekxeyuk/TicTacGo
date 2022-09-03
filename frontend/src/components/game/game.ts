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
        changeCell: (index: number, sign: string) => update((state: typeof defaultState) => {
            // const newBoard = [...state.board];
            state.board[index] = sign === 'PLAYER_X' ? 'X' : 'O';
            return {
                board: state.board,
                isXNext: !state.isXNext,
            };
        }),
        updateBoard: (board: number[]) => update((state: typeof defaultState) => {
            board.forEach((sign, index) => {
                state.board[index] = sign === 0 ? '' : sign === 1 ? 'O' : 'X';
            });
            return {
                board: state.board,
                isXNext: !state.isXNext,
            };
        }),
    };
}

async function getRandomRoom(): Promise<{ uuid: string; state: number[] }> {
    const response = await fetch("http://localhost:80/room/random", { credentials: 'include' });
    return response.json();
}

// Game class definition
export class Game {
    public store = createStore();
    private mutex = false;
    private channelId = '';

    async init(): Promise<string> {
        const roomData = await getRandomRoom();
        this.channelId = roomData.uuid;
        this.store.updateBoard(roomData.state);
        return this.channelId;
    }

    sendClick(index: number) {
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
    }

    reset() {
        this.store.reset();
    }

    placeSignAtIndex(index: number, sign: string) {
        this.store.changeCell(index, sign);
    }

    updateBoard(board: number[]) {
        this.store.updateBoard(board);
    }

    subscribe(callback: (state: typeof defaultState) => void) {
        return this.store.subscribe(callback);
    }
}
