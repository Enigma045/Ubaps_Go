const modal = document.getElementById("modal");
const closeBtn = document.getElementById("closeModal");
const openBtn = document.getElementById("openModal");

openBtn.addEventListener("click", () => {
  modal.classList.add("active");
});

closeBtn.addEventListener("click", () => {
  modal.classList.remove("active");
});

// Close when clicking outside card
modal.addEventListener("click", (e) => {
  if (e.target === modal) {
    modal.classList.remove("active");
  }
});