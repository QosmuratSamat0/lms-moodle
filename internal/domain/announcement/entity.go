package announcement

import "time"

type Announcement struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Pinned    bool      `json:"pinned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAnnouncementInput struct {
	CourseID string `json:"course_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Pinned   bool   `json:"pin"`
}

type UpdateAnnouncementInput struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Pinned  *bool   `json:"pin"`
}
