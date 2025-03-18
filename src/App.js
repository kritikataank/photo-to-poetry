import React, { useState } from "react";
import CameraCapture from "./components/CameraCapture";

function App() {
  const [capturedImage, setCapturedImage] = useState(null);

  // Function to send image to backend and redirect
  const sendToBackend = async () => {
    if (!capturedImage) {
      alert("Capture an image first!");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ image: capturedImage }),
      });

      // ✅ If response is a redirect, navigate to new page
      if (response.redirected) {
        window.location.href = response.url; // Redirect user to /uploaded
      }
    } catch (error) {
      console.error("❌ Error sending image:", error);
    }
  };

  return (
    <div className="App">
      <h1>Photo-to-Poetry Web App</h1>
      <CameraCapture onCapture={setCapturedImage} />

      {capturedImage && (
        <div>
          <h2>Captured Image:</h2>
          <img src={capturedImage} alt="Captured" />

          {/* ✅ Upload button to send image & redirect */}
          <button onClick={sendToBackend}>Upload to Server</button>
        </div>
      )}
    </div>
  );
}

export default App;