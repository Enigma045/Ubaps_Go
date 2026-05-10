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
let currentMode = 'student'; // 'student' or 'scheme'
let schemeData = null;

document.addEventListener("DOMContentLoaded", function () {
    tabs = document.querySelectorAll(".tabs input");
    tables = document.querySelectorAll(".log-table");
    filter = document.querySelector(".tabs")?.querySelectorAll("input") || [];

    setupModeToggle();

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
            updatePayloadAndFetch(); 
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
            if (currentMode === 'student') {
                generatePDFReport();
            } else {
                generateSchemePDFReport();
            }
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
            window.activeAdvancedFilters = {};
            currentPage = 1;
            updatePayloadAndFetch();
        });
    }

    fetchReportStats();
});

function setupModeToggle() {
    const studentBtn = document.getElementById('mode-student');
    const schemeBtn = document.getElementById('mode-scheme');
    const studentView = document.getElementById('student-reports-view');
    const schemeView = document.getElementById('scheme-reports-view');
    const studentOnlyElements = document.querySelectorAll('.student-only');
    const filterBar = document.querySelector('.filter-bar');
    const tabsArea = document.querySelector('.tabs');
    const pagination = document.getElementById('pagination-controls');

    const switchMode = (mode) => {
        currentMode = mode;
        if (mode === 'student') {
            studentBtn.classList.add('active');
            schemeBtn.classList.remove('active');
            studentView.style.display = 'block';
            schemeView.style.display = 'none';
            filterBar.style.display = 'flex';
            tabsArea.style.display = 'flex';
            pagination.style.display = 'flex';
            studentOnlyElements.forEach(el => el.style.display = 'flex');
            
            updateMainTableHead('student');
            updatePayloadAndFetch();
        } else {
            schemeBtn.classList.add('active');
            studentBtn.classList.remove('active');
            schemeView.style.display = 'block';
            studentView.style.display = 'none';
            filterBar.style.display = 'none';
            tabsArea.style.display = 'none';
            pagination.style.display = 'none';
            studentOnlyElements.forEach(el => el.style.display = 'none');
            
            fetchSchemeReports();
        }
    };

    studentBtn.addEventListener('click', () => switchMode('student'));
    schemeBtn.addEventListener('click', () => switchMode('scheme'));
}

function updateMainTableHead(mode) {
    const thead = document.getElementById('main-table-head');
    if (!thead) return;

    if (mode === 'student') {
        thead.innerHTML = `
            <tr>
                <th>Student Name</th>
                <th>Reg. Number</th>
                <th>App. Date</th>
                <th>Sex</th>
                <th>Parent Status</th>
                <th>Employment</th>
                <th>Relative Support</th>
                <th>Bursary Scheme</th>
                <th>Action</th>
            </tr>
        `;
    } else {
        thead.innerHTML = `
            <tr>
                <th>Benefactor</th>
                <th>Scheme Name</th>
                <th>Total Fund</th>
                <th>Committed</th>
                <th>Remaining</th>
                <th>Usage %</th>
                <th>Applicants</th>
                <th>Status</th>
            </tr>
        `;
    }
}

async function fetchSchemeReports() {
    try {
        const res = await fetch("/api/reports/schemes");
        if (!res.ok) throw new Error("Failed to fetch scheme reports");
        schemeData = await res.json();
        
        updateMainTableHead('scheme');
        renderSchemeSummaryReport();
        
    } catch (err) {
        console.error("Error fetching scheme reports:", err);
        showToast("Failed to load scheme analytics", "error");
    }
}

function renderSchemeSummaryReport() {
    if (!schemeData || !schemeData.summary) return;
    const tbody = document.getElementById("judge_tbody");
    tbody.innerHTML = "";

    schemeData.summary.forEach(s => {
        const tr = document.createElement("tr");
        const usageClass = s.usage_percent > 90 ? 'color: #ef4444;' : (s.usage_percent > 70 ? 'color: #f59e0b;' : 'color: #2563eb;');
        const statusClass = s.status.toLowerCase() === 'active' ? 'background: #dcfce7; color: #166534;' : 'background: #fee2e2; color: #991b1b;';
        
        tr.innerHTML = `
            <td><strong>${s.benefactor_name}</strong></td>
            <td>${s.scheme_name}</td>
            <td>MWK ${s.total_fund.toLocaleString()}</td>
            <td>MWK ${s.committed.toLocaleString()}</td>
            <td>MWK ${s.remaining.toLocaleString()}</td>
            <td style="${usageClass} font-weight: bold;">${s.usage_percent.toFixed(1)}%</td>
            <td style="text-align: center;">${s.number_of_applicants}</td>
            <td><span class="status-pill" style="padding: 2px 8px; border-radius: 12px; font-size: 0.75rem; font-weight: 600; ${statusClass}">${s.status}</span></td>
        `;
        tbody.appendChild(tr);
    });
}

