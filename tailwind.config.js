/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.{html,js}",
    "./node_modules/flowbite/**/*.js"
  ],
  plugins: [require("flowbite/plugin")],
  darkMode: 'media',
  theme: {
    extend: {
      colors: {
        primary: {"50":"#ecfeff","100":"#cffafe","200":"#a5f3fc","300":"#67e8f9","400":"#22d3ee","500":"#06b5d4","600":"#0890b2","700":"#0e7490","800":"#155f75","900":"#164e63","950":"#0d3b4d"}
      }
    },
  },
}