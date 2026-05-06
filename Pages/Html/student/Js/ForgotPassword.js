document.getElementById('forgotPasswordForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const submitBtn = e.target.querySelector('button');
    
    submitBtn.disabled = true;
    submitBtn.textContent = 'Sending...';

    try {
        const response = await fetch('/api/forgot-password', {
            method: 'POST',
            body: formData
        });

        const result = await response.json();

        if (response.ok) {
            showToast(result.message || 'Reset link sent!', 'success');
            e.target.reset();
        } else {
            showToast(result.error || 'Failed to send reset link.', 'error');
        }
    } catch (error) {
        showToast('Connection error. Please try again.', 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Send Link';
    }
});
