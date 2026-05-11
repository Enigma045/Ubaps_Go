/**
 * Shared logic for Application Comments.
 * Requires:
 * - commentModal HTML structure in the page.
 * - window.currentApplicants to be populated.
 *
 * The submit button handler is bound ONCE on DOMContentLoaded,
 * NOT inside initCommentLogic(), so it is never lost on table refresh.
 */

window.activeStudent = null;

// ─── Open modal & populate comments ─────────────────────────────────────────
window.openCommentModal = function(student) {
    console.log("[Comment] openCommentModal called for:", student && student.registration_number);
    window.activeStudent = student;

    const commentModal = document.getElementById("commentModal");
    if (!commentModal) {
        console.error("[Comment] commentModal element NOT found in DOM!");
        return;
    }

    let comments = student.comments || [];
    if (typeof comments === 'string') {
        try { comments = JSON.parse(comments); } catch(e) { comments = []; }
    }
    // pgx returns jsonb as an object/array directly — handle both
    if (!Array.isArray(comments)) {
        try { comments = JSON.parse(JSON.stringify(comments)); } catch(e) { comments = []; }
        if (!Array.isArray(comments)) comments = [];
    }

    const list = document.getElementById("comments-list");
    if (list) {
        list.innerHTML = "";
        if (!comments || comments.length === 0) {
            list.innerHTML = `<p style="text-align: center; color: var(--sr-slate-400); font-size: 0.85rem; margin-top: 20px;">No comments yet.</p>`;
        } else {
            comments.forEach(c => {
                const date = new Date(parseFloat(c.date) * 1000).toLocaleString();
                const div = document.createElement("div");
                div.style.background = "white";
                div.style.padding = "10px";
                div.style.borderRadius = "6px";
                div.style.borderLeft = "3px solid var(--sr-blue)";
                div.style.boxShadow = "0 1px 2px rgba(0,0,0,0.05)";
                div.innerHTML = `
                    <div style="display: flex; justify-content: space-between; margin-bottom: 5px;">
                        <span style="font-weight: 700; font-size: 0.8rem; color: var(--sr-slate-700);">${c.name} (${c.role})</span>
                        <span style="font-size: 0.7rem; color: var(--sr-slate-400);">${date}</span>
                    </div>
                    <p style="margin: 0; font-size: 0.85rem; color: var(--sr-slate-600); line-height: 1.4;">${c.text}</p>
                `;
                list.appendChild(div);
            });
            list.scrollTop = list.scrollHeight;
        }
    }

    // Clear the textarea each time the modal opens
    const textarea = document.getElementById("new-comment-text");
    if (textarea) textarea.value = "";

    commentModal.classList.add("active");
    console.log("[Comment] modal opened successfully.");
};

// ─── Called after each table render (just wires close buttons) ───────────────
window.initCommentLogic = function() {
    const commentModal = document.getElementById("commentModal");
    if (!commentModal) return;

    // Re-bind close buttons (safe to call multiple times)
    document.querySelectorAll("#closeCommentModal, #cancelComment").forEach(btn => {
        btn.onclick = () => commentModal.classList.remove("active");
    });
};

// ─── Submit button — bound ONCE when the page loads ──────────────────────────
document.addEventListener("DOMContentLoaded", function() {
    const submitBtn = document.getElementById("submitCommentBtn");
    if (!submitBtn) {
        console.warn("[Comment] submitCommentBtn not found on DOMContentLoaded — will retry on first click.");
        return;
    }

    submitBtn.addEventListener("click", function() {
        console.log("[Comment] Post button clicked.");
        console.log("[Comment] activeStudent:", window.activeStudent);

        const text = document.getElementById("new-comment-text")?.value.trim();
        console.log("[Comment] comment text:", text);

        if (!text) {
            alert("Please write a comment before submitting.");
            return;
        }
        if (!window.activeStudent) {
            alert("No student selected. Please open a student's comment modal first.");
            return;
        }

        const regNum = window.activeStudent.registration_number;
        console.log("[Comment] Sending POST /api/add-comment for student:", regNum);

        submitBtn.innerHTML = "Posting...";
        submitBtn.disabled = true;

        fetch("/api/add-comment", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ id: regNum, comment: text })
        }).then(async response => {
            if (!response.ok) {
                const errText = await response.text();
                console.error("[Comment] Server error:", errText);
                throw new Error(errText || "Failed to add comment");
            }
            return response.json();
        }).then(data => {
            console.log("[Comment] Success:", data);
            if (typeof showToast === "function") {
                showToast("Comment posted successfully!", "success");
            } else {
                alert("Comment posted successfully!");
            }

            // Clear the textarea
            document.getElementById("new-comment-text").value = "";

            // Optimistically update the in-memory student object
            let comments = window.activeStudent.comments || [];
            if (typeof comments === 'string') {
                try { comments = JSON.parse(comments); } catch(e) { comments = []; }
            }
            if (!Array.isArray(comments)) comments = [];
            comments.push({
                name: "You",
                role: "Staff",
                date: (Date.now() / 1000).toString(),
                text: text
            });
            window.activeStudent.comments = comments;

            // Re-render the modal to show the new comment immediately
            window.openCommentModal(window.activeStudent);

            // Refresh the table in background to get server-persisted data
            if (typeof window.fetchApplicants === "function" && window.currentPage) {
                window.fetchApplicants(window.currentPage);
            }
        }).catch(error => {
            console.error("[Comment] Fetch error:", error);
            if (typeof showToast === "function") {
                showToast("Error: " + error.message, "error");
            } else {
                alert("Error: " + error.message);
            }
        }).finally(() => {
            submitBtn.innerHTML = "Post Comment";
            submitBtn.disabled = false;
        });
    });

    console.log("[Comment] Submit button listener attached successfully.");
});
