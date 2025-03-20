import React, { useState } from "react";
import CameraCapture from "./components/CameraCapture";

function App() {
  const [capturedImage, setCapturedImage] = useState(null);
  const [imageName, setImageName] = useState(null);
  const [caption, setCaption] = useState("");
  const [poem, setPoem] = useState("");

  // Handle File Upload
  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setCapturedImage(reader.result);
      };
      reader.readAsDataURL(file);
    }
  };

  // Upload Image to Backend
  const uploadToBackend = async () => {
    if (!capturedImage) {
      alert("Capture or upload an image first!");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ image: capturedImage }),
      });

      const data = await response.json();

      if (data.image_name) {
        setImageName(data.image_name);
        setCaption(""); 
        setPoem(""); 
      } else {
        alert("Image upload failed.");
      }
    } catch (error) {
      console.error("❌ Error uploading image:", error);
      alert("Error uploading image.");
    }
  };

  // Fetch Caption from Backend
  const fetchCaption = async () => {
    if (!imageName) {
      alert("Upload an image first!");
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/caption/${imageName}`);
      const data = await response.json();

      if (data.caption) {
        setCaption(data.caption);
      } else {
        alert("Failed to generate caption.");
      }
    } catch (error) {
      console.error("❌ Error fetching caption:", error);
      alert("Error retrieving caption.");
    }
  };

  // Convert Caption to Poetry
  const convertToPoetry = async () => {
    if (!caption) {
      alert("Generate a caption first!");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/convert", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ caption }),
      });

      const data = await response.json();

      if (data.poem) {
        setPoem(data.poem);
      } else {
        alert("Failed to generate poetry.");
      }
    } catch (error) {
      console.error("❌ Error converting text:", error);
      alert("Error generating poetry.");
    }
  };

  return (
    <div className="App">
      <h1>Photo-to-Poetry Web App</h1>

      <CameraCapture onCapture={setCapturedImage} />
      <input type="file" accept="image/*" onChange={handleFileUpload} />
      <br />

      {capturedImage && (
        <div>
          <h2>Selected Image:</h2>
          <img src={capturedImage} alt="Captured" style={{ maxWidth: "300px" }} />
          <br />
          <button onClick={uploadToBackend}>Upload to Server</button>
        </div>
      )}

      {imageName && (
        <div>
          <h2>Uploaded Image:</h2>
          <img src={`http://localhost:8080/image/${imageName}`} alt="Uploaded" style={{ maxWidth: "300px" }} />
          <br />
          <button onClick={fetchCaption}>Generate Caption</button>
          {caption && (
            <div>
              <h3>Generated Caption:</h3>
              <p>{caption}</p>
              <button onClick={convertToPoetry}>Convert to Poetry</button>
            </div>
          )}
          {poem && (
            <div>
              <h3>Generated Poem:</h3>
              <p>{poem}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default App;
