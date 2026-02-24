/* Dean Dashboard Logic */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("Dean Dashboard Initialized");


    // Load Stats
    if (typeof fetchDeanStats === 'function') {
        const stats = await fetchDeanStats();
        if (stats) {
            const pendingEl = document.getElementById('num-pending');
            const selectedEl = document.getElementById('num-selected');
            const rejectedEl = document.getElementById('num-rejected');
            const lettersEl = document.getElementById('num-letters');

            if (pendingEl) pendingEl.textContent = stats.pending_applications + stats.considering_applications;
            if (selectedEl) selectedEl.textContent = stats.selected_students;
            if (rejectedEl) rejectedEl.textContent = stats.rejected_students;
            if (lettersEl) lettersEl.textContent = stats.pending_letters;
        }
    }
});