async function generateSchemePDFReport() {
    showToast("Generating Scheme Summary Report...", "info");
    try {
        const { jsPDF } = window.jspdf;
        const doc = new jsPDF('l', 'mm', 'a4');

        doc.setFontSize(22);
        doc.setTextColor(30, 41, 59);
        doc.text("UBAPS", 14, 20);
        doc.setFontSize(10);
        doc.setTextColor(37, 99, 235);
        doc.text("Comprehensive Bursary Scheme Summary Report", 44, 19);
        doc.setDrawColor(226, 232, 240);
        doc.line(14, 25, 283, 25);
        
        doc.setFontSize(12);
        doc.setTextColor(30, 41, 59);
        doc.text(`Generated: ${new Date().toLocaleString()}`, 283, 32, { align: 'right' });

        const headers = [['Benefactor', 'Scheme Name', 'Total Fund', 'Committed', 'Remaining', 'Usage %', 'Applicants', 'Status']];
        const body = schemeData.summary.map(s => [
            s.benefactor_name,
            s.scheme_name,
            s.total_fund.toLocaleString(),
            s.committed.toLocaleString(),
            s.remaining.toLocaleString(),
            s.usage_percent.toFixed(1) + '%',
            s.number_of_applicants,
            s.status
        ]);

        doc.autoTable({
            startY: 40,
            head: headers,
            body: body,
            theme: 'grid',
            headStyles: { fillColor: [37, 99, 235] },
            styles: { fontSize: 8 }
        });

        doc.save(`UBAPS_Scheme_Summary_Report.pdf`);
        showToast("Scheme Summary Report Exported!", "success");
    } catch (err) {
        console.error(err);
        showToast("PDF Export failed", "error");
    }
}

// Student Report Logic
async function generatePDFReport() {
    const yearVal = document.getElementById('period-year')?.value || '';
    const semesterVal = document.getElementById('period-semester')?.value || '';
    const monthVal = document.getElementById('period-month')?.value || '';

    showToast("Generating PDF Report...", "info");

    try {
        const url = `/api/reports/comprehensive?year=${yearVal}&semester=${semesterVal}&month=${monthVal}&limit=5000`;
        const res = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(currentPayload)
        });

        const json = await res.json();
        const data = json.data || [];

        const { jsPDF } = window.jspdf;
        const doc = new jsPDF('l', 'mm', 'a4');

        doc.setFontSize(22);
        doc.text("UBAPS", 14, 20);
        doc.setFontSize(10);
        doc.text("University Bursary Application & Payment System", 44, 19);
        doc.line(14, 25, 283, 25);

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

        doc.autoTable({
            startY: 55,
            head: [['Reg. Number', 'App. Date', 'Student Name', 'Sex', 'Department', 'Status', 'Bursary Scheme', 'Amount (MWK)']],
            body: tableData,
            theme: 'grid',
            headStyles: { fillColor: [30, 41, 59] }
        });

        doc.save(`UBAPS_Student_Report.pdf`);
        showToast("Student Report Exported!", "success");
    } catch (err) {
        showToast("PDF Export failed", "error");
    }
}

async function populateSchemes() {
    try {
        const res = await fetch("/getschemes");
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
    } catch (err) {}
}

function updatePayloadAndFetch() {
    let statusList = [];
    if (filter[0]?.checked) { statusList.push(filter[0].value); statusList.push("considering"); }
    if (filter[1]?.checked) statusList.push(filter[1].value);
    if (filter[2]?.checked) statusList.push(filter[2].value);
    if (filter[3]?.checked) statusList.push(filter[3].value);

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

    fetchApplicants(currentPage);
    fetchReportStats();
}

async function fetchApplicants(page) {
    if (currentMode === 'scheme') return;
    const year = document.getElementById('period-year')?.value || '';
    const semester = document.getElementById('period-semester')?.value || '';
    const month = document.getElementById('period-month')?.value || '';

    try {
        const response = await fetch(`/applicants?page=${page}&limit=${limit}&year=${year}&semester=${semester}&month=${month}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(currentPayload)
        });
        const result = await response.json();
        const tbody = document.getElementById("judge_tbody");
        tbody.innerHTML = "";

        if (result.data?.length > 0) {
            result.data.forEach(item => {
                const tr = document.createElement("tr");
                const status = (item.status || '').toLowerCase();
                const schemeDisplay = (status === 'selected' || status === 'paid') ? (item.scheme_name || 'Not Assigned') : 'Not Assigned';

                tr.innerHTML = `
                    <td>${item.first_name} ${item.last_name}</td>
                    <td><span class="reg-badge">${item.registration_number}</span></td>
                    <td>${new Date(item.application_date).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })}</td>
                    <td>${item.gender}</td>
                    <td>${item.parent_guardian_status}</td>
                    <td>${item.guardian_employment_status}</td>
                    <td>${item.relative_support}</td>
                    <td>${schemeDisplay}</td>
                    <td><button class="review-btn" data-id="${item.registration_number}">Review</button></td>
                `;
                tbody.appendChild(tr);
            });
        }
        renderPagination(result.total, page);
    } catch (error) {}
}

function renderPagination(total, page) {
    const container = document.getElementById("pagination-controls");
    if (!container) return;
    container.innerHTML = "";
    const totalPages = Math.ceil(total / limit);
    if (totalPages <= 1) return;

    for (let i = 1; i <= totalPages; i++) {
        const btn = document.createElement("button");
        btn.innerText = i;
        if (i === page) btn.classList.add("active");
        btn.onclick = () => { currentPage = i; fetchApplicants(currentPage); };
        container.appendChild(btn);
    }
}

async function fetchReportStats() {
    const year = document.getElementById('period-year')?.value || '';
    const semester = document.getElementById('period-semester')?.value || '';
    const month = document.getElementById('period-month')?.value || '';

    try {
        const res = await fetch(`/api/reports/stats?year=${year}&semester=${semester}&month=${month}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(currentPayload)
        });
        const s = await res.json();
        const set = (id, val) => { if (document.getElementById(id)) document.getElementById(id).textContent = val; };
        set('stat-total', s.total_applications || 0);
        set('stat-approved', s.approved || 0);
        set('stat-paid', s.paid || 0);
        set('stat-pending', s.pending || 0);
        set('stat-rejected', s.rejected || 0);
    } catch (err) {}
}
