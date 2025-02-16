Below is a comprehensive Markdown document outlining the complete registration and login process for your production-grade hybrid system. This document details the architecture, API flows, and the database schema so you can start building the solution together.

---

```markdown
# Production-Grade Registration & Login System

This document outlines the high-level steps, architecture, API flows, and database schema for building a secure registration and login system using a hybrid approach:
- **Firebase** is used for OTP-based phone verification and password resets.
- **Golang Backend** handles token verification, custom business logic, and API management.
- **MySQL Database** is used to store persistent user data and related records.
- **Flutter Frontend** (not covered in detail here) integrates with Firebase for the client-side authentication flow.

---

## Table of Contents

1. Overview
2. System Architecture & Flow
3. API Endpoints
4. Database Schema
5. Deployment & Production Considerations
6. Summary

---

## 1. Overview

Our hybrid authentication system leverages:
- **Firebase**: For secure OTP verification and password reset flows.
- **Golang Backend**: To verify Firebase tokens, enforce custom business logic, and generate custom JWTs for session management.
- **MySQL Database**: To persist user records and maintain relationships (using the Firebase UID as a unique link).

This approach minimizes sensitive logic on the client side while ensuring you maintain full control over your business logic and data management.

---

## 2. System Architecture & Flow

### Registration Flow

1. **User Input Collection**  
   - The user enters details such as full name, phone number, email, academic year, and role (teacher/student) on the registration screen.

2. **Firebase OTP Verification**  
   - The Flutter app triggers Firebase’s OTP verification (using methods like `verifyPhoneNumber`).
   - Once the OTP is verified, Firebase returns a verified user object containing a persistent Firebase UID and an ID token.

3. **Token Hand-off to Backend**  
   - The Flutter app sends the verified Firebase ID token along with the registration details to your Golang backend (e.g., via an endpoint `/auth/register`).

4. **Backend Token Verification & User Creation**  
   - The Golang backend verifies the Firebase token using the Firebase Admin SDK.
   - It extracts the Firebase UID, checks if the user already exists, and then creates a new user record in the MySQL database.
   - Optionally, a custom JWT is generated for subsequent API calls.

### Login Flow

1. **User Login Initiation**  
   - The user logs in via OTP (or password-based, if supported).
   - For OTP-based login, the Flutter app repeats the Firebase OTP process and obtains a verified token.
   - For password-based login, credentials are verified directly against the MySQL records.

2. **Token Verification & Session Management**  
   - The Golang backend verifies the token or credentials and issues a custom JWT or session token.
   - This token is used to authorize further API requests.

---

## 3. API Endpoints

### Authentication Endpoints

- **POST /auth/register**  
  - **Input:** Firebase ID token, full name, phone, email, academic year, role, etc.  
  - **Process:** Verify the token using Firebase Admin SDK, extract the UID, and create or update a user record in MySQL.  
  - **Output:** Success response along with a custom JWT.

- **POST /auth/login**  
  - **Input:** traditional credentials(Teacher-username,password & Student- Roll Number,password).  
  - **Process:** Verify credentials and generate a session token.  
  - **Output:** Custom JWT/session token for subsequent requests.

- **POST /auth/verify-otp**  
  - **Purpose:** (Optional) Additional endpoint if you manage OTP flows separately for actions like password resets.

### Additional Endpoints

- **POST /auth/reset-password**  
  - Initiates a password reset process via Firebase’s built-in email functionalities.
- **GET /user/profile**  
  - Retrieves the authenticated user’s profile information.

- **PUT /user/update**  
  - Allows users to update their profile or registration details.

---

## 4. Database Schema

The following schema provides the core tables required to support the registration and login flows.

### Users Table

This table stores user details and links each record to Firebase using the `firebase_uid`.

```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,                  -- Internal unique identifier
    firebase_uid VARCHAR(128) UNIQUE NOT NULL,     -- Unique Firebase UID (persistent per user)
    full_name VARCHAR(100) NOT NULL,               -- User's full name
    username VARCHAR(100) UNIQUE,                  -- Optional username (for teachers)
    email VARCHAR(255),                            -- Email address (if provided)
    phone VARCHAR(20) NOT NULL,                    -- Phone number
    academic_year VARCHAR(50),                     -- Academic year (for students)
    role ENUM('teacher', 'student') NOT NULL,       -- User role
    password_hash VARCHAR(255),                    -- For password-based login (if applicable)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### OTP Verifications Table

This table is used for any custom OTP verification flows (e.g., password reset OTPs).

```sql
CREATE TABLE otp_verifications (
    id VARCHAR(36) PRIMARY KEY,                    -- Unique OTP verification record
    user_id VARCHAR(36) NOT NULL,                  -- Reference to the user (users.id)
    otp_code VARCHAR(10) NOT NULL,                 -- OTP code sent to the user
    expires_at TIMESTAMP NOT NULL,                 -- OTP expiration time
    verified BOOLEAN DEFAULT FALSE,                -- OTP verification status
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

### (Optional) Audit Logs Table

Tracks critical actions for security and debugging purposes.

```sql
CREATE TABLE audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    action VARCHAR(100) NOT NULL,
    details JSON,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

## 5. Deployment & Production Considerations

- **Environment Management:**  
  - Securely store API keys, Firebase service account credentials, and database connection details using environment variables or a secrets manager.
  - Maintain separate configurations for development, staging, and production.

- **Security Enhancements:**  
  - Use HTTPS for all communications.
  - Implement rate limiting, secure error handling, and input validation on all endpoints.
  - Regularly audit your code and update dependencies.

- **Monitoring & Logging:**  
  - Set up real-time monitoring and logging (using tools like Prometheus, Grafana, or Firebase Crashlytics).
  - Use automated alerts to respond to potential issues.

- **CI/CD Pipelines:**  
  - Automate testing, building, and deployment processes to ensure smooth production rollouts.
  - Use staging environments to thoroughly test new features before production deployment.

---

## 6. Summary

This document outlines a complete production-grade architecture for a registration and login system that:
- **Uses Firebase** for secure OTP-based verification and password reset flows.
- **Leverages a Golang backend** to verify tokens, manage custom business logic, and generate session tokens.
- **Stores persistent user data** in a MySQL database with a permanent link via the Firebase UID.
- **Supports modular API endpoints** for authentication, user management, and additional features.

With this guide and the provided database schema, you now have a detailed blueprint to begin developing your registration and login process. The structure is designed for security, scalability, and ease of maintenance in a production environment.
```

---

This Markdown file is intended to serve as your complete guide to starting the development of your hybrid authentication system. You can further refine or extend the sections as your project evolves.
