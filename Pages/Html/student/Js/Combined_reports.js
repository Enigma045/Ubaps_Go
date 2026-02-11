// Combined Reports Page - Backend Integration

// API endpoint configuration
const API_BASE_URL = '/general_report'; // Adjust based on your backend configuration

/**
 * Fetch combined reports from the backend
 * @param {Object} filters - Optional filters (startDate, endDate, scheme, status)
 */
async function fetchCombinedReports(filters = {}) {
    try {
        const queryParams = new URLSearchParams(filters).toString();
        const url = `${API_BASE_URL}/combined-reports${queryParams ? '?' + queryParams : ''}`;

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
        console.error('Error fetching combined reports:', error);
        return [];
    }
}

/**
 * Apply all filters (date range, scheme, status)
 */
function applyFilters() {
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;
    const scheme = document.getElementById('schemeFilter').value;
    const status = document.getElementById('statusFilter').value;

    const filters = {};
    if (startDate) filters.startDate = startDate;
    if (endDate) filters.endDate = endDate;
    if (scheme) filters.scheme = scheme;
    if (status) filters.status = status;

    fetchCombinedReports(filters).then(reports => {
        populateReportsTable(reports);
    });
}

/**
 * View combined report details
 * @param {string} reportId - The ID of the report to view
 */
async function viewReportDetails(reportId) {
    try {
        const response = await fetch(`${API_BASE_URL}/combined-reports/${reportId}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const reportData = await response.json();
        console.log('Report details:', reportData);
        // Handle displaying the report details (e.g., open in modal or new page)
        return reportData;
    } catch (error) {
        console.error('Error viewing report details:', error);
    }
}

/**
 * Export combined report to CSV/PDF
 */
async function exportReport() {
    try {
        const startDate = document.getElementById('startDate').value;
        const endDate = document.getElementById('endDate').value;
        const scheme = document.getElementById('schemeFilter').value;
        const status = document.getElementById('statusFilter').value;

        const filters = {};
        if (startDate) filters.startDate = startDate;
        if (endDate) filters.endDate = endDate;
        if (scheme) filters.scheme = scheme;
        if (status) filters.status = status;

        const queryParams = new URLSearchParams(filters).toString();
        const url = `${API_BASE_URL}/combined-reports/export${queryParams ? '?' + queryParams : ''}`;

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
        a.download = `combined_report_${new Date().toISOString().split('T')[0]}.csv`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(downloadUrl);
        document.body.removeChild(a);
    } catch (error) {
        console.error('Error exporting report:', error);
        alert('Failed to export report. Please try again.');
    }
}

/**
 * Search reports
 * @param {string} searchTerm - The search term
 */
async function searchReports(searchTerm) {
    try {
        const response = await fetch(`${API_BASE_URL}/combined-reports/search?q=${encodeURIComponent(searchTerm)}`, {
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
        return [];
    }
}

/**
 * Populate the table with reports data
 * @param {Array} reports - Array of report objects
 */
function populateReportsTable(reports) {
    const tbody = document.querySelector('#combined-reports tbody');
    if (!tbody) return;

    tbody.innerHTML = '';

    reports.forEach(report => {
        const row = document.createElement('tr');
        row.className = report.financialStatus === 'disbursed' ? 'success' : 'warning';

        const appStatusClass = report.applicationStatus === 'approved' ? 'approved' :
            report.applicationStatus === 'rejected' ? 'rejected' : 'pending';

        const finStatusClass = report.financialStatus === 'disbursed' ? 'disbursed' : 'pending';

        row.innerHTML = `
      <td>${report.studentName}</td>
      <td>${report.registrationNumber}</td>
      <td>${report.schemeName}</td>
      <td>MWK ${report.amount.toLocaleString()}</td>
      <td><span class="status-badge ${appStatusClass}">${report.applicationStatus}</span></td>
      <td><span class="status-badge ${finStatusClass}">${report.financialStatus}</span></td>
      <td>${report.date}</td>
      <td>
        <button class="view-btn" onclick="viewReportDetails('${report.id}')">View Details</button>
      </td>
    `;

        tbody.appendChild(row);
    });
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    // Load reports on page load
    fetchCombinedReports().then(reports => {
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
                fetchCombinedReports().then(reports => {
                    populateReportsTable(reports);
                });
            }
        });
    }
});
