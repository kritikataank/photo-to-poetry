import React, { useRef } from "react";
import Webcam from "react-webcam";

const CameraCapture = ({ onCapture }) => {
  const webcamRef = useRef(null);

  const capture = () => {
    const imageSrc = webcamRef.current.getScreenshot();
    onCapture(imageSrc); // Send image data to parent (App.js)
  };

  return (
    <div className="camera-container">
      <Webcam ref={webcamRef} screenshotFormat="image/png" width={640} height={480} />
      <button onClick={capture}>Capture Photo</button>
    </div>
  );
};

export default CameraCapture;