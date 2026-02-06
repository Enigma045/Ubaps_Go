const tabs = document.querySelectorAll(".tabs button");
const tables = document.querySelectorAll(".log-table");
let app = false;
let pay = false;
let user = false;
//


tabs.forEach(tab => {
  tab.addEventListener("click", () => {

    // Remove active state from tabs
    tabs.forEach(t => t.classList.remove("active"));
    tab.classList.add("active");

    // Hide all tables
    tables.forEach(table => table.classList.remove("active"));

    // Show selected table
    const target = tab.dataset.tab;
    console.log(target)
    document.getElementById(target).classList.add("active");

    if (target == "applications"){
    

    console.log(target)
    }else if (target == "payments"){
    

    console.log(target)
    }else if (target == "users"){
    

    console.log(target)
    }
  });
});


// users
    fetch("/applicationlog",{
      method:"GET",
      headers:{
        'Content-Type':'application/json'
      }
    }).then(response => {
    if (response.ok) {
      //location.reload(); // Reload the page to reflect changes
      response.json().then(data => {
        console.log(data);
        if (app == false){
       data.forEach(log => {
        const tableBody = document.getElementById("application-body")
        const row = document.createElement("tr");
        row.innerHTML = `
        <tr class="success">
          <td>${log.occurred_at}</td>
          <td>${log.user_role}</td>
          <td>${log.user_id}</td>
          <td>${log.action}</td>
          <td>${log.application}</td>
          
          <td>${log.status}</td>
          <td>${log.duration_ms}ms</td>
        </tr>
        `;
        //<td>${log.amount}</td>
        tableBody.appendChild(row);
        //      
      })
      app = true;
      }
      });
    } else {
      console.error("Failed to retrieve application logs");
    }
  }).catch(error => {
    console.error("Error:", error);
  });
  //
  // users
    fetch("/paymentlog",{
      method:"GET",
      headers:{
        'Content-Type':'application/json'
      }
    }).then(response => {
    if (response.ok) {
      //location.reload(); // Reload the page to reflect changes
      response.json().then(data => {
        console.log(data);
        if (pay == false){
       data.forEach(log => {
        const tableBody = document.getElementById("payment-body");
        const row = document.createElement("tr");
        row.innerHTML = `
        <tr class="success">
          <td>${log.occurred_at}</td>
          <td>${log.user_role}</td>
          <td>${log.user_id}</td>
          <td>${log.action}</td>
          <td>${log.target}</td>
          <td>${log.amount}</td>
          <td>${log.status}</td>
          <td>${log.duration_ms}ms</td>
        </tr>
        `;
        
        tableBody.appendChild(row);
        //      
      })
      pay = true;
      }
      });
    } else {
      console.error("Failed to retrieve payment logs");
    }
  }).catch(error => {
    console.error("Error:", error);
  });
  //
  // users
    fetch("/userlog",{
      method:"GET",
      headers:{
        'Content-Type':'application/json'
      },
      
    }).then(response => {
    if (response.ok) {
      //location.reload(); // Reload the page to reflect changes
      response.json().then(data => {
        console.log(data);
        if (user == false){
       data.forEach(log => {
        const tableBody = document.getElementById("user-body");
        const row = document.createElement("tr");
        row.innerHTML = `
        <tr class="success">
          <td>${log.occurred_at}</td>
          <td>${log.user_role}</td>
          <td>${log.user_id}</td>
          <td>${log.action}</td>
          <td>${log.target}</td>
          <td>${log.status}</td>
          <td>${log.duration_ms}ms</td>
        </tr>
        `;
        
        tableBody.appendChild(row);
        //      
      })
      user = true;
      }
      });
    } else {
      console.error("Failed to retrieve user logs");
    }
  }).catch(error => {
    console.error("Error:", error);
  });
  //