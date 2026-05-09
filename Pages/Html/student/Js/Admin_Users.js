document.addEventListener("DOMContentLoaded", () => {
  const userForm = document.getElementById("UserForm");
  if (!userForm) return;

  userForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    // REGEX
    const nameRegex = /^[A-Za-z]{2,}$/;
    const emailRegex = /^[a-zA-Z0-9._%+-]+@unilia\.ac\.mw$/;
    const phoneRegex = /^(09|08)\d{8}$/;
    const passwordRegex = /^(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*?&]).{8,}$/;

    const firstName = document.getElementById("firstName").value.trim();
    const lastName = document.getElementById("lastName").value.trim();
    const email = document.getElementById("email").value.trim();
    const phone = document.getElementById("phone").value.trim();
    const password = document.getElementById("password").value;
    const confirmPassword = document.getElementById("confirmPassword").value;
    const userTypeElem = document.getElementById("userType");
    const role = userTypeElem ? userTypeElem.value : "";

    // VALIDATION
    if (!firstName || !lastName || !email || !phone || !password) {
        showToast("Please fill in all required fields", "error");
        return;
    }

    if (!nameRegex.test(firstName) || !nameRegex.test(lastName)) {
      showToast("Names must contain letters only (min 2 characters)", "error");
      return;
    }

    if (!emailRegex.test(email)) {
      showToast("Email must be a valid @unilia.ac.mw address", "error");
      return;
    }

    if (!phoneRegex.test(phone)) {
      showToast("Phone must be 10 digits and start with 08 or 09", "error");
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
      showToast("Please select a valid user role", "error");
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

    const submitBtn = document.getElementById("submit");
    if (submitBtn) {
        submitBtn.disabled = true;
        submitBtn.textContent = "Creating...";
    }

    try {
        const res = await fetch("/createuser", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        });

        const result = await res.json();
        if (res.ok) {
            showToast(result.message || "User created successfully!", "success");
            userForm.reset();
        } else {
            showToast(result.message || "Failed to create user.", "error");
        }
    } catch (err) {
        console.error("Fetch error:", err);
        showToast("Connection error. Please check your network.", "error");
    } finally {
        if (submitBtn) {
            submitBtn.disabled = false;
            submitBtn.textContent = "Create User";
        }
    }
  });
});
