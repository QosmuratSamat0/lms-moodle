package repository

import (
	"sync"
	"time"

	"ap1-final-mini-moodle/internal/entity"
)

type InMemoryStore struct {
	mu             sync.RWMutex
	nextCourseID   int64
	nextAssignment int64
	nextSubmission int64
	courses        map[int64]entity.Course
	assignments    map[int64]entity.Assignment
	submissions    map[int64]entity.Submission
	grades         map[int64]entity.Grade
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		nextCourseID:   100,
		nextAssignment: 500,
		nextSubmission: 9000,
		courses:        make(map[int64]entity.Course),
		assignments:    make(map[int64]entity.Assignment),
		submissions:    make(map[int64]entity.Submission),
		grades:         make(map[int64]entity.Grade),
	}
}

func (s *InMemoryStore) CreateCourse(course entity.Course) (entity.Course, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextCourseID++
	course.ID = s.nextCourseID
	if course.CreatedAt.IsZero() {
		course.CreatedAt = time.Now().UTC()
	}
	s.courses[course.ID] = course
	return course, nil
}

func (s *InMemoryStore) GetCourse(id int64) (entity.Course, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	course, ok := s.courses[id]
	return course, ok
}

func (s *InMemoryStore) ListCourses() []entity.Course {
	s.mu.RLock()
	defer s.mu.RUnlock()

	courses := make([]entity.Course, 0, len(s.courses))
	for _, course := range s.courses {
		courses = append(courses, course)
	}
	return courses
}

func (s *InMemoryStore) CreateAssignment(assignment entity.Assignment) (entity.Assignment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextAssignment++
	assignment.ID = s.nextAssignment
	if assignment.CreatedAt.IsZero() {
		assignment.CreatedAt = time.Now().UTC()
	}
	s.assignments[assignment.ID] = assignment
	return assignment, nil
}

func (s *InMemoryStore) GetAssignment(id int64) (entity.Assignment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	assignment, ok := s.assignments[id]
	return assignment, ok
}

func (s *InMemoryStore) CreateSubmission(submission entity.Submission) (entity.Submission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextSubmission++
	submission.ID = s.nextSubmission
	if submission.SubmittedAt.IsZero() {
		submission.SubmittedAt = time.Now().UTC()
	}
	s.submissions[submission.ID] = submission
	return submission, nil
}

func (s *InMemoryStore) GetSubmission(id int64) (entity.Submission, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	submission, ok := s.submissions[id]
	return submission, ok
}

func (s *InMemoryStore) ListSubmissionsByAssignment(assignmentID int64) []entity.Submission {
	s.mu.RLock()
	defer s.mu.RUnlock()

	submissions := make([]entity.Submission, 0)
	for _, submission := range s.submissions {
		if submission.AssignmentID == assignmentID {
			submissions = append(submissions, submission)
		}
	}
	return submissions
}

func (s *InMemoryStore) UpsertGrade(grade entity.Grade) (entity.Grade, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if grade.GradedAt.IsZero() {
		grade.GradedAt = time.Now().UTC()
	}
	s.grades[grade.SubmissionID] = grade
	return grade, nil
}

func (s *InMemoryStore) GetGrade(submissionID int64) (entity.Grade, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	grade, ok := s.grades[submissionID]
	return grade, ok
}

func (s *InMemoryStore) CountStudentsInCourse(courseID int64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	assignmentIDs := make(map[int64]struct{})
	for _, assignment := range s.assignments {
		if assignment.CourseID == courseID {
			assignmentIDs[assignment.ID] = struct{}{}
		}
	}

	students := make(map[int64]struct{})
	for _, submission := range s.submissions {
		if _, ok := assignmentIDs[submission.AssignmentID]; ok {
			students[submission.StudentID] = struct{}{}
		}
	}

	return len(students)
}
