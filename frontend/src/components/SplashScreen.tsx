import { useState, useEffect } from "react";
import { Navigate } from "react-router-dom";

const SplashScreen = () => {
  const [redirect, setRedirect] = useState(false);
  const [blurAmount, setBlurAmount] = useState(0);
  const [textOpacity, setTextOpacity] = useState(1);

  useEffect(() => {
    // Start blur transition after 3 seconds
    const blurTimer = setTimeout(() => {
      setBlurAmount(8); // Target blur amount
      setTextOpacity(0); // Fade out text
    }, 3000);

    // Redirect after blur transition completes (2 additional seconds)
    const redirectTimer = setTimeout(() => {
      setRedirect(true);
    }, 5000);

    return () => {
      clearTimeout(blurTimer);
      clearTimeout(redirectTimer);
    };
  }, []);

  if (redirect) {
    return <Navigate to="/inbox" />;
  }

  return (
    <div className="flex items-center justify-center h-screen w-screen">
      <div
        className="text-center text-white"
        style={{
          opacity: textOpacity,
          transition: "opacity 2s ease-in-out",
        }}
      >
        <h1 className="text-5xl font-bold mb-4">Welcome to Palm</h1>
        <p className="text-xl">Loading your experience...</p>
      </div>
    </div>
  );
};

export default SplashScreen;
