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
  },
  {
    title: "Not Selected",
    description: "Unfortunately, you have not been selected for a bursary award at this time."
  }
];

const points = document.querySelectorAll(".point:not(.future)");
const progress = document.querySelector(".timeline-progress");

const title = document.getElementById("title");
const date = document.getElementById("date");
const description = document.getElementById("description");

function updateProgress(index) {
  if (points.length === 0) return;

  if (index === 4){
  // Update bar percentage
  if (progress) {
    const percentage = (3 / (points.length - 1)) * 100;
    progress.style.width = `${percentage}%`;
  }
  }else{
    if (progress) {
    const percentage = (index / (points.length - 1)) * 100;
    progress.style.width = `${percentage}%`;
    }
  }
  // Update dots
  points.forEach((p, i) => {
    if (i <= index) p.classList.add("active");
    else p.classList.remove("active");
  });

  // Update content (Automatic update)
  const event = events[index];
  console.log(event);
  if (event) {
    if (title) title.textContent = event.title;
    if (description) description.textContent = event.description;
    // Removing date update as it might be static or not needed in this context
  }
}

points.forEach((point, index) => {
  point.addEventListener("click", () => {
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

      if (lowerStatus.includes('not selected')) {
        index = 4;

        const lastPointLabel = points[3]?.querySelector('.date');
        if (lastPointLabel) lastPointLabel.textContent = "Rejected";

      } else if (lowerStatus.includes('selected')) {
        index = 3;
        // Update the label on the dot itself
        
      } else if (lowerStatus.includes('considering')) {
        index = 2;
      }else if (lowerStatus.includes('not submitted')) {
        index = 0;
      }else if (lowerStatus.includes('submitted')) {
        index = 1;
      } else {
        index = 0;
      }

      // Update timeline UI
      updateProgress(index);
    }
  }
});