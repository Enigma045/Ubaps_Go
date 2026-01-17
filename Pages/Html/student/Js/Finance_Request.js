document.addEventListener("DOMContentLoaded", () => {
  console.log("hello world");
  const modal = document.getElementById("modal");
  let but = document.getElementsByClassName('openModal')
 let selectedStudentId = null;

 // Capture student_id from clicked row
   Array.from(but).forEach(h=>{
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
      alert("Please select a row first");
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
      console.log("Success");
    } else {
      console.log(await res.text());
    }
  };
});
