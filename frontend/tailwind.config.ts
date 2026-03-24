import type { Config } from 'tailwindcss'

export default {
  content: [
    './app/**/*.{vue,ts,tsx}',
    './app/pages/**/*.vue',
    './app/components/**/*.vue',
    './app/layouts/**/*.vue',
  ],
  theme: {
    extend: {
      colors: {
        'r-bg':        '#08080d',
        'r-surface':   '#0f0f17',
        'r-card':      '#15151f',
        'r-border':    '#22223a',
        'r-accent':    '#e63946', // racing red
        'r-yellow':    '#ffd166', // caution / warning
        'r-green':     '#06d6a0', // safe / go
        'r-blue':      '#4cc9f0', // info
        'r-purple':    '#c084fc', // personal best / fastest lap (vivid violet)
        'r-gold':      '#fbbf24', // P1
        'r-silver':    '#94a3b8', // P2
        'r-bronze':    '#f97316', // P3
        'r-text':      '#e8e8f0',
        'r-muted':     '#7070a0',
        'r-dim':       '#3a3a55',
      },
      fontFamily: {
        mono:    ['JetBrains Mono', 'Consolas', 'monospace'],
        display: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'pulse-fast': 'pulse 0.6s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'fade-in':    'fadeIn 0.3s ease-in-out',
        'slide-up':   'slideUp 0.3s ease-out',
      },
      keyframes: {
        fadeIn:  { from: { opacity: '0' }, to: { opacity: '1' } },
        slideUp: { from: { transform: 'translateY(10px)', opacity: '0' }, to: { transform: 'translateY(0)', opacity: '1' } },
      },
    },
  },
  plugins: [],
} satisfies Config
