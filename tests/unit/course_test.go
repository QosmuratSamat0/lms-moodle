package unit

import (
	"testing"

	"github.com/ap1-final-mini-moodle/internal/domain/course"
	courseUC "github.com/ap1-final-mini-moodle/internal/usecase/course"
)

type mockCourseRepo struct {
	courses map[string]*course.Course
}

func (m *mockCourseRepo) Create(c *course.Course) error {
	m.courses[c.ID] = c
	return nil
}

func (m *mockCourseRepo) GetByID(id string) (*course.Course, error) {
	if c, ok := m.courses[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *mockCourseRepo) GetByCode(code string) (*course.Course, error) {
	for _, c := range m.courses {
		if c.Code == code {
			return c, nil
		}
	}
	return nil, nil
}

func (m *mockCourseRepo) Update(c *course.Course) error {
	m.courses[c.ID] = c
	return nil
}

func (m *mockCourseRepo) List(skip, take int) ([]*course.Course, error) {
	return nil, nil
}

func (m *mockCourseRepo) ListByTeacher(teacherID string, skip, take int) ([]*course.Course, error) {
	return nil, nil
}

func (m *mockCourseRepo) Delete(id string) error {
	delete(m.courses, id)
	return nil
}

func TestCreateCourse(t *testing.T) {
	repo := &mockCourseRepo{courses: make(map[string]*course.Course)}
	service := courseUC.NewService(repo)

	input := &course.CreateCourseInput{
		Code:        "CS101",
		Title:       "Introduction to CS",
		Description: "Basic CS course",
		TeacherID:   "teacher1",
		MaxPoints:   100,
	}

	c, err := service.Create(input)
	if err != nil {
		t.Fatalf("failed to create course: %v", err)
	}

	if c.Code != input.Code {
		t.Errorf("expected code %s, got %s", input.Code, c.Code)
	}
	if c.MaxPoints != input.MaxPoints {
		t.Errorf("expected max points %d, got %d", input.MaxPoints, c.MaxPoints)
	}
}

func TestUpdateCourse(t *testing.T) {
	repo := &mockCourseRepo{courses: make(map[string]*course.Course)}
	service := courseUC.NewService(repo)

	input := &course.CreateCourseInput{
		Code:        "CS101",
		Title:       "Introduction to CS",
		Description: "Basic CS course",
		TeacherID:   "teacher1",
		MaxPoints:   100,
	}

	c, _ := service.Create(input)

	newTitle := "Advanced CS"
	updateInput := &course.UpdateCourseInput{
		Title: &newTitle,
	}

	updated, err := service.Update(c.ID, updateInput)
	if err != nil {
		t.Fatalf("failed to update course: %v", err)
	}

	if updated.Title != newTitle {
		t.Errorf("expected title %s, got %s", newTitle, updated.Title)
	}
}
