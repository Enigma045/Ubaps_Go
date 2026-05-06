let currentPage = 1;
const limit = 20;
let currentType = 'applications';

document.addEventListener("DOMContentLoaded", () => {
    loadLogs(currentType, currentPage);

    document.getElementById("logType").addEventListener("change", (e) => {
        currentType = e.target.value;
        currentPage = 1;
        loadLogs(currentType, currentPage);
    });

    document.getElementById("apply-filters").addEventListener("click", () => {
        currentPage = 1;
        loadLogs(currentType, currentPage);
    });

    document.getElementById("export-report").addEventListener("click", () => {
        exportComprehensiveReport();
    });

    // Set current time
    const now = new Date();
    document.getElementById("current-time").innerText = now.toLocaleString();
});

async function loadLogs(type, page) {
    const tbody = document.getElementById("table-body");
    const thead = document.getElementById("table-head");
    
    // Set headers
    if (type === 'applications') {
        thead.innerHTML = `
            <tr>
                <th>Time</th>
                <th>Role</th>
                <th>User</th>
                <th>Action</th>
                <th>Target</th>
                <th>Status</th>
                <th>Duration</th>
            </tr>
        `;
    } else if (type === 'payments') {
        thead.innerHTML = `
            <tr>
                <th>Time</th>
                <th>Role</th>
                <th>User</th>
                <th>Action</th>
                <th>Target</th>
                <th>Amount</th>
                <th>Status</th>
                <th>Duration</th>
            </tr>
        `;
    } else if (type === 'users') {
        thead.innerHTML = `
            <tr>
                <th>Time</th>
                <th>Role</th>
                <th>User</th>
                <th>Action</th>
                <th>Target</th>
                <th>Status</th>
                <th>Duration</th>
            </tr>
        `;
    }

    try {
        const endpoint = type === 'applications' ? '/applicationlog' : (type === 'payments' ? '/paymentlog' : '/userlog');
        const response = await fetch(`${endpoint}?page=${page}&limit=${limit}`);
        const result = await response.json();
        
        tbody.innerHTML = '';
        
        if (result.data && result.data.length > 0) {
            result.data.forEach(log => {
                const row = document.createElement("tr");
                row.className = log.status === 'SUCCESS' ? 'success' : (log.status === 'FAILED' ? 'error' : 'warning');
                
                row.onclick = () => showLogDetails(log, type);

                if (type === 'payments') {
                    row.innerHTML = `
                        <td>${new Date(log.occurred_at).toLocaleString()}</td>
                        <td>${log.user_role}</td>
                        <td>User #${log.user_id || 'System'}</td>
                        <td>${log.action}</td>
                        <td>${log.target}</td>
                        <td>MWK ${log.amount.toLocaleString()}</td>
                        <td>${log.status}</td>
                        <td>${log.duration_ms}ms</td>
                    `;
                } else {
                    row.innerHTML = `
                        <td>${new Date(log.occurred_at).toLocaleString()}</td>
                        <td>${log.user_role}</td>
                        <td>User #${log.user_id || 'System'}</td>
                        <td>${log.action}</td>
                        <td>${log.target}</td>
                        <td>${log.status}</td>
                        <td>${log.duration_ms}ms</td>
                    `;
                }
                tbody.appendChild(row);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="10" style="text-align:center;">No logs found</td></tr>';
        }

        renderPagination(result.total, page);

    } catch (error) {
        console.error("Error loading logs:", error);
        showToast("Error loading logs", "error");
    }
}

function showLogDetails(log, type) {
    const modal = document.getElementById("log-detail-modal");
    const body = document.getElementById("modal-details-body");
    const headerTitle = document.querySelector("#modal-header h3");
    
    let title = "Log Details";
    let specializedContent = "";

    if (type === 'applications') {
        title = "Application Audit Detail";
        specializedContent = `
            <div class="detail-row highlight">
                <span class="detail-label">Target Application ID</span>
                <span class="detail-value">${log.application || 'N/A'}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Target User ID</span>
                <span class="detail-value">#${log.target_user_id || 'N/A'}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Action Performed</span>
                <span class="detail-value">${log.action}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Associated Scheme</span>
                <span class="detail-value">${log.target.includes('Scheme:') ? log.target.split('Scheme:')[1].split(',')[0] : 'Standard Processing'}</span>
            </div>
        `;
    } else if (type === 'payments') {
        title = "Payment Disbursement Detail";
        specializedContent = `
            <div class="detail-row highlight-amount">
                <span class="detail-label">Amount Processed</span>
                <span class="detail-value">MWK ${log.amount ? log.amount.toLocaleString() : '0'}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Target Application ID</span>
                <span class="detail-value">${log.application || 'N/A'}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Target User ID</span>
                <span class="detail-value">#${log.target_user_id || 'N/A'}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Payment Action</span>
                <span class="detail-value">${log.action}</span>
            </div>
        `;
    } else {
        title = "System & User Activity";
        specializedContent = `
            <div class="detail-row">
                <span class="detail-label">System Action</span>
                <span class="detail-value">${log.action}</span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Target Entity</span>
                <span class="detail-value">${log.target}</span>
            </div>
            ${log.target_user_id ? `
            <div class="detail-row">
                <span class="detail-label">Target User ID</span>
                <span class="detail-value">#${log.target_user_id}</span>
            </div>` : ''}
        `;
    }

    headerTitle.innerText = title;

    body.innerHTML = `
        ${specializedContent}
        <div class="detail-row separator"></div>
        <div class="detail-row">
            <span class="detail-label">Occurred At</span>
            <span class="detail-value">${new Date(log.occurred_at).toLocaleString()}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Performing Role</span>
            <span class="detail-value">${log.user_role}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Admin ID</span>
            <span class="detail-value">#${log.user_id || 'System'}</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">Execution Time</span>
            <span class="detail-value">${log.duration_ms}ms</span>
        </div>
        <div class="detail-row">
            <span class="detail-label">System Status</span>
            <span class="status-pill ${log.status}">${log.status}</span>
        </div>
    `;

    modal.style.display = 'flex';

    // Close logic
    const closeBtn = document.getElementById("close-modal");
    closeBtn.onclick = () => modal.style.display = 'none';
    window.onclick = (event) => {
        if (event.target == modal) modal.style.display = 'none';
    };
}

function renderPagination(total, page) {
    const container = document.getElementById("pagination-controls");
    container.innerHTML = '';
    
    const totalPages = Math.ceil(total / limit);
    if (totalPages <= 1) return;

    const prevBtn = document.createElement("button");
    prevBtn.innerText = "Prev";
    prevBtn.disabled = page === 1;
    prevBtn.onclick = () => {
        currentPage--;
        loadLogs(currentType, currentPage);
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
            loadLogs(currentType, currentPage);
        };
        container.appendChild(btn);
    }

    const nextBtn = document.createElement("button");
    nextBtn.innerText = "Next";
    nextBtn.disabled = page === totalPages;
    nextBtn.onclick = () => {
        currentPage++;
        loadLogs(currentType, currentPage);
    };
    container.appendChild(nextBtn);
}

async function exportComprehensiveReport() {
    try {
        const response = await fetch('/api/reports/comprehensive?limit=1000'); // Get a large chunk for export
        const result = await response.json();
        
        if (!result.data || result.data.length === 0) {
            showToast("No data to export", "warning");
            return;
        }

        // Convert JSON to CSV for download
        const data = result.data;
        const csvRows = [];
        const headers = Object.keys(data[0]);
        csvRows.push(headers.join(','));

        for (const row of data) {
            const values = headers.map(header => {
                const val = row[header];
                return `"${val}"`;
            });
            csvRows.push(values.join(','));
        }

        const csvContent = csvRows.join('\n');
        const blob = new Blob([csvContent], { type: 'text/csv' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.setAttribute('hidden', '');
        a.setAttribute('href', url);
        a.setAttribute('download', `comprehensive_report_${new Date().toISOString().slice(0,10)}.csv`);
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        
        showToast("Report exported successfully", "success");
    } catch (error) {
        console.error("Export error:", error);
        showToast("Failed to export report", "error");
    }
}

function showToast(message, type) {
    if (window.showToast) {
        window.showToast(message, type);
    } else {
        alert(message);
    }
}
