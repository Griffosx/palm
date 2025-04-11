import {
  darkAquamarine,
  orange,
  peach,
} from "../../../styles/themes";

const colors = [darkAquamarine, orange, peach];

// Function to generate a hash from a string
const stringToHash = (str: string): number => {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = (hash << 5) - hash + str.charCodeAt(i);
    hash |= 0; // Convert to 32bit integer
  }
  return Math.abs(hash); // Ensure positive number
};

// Deterministically select a background color based on a string (e.g., email address)
export const getDeterministicColor = (inputString: string | null | undefined): string => {
  if (!inputString) {
    // Fallback if inputString is somehow missing
    return colors[0];
  }
  const hash = stringToHash(inputString);
  const colorIndex = hash % colors.length;
  return colors[colorIndex];
};