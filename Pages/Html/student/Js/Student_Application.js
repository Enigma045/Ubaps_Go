const steps = document.querySelectorAll(".step");
const stepIndicators = document.querySelectorAll(".progress-container li");
const progress = document.querySelector(".progress");


let currentStep = 0;

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
    alert("Please complete all required fields.");
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
    alert("Please complete all required fields.");
    return;
  }
  //between the lines
  console.log(Array.from(formData.entries()));
  let res = await fetch("/SubmitForm",{
    method:"POST",
    body:formData,
    credentials: "include" // send cookies
  })
 
  const text = await res.text();
  if(!res.ok){
    console.log("Server response:",text)
    return
  } 
  console.log("Server response:",text)
  alert("Application submitted successfully!");
  
  // Optional: Reset form and go back to first step
  // document.getElementById("bursaryForm").reset();
  // currentStep = 0;
  // showStep(currentStep);
};

// Initialize the form
showStep(currentStep);

