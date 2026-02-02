# Mini-Moodle API (MVP Draft)

This document describes the **planned** MVP API for teacher workflows and admin monitoring.
Current repository does not yet include runnable HTTP handlers, so these endpoints are a
contract for implementation and manual review.

## Conventions
- Base path: `/api`
- Auth: `Authorization: Bearer <token>`
- Roles: `teacher`, `admin`
- Dates: ISO-8601 (`2026-02-02T12:30:00Z`)

---

## Teacher-side API

### 1) Create course
**Endpoint:** `POST /api/teacher/courses`  
**Role:** teacher

**Request body**
```json
{
  "title": "Intro to Databases",
  "description": "Basics of relational DBs"
}
```

**Response**
```json
{
  "id": 101,
  "title": "Intro to Databases",
  "description": "Basics of relational DBs",
  "teacher_id": 7,
  "created_at": "2026-02-02T12:30:00Z"
}
```

---

### 2) Create assignment
**Endpoint:** `POST /api/teacher/courses/{courseId}/assignments`  
**Role:** teacher (must own the course)

**Request body**
```json
{
  "title": "Homework 1",
  "description": "ERD + normalization",
  "due_at": "2026-02-10T23:59:59Z",
  "max_points": 100
}
```

**Response**
```json
{
  "id": 501,
  "course_id": 101,
  "title": "Homework 1",
  "description": "ERD + normalization",
  "due_at": "2026-02-10T23:59:59Z",
  "max_points": 100,
  "created_at": "2026-02-02T12:45:00Z"
}
```

---

### 3) View submissions for an assignment
**Endpoint:** `GET /api/teacher/assignments/{assignmentId}/submissions`  
**Role:** teacher (must own the course)

**Response**
```json
[
  {
    "id": 9001,
    "assignment_id": 501,
    "student_id": 55,
    "submitted_at": "2026-02-09T20:01:00Z",
    "content_ref": "s3://submissions/9001.pdf",
    "grade": null
  },
  {
    "id": 9002,
    "assignment_id": 501,
    "student_id": 56,
    "submitted_at": "2026-02-09T21:15:00Z",
    "content_ref": "s3://submissions/9002.pdf",
    "grade": {
      "score": 87,
      "feedback": "Good structure, fix 2NF on table X"
    }
  }
]
```

---

### 4) Grade a submission (score + feedback)
**Endpoint:** `POST /api/teacher/submissions/{submissionId}/grade`  
**Role:** teacher (must own the course)

**Request body**
```json
{
  "score": 87,
  "feedback": "Good structure, fix 2NF on table X"
}
```

**Response**
```json
{
  "submission_id": 9002,
  "graded_by": 7,
  "score": 87,
  "feedback": "Good structure, fix 2NF on table X",
  "graded_at": "2026-02-10T08:10:00Z"
}
```

---

## Teacher flow (words)
1. Teacher logs in and calls **Create course**.
2. Teacher adds assignments to that course via **Create assignment**.
3. When students submit, teacher loads **View submissions** for the assignment.
4. Teacher assigns **grade + feedback** per submission.

---

## Admin monitoring API (MVP)

### List courses with teacher and student count
**Endpoint:** `GET /api/admin/courses`  
**Role:** admin

**Response**
```json
[
  {
    "id": 101,
    "title": "Intro to Databases",
    "teacher": {
      "id": 7,
      "name": "Dr. Smith"
    },
    "student_count": 28
  },
  {
    "id": 102,
    "title": "Algorithms",
    "teacher": {
      "id": 9,
      "name": "Prof. Lee"
    },
    "student_count": 32
  }
]
```

