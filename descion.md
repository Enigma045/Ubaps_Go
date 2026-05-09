# Analysis of Bursary Decision Pages

This document provides a detailed analysis of the three primary review and decision-making interfaces within the UBAPS (Unilia Bursary Award Processing System).

## 1. Overview of Pages

| Page | Primary User Role | Navigation Focus |
| :--- | :--- | :--- |
| `Deans_Decisions.html` | Dean of Students / Faculty | Review and intermediate approval |
| `Committe_Review.html` | Registrar / Committee Head | Final review and Bursary Assignment |
| `Financial_review.html` | Finance Office | Financial verification and audit |

---

## 2. Common Architectural Elements

All three pages share a unified design system and core logic, ensuring a consistent experience for different administrative roles.

### A. Shared UI Components
- **Sidebar & Header**: Standard UNILIA sidebar with role-based navigation links and a header featuring breadcrumbs, notification badges, and user avatars.
- **Tabbed Interface**: Three primary views for data management:
    - **Applicants**: Active submissions (`submitted` / `considering`).
    - **Selected**: Students marked for bursaries (`selected`).
    - **Not Selected**: Rejected applications (`not selected`).
- **Advanced Filter Bar**: Identical filtering capabilities across all roles:
    - Search (Name/Reg #)
    - Department (CEN, BPH, FSN, EDU)
    - Priority (High, Medium, Low)
    - Parent Status (Deceased/Alive)
    - Employment (Unemployed, Formal, Informal, etc.)
- **Dynamic Table**: All pages use the `#judge_tbody` for AJAX-powered row injection.

### B. Shared JavaScript Logic
- `Committe_Review.js`: Handles tab switching, pagination, and fetching data from `/applicants`.
- `Filter_Feature.js`: Provides real-time client-side filtering logic based on table headers.
- `card.js` & `Dossier_Feature.js`: Powers the complex "Review Card" and "Financial Dossier" modals.

---

## 3. Functional Analysis & Role Specifics

### `Deans_Decisions.html`
- **Purpose**: Acts as the first line of academic/situational vetting.
- **Key Logic**: Deans "Approve" students, which marks them as `considering` in the database. This doesn't award money but signals to the committee that the student is a valid candidate.
- **Review Card**: Uses the `no-sidebar` layout. Deans can see the "Approval Status Tracker" to see if other deans or the registrar have weighed in, but they **cannot** assign money.

### `Committe_Review.html`
- **Purpose**: The "Command Center" for the Registrar.
- **Key Logic**: This is the most powerful page. It is the **only** page that includes the `sidebar-action` inside the Review Card.
- **Bursary Assignment**: Once all required approvals (Registrar, Deans, Finance) are marked as "Approved" in the tracker, the **Assign Bursary** sidebar appears.
- **Fields**: Features a Scheme Selector and Amount Input (MWK). Clicking "Finalize Assignment" triggers the `/schemeinfo` API to officially award the bursary.

### `Financial_review.html`
- **Purpose**: Final financial sanity check.
- **Key Logic**: Finance officers review the "Financial Dossier" (payment history) to ensure the student actually needs support or has a valid debt.
- **Review Card**: Uses the `no-sidebar` layout. Finance officers focus on "Verification Documents" (Requesting Statements/Viewing Dossiers).

---

## 4. Modal System Analysis

The system uses a multi-layered modal approach to handle complex data without leaving the dashboard.

1.  **Approve/Reject Modals**: Simple confirmation dialogs for quick status updates.
2.  **Request Statement Modal**: Triggers a notification to the student to upload their latest bank/tuition statements.
3.  **Review Card (`modal-overlay`)**: A comprehensive 360-degree view of the student, including:
    - Academic Info (Course/Reg #).
    - Social Info (Parents/Employment/Support).
    - **Approval Tracker**: Real-time status of the 5 required approvals.
4.  **Dossier Modal**: A deep-dive table showing every semester's transaction history for the student, fetched via `/api/get-financial-history`.

---

## 5. Technical Commonalities

- **CSS Layers**: All pages load `Navbar.css`, `Layout.css`, `Dashboard.css`, `Notification.css`, `Bursary_Scheme.css`, `card.css`, and `Filter_Bar.css`.
- **API Strategy**: All pages interact with `/applicants` (POST) to populate the main table and `/getapplicationstatus` for the tracker.
- **MIME/Static Handling**: All assets are served through the central Go `FileServer` defined in `main.go`.
