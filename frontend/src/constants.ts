import { PUBLIC_API_PORT, PUBLIC_API_HOST, PUBLIC_EXPOSED_API_HOST } from '$env/static/public';

export const API_PORT = PUBLIC_API_PORT;
export const API_HOST = PUBLIC_API_HOST;
export const BASE_API_URL = `http://${API_HOST}:${API_PORT}`;
export const EXPOSED_BASE_API_URL = `http://${PUBLIC_EXPOSED_API_HOST}:${API_PORT}`;
// Using `.MODE` since it is set when changing the `--mode` flag in `vite`
// Seems like the `.DEV` variable is always `false` when running `npm run build`
export const IS_DEV_MODE = import.meta.env.MODE === 'development';
export const PAGE_SIZE = 25;
