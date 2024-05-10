/** @type {import('tailwindcss').Config} */
import type { Config } from "tailwindcss";

import { colors } from "./ui/theme/colors";
import { shadows } from "./ui/theme/shadows";

export default {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./ui/**/*.{js,ts,jsx,tsx,mdx}",
    "./styles/**/*.{css,scss}",
  ],
  theme: {
    fontFamily: {
      heading: ["var(--font-barlow)"],
      body: ["var(--font-barlow)"],
      sticky: ["var(--font-merriweather)"],
    },
    boxShadow: shadows,
    colors: {
      ...colors,
      transparent: "transparent",
      current: "currentColor",
    },
    extend: {
      keyframes: {
        pulseOpacity: {
          from: { opacity: "0.3" },
          to: { opacity: "0.7" },
        },
        wave: {
          "0%": { transform: "rotate(0deg)" },
          "100%": { transform: "rotate(360deg)" },
        },
        slideDownAndFade: {
          from: { opacity: "0", transform: "translateY(-2px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
        slideLeftAndFade: {
          from: { opacity: "0", transform: "translateX(2px)" },
          to: { opacity: "1", transform: "translateX(0)" },
        },
        slideUpAndFade: {
          from: { opacity: "0", transform: "translateY(2px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
        slideRightAndFade: {
          from: { opacity: "0", transform: "translateX(-2px)" },
          to: { opacity: "1", transform: "translateX(0)" },
        },
        overlayShow: {
          from: { opacity: "0" },
          to: { opacity: "1" },
        },
        contentShowTop: {
          from: {
            opacity: "0",
            transform: "translate(-50%, -4%) scale(0.96)",
          },
          to: { opacity: "1", transform: "translate(-50%, 0%) scale(1)" },
        },
        contentShowCenter: {
          from: {
            opacity: "0",
            transform: "translate(-50%, -48%) scale(0.96)",
          },
          to: { opacity: "1", transform: "translate(-50%, -50%) scale(1)" },
        },
        collapseDown: {
          from: {
            height: 0,
            opacity: 0,
          },
          to: {
            height: "var(--radix-collapsible-content-height)",
            opacity: 1,
          },
        },
        collapseUp: {
          from: {
            height: "var(--radix-collapsible-content-height)",
            opacity: "1",
          },
          to: {
            height: "0",
            opacity: "0",
          },
        },
        slideDown: {
          from: {
            opacity: "0",
            transform: "translateY(-10px)",
          },
          to: {
            opacity: "1",
            transform: "translateY(0)",
          },
        },
        slideUp: {
          from: {
            opacity: "0",
            transform: "translateY(10px)",
          },
          to: {
            opacity: "1",
            transform: "translateY(0)",
          },
        },
      },

      animation: {
        slideDownAndFade:
          "slideDownAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)",
        slideLeftAndFade:
          "slideLeftAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)",
        slideUpAndFade: "slideUpAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)",
        slideRightAndFade:
          "slideRightAndFade 400ms cubic-bezier(0.16, 1, 0.3, 1)",
        overlayShow: "overlayShow 150ms cubic-bezier(0.16, 1, 0.3, 1)",
        contentShowTop: "contentShowTop 150ms cubic-bezier(0.16, 1, 0.3, 1)",
        contentShowCenter:
          "contentShowCenter 150ms cubic-bezier(0.16, 1, 0.3, 1)",
        collapseDown: "collapseDown 400ms cubic-bezier(0.16, 1, 0.3, 1)",
        collapseUp: "collapseUp 400ms cubic-bezier(0.16, 1, 0.3, 1)",
        pulseOpacity: "pulseOpacity 0.7s infinite alternate ease-in-out",
        slideUp: "slideUp 400ms cubic-bezier(0.16, 1, 0.3, 1)",
        slideDown: "slideDown 400ms cubic-bezier(0.16, 1, 0.3, 1)",
      },
      flex: {
        2: "2 2 0%",
        3: "3 3 0%",
        4: "4 4 0%",
        5: "5 5 0%",
        6: "6 6 0%",
        7: "7 7 0%",
        8: "8 8 0%",
        9: "9 9 0%",
      },
    },
  },
  plugins: [],
} satisfies Config;
