let currentPage = 1;
const limit = 20;

document.addEventListener("DOMContentLoaded", function () {
    loadUsers(currentPage);
});

async function loadUsers(page) {
    try {
        const response = await fetch(`/getuserdetails?page=${page}&limit=${limit}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        const result = await response.json();
        const tbody = document.querySelector(".tbody");
        tbody.innerHTML = "";

        if (result.data && result.data.length > 0) {
            result.data.forEach(user => {
                const tr = document.createElement("tr");
                tr.dataset.user = JSON.stringify(user);

                const name = document.createElement("td");
                name.className = "name";
                name.textContent = user.first + " " + user.last;

                const role = document.createElement("td");
                role.textContent = user.role;

                const email = document.createElement("td");
                email.className = "email";
                email.textContent = user.email;

                const phone = document.createElement("td");
                phone.textContent = user.phone;

                const actions = document.createElement("td");

                const edit = document.createElement("button");
                edit.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#000"><path d="M200-200h57l391-391-57-57-391 391v57Zm-80 80v-170l528-527q12-11 26.5-17t30.5-6q16 0 31 6t26 18l55 56q12 11 17.5 26t5.5 30q0 16-5.5 30.5T817-647L290-120H120Zm640-584-56-56 56 56Zm-141 85-28-29 57 57-29-28Z"/></svg>`;
                // Add click listener for edit (assuming Edit_User.js handles it or we need to trigger it)
                edit.onclick = () => { if (typeof openEditModal === 'function') openEditModal(user); };

                const del = document.createElement("button");
                del.className = "openModal";
                del.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#000"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg>`;
                // Add click listener for delete (assuming Test_Delete.js handles it or we need to trigger it)
                del.onclick = () => { if (typeof openDeleteModal === 'function') openDeleteModal(user.email); };

                actions.appendChild(edit);
                actions.appendChild(del);

                tr.appendChild(name);
                tr.appendChild(role);
                tr.appendChild(email);
                tr.appendChild(phone);
                tr.appendChild(actions);

                tbody.appendChild(tr);
            });
        } else {
            tbody.innerHTML = '<tr><td colspan="5" style="text-align:center;">No users found</td></tr>';
        }

        renderPagination(result.total, page);

    } catch (error) {
        console.error("Error loading users:", error);
        if (window.showToast) showToast("Failed to load user details.", "error");
    }
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
        loadUsers(currentPage);
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
            loadUsers(currentPage);
        };
        container.appendChild(btn);
    }

    const nextBtn = document.createElement("button");
    nextBtn.innerText = "Next";
    nextBtn.disabled = page === totalPages;
    nextBtn.onclick = () => {
        currentPage++;
        loadUsers(currentPage);
    };
    container.appendChild(nextBtn);
}



