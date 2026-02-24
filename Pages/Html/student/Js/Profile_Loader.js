/* Global Profile Loader Logic */

document.addEventListener('DOMContentLoaded', async () => {
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
                    // Check if it already has text or if we should just prepend/replace
                    // Standardizing to "Welcome, [Name]"
                    welcomeMsg.textContent = `Welcome, ${profile.name}`;
                }
            }
        } catch (error) {
            console.error("Profile Loader Error:", error);
        }
    } else {
        console.warn("fetchUserProfile function not found. Ensure Dashboard.js is included.");
    }
});
