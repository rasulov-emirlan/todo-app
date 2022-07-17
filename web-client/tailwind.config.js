/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
	theme: {
		extend: {
			colors: {
				darkGray: '#222831',
				mediumGray: '#393E46',
				lightGray: '#545b66',
				darkTeal: '#00666b',
				mediumTeal: '#00ADB5',
				lightWhite: '#EEEEEE',
			}
		},
	},
	plugins: [],
};
