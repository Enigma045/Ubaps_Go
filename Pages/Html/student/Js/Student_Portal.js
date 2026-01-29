const dropZone = document.getElementById("dropZone");
const browseBtn = document.getElementById("browseBtn");
const fileInput = document.getElementById("fileInput");
const fileList = document.getElementById("fileList");

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
  [...files].forEach(file => {
    createFileItem(file);
  });
}

/* ==============================
   CREATE FILE ITEM UI
================================ */
function createFileItem(file) {
  const item = document.createElement("div");
  item.className = "file-item";

  const ext = file.name.split(".").pop().toUpperCase();
  const isPDF = ext === "PDF";

  item.innerHTML = `
    <div class="file-icon ${isPDF ? "pdf" : "txt"}">${ext}</div>

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

  const progressBar = item.querySelector(".progress-bar span");
  const percentText = item.querySelector(".percent");
  const removeBtn = item.querySelector(".remove");

  simulateUpload(progressBar, percentText, item);

  removeBtn.addEventListener("click", () => {
    item.remove();
  });
}

/* ==============================
   SIMULATE UPLOAD PROGRESS
================================ */
function simulateUpload(bar, percentText, item) {
  let progress = 0;

  const interval = setInterval(() => {
    progress += Math.random() * 10;

    if (progress >= 100) {
      progress = 100;
      clearInterval(interval);

      percentText.textContent = "✔";
      percentText.classList.add("check");

      bar.style.width = "100%";
      bar.style.background = "#5fc77a";
      return;
    }

    bar.style.width = progress + "%";
    percentText.textContent = Math.floor(progress) + "%";
  }, 300);
}
