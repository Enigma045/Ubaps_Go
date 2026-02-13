document.addEventListener('DOMContentLoaded', () => {
    const reviewModal = document.getElementById('reviewModalOverlay');
    const closeBtn = document.getElementById('closeReviewCard');
    const tbody = document.getElementById('judge_tbody');

    // Data placeholders in modal
    const fields = {
        name: document.getElementById('review-name'),
        id: document.getElementById('review-id'),
        course: document.getElementById('review-course'),
        parents: document.getElementById('review-parents'),
        employment: document.getElementById('review-employment'),
        income: document.getElementById('review-income'),
        priority: document.getElementById('review-priority'),
        residence: document.getElementById('review-residence'), // We'll need to infer or add this
        guardian: document.getElementById('review-guardian'),
        select: document.getElementById('review-scheme-select'),
        amount: document.getElementById('review-amount-input'),
    };

    // Open logic
    tbody.addEventListener('click', async (e) => {
        const reviewBtn = e.target.closest('.review-btn');
        if (!reviewBtn) return;

        const row = reviewBtn.closest('tr');
        const cells = row.cells;

        // Extract data (based on the table structure in Committe_Review.html)
        // 0: Name, 1: Reg #, 2: Parent Status, 3: Employment, 4: Income, 5: Priority
        const data = {
            name: cells[0].textContent.trim(),
            id: cells[1].textContent.trim(),
            parents: cells[2].textContent.trim(),
            employment: cells[3].textContent.trim(),
            income: cells[4].textContent.trim(),
            priority: cells[5].textContent.trim()
        };

        // Populate modal
        fields.name.textContent = data.name;
        fields.id.textContent = data.id;
        fields.parents.textContent = data.parents;
        fields.employment.textContent = data.employment;
        fields.income.textContent = data.income;
        fields.priority.textContent = data.priority;

        // These might need mapping or fallback
        fields.course.textContent = "Student @ UNILIA";
        fields.residence.textContent = "N/A";
        fields.guardian.textContent = data.parents.includes('Father') ? 'Father' : 'N/A';
       // document.createElement('option').value = "Loans";
       //temp.textContent = "Loans"
       //
       let road = await getScheme();
        console.log(road);

        road.forEach(item => {
            const option = document.createElement('option');
            option.value = item.scheme_name;
            option.textContent = item.scheme_name;
            fields.select.appendChild(option);
        });

        reviewModal.classList.add('active');
        
    });

    // Close logic
    closeBtn.addEventListener('click', () => {
        reviewModal.classList.remove('active');
    });

    window.addEventListener('click', (e) => {
        if (e.target === reviewModal) {
            reviewModal.classList.remove('active');
        }
    });

    // Finalize Assignment Logic (Placeholder)
    document.getElementById('review-assign-btn').addEventListener('click', () => {
        const scheme = document.getElementById('review-scheme-select').value;
        const amount = document.getElementById('review-amount-input').value;
        
        const studentId = fields.id.textContent;

        if (scheme === 'none' || !amount) {
            alert('Please select a scheme and enter an amount.');
            return;
        }

        console.log(`Finalizing assignment for ${studentId}: ${scheme} - MWK ${amount}`);
        // Here you would typically perform a fetch to save the data
        fetch("/schemeinfo",{
            method:"POST",
            headers:{
                "Content-Type":"application/json"
            },
            body:JSON.stringify({
                reg : studentId,
                amount: amount,
                scheme: scheme
            })
        }).then(res => {
            if (!res.ok) {
                throw new Error("Failed to send scheme info");
            }
            return res.json();
        }).then(data => {
            console.log(data);
            alert(`Bursary assigned successfully to ${fields.name.textContent}`);
            reviewModal.classList.remove('active');
        }).catch(err => {
            console.error(err);
            alert("Failed to send scheme info");
        });
        
    });
});

async function getScheme() {
  try {
    const res = await fetch("/getschemes");

    if (!res.ok) {
      throw new Error("Failed to fetch schemes");
    }

    const data = await res.json();
    return data;

  } catch (err) {
    console.error(err);
    throw err;
  }
}
