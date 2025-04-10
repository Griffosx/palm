import React from "react";
import { noiseOpacity as defaultNoiseOpacity } from "../styles/themes";

interface NoiseOverlayProps {
  opacity?: number;
  className?: string;
}

const NoiseOverlay: React.FC<NoiseOverlayProps> = ({
  opacity,
  className = "",
}) => {
  // Use custom opacity if provided, otherwise use the default from themes
  const finalOpacity = opacity !== undefined ? opacity : defaultNoiseOpacity;

  return (
    <div
      className={`absolute inset-0 w-full h-full ${className}`}
      style={{
        backgroundImage:
          "url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noiseFilter%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.65%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url%28%23noiseFilter%29%22/%3E%3C/svg%3E')",
        backgroundPosition: "0 0",
        backgroundSize: "200px 200px",
        opacity: finalOpacity,
      }}
    />
  );
};

export default NoiseOverlay;
