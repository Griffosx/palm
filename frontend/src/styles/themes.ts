// Orange
export const orange = "#EC8234";
// Peach
export const lightPeach = "#F9E2BE";
export const peach = "#F9D59B";
// Aquamarine
export const lightAquamarine = "#E8EFDC";
export const aquamarine = "#CAE7D4";
export const darkAquamarine = "#97CAB5";
// Brown
export const brown = "#BCA482";
export const darkBrown = "#62533C";
// Cream
export const cream = "#FEF6E5";
export const borderCream = "#E2C59C";

export const gradients = {
  main: {
    colors: [lightPeach + " 20%", lightPeach, aquamarine],
    direction: "to bottom",
  },
  alternate: {
    colors: [aquamarine, lightPeach],
    direction: "to right",
  },
};

export const noiseOpacity = 0.25;

export function hexToRgba(hex: string, alpha = 1) {
  // Remove the hash if it exists
  hex = hex.replace("#", "");

  // Handle 3-digit hex codes (like #FFF)
  if (hex.length === 3) {
    hex = hex.split("").map((char: string) => char + char).join("");
  }

  // Parse the hex values to RGB
  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);

  // Return the RGBA string
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}