// Letters Receive Page - Backend Integration

// API endpoint configuration
const API_BASE_URL = '/api'; 

/**
 * Fetch all letters from the backend
 */
async function fetchLetters() {
  try {
    const response = await fetch(`${API_BASE_URL}/get-letters`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const data = await response.json();
    return data; // Array of LetterMetadata
  } catch (error) {
    console.error('Error fetching letters:', error);
    showToast("Failed to load letters.", "error");
    return [];
  }
}

/**
 * Download a letter
 * @param {string} letterId - The ID of the letter to download
 */
async function downloadLetter(letterId, filename) {
  try {
    // Backend expects ?id=...
    const downloadUrl = `${API_BASE_URL}/download-letter?id=${letterId}`;
    
    // Simply redirecting or opening in new tab is often easiest for downloads
    // But using fetch allows for better error handling
    const response = await fetch(downloadUrl);
    if (!response.ok) throw new Error("File not found");
    
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename || `letter_${letterId}`; // Use the provided filename
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
 * Send a letter to the benefactor
 * @param {string} letterId - The ID of the letter to send
 */
async function sendLetter(letterId) {
  try {
    const response = await fetch(`${API_BASE_URL}/send-letter?id=${letterId}`, {
      method: 'GET',
    });

    if (!response.ok) {
      const errorMsg = await response.text();
      throw new Error(errorMsg || "Failed to send letter");
    }

    showToast("Letter sent to benefactor successfully!", "success");
  } catch (error) {
    console.error('Error sending letter:', error);
    showToast(error.message || "Failed to send letter.", "error");
  }
}

/**
 * View a specific letter (Proxy for download/opening in new tab)
 */
function viewLetter(letterId) {
  const viewUrl = `${API_BASE_URL}/download-letter?id=${letterId}`;
  window.open(viewUrl, '_blank');
}

/**
 * Populate the table with letters data
 */
function populateLettersTable(letters) {
  const tbody = document.querySelector('#letters tbody');
  if (!tbody) return;

  tbody.innerHTML = '';

  if (letters.length === 0) {
    tbody.innerHTML = '<tr><td colspan="6" style="text-align:center;">No letters found.</td></tr>';
    return;
  }

  letters.forEach(letter => {
    const row = document.createElement('tr');
    // Using green for all received letters for now
    row.className = 'success';

    row.innerHTML = `
      <td>🟢</td>
      <td>${letter.studentName}</td>
      <td>${letter.registrationNumber}</td>
      <td>${letter.letterType}</td>
      <td>${letter.dateSubmitted}</td>
      <td>
        <button class="view-btn" onclick="viewLetter('${letter.id}')"><img class="svg-btn" src="/Image/svgviewer-output (15).svg" alt=""></button>
        <button class="download-btn" onclick="sendLetter('${letter.id}')">Send</button>
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

  // Search functionality (local filtering for now as API search isn't implemented)
  const searchInput = document.querySelector('.search-bar input');
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      const searchTerm = e.target.value.toLowerCase();
      fetchLetters().then(allLetters => {
        const filtered = allLetters.filter(l => 
          l.studentName.toLowerCase().includes(searchTerm) || 
          l.registrationNumber.toLowerCase().includes(searchTerm)
        );
        populateLettersTable(filtered);
      });
    });
  }
});
