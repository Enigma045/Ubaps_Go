let currentPage = 1;
const limit = 20;

document.addEventListener("DOMContentLoaded", () => {
    loadSchemes(currentPage);

    document.getElementById("scheme_info").onsubmit = async e => {
        e.preventDefault();
        const formData = new FormData(e.target);

        const res = await fetch("/benefactor", {
            method: "POST",
            body: formData,
            credentials: "include"
        });

        if (res.ok) {
            showToast("Bursary scheme created successfully!", "success");
            setTimeout(() => {
                location.reload();
            }, 2000);
        } else {
            const text = await res.text();
            showToast(text || "Failed to create scheme.", "error");
        }
    };
});

async function loadSchemes(page) {
    try {
        const res = await fetch(`/getbenefactor?page=${page}&limit=${limit}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        });
        
        const result = await res.json();
        displayBenefactor(result.data);
        renderPagination(result.total, page);
    } catch (err) {
        console.error(err);
        if (window.showToast) showToast("Failed to load bursary schemes.", "error");
    }
}

function displayBenefactor(data) {
    const tbody = document.getElementById("bursarytbody");
    tbody.innerHTML = "";
    if (!data || data.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;">No schemes found</td></tr>';
        return;
    }
    data.forEach(item => {
        const tr = document.createElement("tr");
        tr.innerHTML = `
            <td>🟡</td>
            <td class="name">${item.scheme_name}</td>
            <td>${item.benefactor_name}</td>
            <td class="email">${item.benefactor_email}</td>
            <td>${item.total_fund_amount}</td>
            <td>${item.available_balance}</td>
            <td>
                <button title="Review"><img class="action" src="/Image/svgviewer-output (18).svg" alt=""></button>
                <button title="delete" class="openModal2"><img class="action" src="/Image/svgviewer-output (21).svg" alt=""></button>
            </td>
        `;
        tbody.appendChild(tr);
    });
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
        loadSchemes(currentPage);
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
            loadSchemes(currentPage);
        };
        container.appendChild(btn);
    }

    const nextBtn = document.createElement("button");
    nextBtn.innerText = "Next";
    nextBtn.disabled = page === totalPages;
    nextBtn.onclick = () => {
        currentPage++;
        loadSchemes(currentPage);
    };
    container.appendChild(nextBtn);
}