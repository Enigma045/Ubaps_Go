/**
 * Updates the bursary status badge and amount on the student review card.
 * Fetches accurate data from the backend using the registration number.
 * @param {string} regNumber - The registration number of the student.
 */
function updateBursaryCardInfo(regNumber) {
    if (!regNumber) return;

    // Elements in the review card
    const statusElem = document.getElementById('review-status');
    const amountElem = document.getElementById('review-amount');

    fetch("/getapplicationstatus", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ reg: regNumber })
    })
        .then(res => {
            if (!res.ok) throw new Error("Failed to fetch accurate bursary info");
            return res.json();
        })
        .then(data => {
            console.log("Accurate bursary data retrieved:", data);

            if (statusElem) {
                // Determine display text: Use Scheme Name if selected, otherwise overall Status
                let displayText = "Pending";
                const statusStr = (data.status || "not submitted").toLowerCase().trim();

                if (statusStr === "selected" && data.scheme_name && data.scheme_name.Valid) {
                    displayText = data.scheme_name.String;
                } else {
                    displayText = statusStr.charAt(0).toUpperCase() + statusStr.slice(1);
                }

                statusElem.textContent = displayText;
                statusElem.className = "status-badge";

                if (statusStr === "selected" || statusStr === "approved") {
                    statusElem.classList.add("approved");
                } else if (statusStr === "not selected" || statusStr === "rejected") {
                    statusElem.classList.add("rejected");
                } else {
                    statusElem.classList.add("pending");
                }
            }

            if (amountElem) {
                const amountStr = String(data.bursary_amount || "0").trim();
                const amount = parseFloat(amountStr) || 0;
                amountElem.textContent = `MWK ${amount.toLocaleString(undefined, { minimumFractionDigits: 2 })}`;
            }
        })
        .catch(err => {
            console.error("Error updating bursary card info:", err);
        });
}
