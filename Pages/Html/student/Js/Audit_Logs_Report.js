let currentPage = 1;
const limit = 20;
let selectedTypes = ['applications', 'payments', 'users'];

document.addEventListener("DOMContentLoaded", () => {
    setupMultiSelect();
    loadLogs(currentPage);

    document.getElementById("apply-filters").addEventListener("click", () => {
        currentPage = 1;
        loadLogs(currentPage);
    });

    const exportBtn = document.getElementById("export-pdf-btn");
    if (exportBtn) {
        exportBtn.addEventListener("click", () => {
            generatePDFReport();
        });
    }

    // Close modal on click outside
    const modal = document.getElementById("log-detail-modal");
    if (modal) {
        window.onclick = (event) => {
            if (event.target == modal) modal.style.display = 'none';
        };
        const closeBtn = document.getElementById("close-modal");
        if (closeBtn) {
            closeBtn.onclick = () => {
                modal.style.display = 'none';
            };
        }
    }
});

function setupMultiSelect() {
    const selectBox = document.getElementById('selectBox');
    const optionsList = document.getElementById('optionsList');
    const checkboxes = document.querySelectorAll('.type-checkbox');
    const countText = document.getElementById('selectedCountText');

    if (!selectBox) return;

    selectBox.addEventListener('click', (e) => {
        e.stopPropagation();
        optionsList.classList.toggle('show');
    });

    document.addEventListener('click', () => {
        optionsList.classList.remove('show');
    });

    optionsList.addEventListener('click', (e) => {
        e.stopPropagation();
    });

    checkboxes.forEach(cb => {
        cb.addEventListener('change', () => {
            updateSelectedTypes();
        });
    });

    function updateSelectedTypes() {
        selectedTypes = Array.from(checkboxes)
            .filter(cb => cb.checked)
            .map(cb => cb.value);

        if (selectedTypes.length === 0) {
            countText.innerText = "Select Type";
        } else if (selectedTypes.length === 3) {
            countText.innerText = "All Categories";
        } else {
            countText.innerText = `${selectedTypes.length} Selected`;
        }
    }
}

async function loadLogs(page) {
    const tbody = document.getElementById("table-body");
    const thead = document.getElementById("table-head");

    // Get filter values
    const userId = document.getElementById("userId").value;
    const startDate = document.getElementById("startDate").value;
    const endDate = document.getElementById("endDate").value;

    // Set Unified Header
    thead.innerHTML = `
        <tr>
            <th>Time</th>
            <th>Role</th>
            <th>Performer</th>
            <th>Action</th>
            <th>Target Details</th>
            <th>Status</th>
            <th>Execution</th>
        </tr>
    `;

    try {
        const typesQuery = selectedTypes.join(',');
        let url = `/api/audit/all?page=${page}&limit=${limit}&types=${typesQuery}`;

        if (userId) url += `&user_id=${userId}`;
        if (startDate) url += `&start_date=${startDate}`;
        if (endDate) url += `&end_date=${endDate}`;

        const response = await fetch(url);
        const result = await response.json();

        tbody.innerHTML = "";

        if (result.data && result.data.length > 0) {
            result.data.forEach(log => {
                const tr = document.createElement("tr");
                const time = new Date(log.occurred_at).toLocaleString('en-GB', {
                    day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit'
                });

                // Determine what to show in Target Details
                let details = log.target || '—';
                if (log.amount && log.amount.Valid) {
                    details = `MWK ${log.amount.Float64.toLocaleString()}`;
                } else if (log.application && log.application.Valid) {
                    details = `App: ${log.application.String}`;
                }

                const statusClass = log.status === 'SUCCESS' ? 'success' : 'error';
                const performer = log.user_id ? `Admin #${log.user_id}` : 'SYSTEM';

                tr.innerHTML = `
                    <td><strong>${time}</strong></td>
                    <td><span class="role-badge">${log.user_role}</span></td>
                    <td>${performer}</td>
                    <td>${log.action}</td>
                    <td><code>${details}</code></td>
                    <td><span class="status-pill ${statusClass}">${log.status}</span></td>
                    <td>${log.duration_ms}ms</td>
                `;

                tr.onclick = () => showLogDetail(log);
                tbody.appendChild(tr);
            });

            updateAnalytics(result.data, result.total);
            renderPagination(result.total, page);
        } else {
            tbody.innerHTML = "<tr><td colspan='7' style='text-align:center; padding: 40px;'>No logs found matching your criteria.</td></tr>";
            document.getElementById("pagination-controls").innerHTML = "";
            document.getElementById("stat-total").innerText = 0;
            document.getElementById("stat-success").innerText = 0;
            document.getElementById("stat-failed").innerText = 0;
            document.getElementById("stat-rate").innerText = '0%';
        }
    } catch (error) {
        console.error("Error loading logs:", error);
    }
}

