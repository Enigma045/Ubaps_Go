const tabs = document.querySelectorAll(".tabs button");
const tables = document.querySelectorAll(".log-table");

tabs.forEach(tab => {
  tab.addEventListener("click", () => {

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
