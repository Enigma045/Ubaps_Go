document.addEventListener('DOMContentLoaded', () => {
    const reviewModal = document.getElementById('reviewModalOverlay');
    const closeBtn = document.getElementById('closeReviewCard');
    const tbody = document.getElementById('judge_tbody');

    // Data placeholders in modal
    const fields = {
        name: document.getElementById('review-name'),
        id: document.getElementById('review-id'),
        course: document.getElementById('review-course'),
        gender: document.getElementById('review-gender'),
        parents: document.getElementById('review-parents'),
        employment: document.getElementById('review-employment'),
        income: document.getElementById('review-income'),
        priority: document.getElementById('review-priority'),
        residence: document.getElementById('review-residence'), // We'll need to infer or add this
        guardian: document.getElementById('review-guardian'),
        select: document.getElementById('review-scheme-select'),
        amount: document.getElementById('review-amount-input'),
        reviewStatus: document.getElementById('review-status'),
        reviewAmount: document.getElementById('review-amount'),
    };

    const assignBtn = document.getElementById('review-assign-btn');

    // Open logic
    if (tbody) {
        tbody.addEventListener('click', async (e) => {
            const reviewBtn = e.target.closest('.review-btn');
            if (!reviewBtn) return;

        const row = reviewBtn.closest('tr');
        const cells = row.cells;

        // Extract data (based on the table structure in Committe_Review.html)
        // 0: Name, 1: Reg #, 2: Gender, 3: Parent Status, 4: Employment, 5: Income, 6: Priority
        const data = {
            name: cells[0].textContent.trim(),
            id: cells[1].textContent.trim(),
            gender: cells[2].textContent.trim(),
            parents: cells[3].textContent.trim(),
            employment: cells[4].textContent.trim(),
            income: cells[5].textContent.trim(),
            priority: cells[6].textContent.trim()
        };

        // Populate modal
        fields.name.textContent = data.name;
        fields.id.textContent = data.id;
        fields.gender.textContent = data.gender;
        fields.parents.textContent = data.parents;
        fields.employment.textContent = data.employment;
        fields.income.textContent = data.income;
        fields.priority.textContent = data.priority;

        // These might need mapping or fallback
        fields.course.textContent = "Student @ UNILIA";
        fields.residence.textContent = "N/A";
        fields.guardian.textContent = data.parents.includes('Father') ? 'Father' : 'N/A';
        // document.createElement('option').value = "Loans";
        //temp.textContent = "Loans"
        //
        if (fields.select) {
            let road = await getScheme();
            console.log(road);

            road.forEach(item => {
                const option = document.createElement('option');
                option.value = item.scheme_name;
                option.textContent = item.scheme_name;
                fields.select.appendChild(option);
            });
        }

        // Hide sidebar until approvals are verified
        const sidebar = document.querySelector('.sidebar-action');
        if (sidebar) sidebar.style.display = 'none';

        reviewModal.classList.add('active');

        // Fetch and update approval status tracker
        fetchApprovalStatus(data.id);
    });
}

    // Close logic
    if (closeBtn) {
        closeBtn.addEventListener('click', () => {
            if (reviewModal) reviewModal.classList.remove('active');
        });
    }

    if (reviewModal) {
        window.addEventListener('click', (e) => {
            if (e.target === reviewModal) {
                reviewModal.classList.remove('active');
            }
        });
    }

    // Toast Feedback for card buttons
    const statementBtn = document.querySelector('.card-btn.primary');
    const dossierBtn = document.querySelector('.card-btn.outline');

    // Statement Request Modal Elements
    const statementModal = document.getElementById('requestStatementModal');
    const closeStatementBtn = document.getElementById('closeStatementModal');
    const cancelStatementBtn = document.getElementById('cancelStatement');
    const confirmStatementBtn = document.getElementById('confirmStatementRequest');
    const statementStudentName = document.getElementById('statementStudentName');

    if (statementBtn) {
        statementBtn.addEventListener('click', () => {
            console.log("Request Statements button clicked");
            const studentId = fields.id.textContent.trim();
            console.log("Current Student ID in modal:", studentId);

            if (!studentId || studentId === "ID Number") {
                showToast("No student selected.", "error");
                return;
            }
            // Populate and show the confirmation modal
            if (statementStudentName) statementStudentName.textContent = fields.name.textContent;
            console.log("Showing statement confirmation modal for:", fields.name.textContent);
            if (statementModal) {
                statementModal.classList.add('active');
            } else {
                console.error("Statement confirmation modal not found in DOM");
                showToast("Confirmation modal error.", "error");
            }
        });
    }

    const closeStatementModals = () => {
        if (statementModal) statementModal.classList.remove('active');
    };

    if (closeStatementBtn) closeStatementBtn.addEventListener('click', closeStatementModals);
    if (cancelStatementBtn) cancelStatementBtn.addEventListener('click', closeStatementModals);

    if (confirmStatementBtn) {
        confirmStatementBtn.addEventListener('click', () => {
            const studentId = fields.id.textContent.trim();
            console.log("Confirm button clicked. Requesting statement for:", studentId);
            
            fetch("/requeststatement", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ reg: studentId })
            }).then(async res => {
                console.log("Response status:", res.status);
                if (!res.ok) {
                    const errorText = await res.text();
                    throw new Error(errorText || "Failed to request statement");
                }
                return res.json();
            }).then(data => {
                console.log("Success data:", data);
                showToast("Statement request sent successfully.", "success");
                closeStatementModals();
            }).catch(err => {
                console.error("Statement request fetch error:", err);
                showToast(err.message || "Failed to request statement.", "error");
                closeStatementModals();
            });
        });
    }

    if (dossierBtn) {
        dossierBtn.addEventListener('click', () => {
            const studentId = fields.id.textContent.trim();
            const studentName = fields.name.textContent.trim();
            openDossier(studentId, studentName);
        });
    }


    // Finalize Assignment Logic (only if sidebar exists)
    if (assignBtn) {
        assignBtn.addEventListener('click', () => {
            const scheme = document.getElementById('review-scheme-select').value;
            const amount = document.getElementById('review-amount-input').value;

            const studentId = fields.id.textContent;

            if (scheme === 'none' || !amount) {
                showToast("Please select a scheme and enter an amount.", "error");
                return;
            }

            console.log(`Finalizing assignment for ${studentId}: ${scheme} - MWK ${amount}`);
            fetch("/schemeinfo", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    reg: studentId,
                    amount: amount,
                    scheme: scheme
                })
            }).then(res => {
                if (!res.ok) {
                    throw new Error("Failed to send scheme info");
                }
                return res.json();
            }).then(data => {
                console.log(data);
                showToast(`Bursary assigned successfully to ${fields.name.textContent}`, "success");
                reviewModal.classList.remove('active');
            }).catch(err => {
                console.error(err);
                showToast(err.message || "Failed to finalize assignment.", "error");
            });
        });
    }

    async function getScheme() {
        try {
            const res = await fetch("/getschemes");

            if (!res.ok) {
                throw new Error("Failed to fetch schemes");
            }

            const data = await res.json();
            return data;

        } catch (err) {
            console.error(err);
            throw err;
        }
    }

    // ─── Approval Status Tracker ────────────────────────────────
    function fetchApprovalStatus(regNumber) {
        fetch("/getapplicationstatus", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ reg: regNumber })
        })
            .then(res => {
                if (!res.ok) throw new Error("Failed to fetch approval status");
                return res.json();
            })
            .then(data => {
                console.log("Approval status:", data);

                // Map response keys to badge element IDs
                const mapping = {
                    registrar_approval_status: "approval-registrar",
                    dean_of_student_approval_status: "approval-dean_of_student",
                    dean_of_facult_approval_status: "approval-dean_of_facult",
                    dean_of_science_approval_status: "approval-dean_of_science",
                    finance_office_approval_status: "approval-finance_office"
                };

                // Keys that must all be "approved" before the sidebar is shown
                const requiredApprovals = [
                    "dean_of_student_approval_status",
                    "dean_of_facult_approval_status",
                    "dean_of_science_approval_status",
                    "finance_office_approval_status"
                ];

                let allRequiredApproved = true;

                for (const [key, badgeId] of Object.entries(mapping)) {
                    const badge = document.getElementById(badgeId);
                    if (!badge) continue;

                    // Handle sql.NullString: {String: "approved", Valid: true}
                    let status = "pending";
                    if (data[key] && data[key].Valid) {
                        status = data[key].String.toLowerCase().trim();
                    } else if (typeof data[key] === "string" && data[key] !== "") {
                        status = data[key].toLowerCase().trim();
                    }

                    // Check if this required role is not yet approved
                    if (requiredApprovals.includes(key) && status !== "approved") {
                        allRequiredApproved = false;
                    }

                    // Reset classes
                    badge.className = "approval-badge";

                    if (status === "approved") {
                        badge.textContent = "Approved";
                        badge.classList.add("approved");
                    } else if (status === "rejected" || status === "not selected") {
                        badge.textContent = "Rejected";
                        badge.classList.add("rejected");
                    } else {
                        badge.textContent = "Pending";
                        badge.classList.add("pending");
                    }
                }

                // Update overall status and bursary amount using the new autonomous utility
                updateBursaryCardInfo(regNumber);

                // Show the bursary sidebar only when all required roles approved
                const sidebar = document.querySelector('.sidebar-action');
                if (sidebar) {
                    sidebar.style.display = allRequiredApproved ? '' : 'none';
                }
            })
            .catch(err => {
                console.error("Approval status error:", err);
            });
    }
});
