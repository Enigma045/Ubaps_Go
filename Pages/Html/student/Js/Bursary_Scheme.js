document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("scheme_info").onsubmit = async e => {
    e.preventDefault();

    const formData = new FormData(e.target);

    // ✅ Correctly log form values
    console.log(Array.from(formData.entries()));

    const res = await fetch("/benefactor", {
      method: "POST",
      body: formData,
      credentials: "include" // send cookies
    });

    if (res.ok) {
      showToast("Bursary scheme created successfully!", "success");
      setTimeout(() => {
        window.location.reload();
      }, 2000);
    } else {
      const text = await res.text();
      console.log(text);
      showToast(text || "Failed to create scheme. Please check your data.", "error");
      // Optionally don't reload on error to let user fix form
    }
  };

  const res = fetch("/getbenefactor", {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    }, // send cookies
  }).then((res) => {

    return res.json()
  }).then(data => {
    console.log(data)
    displayBenefactor(data)
  }).catch(err => {
    console.log(err);
    showToast("Failed to load bursary schemes.", "error");
  })
});

function displayBenefactor(data) {
  const tbody = document.getElementById("bursarytbody");
  tbody.innerHTML = "";
  data.forEach(item => {
    const tr = document.createElement("tr");
    tr.innerHTML = `
                <td>🟡</td>
                <td class="name">${item.scheme_name}</td>
                <td>${item.benefactor_name}</td>
                <td class="email">${item.benefactor_email}</td>
                <td>${item.total_fund_amount}</td>
                <td>${item.available_balance}</td>
                <td ><button title="Review"><img class="action" src="/Image/svgviewer-output (18).svg" alt=""></button>
                   <button title="delete" class="openModal2"><img class="action" src="/Image/svgviewer-output (21).svg" alt=""></button>
                </td>
            `;
    tbody.appendChild(tr);
  });
}
//<tbody id="bursarytbody"></tbody>
// <tr class="success">
//             <td>🟡</td>
//             <td>Loans Board</td>
//             <td>Richard</td>
//             <td>Lemillion045@gmail.com</td>
//             <td>2000000</td>
//             <td>200000</td>
//             <td><button title="Review"><img class="action" src="/Image/svgviewer-output (18).svg" alt=""></button>
//               <button title="delete"><img class="action" src="/Image/svgviewer-output (21).svg" alt=""></button>
//           </tr>