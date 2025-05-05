/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#ec4c21", // Orange-red
        secondary: "#322f2e", // Dark gray
        background: "#fefdfa", // Off-white
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        tastygo: {
          primary: "#ec4c21",
          secondary: "#322f2e",
          accent: "#ec4c21",
          neutral: "#322f2e",
          "base-100": "#fefdfa",
          info: "#3abff8",
          success: "#36d399",
          warning: "#fbbd23",
          error: "#f87272",
        },
      },
    ],
  },
};