function updateAnalytics(logs, totalCount) {
    // We only have the current page here, but for simple analytics we use what we have or fetch separately.
    // Ideally the backend returns total success/fail counts.
    const total = logs.length;
    const success = logs.filter(l => l.status === 'SUCCESS').length;
    const failed = total - success;
    const rate = total > 0 ? Math.round((success / total) * 100) : 0;

    document.getElementById("stat-total").innerText = totalCount || total;
    document.getElementById("stat-success").innerText = success; // Note: this is per page, ideally get from backend
    document.getElementById("stat-failed").innerText = failed;
    document.getElementById("stat-rate").innerText = rate + '%';
}

function renderPagination(total, page) {
    const container = document.getElementById("pagination-controls");
    container.innerHTML = "";
    const totalPages = Math.ceil(total / limit);

    if (totalPages <= 1) return;

    // Show at most 5 pages or dots
    for (let i = 1; i <= totalPages; i++) {
        const btn = document.createElement("button");
        btn.innerText = i;
        if (i === page) btn.classList.add("active");
        btn.onclick = () => {
            currentPage = i;
            loadLogs(currentPage);
        };
        container.appendChild(btn);
    }
}

function showLogDetail(log) {
    const modal = document.getElementById("log-detail-modal");
    const body = document.getElementById("modal-details-body");

    let extra = '';
    if (log.application && log.application.Valid) extra += `<div class="detail-row"><span class="detail-label">Application</span><span class="detail-value">${log.application.String}</span></div>`;
    if (log.amount && log.amount.Valid) extra += `<div class="detail-row highlight-amount"><span class="detail-label">Amount</span><span class="detail-value">MWK ${log.amount.Float64.toLocaleString()}</span></div>`;
    if (log.target_user_id) extra += `<div class="detail-row"><span class="detail-label">Target User ID</span><span class="detail-value">${log.target_user_id}</span></div>`;

    body.innerHTML = `
        <div class="detail-row highlight">
            <span class="detail-label">Action Performed</span>
            <span class="detail-value">${log.action}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Occurred At</span>
            <span class="detail-value">${new Date(log.occurred_at).toLocaleString()}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Performer Role</span>
            <span class="detail-value">${log.user_role}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Performer ID</span>
            <span class="detail-value">${log.user_id || 'System'}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Entity Target</span>
            <span class="detail-value">${log.target}</span>
        </div>
        ${extra}
        <div class="detail-row">
            <span class="detail-label">Status</span>
            <span class="detail-value"><span class="status-pill ${log.status === 'SUCCESS' ? 'SUCCESS' : 'FAILED'}">${log.status}</span></span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Execution Time</span>
            <span class="detail-value">${log.duration_ms}ms</span>
        </div>
    `;

    modal.style.display = 'flex';
}

async function generatePDFReport() {
    if (!window.jspdf) {
        alert("PDF library not loaded yet.");
        return;
    }

    const { jsPDF } = window.jspdf;
    const doc = new jsPDF('l', 'mm', 'a4');

    // Title
    doc.setFontSize(22);
    doc.setTextColor(30, 41, 59);
    doc.text("UBAPS Audit Log Report", 14, 20);

    doc.setFontSize(10);
    doc.setTextColor(100, 116, 139);
    doc.text(`Generated on: ${new Date().toLocaleString()}`, 14, 28);
    doc.text(`Categories: ${selectedTypes.join(', ')}`, 14, 34);

    doc.line(14, 38, 283, 38);

    // Fetch all logs for the report (up to 500)
    const typesQuery = selectedTypes.join(',');
    const userId = document.getElementById("userId").value;
    const startDate = document.getElementById("startDate").value;
    const endDate = document.getElementById("endDate").value;

    let url = `/api/audit/all?page=1&limit=500&types=${typesQuery}`;
    if (userId) url += `&user_id=${userId}`;
    if (startDate) url += `&start_date=${startDate}`;
    if (endDate) url += `&end_date=${endDate}`;

    try {
        const response = await fetch(url);
        const result = await response.json();
        const logs = result.data || [];

        const tableData = logs.map(l => [
            new Date(l.occurred_at).toLocaleString(),
            l.user_role,
            l.user_id || 'SYSTEM',
            l.action,
            l.target + (l.amount && l.amount.Valid ? ` (MWK ${l.amount.Float64})` : '') + (l.application && l.application.Valid ? ` (App: ${l.application.String})` : ''),
            l.status,
            l.duration_ms + 'ms'
        ]);

        doc.autoTable({
            startY: 45,
            head: [['Time', 'Role', 'User', 'Action', 'Target/Amount', 'Status', 'Duration']],
            body: tableData,
            theme: 'grid',
            headStyles: { fillColor: [15, 23, 42] },
            styles: { fontSize: 8 }
        });

        doc.save(`UBAPS_Audit_Report_${new Date().getTime()}.pdf`);
    } catch (error) {
        console.error("PDF Export failed:", error);
        alert("Failed to generate PDF report.");
    }
}
