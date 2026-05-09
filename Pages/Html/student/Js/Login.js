document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("loginForm").onsubmit = async e => {
    e.preventDefault();

    const formData = new FormData(e.target);

    // ✅ Correctly log form values
    console.log(Array.from(formData.entries()));

    const payload = Object.fromEntries(formData.entries());

    fetch("/Authorize", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(payload),
      credentials: "include"
    }).then(async res => {
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.message || "Invalid credentials");
      }
      return data;
    }).then(data => {
      showToast("Login successful! Redirecting...", "success");

      // Delay redirection slightly so they see the toast
      setTimeout(() => {
        if (data.role == "registrar") {
          location.href = "/registrardashboard";
        } else if (data.role == "student") {
          location.href = "/dashboard";
        } else if (data.role == "admin") {
          location.href = "/admindashboard";
        } else if (data.role == "dean_of_student" || data.role == "dean_of_facult" || data.role == "dean_of_science") {
          location.href = "/deandashboard";
        } else if (data.role == "finance_office") {
          location.href = "/financialdashboard";
        }
      }, 1000);

      console.log(data);
    }).catch(error => {
      console.error(error);
      showToast(error.message || "Login failed. Please try again.", "error");
    });

    // if (res.ok){

    //     if(res.text == "registrar"){
    //     location.href = "/admindashboard";
    //     console.log("here")
    //     }else if(res.text == "student"){
    //     location.href = "/dashboard";
    //     console.log("here")
    //     }else if(res.text == "admin"){
    //     location.href = "/admindashboard";
    //      console.log("here")
    //     }else if(res.text == "dean_of_student"){

    //     }else if(res.text == "finance_office"){
    //      location.href = "/financial";
    //      console.log("here")  
    //     }

    //     }else{
    //         const text = await res.text();
    //         console.log(text);
    //     } // change to your page
  };
});
