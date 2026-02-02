// Package main - Data seeder for LMS
// Generates realistic test data via API endpoints
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Configuration
var (
	baseURL = getEnv("API_URL", "http://localhost:8085/api/v1")
	client  = &http.Client{Timeout: 30 * time.Second}
)

// ============================================================================
// Data Models
// ============================================================================

// APIResponse wraps all API responses
type APIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

type CourseResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type AssignmentResponse struct {
	ID       string `json:"id"`
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
}

type SubmissionResponse struct {
	ID           string `json:"id"`
	AssignmentID string `json:"assignment_id"`
	StudentID    string `json:"student_id"`
}

type EnrollmentResponse struct {
	ID        string `json:"id"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}

type RoomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type EventResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// ============================================================================
// Sample Data
// ============================================================================

var (
	// Admin user (created first)
	adminUser = map[string]interface{}{
		"email":      "admin@lms.local",
		"password":   "Admin123!@#",
		"role":       "manager",
		"first_name": "System",
		"last_name":  "Administrator",
	}

	// Teachers
	teachers = []map[string]interface{}{
		{"email": "john.smith@lms.local", "password": "Teacher123!", "role": "teacher", "first_name": "John", "last_name": "Smith", "department": "Computer Science"},
		{"email": "maria.garcia@lms.local", "password": "Teacher123!", "role": "teacher", "first_name": "Maria", "last_name": "Garcia", "department": "Mathematics"},
		{"email": "david.johnson@lms.local", "password": "Teacher123!", "role": "teacher", "first_name": "David", "last_name": "Johnson", "department": "Physics"},
		{"email": "sarah.williams@lms.local", "password": "Teacher123!", "role": "teacher", "first_name": "Sarah", "last_name": "Williams", "department": "Computer Science"},
		{"email": "michael.brown@lms.local", "password": "Teacher123!", "role": "teacher", "first_name": "Michael", "last_name": "Brown", "department": "Engineering"},
	}

	// Students
	students = []map[string]interface{}{
		{"email": "alice.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Alice", "last_name": "Anderson", "group_name": "CS-101"},
		{"email": "bob.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Bob", "last_name": "Baker", "group_name": "CS-101"},
		{"email": "charlie.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Charlie", "last_name": "Clark", "group_name": "CS-101"},
		{"email": "diana.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Diana", "last_name": "Davis", "group_name": "CS-102"},
		{"email": "evan.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Evan", "last_name": "Edwards", "group_name": "CS-102"},
		{"email": "fiona.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Fiona", "last_name": "Foster", "group_name": "MATH-101"},
		{"email": "george.student@lms.local", "password": "Student123!", "role": "student", "first_name": "George", "last_name": "Green", "group_name": "MATH-101"},
		{"email": "hannah.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Hannah", "last_name": "Hill", "group_name": "PHYS-101"},
		{"email": "ivan.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Ivan", "last_name": "Ivanov", "group_name": "PHYS-101"},
		{"email": "julia.student@lms.local", "password": "Student123!", "role": "student", "first_name": "Julia", "last_name": "Jones", "group_name": "ENG-101"},
	}

	// Courses with descriptions
	courses = []map[string]interface{}{
		{
			"title":       "Introduction to Programming",
			"description": "Learn the fundamentals of programming using Python. This course covers variables, control structures, functions, and basic data structures.",
		},
		{
			"title":       "Data Structures and Algorithms",
			"description": "Advanced programming concepts including arrays, linked lists, trees, graphs, sorting, and searching algorithms.",
		},
		{
			"title":       "Calculus I",
			"description": "Introduction to differential calculus covering limits, derivatives, and their applications.",
		},
		{
			"title":       "Linear Algebra",
			"description": "Study of vectors, matrices, linear transformations, and systems of linear equations.",
		},
		{
			"title":       "Classical Mechanics",
			"description": "Fundamentals of physics including motion, forces, energy, and momentum.",
		},
		{
			"title":       "Digital Electronics",
			"description": "Introduction to digital logic, Boolean algebra, and basic circuit design.",
		},
		{
			"title":       "Database Systems",
			"description": "Database design, SQL, normalization, and transaction management.",
		},
		{
			"title":       "Web Development",
			"description": "Full-stack web development with HTML, CSS, JavaScript, and backend technologies.",
		},
	}

	// Assignments templates
	assignmentTemplates = []map[string]interface{}{
		{"title": "Homework 1: Basic Concepts", "description": "Complete exercises on fundamental concepts covered in weeks 1-2.", "max_points": 100},
		{"title": "Homework 2: Practice Problems", "description": "Solve the given practice problems to reinforce your understanding.", "max_points": 100},
		{"title": "Lab Assignment 1", "description": "Hands-on lab exercise to apply theoretical knowledge.", "max_points": 50},
		{"title": "Midterm Project", "description": "Individual project demonstrating mastery of first half of course material.", "max_points": 200},
		{"title": "Homework 3: Advanced Topics", "description": "Exercises covering advanced topics from weeks 5-7.", "max_points": 100},
		{"title": "Group Project", "description": "Collaborative project working in teams of 2-3 students.", "max_points": 300},
		{"title": "Final Exam Prep", "description": "Practice problems to prepare for the final examination.", "max_points": 50},
	}

	// Schedule event types
	eventTypes = []string{"class", "lab", "exam", "office_hours"}

	// Chat message templates
	chatMessages = []string{
		"Hello everyone!",
		"Can someone help me with the homework?",
		"When is the next assignment due?",
		"Great lecture today!",
		"I'm having trouble understanding this concept.",
		"Does anyone want to form a study group?",
		"Thanks for the explanation!",
		"See you in class tomorrow.",
		"Has anyone started the project yet?",
		"What time is office hours?",
	}

	// Submission content templates
	submissionContents = []string{
		"Here is my solution to the assignment. I approached the problem by first analyzing the requirements and then implementing the solution step by step.",
		"My submission includes the completed exercises. I found this assignment challenging but educational.",
		"Attached is my work for this assignment. Please let me know if you have any questions.",
		"I have completed all the required tasks. The solution follows the guidelines provided in class.",
		"This is my submission for the assignment. I would appreciate any feedback on my approach.",
	}

	// Grade feedback templates
	gradeFeedbacks = []string{
		"Excellent work! Your solution demonstrates a strong understanding of the concepts.",
		"Good effort. Consider reviewing the lecture notes for better optimization.",
		"Well done. Your code is clean and well-documented.",
		"Satisfactory work. Please pay more attention to edge cases.",
		"Great improvement from the previous assignment!",
		"Nice work overall, but there are some minor issues to address.",
	}

	// Attendance notes templates
	attendanceNotes = []string{
		"",
		"Arrived 5 minutes late",
		"Left early due to medical appointment",
		"Participated actively in discussion",
		"Technical issues with connection",
	}
)

// ============================================================================
// State (tracked during seeding)
// ============================================================================

type SeederState struct {
	AdminToken     string
	TeacherTokens  map[string]string // email -> token
	StudentTokens  map[string]string // email -> token
	TeacherIDs     map[string]string // email -> user_id
	StudentIDs     map[string]string // email -> user_id
	Courses        []CourseResponse
	Assignments    []AssignmentResponse
	Enrollments    []EnrollmentResponse
	Submissions    []SubmissionResponse
	ChatRooms      []RoomResponse
	ScheduleEvents []EventResponse
}

func NewSeederState() *SeederState {
	return &SeederState{
		TeacherTokens: make(map[string]string),
		StudentTokens: make(map[string]string),
		TeacherIDs:    make(map[string]string),
		StudentIDs:    make(map[string]string),
	}
}

// ============================================================================
// HTTP Client Helpers
// ============================================================================

func doRequest(method, endpoint string, body interface{}, token string) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

func post(endpoint string, body interface{}, token string) ([]byte, error) {
	respBody, status, err := doRequest("POST", endpoint, body, token)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("POST %s failed with status %d: %s", endpoint, status, string(respBody))
	}
	return respBody, nil
}

func get(endpoint string, token string) ([]byte, error) {
	respBody, status, err := doRequest("GET", endpoint, nil, token)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("GET %s failed with status %d: %s", endpoint, status, string(respBody))
	}
	return respBody, nil
}

func put(endpoint string, body interface{}, token string) ([]byte, error) {
	respBody, status, err := doRequest("PUT", endpoint, body, token)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("PUT %s failed with status %d: %s", endpoint, status, string(respBody))
	}
	return respBody, nil
}

// parseResponse extracts data from wrapped API response
func parseResponse(respBody []byte, target interface{}) error {
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("parse API response: %w", err)
	}

	if !apiResp.Success {
		if apiResp.Error != nil {
			return fmt.Errorf("API error: %s - %s", apiResp.Error.Code, apiResp.Error.Message)
		}
		return fmt.Errorf("API returned success=false")
	}

	if target != nil && apiResp.Data != nil {
		if err := json.Unmarshal(apiResp.Data, target); err != nil {
			return fmt.Errorf("parse response data: %w", err)
		}
	}

	return nil
}

// ============================================================================
// Seeder Functions
// ============================================================================

func (s *SeederState) registerAndLogin(user map[string]interface{}) (*LoginResponse, error) {
	// Register
	_, err := post("/auth/register", user, "")
	if err != nil {
		// User might already exist, try login
		log.Printf("Registration failed (might exist): %v", err)
	}

	// Login
	loginBody := map[string]interface{}{
		"email":    user["email"],
		"password": user["password"],
	}
	respBody, err := post("/auth/login", loginBody, "")
	if err != nil {
		return nil, fmt.Errorf("login failed for %s: %w", user["email"], err)
	}

	var loginResp LoginResponse
	if err := parseResponse(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("parse login: %w", err)
	}

	return &loginResp, nil
}

func (s *SeederState) seedUsers() error {
	log.Println("📝 Seeding users...")

	// Register and login admin
	adminResp, err := s.registerAndLogin(adminUser)
	if err != nil {
		return fmt.Errorf("admin: %w", err)
	}
	s.AdminToken = adminResp.AccessToken
	log.Printf("  ✓ Admin: %s", adminUser["email"])

	// Register and login teachers
	for _, teacher := range teachers {
		resp, err := s.registerAndLogin(teacher)
		if err != nil {
			log.Printf("  ✗ Teacher %s: %v", teacher["email"], err)
			continue
		}
		s.TeacherTokens[teacher["email"].(string)] = resp.AccessToken
		s.TeacherIDs[teacher["email"].(string)] = resp.User.ID
		log.Printf("  ✓ Teacher: %s", teacher["email"])
	}

	// Register and login students
	for _, student := range students {
		resp, err := s.registerAndLogin(student)
		if err != nil {
			log.Printf("  ✗ Student %s: %v", student["email"], err)
			continue
		}
		s.StudentTokens[student["email"].(string)] = resp.AccessToken
		s.StudentIDs[student["email"].(string)] = resp.User.ID
		log.Printf("  ✓ Student: %s", student["email"])
	}

	return nil
}

func (s *SeederState) seedCourses() error {
	log.Println("📚 Seeding courses...")

	teacherEmails := make([]string, 0, len(s.TeacherTokens))
	for email := range s.TeacherTokens {
		teacherEmails = append(teacherEmails, email)
	}

	if len(teacherEmails) == 0 {
		return fmt.Errorf("no teachers available")
	}

	for i, course := range courses {
		// Distribute courses among teachers
		teacherEmail := teacherEmails[i%len(teacherEmails)]
		token := s.TeacherTokens[teacherEmail]

		respBody, err := post("/courses", course, token)
		if err != nil {
			log.Printf("  ✗ Course '%s': %v", course["title"], err)
			continue
		}

		var courseResp CourseResponse
		if err := parseResponse(respBody, &courseResp); err != nil {
			log.Printf("  ✗ Parse course response: %v", err)
			continue
		}

		s.Courses = append(s.Courses, courseResp)
		log.Printf("  ✓ Course: %s (Teacher: %s)", course["title"], teacherEmail)
	}

	return nil
}

func (s *SeederState) seedEnrollments() error {
	log.Println("📋 Seeding enrollments...")

	if len(s.Courses) == 0 {
		return fmt.Errorf("no courses available")
	}

	for studentEmail, token := range s.StudentTokens {
		// Enroll each student in 2-4 random courses
		numCourses := 2 + rand.Intn(3)
		enrolledCourses := make(map[int]bool)

		for i := 0; i < numCourses && i < len(s.Courses); i++ {
			courseIdx := rand.Intn(len(s.Courses))
			if enrolledCourses[courseIdx] {
				continue
			}
			enrolledCourses[courseIdx] = true

			course := s.Courses[courseIdx]
			enrollBody := map[string]interface{}{
				"course_id": course.ID,
			}

			respBody, err := post("/enrollments", enrollBody, token)
			if err != nil {
				log.Printf("  ✗ Enroll %s in %s: %v", studentEmail, course.Title, err)
				continue
			}

			var enrollResp EnrollmentResponse
			if err := parseResponse(respBody, &enrollResp); err != nil {
				log.Printf("  ✗ Parse enrollment response: %v", err)
				continue
			}

			s.Enrollments = append(s.Enrollments, enrollResp)
			log.Printf("  ✓ Enrolled: %s -> %s", studentEmail, course.Title)
		}
	}

	// Approve enrollments (teachers)
	log.Println("  Approving enrollments...")
	for _, enrollment := range s.Enrollments {
		// Find teacher token for this course
		var teacherToken string
		for _, c := range s.Courses {
			if c.ID == enrollment.CourseID {
				for email, token := range s.TeacherTokens {
					teacherToken = token
					_ = email
					break
				}
				break
			}
		}

		if teacherToken == "" {
			continue
		}

		statusBody := map[string]interface{}{
			"status": "active",
		}
		_, err := put("/enrollments/"+enrollment.ID+"/status", statusBody, teacherToken)
		if err != nil {
			log.Printf("  ✗ Approve enrollment %s: %v", enrollment.ID, err)
		}
	}

	return nil
}

func (s *SeederState) seedAssignments() error {
	log.Println("📝 Seeding assignments...")

	if len(s.Courses) == 0 {
		return fmt.Errorf("no courses available")
	}

	for _, course := range s.Courses {
		// Find teacher token for this course
		var teacherToken string
		for _, token := range s.TeacherTokens {
			teacherToken = token
			break
		}

		if teacherToken == "" {
			continue
		}

		// Create 3-5 assignments per course
		numAssignments := 3 + rand.Intn(3)
		for i := 0; i < numAssignments && i < len(assignmentTemplates); i++ {
			template := assignmentTemplates[i]

			// Set due date 1-4 weeks from now
			dueDate := time.Now().AddDate(0, 0, 7+rand.Intn(21))

			assignmentBody := map[string]interface{}{
				"course_id":   course.ID,
				"title":       template["title"],
				"description": template["description"],
				"max_points":  template["max_points"],
				"due_at":      dueDate.Format(time.RFC3339),
			}

			respBody, err := post("/assignments", assignmentBody, teacherToken)
			if err != nil {
				log.Printf("  ✗ Assignment '%s' for %s: %v", template["title"], course.Title, err)
				continue
			}

			var assignmentResp AssignmentResponse
			if err := parseResponse(respBody, &assignmentResp); err != nil {
				log.Printf("  ✗ Parse assignment response: %v", err)
				continue
			}

			s.Assignments = append(s.Assignments, assignmentResp)
			log.Printf("  ✓ Assignment: %s -> %s", template["title"], course.Title)
		}
	}

	return nil
}

func (s *SeederState) seedSubmissions() error {
	log.Println("📤 Seeding submissions...")

	if len(s.Assignments) == 0 {
		return fmt.Errorf("no assignments available")
	}

	// Get enrolled students for each course
	courseStudents := make(map[string][]string) // courseID -> []studentEmail
	for _, enrollment := range s.Enrollments {
		for email, id := range s.StudentIDs {
			if id == enrollment.StudentID {
				courseStudents[enrollment.CourseID] = append(courseStudents[enrollment.CourseID], email)
				break
			}
		}
	}

	for _, assignment := range s.Assignments {
		students := courseStudents[assignment.CourseID]
		if len(students) == 0 {
			continue
		}

		// 60-80% of students submit
		numSubmissions := len(students) * (60 + rand.Intn(20)) / 100
		if numSubmissions == 0 {
			numSubmissions = 1
		}

		for i := 0; i < numSubmissions && i < len(students); i++ {
			studentEmail := students[i]
			token := s.StudentTokens[studentEmail]

			content := submissionContents[rand.Intn(len(submissionContents))]
			submissionBody := map[string]interface{}{
				"assignment_id": assignment.ID,
				"content_text":  content,
			}

			respBody, err := post("/submissions", submissionBody, token)
			if err != nil {
				log.Printf("  ✗ Submission by %s for %s: %v", studentEmail, assignment.Title, err)
				continue
			}

			var submissionResp SubmissionResponse
			if err := parseResponse(respBody, &submissionResp); err != nil {
				log.Printf("  ✗ Parse submission response: %v", err)
				continue
			}

			s.Submissions = append(s.Submissions, submissionResp)
			log.Printf("  ✓ Submission: %s -> %s", studentEmail, assignment.Title)
		}
	}

	return nil
}

func (s *SeederState) seedGrades() error {
	log.Println("📊 Seeding grades...")

	if len(s.Submissions) == 0 {
		return fmt.Errorf("no submissions available")
	}

	// Grade 70-90% of submissions
	numToGrade := len(s.Submissions) * (70 + rand.Intn(20)) / 100

	for i := 0; i < numToGrade && i < len(s.Submissions); i++ {
		submission := s.Submissions[i]

		// Find a teacher token
		var teacherToken string
		for _, token := range s.TeacherTokens {
			teacherToken = token
			break
		}

		if teacherToken == "" {
			continue
		}

		// Find assignment to get max points
		var maxPoints int = 100
		for _, a := range s.Assignments {
			if a.ID == submission.AssignmentID {
				maxPoints = 100 // Default
				break
			}
		}

		// Random score 60-100% of max points
		score := maxPoints * (60 + rand.Intn(40)) / 100
		feedback := gradeFeedbacks[rand.Intn(len(gradeFeedbacks))]

		gradeBody := map[string]interface{}{
			"submission_id": submission.ID,
			"score":         score,
			"feedback":      feedback,
		}

		_, err := post("/grades", gradeBody, teacherToken)
		if err != nil {
			log.Printf("  ✗ Grade submission %s: %v", submission.ID, err)
			continue
		}

		log.Printf("  ✓ Graded: %s (%d points)", submission.ID[:8], score)
	}

	return nil
}

func (s *SeederState) seedAttendance() error {
	log.Println("📅 Seeding attendance sessions...")

	if len(s.Courses) == 0 {
		return fmt.Errorf("no courses available")
	}

	// Get enrolled students for each course
	courseStudents := make(map[string][]string) // courseID -> []studentID
	for _, enrollment := range s.Enrollments {
		courseStudents[enrollment.CourseID] = append(courseStudents[enrollment.CourseID], enrollment.StudentID)
	}

	for _, course := range s.Courses {
		// Find teacher token
		var teacherToken string
		for _, token := range s.TeacherTokens {
			teacherToken = token
			break
		}

		if teacherToken == "" {
			continue
		}

		// Create 3-5 attendance sessions per course (past dates)
		numSessions := 3 + rand.Intn(3)
		for i := 0; i < numSessions; i++ {
			sessionDate := time.Now().AddDate(0, 0, -7*(numSessions-i))

			sessionBody := map[string]interface{}{
				"course_id":    course.ID,
				"title":        fmt.Sprintf("Session %d - %s", i+1, course.Title),
				"session_date": sessionDate.Format("2006-01-02"),
				"start_time":   "09:00",
				"end_time":     "10:30",
			}

			respBody, err := post("/attendance/sessions", sessionBody, teacherToken)
			if err != nil {
				log.Printf("  ✗ Attendance session for %s: %v", course.Title, err)
				continue
			}

			var sessionResp struct {
				ID string `json:"id"`
			}
			if err := parseResponse(respBody, &sessionResp); err != nil {
				log.Printf("  ✗ Parse session response: %v", err)
				continue
			}

			log.Printf("  ✓ Session: %s", sessionBody["title"])

			// Mark attendance for enrolled students
			students := courseStudents[course.ID]
			if len(students) == 0 {
				continue
			}

			marks := make([]map[string]interface{}, 0)
			statuses := []string{"present", "present", "present", "present", "late", "absent", "excused"}

			for _, studentID := range students {
				status := statuses[rand.Intn(len(statuses))]
				note := attendanceNotes[rand.Intn(len(attendanceNotes))]

				mark := map[string]interface{}{
					"student_id": studentID,
					"status":     status,
				}
				if note != "" {
					mark["notes"] = note
				}
				marks = append(marks, mark)
			}

			bulkBody := map[string]interface{}{
				"marks": marks,
			}

			_, err = post("/attendance/sessions/"+sessionResp.ID+"/bulk-mark", bulkBody, teacherToken)
			if err != nil {
				log.Printf("  ✗ Mark attendance: %v", err)
			}
		}
	}

	return nil
}

func (s *SeederState) seedSchedule() error {
	log.Println("🗓️ Seeding schedule events...")

	if len(s.Courses) == 0 {
		return fmt.Errorf("no courses available")
	}

	for _, course := range s.Courses {
		// Find teacher token
		var teacherToken string
		for _, token := range s.TeacherTokens {
			teacherToken = token
			break
		}

		if teacherToken == "" {
			continue
		}

		// Create weekly class events
		for week := 0; week < 4; week++ {
			eventDate := time.Now().AddDate(0, 0, 7*week+rand.Intn(5))

			eventBody := map[string]interface{}{
				"course_id":   course.ID,
				"title":       fmt.Sprintf("Lecture: %s", course.Title),
				"description": "Regular class session",
				"event_type":  eventTypes[rand.Intn(len(eventTypes))],
				"start_time":  eventDate.Format("2006-01-02") + "T09:00:00Z",
				"end_time":    eventDate.Format("2006-01-02") + "T10:30:00Z",
				"location":    fmt.Sprintf("Room %d", 100+rand.Intn(200)),
				"is_online":   rand.Intn(2) == 1,
				"recurrence":  "none",
			}

			respBody, err := post("/schedule/events", eventBody, teacherToken)
			if err != nil {
				log.Printf("  ✗ Schedule event for %s: %v", course.Title, err)
				continue
			}

			var eventResp EventResponse
			if err := parseResponse(respBody, &eventResp); err != nil {
				log.Printf("  ✗ Parse event response: %v", err)
				continue
			}

			s.ScheduleEvents = append(s.ScheduleEvents, eventResp)
			log.Printf("  ✓ Event: %s (Week %d)", course.Title, week+1)
		}
	}

	return nil
}

func (s *SeederState) seedChat() error {
	log.Println("💬 Seeding chat rooms and messages...")

	if len(s.Courses) == 0 {
		return fmt.Errorf("no courses available")
	}

	// Create course chat rooms
	for _, course := range s.Courses {
		// Find teacher token
		var teacherToken string
		for _, token := range s.TeacherTokens {
			teacherToken = token
			break
		}

		if teacherToken == "" {
			continue
		}

		roomBody := map[string]interface{}{
			"name":      fmt.Sprintf("%s Discussion", course.Title),
			"type":      "course",
			"course_id": course.ID,
		}

		respBody, err := post("/chat/rooms", roomBody, teacherToken)
		if err != nil {
			log.Printf("  ✗ Chat room for %s: %v", course.Title, err)
			continue
		}

		var roomResp RoomResponse
		if err := parseResponse(respBody, &roomResp); err != nil {
			log.Printf("  ✗ Parse room response: %v", err)
			continue
		}

		s.ChatRooms = append(s.ChatRooms, roomResp)
		log.Printf("  ✓ Chat room: %s", roomBody["name"])

		// Add some messages
		tokens := make([]string, 0)
		for _, t := range s.TeacherTokens {
			tokens = append(tokens, t)
		}
		for _, t := range s.StudentTokens {
			tokens = append(tokens, t)
		}

		numMessages := 5 + rand.Intn(10)
		for i := 0; i < numMessages && len(tokens) > 0; i++ {
			token := tokens[rand.Intn(len(tokens))]
			message := chatMessages[rand.Intn(len(chatMessages))]

			msgBody := map[string]interface{}{
				"content": message,
			}

			_, err := post("/chat/rooms/"+roomResp.ID+"/messages", msgBody, token)
			if err != nil {
				log.Printf("  ✗ Send message: %v", err)
				continue
			}
		}
		log.Printf("    ✓ Added %d messages", numMessages)
	}

	// Create some direct chat rooms between random users
	studentEmails := make([]string, 0, len(s.StudentTokens))
	for email := range s.StudentTokens {
		studentEmails = append(studentEmails, email)
	}

	for i := 0; i < 3 && len(studentEmails) >= 2; i++ {
		student1 := studentEmails[rand.Intn(len(studentEmails))]
		student2 := studentEmails[rand.Intn(len(studentEmails))]
		if student1 == student2 {
			continue
		}

		token := s.StudentTokens[student1]
		directBody := map[string]interface{}{
			"user_id": s.StudentIDs[student2],
		}

		respBody, err := post("/chat/rooms/direct", directBody, token)
		if err != nil {
			log.Printf("  ✗ Direct room %s <-> %s: %v", student1, student2, err)
			continue
		}

		var roomResp RoomResponse
		if err := parseResponse(respBody, &roomResp); err != nil {
			log.Printf("  ✗ Parse direct room response: %v", err)
			continue
		}

		s.ChatRooms = append(s.ChatRooms, roomResp)
		log.Printf("  ✓ Direct chat: %s <-> %s", student1, student2)

		// Add some messages
		for j := 0; j < 3; j++ {
			token := s.StudentTokens[student1]
			if j%2 == 1 {
				token = s.StudentTokens[student2]
			}

			msgBody := map[string]interface{}{
				"content": chatMessages[rand.Intn(len(chatMessages))],
			}

			_, _ = post("/chat/rooms/"+roomResp.ID+"/messages", msgBody, token)
		}
	}

	return nil
}

func (s *SeederState) printSummary() {
	log.Println("")
	log.Println("=" + "==========================================================")
	log.Println("📊 SEEDING SUMMARY")
	log.Println("=" + "==========================================================")
	log.Printf("  👥 Teachers:     %d", len(s.TeacherTokens))
	log.Printf("  👨‍🎓 Students:     %d", len(s.StudentTokens))
	log.Printf("  📚 Courses:      %d", len(s.Courses))
	log.Printf("  📋 Enrollments:  %d", len(s.Enrollments))
	log.Printf("  📝 Assignments:  %d", len(s.Assignments))
	log.Printf("  📤 Submissions:  %d", len(s.Submissions))
	log.Printf("  💬 Chat rooms:   %d", len(s.ChatRooms))
	log.Printf("  🗓️ Schedule events: %d", len(s.ScheduleEvents))
	log.Println("=" + "==========================================================")
	log.Println("")
	log.Println("🔑 TEST CREDENTIALS:")
	log.Println("  Admin:   admin@lms.local / Admin123!@#")
	log.Println("  Teacher: john.smith@lms.local / Teacher123!")
	log.Println("  Student: alice.student@lms.local / Student123!")
	log.Println("")
}

// ============================================================================
// Main
// ============================================================================

func main() {
	rand.Seed(time.Now().UnixNano())

	// Check for university mode
	if len(os.Args) > 1 && os.Args[1] == "university" {
		log.Println("🎓 Running University Seeder (SE-2430)...")
		log.Printf("   API URL: %s", baseURL)
		log.Println("")

		// Check API health
		_, err := get("/health", "")
		if err != nil {
			log.Fatalf("❌ API not available: %v", err)
		}
		log.Println("✓ API is healthy")

		RunUniversitySeeder()
		log.Println("✅ University seeding completed!")
		return
	}

	log.Println("🚀 Starting LMS Data Seeder")
	log.Printf("   API URL: %s", baseURL)
	log.Println("")
	log.Println("💡 Tip: Run 'go run . university' for SE-2430 demo data")
	log.Println("")

	state := NewSeederState()

	// Check API health
	_, err := get("/health", "")
	if err != nil {
		log.Fatalf("❌ API not available: %v", err)
	}
	log.Println("✓ API is healthy")
	log.Println("")

	// Run seeders in order
	seeders := []struct {
		name string
		fn   func() error
	}{
		{"Users", state.seedUsers},
		{"Courses", state.seedCourses},
		{"Enrollments", state.seedEnrollments},
		{"Assignments", state.seedAssignments},
		{"Submissions", state.seedSubmissions},
		{"Grades", state.seedGrades},
		{"Attendance", state.seedAttendance},
		{"Schedule", state.seedSchedule},
		{"Chat", state.seedChat},
	}

	for _, seeder := range seeders {
		if err := seeder.fn(); err != nil {
			log.Printf("⚠️ %s seeder failed: %v", seeder.name, err)
		}
		log.Println("")
	}

	state.printSummary()
	log.Println("✅ Seeding completed!")
}

// ============================================================================
// Helpers
// ============================================================================

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
