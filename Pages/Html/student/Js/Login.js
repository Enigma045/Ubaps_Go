document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("loginForm").onsubmit = async e => {
    e.preventDefault();

    const formData = new FormData(e.target);

    // ✅ Correctly log form values
    console.log(Array.from(formData.entries()));

    fetch("/Authorize", {
      method: "POST",
      body: formData,
      credentials: "include" // send cookies
    }).then(res => {
      return res.json()
    }).then(data => {
      if(data.role == "registrar"){
        location.href = "/decision";
        console.log("here")
        }else if(data.role == "student"){
        location.href = "/dashboard";
        console.log("here")
        }else if(data.role == "admin"){
        location.href = "/admindashboard";
         console.log("here")
        }else if(data.role == "dean_of_student"){

        }else if(data.role == "finance_office"){
         location.href = "/financial";
         console.log("here")  
        }
      console.log(data)
    }).catch(error => {
      console.error(error)
    })

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
