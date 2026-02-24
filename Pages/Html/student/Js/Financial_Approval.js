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

    // ✅ Correctly log form values
    console.log(Array.from(formData.entries()));

    const res = await fetch("/sendrequest", {
      method: "POST",
      body: formData,
      credentials: "include" // send cookies
    });

    if (res.ok) {
      showToast("Financial request processed successfully!", "success");
      setTimeout(() => window.location.reload(), 2000);
    } else {
      const text = await res.text();
      console.log(text);
      showToast(text || "Failed to process request.", "error");
    }
  };
});