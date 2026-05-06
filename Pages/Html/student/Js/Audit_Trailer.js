const tabs = document.querySelectorAll(".tabs button");
const tables = document.querySelectorAll(".log-table");
let currentTab = "applications";
const pages = {
    applications: 1,
    payments: 1,
    users: 1
};
const limit = 20;

document.addEventListener("DOMContentLoaded", () => {
    // Initial loads
    renderTab(currentTab);

    tabs.forEach(tab => {
        tab.addEventListener("click", () => {
            tabs.forEach(t => t.classList.remove("active"));
            tab.classList.add("active");

            currentTab = tab.dataset.tab;
            tables.forEach(table => table.classList.toggle("active", table.id === currentTab));
            
            renderTab(currentTab);
        });
    });
});

async function renderTab(tabName) {
    const endpointMap = {
        applications: "/applicationlog",
        payments: "/paymentlog",
        users: "/userlog"
    };

    const tbodyMap = {
        applications: "application-body",
        payments: "payment-body",
        users: "user-body"
    };

    const page = pages[tabName];
    const endpoint = endpointMap[tabName];
    const tbody = document.getElementById(tbodyMap[tabName]);

    try {
        const response = await fetch(`${endpoint}?page=${page}&limit=${limit}`);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        
        const result = await response.json();
        const logs = result.data || [];
        const total = result.total || 0;

        tbody.innerHTML = "";
        logs.forEach(log => {
            const row = document.createElement("tr");
            row.className = log.status === 'SUCCESS' ? 'success' : (log.status === 'FAILED' ? 'error' : 'warning');
            
            row.onclick = () => showLogDetails(log, tabName);

            if (tabName === "payments") {
                row.innerHTML = `
                    <td>${new Date(log.occurred_at).toLocaleString()}</td>
                    <td>${log.user_role}</td>
                    <td>${log.user_id || 'System'}</td>
                    <td>${log.action}</td>
                    <td>${log.target}</td>
                    <td>MWK ${log.amount ? log.amount.toLocaleString() : '0'}</td>
                    <td>${log.status}</td>
                    <td>${log.duration_ms}ms</td>
                `;
            } else {
                row.innerHTML = `
                    <td>${new Date(log.occurred_at).toLocaleString()}</td>
                    <td>${log.user_role}</td>
                    <td>${log.user_id || 'System'}</td>
                    <td>${log.action}</td>
                    <td>${log.target}</td>
                    <td>${log.status}</td>
                    <td>${log.duration_ms}ms</td>
                `;
            }
            tbody.appendChild(row);
        });

        if (logs.length === 0) {
            tbody.innerHTML = `<tr><td colspan="10" style="text-align:center;">No logs found</td></tr>`;
        }

        renderPaginationControls(tabName, total);

    } catch (error) {
        console.error(`Error loading ${tabName} logs:`, error);
        if (window.showToast) window.showToast(`Failed to load ${tabName} logs`, "error");
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

function renderPaginationControls(tabName, total) {
    const container = document.getElementById("pagination-controls");
    container.innerHTML = "";

    const totalPages = Math.ceil(total / limit);
    if (totalPages <= 1) return;

    const page = pages[tabName];

    const prevBtn = document.createElement("button");
    prevBtn.innerText = "Prev";
    prevBtn.disabled = page === 1;
    prevBtn.onclick = () => {
        pages[tabName]--;
        renderTab(tabName);
    };
    container.appendChild(prevBtn);

    const start = Math.max(1, page - 2);
    const end = Math.min(totalPages, page + 2);

    for (let i = start; i <= end; i++) {
        const btn = document.createElement("button");
        btn.innerText = i;
        if (i === page) btn.classList.add("active");
        btn.onclick = () => {
            pages[tabName] = i;
            renderTab(tabName);
        };
        container.appendChild(btn);
    }

    const nextBtn = document.createElement("button");
    nextBtn.innerText = "Next";
    nextBtn.disabled = page === totalPages;
    nextBtn.onclick = () => {
        pages[tabName]++;
        renderTab(tabName);
    };
    container.appendChild(nextBtn);
}