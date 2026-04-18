let currentPage = 1;
const limit = 20;

document.addEventListener('DOMContentLoaded', () => {
    loadReports(currentPage);
});

async function loadReports(page) {
    try {
        const response = await fetch(`/api/reports/comprehensive?page=${page}&limit=${limit}`);
        if (!response.ok) throw new Error("Failed to fetch reports");
        
        const result = await response.json();
        const tbody = document.querySelector('#student-reports tbody');
        tbody.innerHTML = '';

        if (result.data && result.data.length > 0) {
            result.data.forEach(report => {
                const row = document.createElement('tr');
                row.className = report.status === 'Approved' ? 'success' : 'warning';

                row.innerHTML = `
                    <td>${report.name} ${report.surname}</td>
                    <td>${report.application_id}</td>
                    <td>Computer Science</td>
                    <td><span class="status-badge approved">${report.status}</span></td>
                    <td>MWK ${report.bursary_amount.toLocaleString()}</td>
                    <td>2025/2026</td>
                    <td>
                        <button class="view-btn">View Details</button>
                    </td>
                `;
                tbody.appendChild(row);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;">No reports found</td></tr>';
        }

        renderPagination(result.total, page);

    } catch (error) {
        console.error('Error loading reports:', error);
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

    for (let i = 1; i <= totalPages; i++) {
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

