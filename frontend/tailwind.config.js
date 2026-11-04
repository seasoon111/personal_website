/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,ts}'],
  theme: {
    extend: {
      colors: {
        paper: {
          DEFAULT: '#FAF5EC',
          deep: '#F1E7D6',
          dark: '#E8DBC6',
        },
        ink: {
          DEFAULT: '#3B332C',
          soft: '#7A6E62',
          faint: '#A79B8D',
        },
        terracotta: {
          DEFAULT: '#C4633F',
          soft: '#E8B39C',
          deep: '#A44E2E',
        },
        sage: {
          DEFAULT: '#8A9A7B',
          soft: '#C3CDB6',
        },
        line: '#E3D6C2',
      },
      fontFamily: {
        hand: ['"LXGW WenKai"', '"Kaiti SC"', 'KaiTi', 'STKaiti', 'serif'],
        sans: ['system-ui', '-apple-system', '"PingFang SC"', '"Microsoft YaHei"', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'Consolas', 'Menlo', 'monospace'],
      },
      boxShadow: {
        hand: '2px 3px 0 0 rgba(59, 51, 44, 0.08)',
        'hand-lg': '4px 6px 0 0 rgba(59, 51, 44, 0.10)',
      },
      maxWidth: {
        content: '1120px',
        prose: '720px',
      },
    },
  },
  plugins: [require('@tailwindcss/typography')],
}
