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
    showToast("Names must contain letters only (min 2 characters)", "error");
    return;
  }

  if (!emailRegex.test(email)) {
    showToast("Email must be a valid @unilia.ac address", "error");
    return;
  }

  if (!phoneRegex.test(phone)) {
    showToast("Phone must be in format 09XXXXXXXX", "error");
    return;
  }

  if (!passwordRegex.test(password)) {
    showToast("Password must be 8+ chars with upper, lower, number & symbol", "error");
    return;
  }

  if (password !== confirmPassword) {
    showToast("Passwords do not match", "error");
    return;
  }

  if (!role) {
    showToast("Select a user role", "error");
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
  if (res.ok) {
    showToast(result.message || "User created successfully!", "success");
    e.target.reset();
  } else {
    showToast(result.message || "Failed to create user.", "error");
  }
});
