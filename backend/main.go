package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Teacher represents a teacher user
type Teacher struct {
	ID            string `json:"id"`
	FullName      string `json:"full_name"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	HighestDegree string `json:"highest_degree"`
	Experience    string `json:"experience"`
	Password      string `json:"password"`
	Verified      bool   `json:"verified"`
}

// Student represents a student user
type Student struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	RollNumber   string `json:"roll_number"`
	Course       string `json:"course"`
	AcademicYear string `json:"academic_year"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	Verified     bool   `json:"verified"`
}

var (
	teachers       = make(map[string]Teacher)
	students       = make(map[string]Student)
	teacherOTP     = make(map[string]string) // maps teacher ID to OTP
	studentOTP     = make(map[string]string) // maps student ID to OTP
	teacherIDCount = 1
	studentIDCount = 1
	mu             sync.Mutex
)

// generateOTP creates a 6-digit OTP
func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// teacherRegisterHandler handles teacher registration
func teacherRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData map[string]string
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	teacherID := strconv.Itoa(teacherIDCount)
	teacherIDCount++
	mu.Unlock()

	teacher := Teacher{
		ID:            teacherID,
		FullName:      reqData["full_name"],
		Username:      reqData["username"],
		Email:         reqData["email"],
		Phone:         reqData["phone"],
		HighestDegree: reqData["highest_degree"],
		Experience:    reqData["experience"],
		Password:      reqData["password"],
		Verified:      false,
	}

	otp := generateOTP()

	mu.Lock()
	teachers[teacherID] = teacher
	teacherOTP[teacherID] = otp
	mu.Unlock()

	// Simulate sending OTP via SMS (in production, call the Twilio SendSMS function)
	go func(phone, otp string) {
		log.Printf("Sending OTP %s to teacher phone %s", otp, phone)
	}(teacher.Phone, otp)

	res := map[string]string{"userId": teacherID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// teacherLoginHandler handles teacher login
func teacherLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, password := reqData["email"], reqData["password"]

	mu.Lock()
	defer mu.Unlock()
	for _, t := range teachers {
		if t.Email == email && t.Password == password {
			if !t.Verified {
				http.Error(w, "Teacher not verified via OTP", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// studentRegisterHandler handles student registration
func studentRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	studentID := strconv.Itoa(studentIDCount)
	studentIDCount++
	mu.Unlock()

	student := Student{
		ID:           studentID,
		FullName:     reqData["full_name"],
		RollNumber:   reqData["roll_number"],
		Course:       reqData["course"],
		AcademicYear: reqData["academic_year"],
		Phone:        reqData["phone"],
		Password:     reqData["password"],
		Verified:     false,
	}

	otp := generateOTP()

	mu.Lock()
	students[studentID] = student
	studentOTP[studentID] = otp
	mu.Unlock()

	// Simulate sending OTP via SMS
	go func(phone, otp string) {
		log.Printf("Sending OTP %s to student phone %s", otp, phone)
	}(student.Phone, otp)

	res := map[string]string{"userId": studentID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// studentLoginHandler handles student login
func studentLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rollNumber, password := reqData["roll_number"], reqData["password"]

	mu.Lock()
	defer mu.Unlock()
	for _, s := range students {
		if s.RollNumber == rollNumber && s.Password == password {
			if !s.Verified {
				http.Error(w, "Student not verified via OTP", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(s)
			return
		}
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// resetPasswordHandler handles password reset using OTP
func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := reqData["user_id"]
	otp := reqData["otp"]
	newPassword := reqData["new_password"]

	mu.Lock()
	defer mu.Unlock()
	if teacher, ok := teachers[userId]; ok {
		if storedOTP, exists := teacherOTP[userId]; exists && storedOTP == otp {
			teacher.Password = newPassword
			teachers[userId] = teacher
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
			return
		}
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	if student, ok := students[userId]; ok {
		if storedOTP, exists := studentOTP[userId]; exists && storedOTP == otp {
			student.Password = newPassword
			students[userId] = student
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "success"})
			return
		}
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

// verifyOTPHandler handles OTP verification for both teachers and students
func verifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := reqData["user_id"]
	otp := reqData["otp"]

	mu.Lock()
	defer mu.Unlock()
	if teacher, ok := teachers[userId]; ok {
		if storedOTP, exists := teacherOTP[userId]; exists && storedOTP == otp {
			teacher.Verified = true
			teachers[userId] = teacher
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "verified"})
			return
		}
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	if student, ok := students[userId]; ok {
		if storedOTP, exists := studentOTP[userId]; exists && storedOTP == otp {
			student.Verified = true
			students[userId] = student
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "verified"})
			return
		}
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	http.HandleFunc("/api/v1/auth/register/teacher", teacherRegisterHandler)
	http.HandleFunc("/api/v1/auth/login/teacher", teacherLoginHandler)
	http.HandleFunc("/api/v1/auth/register/student", studentRegisterHandler)
	http.HandleFunc("/api/v1/auth/login/student", studentLoginHandler)
	http.HandleFunc("/api/v1/auth/reset-password", resetPasswordHandler)
	http.HandleFunc("/api/v1/auth/verify-otp", verifyOTPHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
