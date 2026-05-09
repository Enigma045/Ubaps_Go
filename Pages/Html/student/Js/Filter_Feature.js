/**
 * Filter_Feature.js
 * Advanced Table Filtering Utility
 */

document.addEventListener('DOMContentLoaded', () => {
    const filterSearch = document.getElementById('filter-search');
    const filterDept = document.getElementById('filter-dept');
    const filterPriority = document.getElementById('filter-priority');
    const filterScheme = document.getElementById('filter-scheme');
    const filterParent = document.getElementById('filter-parent');
    const filterEmployment = document.getElementById('filter-employment');
    const filterGender = document.getElementById('filter-gender');
    const clearBtn = document.getElementById('clear-filters');

    const inputs = [filterSearch, filterDept, filterPriority, filterScheme, filterParent, filterEmployment, filterGender];

    function applyFilters() {
        const searchValue = filterSearch ? filterSearch.value.toLowerCase() : '';
        const deptValue = filterDept ? filterDept.value.toUpperCase() : '';
        const priorityValue = filterPriority ? filterPriority.value.toLowerCase() : '';
        const schemeValue = filterScheme ? filterScheme.value.toLowerCase() : '';
        const parentValue = filterParent ? filterParent.value.toLowerCase() : '';
        const employmentValue = filterEmployment ? filterEmployment.value.toLowerCase() : '';
        const genderValue = filterGender ? filterGender.value.toLowerCase() : '';

        // Support both ID based table finding or just all visible tables
        const tables = document.querySelectorAll('.log-table');
        
        tables.forEach(table => {
            const rows = table.querySelectorAll('tbody tr');
            
            rows.forEach(row => {
                const text = row.innerText.toLowerCase();
                const cells = row.cells;
                
                // If there's no data in cells, it might be a loading row or empty state
                if (cells.length < 2) return;

                // Column Mapping (Approximated based on provided HTML files)
                // Name (0), Reg # (1), Parent (2), Employment (3), Scheme (varies), Priority (varies)
                
                let matchesSearch = text.includes(searchValue);
                let matchesDept = deptValue === '' || cells[1].innerText.toUpperCase().includes(deptValue);
                
                // Find Priority column (often named "Priority" in header)
                let priorityCell = findCellByHeader(table, "Priority", row);
                let matchesPriority = priorityValue === '' || (priorityCell && priorityCell.innerText.toLowerCase().includes(priorityValue));

                // Find Parent column
                let parentCell = findCellByHeader(table, "Parent Status", row);
                let matchesParent = parentValue === '' || (parentCell && parentCell.innerText.toLowerCase().includes(parentValue));

                // Find Employment column
                let employmentCell = findCellByHeader(table, "Employment Status", row);
                let matchesEmployment = employmentValue === '' || (employmentCell && employmentCell.innerText.toLowerCase().includes(employmentValue));

                // Find Gender column
                let genderCell = findCellByHeader(table, "Gender", row);
                let matchesGender = genderValue === '' || (genderCell && genderCell.innerText.toLowerCase() === genderValue);

                // Find Scheme column
                let schemeCell = findCellByHeader(table, "Bursary Scheme", row) || findCellByHeader(table, "Letter Type", row);
                let matchesScheme = schemeValue === '' || (schemeCell && schemeCell.innerText.toLowerCase().includes(schemeValue));

                if (matchesSearch && matchesDept && matchesPriority && matchesParent && matchesEmployment && matchesScheme && matchesGender) {
                    row.style.display = '';
                } else {
                    row.style.display = 'none';
                }
            });
        });
    }

    function findCellByHeader(table, headerText, row) {
        const headers = Array.from(table.querySelectorAll('thead th'));
        const index = headers.findIndex(h => h.innerText.toLowerCase().trim() === headerText.toLowerCase().trim());
        return index !== -1 ? row.cells[index] : null;
    }

    inputs.forEach(input => {
        if (input) {
            input.addEventListener('input', applyFilters);
        }
    });

    if (clearBtn) {
        clearBtn.addEventListener('click', () => {
            inputs.forEach(input => { if (input) input.value = ''; });
            applyFilters();
        });
    }

    // Expose globally so other scripts can trigger it after dynamic loading
    window.triggerTableFilter = applyFilters;
});
