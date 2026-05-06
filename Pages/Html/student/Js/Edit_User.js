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
    const statusSelect = document.getElementById("editStatus");
    const triggerResetBtn = document.getElementById("triggerReset");

    let originalEmail = null;

    // Use event delegation on the table body to catch "Edit" button clicks
    document.querySelector(".tbody").addEventListener("click", (e) => {
        const editButton = e.target.closest("button");
        if (!editButton) return;

        // Ensure it's not the delete button
        if (editButton.classList.contains("openModal")) return;

        const row = editButton.closest("tr");
        if (!row) return;

        // Retrieve user data from data attribute (we'll set this in Admin_Modification.js)
        const user = JSON.parse(row.dataset.user || "{}");
        
        // Populate fields
        firstNameInput.value = user.first || "";
        lastNameInput.value = user.last || "";
        emailInput.value = user.email || "";
        phoneInput.value = user.phone || "";
        statusSelect.value = user.status || "active";
        originalEmail = user.email;

        editModal.classList.add("active");
    });

    // Handle Password Reset Trigger
    triggerResetBtn.addEventListener("click", async () => {
        if (!originalEmail) return;
        
        if (!confirm(`Are you sure you want to trigger a password reset for ${originalEmail}?`)) return;

        triggerResetBtn.disabled = true;
        triggerResetBtn.textContent = "Triggering...";

        try {
            const response = await fetch("/api/admin/trigger-reset", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email: originalEmail })
            });

            if (response.ok) {
                showToast("Password reset link sent to admin/user.", "success");
            } else {
                const text = await response.text();
                showToast(text || "Failed to trigger reset.", "error");
            }
        } catch (err) {
            showToast("Connection error.", "error");
        } finally {
            triggerResetBtn.disabled = false;
            triggerResetBtn.textContent = "Trigger Password Reset";
        }
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
            phone: phoneInput.value,
            status: statusSelect.value
        };

        const saveBtn = e.target.querySelector(".edit-btn-save");
        saveBtn.disabled = true;
        saveBtn.textContent = "Saving...";

        fetch("/api/admin/update-user", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(updatedData)
        })
            .then(async response => {
                if (response.ok) {
                    showToast("User updated successfully!", "success");
                    setTimeout(() => location.reload(), 1000);
                } else {
                    const text = await response.text();
                    showToast(text || "Failed to update user.", "error");
                }
            })
            .catch(err => {
                console.error("Error updating user:", err);
                showToast("An error occurred while updating the user.", "error");
            })
            .finally(() => {
                saveBtn.disabled = false;
                saveBtn.textContent = "Save Changes";
            });
    });
});
