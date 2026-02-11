document.addEventListener('DOMContentLoaded', function() {
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

    // Toggle function
    bellIcon.addEventListener('click', function(e) {
        e.stopPropagation();
        container.classList.toggle('active');
        if (container.classList.contains('active')) {
            fetchNotifications();
        }
    });

    // Close when clicking outside
    document.addEventListener('click', function(e) {
        if (!container.contains(e.target) && !bellIcon.contains(e.target)) {
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
            list.innerHTML = ''; // Clear loading/old

            if (!data || data.length === 0) {
                list.innerHTML = '<div class="n-empty">No new notifications</div>';
                countDisplay.textContent = '0 New';
                return;
            }

            countDisplay.textContent = `${data.length} New`;
            
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
            list.innerHTML = '<div class="n-empty">Failed to load</div>';
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
                 // Assuming badge styling in CSS handles visibility (e.g. empty or 0 might be hidden)
                 // Or typically JS handles display: block/none
                 if (count > 0) {
                     badge.style.display = 'inline-block'; // or flex, depends on CSS. Unified_Dashboard said absolute/flex
                     badge.style.display = 'flex'; // Force flex for centering
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
});
