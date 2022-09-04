import { writable } from "svelte/store";

export const createChannelStore = (channelId: string, eventName: string, withAuth: boolean) => {
    const { subscribe, set } = writable('');

    const eventSource = new EventSource(
        `https://tictacgo-production.up.railway.app/room/${channelId}/stream`,
        { withCredentials: withAuth }
    );

    eventSource.addEventListener(eventName, (event) => {
        set(event.data);
    });

    return {
        subscribe,
        reset: () => set(''),
        close: eventSource.close,
    };
};