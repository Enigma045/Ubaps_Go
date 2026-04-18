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