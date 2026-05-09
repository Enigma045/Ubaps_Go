function openDossier(studentId, studentName) {
    if (!studentId || studentId === "ID Number") {
        showToast("No student selected.", "error");
        return;
    }

    console.log("Fetching dossier for:", studentId);
    showToast("Fetching student dossier...", "info");

    fetch("/api/get-financial-history", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ reg: studentId })
    })
    .then(res => {
        if (!res.ok) throw new Error("Failed to fetch financial history");
        return res.json();
    })
    .then(data => {
        const dossierModal = document.getElementById('dossierModal');
        const dossierTbody = document.getElementById('dossier-history-tbody');
        const dossierName = document.getElementById('dossier-student-name');
        const dossierReg = document.getElementById('dossier-student-reg');

        if (dossierName) dossierName.textContent = studentName || "Student";
        if (dossierReg) dossierReg.textContent = studentId;

        if (dossierTbody) {
            dossierTbody.innerHTML = "";
            if (!data || data.length === 0) {
                dossierTbody.innerHTML = '<tr><td colspan="4" style="text-align:center; padding: 20px;">No financial history found for this student.</td></tr>';
            } else {
                data.forEach(item => {
                    const date = new Date(item.date).toLocaleDateString();
                    const tr = document.createElement('tr');
                    tr.innerHTML = `
                        <td style="padding: 10px; border-bottom: 1px solid #f1f5f9;">${item.semester}</td>
                        <td style="padding: 10px; border-bottom: 1px solid #f1f5f9;">${date}</td>
                        <td style="padding: 10px; border-bottom: 1px solid #f1f5f9;">${item.details}</td>
                        <td style="padding: 10px; border-bottom: 1px solid #f1f5f9; text-align: right;">${item.amount.toLocaleString()}</td>
                    `;
                    dossierTbody.appendChild(tr);
                });
            }
        }

        if (dossierModal) dossierModal.classList.add('active');
    })
    .catch(err => {
        console.error("Dossier fetch error:", err);
        showToast("Error loading dossier: " + err.message, "error");
    });
}

document.addEventListener('DOMContentLoaded', () => {
    // Dossier Modal Close Logic
    const closeDossierBtn = document.getElementById('closeDossier');
    const closeDossierFooter = document.getElementById('closeDossierFooter');
    const dossierModal = document.getElementById('dossierModal');

    const closeDossierModals = () => {
        if (dossierModal) dossierModal.classList.remove('active');
    };

    if (closeDossierBtn) closeDossierBtn.addEventListener('click', closeDossierModals);
    if (closeDossierFooter) closeDossierFooter.addEventListener('click', closeDossierModals);
    
    // Close on outside click
    window.addEventListener('click', (e) => {
        if (e.target === dossierModal) {
            closeDossierModals();
        }
    });
});
