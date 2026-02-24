// Student Reports Page - Backend Integration

// API endpoint configuration
const API_BASE_URL = '/student_report'; // Adjust based on your backend configuration

/**
 * Fetch student reports from the backend
 * @param {Object} filters - Optional filters (program, status)
 */
async function fetchStudentReports(filters = {}) {
    try {
        const queryParams = new URLSearchParams(filters).toString();
        const url = `${API_BASE_URL}/student-reports${queryParams ? '?' + queryParams : ''}`;

        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching student reports:', error);
        showToast("Failed to load student reports.", "error");
        return [];
    }
}

/**
 * Apply filters (program and status)
 */
function applyFilters() {
    const program = document.getElementById('programFilter').value;
    const status = document.getElementById('statusFilter').value;

    const filters = {};
    if (program) filters.program = program;
    if (status) filters.status = status;

    fetchStudentReports(filters).then(reports => {
        populateReportsTable(reports);
    });
}

/**
 * View student details
 * @param {string} studentId - The ID of the student to view
 */
async function viewStudentDetails(studentId) {
    try {
        const response = await fetch(`${API_BASE_URL}/student-reports/${studentId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const studentData = await response.json();
        console.log('Student details:', studentData);
        // Handle displaying the student details (e.g., open in modal or new page)
        return studentData;
    } catch (error) {
        console.error('Error viewing student details:', error);
        showToast("Failed to load student details.", "error");
    }
}

/**
 * Export report to CSV/PDF
 */
async function exportReport() {
    try {
        const program = document.getElementById('programFilter').value;
        const status = document.getElementById('statusFilter').value;

        const filters = {};
        if (program) filters.program = program;
        if (status) filters.status = status;

        const queryParams = new URLSearchParams(filters).toString();
        const url = `${API_BASE_URL}/student-reports/export${queryParams ? '?' + queryParams : ''}`;

        const response = await fetch(url, {
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const blob = await response.blob();
        const downloadUrl = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = downloadUrl;
        a.download = `student_report_${new Date().toISOString().split('T')[0]}.csv`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(downloadUrl);
        document.body.removeChild(a);
    } catch (error) {
        console.error('Error exporting report:', error);
        showToast('Failed to export report. Please try again.', 'error');
    }
}

/**
 * Search reports
 * @param {string} searchTerm - The search term
 */
async function searchReports(searchTerm) {
    try {
        const response = await fetch(`${API_BASE_URL}/student-reports/search?q=${encodeURIComponent(searchTerm)}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error searching reports:', error);
        showToast("Failed to search reports.", "error");
        return [];
    }
}

/**
 * Populate the table with reports data
 * @param {Array} reports - Array of report objects
 */
function populateReportsTable(reports) {
    const tbody = document.querySelector('#student-reports tbody');
    if (!tbody) return;

    tbody.innerHTML = '';

    reports.forEach(report => {
        const row = document.createElement('tr');
        row.className = report.bursaryStatus === 'approved' ? 'success' : 'warning';

        const statusBadgeClass = report.bursaryStatus === 'approved' ? 'approved' :
            report.bursaryStatus === 'rejected' ? 'rejected' : 'pending';

        row.innerHTML = `
      <td>${report.studentName}</td>
      <td>${report.registrationNumber}</td>
      <td>${report.program}</td>
      <td><span class="status-badge ${statusBadgeClass}">${report.bursaryStatus}</span></td>
      <td>MWK ${report.amountReceived.toLocaleString()}</td>
      <td>${report.academicYear}</td>
      <td>
        <button class="view-btn" onclick="viewStudentDetails('${report.id}')">View Details</button>
      </td>
    `;

        tbody.appendChild(row);
    });
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    // Load reports on page load
    fetchStudentReports().then(reports => {
        populateReportsTable(reports);
    });

    // Search functionality
    const searchInput = document.querySelector('.search-bar input');
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            const searchTerm = e.target.value;
            if (searchTerm.length > 2) {
                searchReports(searchTerm).then(reports => {
                    populateReportsTable(reports);
                });
            } else if (searchTerm.length === 0) {
                fetchStudentReports().then(reports => {
                    populateReportsTable(reports);
                });
            }
        });
    }
});
