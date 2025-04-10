import React from "react";
import { gradients } from "../styles/themes";
import NoiseOverlay from "./NoiseOverlay";

interface BackgroundOverlayProps {
  gradient?: keyof typeof gradients;
  customGradient?: string;
  solidColor?: string;
  noiseOpacity?: number;
  className?: string;
  style?: React.CSSProperties;
  rounded?: string;
  position?: "absolute" | "fixed" | "relative";
  zIndex?: number;
}

/**
 * A reusable component for background gradients or solid colors with noise overlay
 */
const BackgroundOverlay: React.FC<BackgroundOverlayProps> = ({
  gradient,
  customGradient,
  solidColor,
  noiseOpacity: customNoiseOpacity,
  className = "",
  style = {},
  rounded = "",
  position = "absolute",
  zIndex = 0,
}) => {
  // Determine the background style
  let backgroundStyle = "";

  if (solidColor) {
    backgroundStyle = solidColor;
  } else if (customGradient) {
    backgroundStyle = customGradient;
  } else if (gradient) {
    backgroundStyle = `linear-gradient(${
      gradients[gradient].direction
    }, ${gradients[gradient].colors.join(", ")})`;
  } else {
    // Default to main gradient
    backgroundStyle = `linear-gradient(${
      gradients.main.direction
    }, ${gradients.main.colors.join(", ")})`;
  }

  // Combine styles
  const combinedStyles = {
    ...style,
    background: backgroundStyle,
    position,
    zIndex,
  };

  // Combine class names
  const combinedClassNames = `inset-0 w-full h-full ${rounded} ${className}`;

  return (
    <div className={combinedClassNames} style={combinedStyles}>
      {/* Noise overlay */}
      <NoiseOverlay opacity={customNoiseOpacity} />
    </div>
  );
};

export default BackgroundOverlay;
