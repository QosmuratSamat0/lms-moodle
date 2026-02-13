package dashboard

import (
	"context"
	"math"
	"time"

	"github.com/ap1-final-mini-moodle/internal/domain/dashboard"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) dashboard.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetStudentDashboard(ctx context.Context, studentID string) (*dashboard.StudentDashboard, error) {
	d := &dashboard.StudentDashboard{StudentID: studentID}

	// Get student name
	err := r.db.QueryRow(ctx, `SELECT first_name || ' ' || last_name FROM users WHERE id = $1`, studentID).Scan(&d.StudentName)
	if err != nil {
		return nil, err
	}

	// Get GPA
	d.GPA, _ = r.CalculateGPA(ctx, studentID)

	// Get course counts
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM enrollments WHERE student_id = $1`, studentID).Scan(&d.TotalCourses)
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM enrollments e WHERE e.student_id = $1 AND e.status = 'active'`, studentID).Scan(&d.ActiveCourses)

	// Get assignment counts
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM submissions WHERE student_id = $1`, studentID).Scan(&d.CompletedAssignments)
	r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM assignments a
		JOIN enrollments e ON a.course_id = e.course_id
		WHERE e.student_id = $1
		AND a.due_at > NOW()
		AND NOT EXISTS (SELECT 1 FROM submissions s WHERE s.assignment_id = a.id AND s.student_id = $1)
	`, studentID).Scan(&d.PendingAssignments)

	// Quiz stats - tables don't exist yet, set to 0
	d.TotalQuizzes = 0
	d.AverageQuizScore = 0

	// Get attendance rate
	d.AttendanceRate, _ = r.GetAttendanceRate(ctx, studentID)

	// Get upcoming deadlines
	d.UpcomingDeadlines, _ = r.GetUpcomingDeadlines(ctx, studentID, 5)

	// Get recent grades
	d.RecentGrades, _ = r.GetRecentGrades(ctx, studentID, 5)

	// Get grade trends
	d.GradeTrends, _ = r.GetGradeTrends(ctx, studentID, 6)

	// Get course progress
	d.CourseProgress, _ = r.GetCourseProgress(ctx, studentID)

	return d, nil
}

