import { writable } from "svelte/store";

const defaultState = {
    board: Array(9).fill(''),
    mySign: '',
    winner: '',
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
                mySign: state.mySign,
                winner: state.winner,
                isXNext: !state.isXNext,
            };
        }),
        setWinner: (winner: string) => update((state: typeof defaultState) => {
            return {
                board: state.board,
                mySign: state.mySign,
                winner: winner === 'PLAYER_O' ? 'O' : 'X',
                isXNext: state.isXNext,
            };
        }),
        updateBoard: (board: number[], sign: string, is_x_next: boolean) => update((state: typeof defaultState) => {
            board.forEach((sign, index) => {
                state.board[index] = sign === 0 ? '' : sign === 1 ? 'O' : 'X';
            });
            return {
                board: state.board,
                mySign: sign === 'PLAYER_X' ? 'O' : 'X',
                winner: state.winner,
                isXNext: is_x_next,
            };
        }),
    };
}

async function getRandomRoom(): Promise<{ uuid: string; state: number[]; sign: string; is_x_next: boolean }> {
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
        this.store.updateBoard(roomData.state, roomData.sign, roomData.is_x_next);
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

    setWinner(sign: string) {
        this.store.setWinner(sign);
    }

    subscribe(callback: (state: typeof defaultState) => void) {
        return this.store.subscribe(callback);
    }
}
