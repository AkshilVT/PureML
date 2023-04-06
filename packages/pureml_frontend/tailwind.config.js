/** @type {import('tailwindcss').Config} */

const { withTV } = require("tailwind-variants/transformer");

module.exports = withTV({
  daisyui: {
    themes: [
      {
        mytheme: {
          primary: "#191F4D",
          secondary: "#ffffff",
          error: "#ef4444",
          "base-100": "#ffffff",
          "--rounded-btn": "0.25rem",
          "--rounded-btn-secondary": "0.25rem",
          "--btn-text-case": "titlecase",
        },
      },
    ],
  },
  content: ["./app/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        brand: {
          200: "#191f4d",
        },
        blue: {
          150: "#2196f3",
          250: "#0E4DDD",
        },
        slate: {
          0: "#ffffff",
          50: "#f8fafc",
          100: "#f1f5f9",
          150: "#eaeef3",
          200: "#e2e8f0",
          300: "#cbd5e1",
          400: "#94a3b8",
          500: "#64748b",
          600: "#475569",
          700: "#334155",
          800: "#1e293b",
          900: "#0f172a",
          950: "#020617",
        },
        orange: {
          50: "#fff3e0",
          100: "#fce4c0",
          200: "#ffce84",
        },
      },
    },
    fontSize: {
      base: "14px",
      lg: "16px",
      xl: "18px",
      "2xl": "20px",
      "3xl": "24px",
      "4xl": "32px",
      "5xl": "48px",
      "6xl": "64px",
    },
    screens: {
      sm: "640px",
      md: "768px",
      lg: "1024px",
      xl: "1280px",
      "2xl": "2560px",
      "4xl": "3100px",
    },
    dropShadow: {
      "3xl": "0 0 30px rgba(203, 213, 225, 0.05)",
    },
  },
  plugins: [
    require("@tailwindcss/forms"),
    require("daisyui"),
    require("tailwindcss-radix")(),
  ],
});
