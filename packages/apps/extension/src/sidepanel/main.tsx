import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App";
import "./styles/tailwind.css";

const root = document.getElementById("root");

if (root) {
  console.log("Creating React root");
  ReactDOM.createRoot(root).render(
    <React.StrictMode>
      <App />
    </React.StrictMode>
  );
} else {
  console.error("Root element not found");
}
