/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./internal/template/**/*.templ'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-geist-sans)'],
        mono: ['var(--font-geist-mono)'],
      },
      keyframes: {
        flash: {
          '0%, 100%': { borderColor: 'red' },
          '50%': { borderColor: 'transparent' },
        }
      },
      animation: {
        'flash-border': 'flash 1s infinite',
      },
    },
  },
  plugins: [],
}
