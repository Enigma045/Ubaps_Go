const modalRequest = document.getElementById("modalRequest");
const openBtnRequest = document.getElementById("openModalRequest");
const closeBtnRequest = document.getElementById("closeModal");
console.log(modalRequest, openBtnRequest, closeBtnRequest);
openBtnRequest.addEventListener("click", () => {
  modalRequest.classList.add("active");
});

closeBtnRequest.addEventListener("click", () => {
  modalRequest.classList.remove("active");
});

// Close when clicking outside card
modalRequest.addEventListener("click", (e) => {
  if (e.target === modalRequest) {
    modalRequest.classList.remove("active");
  }
});