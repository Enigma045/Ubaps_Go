document.addEventListener("DOMContentLoaded", function() {
  let count = 3;
  const badge = document.getElementById("badge-nav");
  const parent = document.getElementById("parentDiv");

  fetch('/countnotifications',{
    headers: {
      'Content-Type': 'application/json'
    }
  }).then(response => response.json())
  .then(data => {
    count = data;
    console.log(count)
  })
  .catch(error => {
    console.error("Error fetching notifications:", error);
  });

  if (count > 0) {
    badge.textContent = count;
    badge.style.display = "inline-block";      
  }
  });