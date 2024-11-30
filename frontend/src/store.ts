import { writable, type Writable } from 'svelte/store';

export const fetchingpage: Writable<boolean> = writable(false);
