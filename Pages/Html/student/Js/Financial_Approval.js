document.addEventListener("DOMContentLoaded", () => {
  const tabs = document.querySelectorAll(".tabs input");
  const tables = document.querySelectorAll(".log-table");

  tabs.forEach(tab => {
    console.log(tab)
    tab.addEventListener("change", () => {

      // Remove active state from tabs
      tabs.forEach(t => t.classList.remove("active"));
      tab.classList.add("active");

      // Hide all tables
      tables.forEach(table => table.classList.remove("active"));

      // Show selected table
      const target = tab.dataset.tab;
      document.getElementById(target).classList.add("active");
    });
  });


  document.getElementById("RequestForm").onsubmit = async e => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const res = await fetch("/sendrequest", {
      method: "POST",
      body: formData,
      credentials: "include"
    });

    if (res.ok) {
      showToast("Financial request processed successfully!", "success");
      setTimeout(() => window.location.reload(), 2000);
    } else {
      const text = await res.text();
      showToast(text || "Failed to process request.", "error");
    }
  };

  // Fetch and render data
  const fetchFinancialRequests = async () => {
    try {
      const response = await fetch("/applicants", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(["selected"]), // Populating with students who have been selected
        credentials: "include"
      });

      if (!response.ok) throw new Error("Failed to fetch applicants");

      const result = await response.json();
      // Handle paginated response structure { data: [], total: x, ... }
      renderTable(result.data || []);
    } catch (error) {
      console.error("Error fetching selected applicants:", error);
      showToast("Failed to load selected students.", "error");
    }
  };

  const fetchTotalAmount = async () => {
    try {
      const response = await fetch("/gettotalamount", {
        method: "GET",
        credentials: "include"
      });

      if (!response.ok) throw new Error("Failed to fetch total amount");

      const amounts = await response.json();
      const total = amounts ? amounts.reduce((sum, amt) => sum + amt, 0) : 0;
      const totalEl = document.querySelector(".tabs h4");
      if (totalEl) totalEl.textContent = `MWK ${total.toLocaleString()}`;
    } catch (error) {
      console.error("Error fetching total amount:", error);
    }
  };

  const renderTable = (applicants) => {
    const tbody = document.getElementById("financial-requests-body");
    if (!tbody) return;
    tbody.innerHTML = "";

    if (!applicants || applicants.length === 0) {
      tbody.innerHTML = '<tr><td colspan="8" style="text-align: center;">No selected students found</td></tr>';
      return;
    }

    applicants.forEach(app => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${app.first_name} ${app.last_name}</td>
        <td>${app.registration_number || "N/A"}</td>
        <td>${app.parent_guardian_status || "N/A"}</td>
        <td><span class="status status-paid">${app.guardian_employment_status || "N/A"}</span></td>
        <td>${app.programme || "N/A"}</td>
        <td><span class="status status-paid">${app.bursary_amount ? parseFloat(app.bursary_amount).toLocaleString() : "0"}</span></td>
        <td><button title="${app.reason ? app.reason.String : "No details provided"}">More Info</button></td>
        <td><button class="pay-btn" data-id="${app.registration_number}" data-name="${app.first_name} ${app.last_name}">Pay Installment</button></td>
      `;
      tbody.appendChild(row);
    });
  };

  // Initial load
  fetchFinancialRequests();
  fetchTotalAmount();
});