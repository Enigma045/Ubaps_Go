/* ============================================================
   Role Selector Widget
   A draggable, floating role selector for quick user navigation.
   ============================================================ */

// ─── Passwords Array ──────────────────────────────────────────
// Place your role passwords here. Index corresponds to the role
// order in ROLE_MAP below.
const ROLE_PASSWORDS = [
    "Student@2026",   // student
    "Rich@2002",   // registrarle
    "Rich@2002248",   // adminf
    "Dean@2002",   // dean_of_student
    "Faculty@2002",   // dean_of_facult
    "Science@2002",   // dean_of_science
    "Rich@2002",   // finance_officer
];

// ─── Emails Array ─────────────────────────────────────────────
// Place the email for each role here (same index as ROLE_PASSWORDS).
const ROLE_EMAILS = [
    "student@unilia.ac.mw",                           // student
    "lemillion045@gmail.com",                             // registrar
    "franciessambo@gmail.com",                             // admin
    "dean@unilia.ac.mw",                             // dean_of_student
    "faculty@unilia.ac.mw",                             // dean_of_facult
    "science@unilia.ac.mw",                             // dean_of_science
    "richardsambo@gmail.com",                             // finance_office
];

// ─── Role → URL Map ──────────────────────────────────────────
const ROLE_MAP = [
    { label: "Student", value: "student", url: "/dashboard" },
    { label: "Registrar", value: "registrar", url: "/registrardashboard" },
    { label: "Admin", value: "admin", url: "/admindashboard" },
    { label: "Dean of Students", value: "dean_of_student", url: "/deandashboard" },
    { label: "Dean of Faculty", value: "dean_of_facult", url: "/deandashboard" },
    { label: "Dean of Science", value: "dean_of_science", url: "/deandashboard" },
    { label: "Finance Office", value: "finance_office", url: "/financialdashboard" },
];

// ─── Build Widget DOM ────────────────────────────────────────
document.addEventListener("DOMContentLoaded", () => {
    const widget = document.createElement("div");
    widget.id = "role-selector-widget";
    widget.innerHTML = `
    <div class="rs-header" id="rs-drag-handle">
      <span class="rs-title">⚙ Switch Role</span>
      <button class="rs-toggle" id="rs-toggle-btn" title="Collapse / Expand">−</button>
    </div>
    <div class="rs-body" id="rs-body">
      <select class="rs-select" id="rs-role-select">
        <option value="" disabled selected>Select a role…</option>
      </select>
      <button class="rs-go-btn" id="rs-go-btn">Go →</button>
    </div>
  `;
    document.body.appendChild(widget);

    // Populate select options
    const select = document.getElementById("rs-role-select");
    ROLE_MAP.forEach((role) => {
        const opt = document.createElement("option");
        opt.value = role.value;
        opt.textContent = role.label;
        select.appendChild(opt);
    });

    // ─── Navigate on "Go" (Login via /Authorize) ──────────────
    document.getElementById("rs-go-btn").addEventListener("click", () => {
        const selected = select.value;
        if (!selected) return;

        const roleIndex = ROLE_MAP.findIndex((r) => r.value === selected);
        if (roleIndex === -1) return;

        const email = ROLE_EMAILS[roleIndex];
        const password = ROLE_PASSWORDS[roleIndex];

        if (!email || !password) {
            console.warn("Role Selector: No email/password configured for", selected);
            showToast("No credentials configured for this role.", "error");
            return;
        }

        // Build FormData just like Login.js
        const formData = new FormData();
        formData.append("email", email);
        formData.append("password", password);

        // Disable button while loading
        const goBtn = document.getElementById("rs-go-btn");
        goBtn.textContent = "Logging in…";
        goBtn.disabled = true;

        fetch("/Authorize", {
            method: "POST",
            body: formData,
            credentials: "include"
        })
            .then(res => res.json())
            .then(data => {
                console.log("Role Selector login response:", data);

                // Use role from response to navigate (mirrors Login.js)
                if (data.role == "registrar") {
                    location.href = "/registrardashboard";
                } else if (data.role == "student") {
                    location.href = "/dashboard";
                } else if (data.role == "admin") {
                    location.href = "/admindashboard";
                } else if (data.role == "dean_of_student") {
                    location.href = "/deandashboard";
                } else if (data.role == "dean_of_facult") {
                    location.href = "/deandashboard";
                } else if (data.role == "dean_of_science") {
                    location.href = "/deandashboard";
                } else if (data.role == "finance_office") {
                    location.href = "/financialdashboard";
                } else {
                    // Fallback: use the URL from ROLE_MAP
                    const entry = ROLE_MAP[roleIndex];
                    if (entry) {
                        location.href = entry.url;
                    }
                }
            })
            .catch(error => {
                console.error("Role Selector login error:", error);
                showToast("Login failed. Please try again.", "error");
                goBtn.textContent = "Go →";
                goBtn.disabled = false;
            });
    });

    // ─── Collapse / Expand ──────────────────────────────────
    const toggleBtn = document.getElementById("rs-toggle-btn");
    const body = document.getElementById("rs-body");
    let collapsed = localStorage.getItem("rs_collapsed") === "true";

    const applyCollapse = () => {
        body.style.display = collapsed ? "none" : "flex";
        toggleBtn.textContent = collapsed ? "+" : "−";
    };
    applyCollapse();

    toggleBtn.addEventListener("click", () => {
        collapsed = !collapsed;
        localStorage.setItem("rs_collapsed", collapsed);
        applyCollapse();
    });

    // ─── Drag Logic ─────────────────────────────────────────
    const handle = document.getElementById("rs-drag-handle");
    let isDragging = false;
    let offsetX = 0;
    let offsetY = 0;

    // Restore saved position
    const savedX = localStorage.getItem("rs_pos_x");
    const savedY = localStorage.getItem("rs_pos_y");
    if (savedX !== null && savedY !== null) {
        widget.style.left = savedX + "px";
        widget.style.top = savedY + "px";
        widget.style.right = "auto";
        widget.style.bottom = "auto";
    }

    handle.addEventListener("mousedown", (e) => {
        isDragging = true;
        const rect = widget.getBoundingClientRect();
        offsetX = e.clientX - rect.left;
        offsetY = e.clientY - rect.top;
        widget.style.transition = "none";
        widget.classList.add("rs-dragging");
        e.preventDefault();
    });

    document.addEventListener("mousemove", (e) => {
        if (!isDragging) return;
        let newX = e.clientX - offsetX;
        let newY = e.clientY - offsetY;

        // Constrain to viewport
        const maxX = window.innerWidth - widget.offsetWidth;
        const maxY = window.innerHeight - widget.offsetHeight;
        newX = Math.max(0, Math.min(newX, maxX));
        newY = Math.max(0, Math.min(newY, maxY));

        widget.style.left = newX + "px";
        widget.style.top = newY + "px";
        widget.style.right = "auto";
        widget.style.bottom = "auto";
    });

    document.addEventListener("mouseup", () => {
        if (!isDragging) return;
        isDragging = false;
        widget.classList.remove("rs-dragging");
        widget.style.transition = "";

        // Save position
        const rect = widget.getBoundingClientRect();
        localStorage.setItem("rs_pos_x", Math.round(rect.left));
        localStorage.setItem("rs_pos_y", Math.round(rect.top));
    });
});
