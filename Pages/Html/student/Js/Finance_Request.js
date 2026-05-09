document.addEventListener("DOMContentLoaded", () => {
    const modal = document.getElementById("modal");
    const closeModal = document.getElementById("closeModal");
    const tableBody = document.getElementById("statement-requests-body");
    let selectedStudentId = null;

    // Fetch and render statement requests
    async function fetchStatementRequests() {
        try {
            const res = await fetch("/getstatementrequests", { credentials: "include" });
            if (!res.ok) throw new Error("Failed to fetch statement requests");
            const data = await res.json();
            renderTable(data);
        } catch (err) {
            console.error(err);
            showToast("Error loading statement requests", "error");
        }
    }

    function renderTable(requests) {
        if (!tableBody) return;
        tableBody.innerHTML = "";

        if (requests.length === 0) {
            tableBody.innerHTML = `<tr><td colspan="7" style="text-align:center;">No pending statement requests</td></tr>`;
            return;
        }

        requests.forEach(req => {
            const tr = document.createElement("tr");
            tr.className = "success";
            tr.innerHTML = `
                <td>${req.first_name} ${req.last_name}</td>
                <td class="student_id">${req.registration_number}</td>
                <td>${req.parent_guardian_status}</td>
                <td>${req.employment_status}</td>
                <td>${req.income}</td>
                <td><span class="status warning">${req.priority}</span></td>
                <td>
                    <div style="display: flex; gap: 5px;">
                        <button class="openModal">Response</button>
                        <button class="viewDossier" style="background: #0ea5e9;">Dossier</button>
                    </div>
                </td>
            `;
            tableBody.appendChild(tr);
        });

        // Re-bind modal events
        tableBody.querySelectorAll(".openModal").forEach(btn => {
            btn.onclick = (e) => {
                const row = e.target.closest("tr");
                selectedStudentId = row.querySelector(".student_id").textContent.trim();
                modal.classList.add("active");
            };
        });

        // Bind Dossier events
        tableBody.querySelectorAll(".viewDossier").forEach(btn => {
            btn.onclick = (e) => {
                const row = e.target.closest("tr");
                const studentId = row.querySelector(".student_id").textContent.trim();
                const studentName = row.cells[0].textContent.trim();
                openDossier(studentId, studentName);
            };
        });

        // Trigger filter after dynamic load
        if (window.triggerTableFilter) window.triggerTableFilter();
    }

    if (closeModal) {
        closeModal.onclick = () => modal.classList.remove("active");
    }

    // Submit form
    document.getElementById("request_info").onsubmit = async e => {
        e.preventDefault();

        if (!selectedStudentId) {
            showToast("Please select a student first", "error");
            return;
        }

        const formData = new FormData(e.target);
        formData.append("student_id", selectedStudentId);

        const res = await fetch("/fees", {
            method: "POST",
            body: formData,
            credentials: "include"
        });

        if (res.ok) {
            showToast("Financial statement submitted successfully!", "success");
            modal.classList.remove("active");
            fetchStatementRequests(); // Refresh table
        } else {
            const text = await res.text();
            showToast(text || "Failed to submit statement.", "error");
        }
    };

    // Initial fetch
    fetchStatementRequests();
});
