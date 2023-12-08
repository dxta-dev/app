/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./internals/templates/**/*.templ'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-geist-sans)'],
        mono: ['var(--font-geist-mono)'],
      },
    },
  },
  plugins: [],
}

