  let count = 0;
  const badge = document.getElementById("badge");

  function increment() {
    count++;
    badge.textContent = count;
    badge.style.display = "inline-block";
  }
