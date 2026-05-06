document.addEventListener('DOMContentLoaded', async () => {
    // Load current profile data
    try {
        const response = await fetch('/user/profile/detailed');
        if (response.ok) {
            const data = await response.json();
            document.getElementById('firstName').value = data.name;
            document.getElementById('lastName').value = data.surname;
            document.getElementById('email').value = data.email;
            document.getElementById('phone').value = data.phone;
        }
    } catch (error) {
        console.error("Error loading profile:", error);
        showToast("Failed to load profile data", "error");
    }

    // Handle Profile Update
    document.getElementById('profileForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = {
            name: document.getElementById('firstName').value,
            surname: document.getElementById('lastName').value,
            email: document.getElementById('email').value,
            phone: document.getElementById('phone').value
        };

        const submitBtn = e.target.querySelector('button');
        submitBtn.disabled = true;
        submitBtn.textContent = "Saving...";

        try {
            const response = await fetch('/api/user/profile/update', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formData)
            });

            if (response.ok) {
                showToast("Profile updated successfully", "success");
                // Refresh profile card initials if needed
                setTimeout(() => window.location.reload(), 1000);
            } else {
                const err = await response.text();
                showToast(err || "Update failed", "error");
            }
        } catch (error) {
            showToast("Connection error", "error");
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = "Save Changes";
        }
    });

    // Handle Password Change
    document.getElementById('passwordForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const oldPassword = document.getElementById('oldPassword').value;
        const newPassword = document.getElementById('newPassword').value;
        const confirmPassword = document.getElementById('confirmPassword').value;

        if (newPassword !== confirmPassword) {
            showToast("New passwords do not match", "error");
            return;
        }

        const submitBtn = e.target.querySelector('button');
        submitBtn.disabled = true;
        submitBtn.textContent = "Changing...";

        try {
            const response = await fetch('/api/user/profile/change-password', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ oldPassword, newPassword })
            });

            if (response.ok) {
                showToast("Password changed successfully", "success");
                e.target.reset();
            } else {
                const err = await response.text();
                showToast(err || "Password change failed", "error");
            }
        } catch (error) {
            showToast("Connection error", "error");
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = "Change Password";
        }
    });
});
