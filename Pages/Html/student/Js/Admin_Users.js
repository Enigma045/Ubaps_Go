document.addEventListener("DOMContentLoaded", () => {
  const userForm = document.getElementById("UserForm");
  if (!userForm) return;

  const firstNameInput = document.getElementById("firstName");
  const lastNameInput = document.getElementById("lastName");
  const emailInput = document.getElementById("email");
  const phoneInput = document.getElementById("phone");
  const passwordInput = document.getElementById("password");
  const confirmPasswordInput = document.getElementById("confirmPassword");
  const userTypeElem = document.getElementById("userType");

  // REGEX (Standardized with Register.js)
  const nameRegex = /^[A-Za-z]{2,}(?:\s[A-Za-z]{2,})*$/;
  const emailRegex = /^[a-zA-Z0-9._%+-]+@unilia\.ac\.mw$/;
  const phoneRegex = /^(09|08)\d{8}$/;
  const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/;

  // Utility: show error
  function showError(input, message) {
    clearError(input);
    const error = document.createElement("div");
    error.className = "error-msg";
    error.innerText = message;
    input.closest(".form-group").appendChild(error);
    input.classList.add("error");
  }

  function clearError(input) {
    input.classList.remove("error");
    const parent = input.closest(".form-group");
    if (!parent) return;
    const err = parent.querySelector(".error-msg");
    if (err) err.remove();
  }

  function clearAllErrors() {
    document.querySelectorAll(".error-msg").forEach(e => e.remove());
    document.querySelectorAll(".error").forEach(e => e.classList.remove("error"));
  }

  userForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    clearAllErrors();

    const firstName = firstNameInput.value.trim();
    const lastName = lastNameInput.value.trim();
    const email = emailInput.value.trim();
    const phone = phoneInput.value.trim();
    const password = passwordInput.value;
    const confirmPassword = confirmPasswordInput.value;
    const role = userTypeElem ? userTypeElem.value : "";

    let valid = true;

    // VALIDATION
    if (!nameRegex.test(firstName)) {
      showError(firstNameInput, "Enter a valid first name");
      valid = false;
    }

    if (!nameRegex.test(lastName)) {
      showError(lastNameInput, "Enter a valid surname");
      valid = false;
    }

    if (!emailRegex.test(email)) {
      showError(emailInput, "Email must be a valid @unilia.ac.mw address");
      valid = false;
    }

    if (!phoneRegex.test(phone)) {
      showError(phoneInput, "Phone must be 10 digits and start with 09 or 08");
      valid = false;
    }

    if (!passwordRegex.test(password)) {
      showError(passwordInput, "Password must be 8+ chars, include upper, lower, number & symbol");
      valid = false;
    }

    if (password !== confirmPassword) {
      showError(confirmPasswordInput, "Passwords do not match");
      valid = false;
    }

    if (!role) {
      showToast("Please select a valid user role", "error");
      valid = false;
    }

    if (!valid) return;

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
