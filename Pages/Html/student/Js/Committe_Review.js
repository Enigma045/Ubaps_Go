const tabs = document.querySelectorAll(".tabs input");
const tables = document.querySelectorAll(".log-table");
const filter = document.querySelector(".tabs").querySelectorAll("input");
let currentPage = 1;
const limit = 20;
let currentPayload = {
    statuses: [],
    search: '',
    department: [],
    scheme: [],
    parent: [],
    employment: [],
    gender: []
};

document.addEventListener("DOMContentLoaded", function () {
    // Initial fetch
    updatePayloadAndFetch();

    filter.forEach(e => {
        e.addEventListener("change", () => {
            currentPage = 1; // Reset to page 1 on filter change
            updatePayloadAndFetch();
        });
    });

    // Advanced filter bar listeners
    const searchEl = document.getElementById('filter-search');
    if (searchEl) {
        searchEl.addEventListener('input', () => {
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }

    const advancedFilter = document.getElementById('filter-advanced');
    if (advancedFilter) {
        // Event is handled in Filter_Feature.js, which calls window.onFilterChange
        window.onFilterChange = () => {
            currentPage = 1;
            updatePayloadAndFetch();
        };
    }

    const clearBtn = document.getElementById('clear-filters');
    if (clearBtn) {
        clearBtn.addEventListener('click', () => {
            if (searchEl) searchEl.value = '';
            if (advancedFilter) advancedFilter.value = '';
            window.activeAdvancedFilters = {}; // Clear persistent state
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }
});

function updatePayloadAndFetch() {
    let statusList = [];
    if (filter[0].checked) {
        statusList.push(filter[0].value);
        statusList.push("considering");
    }
    if (filter[1].checked) {
        statusList.push(filter[1].value);
        statusList.push("paid");
    }
    if (filter[2].checked) {
        statusList.push(filter[2].value);
    }

    const advanced = window.activeAdvancedFilters || {};

    currentPayload = {
        statuses: statusList,
        search: document.getElementById('filter-search')?.value || '',
        department: advanced.dept || [],
        scheme: advanced.scheme || [],
        parent: advanced.parent || [],
        employment: advanced.employment || [],
        gender: advanced.gender || []
    };

    fetchApplicants(currentPage);

    if (window.renderFilterTags) {
        window.renderFilterTags(currentPayload, () => {
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }
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
                    review: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M480-320q75 0 127.5-52.5T660-500q0-75-52.5-127.5T480-680q-75 0-127.5 52.5T300-500q0 75 52.5 127.5T480-320Zm0-72q-45 0-76.5-31.5T372-500q0-45 31.5-76.5T480-608q45 0 76.5 31.5T588-500q0 45-31.5 76.5T480-392Zm0 192q-146 0-266-81.5T35-500q61-137 181-218.5T480-800q146 0 266 81.5T925-500q-61 137-181 218.5T480-200Zm0-300Zm0 220q113 0 207.5-59.5T832-500q-50-101-144.5-160.5T480-720q-113 0-207.5 59.5T128-500q50 101 144.5 160.5T480-280Z"/></svg>`,
                    reject: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M240-840h440v520L400-40l-50-50q-7-7-11.5-19t-4.5-23v-14l44-174H120q-32 0-56-24t-24-56v-80q0-7 2-15t4-15l120-282q9-20 30-34t44-14Zm360 80H240L120-480v80h360l-54 220 174-174v-406Zm0 406v-406 406Zm80 34v-80h120v-360H680v-80h200v520H680Z" /></svg>`,
                    approve: `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M720-120H280v-520l280-280 50 50q7 7 11.5 19t4.5 23v14l-44 174h258q32 0 56 24t24 56v80q0 7-2 15t-4 15L794-168q-9 20-30 34t-44 14Zm-360-80h360l120-280v-80H480l54-220-174 174v406Zm0-406v406-406Zm-80-34v80H160v360h120v80H80v-520h200Z" /></svg>`,
                    rollback: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path><path d="M3 3v5h5"></path></svg>`,
                    comment: `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>`
                };

                let actions = `<button title="Review" class="action-btn review-btn" data-id="${item.registration_number}">${icons.review}</button>`;
                actions += `<button title="Comment" class="action-btn comment-btn" data-id="${item.registration_number}" data-name="${item.first_name} ${item.last_name}" data-comments='${JSON.stringify(item.comments || [])}'>${icons.comment}</button>`;
                
                if (item.status === "submitted" || item.status === "considering") {
                    actions += `<button title="Approve" class="action-btn approve-btn" data-name="${item.first_name} ${item.last_name}" data-id="${item.registration_number}">${icons.approve}</button>`;
                }
                if (["submitted", "selected", "considering"].includes(item.status)) {
                    actions += `<button title="Reject" class="action-btn reject-btn" data-name="${item.first_name} ${item.last_name}" data-id="${item.registration_number}">${icons.reject}</button>`;
                }
                if (item.status === "selected") {
                    actions += `<button title="Rollback" class="action-btn rollback-btn" data-name="${item.first_name} ${item.last_name}" data-id="${item.registration_number}">${icons.rollback}</button>`;
                }

                tr.innerHTML = `
                    <td>${item.first_name} ${item.last_name}</td>
                    <td>${item.registration_number}</td>
                    <td>${item.gender}</td>
                    <td>${item.parent_guardian_status}</td>
                    <td>${item.guardian_employment_status}</td>
                    <td>${item.relative_support}</td>
                    <td>low</td>
                    <td>${actions}</td>
                `;
                tbody.appendChild(tr);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="8" style="text-align:center;">No applicants found</td></tr>';
        }

        renderPagination(result.total, page);
        initModalLogic();

        // Trigger filter after dynamic load
        if (window.triggerTableFilter) window.triggerTableFilter();

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
    const rollbackModal = document.getElementById("rollbackModal");
    const commentModal = document.getElementById("commentModal");
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
        } else if (btn.classList.contains("rollback-btn")) {
            document.getElementById("rollbackName").textContent = name;
            rollbackModal.classList.add("active");
        } else if (btn.classList.contains("comment-btn")) {
            const rawComments = btn.getAttribute("data-comments");
            openCommentModal(id, name, rawComments);
        } else if (btn.classList.contains("review-btn")) {
            // Review logic is handled by other scripts like card.js
        }
    };

    const closeModals = () => {
        approveModal.classList.remove("active");
        rejectModal.classList.remove("active");
        if (rollbackModal) rollbackModal.classList.remove("active");
        if (commentModal) commentModal.classList.remove("active");
    };

    document.querySelectorAll(".close-modal, .cancel, .btn-cancel").forEach(btn => {
        btn.addEventListener("click", closeModals);
    });

    function openCommentModal(id, name, rawComments) {
        let comments = [];
        try {
            comments = JSON.parse(rawComments) || [];
        } catch (e) {
            console.error("Failed to parse comments", e);
        }

        const list = document.getElementById("comments-list");
        list.innerHTML = "";

        if (comments.length === 0) {
            list.innerHTML = `<p style="text-align: center; color: var(--sr-slate-400); font-size: 0.85rem; margin-top: 20px;">No comments yet.</p>`;
        } else {
            comments.forEach(c => {
                const date = new Date(parseFloat(c.date) * 1000).toLocaleString();
                const div = document.createElement("div");
                div.style.background = "white";
                div.style.padding = "10px";
                div.style.borderRadius = "6px";
                div.style.borderLeft = "3px solid var(--sr-blue)";
                div.style.boxShadow = "0 1px 2px rgba(0,0,0,0.05)";
                div.innerHTML = `
                    <div style="display: flex; justify-content: space-between; margin-bottom: 5px;">
                        <span style="font-weight: 700; font-size: 0.8rem; color: var(--sr-slate-700);">${c.name} (${c.role})</span>
                        <span style="font-size: 0.7rem; color: var(--sr-slate-400);">${date}</span>
                    </div>
                    <p style="margin: 0; font-size: 0.85rem; color: var(--sr-slate-600); line-height: 1.4;">${c.text}</p>
                `;
                list.appendChild(div);
            });
            list.scrollTop = list.scrollHeight;
        }

        commentModal.classList.add("active");
    }

    document.getElementById("submitCommentBtn").onclick = () => {
        const text = document.getElementById("new-comment-text").value.trim();
        if (!text) return;

        fetch("/api/add-comment", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: activeStudent.id, comment: text })
        }).then(async response => {
            if (!response.ok) throw new Error(await response.text() || "Failed to add comment");
            return response.json();
        }).then(data => {
            showToast("Comment posted!", "success");
            document.getElementById("new-comment-text").value = "";
            fetchApplicants(currentPage);
        }).catch(error => {
            showToast(error.message, "error");
        });
        closeModals();
    };

    document.getElementById("confirmRollback").onclick = () => {
        fetch("/api/rollback-selection", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: activeStudent.id })
        }).then(async response => {
            if (!response.ok) throw new Error(await response.text() || "Rollback failed");
            return response.json();
        }).then(data => {
            showToast(`Selection for ${activeStudent.name} rolled back!`, "success");
            fetchApplicants(currentPage);
        }).catch(error => {
            showToast(error.message, "error");
        });
        closeModals();
    };

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

