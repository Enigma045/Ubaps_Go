/* Global Profile Loader Logic */

document.addEventListener('DOMContentLoaded', async () => {
    // Dynamically load profile card styles
    if (!document.getElementById('profile-card-styles')) {
        const link = document.createElement('link');
        link.id = 'profile-card-styles';
        link.rel = 'stylesheet';
        link.href = '/Css/Profile_Card.css';
        document.head.appendChild(link);
    }


    // Wait slightly to ensure Dashboard.js is loaded if scripts are not in order
    // though ideally they should be correctly ordered in HTML.
    if (typeof fetchUserProfile === 'function') {
        try {
            const profile = await fetchUserProfile();
            if (profile && profile.name) {
                // Update Badge Initials
                const avatar = document.getElementById('user-avatar');
                if (avatar) {
                    const names = profile.name.trim().split(/\s+/);
                    let initials = "";
                    if (names.length > 0) {
                        initials = names[0][0];
                        if (names.length > 1) {
                            initials += names[names.length - 1][0];
                        }
                    }
                    avatar.textContent = initials.toUpperCase();
                    avatar.title = profile.name; // Full name on hover
                }

                // Update Welcome Message if it exists on the page (common in dashboards)
                const welcomeMsg = document.getElementById('welcome-message');
                if (welcomeMsg) {
                    welcomeMsg.textContent = `Welcome, ${profile.name}`;
                }

                // Create Profile Card Container if not exists
                if (!document.querySelector('.profile-container')) {
                    const profileContainer = document.createElement('div');
                    profileContainer.className = 'profile-container';

                    // Format role for display
                    const displayRole = profile.role.replace(/_/g, ' ');

                    profileContainer.innerHTML = `
                        <div class="profile-card">
                            <div class="profile-info">
                                <div class="p-avatar">${avatar.textContent}</div>
                                <h3>${profile.name}</h3>
                                <p>${displayRole}</p>
                            </div>
                            <div class="profile-actions">
                                <a href="/user/profile-settings" class="profile-link">
                                    <span>My Profile</span>
                                </a>
                                <a href="/logout" class="profile-link">
                                    <span>Logout</span>
                                </a>
                            </div>
                        </div>
                    `;
                    document.body.appendChild(profileContainer);
                }

                const profileCard = document.querySelector('.profile-container');

                // Toggle Profile Card
                avatar.addEventListener('click', (e) => {
                    e.stopPropagation();
                    profileCard.classList.toggle('active');
                });

                // Close when clicking outside
                document.addEventListener('click', (e) => {
                    if (!profileCard.contains(e.target) && !avatar.contains(e.target)) {
                        profileCard.classList.remove('active');
                    }
                });
            }
        } catch (error) {
            console.error("Profile Loader Error:", error);
        }
    } else {
        console.warn("fetchUserProfile function not found. Ensure Dashboard.js is included.");
    }
});
