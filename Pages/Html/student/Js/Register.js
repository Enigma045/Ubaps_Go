document.addEventListener("DOMContentLoaded", () => {
  const form = document.querySelector("form");

  // Inputs
  const nameInput = document.querySelector('input[placeholder="Name"]');
  const surnameInput = document.querySelector('input[placeholder="SurName"]');
  const phoneInput = document.querySelector('input[placeholder="Phone"]');
  const passwordInput = document.querySelector('input[placeholder="Password"]');
  const confirmPasswordInput = document.querySelector('input[placeholder="Confirm Password"]');

  const programSelect = document.getElementById("pro");
  const numberInputs = document.querySelectorAll(".option-group input[type='number']");

  // Regex rules
  const nameRegex = /^[A-Za-z]{2,}(?:\s[A-Za-z]{2,})*$/;
  const phoneRegex = /^\d{7,15}$/;
  const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?#&])[A-Za-z\d@$!%*?#&]{8,}$/;
  const twoDigitRegex = /^\d{2}$/;
  const regNumberRegex = /^(CEN|BPH|EDH|FSN)-\d{2}-\d{2}-\d{2}$/;

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

  form.addEventListener("submit", (e) => {
    e.preventDefault();
    clearAllErrors();

    let valid = true;

    // Name
    if (!nameRegex.test(nameInput.value.trim())) {
      showError(nameInput, "Enter a valid name");
      valid = false;
    }

    // Surname
    if (!nameRegex.test(surnameInput.value.trim())) {
      showError(surnameInput, "Enter a valid surname");
      valid = false;
    }

    // Phone
    if (!phoneRegex.test(phoneInput.value.trim())) {
      showError(phoneInput, "Enter a valid phone number");
      valid = false;
    }

    // Password
    if (!passwordRegex.test(passwordInput.value)) {
      showError(
        passwordInput,
        "Password must be 8+ chars, include upper, lower, number & symbol"
      );
      valid = false;
    }

    // Confirm Password
    if (confirmPasswordInput.value !== passwordInput.value) {
      showError(confirmPasswordInput, "Passwords do not match");
      valid = false;
    }

    // Reg Number parts
    const program = programSelect.value;
    const nums = [];

    numberInputs.forEach((input) => {
      if (!twoDigitRegex.test(input.value)) {
        showError(input, "Must be exactly 2 digits");
        valid = false;
      }
      nums.push(input.value.padStart(2, "0"));
    });

    // Build Reg Number
    const regNumber = `${program}-${nums.join("-")}`;

    if (!regNumberRegex.test(regNumber)) {
      alert("Invalid Registration Number format");
      valid = false;
    }

    if (!valid) return;

    // Final payload
    const payload = {
      name:nameInput.value.trim(),
      surname: surnameInput.value.trim(),
      phone: phoneInput.value.trim(),
      password: passwordInput.value,
      reg_number: regNumber
    };

    console.log("✅ Registration payload:", payload);

    // TODO: send to backend
    fetch("/register",{
      method: "POST",
      headers:{
        "Content-Type":"application"
      },
      body: JSON.stringify(payload)

    }).then(res => {
      if(!res.ok){
         return res.json().then(err => {
          throw new Error(err.message || "Registration failed");
         })
      }
      return res.json();
    })
    .then(data => {
      alert("Registration sent Successfully");
      console.log("Server response:", data);

      window.location.href = "";
    }).catch(err =>{
      console.error("Register error:",err.message);
      alert(err.message);
    }
    );
    // fetch("/register", { ... })

    alert("Registration successful!\nReg Number: " + regNumber);
    form.reset();
  });
});
