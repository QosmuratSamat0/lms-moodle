// Package main - University Data Seeder for LMS
// Generates realistic university data: SE-2430 group, 6 subjects, 21 students, 10-week structure
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
	apiURL     = getEnvUniv("API_URL", "http://localhost:8085/api/v1")
	httpClient = &http.Client{Timeout: 30 * time.Second}
)

func getEnvUniv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// ============================================================================
// Data Models
// ============================================================================

type UnivAPIResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type UnivLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

type UnivCourseResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type UnivAssignmentResponse struct {
	ID       string `json:"id"`
	CourseID string `json:"course_id"`
	Title    string `json:"title"`
}

type UnivEnrollmentResponse struct {
	ID        string `json:"id"`
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}

// ============================================================================
// University Data - SE-2430 Software Engineering Group
// ============================================================================

// Admin
var univAdmin = map[string]interface{}{
	"email":      "admin@aitu.edu.kz",
	"password":   "Admin123!@#",
	"role":       "manager",
	"first_name": "System",
	"last_name":  "Administrator",
}

// Teachers for SE-2430 (6 subjects)
var univTeachers = []map[string]interface{}{
	{
		"email":      "khaimuldin.n@aitu.edu.kz",
		"password":   "Teacher123!",
		"role":       "teacher",
		"first_name": "Nursultan",
		"last_name":  "Khaimuldin",
		"department": "Software Engineering",
	},
	{
		"email":      "nurgaliyev.k@aitu.edu.kz",
		"password":   "Teacher123!",
		"role":       "teacher",
		"first_name": "Kenzhegali",
		"last_name":  "Nurgaliyev",
		"department": "Computer Science",
	},
	{
		"email":      "kalzhan.b@aitu.edu.kz",
		"password":   "Teacher123!",
		"role":       "teacher",
		"first_name": "Bakhtiyar",
		"last_name":  "Kalzhan",
		"department": "Mathematics",
	},
	{
		"email":      "alkhabay.b@aitu.edu.kz",
		"password":   "Teacher123!",
		"role":       "teacher",
		"first_name": "Bakgeldi",
		"last_name":  "Alkhabay",
		"department": "Computer Engineering",
	},
	{
		"email":      "tankeyev.s@aitu.edu.kz",
		"password":   "Teacher123!",
		"role":       "teacher",
		"first_name": "Samat",
		"last_name":  "Tankeyev",
		"department": "Software Engineering",
	},
	{
		"email":      "akhmetvalieva.i@aitu.edu.kz",
		"password":   "Teacher123!",
		"role":       "teacher",
		"first_name": "Irina",
		"last_name":  "Akhmetvalieva",
		"department": "Languages",
	},
}

// Students in SE-2430 (21 students)
var univStudents = []map[string]interface{}{
	{"email": "nurmukhammed.abdikarim@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Nurmukhammed", "last_name": "Abdikarim", "group_name": "SE-2430"},
	{"email": "dauren.akhmetov@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Dauren", "last_name": "Akhmetov", "group_name": "SE-2430"},
	{"email": "abylaikhan.alpamys@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Abylaikhan", "last_name": "Alpamys", "group_name": "SE-2430"},
	{"email": "madias.bek@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Madias", "last_name": "Bek", "group_name": "SE-2430"},
	{"email": "ulykbek.dovutbekov@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Ulykbek", "last_name": "Dovutbekov", "group_name": "SE-2430"},
	{"email": "abylaikhan.janmolda@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Abylaikhan", "last_name": "Janmolda", "group_name": "SE-2430"},
	{"email": "daryn.kaber@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Daryn", "last_name": "Kaber", "group_name": "SE-2430"},
	{"email": "zarina.kerimbay@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Zarina", "last_name": "Kerimbay", "group_name": "SE-2430"},
	{"email": "samat.kosmurat@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Samat", "last_name": "Kosmurat", "group_name": "SE-2430"},
	{"email": "madi.kuanay@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Madi", "last_name": "Kuanay", "group_name": "SE-2430"},
	{"email": "balaussa.kulmakhanbet@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Balaussa", "last_name": "Kulmakhanbet", "group_name": "SE-2430"},
	{"email": "artur.kupzhassarov@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Artur", "last_name": "Kupzhassarov", "group_name": "SE-2430"},
	{"email": "aisultan.nuriman@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Aisultan", "last_name": "Nuriman", "group_name": "SE-2430"},
	{"email": "bauyrzhan.nurzhanov@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Bauyrzhan", "last_name": "Nurzhanov", "group_name": "SE-2430"},
	{"email": "madina.rakhmetulla@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Madina", "last_name": "Rakhmetulla", "group_name": "SE-2430"},
	{"email": "ernar.sadenov@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Ernar", "last_name": "Sadenov", "group_name": "SE-2430"},
	{"email": "anton.shaubert@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Anton", "last_name": "Shaubert", "group_name": "SE-2430"},
	{"email": "yerkebulan.sovet@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Yerkebulan", "last_name": "Sovet", "group_name": "SE-2430"},
	{"email": "ali.turaliyev@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Ali", "last_name": "Turaliyev", "group_name": "SE-2430"},
	{"email": "shugyla.turganbek@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Shugyla", "last_name": "Turganbek", "group_name": "SE-2430"},
	{"email": "kamilla.yensebay@aitu.edu.kz", "password": "Student123!", "role": "student", "first_name": "Kamilla", "last_name": "Yensebay", "group_name": "SE-2430"},
}

