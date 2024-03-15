/** @type {import('tailwindcss').Config} */

import { colors } from './ui/theme/colors';
import { shadows } from './ui/theme/shadows';

module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './ui/**/*.{js,ts,jsx,tsx,mdx}',
    './styles/**/*.{css,scss}',
  ],
  theme: {
    boxShadow: shadows,
    colors,
    extend: {},
  },
  plugins: [require('tailwindcss-animate')],
};
