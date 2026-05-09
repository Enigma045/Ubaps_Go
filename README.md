# UBAPS - University Bursary Application and Processing System

UBAPS is a robust, web-based management system designed to streamline the entire lifecycle of bursary applications within a university ecosystem. It provides a centralized platform for students to apply for financial aid and for university staff (Registrars, Finance Officers, Deans, and Admins) to review, process, and manage these applications efficiently.

## 🚀 Key Features

### 🎓 For Students
- **Digital Application**: Submit bursary applications with all required documentation online.
- **Application Tracking**: Real-time status updates on application progress.
- **Financial History**: View past disbursements, installments, and payment logs.
- **Notifications**: Stay updated with system alerts regarding application status and requirements.
- **Profile Management**: Maintain personal and academic information securely.

### 📜 For Registrars
- **Application Review**: Comprehensive tools to evaluate student eligibility and documentation.
- **Scheme Management**: Create and manage different bursary schemes and criteria.
- **Letter Processing**: Generate and send official communication/letters to students.
- **Advanced Statistics**: Visualized data on applicant demographics and trends.

### 💰 For Finance Office
- **Payment Processing**: Manage installments and disbursement requests.
- **Financial Reviews**: Audit financial eligibility and history of applicants.
- **Total Amount Tracking**: Monitor bursary fund utilization and remaining balances.
- **Comprehensive Reporting**: Export financial data for auditing and planning.

### 🏛️ For Deans & Committee
- **Decision Workflow**: Review registrar recommendations and make final approvals/rejections.
- **Departmental Reports**: View bursary statistics specific to their faculty or school.

### 🛠️ For Administrators
- **User Management**: Create, update, and manage accounts for all system roles.
- **Audit Trails**: Detailed logs of user actions, payments, and application changes for security and accountability.
- **System Monitoring**: Access to system-wide logs and performance metrics.

## 🛠️ Tech Stack

- **Backend**: [Go (Golang)](https://golang.org/) 1.25.5+
- **Database**: [PostgreSQL](https://www.postgresql.org/) (via `pgx/v5`)
- **Frontend**: Vanilla HTML5, CSS3 (Custom Design), and Modern JavaScript
- **Security**: 
    - Custom Role-Based Access Control (RBAC) middleware
    - Password hashing with `bcrypt`
    - UUID-based session/record tracking
- **Integrations**: SMS bridge via PHP integration for notifications.

## 📂 Project Structure

```text
Ubaps/
├── main.go             # Entry point - Server configuration and routing
├── Db/                 # Database connection and schema management
├── Routes/             # API and Page route handlers
├── Handles/            # Core business logic and data processing handlers
├── Middleware/         # Authentication and Role validation logic
├── Models/             # Database entity structs and data models
├── Pages/              # Frontend assets (HTML, CSS, JS, Images)
├── Notifications/      # System-wide notification logic
├── services/           # External service integrations
├── sms/                # SMS gateway integration scripts
├── utils/              # Shared helper functions and constants
└── Audit_logs/         # System audit trail storage
```

## ⚙️ Setup & Installation

### Prerequisites
- Go 1.25.5 or higher
- PostgreSQL 15+
- A modern web browser

### Steps
1. **Clone the repository**:
   ```bash
   git clone https://github.com/Enigma045/Ubaps_Go.git
   cd Ubaps_Go
   ```

2. **Database Setup**:
   - Create a PostgreSQL database named `Ubaps`.
   - The application currently expects a connection string in `Db/DB.go`. Ensure your PostgreSQL user `postgres` has the password `characte2002` or update the `dsn` variable in `Db/DB.go` to match your local setup.
   - *Note: It is recommended to use environment variables for sensitive credentials in production.*

3. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the Application**:
   ```bash
   go run main.go
   ```
   The server will start at `http://localhost:8080`.

## 🔒 Security

UBAPS implements a strict security model:
- **Authentication**: All sensitive routes are protected by a session-based auth middleware.
- **RBAC**: Access to specific dashboards (Admin, Registrar, Finance) is strictly enforced based on the user's assigned role.
- **Data Integrity**: Uses PostgreSQL transactions for critical operations like payment processing and application status changes.

## 📄 License

This project is currently for internal use.

