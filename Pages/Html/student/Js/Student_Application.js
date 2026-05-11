const steps = document.querySelectorAll(".step");
const stepIndicators = document.querySelectorAll(".progress-container li");
const progress = document.querySelector(".progress");


let currentStep = 0;

// Restrict date of birth to past dates
const dobInput = document.getElementById("dob");
if (dobInput) {
  const today = new Date().toISOString().split('T')[0];
  dobInput.setAttribute('max', today);
}

function showStep(index) {
  // Remove active class from all steps
  steps.forEach(step => step.classList.remove("active"));
  steps[index].classList.add("active");

  // Update step indicators
  stepIndicators.forEach((li, i) => {
    li.classList.remove("done", "current");
    if (i < index) li.classList.add("done");
    if (i === index) li.classList.add("current");
  });

  // Update progress bar
  const percent = index / (steps.length - 1);
  progress.style.transform = `translateY(-50%) scaleX(${percent})`;

  // Update button visibility
  document.getElementById("prev").style.display = index === 0 ? "none" : "inline-block";
  document.getElementById("next").style.display = index === steps.length - 1 ? "none" : "inline-block";
  document.getElementById("submit").style.display = index === steps.length - 1 ? "inline-block" : "none";
}

function validateStep() {
  const inputs = steps[currentStep].querySelectorAll("input, select");

  // If no inputs (like review step), return true
  if (inputs.length === 0) return true;

  let isValid = true;

  inputs.forEach(input => {
    if (!input.checkValidity()) {
      isValid = false;
      // Add visual feedback - red border for invalid fields
      input.style.borderColor = "#ef4444";
      input.style.boxShadow = "0 0 0 3px rgba(239, 68, 68, 0.1)";
    } else {
      // Reset to default border
      input.style.borderColor = "#9ca3af";
      input.style.boxShadow = "none";
    }
  });

  return isValid;
}

// Clear validation styling when user starts typing
document.addEventListener('input', (e) => {
  if (e.target.matches('input, select')) {
    if (e.target.checkValidity()) {
      e.target.style.borderColor = "#9ca3af";
      e.target.style.boxShadow = "none";
    }
  }
});

document.getElementById("next").onclick = () => {
  if (!validateStep()) {
    showToast("Please complete all required fields.", "error");
    return;
  }

  // Prevent going beyond last step
  if (currentStep < steps.length - 1) {
    currentStep++;
    showStep(currentStep);
  }
};

document.getElementById("prev").onclick = () => {
  // Prevent going before first step
  if (currentStep > 0) {
    currentStep--;
    showStep(currentStep);
  }
};

document.getElementById("bursaryForm").onsubmit = async e => {
  e.preventDefault();
  const formData = new FormData(e.target)
  // Final validation check
  if (!validateStep()) {
    showToast("Please complete all required fields.", "error");
    return;
  }
  //between the lines
  console.log(Array.from(formData.entries()));
  let res = await fetch("/SubmitForm", {
    method: "POST",
    body: formData,
    credentials: "include" // send cookies
  })

  const text = await res.text();
  if (!res.ok) {
    console.log("Server response:", text);
    showToast(text || "Failed to submit application. Please check your data.", "error");
    return;
  }
  console.log("Server response:", text);
  showToast("Application submitted successfully!", "success");

  // Optional: Redirect or reset
  // setTimeout(() => window.location.href = "/dashboard", 2000);
};

// Initialize the form
async function initForm() {
  try {
    const res = await fetch("/api/application/my");
    if (!res.ok) return;
    const data = await res.json();

    if (data.status && data.status !== 'not submitted' && data.status !== 'submitted') {
      // Application is closed for updates
      const form = document.getElementById("bursaryForm");
      const inputs = form.querySelectorAll("input, select, textarea, button");
      inputs.forEach(i => {
        if (i.id !== 'prev' && i.id !== 'next') i.disabled = true;
      });
      
      const submitBtn = document.getElementById("submit");
      if (submitBtn) {
        submitBtn.style.display = 'none';
        const msg = document.createElement("p");
        msg.style.color = "#ef4444";
        msg.style.fontWeight = "bold";
        msg.style.textAlign = "center";
        msg.style.marginTop = "20px";
        msg.innerText = "This application is currently being processed and cannot be updated.";
        submitBtn.parentNode.appendChild(msg);
      }
    }

    // Pre-fill fields
    if (data.dob) {
      const dob = new Date(data.dob).toISOString().split('T')[0];
      const dobEl = document.getElementById("dob");
      if (dobEl) dobEl.value = dob;
    }
    
    if (data.gender) {
      const genderEl = document.querySelector(`input[name="gender"][value="${data.gender}"]`);
      if (genderEl) genderEl.checked = true;
    }

    if (data.home_district) {
      const distEl = document.getElementById("HomeDistrict");
      if (distEl) distEl.value = data.home_district;
    }

    if (data.accommodation) {
      const accEl = document.querySelector(`input[name="Accomodation"][value="${data.accommodation}"]`);
      if (accEl) accEl.checked = true;
    }

    if (data.guardian_status) {
      const gsEl = document.querySelector(`select[name="Gurdian Status"]`);
      if (gsEl) gsEl.value = data.guardian_status;
    }

    if (data.guardian_employment_status) {
      const gesEl = document.querySelector(`select[name="Guardian Employment Status"]`);
      if (gesEl) gesEl.value = data.guardian_employment_status;
    }

    if (data.other_support) {
      const osEl = document.getElementById("otherSupport");
      if (osEl) osEl.value = data.other_support;
    }

    if (data.fee_responsibility) {
      const frEl = document.getElementById("feeResponsibility");
      if (frEl) frEl.value = data.fee_responsibility;
    }

    if (data.financial_hardship) {
      const fhEl = document.getElementById("financialHardship");
      if (fhEl) fhEl.value = data.financial_hardship;
    }

    if (data.impact_of_no_bursary) {
      const inbEl = document.getElementById("impactOfNoBursary");
      if (inbEl) inbEl.value = data.impact_of_no_bursary;
    }

    if (data.reason) {
      const reasonEl = document.getElementById("Reason");
      if (reasonEl) reasonEl.value = data.reason;
    }

  } catch (err) {
    console.error("Error fetching application:", err);
  }
}

initForm();
showStep(currentStep);
