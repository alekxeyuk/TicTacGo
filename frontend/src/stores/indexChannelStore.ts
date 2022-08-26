import { writable } from "svelte/store";

export const createChannelStore = (channelId: string) => {
    const { subscribe, set } = writable('');

    const eventSource = new EventSource(
        `http://localhost:80/stream/${channelId}`
    );

    eventSource.addEventListener('time', (e) => {
        set(e.data);
    });

    return {
        subscribe,
        reset: () => set(''),
        close: eventSource.close,
    };
};