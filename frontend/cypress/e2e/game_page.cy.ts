/// <reference types="cypress" />

import { registerCommand } from 'cypress-wait-for-stable-dom';
registerCommand();

Cypress.config('defaultCommandTimeout', 6000);

function dndSelector(selector: string, target: string) {
	const dataTransfer = new DataTransfer();

	cy.get(selector).trigger('dragstart', {
		dataTransfer
	});

	cy.get(target).trigger('drop', {
		dataTransfer
	});
}

function dndElement(source: Cypress.Chainable, target: string) {
	const dataTransfer = new DataTransfer();

	source.trigger('dragstart', {
		dataTransfer
	});

	cy.get(target).trigger('drop', {
		dataTransfer
	});
}

async function login(credential: string, password: string) {
	cy.get('#credential').type(credential);
	cy.get('#password').type(password);
	cy.get('#loginBtn').click();
}

let body: Record<string, unknown> = {};

describe('Game page Site', () => {
	beforeEach(() => {
		cy.visit(Cypress.env('baseUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 10000 });
		login('test', 'rootroot');
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 10000 });

		cy.intercept('GET', Cypress.env('apiFirstLevelPageUrl'), (req) => {
			req.continue((res) => {
				body = res.body;
			});
		});
		cy.visit(Cypress.env('firstLevelPageUrl'));
		cy.waitForStableDOM({ pollInterval: 1000, timeout: 10000 });
	});

	it('I can see the all the blocks', () => {
		const grid = cy.get('#block-list');
		const grid_blocks = grid.children();
		grid_blocks.should('have.length', 11);
		grid_blocks.each(($blk, index) => {
			expect($blk.text()).to.equal((body.blocks as string[])[index]);
		});
	});

	it('I can return to the landing page using the quit button', () => {
		cy.get('#overlayMenuBtn').click();
		cy.get('#quitBtn').click();
		cy.url().should('eq', Cypress.env('baseUrl'));
	});

	it('I can see the game with a correct layout', () => {
		cy.get('#input-area')
			.children()
			.should('be.visible')
			.should('have.length', 7)
			.each(($blk, index) => {
				if ($blk.text().trim() !== 'fill') {
					expect($blk.text()).to.equal((body.skeleton as Record<string, string>)[index.toString()]);
				}
			});
	});

	it('I can see the story', () => {
		const story = cy.get('#story');
		story.click();
		cy.get('#sidebar').should('exist').should('contain', 'EMAIL');
	});

	it('I can see the cheatsheet', () => {
		const cheatsheet = cy.get('#cheatsheet');
		cheatsheet.click();
		cy.get('#sidebar').should('exist').should('contain', 'Introduction');
	});

	it('I can see blocks correctly displayed in their area', () => {
		const BLOCK_IN_LIST_SELECTOR = '#block-list > :nth-child(1)';
		const BLOCK_IN_BLOCKS_SELECTOR = '#input-area > :nth-child(1)';

		cy.get(BLOCK_IN_LIST_SELECTOR)
			.invoke('text')
			.then((value) => {
				dndSelector(BLOCK_IN_LIST_SELECTOR, BLOCK_IN_BLOCKS_SELECTOR);

				// the set of possibilities shouldn't contain the block anymore
				cy.get(BLOCK_IN_LIST_SELECTOR).should('not.contain.text', value);

				// but the input area should
				cy.get(BLOCK_IN_BLOCKS_SELECTOR).should('contain.text', value);
			});
	});

	it("I can't drop a block in a non-droppable position", () => {
		const BLOCK_IN_LIST_SELECTOR = '#block-list > :nth-child(1)';
		const BLOCK_IN_BLOCKS_SELECTOR = '#input-area > :nth-child(2)';

		cy.get(BLOCK_IN_LIST_SELECTOR)
			.invoke('text')
			.then((value) => {
				dndSelector(BLOCK_IN_LIST_SELECTOR, BLOCK_IN_BLOCKS_SELECTOR);

				cy.get('#block-list').should('contain.text', value);
				cy.get(BLOCK_IN_BLOCKS_SELECTOR).should('not.contain.text', value);
			});
	});

	it('I can move a block in a position to fill', () => {
		const BLOCK_IN_LIST_SELECTOR = '#block-list > :nth-child(1)';
		const BLOCK_IN_BLOCKS_SELECTOR = '#input-area > :nth-child(1)';
		const THIRD_BLOCK_IN_BLOCKS_SELECTOR = '#input-area > :nth-child(3)';

		cy.get(BLOCK_IN_LIST_SELECTOR)
			.invoke('text')
			.then((value) => {
				dndSelector(BLOCK_IN_LIST_SELECTOR, BLOCK_IN_BLOCKS_SELECTOR);

				cy.get(BLOCK_IN_LIST_SELECTOR).should('not.contain.text', value);
				cy.get(BLOCK_IN_BLOCKS_SELECTOR).should('contain.text', value);

				dndSelector(BLOCK_IN_BLOCKS_SELECTOR, THIRD_BLOCK_IN_BLOCKS_SELECTOR);

				cy.get(BLOCK_IN_BLOCKS_SELECTOR).should('not.contain.text', value);
				cy.get(BLOCK_IN_BLOCKS_SELECTOR).should('contain.text', 'fill');
				cy.get(THIRD_BLOCK_IN_BLOCKS_SELECTOR).should('contain.text', value);
			});
	});

	it("I can't move a block in a already filled position", () => {
		const BLOCK_IN_LIST_SELECTOR = '#block-list > :nth-child(1)';
		const BLOCK_IN_BLOCKS_SELECTOR = '#input-area > :nth-child(1)';
		const THIRD_BLOCK_IN_BLOCKS_SELECTOR = '#input-area > :nth-child(3)';

		cy.get(BLOCK_IN_LIST_SELECTOR)
			.invoke('text')
			.then((value) => {
				dndSelector(BLOCK_IN_LIST_SELECTOR, BLOCK_IN_BLOCKS_SELECTOR);

				cy.get(BLOCK_IN_LIST_SELECTOR).should('not.contain.text', value);
				cy.get(BLOCK_IN_BLOCKS_SELECTOR).should('contain.text', value);

				dndSelector(BLOCK_IN_LIST_SELECTOR, THIRD_BLOCK_IN_BLOCKS_SELECTOR);

				dndSelector(BLOCK_IN_BLOCKS_SELECTOR, THIRD_BLOCK_IN_BLOCKS_SELECTOR);

				cy.get(BLOCK_IN_BLOCKS_SELECTOR).should('contain.text', value);
				cy.get(THIRD_BLOCK_IN_BLOCKS_SELECTOR).should('not.contain.text', value);
			});
	});

	it('I can remove blocks correctly with the button', () => {
		cy.get('#block-list > :nth-child(1)')
			.invoke('text')
			.then((value) => {
				// move the block to the input area
				dndSelector('#block-list > :nth-child(1)', '#input-area > :nth-child(1)');

				// then remove it using the button
				cy.get('#input-area > :nth-child(1) > .opacity-0').click();

				cy.get('#block-list').should('contain.text', value);
				cy.get('#input-area > :nth-child(1) > .text-black').should('not.contain.text', value);
				cy.get('#block-list').should('contain.text', value);
			});
	});

	it('I can drag and drop blocks correctly', () => {
		cy.get('#block-list > :nth-child(1)')
			.invoke('text')
			.then((value) => {
				// move the block to the input area
				dndSelector('#block-list > :nth-child(1)', '#input-area > :nth-child(1)');

				// then remove it using dnd
				dndSelector('#input-area > :nth-child(1)', '#block-list');

				cy.get('#block-list').should('contain.text', value);
				cy.get('#input-area > :nth-child(1) > .text-black').should('not.contain.text', value);
			});
	});

	it('I can see an alert saying that some blocks are missing', () => {
		cy.get('#inject-button').should('be.visible').click();
		cy.on('window:alert', (text) => {
			expect(text).to.equal('Please fill all the boxes');
		});
	});

	it('I can see that I won', () => {
		dndElement(cy.contains('button', '<iframe'), '#input-area > :nth-child(1)');
		dndElement(cy.contains('button', 'evilcompany'), '#input-area > :nth-child(3)');
		dndElement(cy.contains('button', '100%'), '#input-area > :nth-child(5)');
		dndElement(cy.contains('button', '></iframe'), '#input-area > :nth-child(7)');

		cy.get('#inject-button').should('be.visible').click();
		cy.get('#game-dialog').should('be.visible').contains('Success!');
	});

	it('I can open and close the the game dialog', () => {
		dndElement(cy.contains('button', '<iframe'), '#input-area > :nth-child(1)');
		dndElement(cy.contains('button', 'evilcompany'), '#input-area > :nth-child(3)');
		dndElement(cy.contains('button', '100%'), '#input-area > :nth-child(5)');
		dndElement(cy.contains('button', '></iframe'), '#input-area > :nth-child(7)');

		cy.get('#open-game-dialog-btn').should('not.exist');

		cy.get('#inject-button').should('be.visible').click();
		cy.get('#game-dialog').should('be.visible').contains('Success!');

		cy.get('#game-dialog > #close-button').should('be.visible').click();
		cy.get('#open-game-dialog-btn').should('be.visible').contains('see result').click();
		cy.get('#game-dialog').should('be.visible').contains('Success!');
	});

	it('I can see that I lost', () => {
		dndElement(cy.contains('button', '<iframe'), '#input-area > :nth-child(1)');
		dndElement(cy.contains('button', '100%'), '#input-area > :nth-child(3)');
		dndElement(cy.contains('button', 'evilcompany'), '#input-area > :nth-child(5)');
		dndElement(cy.contains('button', '></iframe'), '#input-area > :nth-child(7)');

		cy.get('#inject-button').should('be.visible').click();
		cy.get('#game-dialog').should('be.visible').contains('You lose!');
	});

	it('I can return to the home page after I won', () => {
		dndElement(cy.contains('button', '<iframe'), '#input-area > :nth-child(1)');
		dndElement(cy.contains('button', 'evilcompany'), '#input-area > :nth-child(3)');
		dndElement(cy.contains('button', '100%'), '#input-area > :nth-child(5)');
		dndElement(cy.contains('button', '></iframe'), '#input-area > :nth-child(7)');

		cy.get('#inject-button').should('be.visible').click();

		cy.get('#game-dialog').should('be.visible').contains('Success!');
		cy.contains('a', 'exit the game').click();
		cy.url().should('eq', Cypress.env('baseUrl'));
	});

	it('I can go to the next level after I won', () => {
		dndElement(cy.contains('button', '<iframe'), '#input-area > :nth-child(1)');
		dndElement(cy.contains('button', 'evilcompany'), '#input-area > :nth-child(3)');
		dndElement(cy.contains('button', '100%'), '#input-area > :nth-child(5)');
		dndElement(cy.contains('button', '></iframe'), '#input-area > :nth-child(7)');

		cy.get('#inject-button').should('be.visible').click();

		cy.get('#game-dialog').should('be.visible').contains('Success!');
		cy.contains('a', 'next level').click();
		cy.url().should('eq', Cypress.env('baseUrl') + 'game/05732286-9fa5-45d4-bef3-13ae0d481afa');
	});

	it('I can return to the home page after I lose', () => {
		dndElement(cy.contains('button', '<iframe'), '#input-area > :nth-child(1)');
		dndElement(cy.contains('button', '100%'), '#input-area > :nth-child(3)');
		dndElement(cy.contains('button', 'evilcompany'), '#input-area > :nth-child(5)');
		dndElement(cy.contains('button', '></iframe'), '#input-area > :nth-child(7)');

		cy.get('#inject-button').should('be.visible').click();

		cy.get('#game-dialog').should('be.visible').contains('You lose!');
		cy.contains('a', 'exit the game').click();
		cy.url().should('eq', Cypress.env('baseUrl'));
	});

	it('I can try again the game after I lost', () => {
		dndElement(cy.contains('button', '<iframe'), '#input-area > :nth-child(1)');
		dndElement(cy.contains('button', '100%'), '#input-area > :nth-child(3)');
		dndElement(cy.contains('button', 'evilcompany'), '#input-area > :nth-child(5)');
		dndElement(cy.contains('button', '></iframe'), '#input-area > :nth-child(7)');

		cy.get('#inject-button').should('be.visible').click();

		cy.get('#game-dialog').should('be.visible').contains('You lose!');
		cy.contains('button', 'try again').click();
		cy.url().should('eq', Cypress.env('firstLevelPageUrl'));
	});

	it('I can see the how to play in the game dialog', () => {
		cy.get('#overlayMenuBtn').should('be.visible').click();

		cy.get('#tutorialBtn').should('be.visible').click();
		cy.get('#game-dialog').should('be.visible').contains('How to play');
	});

	it('I can see the hints being displayed', () => {
		const hints = cy.get('#hints');
		hints.click();
		cy.get('h1').contains('Coins').should('be.visible');
		cy.get('button').contains('Fill Block').should('be.visible');
		cy.get('button').contains('Freeze time').should('be.visible');
		cy.get('button').contains('Textual hint').should('be.visible');
	});

	it('I can see the timer being displayed', () => {
		const menu_button = cy.get('#overlayMenuBtn');
		menu_button.click();
		const timer = cy.get('#timer');
		timer.should('be.visible');
	});

	it('I can see the snowflake icon when the time is frozen', () => {
		const hints = cy.get('#hints');
		hints.click();
		cy.get('button').contains('Freeze time').click();
		const score = cy.get('#timeFreezeAd');
		score.should('be.visible');
		const menu_button = cy.get('#overlayMenuBtn');
		menu_button.click();
		const snowflake_r = cy.get('#snowflake_side_r');
		snowflake_r.should('be.visible');
		const snowflake_l = cy.get('#snowflake_side_l');
		snowflake_l.should('be.visible');
	});

	it('I can see the first block filled by the hint', () => {
		const hints = cy.get('#hints');
		hints.click();
		cy.get('button').contains('Fill Block').click();
		cy.get('#input-area > div:nth-child(1) > button').click();
		cy.get('#block-list').should('not.contain', '<iframe');
	});

	it('I can see the textual hint being displayed', () => {
		const hints = cy.get('#hints');
		hints.click();
		cy.get('button').contains('Textual hint').click();
		const score = cy.get('#textualHintAd');
		score.should('be.visible');
	});
});
