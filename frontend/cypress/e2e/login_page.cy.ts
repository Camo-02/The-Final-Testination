/// <reference types="cypress" />

import { registerCommand } from 'cypress-wait-for-stable-dom';
registerCommand();

Cypress.config('defaultCommandTimeout', 6000);

async function login(credential: string, password: string) {
	cy.get('#credential').type(credential);
	cy.get('#password').type(password);
	cy.get('#loginBtn').click();
}

describe('Login ', () => {
	beforeEach(() => {
		cy.visit(Cypress.env('loginPageUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
	});

	it('A user can login with username', () => {
		login('admin', 'rootroot');
		cy.url().should('eq', Cypress.env('baseUrl'));
	});

	it('A user can login with email', () => {
		login('admin@testination.com', 'rootroot');
		cy.url().should('eq', Cypress.env('baseUrl'));
	});

	it("A user can't login with wrong credentials", () => {
		login('admin', 'wrongpassword');
		cy.on('window:alert', (text) => {
			expect(text).to.equal('wrong credentials');
		});
		cy.url().should('eq', Cypress.env('loginPageUrl'));
	});
});
