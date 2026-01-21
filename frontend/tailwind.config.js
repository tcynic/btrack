/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        ld: {
          bg: 'var(--ld-bg)',
          surface: 'var(--ld-surface)',
          surface2: 'var(--ld-surface-2)',
          border: 'var(--ld-border)',
          text: 'var(--ld-text)',
          muted: 'var(--ld-text-muted)',
          primary: 'var(--ld-primary)',
          cyan: 'var(--ld-cyan)',
          blue: 'var(--ld-blue)',
          powder: 'var(--ld-powder)',
          purple: 'var(--ld-purple)',
          pink: 'var(--ld-pink)',
          orange: 'var(--ld-orange)',
          yellow: 'var(--ld-yellow)',
          green: 'var(--ld-green)'
        }
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'Avenir', 'Helvetica', 'Arial', 'sans-serif']
      }
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
