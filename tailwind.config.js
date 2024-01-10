/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./internal/templates/**/*.templ'],
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

