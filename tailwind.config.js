/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./views/**/*.templ"],
    theme: {
        extend: {},
    },
    plugins: [],
    // safelist: [{ pattern: /(border|text|bg)-(red|green|blue)-(100|200|900)/ }],
    safelist: ["hidden"],
};
