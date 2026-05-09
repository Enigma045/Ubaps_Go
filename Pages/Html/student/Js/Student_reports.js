let tabs, tables, filter;
let currentPage = 1;
const limit = 20;
let currentPayload = {
    statuses: [],
    search: '',
    department: '',
    scheme: '',
    parent: '',
    employment: '',
    gender: ''
};

document.addEventListener("DOMContentLoaded", function () {
    tabs = document.querySelectorAll(".tabs input");
    tables = document.querySelectorAll(".log-table");
    filter = document.querySelector(".tabs")?.querySelectorAll("input") || [];

    // Initial fetch
    updatePayloadAndFetch();
    populateSchemes();

    if (filter.length > 0) {
        filter.forEach(e => {
            e.addEventListener("change", () => {
                currentPage = 1;
                updatePayloadAndFetch();
            });
        });
    }

    // Period selectors trigger stats refresh
    ['period-year','period-semester','period-month'].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.addEventListener('change', () => {
            fetchReportStats();
            updatePayloadAndFetch(); // Refresh table too
        });
    });

    // Generate Report button
    const genBtn = document.getElementById('generate-report-btn');
    if (genBtn) {
        genBtn.onclick = () => {
            if (!window.jspdf) {
                showToast("PDF library not loaded yet. Please wait or refresh.", "error");
                return;
            }
            generatePDFReport();
        };
    }

    // Advanced filter bar listeners
    const searchEl = document.getElementById('filter-search');
    if (searchEl) {
        searchEl.addEventListener('input', () => {
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }

    ['filter-dept', 'filter-scheme', 'filter-parent', 'filter-employment', 'filter-gender'].forEach(id => {
        const el = document.getElementById(id);
        if (el) el.addEventListener('change', () => {
            currentPage = 1;
            updatePayloadAndFetch();
        });
    });

    const clearBtn = document.getElementById('clear-filters');
    if (clearBtn) {
        clearBtn.addEventListener('click', () => {
            if (searchEl) searchEl.value = '';
            ['filter-dept', 'filter-scheme', 'filter-parent', 'filter-employment', 'filter-gender'].forEach(id => {
                const el = document.getElementById(id);
                if (el) el.value = '';
            });
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }

    // Initial stats load
    fetchReportStats();

});

async function generatePDFReport() {
    console.log("Starting PDF generation...");
    const yearVal     = document.getElementById('period-year')?.value     || '';
    const semesterVal = document.getElementById('period-semester')?.value || '';
    const monthVal    = document.getElementById('period-month')?.value    || '';

    const yearLabel     = yearVal || 'All Years';
    const semesterLabel = semesterVal ? `Semester ${semesterVal}` : 'All Semesters';
    const monthLabel    = monthVal ? `Month ${monthVal}` : 'All Months';

    showToast("Generating PDF Report...", "info");

    try {
        const url = `/api/reports/comprehensive?year=${yearVal}&semester=${semesterVal}&month=${monthVal}&limit=5000`;
        console.log("Fetching report data from:", url, "with payload:", currentPayload);
        
        const res = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(currentPayload)
        });

        if (!res.ok) {
            const errText = await res.text();
            throw new Error(`Server Error: ${res.status} ${errText}`);
        }
        
        const json = await res.json();
        const data = json.data || [];
        console.log(`Fetched ${data.length} records.`);

        if (data.length === 0) {
            showToast("No records found for the selected period", "warning");
            return;
        }

        if (!window.jspdf || !window.jspdf.jsPDF) {
            throw new Error("PDF library (jsPDF) is not fully loaded. Please refresh.");
        }

        const { jsPDF } = window.jspdf;
        const doc = new jsPDF('l', 'mm', 'a4');

        console.log("Constructing PDF header...");
        // Header Aesthetics
        doc.setFontSize(22);
        doc.setTextColor(30, 41, 59);
        doc.text("UBAPS", 14, 20);
        
        doc.setFontSize(10);
        doc.setTextColor(37, 99, 235);
        doc.text("University Bursary Application & Payment System", 44, 19);
        
        doc.setDrawColor(226, 232, 240);
        doc.line(14, 25, 283, 25);

        doc.setFontSize(14);
        doc.setTextColor(30, 41, 59);
        doc.text(`Comprehensive Administrative Report`, 14, 35);
        
        doc.setFontSize(10);
        doc.setTextColor(100, 116, 139);
        doc.text(`Period: ${yearLabel} / ${semesterLabel} / ${monthLabel}`, 14, 42);
        doc.text(`Total Records: ${data.length}`, 14, 47);
        doc.text(`Generated: ${new Date().toLocaleString()}`, 283, 47, { align: 'right' });

        console.log("Mapping table data...");
        const tableData = data.map(r => [
            r.registration_number || '—',
            `${r.name} ${r.surname}`,
            r.gender || '—',
            r.programme || '—',
            (r.status || '').toUpperCase(),
            r.scheme_name || 'No Scheme',
            Number(r.bursary_amount || 0).toLocaleString()
        ]);

        console.log("Generating AutoTable...");
        doc.autoTable({
            startY: 55,
            head: [['Reg. Number', 'Student Name', 'Sex', 'Department', 'Status', 'Bursary Scheme', 'Amount (MWK)']],
            body: tableData,
            theme: 'grid',
            headStyles: { 
                fillColor: [30, 41, 59], 
                textColor: 255, 
                fontSize: 9, 
                fontStyle: 'bold',
                halign: 'center'
            },
            styles: { 
                fontSize: 8, 
                cellPadding: 3, 
                lineColor: [226, 232, 240], 
                lineWidth: 0.1 
            },
            columnStyles: {
                0: { fontStyle: 'bold', halign: 'center' },
                2: { halign: 'center' },
                4: { halign: 'center' },
                6: { halign: 'right', fontStyle: 'bold' }
            },
            alternateRowStyles: { fillColor: [248, 250, 252] }
        });

        console.log("Finalizing PDF...");
        const finalY = doc.lastAutoTable.finalY || 55;
        if (finalY < 160) {
            doc.setFontSize(10);
            doc.setTextColor(30, 41, 59);
            doc.text("Authorized Signature: ___________________________", 14, finalY + 25);
            doc.text("Date of Approval: ___________________________", 14, finalY + 35);
        }

        const pageCount = doc.internal.getNumberOfPages();
        for(let i = 1; i <= pageCount; i++) {
            doc.setPage(i);
            doc.setFontSize(8);
            doc.setTextColor(148, 163, 184);
            doc.text(`Page ${i} of ${pageCount}`, 283, 200, { align: 'right' });
            doc.text("This is a system-generated document. Hash: " + btoa(Date.now().toString()).substring(0, 8), 14, 200);
        }

        console.log("Saving PDF...");
        doc.save(`UBAPS_Report_${yearLabel.replace(/\//g,'-')}.pdf`);
        showToast("Professional PDF Generated!", "success");

    } catch (err) {
        console.error("PDF Gen Error Details:", err);
        showToast(`Failed: ${err.message}`, "error");
    }
}

