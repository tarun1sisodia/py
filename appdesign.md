
---

```markdown
# Smart Attendance App Specification



- [Introduction](#introduction)
- [Technology Stack](#technology-stack)
- [Application Architecture](#application-architecture)
- [User Roles and Flows](#user-roles-and-flows)
  - [Common Authentication & Registration](#common-authentication--registration)
  - [Teacher Workflow](#teacher-workflow)
  - [Student Workflow](#student-workflow)
- [Core Features](#core-features)
- [Backend & AI Processing](#backend--ai-processing)
- [Database Design and Schema](#database-design-and-schema)
- [Error Handling and Security Measures](#error-handling-and-security-measures)
  - [Enhanced API Security](#enhanced-api-security)
- [Implementation Details](#implementation-details)
- [Step-by-Step Implementation Plan](#step-by-step-implementation-plan)
  - [Phase 11: API Testing with Postman Locally](#phase-11-api-testing-with-postman-locally)
- [Project Structure](#project-structure)
- [Conclusion](#conclusion)

---

## Introduction

The Smart Attendance App is designed to replace traditional paper-based attendance with a digital, secure, and automated solution for educational institutions. The app offers distinct interfaces for teachers and students, featuring robust verification methods (location, Wi-Fi, device binding, and developer option checks) and AI-powered anomaly detection using Gemini AI. Only students whose registered academic year matches the academic year of the teacher's active session are allowed to mark attendance.

> **Updated:** Authentication now leverages Firebase to handle phone verification via SMS and password resets via email. This ensures robust, secure, and scalable user verification.

---

## Technology Stack

- **Frontend: Flutter**
  - State management using the BLoC pattern.
  - Animated onboarding screens and secure device binding.
- **Backend: Golang**
  - Developed with the Gin framework for optimal API development.
  - Secure session management and business logic implementation.
  - **Firebase Admin SDK:** Used for verifying Firebase ID tokens and linking Firebase users to MySQL records.
- **Database: MySQL**
  - Optimized schema to manage users, sessions, and attendance records.
- **AI Processing: Gemini AI**
  - Provides anomaly detection and secures attendance verification.
- **Additional Security Features:**
  - OTP-based verification for registration via SMS.
  - Device binding mechanism.
  - Developer option detection.
  - Real-time verification (Wi-Fi, GPS, Device ID).
  - **SMS Code for Phone Number Verification & Email for Password Reset**

---

## Application Architecture

- **Client Side (Flutter):**
  - **User Interfaces:** Separate dashboards and flows for teachers and students.
  - **Authentication:** Role-based Sign In/Sign Up for teachers and students.
    - **Firebase Integration:** Uses Firebase Authentication to send SMS codes for OTP verification and secure email links for password resets.
  - **Real-Time Updates:** Dynamic session timers, status indicators, and an animated onboarding process.
  - **Device Integration:** Accesses GPS and Wi-Fi APIs for verification.
  
- **Server Side (Golang):**
  - **RESTful API:** Handles authentication, session creation, attendance logging, OTP verification, and password resets.
  - **Business Logic:** Validates user credentials, academic year, location, Wi-Fi connectivity, and device binding.
  - **Persistent Storage:** Interfaces with MySQL to store user profiles, session data, and attendance records.
  - **Firebase Admin Integration:** Verifies Firebase ID tokens to securely authenticate users and create or update records based on their persistent Firebase UID.
  - **AI Integration:** Uses Gemini AI to detect proxy attendance and analyze attendance patterns.

---

## User Roles and Flows

### Common Authentication & Registration

- **Welcome Screen:**
  - Upon launch, users select their role as either **Teacher** or **Student** via clearly marked icons or buttons.
- **Registration (Sign Up):**
  - **Teacher Registration:**
    - **Input Fields:** Full Name, Username (primary key), Email, Phone Number, Highest Degree, Password, Confirm Password, Experience.
    - After form submission, an OTP SMS code is sent to the provided phone number to verify the phone. In addition, a verification email may be sent for further confirmation. Once verified via Firebase, the teacher is redirected to the login page.
  - **Student Registration:**
    - **Input Fields:** Full Name, Roll Number (primary key), Course, Academic Year, Phone Number, Password, Confirm Password.
    - An OTP SMS code is sent for phone number verification. After entering and validating the SMS code (via Firebase), the student is redirected to the login page.
- **Login:**
  - **Teacher Login:** Teachers log in using their **Email** and **Password**.
  - **Student Login:** Students log in using their **Roll Number** and **Password**.
- **Password Reset:**
  - When a user requests a password reset, a secure password reset link is sent to their registered email address (handled by Firebase). The user follows the link to reset their password and is then redirected to the login page.

---

### Teacher Workflow

1. **Teacher Dashboard:**
   - After logging in, the teacher sees a personalized dashboard displaying their name.
   - The dashboard includes dropdown menus for selecting the **Academic Year** (e.g., BCA 1st Year, 2nd Year, 3rd Year, etc.) and **Subject**.
   
2. **Session Initiation:**
   - A prominent **"Start Attendance"** button is available.
   - Prior to starting a session, the teacher selects a countdown timer option (30 seconds, 1 minute, or 3 minutes).
   - When the button is tapped, an API call creates a new attendance session. The session details include:
     - Teacher ID.
     - Selected Academic Year and Subject.
     - Countdown timer duration.
   - The session is marked as "active" and the countdown timer begins.

3. **Session Monitoring:**
   - The dashboard updates in real time to show the active session status.
   - Teachers can view live attendance records as students mark their attendance.
   - Additional controls may allow ending the session early or reviewing detailed statistics.

---

### Student Workflow

1. **Student Dashboard:**
   - Upon login, the student sees a personalized home screen displaying their name and a **"Mark Attendance"** button.
   - The button remains disabled until an active session is detected.
   
2. **Academic Year Matching:**
   - During registration, the student provides their Academic Year.
   - When a session is active, the backend verifies that the student’s Academic Year matches the teacher’s session Academic Year. Only matching students are allowed to mark attendance.

3. **Session Detection and Onboarding:**
   - The app continuously polls (or uses a WebSocket) for an active session.
   - Once a valid session is detected and the academic year check passes, an animated onboarding process is initiated. This process displays:
     - Wi-Fi connection verification (ensuring connection to the institution’s network).
     - GPS location verification (confirming the student is on campus).
     - Device binding confirmation (ensuring the device is registered).
     - Developer option check (ensuring the developer mode is disabled).
   
4. **Attendance Marking Process:**
   - After successful onboarding and verification, the **"Mark Attendance"** button is activated.
   - When tapped, an API call is sent with:
     - Student ID.
     - Active Session ID.
     - Verification data (GPS coordinates, Wi-Fi SSID/BSSID, device information).
   - The backend validates that:
     - The student is within campus bounds.
     - The Wi-Fi network matches the institution’s credentials.
     - The device is bound and the developer option is disabled.
     - The student has not already marked attendance for the subject on that day.
   - On successful validation, the attendance record is created, and a success notification is displayed.

---

## Core Features

- **Role-Based Dashboards:** Tailored interfaces for teachers and students.
- **Session-Based Attendance:** Attendance is recorded only during active sessions.
- **Robust Verification:** Multiple layers of verification (location, Wi-Fi, device binding, developer option check).
- **Academic Year Matching:** Ensures that only students with the same Academic Year as the session can mark attendance.
- **OTP Verification via SMS:** Phone numbers are verified using an SMS code (managed by Firebase).
- **Password Reset via Email:** Secure password reset links are sent to users’ email addresses (handled by Firebase).
- **Scalable Architecture:** Built using Flutter, Golang, and MySQL.
- **AI-Powered Security:** Gemini AI detects and prevents proxy attendance.

---

## Backend & AI Processing

- **Authentication & User Management:**
  - API endpoints for teacher and student registration, OTP (SMS) verification, login, and password reset (email-based).
  - JWT tokens are used for session management.
  - **Firebase Integration:** The backend uses the Firebase Admin SDK to verify Firebase ID tokens. The unique Firebase UID is stored in the MySQL database to link Firebase authentication with custom user data.
  
- **Session Management:**
  - API for creating and managing attendance sessions.
  - Session data includes teacher ID, academic year, subject, start/end times, and timer duration.
  
- **Attendance Recording:**
  - API to record student attendance after verification.
  - Enforces one attendance mark per subject per day.
  
- **AI Integration with Gemini AI:**
  - Analyzes attendance data to detect anomalies and potential proxy attendance.
  - Provides detailed reports and insights.
  
- **Logging and Security:**
  - Comprehensive logging for all API calls and error handling.
  - Secure data transmission via HTTPS and encrypted storage.

---

## Database Design and Schema

### Users Table
- **Teacher Users:**  
  - Fields: User ID, Full Name, Username (primary key), Email, Phone Number, Highest Degree, Password Hash, Experience, Role, Created/Updated Timestamps.
- **Student Users:**  
  - Fields: User ID, Full Name, Roll Number (primary key), Course, Academic Year, Phone Number, Password Hash, Role, Created/Updated Timestamps.

### Sessions Table
- Stores session records with:
  - Session ID, Teacher ID, Subject ID, Academic Year, Start Time, End Time, Countdown Duration, Status, Wi-Fi Credentials, Location Coordinates, Created Timestamp.

### Attendance Records Table
- Logs attendance with:
  - Record ID, Session ID, Student ID, Marked Timestamp, Verification Methods, Device Info, Location Data, Wi-Fi Data, Created Timestamp.

### Additional Tables
- **OTP Verification Table:**  
  - Stores temporary OTPs and validation statuses for registration and password reset.
- **Audit Logs Table:**  
  - Logs critical actions and security events.

#### Example Database Schema (SQL)

```sql
-- Users Table (Teachers and Students)
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    role ENUM('teacher', 'student') NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    username VARCHAR(100) UNIQUE,       -- For teachers
    roll_number VARCHAR(100) UNIQUE,      -- For students
    email VARCHAR(255),                   -- Used for teacher login and password reset
    course VARCHAR(100),                  -- For student registration
    academic_year VARCHAR(50),            -- For student registration and matching with session
    phone VARCHAR(20) NOT NULL,
    highest_degree VARCHAR(100),          -- For teachers
    experience VARCHAR(50),               -- For teachers
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Sessions Table
CREATE TABLE attendance_sessions (
    id VARCHAR(36) PRIMARY KEY,
    teacher_id VARCHAR(36) NOT NULL,
    subject_id VARCHAR(36) NOT NULL,
    academic_year VARCHAR(50) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    countdown_duration ENUM('30s', '1m', '3m') NOT NULL,
    status ENUM('active', 'completed', 'cancelled') DEFAULT 'active',
    wifi_ssid VARCHAR(100) NOT NULL,
    wifi_bssid VARCHAR(100) NOT NULL,
    location_lat DECIMAL(10, 8) NOT NULL,
    location_long DECIMAL(11, 8) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (teacher_id) REFERENCES users(id)
);

