document.addEventListener("DOMContentLoaded", () => {
  const modal = document.getElementById("modal2");
  const openBtn = document.querySelectorAll(".openModal2");
  const closeBtn = document.getElementById("closeModal2");
  const cancelBtn = document.getElementById("cancel");
  let identity = null;
  let email = null;
  console.log(modal, openBtn, closeBtn);

  // ✅ Event Delegation for Delete Button
  document.querySelector(".bursarytbody").addEventListener("click", function (e) {
    if (e.target.closest(".openModal2")) {
      const row = e.target.closest("tr");
      identity = row.querySelector(".name").innerText;
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

    fetch("/deletebenefactor", {
      method: "POST",
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: identity,
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