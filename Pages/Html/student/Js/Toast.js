/**
 * Toast Notification Utility
 * 
 * Usage: showToast("Your message", "success" | "error", 3000)
 */

function showToast(message, type = 'success', duration = 3000) {
    // 1. Ensure CSS is loaded
    if (!document.getElementById('toast-styles')) {
        const link = document.createElement('link');
        link.id = 'toast-styles';
        link.rel = 'stylesheet';
        link.href = '/Css/Toast.css';
        document.head.appendChild(link);
    }

    // 2. Ensure Container exists
    let container = document.querySelector('.toast-container');
    if (!container) {
        container = document.createElement('div');
        container.className = 'toast-container';
        document.body.appendChild(container);
    }

    // 3. Create Toast Element
    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;

    // Icon selection
    const icon = type === 'success' ? '✅' : '❌';

    toast.innerHTML = `
        <div class="toast-icon">${icon}</div>
        <div class="toast-message">${message}</div>
    `;

    // 4. Add to container
    container.appendChild(toast);

    // 5. Automatic Removal
    setTimeout(() => {
        toast.classList.add('toast-closing');
        toast.addEventListener('animationend', () => {
            toast.remove();
            // Remove container if empty
            if (container.children.length === 0) {
                container.remove();
            }
        });
    }, duration);
}

// Make it globally accessible
window.showToast = showToast;
