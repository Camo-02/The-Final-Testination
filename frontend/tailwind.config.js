/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			colors: {
				theme: '#CF19E0', //bg pink
				primary: '#E41BF6', //primary pink
				'game-green': '#01DC17',
				'game-gray': '#414141',
				'platform-black': '#360A14',
				'platform-white': '#F6E8EA',
				'game-red': '#DF2B40',
				'game-cyan': '#62B4FF',
				'game-blue': '#155079'
			},
			fontFamily: {
				platform: ['Montserrat', 'sans-serif'],
				game: ['JetBrains Mono', 'monospace']
			}
		}
	},
	plugins: []
};
