/// <reference types="cypress" />

import { registerCommand } from 'cypress-wait-for-stable-dom';
registerCommand();

Cypress.config('defaultCommandTimeout', 6000);

describe('Redirect if access without login', () => {
	it('should redirect to the login page from game', () => {
		cy.visit(Cypress.env('firstLevelPageUrl'));
		cy.url().should('eq', Cypress.env('loginPageUrl'));
	});

	it('should redirect to the login page from landing page', () => {
		cy.visit(Cypress.env('baseUrl'));
		cy.url().should('eq', Cypress.env('loginPageUrl'));
	});
	it('should redirect to the login page from profile page', () => {
		cy.visit(Cypress.env('profileUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.url().should('eq', Cypress.env('loginPageUrl'));
	});
});
