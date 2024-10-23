// tailwind.config.js

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,jsx,ts,tsx}', './public/index.html'],
  theme: {
    extend: {},
  },
  variants: {
    extend: {},
  },
  plugins: [],
  theme: {
    extend: {
      fontSize: {
        base: '0.875rem',
    
      },
      maxWidth: {
        '5xl': '64rem', 
      },
    },
  },
}

