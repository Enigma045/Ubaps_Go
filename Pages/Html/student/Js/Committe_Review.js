const tabs = document.querySelectorAll(".tabs input");
const tables = document.querySelectorAll(".log-table");
const filter = document.querySelector(".tabs").querySelectorAll("input");
let currentPage = 1;
const limit = 20;
let currentPayload = [];

document.addEventListener("DOMContentLoaded", function () {
    // Initial fetch
    updatePayloadAndFetch();

    filter.forEach(e => {
        e.addEventListener("change", () => {
            currentPage = 1; // Reset to page 1 on filter change
            updatePayloadAndFetch();
        });
    });
});

function updatePayloadAndFetch() {
    currentPayload = [];
    if (filter[0].checked) {
        currentPayload.push(filter[0].value);
        currentPayload.push("considering");
    }
    if (filter[1].checked) {
        currentPayload.push(filter[1].value);
    }
    if (filter[2].checked) {
        currentPayload.push(filter[2].value);
    }
    
    fetchApplicants(currentPage);
}

async function fetchApplicants(page) {
    try {
        const response = await fetch(`/applicants?page=${page}&limit=${limit}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(currentPayload)
        });
        
        const result = await response.json();
        const tbody = document.getElementById("judge_tbody");
        tbody.innerHTML = "";

        if (result.data && result.data.length > 0) {
            result.data.forEach(item => {
                const tr = document.createElement("tr");

                const icons = {
                    review: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M560-680v-80h320v80H560Zm0 160v-80h320v80H560Zm0 160v-80h320v80H560Zm-240-40q-50 0-85-35t-35-85q0-50 35-85t85-35q50 0 85 35t35 85q0 50-35 85t-85 35ZM80-160v-76q0-21 10-40t28-30q45-27 95.5-40.5T320-360q56 0 106.5 13.5T522-306q18 11 28 30t10 40v76H80Zm86-80h308q-35-20-74-30t-80-10q-41 0-80 10t-74 30Zm154-240q17 0 28.5-11.5T360-520q0-17-11.5-28.5T320-560q-17 0-28.5 11.5T280-520q0 17 11.5 28.5T320-480Zm0-40Zm0 280Z" /></svg>`,
                    reject: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M240-840h440v520L400-40l-50-50q-7-7-11.5-19t-4.5-23v-14l44-174H120q-32 0-56-24t-24-56v-80q0-7 2-15t4-15l120-282q9-20 30-34t44-14Zm360 80H240L120-480v80h360l-54 220 174-174v-406Zm0 406v-406 406Zm80 34v-80h120v-360H680v-80h200v520H680Z" /></svg>`,
                    approve: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M720-120H280v-520l280-280 50 50q7 7 11.5 19t4.5 23v14l-44 174h258q32 0 56 24t24 56v80q0 7-2 15t-4 15L794-168q-9 20-30 34t-44 14Zm-360-80h360l120-280v-80H480l54-220-174 174v406Zm0-406v406-406Zm-80-34v80H160v360h120v80H80v-520h200Z" /></svg>`
                };

                let actions = `<button title="Review" class="action-btn review-btn" data-id="${item.registration_number}">${icons.review}</button>`;
                if (item.status === "submitted" || item.status === "considering") {
                    actions += `<button title="Approve" class="action-btn approve-btn" data-name="${item.first_name} ${item.last_name}" data-id="${item.registration_number}">${icons.approve}</button>`;
                }
                if (["submitted", "selected", "considering"].includes(item.status)) {
                    actions += `<button title="Reject" class="action-btn reject-btn" data-name="${item.first_name} ${item.last_name}" data-id="${item.registration_number}">${icons.reject}</button>`;
                }

                tr.innerHTML = `
                    <td>${item.first_name} ${item.last_name}</td>
                    <td>${item.registration_number}</td>
                    <td>${item.parent_guardian_status}</td>
                    <td>${item.guardian_employment_status}</td>
                    <td>${item.relative_support}</td>
                    <td>low</td>
                    <td>${actions}</td>
                `;
                tbody.appendChild(tr);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;">No applicants found</td></tr>';
        }

        renderPagination(result.total, page);
        initModalLogic();

    } catch (error) {
        console.error("Error fetching applicants:", error);
    }
}

function renderPagination(total, page) {
    const container = document.getElementById("pagination-controls");
    if (!container) return;
    container.innerHTML = "";

    const totalPages = Math.ceil(total / limit);
    if (totalPages <= 1) return;

    const prevBtn = document.createElement("button");
    prevBtn.innerText = "Prev";
    prevBtn.disabled = page === 1;
    prevBtn.onclick = () => {
        currentPage--;
        fetchApplicants(currentPage);
    };
    container.appendChild(prevBtn);

    const start = Math.max(1, page - 2);
    const end = Math.min(totalPages, page + 2);

    for (let i = start; i <= end; i++) {
        const btn = document.createElement("button");
        btn.innerText = i;
        if (i === page) btn.classList.add("active");
        btn.onclick = () => {
            currentPage = i;
            fetchApplicants(currentPage);
        };
        container.appendChild(btn);
    }

    const nextBtn = document.createElement("button");
    nextBtn.innerText = "Next";
    nextBtn.disabled = page === totalPages;
    nextBtn.onclick = () => {
        currentPage++;
        fetchApplicants(currentPage);
    };
    container.appendChild(nextBtn);
}

function initModalLogic() {
    const tbody = document.getElementById("judge_tbody");
    const approveModal = document.getElementById("approveModal");
    const rejectModal = document.getElementById("rejectModal");
    let activeStudent = null;

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
        } else if (btn.classList.contains("review-btn")) {
            // Review logic is handled by other scripts like card.js
        }
    };

    const closeModals = () => {
        approveModal.classList.remove("active");
        rejectModal.classList.remove("active");
    };

    document.querySelectorAll(".close-modal, .cancel").forEach(btn => {
        btn.addEventListener("click", closeModals);
    });

    document.getElementById("confirmApprove").onclick = () => {
        fetch("/considerstudent", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(activeStudent)
        }).then(async response => {
            if (!response.ok) throw new Error(await response.text() || "Approval failed");
            return response.json();
        }).then(data => {
            showToast(`${activeStudent.name} marked for consideration!`, "success");
            fetchApplicants(currentPage);
        }).catch(error => {
            showToast(error.message, "error");
        });
        closeModals();
    };

    document.getElementById("confirmReject").onclick = () => {
        fetch("/rejectstudent", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(activeStudent)
        }).then(async response => {
            if (!response.ok) throw new Error(await response.text() || "Rejection failed");
            return response.json();
        }).then(data => {
            showToast(`${activeStudent.name} application rejected.`, "success");
            fetchApplicants(currentPage);
        }).catch(error => {
            showToast(error.message, "error");
        });
        closeModals();
    };
}

