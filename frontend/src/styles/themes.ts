/**
 * Color themes and gradients for the application
 */

const peachBackground = "#F9E2BE";
const lightBlueBackground = "#CAE7D4";
const creamBackground = "#FEF6E5";

export const gradients = {
  main: {
    colors: [peachBackground + " 20%", peachBackground, lightBlueBackground],
    direction: "to bottom",
  },
  alternate: {
    colors: [lightBlueBackground, peachBackground],
    direction: "to right",
  },
};

export const solidColors = {
  cream: creamBackground,
  peach: peachBackground,
  lightBlue: lightBlueBackground,
};

export const colors = {
  primary: peachBackground,
  secondary: lightBlueBackground,
  cream: creamBackground,
  text: {
    primary: "#333333",
    secondary: "#666666",
    light: "#FFFFFF",
  },
  overlay: {
    light: "rgba(255, 255, 255, 0.1)",
    medium: "rgba(255, 255, 255, 0.3)",
    dark: "rgba(0, 0, 0, 0.1)",
  },
};

export const noiseOpacity = 0.3;