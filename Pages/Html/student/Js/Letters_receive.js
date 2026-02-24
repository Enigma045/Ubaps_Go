// Letters Receive Page - Backend Integration

// API endpoint configuration
const API_BASE_URL = '/receive_letters'; // Adjust based on your backend configuration

/**
 * Fetch all letters from the backend
 */
async function fetchLetters() {
  try {
    const response = await fetch(`${API_BASE_URL}/letters`, {
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
    console.error('Error fetching letters:', error);
    showToast("Failed to load letters.", "error");
    return [];
  }
}

/**
 * View a specific letter
 * @param {string} letterId - The ID of the letter to view
 */
async function viewLetter(letterId) {
  try {
    const response = await fetch(`${API_BASE_URL}/letters/${letterId}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const letterData = await response.json();
    // Handle displaying the letter (e.g., open in modal or new page)
    console.log('Letter data:', letterData);
    return letterData;
  } catch (error) {
    console.error('Error viewing letter:', error);
    showToast("Failed to load letter details.", "error");
  }
}

/**
 * Download a letter
 * @param {string} letterId - The ID of the letter to download
 */
async function downloadLetter(letterId) {
  try {
    const response = await fetch(`${API_BASE_URL}/letters/${letterId}/download`, {
      method: 'GET',
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `letter_${letterId}.pdf`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
    showToast("Download started!", "success");
  } catch (error) {
    console.error('Error downloading letter:', error);
    showToast("Failed to download letter.", "error");
  }
}

/**
 * Search letters
 * @param {string} searchTerm - The search term
 */
async function searchLetters(searchTerm) {
  try {
    const response = await fetch(`${API_BASE_URL}/letters/search?q=${encodeURIComponent(searchTerm)}`, {
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
    console.error('Error searching letters:', error);
    return [];
  }
}

/**
 * Populate the table with letters data
 * @param {Array} letters - Array of letter objects
 */
function populateLettersTable(letters) {
  const tbody = document.querySelector('#letters tbody');
  if (!tbody) return;

  tbody.innerHTML = '';

  letters.forEach(letter => {
    const row = document.createElement('tr');
    row.className = letter.status === 'pending' ? 'warning' : 'success';

    row.innerHTML = `
      <td>${letter.status === 'pending' ? '🟡' : '🟢'}</td>
      <td>${letter.studentName}</td>
      <td>${letter.registrationNumber}</td>
      <td>${letter.letterType}</td>
      <td>${letter.dateSubmitted}</td>
      <td>
        <button class="view-btn" onclick="viewLetter('${letter.id}')">View</button>
        <button class="download-btn" onclick="downloadLetter('${letter.id}')">Download</button>
      </td>
    `;

    tbody.appendChild(row);
  });
}

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
  // Load letters on page load
  fetchLetters().then(letters => {
    populateLettersTable(letters);
  });

  // Search functionality
  const searchInput = document.querySelector('.search-bar input');
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      const searchTerm = e.target.value;
      if (searchTerm.length > 2) {
        searchLetters(searchTerm).then(letters => {
          populateLettersTable(letters);
        });
      } else if (searchTerm.length === 0) {
        fetchLetters().then(letters => {
          populateLettersTable(letters);
        });
      }
    });
  }
});
