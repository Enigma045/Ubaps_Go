document.addEventListener("DOMContentLoaded", () => {
  const modal = document.getElementById("modal");
  const openBtn = document.querySelectorAll(".openModal");
  const closeBtn = document.getElementById("closeModal2");
  const cancelBtn = document.getElementById("cancel");
  let email = null;
  console.log(modal, openBtn, closeBtn);

  // ✅ Event Delegation for Delete Button
  document.querySelector(".tbody").addEventListener("click", function (e) {
    if (e.target.closest(".openModal")) {
      const row = e.target.closest("tr");
      const identity = row.querySelector(".name").innerText;
      email = row.querySelector(".email").innerText;

      document.querySelector(".delete").textContent =
        "Are you sure you want to delete " + identity + "'s account?";

      modal.classList.add("active");
    }
  });


  const closeModal = (e) => {
    e.preventDefault();
    modal.classList.remove("active");
  };

  closeBtn.addEventListener("click", closeModal);
  if (cancelBtn) cancelBtn.addEventListener("click", closeModal);

  // Close when clicking outside card
  modal.addEventListener("click", (e) => {
    if (e.target === modal) {
      modal.classList.remove("active");
    }
  });

  document.querySelector('.DeleteAccount').addEventListener('click', e => {
    e.preventDefault()
    console.log(email)

    fetch("/deleteaccount", {
      method: "POST",
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        email: email
      })
    }).then(response => {
      if (response.ok) {
        location.reload(); // Reload the page to reflect changes
      } else {
        console.error("Failed to delete account");
      }
    }).catch(error => {
      console.error("Error:", error);
    });
  })
})