// Courses (6 subjects for SE-2430)
type CourseConfig struct {
	Title          string
	Description    string
	TeacherEmail   string
	WeeklyContent  []WeekContent
}

type WeekContent struct {
	Week         int
	LectureTitle string
	LectureDesc  string
	HasQuiz      bool
	QuizTitle    string
	Assignment   *AssignmentConfig
}

type AssignmentConfig struct {
	Title       string
	Description string
	MaxPoints   int
	DaysFromNow int // Deadline in days from now
}

var univCourses = []CourseConfig{
	{
		Title:        "Advanced Programming 1",
		Description:  "Advanced programming concepts in Java/Python. Topics include OOP, design patterns, data structures, algorithms, and software development best practices.",
		TeacherEmail: "khaimuldin.n@aitu.edu.kz",
		WeeklyContent: []WeekContent{
			{Week: 1, LectureTitle: "Introduction to Advanced Programming", LectureDesc: "Course overview, development environment setup, review of basic programming concepts", HasQuiz: true, QuizTitle: "Basics Review Quiz"},
			{Week: 2, LectureTitle: "Object-Oriented Programming Deep Dive", LectureDesc: "Classes, objects, inheritance, polymorphism, encapsulation", HasQuiz: true, QuizTitle: "OOP Concepts Quiz", Assignment: &AssignmentConfig{Title: "Assignment 1", Description: "Provide solutions to the tasks at the end of Chapter 2.\n\nDEADLINE: End of Week 3. Late submissions will not be accepted!\n\nSubmit only one file in one of the following format:\n• .py format where each problem is neatly separated (e.g. # ----Problem 1-------, #----Problem 2-------, etc.)\n• .ipynb format where each cell is dedicated to one problem", MaxPoints: 100, DaysFromNow: 14}},
			{Week: 3, LectureTitle: "Design Patterns - Creational", LectureDesc: "Singleton, Factory, Abstract Factory, Builder patterns", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 2", Description: "Provide solutions to the tasks at the end of Chapter 3.\n\nImplement the design patterns covered in lectures.", MaxPoints: 100, DaysFromNow: 21}},
			{Week: 4, LectureTitle: "Design Patterns - Structural", LectureDesc: "Adapter, Bridge, Composite, Decorator patterns", HasQuiz: true, QuizTitle: "Design Patterns Quiz 1"},
			{Week: 5, LectureTitle: "Design Patterns - Behavioral", LectureDesc: "Observer, Strategy, Command, State patterns", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Midterm Project", Description: "Design and implement a small application using at least 3 design patterns.\n\nRequirements:\n• Clear documentation\n• UML diagrams\n• Working code with unit tests", MaxPoints: 200, DaysFromNow: 35}},
			{Week: 6, LectureTitle: "Data Structures - Advanced", LectureDesc: "Trees, graphs, hash tables implementation", HasQuiz: true, QuizTitle: "Data Structures Quiz"},
			{Week: 7, LectureTitle: "Algorithm Analysis", LectureDesc: "Time complexity, space complexity, Big O notation", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 3", Description: "Solve algorithm problems from LeetCode (provided list).\n\nSubmit solutions with complexity analysis.", MaxPoints: 100, DaysFromNow: 49}},
			{Week: 8, LectureTitle: "Sorting and Searching Algorithms", LectureDesc: "QuickSort, MergeSort, Binary Search variations", HasQuiz: true, QuizTitle: "Algorithms Quiz"},
			{Week: 9, LectureTitle: "Graph Algorithms", LectureDesc: "BFS, DFS, Dijkstra, shortest path algorithms", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 4", Description: "Implement graph algorithms and solve path-finding problems.", MaxPoints: 100, DaysFromNow: 63}},
			{Week: 10, LectureTitle: "Final Review & Exam Preparation", LectureDesc: "Course summary, exam preparation, Q&A session", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Final Project", Description: "Comprehensive project demonstrating all concepts learned.\n\nDeadline: Final exam week.", MaxPoints: 300, DaysFromNow: 70}},
		},
	},
	{
		Title:        "Storage Systems",
		Description:  "Database management systems, SQL, NoSQL, data modeling, and storage architecture.",
		TeacherEmail: "nurgaliyev.k@aitu.edu.kz",
		WeeklyContent: []WeekContent{
			{Week: 1, LectureTitle: "Introduction to Storage Systems", LectureDesc: "Overview of data storage, file systems, database concepts", HasQuiz: true, QuizTitle: "Storage Basics Quiz"},
			{Week: 2, LectureTitle: "Relational Databases", LectureDesc: "RDBMS concepts, normalization, ER diagrams", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Lab 1: ER Diagrams", Description: "Create ER diagrams for the given scenarios.\n\nSubmit as PDF.", MaxPoints: 50, DaysFromNow: 14}},
			{Week: 3, LectureTitle: "SQL Fundamentals", LectureDesc: "SELECT, INSERT, UPDATE, DELETE, JOINs", HasQuiz: true, QuizTitle: "SQL Basics Quiz", Assignment: &AssignmentConfig{Title: "Assignment 1: SQL Queries", Description: "Write SQL queries for the provided database schema.", MaxPoints: 100, DaysFromNow: 21}},
			{Week: 4, LectureTitle: "Advanced SQL", LectureDesc: "Subqueries, views, stored procedures, triggers", HasQuiz: false},
			{Week: 5, LectureTitle: "Database Design", LectureDesc: "Normalization (1NF-BCNF), denormalization, indexing", HasQuiz: true, QuizTitle: "Database Design Quiz", Assignment: &AssignmentConfig{Title: "Midterm Project", Description: "Design a complete database system for a real-world scenario.", MaxPoints: 200, DaysFromNow: 35}},
			{Week: 6, LectureTitle: "NoSQL Databases", LectureDesc: "MongoDB, document stores, key-value stores", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Lab 2: MongoDB", Description: "Implement CRUD operations using MongoDB.", MaxPoints: 50, DaysFromNow: 42}},
			{Week: 7, LectureTitle: "Distributed Storage", LectureDesc: "Replication, sharding, CAP theorem", HasQuiz: true, QuizTitle: "Distributed Systems Quiz"},
			{Week: 8, LectureTitle: "Cloud Storage Solutions", LectureDesc: "AWS S3, Azure Blob, Google Cloud Storage", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 2: Cloud Storage", Description: "Implement file storage using cloud services.", MaxPoints: 100, DaysFromNow: 56}},
			{Week: 9, LectureTitle: "Data Backup and Recovery", LectureDesc: "Backup strategies, disaster recovery, data integrity", HasQuiz: false},
			{Week: 10, LectureTitle: "Storage Security & Final Review", LectureDesc: "Encryption, access control, course summary", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Final Project", Description: "Build a complete storage solution with security features.", MaxPoints: 250, DaysFromNow: 70}},
		},
	},
	{
		Title:        "Computational Mathematics",
		Description:  "Numerical methods, linear algebra, calculus applications in computing, optimization algorithms.",
		TeacherEmail: "kalzhan.b@aitu.edu.kz",
		WeeklyContent: []WeekContent{
			{Week: 1, LectureTitle: "Mathematical Foundations", LectureDesc: "Review of calculus, linear algebra basics", HasQuiz: true, QuizTitle: "Math Foundations Quiz"},
			{Week: 2, LectureTitle: "Approximation Methods", LectureDesc: "Taylor series, interpolation, curve fitting", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 1", Description: "Provide solutions to the tasks at the end of Chapter 2.\n\nDEADLINE: End of Week 3, 23:59. Late submissions will not be accepted!\n\nSubmit only one file in one of the following format:\n• .py format where each problem is neatly separated\n• .ipynb format where each cell is dedicated to one problem", MaxPoints: 100, DaysFromNow: 14}},
			{Week: 3, LectureTitle: "Numerical Differentiation", LectureDesc: "Finite differences, error analysis", HasQuiz: true, QuizTitle: "Differentiation Quiz", Assignment: &AssignmentConfig{Title: "Assignment 2", Description: "Implement numerical differentiation algorithms.", MaxPoints: 100, DaysFromNow: 21}},
			{Week: 4, LectureTitle: "Numerical Integration", LectureDesc: "Trapezoidal rule, Simpson's rule, Gaussian quadrature", HasQuiz: false},
			{Week: 5, LectureTitle: "Root Finding Methods", LectureDesc: "Bisection, Newton-Raphson, secant method", HasQuiz: true, QuizTitle: "Root Finding Quiz", Assignment: &AssignmentConfig{Title: "Midterm Exam", Description: "Comprehensive exam covering weeks 1-5.\n\nLocation: Room 301\nTime: 10:00 - 12:00", MaxPoints: 200, DaysFromNow: 35}},
			{Week: 6, LectureTitle: "Linear Systems", LectureDesc: "Gaussian elimination, LU decomposition, iterative methods", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 3", Description: "Solve systems of linear equations using different methods.", MaxPoints: 100, DaysFromNow: 42}},
			{Week: 7, LectureTitle: "Eigenvalue Problems", LectureDesc: "Power method, QR algorithm", HasQuiz: true, QuizTitle: "Linear Algebra Quiz"},
			{Week: 8, LectureTitle: "Optimization Methods", LectureDesc: "Gradient descent, Newton's method for optimization", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 4", Description: "Implement optimization algorithms for given functions.", MaxPoints: 100, DaysFromNow: 56}},
			{Week: 9, LectureTitle: "Differential Equations", LectureDesc: "Euler method, Runge-Kutta methods", HasQuiz: false},
			{Week: 10, LectureTitle: "Final Review", LectureDesc: "Course summary, final exam preparation", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Final Exam", Description: "Comprehensive final examination.\n\nCovers all topics from the course.", MaxPoints: 300, DaysFromNow: 70}},
		},
	},
	{
		Title:        "Computer Organisation and Architecture",
		Description:  "Computer hardware organization, CPU architecture, memory hierarchy, I/O systems, and assembly programming.",
		TeacherEmail: "alkhabay.b@aitu.edu.kz",
		WeeklyContent: []WeekContent{
			{Week: 1, LectureTitle: "Introduction to Computer Architecture", LectureDesc: "History of computing, von Neumann architecture", HasQuiz: true, QuizTitle: "Computer History Quiz"},
			{Week: 2, LectureTitle: "Digital Logic Fundamentals", LectureDesc: "Boolean algebra, logic gates, combinational circuits", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Lab 1: Logic Gates", Description: "Design and simulate basic logic circuits using Logisim.", MaxPoints: 50, DaysFromNow: 14}},
			{Week: 3, LectureTitle: "Sequential Circuits", LectureDesc: "Flip-flops, registers, counters", HasQuiz: true, QuizTitle: "Digital Logic Quiz", Assignment: &AssignmentConfig{Title: "Assignment 1", Description: "Design sequential circuits for given specifications.", MaxPoints: 100, DaysFromNow: 21}},
			{Week: 4, LectureTitle: "CPU Architecture", LectureDesc: "ALU, control unit, datapath design", HasQuiz: false},
			{Week: 5, LectureTitle: "Instruction Set Architecture", LectureDesc: "RISC vs CISC, instruction formats, addressing modes", HasQuiz: true, QuizTitle: "ISA Quiz", Assignment: &AssignmentConfig{Title: "Midterm Exam", Description: "Written examination on computer architecture fundamentals.", MaxPoints: 200, DaysFromNow: 35}},
			{Week: 6, LectureTitle: "Assembly Programming", LectureDesc: "x86/ARM assembly basics, programming exercises", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Lab 2: Assembly", Description: "Write assembly programs for basic algorithms.", MaxPoints: 100, DaysFromNow: 42}},
			{Week: 7, LectureTitle: "Memory Hierarchy", LectureDesc: "Cache memory, virtual memory, memory management", HasQuiz: true, QuizTitle: "Memory Systems Quiz"},
			{Week: 8, LectureTitle: "Pipelining", LectureDesc: "Instruction pipelining, hazards, branch prediction", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 2", Description: "Analyze pipeline performance for given instruction sequences.", MaxPoints: 100, DaysFromNow: 56}},
			{Week: 9, LectureTitle: "I/O Systems", LectureDesc: "I/O interfaces, interrupts, DMA", HasQuiz: false},
			{Week: 10, LectureTitle: "Modern Architectures & Final Review", LectureDesc: "Multi-core, GPU computing, course summary", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Final Exam", Description: "Comprehensive final examination.", MaxPoints: 300, DaysFromNow: 70}},
		},
	},
	{
		Title:        "WEB Technologies 2 (Back End)",
		Description:  "Server-side web development, REST APIs, databases, authentication, and deployment.",
		TeacherEmail: "tankeyev.s@aitu.edu.kz",
		WeeklyContent: []WeekContent{
			{Week: 1, LectureTitle: "Introduction to Backend Development", LectureDesc: "Client-server architecture, HTTP protocol, REST principles", HasQuiz: true, QuizTitle: "Backend Basics Quiz"},
			{Week: 2, LectureTitle: "Node.js/Go Fundamentals", LectureDesc: "Setting up development environment, basic server creation", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Lab 1: Hello Server", Description: "Create a basic HTTP server that responds to requests.", MaxPoints: 50, DaysFromNow: 14}},
			{Week: 3, LectureTitle: "REST API Design", LectureDesc: "RESTful conventions, API versioning, documentation", HasQuiz: true, QuizTitle: "REST API Quiz", Assignment: &AssignmentConfig{Title: "Assignment 1: CRUD API", Description: "Build a complete CRUD API for a resource of your choice.", MaxPoints: 100, DaysFromNow: 21}},
			{Week: 4, LectureTitle: "Database Integration", LectureDesc: "ORM/query builders, PostgreSQL integration", HasQuiz: false},
			{Week: 5, LectureTitle: "Authentication & Authorization", LectureDesc: "JWT, OAuth, session management, RBAC", HasQuiz: true, QuizTitle: "Security Quiz", Assignment: &AssignmentConfig{Title: "Midterm Project", Description: "Build an authenticated API with user management.", MaxPoints: 200, DaysFromNow: 35}},
			{Week: 6, LectureTitle: "File Upload & Storage", LectureDesc: "Handling file uploads, cloud storage integration", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Lab 2: File Upload", Description: "Implement file upload functionality in your API.", MaxPoints: 50, DaysFromNow: 42}},
			{Week: 7, LectureTitle: "Real-time Communication", LectureDesc: "WebSockets, Server-Sent Events, chat implementation", HasQuiz: true, QuizTitle: "Real-time Quiz"},
			{Week: 8, LectureTitle: "Testing & Documentation", LectureDesc: "Unit testing, integration testing, API documentation", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Assignment 2", Description: "Write tests for your API and create OpenAPI documentation.", MaxPoints: 100, DaysFromNow: 56}},
			{Week: 9, LectureTitle: "Deployment & DevOps", LectureDesc: "Docker, CI/CD, cloud deployment", HasQuiz: false},
			{Week: 10, LectureTitle: "Final Project Presentations", LectureDesc: "Student project demonstrations, course wrap-up", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Final Project", Description: "Build a complete backend application with:\n• User authentication\n• CRUD operations\n• File upload\n• Real-time features\n• Deployed to cloud", MaxPoints: 300, DaysFromNow: 70}},
		},
	},
	{
		Title:        "Russian Language 2 (B1)",
		Description:  "Intermediate Russian language course focusing on grammar, vocabulary, reading, and conversation skills.",
		TeacherEmail: "akhmetvalieva.i@aitu.edu.kz",
		WeeklyContent: []WeekContent{
			{Week: 1, LectureTitle: "Повторение (Review)", LectureDesc: "Review of A2 level grammar and vocabulary", HasQuiz: true, QuizTitle: "Placement Quiz"},
			{Week: 2, LectureTitle: "Падежи (Cases) - Часть 1", LectureDesc: "Nominative, genitive, dative cases review and practice", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Домашнее задание 1", Description: "Complete exercises on Russian cases.\n\nWorkbook pages 15-20.", MaxPoints: 50, DaysFromNow: 14}},
			{Week: 3, LectureTitle: "Падежи (Cases) - Часть 2", LectureDesc: "Accusative, instrumental, prepositional cases", HasQuiz: true, QuizTitle: "Cases Quiz", Assignment: &AssignmentConfig{Title: "Домашнее задание 2", Description: "Grammar exercises and short essay (150 words).", MaxPoints: 50, DaysFromNow: 21}},
			{Week: 4, LectureTitle: "Глаголы движения", LectureDesc: "Verbs of motion: идти/ходить, ехать/ездить", HasQuiz: false},
			{Week: 5, LectureTitle: "Виды глагола", LectureDesc: "Verbal aspects: perfective and imperfective", HasQuiz: true, QuizTitle: "Verbs Quiz", Assignment: &AssignmentConfig{Title: "Midterm Test", Description: "Written test covering weeks 1-5.\n\nGrammar, vocabulary, reading comprehension.", MaxPoints: 100, DaysFromNow: 35}},
			{Week: 6, LectureTitle: "Причастия и деепричастия", LectureDesc: "Participles and verbal adverbs", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Домашнее задание 3", Description: "Exercises on participles.\nReading: short story analysis.", MaxPoints: 50, DaysFromNow: 42}},
			{Week: 7, LectureTitle: "Условные предложения", LectureDesc: "Conditional sentences, subjunctive mood", HasQuiz: true, QuizTitle: "Grammar Quiz 2"},
			{Week: 8, LectureTitle: "Чтение и аудирование", LectureDesc: "Reading comprehension and listening practice", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Домашнее задание 4", Description: "Listening comprehension exercises.\nWritten response to audio materials.", MaxPoints: 50, DaysFromNow: 56}},
			{Week: 9, LectureTitle: "Разговорная практика", LectureDesc: "Conversation practice, role-playing scenarios", HasQuiz: false},
			{Week: 10, LectureTitle: "Итоговый экзамен", LectureDesc: "Final exam preparation and review", HasQuiz: false, Assignment: &AssignmentConfig{Title: "Final Exam", Description: "Comprehensive exam:\n• Grammar (40%)\n• Reading (20%)\n• Listening (20%)\n• Speaking (20%)", MaxPoints: 200, DaysFromNow: 70}},
		},
	},
}

// Attendance statuses
var attendanceStatuses = []string{"present", "absent", "late", "excused"}

// Grade feedback templates
var univGradeFeedbacks = []string{
	"Excellent work! Shows deep understanding of the material.",
	"Good effort. Some minor issues but overall well done.",
	"Satisfactory. Please review the lecture notes for improvement.",
	"Needs improvement. Visit office hours for additional help.",
	"Well done! Keep up the great work.",
	"Good progress. Consider the feedback for future assignments.",
}

// ============================================================================
// University Seeder State
// ============================================================================

type UnivSeederState struct {
	AdminToken    string
	TeacherTokens map[string]string // email -> token
	StudentTokens map[string]string // email -> token
	TeacherIDs    map[string]string // email -> user_id
	StudentIDs    map[string]string // email -> user_id
	Courses       []UnivCourseResponse
	Assignments   []UnivAssignmentResponse
	Enrollments   []UnivEnrollmentResponse
}

func NewUnivSeederState() *UnivSeederState {
	return &UnivSeederState{
		TeacherTokens: make(map[string]string),
		StudentTokens: make(map[string]string),
		TeacherIDs:    make(map[string]string),
		StudentIDs:    make(map[string]string),
	}
}

// ============================================================================
// HTTP Helpers
// ============================================================================

func univDoRequest(method, endpoint string, body interface{}, token string) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, apiURL+endpoint, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := httpClient.Do(req)
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

func univPost(endpoint string, body interface{}, token string) ([]byte, error) {
	respBody, status, err := univDoRequest("POST", endpoint, body, token)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("POST %s failed with status %d: %s", endpoint, status, string(respBody))
	}
	return respBody, nil
}

func univPut(endpoint string, body interface{}, token string) ([]byte, error) {
	respBody, status, err := univDoRequest("PUT", endpoint, body, token)
	if err != nil {
		return nil, err
	}
	if status >= 400 {
		return nil, fmt.Errorf("PUT %s failed with status %d: %s", endpoint, status, string(respBody))
	}
	return respBody, nil
}

func univParseResponse(respBody []byte, target interface{}) error {
	var apiResp UnivAPIResponse
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
// Seeding Functions
// ============================================================================

func (s *UnivSeederState) registerAndLogin(user map[string]interface{}) (*UnivLoginResponse, error) {
	// Register
	_, err := univPost("/auth/register", user, "")
	if err != nil {
		log.Printf("    Registration note (might exist): %v", err)
	}

	// Login
	loginBody := map[string]interface{}{
		"email":    user["email"],
		"password": user["password"],
	}
	respBody, err := univPost("/auth/login", loginBody, "")
	if err != nil {
		return nil, fmt.Errorf("login failed for %s: %w", user["email"], err)
	}

	var loginResp UnivLoginResponse
	if err := univParseResponse(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("parse login: %w", err)
	}

	return &loginResp, nil
}

func (s *UnivSeederState) seedUsers() error {
	log.Println("\n📝 Seeding University Users...")
	log.Println("=" + string(make([]byte, 50)))

	// Admin
	adminResp, err := s.registerAndLogin(univAdmin)
	if err != nil {
		return fmt.Errorf("admin: %w", err)
	}
	s.AdminToken = adminResp.AccessToken
	log.Printf("  ✓ Admin: %s", univAdmin["email"])

	// Teachers
	log.Println("\n  📚 Teachers:")
	for _, teacher := range univTeachers {
		resp, err := s.registerAndLogin(teacher)
		if err != nil {
			log.Printf("    ✗ %s %s: %v", teacher["first_name"], teacher["last_name"], err)
			continue
		}
		email := teacher["email"].(string)
		s.TeacherTokens[email] = resp.AccessToken
		s.TeacherIDs[email] = resp.User.ID
		log.Printf("    ✓ %s %s (%s)", teacher["first_name"], teacher["last_name"], teacher["department"])
	}

	// Students
	log.Println("\n  👨‍🎓 Students (SE-2430):")
	for _, student := range univStudents {
		resp, err := s.registerAndLogin(student)
		if err != nil {
			log.Printf("    ✗ %s %s: %v", student["first_name"], student["last_name"], err)
			continue
		}
		email := student["email"].(string)
		s.StudentTokens[email] = resp.AccessToken
		s.StudentIDs[email] = resp.User.ID
		log.Printf("    ✓ %s %s", student["first_name"], student["last_name"])
	}

	return nil
}

func (s *UnivSeederState) seedCourses() error {
	log.Println("\n📚 Seeding Courses...")
	log.Println("=" + string(make([]byte, 50)))

	for _, courseConfig := range univCourses {
		token, ok := s.TeacherTokens[courseConfig.TeacherEmail]
		if !ok {
			log.Printf("  ✗ No token for teacher %s", courseConfig.TeacherEmail)
			continue
		}

		courseBody := map[string]interface{}{
			"title":       courseConfig.Title,
			"description": courseConfig.Description,
		}

		respBody, err := univPost("/courses", courseBody, token)
		if err != nil {
			log.Printf("  ✗ Course '%s': %v", courseConfig.Title, err)
			continue
		}

		var courseResp UnivCourseResponse
		if err := univParseResponse(respBody, &courseResp); err != nil {
			log.Printf("  ✗ Parse course: %v", err)
			continue
		}

		s.Courses = append(s.Courses, courseResp)

		// Find teacher name
		teacherName := ""
		for _, t := range univTeachers {
			if t["email"] == courseConfig.TeacherEmail {
				teacherName = fmt.Sprintf("%s %s", t["first_name"], t["last_name"])
				break
			}
		}
		log.Printf("  ✓ %s | %s", courseConfig.Title, teacherName)
	}

	return nil
}

func (s *UnivSeederState) seedEnrollments() error {
	log.Println("\n📋 Enrolling SE-2430 Students...")
	log.Println("=" + string(make([]byte, 50)))

	if len(s.Courses) == 0 {
		return fmt.Errorf("no courses available")
	}

	// Enroll all students in all courses
	for studentEmail, token := range s.StudentTokens {
		studentName := ""
		for _, st := range univStudents {
			if st["email"] == studentEmail {
				studentName = fmt.Sprintf("%s %s", st["first_name"], st["last_name"])
				break
			}
		}

		for _, course := range s.Courses {
			enrollBody := map[string]interface{}{
				"course_id": course.ID,
			}

			respBody, err := univPost("/enrollments", enrollBody, token)
			if err != nil {
				log.Printf("  ✗ %s -> %s: %v", studentName, course.Title, err)
				continue
			}

			var enrollResp UnivEnrollmentResponse
			if err := univParseResponse(respBody, &enrollResp); err != nil {
				continue
			}

			s.Enrollments = append(s.Enrollments, enrollResp)
		}
		log.Printf("  ✓ %s enrolled in %d courses", studentName, len(s.Courses))
	}

	// Approve all enrollments
	log.Println("\n  Approving enrollments...")
	for _, enrollment := range s.Enrollments {
		// Find teacher token for this course
		var teacherToken string
		for i, c := range s.Courses {
			if c.ID == enrollment.CourseID {
				teacherEmail := univCourses[i].TeacherEmail
				teacherToken = s.TeacherTokens[teacherEmail]
				break
			}
		}

		if teacherToken == "" {
			continue
		}

		statusBody := map[string]interface{}{
			"status": "active",
		}
		univPut("/enrollments/"+enrollment.ID+"/status", statusBody, teacherToken)
	}
	log.Printf("  ✓ Approved %d enrollments", len(s.Enrollments))

	return nil
}

func (s *UnivSeederState) seedAssignments() error {
	log.Println("\n📝 Seeding Assignments & Content...")
	log.Println("=" + string(make([]byte, 50)))

	for i, courseConfig := range univCourses {
		if i >= len(s.Courses) {
			break
		}
		course := s.Courses[i]
		token := s.TeacherTokens[courseConfig.TeacherEmail]

		log.Printf("\n  📖 %s:", course.Title)

		for _, week := range courseConfig.WeeklyContent {
			// Create assignment if exists for this week
			if week.Assignment != nil {
				dueDate := time.Now().AddDate(0, 0, week.Assignment.DaysFromNow)

				assignmentBody := map[string]interface{}{
					"course_id":   course.ID,
					"title":       fmt.Sprintf("Week %d: %s", week.Week, week.Assignment.Title),
					"description": week.Assignment.Description,
					"max_points":  week.Assignment.MaxPoints,
					"due_at":      dueDate.Format(time.RFC3339),
				}

				respBody, err := univPost("/assignments", assignmentBody, token)
				if err != nil {
					log.Printf("    ✗ Week %d Assignment: %v", week.Week, err)
					continue
				}

				var assignmentResp UnivAssignmentResponse
				if err := univParseResponse(respBody, &assignmentResp); err != nil {
					continue
				}

				s.Assignments = append(s.Assignments, assignmentResp)
				log.Printf("    ✓ Week %d: %s (Due: %s, %d pts)",
					week.Week, week.Assignment.Title,
					dueDate.Format("Jan 02"), week.Assignment.MaxPoints)
			}
		}
	}

	return nil
}

func (s *UnivSeederState) seedAttendance() error {
	log.Println("\n📅 Seeding Attendance Sessions...")
	log.Println("=" + string(make([]byte, 50)))

	// Create attendance sessions for past weeks
	for i, courseConfig := range univCourses {
		if i >= len(s.Courses) {
			break
		}
		course := s.Courses[i]
		token := s.TeacherTokens[courseConfig.TeacherEmail]

		log.Printf("  📖 %s:", course.Title)

		// Create sessions for weeks 1-3 (past weeks)
		for week := 1; week <= 3; week++ {
			sessionDate := time.Now().AddDate(0, 0, -(10-week)*7)

			sessionBody := map[string]interface{}{
				"course_id":    course.ID,
				"title":        fmt.Sprintf("Week %d - Lecture", week),
				"session_date": sessionDate.Format("2006-01-02"),
				"start_time":   sessionDate.Format(time.RFC3339),
			}

			respBody, err := univPost("/attendance/sessions", sessionBody, token)
			if err != nil {
				log.Printf("    ✗ Week %d session: %v", week, err)
				continue
			}

			var sessionResp struct {
				ID string `json:"id"`
			}
			if err := univParseResponse(respBody, &sessionResp); err != nil {
				continue
			}

			// Mark attendance for all students
			presentCount := 0
			for _, enrollment := range s.Enrollments {
				if enrollment.CourseID != course.ID {
					continue
				}

				// Random attendance status (80% present, 10% late, 5% absent, 5% excused)
				r := rand.Float32()
				var status string
				if r < 0.80 {
					status = "present"
					presentCount++
				} else if r < 0.90 {
					status = "late"
					presentCount++
				} else if r < 0.95 {
					status = "absent"
				} else {
					status = "excused"
				}

				markBody := map[string]interface{}{
					"session_id": sessionResp.ID,
					"student_id": enrollment.StudentID,
					"status":     status,
				}

				univPost("/attendance/marks", markBody, token)
			}
			log.Printf("    ✓ Week %d: %d/%d present", week, presentCount, 21)
		}
	}

	return nil
}

func (s *UnivSeederState) seedSubmissionsAndGrades() error {
	log.Println("\n📤 Seeding Submissions & Grades...")
	log.Println("=" + string(make([]byte, 50)))

	// Get course-student mapping from enrollments
	courseStudents := make(map[string][]string) // courseID -> []studentEmail
	for _, enrollment := range s.Enrollments {
		for email, id := range s.StudentIDs {
			if id == enrollment.StudentID {
				courseStudents[enrollment.CourseID] = append(courseStudents[enrollment.CourseID], email)
				break
			}
		}
	}

	submissionCount := 0
	gradedCount := 0

	for _, assignment := range s.Assignments {
		students := courseStudents[assignment.CourseID]
		if len(students) == 0 {
			continue
		}

		// Find teacher token for grading
		var teacherToken string
		for i, c := range s.Courses {
			if c.ID == assignment.CourseID {
				teacherToken = s.TeacherTokens[univCourses[i].TeacherEmail]
				break
			}
		}

		// 70-90% of students submit past assignments
		numSubmissions := len(students) * (70 + rand.Intn(20)) / 100
		if numSubmissions == 0 {
			numSubmissions = 1
		}

		for i := 0; i < numSubmissions && i < len(students); i++ {
			studentEmail := students[i]
			token := s.StudentTokens[studentEmail]

			submissionBody := map[string]interface{}{
				"assignment_id": assignment.ID,
				"content_text":  "Here is my submission for this assignment. I have completed all required tasks as specified in the assignment description.",
			}

			respBody, err := univPost("/submissions", submissionBody, token)
			if err != nil {
				continue
			}

			var submissionResp struct {
				ID string `json:"id"`
			}
			if err := univParseResponse(respBody, &submissionResp); err != nil {
				continue
			}
			submissionCount++

			// Grade 80% of submissions
			if rand.Float32() < 0.8 && teacherToken != "" {
				score := 60 + rand.Intn(40) // Score between 60-100%
				feedback := univGradeFeedbacks[rand.Intn(len(univGradeFeedbacks))]

				gradeBody := map[string]interface{}{
					"submission_id": submissionResp.ID,
					"score":         score,
					"feedback":      feedback,
				}

				_, err := univPost("/grades", gradeBody, teacherToken)
				if err == nil {
					gradedCount++
				}
			}
		}
	}

	log.Printf("  ✓ Created %d submissions", submissionCount)
	log.Printf("  ✓ Graded %d submissions", gradedCount)

	return nil
}

// ============================================================================
// Main - Run University Seeder
// ============================================================================

func RunUniversitySeeder() {
	log.Println("🎓 ====================================================")
	log.Println("   AITU LMS University Seeder")
	log.Println("   SE-2430 Software Engineering Group")
	log.Println("🎓 ====================================================")

	rand.Seed(time.Now().UnixNano())
	state := NewUnivSeederState()

	// Seed in order
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Users", state.seedUsers},
		{"Courses", state.seedCourses},
		{"Enrollments", state.seedEnrollments},
		{"Assignments", state.seedAssignments},
		{"Attendance", state.seedAttendance},
		{"Submissions & Grades", state.seedSubmissionsAndGrades},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			log.Printf("❌ %s failed: %v", step.name, err)
		}
	}

	// Summary
	log.Println("\n🎓 ====================================================")
	log.Println("   SEEDING COMPLETE!")
	log.Println("🎓 ====================================================")
	log.Printf("   📚 Teachers: %d", len(state.TeacherTokens))
	log.Printf("   👨‍🎓 Students: %d", len(state.StudentTokens))
	log.Printf("   📖 Courses: %d", len(state.Courses))
	log.Printf("   📝 Assignments: %d", len(state.Assignments))
	log.Printf("   📋 Enrollments: %d", len(state.Enrollments))
	log.Println("")
	log.Println("   Test Credentials:")
	log.Println("   ─────────────────────────────────────────────────")
	log.Println("   Admin:    admin@aitu.edu.kz / Admin123!@#")
	log.Println("   Teacher:  khaimuldin.n@aitu.edu.kz / Teacher123!")
	log.Println("   Student:  bauyrzhan.nurzhanov@aitu.edu.kz / Student123!")
	log.Println("🎓 ====================================================")
}
