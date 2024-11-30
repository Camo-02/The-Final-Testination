/// <reference types="cypress" />

import { registerCommand } from 'cypress-wait-for-stable-dom';
registerCommand();

Cypress.config('defaultCommandTimeout', 6000);

async function login(credential: string, password: string) {
	cy.get('#credential').type(credential);
	cy.get('#password').type(password);
	cy.get('#loginBtn').click();
}

describe('Landing page Site', () => {
	beforeEach(() => {
		cy.visit(Cypress.env('baseUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		login('test', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
	});

	it('Should have a link to the leaderboard', () => {
		cy.get('#navbar > :nth-child(5)')
			.should('be.visible')
			.should('have.attr', 'href', '/leaderboard');
	});

	it('Should show the available levels', () => {
		cy.get('#gameMap > div').should('be.visible');
		cy.get('#gameMap > div').children('div').should('have.length', 3);
	});

	it('Should show the first level as completed', () => {
		cy.get('#gameMap > div > div:nth-child(2) > img').should('be.visible');
		cy.get('#gameMap > div > div:nth-child(2) > button')
			.should('have.text', '1')
			.should('not.be.disabled')
			.click();
		cy.get('#modal').should('be.visible');
		cy.get('#level-title').should('be.visible').should('have.text', 'XSS Attack');
		cy.get('#level-description').should('be.visible').should('contain.text', 'A conspiracy');
		cy.get('#level-description').should('contain.text', 'Score Achieved: 100 out of 100');
		cy.get('#play-game-btn')
			.should('be.visible')
			.should('not.be.disabled')
			.should('have.text', 'Play')
			.should('have.attr', 'href', '/game/af8e4754-1b84-4fec-bec4-154a3f894b8f');
		cy.get('#close-modal-btn').should('not.be.disabled').click();
		cy.get('#gameMap').should('be.visible');
	});

	it('Should show the second level as playable', () => {
		cy.get('#gameMap > div > div:nth-child(3) > button')
			.should('have.text', '2')
			.should('not.be.disabled')
			.click();

		cy.get('#modal').should('be.visible');

		cy.get('#play-game-btn').should('be.visible');
		cy.get('#close-modal-btn').should('not.be.disabled').click();
		cy.get('#gameMap').should('be.visible');
	});

	it('Should show the third level as locked', () => {
		cy.get('#gameMap > div > div:nth-child(4) > button').should('not.be.disabled').click();

		cy.get('#modal').should('be.visible');

		cy.get('#level-title').should('be.visible').should('have.text', 'Level 3');
		cy.get('#level-description').should('be.visible').should('have.text', 'Locked');

		cy.get('#play-game-btn').should('not.exist');
		cy.get('#gameMap').should('be.visible');
	});
});
