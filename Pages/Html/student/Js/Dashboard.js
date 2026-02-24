/**
 * Dashboard Statistics Functions
 * 
 * These functions fetch data for the various dashboard stat cards.
 * They are used across different user roles.
 */

/**
 * Common fetch helper to handle responses and common errors
 * @param {string} url - The endpoint to fetch from
 * @returns {Promise<Object>} - The JSON response
 */
async function fetchDashboardData(url) {
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            }
        });

        if (!response.ok) {
            if (response.status === 401) {
                console.error('Unauthorized access. Redirecting to login.');
                window.location.href = '/Login';
                return null;
            }
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return await response.json();
    } catch (error) {
        console.error(`Error fetching data from ${url}:`, error);
        return null;
    }
}

/**
 * Fetches statistics for the student dashboard.
 * Retrieves: application status and bursary scheme.
 */
async function fetchStudentStats() {
    return await fetchDashboardData('/stats/student');
}

/**
 * Fetches statistics for the registrar dashboard.
 * Retrieves: approved amount, number of applicants, and number of bursary schemes.
 */
async function fetchRegistrarStats() {
    return await fetchDashboardData('/stats/registrar');
}

/**
 * Fetches statistics for the admin dashboard.
 * Retrieves: total users, active users, and deactive users.
 */
async function fetchAdminStats() {
    return await fetchDashboardData('/stats/admin');
}


/**
 * Fetches statistics for the dean's dashboard.
 * Retrieves: pending applications, selected students, rejected students, and pending letters.
 */
async function fetchDeanStats() {
    return await fetchDashboardData('/stats/dean');
}

/**
 * Fetches statistics for the financial officer's dashboard.
 * Retrieves: approved amount, number of disbursements made, and financial history requests.
 */
async function fetchFinancialOfficerStats() {
    return await fetchDashboardData('/stats/finance');
}

/**
 * Fetches the current user's name and role information.
 */
async function fetchUserProfile() {
    return await fetchDashboardData('/user/profile');
}

// Example usage / Initialization if needed
// (Specific pages will call these functions as required by their layout)
