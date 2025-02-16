# Deployment and Server Setup Guide for Smart Attendance Backend

This document provides step-by-step instructions on how to set up, configure, update, deploy, and run the backend server for the Smart Attendance project. It also includes troubleshooting tips and common commands you can use to manage your server.

---
Every time you change the elastic ip address, you need to update the .env file with the new ip address. and also injection_container .dart and also which files which are required to be changed. but injection_container will enough to change as per server's configuration.




## 1. Prerequisites

- **Operating System:** Ubuntu or other Linux distro
- **Required Software:**
  - Go (version 1.21 as per go.mod)
  - MySQL
  - Git
  - Systemd (for service management)

Ensure that your environment variables and configuration files (e.g., `.env`) are updated with the correct credentials and values.

---

## 2. Initial Server Setup

### a. Clone the Repository (if not already cloned)

```bash
git clone https://github.com/yourusername/smart_attendance_server.git
cd smart_attendance_server
```

### b. Environment Configuration

- Create a `.env` file in the project root using the provided sample.
- Example `.env` file:

```env
PORT=8080
GIN_MODE=release

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root1234
DB_NAME=smart_attendance

# JWT Configuration
JWT_SECRET=Yw6Pn3PkxZfFvM+vxsxKzH8jQ9xJGJj2fqwHDuYz9AM=
JWT_EXPIRATION=24h

# Server Configuration
SERVER_HOST=localhost
ALLOWED_ORIGINS=*

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@example.com
SMTP_PASSWORD=your-password
FROM_EMAIL=your_email@example.com
```

- Secure the `.env` file:

```bash
chmod 600 .env
```

---

## 3. Database Setup

Use the provided schema file to set up the database.

```bash
mysql -h localhost -P 3306 -u root -p"root1234" < ./db/schema.sql
```

---

## 4. Building and Deploying the Application

### a. Using the Deployment Script

A deployment script (`scripts/deploy.sh`) is provided to automate the deployment process.

To execute the deployment script:

```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh
```

The script performs the following actions:

- Sets up the database (if not already done)
- Builds the Go application and places the binary in the `bin` folder
- Sets up a systemd service for the backend
- Creates or updates the `.env` file with the correct configuration
- Reloads the systemd daemon and starts the service

### b. Manual Deployment Steps (if needed)

1. **Update from GitHub:**
    
    ```bash
git pull origin main
    ```
    
2. **Build the Application:**

    ```bash
mkdir -p bin
go build -o bin/server ./cmd/server/main.go
    ```

3. **Restart the Service:**

    Reload the systemd daemon and restart the service:

    ```bash
sudo systemctl daemon-reload
sudo systemctl restart smart-attendance
    ```

4. **Check Service Status:**

    ```bash
sudo systemctl status smart-attendance
    ```

---

## 5. Common Commands and Operations

### a. Git Operations

- **Pull latest changes:**

  ```bash
git pull origin main
  ```

- **Commit and push changes (if making changes on the server):**

  ```bash
git add .
git commit -m "Your commit message"
git push origin main
  ```

### b. Systemd Service Management

- **Check status:**

  ```bash
sudo systemctl status smart-attendance
  ```

- **Restart service:**

  ```bash
sudo systemctl restart smart-attendance
  ```

- **Enable service on boot:**

  ```bash
sudo systemctl enable smart-attendance
  ```

- **View logs:**

  ```bash
sudo journalctl -u smart-attendance -f
  ```

### c. Building the Application Manually

If you need to make quick changes and rebuild manually:

```bash
mkdir -p bin
go build -o bin/server ./cmd/server/main.go
sudo systemctl restart smart-attendance
  ```

---

## 6. Troubleshooting

### a. Common Errors and How to Resolve Them

- **Error: Schema file not found**
  
  - Ensure that the path `./db/schema.sql` is correct. The deployment script should be run from the project root.

- **Build Errors in Go:**
  
  - Check the error messages in the terminal. Use `go fmt` and `go vet` to catch formatting and lint issues:

    ```bash
    go fmt ./...
    go vet ./...
    ```

- **Systemd Service Fails to Start:**
  
  - Check the service logs using:

    ```bash
    sudo journalctl -u smart-attendance
    ```
  
  - Verify that your environment variables are set correctly in the service file and `.env` file.

- **Network Issues:**
  
  - If you encounter connection timeout or similar network errors, ensure that your server's firewall rules allow traffic on port 8080 (or the port specified in your `.env` and service file).

### b. Updating Files and Deploying New Changes

1. **Make changes locally and push to GitHub:**
    
    ```bash
    git add .
    git commit -m "Describe your changes"
    git push origin main
    ```

2. **On the server, navigate to the project directory and update:**
    
    ```bash
    cd /path/to/smart_attendance_server
    git pull origin main
    ```

3. **Rebuild and restart the service:**
    
    ```bash
    go build -o bin/server ./cmd/server/main.go
    sudo systemctl restart smart-attendance
    ```

4. **Check logs for any errors:**

    ```bash
    sudo journalctl -u smart-attendance -f
    ```

---

## 7. Additional Tips

- **Automate Deployment:** Consider setting up a CI/CD pipeline to automatically deploy changes when new commits are pushed to GitHub.
- **Backup:** Regularly backup your database and configuration files.
- **Security:** Always secure sensitive files like `.env` and ensure that firewall settings allow only the necessary traffic (e.g., ports 80/443 if using a reverse proxy, and 8080 for internal communication if applicable).
- **Environment File Management:** Ensure that your actual `.env` file is not committed to GitHub. Add `.env` to your `.gitignore` and commit only a template file (e.g. `.env.example`). Manage sensitive credentials on the server or through secure secret management solutions. This way, when you run `git pull origin main`, your code updates without exposing sensitive information.
- **Monitoring:** Use tools like Prometheus or Grafana for monitoring, and set up log rotation for your logs.

---

## Conclusion

Following the steps in this guide will help you manage your server setup, deploy updates from GitHub, and troubleshoot common issues. For more detailed commands, refer to this guide and the inline comments in your deployment scripts.

Happy Deploying! 