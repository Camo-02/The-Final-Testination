import { BASE_API_URL } from '$src/constants';
import type { PageServerLoad } from './$types';

export interface LeaderboardEntry {
	username: string;
	score: number;
}

export interface Leaderboard {
	currentPage: number;
	pages: number;
	entries: LeaderboardEntry[];
}

export const load: PageServerLoad = async ({ url }) => {
	let page = parseInt(url.searchParams.get('page') ?? '1');

	if (Number.isNaN(page) || page < 1) {
		page = 1;
	}

	const res = await fetch(`${BASE_API_URL}/leaderboard/${page}`);

	return (await res.json()) as Leaderboard;
};
