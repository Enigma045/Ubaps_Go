document.addEventListener("DOMContentLoaded", () => {
    const editModal = document.getElementById("editUserModal");
    const editForm = document.getElementById("editUserForm");
    const closeBtn = document.querySelector(".edit-modal-close");
    const cancelBtn = document.getElementById("cancelEdit");

    // Elements to populate
    const firstNameInput = document.getElementById("editFirstName");
    const lastNameInput = document.getElementById("editLastName");
    const emailInput = document.getElementById("editEmail");
    const phoneInput = document.getElementById("editPhone");

    let originalEmail = null;

    // Use event delegation on the table body to catch "Edit" button clicks
    document.querySelector(".tbody").addEventListener("click", (e) => {
        const editButton = e.target.closest("button");
        if (!editButton) return;

        // Ensure it's not the delete button
        if (editButton.classList.contains("openModal")) return;

        const row = editButton.closest("tr");
        if (!row) return;

        const fullName = row.querySelector(".name").innerText;
        const email = row.querySelector(".email").innerText;
        const phone = row.querySelectorAll("td")[3].innerText; // Phone is the 4th column (index 3)

        // Split full name (assuming First Last format)
        const nameParts = fullName.split(" ");
        const first = nameParts[0] || "";
        const last = nameParts.slice(1).join(" ") || "";

        // Populate fields
        firstNameInput.value = first;
        lastNameInput.value = last;
        emailInput.value = email;
        phoneInput.value = phone;
        originalEmail = email;

        editModal.classList.add("active");
    });

    const closeEditModal = () => {
        editModal.classList.remove("active");
    };

    closeBtn.addEventListener("click", closeEditModal);
    cancelBtn.addEventListener("click", closeEditModal);

    // Close on click outside
    editModal.addEventListener("click", (e) => {
        if (e.target === editModal) closeEditModal();
    });

    // Handle Save
    editForm.addEventListener("submit", (e) => {
        e.preventDefault();

        const updatedData = {
            originalEmail: originalEmail,
            first: firstNameInput.value,
            last: lastNameInput.value,
            email: emailInput.value,
            phone: phoneInput.value
        };

        console.log("Saving user data:", updatedData);

        // API Call (Endpoint assumed based on common patterns, though user said don't touch Go)
        fetch("/updateuser", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(updatedData)
        })
            .then(response => {
                if (response.ok) {
                    location.reload();
                } else {
                    alert("Failed to update user.");
                }
            })
            .catch(err => {
                console.error("Error updating user:", err);
                alert("An error occurred.");
            });
    });
});
