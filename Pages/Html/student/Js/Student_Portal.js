const dropZone = document.getElementById("dropZone");
const browseBtn = document.getElementById("browseBtn");
const fileInput = document.getElementById("fileInput");
const fileList = document.getElementById("fileList");
const sendBtn = document.getElementById("sendBtn");

let selectedFile = null;

/* ==============================
   OPEN FILE PICKER
================================ */
browseBtn.addEventListener("click", () => {
  fileInput.click();
});

fileInput.addEventListener("change", () => {
  handleFiles(fileInput.files);
});

/* ==============================
   DRAG & DROP
================================ */
dropZone.addEventListener("dragover", (e) => {
  e.preventDefault();
  dropZone.classList.add("dragging");
});

dropZone.addEventListener("dragleave", () => {
  dropZone.classList.remove("dragging");
});

dropZone.addEventListener("drop", (e) => {
  e.preventDefault();
  dropZone.classList.remove("dragging");
  handleFiles(e.dataTransfer.files);
});

/* ==============================
   HANDLE FILES
================================ */
function handleFiles(files) {
  if (files.length === 0) return;
  const file = files[0];

  // Basic validation
  const ext = file.name.substring(file.name.lastIndexOf('.') + 1).toLowerCase();
  if (ext !== 'pdf' && ext !== 'doc' && ext !== 'docx') {
    showToast("Invalid file type. Only PDF and Word allowed.", "error");
    return;
  }

  selectedFile = file;
  fileList.innerHTML = '';
  createFileItem(file);
  
  // Show send button
  sendBtn.style.display = 'inline-block';
  sendBtn.disabled = false;
  sendBtn.style.opacity = '1';
}

/* ==============================
   CREATE FILE ITEM UI
================================ */
function createFileItem(file) {
  const item = document.createElement("div");
  item.className = "file-item";

  const ext = file.name.split(".").pop().toUpperCase();
  const isPDF = ext === "PDF";
  const iconClass = isPDF ? "pdf" : (ext === "DOC" || ext === "DOCX" ? "doc" : "txt");

  item.innerHTML = `
    <div class="file-icon ${iconClass}">${ext}</div>

    <div class="file-info">
      <p>${file.name}</p>
      <div class="progress-bar">
        <span></span>
      </div>
    </div>

    <span class="percent">0%</span>
    <button class="remove">✕</button>
  `;

  fileList.appendChild(item);

  const removeBtn = item.querySelector(".remove");
  removeBtn.addEventListener("click", () => {
    item.remove();
    selectedFile = null;
    sendBtn.style.display = 'none';
  });
}

/* ==============================
   SEND BUTTON CLICK
================================ */
sendBtn.addEventListener("click", () => {
  if (!selectedFile) return;
  
  const item = fileList.querySelector(".file-item");
  const progressBar = item.querySelector(".progress-bar span");
  const percentText = item.querySelector(".percent");
  
  sendBtn.disabled = true;
  sendBtn.style.opacity = '0.5';
  
  uploadFile(selectedFile, progressBar, percentText, item);
});

/* ==============================
   REAL UPLOAD LOGIC
================================ */
async function uploadFile(file, bar, percentText, item) {
  const formData = new FormData();
  formData.append("letter", file);

  try {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/submit-letter", true);

    // Track upload progress
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) {
        const progress = (e.loaded / e.total) * 100;
        bar.style.width = progress + "%";
        percentText.textContent = Math.floor(progress) + "%";
      }
    };

    xhr.onload = () => {
      if (xhr.status === 200) {
        bar.style.width = "100%";
        bar.style.background = "#5fc77a";
        percentText.textContent = "✔";
        percentText.classList.add("check");
        showToast("Letter uploaded successfully!", "success");
      } else {
        bar.style.background = "#ff4d4d";
        percentText.textContent = "✘";
        showToast(xhr.responseText || "Upload failed.", "error");
      }
    };

    xhr.onerror = () => {
      bar.style.background = "#ff4d4d";
      percentText.textContent = "✘";
      showToast("Network error occurred.", "error");
    };

    xhr.send(formData);
  } catch (error) {
    console.error("Upload error:", error);
    showToast("An unexpected error occurred.", "error");
  }
}
