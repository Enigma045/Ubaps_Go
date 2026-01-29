  document.addEventListener("DOMContentLoaded", function() {
  const parent = document.getElementById("parentDiv");

  fetch("/getnotifications",{
    method: "GET",
    headers: {
        "Content-Type": "application/json"
    }
  })
  .then(response => response.json())
  .then(data => {
    console.log(data);
    data.forEach(notification => {
      const notificationDiv = document.createElement("div");
      notificationDiv.classList.add("notification");
      notificationDiv.classList.add("success");
      const badge = document.createElement("span");
      badge.classList.add("badge");
      const time = document.createElement("span");
      time.classList.add("time");
      time.textContent = notification.timestamp;
      const contentDiv = document.createElement("div");
      contentDiv.classList.add("content");
      const title = document.createElement("h4");
      const message = document.createElement("p");
      
      title.textContent = notification.title;
      message.textContent = notification.message;
      
      contentDiv.appendChild(title);
      contentDiv.appendChild(message);
      
      
      notificationDiv.appendChild(badge);
      notificationDiv.appendChild(contentDiv);
      notificationDiv.appendChild(time);
      parent.appendChild(notificationDiv);

    }).catch(error => {
    console.error("Error fetching notifications:", error);
  });
  })
});