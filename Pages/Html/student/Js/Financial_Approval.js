document.addEventListener("DOMContentLoaded", () => {
  const tabs = document.querySelectorAll(".tabs input");
  const modal = document.getElementById("modal");
  const closeModalBtns = document.querySelectorAll("#closeModal, #closeModal2");
  const payForm = modal?.querySelector("form");
  const modalTitle = modal?.querySelector("h3");
  let activeStudent = null;

  // ── Close modal ──────────────────────────────────────────────────────────────
  closeModalBtns.forEach(btn => {
    btn.addEventListener("click", () => modal?.classList.remove("active"));
  });

  // ── Payment confirm ───────────────────────────────────────────────────────────
  if (payForm) {
    payForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      if (!activeStudent) return;

      try {
        const response = await fetch("/payinstallment", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ id: activeStudent.id, name: activeStudent.name }),
          credentials: "include"
        });

        if (!response.ok) throw new Error(await response.text() || "Payment failed");

        showToast(`Payment processed for ${activeStudent.name}`, "success");
        modal?.classList.remove("active");
        fetchFinancialRequests(getInitialStatuses());
      } catch (error) {
        console.error("Payment error:", error);
        showToast(error.message, "error");
      }
    });
  }

  // ── Tab change ────────────────────────────────────────────────────────────────
  tabs.forEach(tab => {
    tab.addEventListener("change", () => {
      fetchFinancialRequests(getInitialStatuses());
    });
  });

  const searchEl = document.getElementById('filter-search');
  if (searchEl) {
      searchEl.addEventListener('input', () => {
          fetchFinancialRequests(getInitialStatuses());
      });
  }

  const advancedFilter = document.getElementById('filter-advanced');
  if (advancedFilter) {
      window.onFilterChange = () => {
          fetchFinancialRequests(getInitialStatuses());
      };
  }

  const clearBtn = document.getElementById('clear-filters');
  if (clearBtn) {
      clearBtn.addEventListener('click', () => {
          if (searchEl) searchEl.value = '';
          if (advancedFilter) advancedFilter.value = '';
          window.activeAdvancedFilters = {};
          fetchFinancialRequests(getInitialStatuses());
      });
  }

  // ── Fetch ─────────────────────────────────────────────────────────────────────
  const fetchFinancialRequests = async (statuses = ["selected"]) => {
    try {
      if (statuses.length === 0) { renderTable([]); return; }

      const advanced = window.activeAdvancedFilters || {};
      const payload = {
        statuses: statuses,
        search: searchEl ? searchEl.value : "",
        department: advanced.dept || [],
        scheme: advanced.scheme || [],
        parent: advanced.parent || [],
        employment: advanced.employment || [],
        gender: advanced.gender || []
      };

      if (window.renderFilterTags) {
          window.renderFilterTags(payload, () => {
              fetchFinancialRequests(getInitialStatuses());
          });
      }

      const response = await fetch("/applicants?page=1&limit=1000", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
        credentials: "include"
      });

      if (!response.ok) throw new Error("Failed to fetch applicants");

      const result = await response.json();
      renderTable(result.data || []);
    } catch (error) {
      console.error("Error fetching applicants:", error);
      showToast("Failed to load students.", "error");
    }
  };

  // ── Render ────────────────────────────────────────────────────────────────────
  const renderTable = (applicants) => {
    const tbody = document.getElementById("financial-requests-body");
    if (!tbody) return;
    tbody.innerHTML = "";

    if (!applicants || applicants.length === 0) {
      tbody.innerHTML = '<tr><td colspan="7" style="text-align:center;padding:24px;">No students found</td></tr>';
      return;
    }

    const infoIcon = `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M560-680v-80h320v80H560Zm0 160v-80h320v80H560Zm0 160v-80h320v80H560Zm-240-40q-50 0-85-35t-35-85q0-50 35-85t85-35q50 0 85 35t35 85q0 50-35 85t-85 35ZM80-160v-76q0-21 10-40t28-30q45-27 95.5-40.5T320-360q56 0 106.5 13.5T522-306q18 11 28 30t10 40v76H80Zm86-80h308q-35-20-74-30t-80-10q-41 0-80 10t-74 30Zm154-240q17 0 28.5-11.5T360-520q0-17-11.5-28.5T320-560q-17 0-28.5 11.5T280-520q0 17 11.5 28.5T320-480Zm0-40Zm0 280Z"/></svg>`;
    const payIcon = `<svg xmlns="http://www.w3.org/2000/svg" height="20" viewBox="0 -960 960 960" width="20" fill="currentColor"><path d="M840-200v-80H200v80h640Zm0-160v-240H200v240h640Zm0-320v-80H200v80h640ZM200-80q-33 0-56.5-23.5T120-160v-640q0-33 23.5-56.5T200-880h640q33 0 56.5 23.5T920-800v640q0 33-23.5 56.5T840-80H200ZM200-800h640v640H200v-640Z"/></svg>`;

    applicants.forEach(app => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${app.first_name} ${app.last_name}</td>
        <td>${app.registration_number || "N/A"}</td>
        <td>${app.parent_guardian_status || "N/A"}</td>
        <td><span class="status status-paid">${app.guardian_employment_status || "N/A"}</span></td>
        <td>${app.scheme_name || "N/A"}</td>
        <td><span class="status status-paid">${app.bursary_amount ? parseFloat(app.bursary_amount).toLocaleString() : "0"}</span></td>
        <td>
          <div style="display:flex;gap:8px;">
            <button class="action-btn info-btn" title="${app.reason ? app.reason.String : "No details"}">${infoIcon}</button>
            <button class="action-btn pay-btn" title="Pay Installment">${payIcon}</button>
          </div>
        </td>
      `;

      // Attach click directly — avoids SVG event-target issues
      row.querySelector(".pay-btn").addEventListener("click", () => {
        activeStudent = {
          id: app.registration_number,
          name: `${app.first_name} ${app.last_name}`
        };
        if (modalTitle) modalTitle.textContent = `Pay ${activeStudent.name} Fees?`;
        modal?.classList.add("active");
      });

      tbody.appendChild(row);
    });

    // Trigger filter after dynamic load
    if (window.triggerTableFilter) window.triggerTableFilter();
  };

  // ── Initial status helper ─────────────────────────────────────────────────────
  const getInitialStatuses = () => {
    const statuses = Array.from(tabs)
      .filter(t => t.checked)
      .map(t => t.value);

    return statuses;
  };

  // ── Boot ──────────────────────────────────────────────────────────────────────
  fetchFinancialRequests(getInitialStatuses());
});