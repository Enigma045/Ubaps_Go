const tabs = document.querySelectorAll(".tabs input");
const tables = document.querySelectorAll(".log-table");

const filter = document.querySelector(".tabs").querySelectorAll("input");
const payload = [];
let reason;

fetch('/gettotalamount', {
  headers: {
    "Content-Type": "application/json"
  }
}).then(response => response.json())
  .then(data => {
    console.log(data)
  }).catch(error => {
    console.error('Error:', error);
  })


filter.forEach(e => {
  e.addEventListener("change", () => {
    payload.length = 0; // Clear the payload array
    if (filter[0].checked == true) {
      //console.log(filter[0].value)
      payload.push(filter[0].value)
    }
    if (filter[1].checked == true) {
      payload.push(filter[1].value)
    }
    if (filter[2].checked == true) {

      payload.push(filter[2].value)
    }
    // if (filter[3].checked == true) {
    //   payload.push(filter[3].value)
    // }

    console.log(JSON.stringify(payload))

    fetch('/getrequest', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    }).then(response => response.json())
      .then(data => {
        console.log(data)
        const tbody = document.getElementById("judge_tbody");
        tbody.innerHTML = ""; // Clear table once before rendering

        data.forEach(item => {
          const tr = document.createElement("tr");

          // Action icons constants
          const icons = {
            review: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M560-680v-80h320v80H560Zm0 160v-80h320v80H560Zm0 160v-80h320v80H560Zm-240-40q-50 0-85-35t-35-85q0-50 35-85t85-35q50 0 85 35t35 85q0 50-35 85t-85 35ZM80-160v-76q0-21 10-40t28-30q45-27 95.5-40.5T320-360q56 0 106.5 13.5T522-306q18 11 28 30t10 40v76H80Zm86-80h308q-35-20-74-30t-80-10q-41 0-80 10t-74 30Zm154-240q17 0 28.5-11.5T360-520q0-17-11.5-28.5T320-560q-17 0-28.5 11.5T280-520q0 17 11.5 28.5T320-480Zm0-40Zm0 280Z" /></svg>`,
            reject: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M240-840h440v520L400-40l-50-50q-7-7-11.5-19t-4.5-23v-14l44-174H120q-32 0-56-24t-24-56v-80q0-7 2-15t4-15l120-282q9-20 30-34t44-14Zm360 80H240L120-480v80h360l-54 220 174-174v-406Zm0 406v-406 406Zm80 34v-80h120v-360H680v-80h200v520H680Z" /></svg>`,
            approve: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M720-120H280v-520l280-280 50 50q7 7 11.5 19t4.5 23v14l-44 174h258q32 0 56 24t24 56v80q0 7-2 15t-4 15L794-168q-9 20-30 34t-44 14Zm-360-80h360l120-280v-80H480l54-220-174 174v406Zm0-406v406-406Zm-80-34v80H160v360h120v80H80v-520h200Z" /></svg>`
          };

          // Review button is always present
          //let actions = `<button title="Review" class="action-btn review-btn">${icons.review}</button>`;
          let actions = "";
          // Add Reject button only for submitted items
          if (["sent", "rejected"].includes(item.status)) {
            actions += `<button title="Approve" class="action-btn approve-btn" data-name="${item.first_name} ${item.last_name}" data-id="${item.id}">${icons.approve}</button>`;
          }

          // Add Approve button for specific statuses
          if (["sent", "approved"].includes(item.status)) {
            actions += `<button title="Reject" class="action-btn reject-btn" data-name="${item.first_name} ${item.last_name}" data-id="${item.id}">${icons.reject}</button>`;
          }


          tr.innerHTML = `
          <td>REQ#${item.id}</td>
          <td>${item.first_name} ${item.last_name}</td>
          <td>${item.payment_amount}</td>
          <td>${item.reason.String}</td>
          <td>${actions}</td>
        `;

          tbody.appendChild(tr);
        });

        // MODAL LOGIC Initialization
        const approveModal = document.getElementById("approveModal");
        const rejectModal = document.getElementById("rejectModal");
        let activeStudent = null;

        // Event Delegation for Table Actions
        tbody.onclick = e => {
          const btn = e.target.closest(".action-btn");
          if (!btn) return;

          const name = btn.getAttribute("data-name");
          const id = btn.getAttribute("data-id");
          activeStudent = { name, id };

          if (btn.classList.contains("approve-btn")) {
            document.getElementById("approveName").textContent = name;
            approveModal.classList.add("active");
          } else if (btn.classList.contains("reject-btn")) {
            document.getElementById("rejectName").textContent = name;
            rejectModal.classList.add("active");
          }
        };

        // Modal Closure Logic
        const closeModals = () => {
          approveModal.classList.remove("active");
          rejectModal.classList.remove("active");
        };

        document.querySelectorAll(".close-modal, .cancel").forEach(btn => {
          btn.addEventListener("click", closeModals);
        });

        // Confirm Logic placeholders
        document.getElementById("confirmApprove").onclick = () => {
          console.log("Approving:", activeStudent);
          //here
          fetch("/acceptrequest", {
            method: "POST",
            headers: {
              "Content-Type": "application/json"
            },
            body: JSON.stringify(activeStudent)
          }).then(async response => {
            if (!response.ok) {
              const errorText = await response.text();
              throw new Error(errorText || "Approval failed");
            }
            return response.json();
          })
            .then(data => {
              console.log(data);
              showToast(`Request from ${activeStudent.name} approved!`, "success");
            }).catch(error => {
              console.error('Error:', error);
              showToast(error.message || "Failed to approve request.", "error");
            });

          closeModals();
        };

        document.getElementById("confirmReject").onclick = () => {
          console.log("Rejecting:", activeStudent);
          //here
          fetch("/rejectrequest", {
            method: "POST",
            headers: {
              "Content-Type": "application/json"
            },
            body: JSON.stringify(activeStudent)
          }).then(async response => {
            if (!response.ok) {
              const errorText = await response.text();
              throw new Error(errorText || "Rejection failed");
            }
            return response.json();
          })
            .then(data => {
              console.log(data);
              showToast(`Request from ${activeStudent.name} rejected.`, "success");
            }).catch(error => {
              console.error('Error:', error);
              showToast(error.message || "Failed to reject request.", "error");
            });
          closeModals();
        };
      }).catch(error => {
        console.error('Error:', error);
      });
  });
});


