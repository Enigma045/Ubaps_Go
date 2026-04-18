/* Registrar Dashboard Logic */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("Registrar Dashboard Initialized");


    // Load Stats
    if (typeof fetchRegistrarStats === 'function') {
        const stats = await fetchRegistrarStats();
        if (stats) {
            const approvedAmountEl = document.getElementById('approved-amount');
            const numApplicantsEl = document.getElementById('num-applicants');
            const numSchemesEl = document.getElementById('num-schemes');

            if (approvedAmountEl) approvedAmountEl.textContent = `MWK ${stats.approved_amount.toLocaleString()}`;
            if (numApplicantsEl) numApplicantsEl.textContent =  stats.pending_applications + stats.considering_applications;
            if (numSchemesEl) numSchemesEl.textContent = stats.number_of_schemes;
        }
    }
});
