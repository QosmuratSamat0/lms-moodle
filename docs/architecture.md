# Mini-Moodl Architecture & Design (Monolith)

## 1. Architecture Overview
Mini-Moodle is built as a **monolithic Go backend** with a clear layered structure. The system exposes a REST API used by clients (web/mobile). Data is stored in **PostgreSQL**. The monolith approach is chosen for the milestone because it simplifies development, deployment, and team coordination while keeping module boundaries clear and scalable for future evolution.

**High-level components:**
- **Client apps**: Web UI / Mobile UI (not fully defined in this milestone)
- **Go Backend (Monolith)**: HTTP API + business logic + data access
- **PostgreSQL**: persistent storage for users, courses, assignments, submissions, grades

## 2. Key Roles and Access Model
The system supports three roles:
- **Student**: enroll in courses, view assignments, submit work, view grades/feedback
- **Teacher**: create courses, create assignments, view submissions, grade and give feedback
- **Admin**: manage users/roles and oversee course-related administration tasks

Access control is enforced on the backend using **role-based authorization** (RBAC). Endpoints are protected via middleware that verifies authentication and checks the required role(s).

## 3. Backend Structure (Layered Design)
The backend is organized into modules by domain. Each domain follows the same internal structure:

- **Handler (HTTP layer)**
  Parses request, validates input (DTO), calls service, returns HTTP response.

- **DTO (Data Transfer Objects)**
  Request/response schemas for API communication. Used for validation and to avoid exposing internal models directly.

- **Service (Business logic layer)**
  Implements core use-cases: rules, permissions, workflow decisions. Service does not talk directly to HTTP or SQL.

- **Repository (Data access layer)**
  Handles database queries and persistence. Service calls repository interfaces to read/write data.

- **Model (Domain entities)**
  Core business objects (e.g., User, Course, Assignment, Submission, Grade). Models are close to database structure but can evolve.

This pattern ensures separation of concerns and makes the project easier to test and maintain.

## 4. Core Modules and Responsibilities
### 4.1 User Module
**Responsibilities:**
- user entity and role management
- authentication support (password hash storage)
- admin-level operations for user creation/role assignment

**Main entities:** `users`

### 4.2 Course Module
**Responsibilities:**
- course creation (teacher)
- course listing and viewing (student/teacher/admin)
- linking course to teacher owner

**Main entities:** `courses`

### 4.3 Enrollment Module
**Responsibilities:**
- student enrollment into courses
- preventing duplicate enrollments
- ensuring only students can enroll

**Main entities:** `enrollments`

### 4.4 Assignment Module
**Responsibilities:**
- teacher creates assignments within a course
- listing assignments by course
- deadlines and max points logic

**Main entities:** `assignments`

### 4.5 Submission Module
**Responsibilities:**
- student submits assignment (one or multiple attempts depending on rules)
- store submission metadata (time, file path / content reference)
- list submissions per assignment for teacher review

**Main entities:** `submissions`

### 4.6 Grade Module
**Responsibilities:**
- teacher grades a submission and provides feedback
- student views grade and feedback
- prevent non-teachers from grading

**Main entities:** `grades`

## 5. Data Flow (Request Lifecycle)
The typical request flow is:

1. **Client** sends HTTP request to backend
2. **Router** matches route → calls a **Handler**
3. **Middleware** verifies authentication (and role if required)
4. **Handler** parses request → maps to **DTO** → validates DTO
5. **Service** executes business logic and permissions
6. **Repository** performs database operations (SQL queries)
7. **Service** returns result
8. **Handler** formats response DTO and returns JSON response

This flow is consistent across all modules.

## 6. Example Workflows
### 6.1 Student submits assignment
1. Student authenticates
2. Student enrolls into a course
3. Student requests assignments list for the course
4. Student submits assignment → `SubmissionService` checks:
   - student is enrolled
   - assignment exists and deadline rules (if enabled)
5. Repository stores submission in `submissions`
6. Student can later view grade/feedback after teacher grading

### 6.2 Teacher grades submission
1. Teacher authenticates
2. Teacher views submissions for an assignment
3. Teacher grades a submission → `GradeService` checks:
   - teacher owns the course (or has permission)
   - submission exists
4. Repository creates/updates `grades` record
5. Student views grade and feedback

## 7. Database Design Summary (ERD Alignment)
The ERD is designed around the core academic workflow:

- `users` stores all accounts with role (student/teacher/admin)
- `courses` belongs to a teacher (`teacher_id → users.id`)
- `enrollments` links students to courses (many-to-many)
- `assignments` belongs to a course
- `submissions` belongs to assignment and student
- `grades` belongs to submission and teacher

This schema supports:
- course ownership by teachers
- enrollment control for students
- assignment publishing
- submission tracking
- grading and feedback

## 8. Non-Functional Considerations (Milestone Level)
- **Maintainability:** clear module boundaries and repeated structure
- **Scalability:** monolith now, but modules can later be extracted if needed
- **Security:** password hashes only (no plain passwords), RBAC middleware
- **Reliability:** health endpoint and basic error handling
- **Extensibility:** chat/notifications can be added as future modules

## 9. Team Work Split (Alignment with Plan)
- **Bauyrzhan**: Student logic + half Admin logic + backend integration
- **Samat**: Teacher logic + half Admin logic + documentation/diagrams
Both contribute to authentication, integration, and final build checks.

---
**Diagrams:** See `docs/diagrams/` for Architecture diagram, Use-Case diagram, ERD, and UML diagram.
