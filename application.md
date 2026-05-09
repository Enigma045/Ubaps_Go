# Student Bursary Application Form - Required Information

This document outlines the exact data points a student must provide when completing the Bursary Application form in UBAPS.

## Step 1: Personal Details
The first section collects basic demographic and residential information.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| **Date of Birth** | Date | Yes | The student's legal date of birth. |
| **Gender** | Radio | Yes | Options: Male or Female. |
| **Home District** | Text | Yes | The student's primary place of residence/origin (e.g., Lilongwe). |
| **Accommodation** | Radio | Yes | Options: On Campus or Off Campus. |

> [!NOTE]
> **Pre-filled Data:** Full Name, Phone Number, and Email are automatically pulled from the student's registration profile and cannot be edited within this form.

---

## Step 2: Financial Information
This section assesses the student's financial need based on family status and external support.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| **Parent/Guardian Status** | Select | Yes | Options: Both Alive, One Alive, or Both Deceased. |
| **Guardian Employment** | Select | Yes | Options: Employed (Formal), Employed (Informal), Self-Employed, or Unemployed. |
| **Other Financial Support**| Select | Yes | Indicates if the student receives aid from other bursaries or sponsors (Yes/No). |

---

## Step 3: Review & Statement
The final step allows the student to provide context for their application.

| Field | Type | Required | Description |
| :--- | :--- | :--- | :--- |
| **Reason for Application** | Textarea| No* | A detailed explanation of why the student requires financial assistance. |

*\*Note: While not technically marked with the `required` attribute in the current HTML, this statement is critical for the committee's decision-making process.*

---

## Technical Submission Details
- **Endpoint**: `/SubmitForm`
- **Method**: `POST`
- **Data Format**: `multipart/form-data`
- **Validation**: Performed client-side before submission to ensure no required fields are left blank.
