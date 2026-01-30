// const tabs = document.querySelectorAll(".tabs button");
// const tables = document.querySelectorAll(".log-table");

// tabs.forEach(tab => {
//   tab.addEventListener("click", () => {

//     // Remove active state from tabs
//     tabs.forEach(t => t.classList.remove("active"));
//     tab.classList.add("active");

//     // Hide all tables
//     tables.forEach(table => table.classList.remove("active"));

//     // Show selected table
//     const target = tab.dataset.tab;
//     document.getElementById(target).classList.add("active");
//   });
// });
document.addEventListener("DOMContentLoaded", function() {
///getuserdetails

    fetch('/getuserdetails', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
      console.log('Success:', data);

      data.forEach(data =>{
      const tbody =  document.querySelector(".tbody")
      const tr = document.createElement("tr")
      tr.className = "success"
      const name = document.createElement("td")
      const role = document.createElement("td")
      const email = document.createElement("td")
      const phone = document.createElement("td")
      const actions = document.createElement("td")
      
      const edit = document.createElement("button")

      edit.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#000"><path d="M200-200h57l391-391-57-57-391 391v57Zm-80 80v-170l528-527q12-11 26.5-17t30.5-6q16 0 31 6t26 18l55 56q12 11 17.5 26t5.5 30q0 16-5.5 30.5T817-647L290-120H120Zm640-584-56-56 56 56Zm-141 85-28-29 57 57-29-28Z"/></svg>`
      const del = document.createElement("button")
      del.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#000"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg>`
      actions.appendChild(edit)
      actions.appendChild(del)

      name.textContent = data.first +" "+ data.last
      role.textContent = data.role
      email.textContent = data.email
      phone.textContent = data.phone


      tr.appendChild(name)
      tr.appendChild(role)
      tr.appendChild(email)
      tr.appendChild(phone)
      tr.appendChild(actions)
      
      tbody.appendChild(tr);
      })
    })
    .catch((error) => {
      console.error('Error:', error);
    });
});

// <tr class="success">
//           <td>Richard sambo</td>
//           <td>Student</td>
//           <td>CEN-01-14-22@unilia.ac.mw</td>
//           <td>characte</td>
//           <td>
//             <Button><svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#000"><path d="M200-200h57l391-391-57-57-391 391v57Zm-80 80v-170l528-527q12-11 26.5-17t30.5-6q16 0 31 6t26 18l55 56q12 11 17.5 26t5.5 30q0 16-5.5 30.5T817-647L290-120H120Zm640-584-56-56 56 56Zm-141 85-28-29 57 57-29-28Z"/></svg></Button>
//             <Button><svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#000"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg></Button>          
//           </td>
//         </tr>