async function populateSchemes() {
    try {
        const res = await fetch("/getschemes");
        if (!res.ok) throw new Error("Failed to fetch schemes");
        const schemes = await res.json();
        const select = document.getElementById('filter-scheme');
        if (select) {
            schemes.forEach(scheme => {
                const option = document.createElement('option');
                option.value = scheme.scheme_name;
                option.textContent = scheme.scheme_name;
                select.appendChild(option);
            });
        }
    } catch (err) {
        console.error("Error populating schemes:", err);
    }
}

function updatePayloadAndFetch() {
    let statusList = [];
    let statusLabel = "All Applicants";

    if (filter[0].checked) {
        statusList.push(filter[0].value);
        statusList.push("considering");
        statusLabel = "Bursary Applicants";
    }
    if (filter[1].checked) {
        statusList.push(filter[1].value);
        statusList.push("paid");
        statusLabel = "Selected Students";
    }
    if (filter[2].checked) {
        statusList.push(filter[2].value);
        statusLabel = "Not Selected Candidates";
    }

    currentPayload = {
        statuses: statusList,
        search: document.getElementById('filter-search')?.value || '',
        department: document.getElementById('filter-dept')?.value || '',
        scheme: document.getElementById('filter-scheme')?.value || '',
        parent: document.getElementById('filter-parent')?.value || '',
        employment: document.getElementById('filter-employment')?.value || '',
        gender: document.getElementById('filter-gender')?.value || ''
    };

    const titleEl = document.getElementById('current-report-title');
    if (titleEl) titleEl.firstChild.textContent = statusLabel + " ";

    updateDynamicHeading();
    fetchApplicants(currentPage);
}

function updateDynamicHeading() {
    const year     = document.getElementById('period-year')?.value     || 'All Years';
    const semester = document.getElementById('period-semester')?.value || '';
    const month    = document.getElementById('period-month')?.value    || '';

    let periodText = `- ${year}`;
    if (semester) periodText += ` | Semester ${semester}`;
    if (month) {
        const monthNames = ["", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
        periodText += ` | ${monthNames[parseInt(month)]}`;
    }

    const periodEl = document.getElementById('current-report-period');
    if (periodEl) periodEl.textContent = periodText;
}

async function fetchApplicants(page) {
    const year     = document.getElementById('period-year')?.value     || '';
    const semester = document.getElementById('period-semester')?.value || '';
    const month    = document.getElementById('period-month')?.value    || '';

    try {
        const response = await fetch(`/applicants?page=${page}&limit=${limit}&year=${year}&semester=${semester}&month=${month}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(currentPayload)
        });

        const result = await response.json();
        
        const countEl = document.getElementById('report-result-count');
        if (countEl) countEl.textContent = `Showing ${result.data?.length || 0} of ${result.total || 0} records`;

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
                    <td>${item.gender}</td>
                    <td>${item.parent_guardian_status}</td>
                    <td>${item.guardian_employment_status}</td>
                    <td>${item.relative_support}</td>
                    <td>${(item.status && (item.status.toLowerCase() === "selected" || item.status.toLowerCase() === "paid")) ? (item.scheme_name && item.scheme_name !== 'No Scheme' ? item.scheme_name : "Not Assigned") : "Not Assigned"}</td>
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

async function fetchReportStats() {
    const year     = document.getElementById('period-year')?.value     || '';
    const semester = document.getElementById('period-semester')?.value || '';
    const month    = document.getElementById('period-month')?.value    || '';

    console.log('Fetching stats for:', { year, semester, month });

    try {
        const res = await fetch(`/api/reports/stats?year=${year}&semester=${semester}&month=${month}`);
        if (!res.ok) {
            const errorText = await res.text();
            throw new Error(`Stats fetch failed: ${res.status} ${errorText}`);
        }
        const s = await res.json();
        console.log('Received stats:', s);

        const set = (id, val) => { 
            const el = document.getElementById(id); 
            if (el) el.textContent = val; 
        };

        set('stat-total',    s.total_applications ?? '0');
        set('stat-approved', s.approved           ?? '0');
        set('stat-pending',  s.pending            ?? '0');
        set('stat-rejected', s.rejected           ?? '0');

    } catch (err) {
        console.error('fetchReportStats error:', err);
    }
}
