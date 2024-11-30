import { defineConfig } from 'cypress';
import installLogsPrinter from 'cypress-terminal-report/src/installLogsPrinter';

const BASE_URL = `http://localhost:${+process.env.FRONTEND_PORT!}`;
const API_BASE_URL = `http://${process.env.PUBLIC_EXPOSED_API_HOST}:${+process.env
	.PUBLIC_API_PORT!}`;
const FIRST_GAME_ID = 'af8e4754-1b84-4fec-bec4-154a3f894b8f';

export default defineConfig({
	env: {
		baseUrl: `${BASE_URL}/`,
		profileUrl: `${BASE_URL}/profile`,
		firstLevelPageUrl: `${BASE_URL}/game/${FIRST_GAME_ID}`,
		leaderBoardPageUrl: `${BASE_URL}/leaderboard`,
		loginPageUrl: `${BASE_URL}/login`,
		registerPageUrl: `${BASE_URL}/register`,
		apiFirstLevelPageUrl: `${API_BASE_URL}/game/${FIRST_GAME_ID}`
	},

	e2e: {
		setupNodeEvents(on, config) {
			installLogsPrinter(on);
		}
	}
});