func (r *PostgresRepository) GetTeacherDashboard(ctx context.Context, teacherID string) (*dashboard.TeacherDashboard, error) {
	d := &dashboard.TeacherDashboard{TeacherID: teacherID}

	// Get teacher name
	err := r.db.QueryRow(ctx, `SELECT first_name || ' ' || last_name FROM users WHERE id = $1`, teacherID).Scan(&d.TeacherName)
	if err != nil {
		return nil, err
	}

	// Get total courses
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM courses WHERE owner_teacher_id = (SELECT id FROM teachers WHERE user_id = $1)`, teacherID).Scan(&d.TotalCourses)

	// Get total students
	r.db.QueryRow(ctx, `
		SELECT COUNT(DISTINCT e.student_id)
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		WHERE c.owner_teacher_id = (SELECT id FROM teachers WHERE user_id = $1)
	`, teacherID).Scan(&d.TotalStudents)

	// Get pending submissions
	r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM submissions s
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		WHERE c.owner_teacher_id = (SELECT id FROM teachers WHERE user_id = $1)
		AND NOT EXISTS (SELECT 1 FROM grades g WHERE g.submission_id = s.id)
	`, teacherID).Scan(&d.PendingSubmissions)

	// Get pending appeals
	r.db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM grade_appeals ga
		JOIN grades g ON ga.grade_id = g.id
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		WHERE c.owner_teacher_id = (SELECT id FROM teachers WHERE user_id = $1) AND ga.status = 'pending'
	`, teacherID).Scan(&d.PendingAppeals)

	// Get recent submissions
	rows, _ := r.db.Query(ctx, `
		SELECT s.id, u.first_name || ' ' || u.last_name, c.id, c.title, a.title, s.submitted_at,
		       EXISTS(SELECT 1 FROM grades g WHERE g.submission_id = s.id) as graded
		FROM submissions s
		JOIN users u ON s.student_id = u.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		WHERE c.owner_teacher_id = (SELECT id FROM teachers WHERE user_id = $1)
		ORDER BY s.submitted_at DESC
		LIMIT 10
	`, teacherID)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var rs dashboard.RecentSubmission
			rows.Scan(&rs.ID, &rs.StudentName, &rs.CourseID, &rs.CourseName, &rs.AssignmentTitle, &rs.SubmittedAt, &rs.Graded)
			d.RecentSubmissions = append(d.RecentSubmissions, rs)
		}
	}

	// Get course stats
	statsRows, _ := r.db.Query(ctx, `
		SELECT c.id, c.title,
		       (SELECT COUNT(*) FROM enrollments WHERE course_id = c.id) as student_count,
		       COALESCE((SELECT AVG(g.score / a.max_points * 100) FROM grades g JOIN submissions s ON g.submission_id = s.id JOIN assignments a ON s.assignment_id = a.id WHERE a.course_id = c.id), 0) as avg_grade,
		       COALESCE((SELECT AVG(CASE WHEN am.status = 'present' THEN 1 ELSE 0 END) * 100 FROM attendance_marks am JOIN attendance_sessions ats ON am.session_id = ats.id WHERE ats.course_id = c.id), 0) as attendance
		FROM courses c
		WHERE c.owner_teacher_id = (SELECT id FROM teachers WHERE user_id = $1)
	`, teacherID)
	if statsRows != nil {
		defer statsRows.Close()
		for statsRows.Next() {
			var cs dashboard.CourseStats
			statsRows.Scan(&cs.CourseID, &cs.CourseName, &cs.StudentCount, &cs.AverageGrade, &cs.AttendanceRate)
			d.CourseStats = append(d.CourseStats, cs)
		}
	}

	return d, nil
}

func (r *PostgresRepository) GetUpcomingDeadlines(ctx context.Context, studentID string, limit int) ([]dashboard.UpcomingDeadline, error) {
	var deadlines []dashboard.UpcomingDeadline

	// Assignments
	rows, err := r.db.Query(ctx, `
		SELECT a.id, 'assignment', c.id, c.title, a.title, a.due_at
		FROM assignments a
		JOIN courses c ON a.course_id = c.id
		JOIN enrollments e ON c.id = e.course_id
		WHERE e.student_id = $1
		AND a.due_at > NOW()
		AND NOT EXISTS (SELECT 1 FROM submissions s WHERE s.assignment_id = a.id AND s.student_id = $1)
		ORDER BY a.due_at
		LIMIT $2
	`, studentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d dashboard.UpcomingDeadline
		rows.Scan(&d.ID, &d.Type, &d.CourseID, &d.CourseName, &d.Title, &d.DueDate)
		d.DaysLeft = int(math.Ceil(time.Until(d.DueDate).Hours() / 24))
		d.Status = "pending"
		deadlines = append(deadlines, d)
	}

	// Quizzes table doesn't exist yet - skip

	return deadlines, nil
}

func (r *PostgresRepository) GetRecentGrades(ctx context.Context, studentID string, limit int) ([]dashboard.RecentGrade, error) {
	var grades []dashboard.RecentGrade

	rows, err := r.db.Query(ctx, `
		SELECT g.id, 'assignment', c.title, a.title, g.score, a.max_points, (g.score / a.max_points * 100), g.graded_at
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		JOIN courses c ON a.course_id = c.id
		WHERE s.student_id = $1
		ORDER BY g.graded_at DESC
		LIMIT $2
	`, studentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g dashboard.RecentGrade
		rows.Scan(&g.ID, &g.Type, &g.CourseName, &g.Title, &g.Score, &g.MaxPoints, &g.Percentage, &g.GradedAt)
		grades = append(grades, g)
	}

	return grades, nil
}

func (r *PostgresRepository) GetGradeTrends(ctx context.Context, studentID string, months int) ([]dashboard.GradeTrend, error) {
	var trends []dashboard.GradeTrend

	rows, err := r.db.Query(ctx, `
		SELECT TO_CHAR(g.graded_at, 'Mon') as month,
		       EXTRACT(YEAR FROM g.graded_at)::int as year,
		       AVG(g.score / a.max_points * 4.0) as avg_gpa
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		WHERE s.student_id = $1
		AND g.graded_at > NOW() - INTERVAL '1 month' * $2
		GROUP BY TO_CHAR(g.graded_at, 'Mon'), EXTRACT(YEAR FROM g.graded_at), DATE_TRUNC('month', g.graded_at)
		ORDER BY DATE_TRUNC('month', g.graded_at)
	`, studentID, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t dashboard.GradeTrend
		rows.Scan(&t.Month, &t.Year, &t.AverageGPA)
		trends = append(trends, t)
	}

	return trends, nil
}

func (r *PostgresRepository) GetCourseProgress(ctx context.Context, studentID string) ([]dashboard.CourseProgress, error) {
	var progress []dashboard.CourseProgress

	rows, err := r.db.Query(ctx, `
		SELECT c.id, c.title,
		       (SELECT COUNT(*) FROM submissions s JOIN assignments a ON s.assignment_id = a.id WHERE a.course_id = c.id AND s.student_id = $1) as completed,
		       (SELECT COUNT(*) FROM assignments WHERE course_id = c.id) as total,
		       COALESCE((SELECT AVG(g.score / a.max_points * 100) FROM grades g JOIN submissions s ON g.submission_id = s.id JOIN assignments a ON s.assignment_id = a.id WHERE a.course_id = c.id AND s.student_id = $1), 0) as current_grade,
		       COALESCE((SELECT u.first_name || ' ' || u.last_name FROM users u JOIN teachers t ON u.id = t.user_id WHERE t.id = c.owner_teacher_id), 'Unknown') as instructor_name
		FROM courses c
		JOIN enrollments e ON c.id = e.course_id
		WHERE e.student_id = $1
	`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p dashboard.CourseProgress
		rows.Scan(&p.CourseID, &p.CourseName, &p.CompletedAssignments, &p.TotalAssignments, &p.CurrentGrade, &p.InstructorName)
		if p.TotalAssignments > 0 {
			p.ProgressPercent = float64(p.CompletedAssignments) / float64(p.TotalAssignments) * 100
		}
		// Compute letter grade
		switch {
		case p.CurrentGrade >= 90:
			p.LetterGrade = "A"
		case p.CurrentGrade >= 80:
			p.LetterGrade = "B"
		case p.CurrentGrade >= 70:
			p.LetterGrade = "C"
		case p.CurrentGrade >= 60:
			p.LetterGrade = "D"
		default:
			p.LetterGrade = "F"
		}
		progress = append(progress, p)
	}

	return progress, nil
}

func (r *PostgresRepository) CalculateGPA(ctx context.Context, studentID string) (float64, error) {
	var gpa float64
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(AVG(g.score / a.max_points * 4.0), 0)
		FROM grades g
		JOIN submissions s ON g.submission_id = s.id
		JOIN assignments a ON s.assignment_id = a.id
		WHERE s.student_id = $1
	`, studentID).Scan(&gpa)
	return math.Round(gpa*100) / 100, err
}

func (r *PostgresRepository) GetAttendanceRate(ctx context.Context, studentID string) (float64, error) {
	var rate float64
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(
			(SELECT COUNT(*) FILTER (WHERE am.status = 'present') * 100.0 / NULLIF(COUNT(*), 0)
			 FROM attendance_marks am
			 WHERE am.student_id = $1
			), 0
		)
	`, studentID).Scan(&rate)
	return math.Round(rate*100) / 100, err
}
