document.addEventListener("DOMContentLoaded", function () {

  fetch('/getuserdetails', {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  })
    .then(response => response.json())
    .then(data => {
      const tbody = document.querySelector(".tbody");

      data.forEach(user => {
        const tr = document.createElement("tr");

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

        const del = document.createElement("button");
        del.className = "openModal";
        del.innerHTML = `<svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px" fill="#000"><path d="M280-120q-33 0-56.5-23.5T200-200v-520h-40v-80h200v-40h240v40h200v80h-40v520q0 33-23.5 56.5T680-120H280Zm400-600H280v520h400v-520ZM360-280h80v-360h-80v360Zm160 0h80v-360h-80v360ZM280-720v520-520Z"/></svg>`;

        actions.appendChild(edit);
        actions.appendChild(del);

        tr.appendChild(name);
        tr.appendChild(role);
        tr.appendChild(email);
        tr.appendChild(phone);
        tr.appendChild(actions);

        tbody.appendChild(tr);
      });
    })
    .catch(error => console.error(error));

});


