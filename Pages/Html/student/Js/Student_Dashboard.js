// Option A uses pure CSS hover
// JavaScript kept minimal for future use


const events = [
  {
    title: "Not Submitted",
    description: "Please submit your application form."
  },
  {
    title: "Submitted",
    description: "You have submitted your application form please wait for the committee to review it."
  },
  {
    title: "Considering",
    description: "You application form is being reviewed"
  },
  {
    title: "Selected",
    description: "You have been selected for the schollaship program please write a thank you write and submit it in the students letter page."
  }
];

const points = document.querySelectorAll(".point:not(.future)");
const progress = document.querySelector(".timeline-progress");

const title = document.getElementById("title");
const date = document.getElementById("date");
const description = document.getElementById("description");

function updateProgress(index) {
  if (!progress || points.length === 0) return;
  const percentage = (index / (points.length - 1)) * 100;
  progress.style.width = `${percentage}%`;
  events[index]
}

points.forEach((point, index) => {
  point.addEventListener("click", () => {
    // Active dot
    points.forEach(p => p.classList.remove("active"));
    point.classList.add("active");

    // Content
    if (title) title.textContent = events[index].title;
    if (date) date.textContent = events[index].date;
    if (description) description.textContent = events[index].description;

    // Animate progress
    updateProgress(index);
  });
});

// Initialize on load
document.addEventListener('DOMContentLoaded', async () => {
  updateProgress(0);


  // Load Stats
  if (typeof fetchStudentStats === 'function') {
    const stats = await fetchStudentStats();
    if (stats) {
      const appStatusEl = document.getElementById('app-status');
      const bursarySchemeEl = document.getElementById('bursary-scheme');

      if (appStatusEl) appStatusEl.textContent = stats.application_status;
      if (bursarySchemeEl) bursarySchemeEl.textContent = stats.bursary_scheme;

      // Update timeline based on status
      let index = 0;
      const lowerStatus = stats.application_status.toLowerCase();

      if (lowerStatus.includes('selected')) index = 3;
      else if (lowerStatus.includes('considering') || lowerStatus.includes('considering')) index = 2;
      else if (lowerStatus.includes('submitted') || lowerStatus.includes('submitted')) index = 1;
      else if (lowerStatus.includes('not submitted') || lowerStatus.includes('submitted')) index = 0;

      // Update timeline UI
      const targetPoint = points[index];
      if (targetPoint) {
        points.forEach(p => p.classList.remove("active"));
        targetPoint.classList.add("active");
        updateProgress(index);
      }
    }
  }
});