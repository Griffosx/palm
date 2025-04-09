import React from "react";
import { gradients, noiseOpacity, solidColors } from "../styles/themes";

interface BackgroundOverlayProps {
  gradient?: keyof typeof gradients;
  customGradient?: string;
  solidColor?: keyof typeof solidColors | string;
  noiseOpacity?: number;
  className?: string;
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
}) => {
  // Determine the background style
  let backgroundStyle = "";

  if (solidColor) {
    // If it's a known color from our theme
    if (solidColor in solidColors) {
      backgroundStyle = solidColors[solidColor as keyof typeof solidColors];
    } else {
      // Otherwise use it as a custom color
      backgroundStyle = solidColor;
    }
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

  // Use custom opacity if provided, otherwise use the default from themes
  const opacity =
    customNoiseOpacity !== undefined ? customNoiseOpacity : noiseOpacity;

  return (
    <div
      className={`absolute inset-0 w-full h-full ${className}`}
      style={{
        background: backgroundStyle,
      }}
    >
      {/* Noise overlay */}
      <div
        className="absolute inset-0 w-full h-full"
        style={{
          backgroundImage:
            "url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noiseFilter%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.65%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url%28%23noiseFilter%29%22/%3E%3C/svg%3E')",
          backgroundPosition: "0 0",
          backgroundSize: "200px 200px",
          opacity,
        }}
      />
    </div>
  );
};

export default BackgroundOverlay;
