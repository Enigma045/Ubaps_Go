// Financial Reports Page - Backend Integration

// API endpoint configuration
const API_BASE_URL = '/financial_reports'; // Adjust based on your backend configuration

/**
 * Fetch financial reports from the backend
 * @param {Object} filters - Optional filters (startDate, endDate)
 */
async function fetchFinancialReports(filters = {}) {
    try {
        const queryParams = new URLSearchParams(filters).toString();
        const url = `${API_BASE_URL}/financial-reports${queryParams ? '?' + queryParams : ''}`;

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
        console.error('Error fetching financial reports:', error);
        showToast("Failed to load financial reports.", "error");
        return [];
    }
}

/**
 * Apply date range filter
 */
function applyDateFilter() {
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;

    const filters = {};
    if (startDate) filters.startDate = startDate;
    if (endDate) filters.endDate = endDate;

    fetchFinancialReports(filters).then(reports => {
        populateReportsTable(reports);
    });
}

/**
 * View report details
 * @param {string} reportId - The ID of the report to view
 */
async function viewReportDetails(reportId) {
    try {
        const response = await fetch(`${API_BASE_URL}/financial-reports/${reportId}`, {
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
        showToast("Failed to load report details.", "error");
    }
}

/**
 * Export report to CSV/PDF
 * @param {string} reportId - The ID of the report to export (optional)
 */
async function exportReport(reportId = null) {
    try {
        const url = reportId
            ? `${API_BASE_URL}/financial-reports/${reportId}/export`
            : `${API_BASE_URL}/financial-reports/export`;

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
        a.download = `financial_report_${reportId || 'all'}_${new Date().toISOString().split('T')[0]}.csv`;
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
        const response = await fetch(`${API_BASE_URL}/financial-reports/search?q=${encodeURIComponent(searchTerm)}`, {
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
    const tbody = document.querySelector('#financial-reports tbody');
    if (!tbody) return;

    tbody.innerHTML = '';

    reports.forEach(report => {
        const row = document.createElement('tr');
        row.className = 'success';

        row.innerHTML = `
      <td>${report.reportId}</td>
      <td>${report.period}</td>
      <td>MWK ${report.totalDisbursed.toLocaleString()}</td>
      <td>MWK ${report.totalPending.toLocaleString()}</td>
      <td>${report.totalSchemes}</td>
      <td>${report.dateGenerated}</td>
      <td>
        <button class="view-btn" onclick="viewReportDetails('${report.id}')">View Details</button>
        <button class="export-btn" onclick="exportReport('${report.id}')">Export</button>
      </td>
    `;

        tbody.appendChild(row);
    });
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    // Load reports on page load
    fetchFinancialReports().then(reports => {
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
                fetchFinancialReports().then(reports => {
                    populateReportsTable(reports);
                });
            }
        });
    }
});