-- Attendance Records Table
CREATE TABLE attendance_records (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL,
    student_id VARCHAR(36) NOT NULL,
    marked_at TIMESTAMP NOT NULL,
    verification_method ENUM('location', 'wifi', 'both') NOT NULL,
    device_info JSON,
    location_lat DECIMAL(10, 8),
    location_long DECIMAL(11, 8),
    wifi_ssid VARCHAR(100),
    wifi_bssid VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES attendance_sessions(id),
    FOREIGN KEY (student_id) REFERENCES users(id)
);

-- OTP Verification Table
CREATE TABLE otp_verifications (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    otp_code VARCHAR(10) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Audit Logs Table (for general system audits)
CREATE TABLE audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    action VARCHAR(100) NOT NULL,
    details JSON,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- ================================
-- Additional Admin Database Schema
-- ================================

-- Admins Table: Stores admin credentials and roles
CREATE TABLE admins (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('super_admin', 'moderator') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Admin Audit Logs Table: Logs actions performed by admin users for accountability
CREATE TABLE admin_audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    admin_id VARCHAR(36) NOT NULL,
    action VARCHAR(100) NOT NULL,
    details JSON,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES admins(id)
);

-- Reports Table: Stores reports generated by admin users (e.g., attendance, user activity, security)
CREATE TABLE reports (
    id VARCHAR(36) PRIMARY KEY,
    report_type ENUM('attendance', 'user_activity', 'security') NOT NULL,
    generated_by VARCHAR(36) NOT NULL,  -- References the admin who generated the report
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data JSON,
    FOREIGN KEY (generated_by) REFERENCES admins(id)
);
```

---

## Error Handling and Security Measures

- **Centralized Error Management:**
  - Global error handling middleware in both Flutter (BLoC) and Golang (Gin) to catch, log, and return meaningful error messages.
- **Security Measures:**
  - OTP verification for registration (via SMS for phone numbers using Firebase).
  - Email-based password reset for secure recovery.
  - Device binding to ensure only registered devices can mark attendance.
  - Developer option detection to prevent misuse.
  - HTTPS for secure data transmission.
  - JWT-based authentication for API access.
  - Audit logging for critical operations.

### Enhanced API Security

- Input validation and sanitization on all endpoints.
- Rate limiting to protect against brute-force attacks.
- Two-factor authentication for sensitive operations (where applicable).
- Regular security audits and vulnerability scanning.

---

## Implementation Details

- **Frontend (Flutter):**
  - Build separate widget trees for teacher and student flows.
  - Use BLoC for managing authentication, session creation, and attendance marking.
  - Integrate device APIs (GPS, Wi-Fi) and include animated onboarding to guide users through verification.
  - **Integrate Firebase Authentication:**  
    • Use SMS code verification for phone numbers during registration via Firebase.  
    • Use email-based password reset flows for secure recovery.
  
- **Backend (Golang):**
  - Develop RESTful APIs using the Gin framework.
  - Implement endpoints for registration, OTP verification, login, session management, attendance recording, and password reset.
  - Enforce business rules such as academic year matching and single attendance per subject per day.
  - **Integrate Firebase Admin SDK:** Verify Firebase ID tokens in endpoints and use the persistent Firebase UID to link user records in MySQL.
  - Integrate Gemini AI for anomaly detection and reporting.

---

## Step-by-Step Implementation Plan

### Phase 1: Project Setup and Basic Structure (1-2 days)

**Frontend (Flutter):**
- Create a new Flutter project using `flutter create smart_attendance_app`.
- Set up the folder structure as defined in the project structure.
- Add essential dependencies in `pubspec.yaml` (flutter_bloc, dio, get_it, etc.).
- Configure basic theming, routing, and dependency injection.

**Backend (Golang):**
- Initialize a new Golang project in the `backend/` folder.
- Set up the folder structure (cmd, config, controllers, models, routes, services, repositories, middleware, utils).
- Configure environment variables and logging.
- Create a basic server entry point using the Gin framework.
- Set up version control and dependency management with `go.mod`.

---

### Phase 2: Authentication Module Implementation (4-5 days)

**Frontend (Flutter):**
- Develop registration screens:
  - **Teacher Registration:** Fields for Full Name, Username, Email, Phone, Highest Degree, Password, Confirm Password, Experience.
  - **Student Registration:** Fields for Full Name, Roll Number, Course, Academic Year, Phone, Password, Confirm Password.
- Implement OTP verification screens and integrate with the registration flow, ensuring that SMS codes (via Firebase) are sent to verify phone numbers.
- Create login screens for teachers and students.
- Implement BLoC for managing authentication states.
- Connect the authentication screens to the backend via API calls using Dio.

**Backend (Golang):**
- Implement API endpoints for:
  - Teacher Registration (`POST /auth/register/teacher`).
  - Student Registration (`POST /auth/register/student`).
  - OTP Verification (`POST /auth/verify-otp`).
  - Teacher Login (`POST /auth/login/teacher`).
  - Student Login (`POST /auth/login/student`).
  - Password Reset (`POST /auth/reset-password`) – which sends a reset link via email.
- Implement OTP generation and validation logic.
- Use JWT for session management.
- Save user data to MySQL.
- Write unit tests for authentication endpoints.

---

### Phase 3: Teacher Features Implementation (4-5 days)

**Frontend (Flutter):**
- Develop the teacher dashboard:
  - Display the teacher's name.
  - Create dropdown menus for Academic Year and Subject selection.
  - Implement the "Start Attendance" button.
- Build the session creation UI:
  - Allow selection of countdown timer options (30s, 1m, 3m).
  - Connect the UI to BLoC to manage session state.
- Implement real-time session monitoring:
  - Display active session status and live attendance records.
  - Provide controls to end sessions early or view statistics.

**Backend (Golang):**
- Create API endpoint for session creation (`POST /sessions/start`):
  - Validate teacher credentials, academic year, and subject.
  - Store session details (teacher ID, academic year, subject, start/end times, countdown duration) in MySQL.
- Implement API endpoint to retrieve active session status (`GET /sessions/active`).
- Enforce business rules and log session actions.
- Write unit tests for session management endpoints.

---

### Phase 4: Student Features Implementation (3-4 days)

**Frontend (Flutter):**
- Develop the student dashboard:
  - Display the student's name.
  - Implement the "Mark Attendance" button (initially disabled).
- Add logic to detect an active session using polling or WebSocket.
- Implement an animated onboarding process:
  - Display real-time verification progress for Wi-Fi, GPS, device binding, and developer option check.
- After onboarding and academic year verification, enable the "Mark Attendance" button.
- On tap, send attendance data (student ID, session ID, verification details) to the backend.

**Backend (Golang):**
- Create API endpoint for attendance marking (`POST /attendance/mark`):
  - Validate that an active session exists.
  - Check that the student's Academic Year matches the session's Academic Year.
  - Validate GPS, Wi-Fi, device binding, and developer option status.
  - Ensure the student has not already marked attendance for the subject on the same day.
  - Record attendance in the database.
- Write tests for attendance verification and recording.

---

### Phase 5: Location, Wi-Fi, and Device Services (2-3 days)

**Frontend (Flutter):**
- Integrate location services using geolocator.
- Implement Wi-Fi verification using network_info_plus.
- Add device binding logic and check for registered devices.
- Implement developer option detection.
- Integrate these verifications into the onboarding animation and attendance marking flow.

**Backend (Golang):**
- Enhance the attendance endpoint to:
  - Validate GPS coordinates against campus geofence.
  - Check Wi-Fi SSID/BSSID against the institution’s credentials.
  - Verify device binding and developer option status.
- Update API documentation for verification processes.
- Write tests simulating various verification scenarios.

---

### Phase 6: Data Synchronization and Offline Support (2-3 days)

**Frontend (Flutter):**
- Implement local data storage using SQLite for offline support.
- Develop background sync tasks to push offline data to the backend.
- Handle conflict resolution and update BLoC with sync status.

**Backend (Golang):**
- Develop API endpoints that support delta updates and synchronization.
- Implement logic to merge offline attendance records upon reconnection.
- Ensure data consistency and write tests for offline scenarios.

---

### Phase 7: Testing, Error Handling, and Security Enhancements (3-4 days)

**Frontend (Flutter):**
- Write unit tests for BLoC components and widget tests.
- Integrate error handling to display user-friendly messages.
- Profile performance using Flutter DevTools.
- Secure sensitive data and tokens with encrypted storage.

**Backend (Golang):**
- Write unit and integration tests for all endpoints.
- Implement global error handling middleware.
- Secure API endpoints with JWT, HTTPS, and input validation.
- Use static code analysis and logging to enforce code quality.

---

### Phase 8: UI/UX Refinement and Performance Optimization (2-3 days)

**Frontend (Flutter):**
- Refine UI designs, animations, and loading states.
- Optimize widget rebuilds and network calls.
- Conduct usability tests and adjust layouts for responsiveness.

**Backend (Golang):**
- Optimize database queries and API response times.
- Implement caching strategies where applicable.
- Refactor code for improved maintainability and performance.
- Perform load testing and optimize concurrency handling.

---

### Phase 9: Deployment and Documentation (2-3 days)

**Frontend (Flutter):**
- Prepare release builds for Android and iOS.
- Set up CI/CD pipelines for automated testing and deployment.
- Document UI components, state management, and integration points.

**Backend (Golang):**
- Containerize the application using Docker.
- Set up deployment pipelines with CI/CD tools.
- Write comprehensive API documentation and deployment guides.
- Monitor staging environments before production release.

---

### Phase 10: Monitoring, Analytics, and Responsive Design (2-3 days)

**Frontend (Flutter):**
- Integrate analytics tools to monitor user engagement and performance.
- Set up error tracking with services like Sentry.
- Ensure responsive design adapts to various device sizes and orientations.

**Backend (Golang):**
- Implement real-time monitoring using tools like Prometheus and Grafana.
- Set up logging and alerting systems for proactive issue detection.
- Monitor API usage patterns and performance metrics.

---

### Phase 11: API Testing with Postman Locally

- **Setup Postman Collections:** Create and organize API endpoints for authentication, session management, and attendance.
- **Local Testing:** Verify endpoint responses, error handling, and security mechanisms locally before deployment.
- **Documentation:** Ensure Postman documentation is updated and shared with the development team.

---

## Project Structure

### Frontend (Flutter)

```
smart_attendance_app/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   ├── config/
│   │   ├── app_config.dart
│   │   ├── theme.dart
│   │   └── routes.dart
│   ├── core/
│   │   ├── constants/
│   │   ├── errors/
│   │   ├── utils/
│   │   └── services/
│   ├── data/
│   │   ├── models/
│   │   ├── repositories/
│   │   └── datasources/
│   ├── domain/
│   │   ├── entities/
│   │   ├── repositories/
│   │   └── usecases/
│   ├── presentation/
│   │   ├── screens/
│   │   │   ├── auth/           // Sign Up, OTP Verification, Login, and Password Reset screens for Teacher and Student
│   │   │   ├── teacher/        // Dashboard, Session Management, and Attendance Monitoring for Teachers
│   │   │   └── student/        // Dashboard, Onboarding, and Attendance Marking for Students
│   │   ├── widgets/
│   │   │   └── common/
│   │   └── bloc/               // BLoC files for authentication, session, and attendance logic
│   └── di/                     // Dependency Injection setup
├── assets/
│   ├── images/
│   ├── fonts/
│   └── icons/
├── test/
│   ├── unit/
│   ├── widget/
│   └── integration/
├── android/
├── ios/
├── web/
├── pubspec.yaml
├── analysis_options.yaml
└── README.md
```

### Backend (Golang)

```
backend/
├── cmd/                       // Main entry points for the application
│   └── server/                // Main server package with initialization
├── config/                    // Configuration files (environment variables, logging, etc.)
├── controllers/               // API endpoint controllers (authentication, sessions, attendance)
├── models/                    // Data models corresponding to the database schema
├── routes/                    // Route definitions using the Gin framework
├── services/                  // Business logic and integrations (Gemini AI, OTP, verification)
├── repositories/              // Database operations and ORM models
├── middleware/                // Middleware for authentication, error handling, logging
├── utils/                     // Utility functions and helper methods
├── Dockerfile                 // Containerization configuration
└── go.mod                     // Go modules and dependency management
```

---

## API Endpoints

### **Authentication & Registration**

| Method | Endpoint                  | Description                                                            |
|--------|---------------------------|------------------------------------------------------------------------|
| POST   | `/auth/register/teacher`  | Register a new teacher                                                 |
| POST   | `/auth/register/student`  | Register a new student                                                 |
| POST   | `/auth/verify-otp`        | Verify OTP (SMS code) during registration                              |
| POST   | `/auth/login/teacher`     | Teacher login                                                          |
| POST   | `/auth/login/student`     | Student login                                                          |
| POST   | `/auth/reset-password`    | Initiate password reset via email                                      |

### **Session Management**

| Method | Endpoint            | Description                                  |
|--------|---------------------|----------------------------------------------|
| POST   | `/sessions/start`   | Start an attendance session                  |
| GET    | `/sessions/active`  | Fetch active session details                 |
| PATCH  | `/sessions/end`     | End an active session manually               |

### **Attendance Marking**

| Method | Endpoint                | Description                                             |
|--------|-------------------------|---------------------------------------------------------|
| POST   | `/attendance/mark`      | Mark attendance (GPS, Wi-Fi verified)                   |
| GET    | `/attendance/status`    | Fetch current attendance status                         |

### **Security & Verification**

| Method | Endpoint                           | Description                                   |
|--------|------------------------------------|-----------------------------------------------|
| GET    | `/security/check-developer-mode`   | Detect if developer mode is enabled           |
| GET    | `/security/check-device-binding`     | Verify if the device is registered             |
| POST   | `/security/report-fraud`            | Report suspicious activity                     |

### **Reports & Analytics**

| Method | Endpoint                   | Description                                          |
|--------|----------------------------|------------------------------------------------------|
| GET    | `/reports/attendance`      | Generate attendance reports                          |
| GET    | `/reports/fraud-detection` | AI-based fraud detection analytics                   |

---

## Conclusion

This comprehensive document outlines the complete flow and implementation plan for the Smart Attendance App—covering every UI step from registration and login (with SMS-based OTP verification via Firebase and email-based password resets) to session creation by teachers and attendance marking by students. The solution leverages Flutter for a dynamic, responsive frontend and Golang for a scalable, secure backend, with user data stored in a MySQL database. Integration with Firebase ensures secure user verification while linking with custom business logic and data management. AI-powered security via Gemini AI and real-time verification mechanisms provide a reliable and secure attendance process.
```

---

This updated version of the universeupdated.md file maintains all previous content while incorporating the new Firebase processes and concepts. You can now use this document as your comprehensive specification for building your production-level Smart Attendance App.