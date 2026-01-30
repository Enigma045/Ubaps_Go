document.addEventListener("DOMContentLoaded", function() {
  const badge = document.getElementById("badge-nav");
  const parent = document.getElementById("parentDiv");

  fetch('/countnotifications',{
    headers: {
      'Content-Type': 'application/json'
    }
  }).then(response => response.json())
  .then(data => {
    let count = data;
    if (count > 0) {
    badge.textContent = count;
    badge.style.display = "inline-block";      
    }
    console.log(count)
  })
  .catch(error => {
    console.error("Error fetching notifications:", error);
  });

  
  });