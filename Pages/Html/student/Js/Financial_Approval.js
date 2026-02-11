// Financial Approval Page - Backend Integration

// API endpoint configuration
const API_BASE_URL = '/financial_approval'; // Adjust based on your backend configuration

/**
 * Fetch all financial requests from the backend
 */
async function fetchFinancialRequests() {
  try {
    const response = await fetch(`${API_BASE_URL}/financial-requests`, {
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
    console.error('Error fetching financial requests:', error);
    return [];
  }
}

/**
 * Approve a financial request
 * @param {string} requestId - The ID of the request to approve
 */
async function approveRequest(requestId) {
  try {
    const response = await fetch(`${API_BASE_URL}/financial-requests/${requestId}/approve`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        approvedBy: 'current_user', // Replace with actual user info
        approvedAt: new Date().toISOString(),
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    console.log('Request approved:', result);

    // Refresh the table
    fetchFinancialRequests().then(requests => {
      populateRequestsTable(requests);
    });

    return result;
  } catch (error) {
    console.error('Error approving request:', error);
    alert('Failed to approve request. Please try again.');
  }
}

/**
 * Reject a financial request
 * @param {string} requestId - The ID of the request to reject
 */
async function rejectRequest(requestId) {
  try {
    const reason = prompt('Please provide a reason for rejection:');
    if (!reason) return; // User cancelled

    const response = await fetch(`${API_BASE_URL}/financial-requests/${requestId}/reject`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        rejectedBy: 'current_user', // Replace with actual user info
        rejectedAt: new Date().toISOString(),
        reason: reason,
      }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    console.log('Request rejected:', result);

    // Refresh the table
    fetchFinancialRequests().then(requests => {
      populateRequestsTable(requests);
    });

    return result;
  } catch (error) {
    console.error('Error rejecting request:', error);
    alert('Failed to reject request. Please try again.');
  }
}

/**
 * Search financial requests
 * @param {string} searchTerm - The search term
 */
async function searchRequests(searchTerm) {
  try {
    const response = await fetch(`${API_BASE_URL}/financial-requests/search?q=${encodeURIComponent(searchTerm)}`, {
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
    console.error('Error searching requests:', error);
    return [];
  }
}

/**
 * Populate the table with financial requests data
 * @param {Array} requests - Array of request objects
 */
function populateRequestsTable(requests) {
  const tbody = document.querySelector('#financial-requests tbody');
  if (!tbody) return;

  tbody.innerHTML = '';

  requests.forEach(request => {
    const row = document.createElement('tr');

    // Set row class based on status
    if (request.status === 'approved') {
      row.className = 'success';
    } else if (request.status === 'rejected') {
      row.className = 'error';
    } else {
      row.className = 'warning';
    }

    // Determine status icon
    let statusIcon = '🟡'; // pending
    if (request.status === 'approved') statusIcon = '🟢';
    if (request.status === 'rejected') statusIcon = '🔴';

    // Build actions column based on status
    let actionsHTML = '';
    if (request.status === 'pending') {
      actionsHTML = `
        <button class="approve-btn" onclick="approveRequest('${request.id}')">Approve</button>
        <button class="reject-btn" onclick="rejectRequest('${request.id}')">Reject</button>
      `;
    } else if (request.status === 'approved') {
      actionsHTML = '<span class="status-badge approved">Approved</span>';
    } else if (request.status === 'rejected') {
      actionsHTML = '<span class="status-badge rejected">Rejected</span>';
    }

    row.innerHTML = `
      <td>${statusIcon}</td>
      <td>${request.studentName}</td>
      <td>${request.registrationNumber}</td>
      <td>MWK ${request.amount.toLocaleString()}</td>
      <td>${request.requestType}</td>
      <td>${request.date}</td>
      <td>${actionsHTML}</td>
    `;

    tbody.appendChild(row);
  });
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
  // Load financial requests on page load
  fetchFinancialRequests().then(requests => {
    populateRequestsTable(requests);
  });

  // Search functionality
  const searchInput = document.querySelector('.search-bar input');
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      const searchTerm = e.target.value;
      if (searchTerm.length > 2) {
        searchRequests(searchTerm).then(requests => {
          populateRequestsTable(requests);
        });
      } else if (searchTerm.length === 0) {
        fetchFinancialRequests().then(requests => {
          populateRequestsTable(requests);
        });
      }
    });
  }
});