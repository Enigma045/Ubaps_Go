/**
 * Filter_Feature.js
 * Advanced Table Filtering Utility (Multi-Value Consolidated Version)
 */

document.addEventListener('DOMContentLoaded', () => {
    const filterSearch = document.getElementById('filter-search');
    const filterAdvanced = document.getElementById('filter-advanced');
    const clearBtn = document.getElementById('clear-filters');

    // Persistent state for advanced filters (Arrays per category)
    window.activeAdvancedFilters = {};

    function applyFilters() {
        const searchValue = filterSearch ? filterSearch.value.toLowerCase() : '';
        const activeFilters = window.activeAdvancedFilters;

        const tables = document.querySelectorAll('.log-table');
        
        tables.forEach(table => {
            const rows = table.querySelectorAll('tbody tr');
            
            rows.forEach(row => {
                const text = row.innerText.toLowerCase();
                const cells = row.cells;
                if (cells.length < 2) return;

                let matchesSearch = text.includes(searchValue);
                let matchesAdvanced = true;

                // Check each active advanced filter category
                for (const [category, values] of Object.entries(activeFilters)) {
                    if (!values || values.length === 0) continue;
                    
                    let matchesCategory = false; // OR logic within category
                    
                    for (const val of values) {
                        const lowercaseVal = val.toLowerCase();
                        let matchesThisVal = false;

                        if (category === 'dept') {
                            matchesThisVal = cells[1].innerText.toUpperCase().includes(val.toUpperCase());
                        } else if (category === 'priority') {
                            let priorityCell = findCellByHeader(table, "Priority", row);
                            matchesThisVal = (priorityCell && priorityCell.innerText.toLowerCase().includes(lowercaseVal));
                        } else if (category === 'parent') {
                            let parentCell = findCellByHeader(table, "Parent Status", row);
                            matchesThisVal = (parentCell && parentCell.innerText.toLowerCase().includes(lowercaseVal));
                        } else if (category === 'employment') {
                            let employmentCell = findCellByHeader(table, "Employment Status", row);
                            matchesThisVal = (employmentCell && employmentCell.innerText.toLowerCase().includes(lowercaseVal));
                        } else if (category === 'gender') {
                            let genderCell = findCellByHeader(table, "Gender", row);
                            matchesThisVal = (genderCell && genderCell.innerText.toLowerCase() === lowercaseVal);
                        } else if (category === 'scheme') {
                            let schemeCell = findCellByHeader(table, "Bursary Scheme", row) || findCellByHeader(table, "Letter Type", row);
                            matchesThisVal = (schemeCell && schemeCell.innerText.toLowerCase().includes(lowercaseVal));
                        }

                        if (matchesThisVal) {
                            matchesCategory = true;
                            break;
                        }
                    }

                    if (!matchesCategory) {
                        matchesAdvanced = false;
                        break; // AND logic across categories
                    }
                }

                if (matchesSearch && matchesAdvanced) {
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

    if (filterSearch) {
        filterSearch.addEventListener('input', applyFilters);
    }
    if (filterAdvanced) {
        filterAdvanced.addEventListener('change', (e) => {
            const val = e.target.value;
            if (val) {
                const [category, filterVal] = val.split(':');
                if (!window.activeAdvancedFilters[category]) {
                    window.activeAdvancedFilters[category] = [];
                }
                // Only add if not already present
                if (!window.activeAdvancedFilters[category].includes(filterVal)) {
                    window.activeAdvancedFilters[category].push(filterVal);
                }
                
                e.target.value = ""; // Reset dropdown
                if (window.onFilterChange) window.onFilterChange();
                applyFilters();
            }
        });
    }

    if (clearBtn) {
        clearBtn.addEventListener('click', () => {
            if (filterSearch) filterSearch.value = '';
            if (filterAdvanced) filterAdvanced.value = '';
            window.activeAdvancedFilters = {};
            applyFilters();
            if (window.renderFilterTags) window.renderFilterTags({});
        });
    }

    // --- Filter Tags Logic ---
    window.renderFilterTags = function(payload, onRemove) {
        const container = document.getElementById('active-filters-container');
        if (!container) return;
        container.innerHTML = "";

        const fields = [
            { key: 'statuses', label: 'Status', value: payload.statuses },
            { key: 'search', label: 'Search', value: payload.search },
            { key: 'department', label: 'Dept', value: window.activeAdvancedFilters.dept },
            { key: 'priority', label: 'Priority', value: window.activeAdvancedFilters.priority },
            { key: 'scheme', label: 'Scheme', value: window.activeAdvancedFilters.scheme },
            { key: 'parent', label: 'Parent', value: window.activeAdvancedFilters.parent },
            { key: 'employment', label: 'Employment', value: window.activeAdvancedFilters.employment },
            { key: 'gender', label: 'Gender', value: window.activeAdvancedFilters.gender }
        ];

        fields.forEach(field => {
            let val = field.value;
            if (val && (Array.isArray(val) ? val.length > 0 : val !== "")) {
                
                const tag = document.createElement('div');
                tag.className = 'filter-tag';
                
                let displayVal = "";
                if (Array.isArray(val)) {
                    // For statuses, we still filter out internal ones
                    if (field.key === 'statuses') {
                        displayVal = val.filter(s => !["considering", "paid"].includes(s)).join(', ');
                    } else {
                        displayVal = val.join(', ');
                    }
                    if (!displayVal) return;
                } else {
                    displayVal = val;
                }

                tag.innerHTML = `
                    <span class="category">${field.label}:</span>
                    <span class="value">${displayVal}</span>
                    <span class="remove-tag" title="Remove">&times;</span>
                `;
                
                tag.querySelector('.remove-tag').onclick = () => {
                    if (field.key === 'statuses') {
                        const checkboxes = document.querySelectorAll('.tabs input[type="checkbox"]');
                        checkboxes.forEach(cb => cb.checked = false);
                    } else if (field.key === 'search') {
                        const el = document.getElementById('filter-search');
                        if (el) el.value = '';
                    } else {
                        const mapping = {
                            'department': 'dept', 'priority': 'priority', 'scheme': 'scheme',
                            'parent': 'parent', 'employment': 'employment', 'gender': 'gender'
                        };
                        delete window.activeAdvancedFilters[mapping[field.key]];
                    }
                    if (onRemove) onRemove(field.key);
                };
                container.appendChild(tag);
            }
        });
    };

    window.triggerTableFilter = applyFilters;
});
