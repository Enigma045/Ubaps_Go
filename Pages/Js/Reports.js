const events = [
  {
    title: "First milestone",
    date: "– January 16th, 2014",
    description: "This marks the beginning of the project with initial planning."
  },
  {
    title: "Second milestone",
    date: "– February 28th, 2014",
    description: "Core development phase with key features implemented."
  },
  {
    title: "Main event",
    date: "– March 20th, 2014",
    description: "Official launch with successful deployment."
  },
  {
    title: "Main event",
    date: "– March 20th, 2014",
    description: "Official launch with successful deployment."
  }
];

const points = document.querySelectorAll(".point:not(.future)");
const progress = document.querySelector(".timeline-progress");

const title = document.getElementById("title");
const date = document.getElementById("date");
const description = document.getElementById("description");

function updateProgress(index) {
  const percentage = (index / (points.length - 1)) * 100;
  progress.style.width = `${percentage}%`;
}

points.forEach((point, index) => {
  point.addEventListener("click", () => {
    // Active dot
    points.forEach(p => p.classList.remove("active"));
    point.classList.add("active");

    // Content
    title.textContent = events[index].title;
    date.textContent = events[index].date;
    description.textContent = events[index].description;

    // Animate progress
    updateProgress(index);
  });
});

// Initialize on load
updateProgress(0);
