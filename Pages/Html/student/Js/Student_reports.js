let tabs, tables, filter;
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
window.currentApplicants = [];

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
    ['period-year', 'period-semester', 'period-month'].forEach(id => {
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

    const advancedFilter = document.getElementById('filter-advanced');
    if (advancedFilter) {
        // Event is handled in Filter_Feature.js, which calls window.onFilterChange
        window.onFilterChange = () => {
            currentPage = 1;
            updatePayloadAndFetch();
        };
    }

    ['filter-date-start', 'filter-date-end', 'filter-cohort-start', 'filter-cohort-end'].forEach(id => {
        const el = document.getElementById(id);
        if (el) {
            el.addEventListener('change', () => {
                currentPage = 1;
                updatePayloadAndFetch();
            });
        }
    });

    const clearBtn = document.getElementById('clear-filters');
    if (clearBtn) {
        clearBtn.addEventListener('click', () => {
            if (searchEl) searchEl.value = '';
            if (advancedFilter) advancedFilter.value = '';
            ['filter-date-start', 'filter-date-end', 'filter-cohort-start', 'filter-cohort-end'].forEach(id => {
                const el = document.getElementById(id);
                if (el) el.value = '';
            });
            window.activeAdvancedFilters = {}; // Clear persistent state
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }

    // Initial stats load
    fetchReportStats();

});

async function generatePDFReport() {
    console.log("Starting PDF generation...");
    const yearVal = document.getElementById('period-year')?.value || '';
    const semesterVal = document.getElementById('period-semester')?.value || '';
    const monthVal = document.getElementById('period-month')?.value || '';

    const yearLabel = yearVal || 'All Years';
    const semesterLabel = semesterVal ? `Semester ${semesterVal}` : 'All Semesters';
    const monthLabel = monthVal ? `Month ${monthVal}` : 'All Months';

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

        const reportTitle = window.currentReportTitle || "Comprehensive Administrative Report";
        doc.text(reportTitle, 14, 35);

        doc.setFontSize(10);
        doc.setTextColor(100, 116, 139);
        doc.text(`Period: ${yearLabel} / ${semesterLabel} / ${monthLabel}`, 14, 42);
        doc.text(`Total Records: ${data.length}`, 14, 47);
        doc.text(`Generated: ${new Date().toLocaleString()}`, 283, 47, { align: 'right' });

        console.log("Mapping table data...");
        const tableData = data.map(r => [
            r.registration_number || '—',
            new Date(r.application_date).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' }),
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
            head: [['Reg. Number', 'App. Date', 'Student Name', 'Sex', 'Department', 'Status', 'Bursary Scheme', 'Amount (MWK)']],
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
                1: { halign: 'center' },
                3: { halign: 'center' },
                5: { halign: 'center' },
                7: { halign: 'right', fontStyle: 'bold' }
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
        for (let i = 1; i <= pageCount; i++) {
            doc.setPage(i);
            doc.setFontSize(8);
            doc.setTextColor(148, 163, 184);
            doc.text(`Page ${i} of ${pageCount}`, 283, 200, { align: 'right' });
            doc.text("This is a system-generated document. Hash: " + btoa(Date.now().toString()).substring(0, 8), 14, 200);
        }

        console.log("Saving PDF...");
        doc.save(`UBAPS_Report_${yearLabel.replace(/\//g, '-')}.pdf`);
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
        const optgroup = document.getElementById('optgroup-scheme');
        if (optgroup) {
            optgroup.innerHTML = "";
            schemes.forEach(scheme => {
                const option = document.createElement('option');
                option.value = `scheme:${scheme.scheme_name}`;
                option.textContent = scheme.scheme_name;
                optgroup.appendChild(option);
            });
        }
    } catch (err) {
        console.error("Error populating schemes:", err);
    }
}

function updatePayloadAndFetch() {
    let statusList = [];
    let statusLabel = "All Applicants";

    if (filter[0] && filter[0].checked) {
        statusList.push(filter[0].value);
        statusList.push("considering");
        statusLabel = "Bursary Applicants";
    }
    if (filter[1] && filter[1].checked) {
        statusList.push(filter[1].value);
        statusList.push("paid");
        statusLabel = "Selected Students";
    }
    if (filter[2] && filter[2].checked) {
        statusList.push(filter[2].value);
        statusLabel = "Not Selected Candidates";
    }

    const advanced = window.activeAdvancedFilters || {};

    currentPayload = {
        statuses: statusList,
        search: document.getElementById('filter-search')?.value || '',
        department: advanced.dept || [],
        scheme: advanced.scheme || [],
        parent: advanced.parent || [],
        employment: advanced.employment || [],
        gender: advanced.gender || [],
        date_start: document.getElementById('filter-date-start')?.value || '',
        date_end: document.getElementById('filter-date-end')?.value || '',
        cohort_start: document.getElementById('filter-cohort-start')?.value || '',
        cohort_end: document.getElementById('filter-cohort-end')?.value || ''
    };

    window.currentReportTitle = statusLabel;
    const titleEl = document.getElementById('current-report-title');
    if (titleEl) titleEl.firstChild.textContent = statusLabel + " ";

    fetchApplicants(currentPage);
    fetchReportStats();

    if (window.renderFilterTags) {
        window.renderFilterTags(currentPayload, () => {
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }
}


async function fetchApplicants(page) {
    const year = document.getElementById('period-year')?.value || '';
    const semester = document.getElementById('period-semester')?.value || '';
    const month = document.getElementById('period-month')?.value || '';

    try {
        const response = await fetch(`/applicants?page=${page}&limit=${limit}&year=${year}&semester=${semester}&month=${month}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(currentPayload)
        });

        const result = await response.json();
        window.currentApplicants = result.data || [];

        const countEl = document.getElementById('report-result-count');
        if (countEl) countEl.textContent = `Showing ${result.data?.length || 0} of ${result.total || 0} records`;

        const tbody = document.getElementById("judge_tbody");
        tbody.innerHTML = "";

        if (result.data && result.data.length > 0) {
            result.data.forEach(item => {
                const tr = document.createElement("tr");

                // Assign status-based accent classes
                const status = (item.status || '').toLowerCase();
                if (status === 'selected' || status === 'paid') {
                    tr.classList.add('row-selected');
                } else if (status === 'not selected') {
                    tr.classList.add('row-rejected');
                } else {
                    tr.classList.add('row-pending');
                }

                const schemeDisplay = (status === 'selected' || status === 'paid')
                    ? (item.scheme_name && item.scheme_name !== 'No Scheme' ? item.scheme_name : 'Not Assigned')
                    : 'Not Assigned';

                // Stamp hover data for the hover-card feature
                tr.dataset.hover = JSON.stringify({
                    name:       `${item.first_name} ${item.last_name}`,
                    reg:        item.registration_number,
                    status:     item.status || '',
                    gender:     item.gender || '—',
                    parent:     item.parent_guardian_status || '—',
                    employment: item.guardian_employment_status || '—',
                    support:    item.relative_support || '—',
                    scheme:     schemeDisplay
                });

                const icons = {
                    review: `<svg xmlns="http://www.w3.org/2000/svg" height="18" viewBox="0 -960 960 960" width="18" fill="currentColor"><path d="M480-320q75 0 127.5-52.5T660-500q0-75-52.5-127.5T480-680q-75 0-127.5 52.5T300-500q0 75 52.5 127.5T480-320Zm0-72q-45 0-76.5-31.5T372-500q0-45 31.5-76.5T480-608q45 0 76.5 31.5T588-500q0 45-31.5 76.5T480-392Zm0 192q-146 0-266-81.5T35-500q61-137 181-218.5T480-800q146 0 266 81.5T925-500q-61 137-181 218.5T480-200Zm0-300Zm0 220q113 0 207.5-59.5T832-500q-50-101-144.5-160.5T480-720q-113 0-207.5 59.5T128-500q50 101 144.5 160.5T480-280Z"/></svg>`,
                    comment: `<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>`
                };

                tr.innerHTML = `
                    <td>${item.first_name} ${item.last_name}</td>
                    <td><span class="reg-badge">${item.registration_number}</span></td>
                    <td>${new Date(item.application_date).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })}</td>
                    <td>${item.gender}</td>
                    <td>${item.parent_guardian_status}</td>
                    <td>${item.guardian_employment_status}</td>
                    <td>${item.relative_support}</td>
                    <td>${schemeDisplay}</td>
                    <td>
                      <div style="display: flex; gap: 6px;">
                        <button title="Review" class="review-btn" data-id="${item.registration_number}">
                          ${icons.review}
                          Review
                        </button>
                        <button title="Comment" class="comment-btn action-btn" data-id="${item.registration_number}">
                          ${icons.comment}
                        </button>
                      </div>
                    </td>
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

function initModalLogic() {
    const tbody = document.getElementById("judge_tbody");
    if (window.initCommentLogic) window.initCommentLogic();

    tbody.onclick = e => {
        const btn = e.target.closest(".action-btn, .comment-btn");
        if (!btn) return;

        const id = btn.getAttribute("data-id");
        const student = window.currentApplicants.find(a => a.registration_number === id);
        if (!student) return;
        window.activeStudent = student;

        if (btn.classList.contains("comment-btn")) {
            if (window.openCommentModal) window.openCommentModal(student);
        }
    };
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


async function fetchReportStats() {
    const year = document.getElementById('period-year')?.value || '';
    const semester = document.getElementById('period-semester')?.value || '';
    const month = document.getElementById('period-month')?.value || '';

    console.log('Fetching stats for:', { year, semester, month, currentPayload });

    try {
        const res = await fetch(`/api/reports/stats?year=${year}&semester=${semester}&month=${month}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(currentPayload || {})
        });
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

        const total = s.total_applications || 0;
        set('stat-total', total);
        set('stat-approved', s.approved || 0);
        set('stat-pending', s.pending || 0);
        set('stat-rejected', s.rejected || 0);

        if (total > 0) {
            set('rate-selection', ((s.approved / total) * 100).toFixed(1) + '%');
            set('rate-pending', ((s.pending / total) * 100).toFixed(1) + '%');
            set('rate-rejection', ((s.rejected / total) * 100).toFixed(1) + '%');
        } else {
            set('rate-selection', '—%');
            set('rate-pending', '—%');
            set('rate-rejection', '—%');
        }

        const grid = document.getElementById('department-insights-grid');
        if (grid) {
            grid.innerHTML = '';

            const deptNames = {
                'cen': 'Computer Engineering',
                'fsn': 'Food Security and Nutrients',
                'edu': 'Education in Humanities',
                'bph': 'Bachelors in Public Health'
            };

            const modalTitle = document.getElementById('dept-modal-title');
            if (modalTitle) {
                if (currentPayload && currentPayload.department && currentPayload.department.length > 0) {
                    if (currentPayload.department.length === 1) {
                        const rawName = currentPayload.department[0].toLowerCase();
                        const dName = deptNames[rawName] || currentPayload.department[0].toUpperCase();
                        modalTitle.textContent = `${dName} Breakdown`;
                    } else {
                        modalTitle.textContent = `Filtered Departments Breakdown (${currentPayload.department.length})`;
                    }
                } else {
                    modalTitle.textContent = `Overall Department Breakdown`;
                }
            }

            if (s.faculty_breakdown && s.faculty_breakdown.length > 0) {
                s.faculty_breakdown.forEach(dept => {
                    const dTotal = dept.count || 0;
                    const dSel = dept.selected_count || 0;

                    const selRate = dTotal > 0 ? ((dSel / dTotal) * 100).toFixed(1) : '0.0';
                    const shareApp = total > 0 ? ((dTotal / total) * 100).toFixed(1) : '0.0';
                    const shareSel = (s.approved || 0) > 0 ? ((dSel / s.approved) * 100).toFixed(1) : '0.0';

                    const rawName = (dept.faculty || '').toLowerCase();
                    const displayName = deptNames[rawName] ? `${deptNames[rawName]} (${dept.faculty.toUpperCase()})` : (dept.faculty || 'Unknown');

                    const card = document.createElement('div');
                    card.className = 'dept-card';
                    card.innerHTML = `
                        <h4>${displayName}</h4>
                        <div class="dept-stat">
                            <span class="dept-stat-label">Selection Rate</span>
                            <span class="dept-stat-val" style="color: #10b981;">${selRate}%</span>
                        </div>
                        <div class="dept-stat">
                            <span class="dept-stat-label">Share of Applicants</span>
                            <span class="dept-stat-val">${shareApp}%</span>
                        </div>
                        <div class="dept-stat">
                            <span class="dept-stat-label">Share of Selections</span>
                            <span class="dept-stat-val" style="color: #2563eb;">${shareSel}%</span>
                        </div>
                    `;
                    grid.appendChild(card);
                });
            }
        }

    } catch (err) {
        console.error('fetchReportStats error:', err);
    }
}
