// Extract token from URL
const urlParams = new URLSearchParams(window.location.search);
const token = urlParams.get('token');

if (!token) {
    showToast('Invalid or missing reset token.', 'error');
    document.getElementById('resetPasswordForm').style.display = 'none';
} else {
    document.getElementById('token').value = token;
}

document.getElementById('resetPasswordForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const password = formData.get('password');
    const confirmPassword = formData.get('confirmPassword');
    const submitBtn = e.target.querySelector('button');

    if (password !== confirmPassword) {
        showToast('Passwords do not match.', 'error');
        return;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = 'Updating...';

    try {
        const response = await fetch('/api/reset-password', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (response.ok) {
            showToast('Password updated successfully!', 'success');
            document.getElementById('resetPasswordForm').innerHTML = `
                <div style="text-align: center; padding: 20px;">
                    <p style="color: var(--success-color); font-weight: 600; margin-bottom: 20px;">Your password has been reset.</p>
                    <a href="/Login" class="register-btn" style="text-decoration: none; display: inline-block;">Login Now</a>
                </div>
            `;
        } else {
            showToast(result.error || 'Failed to reset password.', 'error');
        }
    } catch (error) {
        showToast('Connection error. Please try again.', 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Update Password';
    }
});
