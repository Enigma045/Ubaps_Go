/* Admin Dashboard Logic */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("Admin Dashboard Initialized");


    // Load Stats
    if (typeof fetchAdminStats === 'function') {
        const stats = await fetchAdminStats();
        if (stats) {
            const totalUsersEl = document.getElementById('total-users');
            const activeUsersEl = document.getElementById('active-users');
            const deactiveUsersEl = document.getElementById('deactive-users');

            if (totalUsersEl) totalUsersEl.textContent = stats.total_users;
            if (activeUsersEl) activeUsersEl.textContent = stats.active_users;
            if (deactiveUsersEl) deactiveUsersEl.textContent = stats.deactive_users;
        }
    }
});