//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/MaqsattoTeam/aLMS/golang-service/internal/domain/course"
	"github.com/MaqsattoTeam/aLMS/golang-service/tests/testutil"
	"github.com/google/uuid"
)

func TestCourseRepository_Create(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Seed required data (teacher)
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestTeacherID, "teacher@test.com", "$2a$10$hash", "teacher")
	fixtures.CreateTeacher(ctx, t, testutil.TestTeacherID, "Test", "Teacher", "CS")

	repo := course.NewRepository(pool)

	tests := []struct {
		name    string
		course  *course.Course
		wantErr bool
	}{
		{
			name: "create course with all fields",
			course: &course.Course{
				Title:          "Test Course",
				Description:    stringPtr("A test course description"),
				OwnerTeacherID: &testutil.TestTeacherID,
				IsActive:       true,
			},
			wantErr: false,
		},
		{
			name: "create course with minimal fields",
			course: &course.Course{
				Title:    "Minimal Course",
				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "create course with nil description",
			course: &course.Course{
				Title:          "No Description Course",
				Description:    nil,
				OwnerTeacherID: &testutil.TestTeacherID,
				IsActive:       false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.course)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.course.ID == uuid.Nil {
					t.Error("Create() did not set course ID")
				}
				if tt.course.CreatedAt.IsZero() {
					t.Error("Create() did not set CreatedAt")
				}
				if tt.course.UpdatedAt.IsZero() {
					t.Error("Create() did not set UpdatedAt")
				}
			}
		})
	}
}

func TestCourseRepository_GetByID(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestTeacherID, "teacher@test.com", "$2a$10$hash", "teacher")
	fixtures.CreateTeacher(ctx, t, testutil.TestTeacherID, "Test", "Teacher", "CS")

	desc := "Test description"
	fixtures.CreateCourse(ctx, t, testutil.TestCourseID1, "Test Course", &desc, &testutil.TestTeacherID, true)

	repo := course.NewRepository(pool)

	tests := []struct {
		name      string
		courseID  uuid.UUID
		wantTitle string
		wantErr   bool
	}{
		{
			name:      "get existing course",
			courseID:  testutil.TestCourseID1,
			wantTitle: "Test Course",
			wantErr:   false,
		},
		{
			name:     "get non-existing course",
			courseID: uuid.MustParse("00000000-0000-0000-0000-000000000999"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(ctx, tt.courseID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Title != tt.wantTitle {
					t.Errorf("GetByID() title = %v, want %v", got.Title, tt.wantTitle)
				}
				if got.TeacherFirstName == nil || *got.TeacherFirstName != "Test" {
					t.Errorf("GetByID() teacher first name not properly joined")
				}
			}
		})
	}
}

func TestCourseRepository_List(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestTeacherID, "teacher@test.com", "$2a$10$hash", "teacher")
	fixtures.CreateTeacher(ctx, t, testutil.TestTeacherID, "Test", "Teacher", "CS")

	desc := "Description"
	fixtures.CreateCourse(ctx, t, testutil.TestCourseID1, "Course 1", &desc, &testutil.TestTeacherID, true)
	fixtures.CreateCourse(ctx, t, testutil.TestCourseID2, "Course 2", nil, &testutil.TestTeacherID, true)
	fixtures.CreateCourse(ctx, t, testutil.TestCourseID3, "Inactive Course", nil, &testutil.TestTeacherID, false)

	repo := course.NewRepository(pool)

	tests := []struct {
		name      string
		filter    *course.CourseFilter
		limit     int
		offset    int
		wantCount int
		wantTotal int64
	}{
		{
			name:      "list all courses",
			filter:    nil,
			limit:     10,
			offset:    0,
			wantCount: 3,
			wantTotal: 3,
		},
		{
			name:      "list with pagination",
			filter:    nil,
			limit:     2,
			offset:    0,
			wantCount: 2,
			wantTotal: 3,
		},
		{
			name:      "filter by active status",
			filter:    &course.CourseFilter{IsActive: boolPtr(true)},
			limit:     10,
			offset:    0,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:      "filter by teacher",
			filter:    &course.CourseFilter{TeacherID: &testutil.TestTeacherID},
			limit:     10,
			offset:    0,
			wantCount: 3,
			wantTotal: 3,
		},
		{
			name:      "search by title",
			filter:    &course.CourseFilter{Search: stringPtr("Course 1")},
			limit:     10,
			offset:    0,
			wantCount: 1,
			wantTotal: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courses, total, err := repo.List(ctx, tt.filter, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}

			if len(courses) != tt.wantCount {
				t.Errorf("List() count = %v, want %v", len(courses), tt.wantCount)
			}

			if total != tt.wantTotal {
				t.Errorf("List() total = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func TestCourseRepository_Update(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestTeacherID, "teacher@test.com", "$2a$10$hash", "teacher")
	fixtures.CreateTeacher(ctx, t, testutil.TestTeacherID, "Test", "Teacher", "CS")

	desc := "Original description"
	fixtures.CreateCourse(ctx, t, testutil.TestCourseID1, "Original Title", &desc, &testutil.TestTeacherID, true)

	repo := course.NewRepository(pool)

	t.Run("update course fields", func(t *testing.T) {
		newDesc := "Updated description"
		updated := &course.Course{
			ID:             testutil.TestCourseID1,
			Title:          "Updated Title",
			Description:    &newDesc,
			OwnerTeacherID: &testutil.TestTeacherID,
			IsActive:       false,
		}

		err := repo.Update(ctx, updated)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// Verify update
		got, err := repo.GetByID(ctx, testutil.TestCourseID1)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Title != "Updated Title" {
			t.Errorf("Update() title = %v, want %v", got.Title, "Updated Title")
		}
		if got.Description == nil || *got.Description != "Updated description" {
			t.Error("Update() description not updated correctly")
		}
		if got.IsActive != false {
			t.Error("Update() is_active not updated correctly")
		}
	})
}

func TestCourseRepository_Delete(t *testing.T) {
	pool := getTestPool(t)
	fixtures := testutil.NewFixtures(pool)
	ctx := context.Background()

	// Setup
	testutil.TruncateAll(t, pool)
	fixtures.CreateUser(ctx, t, testutil.TestTeacherID, "teacher@test.com", "$2a$10$hash", "teacher")
	fixtures.CreateTeacher(ctx, t, testutil.TestTeacherID, "Test", "Teacher", "CS")
	fixtures.CreateCourse(ctx, t, testutil.TestCourseID1, "To Delete", nil, &testutil.TestTeacherID, true)

	repo := course.NewRepository(pool)

	t.Run("delete existing course", func(t *testing.T) {
		err := repo.Delete(ctx, testutil.TestCourseID1)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// Verify deletion
		_, err = repo.GetByID(ctx, testutil.TestCourseID1)
		if err == nil {
			t.Error("Delete() course still exists after deletion")
		}
	})

	t.Run("delete non-existing course", func(t *testing.T) {
		// Deleting non-existing should not error (DELETE is idempotent)
		err := repo.Delete(ctx, uuid.MustParse("00000000-0000-0000-0000-000000000999"))
		if err != nil {
			t.Errorf("Delete() error = %v, want nil for non-existing", err)
		}
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
