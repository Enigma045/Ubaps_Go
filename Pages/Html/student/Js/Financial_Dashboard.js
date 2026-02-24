/* Financial Dashboard Logic */

document.addEventListener('DOMContentLoaded', async () => {
    console.log("Financial Dashboard Initialized");


    // Load Stats
    if (typeof fetchFinancialOfficerStats === 'function') {
        const stats = await fetchFinancialOfficerStats();
        if (stats) {
            const approvedAmountEl = document.getElementById('approved-amount');
            const numDisbursementsEl = document.getElementById('num-disbursements');
            const numRequestsEl = document.getElementById('num-requests');

            if (approvedAmountEl) approvedAmountEl.textContent = `MWK ${stats.approved_amount.toLocaleString()}`;
            if (numDisbursementsEl) numDisbursementsEl.textContent = stats.disbursements_made;
            if (numRequestsEl) numRequestsEl.textContent = stats.financial_history_requests;
        }
    }
});
