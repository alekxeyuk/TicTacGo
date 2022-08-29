import { writable } from "svelte/store";

export const createChannelStore = (channelId: string, eventName: string, withAuth: boolean) => {
    const { subscribe, set } = writable('');

    const eventSource = new EventSource(
        `http://localhost:80/room/${channelId}/stream`,
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