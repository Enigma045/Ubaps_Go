# Financial Office Operations Analysis

The Financial Office within the UBAPS ecosystem serves as a critical link between academic selection and fiscal disbursement. Their operations are centered around verifying student financial status, managing bursary payments, and maintaining audit integrity.

## 0. Portal Navigation (Routes)
The Finance Office has access to several specialized views:
*   **`/financialdashboard`**: Main landing page with high-level financial KPIs.
*   **`/financial_request`**: Interface for viewing and managing incoming student statement requests.
*   **`/financial_review`**: The "Judge" interface for approving or rejecting bursary applications from a fiscal perspective.
*   **`/financial`**: The "Financial Request" portal for manually updating student fees and dossier history.
*   **`/financial_reports`**: Analytical interface for generating and exporting fiscal distribution reports (Disbursement, Utilisation, History, and Pending).

## 1. Application Review & Approval
Finance Officers are part of the multi-stage approval workflow for bursary applications.
*   **Approval Status**: Every application has a specific `finance_office_approval_status` field.
*   **Verification**: They review the student's financial background (parental status, guardian employment, and relative support) to ensure the bursary is appropriately targeted.

## 2. Financial Dossier & Fees Management
The system maintains a comprehensive `financial_history` for each student.
*   **Statement Requests**: Students can request their latest fees statement through the portal.
*   **Dossier Updates**: Finance Officers respond to these requests by updating the `financial_history` table with:
    *   Semester information.
    *   Payment dates and exact amounts.
    *   Detailed transaction descriptions.
*   **Status Tracking**: Requests are tracked with statuses like `sent` (requested) and `answered` (provided).

## 3. Bursary Payment Processing
Once a student is selected and approved by all deans, the Finance Office processes the actual disbursement.
*   **Installment Payments**: The office triggers `PayInstallment` operations which update the student's balance in the database.
*   **Automated Notifications**: Upon processing a payment, the system automatically sends:
    *   **In-app Notifications**: Alerting the student of the update.
    *   **SMS Alerts**: Providing immediate confirmation via mobile.
    *   **Email Alerts**: Sending detailed payment summaries to administrators for monitoring.

## 4. Dashboard & Analytics
Finance Officers have a dedicated dashboard (`/financialdashboard`) powered by the `/stats/finance` endpoint.
*   **Disbursement Rates**: Tracking what percentage of allocated bursaries have been paid.
*   **Financial Statistics**: Viewing total applicants vs. approved financial aid amounts.
*   **Departmental Breakdown**: Analyzing financial aid distribution across different faculties (CEN, BPH, FSN, EDU).

## 5. Audit Logging & Security
All financial operations are strictly logged for accountability.
*   **Action Types**: Every success or failure is recorded in the `user_logs` and `payment_logs` tables.
*   **Audit Trail**: Logs include the performer ID, action type (e.g., `BURSARY_PAYMENT_PROCESSED`, `DOSSIER_RESPONSE_SUCCESS`), and a timestamp.
*   **RBAC**: Access to financial data and payment buttons is restricted to users with the `finance_office` role.
