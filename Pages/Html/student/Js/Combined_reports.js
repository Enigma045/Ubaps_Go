let currentPage = 1;
const limit = 20;

document.addEventListener('DOMContentLoaded', () => {
    loadReports(currentPage);

    const searchInput = document.querySelector('.search-bar input');
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            // Future: Implement client-side or server-side search
        });
    }
});

async function loadReports(page) {
    try {
        const response = await fetch(`/api/reports/comprehensive?page=${page}&limit=${limit}`);
        if (!response.ok) throw new Error("Failed to fetch reports");
        
        const result = await response.json();
        const tbody = document.querySelector('#combined-reports tbody');
        tbody.innerHTML = '';

        if (result.data && result.data.length > 0) {
            result.data.forEach(report => {
                const row = document.createElement('tr');
                row.className = report.status === 'Approved' ? 'success' : 'warning';

                row.innerHTML = `
                    <td>${report.name} ${report.surname}</td>
                    <td>${report.application_id}</td>
                    <td>${report.scheme_name}</td>
                    <td>MWK ${report.bursary_amount.toLocaleString()}</td>
                    <td><span class="status-badge approved">${report.status}</span></td>
                    <td><span class="status-badge ${report.request_status === 'Approved' ? 'disbursed' : 'pending'}">${report.request_status}</span></td>
                    <td>${new Date(report.applied_at).toLocaleDateString()}</td>
                    <td>
                        <button class="view-btn">View Details</button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="10" style="text-align:center;">No reports found</td></tr>';
        }

        renderPagination(result.total, page);

    } catch (error) {
        console.error('Error loading reports:', error);
        if (window.showToast) showToast("Failed to load reports.", "error");
    }
}

function renderPagination(total, page) {
    const container = document.getElementById("pagination-controls");
    if (!container) return;
    container.innerHTML = '';
    
    const totalPages = Math.ceil(total / limit);
    if (totalPages <= 1) return;

    const prevBtn = document.createElement("button");
    prevBtn.innerText = "Prev";
    prevBtn.disabled = page === 1;
    prevBtn.onclick = () => {
        currentPage--;
        loadReports(currentPage);
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
            loadReports(currentPage);
        };
        container.appendChild(btn);
    }

    const nextBtn = document.createElement("button");
    nextBtn.innerText = "Next";
    nextBtn.disabled = page === totalPages;
    nextBtn.onclick = () => {
        currentPage++;
        loadReports(currentPage);
    };
    container.appendChild(nextBtn);
}

function applyFilters() {
    currentPage = 1;
    loadReports(currentPage);
}

async function exportReport() {
    try {
        const response = await fetch(`/api/reports/comprehensive?limit=1000`);
        const result = await response.json();
        
        if (!result.data || result.data.length === 0) {
            if (window.showToast) showToast("No data to export", "warning");
            return;
        }

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
        a.setAttribute('download', `report_${new Date().toISOString().slice(0,10)}.csv`);
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        
        if (window.showToast) showToast("Report exported successfully", "success");
    } catch (error) {
        console.error('Export error:', error);
        if (window.showToast) showToast("Failed to export report", "error");
    }
}

