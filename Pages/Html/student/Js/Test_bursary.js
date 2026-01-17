const modal = document.getElementById("modal");
//const openBtn = document.getElementsByClassName("openModal");
const closeBtn = document.getElementById("closeModal");
//console.log( openBtn);
//let selectedStudentId = null;

// Array.from(openBtn).forEach( e =>{
//   e.addEventListener("click",m=>{
    
// modal.classList.add("active");
//     // if (!m.target.classList.contains("open")) return;

//     // const row = m.target.closest("tr");
//     // if (!row) return;
 
//     // selectedStudentId = row.querySelector(".student_id").textContent.trim();
//     // console.log("Selected Student ID:", selectedStudentId);
//   })
// })

// openBtn.addEventListener("click", () => {
//   modal.classList.add("active");
// });

closeBtn.addEventListener("click", () => {
  modal.classList.remove("active");
});

// Close when clicking outside card
modal.addEventListener("click", (e) => {
  if (e.target === modal) {
    modal.classList.remove("active");
  }
});