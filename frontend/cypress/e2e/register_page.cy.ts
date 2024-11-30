/// <reference types="cypress" />

import { registerCommand } from 'cypress-wait-for-stable-dom';
registerCommand();

Cypress.config('defaultCommandTimeout', 6000);

async function register(credential: string, mail:string,password: string) {
	cy.get('#user').type(credential);
	cy.get('#mail').type(mail);
	cy.get('#password').type(password);
	cy.get('#registerBtn').click();
}

describe('Register ', () => {
	beforeEach(() => {
		cy.visit(Cypress.env('registerPageUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
	});

	it('A user can Register with username,email and pwd', () => {
		register('prova', 'prova@gmail.com','rootroot');
		cy.url().should('eq', Cypress.env('loginPageUrl'));
	});

	
});
