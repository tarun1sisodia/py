-- Create database if not exists and use it
CREATE DATABASE IF NOT EXISTS smart_attendance;
USE smart_attendance;

-- Drop tables if they exist (for a clean start)
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS admin_audit_logs;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS otp_verifications;
DROP TABLE IF EXISTS attendance_records;
DROP TABLE IF EXISTS attendance_sessions;
DROP TABLE IF EXISTS users;

-- Create users table (for both teachers and students)
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    role ENUM('teacher', 'student') NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    roll_number VARCHAR(50) UNIQUE,
    course VARCHAR(100),
    academic_year VARCHAR(50),
    phone VARCHAR(20) NOT NULL,
    highest_degree VARCHAR(100),
    experience VARCHAR(100),
    password VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create attendance_sessions table
CREATE TABLE attendance_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    teacher_id INT NOT NULL,
    session_date DATE NOT NULL,
    session_start DATETIME NOT NULL,
    session_end DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (teacher_id) REFERENCES users(id)
);

-- Create attendance_records table
CREATE TABLE attendance_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT NOT NULL,
    student_id INT NOT NULL,
    status ENUM('present', 'absent', 'late') DEFAULT 'absent',
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES attendance_sessions(id),
    FOREIGN KEY (student_id) REFERENCES users(id)
);

-- Create otp_verifications table
CREATE TABLE otp_verifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    otp VARCHAR(10) NOT NULL,
    expires_at DATETIME NOT NULL,
    attempts INT DEFAULT 0,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create audit_logs table
CREATE TABLE audit_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    action VARCHAR(255) NOT NULL,
    details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create admins table
CREATE TABLE admins (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create admin_audit_logs table
CREATE TABLE admin_audit_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    admin_id INT NOT NULL,
    action VARCHAR(255) NOT NULL,
    details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES admins(id)
);

-- Create reports table
CREATE TABLE reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT,
    report_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES attendance_sessions(id)
); 