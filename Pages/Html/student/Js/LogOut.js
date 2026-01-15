document.querySelector('.logout')?.addEventListener('click', () => {
  //alert('You have been logged out.');
  fetch("/logout").then((data)=>{
    if (!data.ok){
      alert('Failed to kill session logged out.');
    }
    return data.text()
  }).then((res)=>{
      alert("Logging Out");
      console.log("Server response:", res);
  }).catch((err)=>{
    alert(err);
    return
  })

  window.location.href = "/Login";
});
