document.getElementById("UserForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  // REGEX
  const nameRegex = /^[A-Za-z]{2,}$/;
  const emailRegex = /^[a-zA-Z0-9._%+-]+@unilia\.ac\.mw$/;
  const phoneRegex = /^09\d{8}$/;
  const passwordRegex = /^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&]).{8,}$/;

  const firstName = document.getElementById("firstName").value.trim();
  const lastName = document.getElementById("lastName").value.trim();
  const email = document.getElementById("email").value.trim();
  const phone = document.getElementById("phone").value.trim();
  const password = document.getElementById("password").value;
  const confirmPassword = document.getElementById("confirmPassword").value;
  const role = document.querySelector('input[name="Users"]:checked')?.value;

  // VALIDATION
  if (!nameRegex.test(firstName) || !nameRegex.test(lastName)) {
    alert("Names must contain letters only (min 2 characters)");
    return;
  }

  if (!emailRegex.test(email)) {
    alert("Email must be a valid @unilia.ac address");
    return;
  }

  if (!phoneRegex.test(phone)) {
    alert("Phone must be in format 09XXXXXXXX");
    return;
  }

  if (!passwordRegex.test(password)) {
    alert("Password must be 8+ chars with upper, lower, number & symbol");
    return;
  }

  if (password !== confirmPassword) {
    alert("Passwords do not match");
    return;
  }

  if (!role) {
    alert("Select a user role");
    return;
  }

  // SEND TO BACKEND
  const payload = {
    first_name: firstName,
    last_name: lastName,
    email: email,
    phone: phone,
    password: password,
    role: role
  };

  console.log("Payload:", payload); // Debugging line

  const res = await fetch("/createuser", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload)
  });

  const result = await res.json();
  alert(result.message);
});
