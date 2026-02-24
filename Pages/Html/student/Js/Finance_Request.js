document.addEventListener("DOMContentLoaded", () => {
  console.log("hello world");
  const modal = document.getElementById("modal");
  let but = document.getElementsByClassName('openModal')
  let selectedStudentId = null;

  // Capture student_id from clicked row
  Array.from(but).forEach(h => {
    h.addEventListener("click", e => {
      if (!e.target.classList.contains('openModal')) return;

      const row = e.target.closest("tr");
      if (!row) return;

      selectedStudentId = row.querySelector(".student_id").textContent.trim();
      console.log("Selected Student ID:", selectedStudentId);
      modal.classList.add("active");
    });
  })


  // Submit form
  document.getElementById("request_info").onsubmit = async e => {
    e.preventDefault();

    if (!selectedStudentId) {
      showToast("Please select a row first", "error");
      return;
    }

    const formData = new FormData(e.target);
    formData.append("student_id", selectedStudentId);

    console.log(Array.from(formData.entries()));

    const res = await fetch("/fees", {
      method: "POST",
      body: formData,
      credentials: "include"
    });

    if (res.ok) {
      showToast("Financial request submitted successfully!", "success");
      setTimeout(() => window.location.reload(), 2000);
    } else {
      const text = await res.text();
      console.log(text);
      showToast(text || "Failed to submit request.", "error");
    }
  };
});
