import { registerCommand } from 'cypress-wait-for-stable-dom';
registerCommand();

Cypress.config('defaultCommandTimeout', 6000);

const MAX_PAGE_LEN = 25;

describe('Leaderboard Page', () => {
	beforeEach(() => {
		cy.visit(Cypress.env('leaderBoardPageUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 1000 });
	});

	it('when I press the go forward button the next page of the leaderboard is displayed', () => {
		cy.get('#contestants_list').children().should('have.length', MAX_PAGE_LEN);
		const go_forward_button = cy.get('#go_forward_button');
		go_forward_button.click();
		cy.get('#contestants_list').children().should('have.length.at.most', MAX_PAGE_LEN);
	});

	it('when I press the go back button the previous page of the leaderboard is displayed', () => {
		const go_forward_button = cy.get('#go_forward_button');
		go_forward_button.click();
		cy.get('#contestants_list').children().should('have.length.at.most', MAX_PAGE_LEN);
		const go_back_button = cy.get('#go_back_button');
		go_back_button.click();
		cy.get('#contestants_list').children().should('have.length', MAX_PAGE_LEN);
	});

	it("when I press the go back button when on the first page of the leaderboard I can't go back anymore", () => {
		cy.get('#go_back_button').should('be.disabled');
	});

	it("when I press the go forward button when on the last page of the leaderboard is displayed I can't go forward anymore", () => {
		const go_forward_button = cy.get('#go_forward_button');
		go_forward_button.click();
		cy.get('#go_forward_button').should('be.disabled');
	});
});
