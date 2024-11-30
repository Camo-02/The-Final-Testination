/// <reference types="cypress" />

import { registerCommand } from 'cypress-wait-for-stable-dom';
registerCommand();

Cypress.config('defaultCommandTimeout', 6000);

async function login(credential: string, password: string) {
	cy.get('#credential').type(credential);
	cy.get('#password').type(password);
	cy.get('#loginBtn').click();
}

describe('Profile page', () => {
	beforeEach(() => {
		cy.visit(Cypress.env('baseUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
	});

	it('Player info should be visible', () => {
		login('admin', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.visit(Cypress.env('profileUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.url().should('eq', Cypress.env('profileUrl'));
		cy.get('#username').should('be.visible');
		cy.get('#icon').should('be.visible');
		cy.get('#coins').should('contain.text', '590 Coins');
	});

	it('A user that won at least a level should see its results', () => {
		login('admin', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.visit(Cypress.env('profileUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.url().should('eq', Cypress.env('profileUrl'));
		// Level title
		cy.get('#level-title-0').should('contain.text', 'XSS Attack');
		// Level Score
		cy.get('#level-score-0').should('contain.text', '100');

		// Start time
		cy.get('#level-started-0').then((el) => {
			expect(new Date(el.text()).toLocaleString()).to.eq(
				new Date(Date.UTC(2024, 1, 24, 10)).toLocaleString()
			);
		});
		// End time
		cy.get('#level-end-0').then((el) => {
			expect(new Date(el.text()).toLocaleString()).to.eq(
				new Date(Date.UTC(2024, 1, 24, 10, 1)).toLocaleString()
			);
		});
		// Freeze hint cost
		cy.get('#level-freeze-0').should('contain.text', '0');
	});

	it("A user that didn't play any level should see a message", () => {
		login('empty', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.visit(Cypress.env('profileUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.url().should('eq', Cypress.env('profileUrl'));
		cy.get('#no-levels-message')
			.should('be.visible')
			.should('contain.text', 'No levels completed yet');
	});

	it("A user that didn't complete a level should see that it's attempt is started and hints haven been used", () => {
		login('profilePageTest', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.visit(Cypress.env('profileUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.url().should('eq', Cypress.env('profileUrl'));

		// Time
		cy.get('#level-started-0').should(($el) => {
			const t = $el.text().toString();
			expect(t.endsWith(new Date(Date.UTC(2024, 1, 24, 10)).toLocaleString())).to.be.true;
		});

		// There should be no end time
		cy.get('#level-end-0').should('not.exist');

		// Freeze hint cost
		cy.get('#level-freeze-0').invoke('text').should('be.equal', '20');
	});

	it('A user should see available icons', () => {
		login('admin', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.visit(Cypress.env('profileUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.url().should('eq', Cypress.env('profileUrl'));
		cy.get('#selectIconButton').click();
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.get('#availableIcons').should('be.visible');
		let availableIcons = cy.get('#availableIcons').find('button');
		availableIcons.should('have.length', 9);
		availableIcons.each((icon, index) => {
			expect(icon).to.have.attr('id', 'icon' + (index + 1));
		});
	});

	it('A user should be able to change its icon', () => {
		login('admin', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.visit(Cypress.env('profileUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
		cy.url().should('eq', Cypress.env('profileUrl'));
		cy.get('#profileIcon')
			.invoke('attr', 'src')
			.then((oldProfPic) => {
				cy.get('#selectIconButton').click();
				cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
				cy.get('#availableIcons').should('be.visible');
				cy.get('#icon2').click();
				cy.get('#saveIconBtn').click();
				cy.get('#profileIcon').should('have.attr', 'src').should('not.equal', oldProfPic);
			});
	});
});
