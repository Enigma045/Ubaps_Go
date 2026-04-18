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
