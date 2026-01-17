document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("scheme_info").onsubmit = async e => {
    e.preventDefault();

    const formData = new FormData(e.target);

    // ✅ Correctly log form values
    console.log(Array.from(formData.entries()));

    const res = await fetch("/benefactor", {
      method: "POST",
      body: formData,
      credentials: "include" // send cookies
    });

    if (res.ok){
        //  location.href = "/dashboard";
         console.log("here")
        }else{
            const text = await res.text();
            console.log(text);
        } // change to your page
  };
});