document.addEventListener('DOMContentLoaded', function () {
    const bellIcon = document.querySelector('.notification-nav');

    // Create notification container if not exists
    if (!document.querySelector('.notification-container')) {
        const container = document.createElement('div');
        container.className = 'notification-container';
        container.innerHTML = `
            <div class="notification-card">
                <div class="notification-header">
                    <h3>Notifications</h3>
                    <span class="notification-count" id="notif-count-display">0 New</span>
                </div>
                <div class="notification-list" id="notification-list">
                    <div class="n-empty">Loading...</div>
                </div>
            </div>
        `;
        // Append to body or header to ensure it's positioned relative to screen/header
        // Best to append to header or main wrapper if relative, but absolute to body works with correct top/right
        document.body.appendChild(container); // or append to .header if header is relative
    }

    const container = document.querySelector('.notification-container');
    const list = document.getElementById('notification-list');
    const countDisplay = document.getElementById('notif-count-display');
    const badge = document.getElementById('badge-nav'); // Existing badge on bell

    if (bellIcon) {
        // Toggle function
        bellIcon.addEventListener('click', function (e) {
            e.stopPropagation();
            container.classList.toggle('active');
            if (container.classList.contains('active')) {
                fetchNotifications();
            }
        });
    }

    // Close when clicking outside
    document.addEventListener('click', function (e) {
        if (container && !container.contains(e.target) && (bellIcon && !bellIcon.contains(e.target))) {
            container.classList.remove('active');
        }
    });

    function fetchNotifications() {
        console.log("Fetching notifications...");
        fetch("/getnotifications", {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        })
            .then(response => response.json())
            .then(data => {
                console.log("Notifications data:", data);
                if (!list) return;
                list.innerHTML = ''; // Clear loading/old

                if (!data || data.length === 0) {
                    list.innerHTML = '<div class="n-empty">No new notifications</div>';
                    if (countDisplay) countDisplay.textContent = '0 New';
                    return;
                }

                if (countDisplay) countDisplay.textContent = `${data.length} New`;

                // Update nav badge too
                if (badge) {
                    badge.textContent = data.length;
                    badge.style.display = data.length > 0 ? 'flex' : 'none'; // Ensure generic badge display
                }

                data.forEach(n => {
                    const item = document.createElement('div');
                    item.className = 'notification-item';
                    item.innerHTML = `
                    <div class="notif-icon">🔔</div>
                    <div class="notif-content">
                        <div class="notif-title">${n.title || 'Notification'}</div>
                        <div class="notif-message">${n.message || 'No description'}</div>
                        <div class="notif-time">${n.timestamp || 'Just now'}</div>
                    </div>
                `;
                    list.appendChild(item);
                });
            })
            .catch(err => {
                console.error(err);
                if (list) list.innerHTML = '<div class="n-empty">Failed to load</div>';
                if (typeof showToast === 'function') showToast("Failed to fetch notifications.", "error");
            });
    }

    // Initial Badge Fetch (Separate from click to show badge on load)
    function updateBadge() {
        fetch('/countnotifications', {
            headers: { 'Content-Type': 'application/json' }
        })
            .then(res => res.json())
            .then(count => {
                if (badge) {
                    badge.textContent = count;
                    if (count > 0) {
                        badge.style.display = 'flex'; 
                        badge.style.alignItems = 'center';
                        badge.style.justifyContent = 'center';
                    } else {
                        badge.style.display = 'none';
                    }
                }
            })
            .catch(e => console.error("Badge error", e));
    }

    // Run badge update on load
    updateBadge();

    // Mapping Dashboard Notification Card to Bell Click
    const dashNotifCard = Array.from(document.querySelectorAll('.action-card')).find(c => {
        const h3 = c.querySelector('h3');
        return h3 && (h3.textContent.trim() === 'Notifications' || h3.textContent.trim() === 'Notification');
    });
    
    if (dashNotifCard && bellIcon) {
        dashNotifCard.addEventListener('click', () => {
            bellIcon.click();
        });
    }
});
