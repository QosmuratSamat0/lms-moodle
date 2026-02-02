package usecase

import "ap1-final-mini-moodle/internal/repository"

type AdminCourseView struct {
	ID           int64
	Title        string
	TeacherID    int64
	TeacherName  string
	StudentCount int
}

type AdminUsecase struct {
	store repository.Store
}

func NewAdminUsecase(store repository.Store) *AdminUsecase {
	return &AdminUsecase{store: store}
}

func (u *AdminUsecase) ListCourses() []AdminCourseView {
	courses := u.store.ListCourses()
	result := make([]AdminCourseView, 0, len(courses))
	for _, course := range courses {
		result = append(result, AdminCourseView{
			ID:           course.ID,
			Title:        course.Title,
			TeacherID:    course.TeacherID,
			TeacherName:  course.TeacherName,
			StudentCount: u.store.CountStudentsInCourse(course.ID),
		})
	}
	return result
}
