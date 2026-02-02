package upload

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the attachment repository interface
type Repository interface {
	Create(ctx context.Context, attachment *Attachment) error
	GetByID(ctx context.Context, id uuid.UUID) (*AttachmentWithUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByReference(ctx context.Context, refType AttachmentType, refID uuid.UUID, limit, offset int) ([]AttachmentWithUser, int64, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]AttachmentWithUser, int64, error)
	GetByPublicID(ctx context.Context, publicID string) (*Attachment, error)
}

type repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new attachment repository
func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, attachment *Attachment) error {
	query := `
		INSERT INTO attachments (uploaded_by, type, reference_id, public_id, url, secure_url, 
			original_name, format, resource_type, size, width, height)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query,
		attachment.UploadedBy,
		attachment.Type,
		attachment.ReferenceID,
		attachment.PublicID,
		attachment.URL,
		attachment.SecureURL,
		attachment.OriginalName,
		attachment.Format,
		attachment.ResourceType,
		attachment.Size,
		attachment.Width,
		attachment.Height,
	).Scan(&attachment.ID, &attachment.CreatedAt)
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*AttachmentWithUser, error) {
	query := `
		SELECT 
			a.id, a.uploaded_by, a.type, a.reference_id, a.public_id, a.url, a.secure_url,
			a.original_name, a.format, a.resource_type, a.size, a.width, a.height, a.created_at,
			COALESCE(s.first_name, t.first_name, m.first_name) as uploader_first_name,
			COALESCE(s.last_name, t.last_name, m.last_name) as uploader_last_name,
			u.email as uploader_email
		FROM attachments a
		JOIN users u ON a.uploaded_by = u.id
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN teachers t ON u.id = t.user_id
		LEFT JOIN managers m ON u.id = m.user_id
		WHERE a.id = $1`

	var att AttachmentWithUser
	err := r.db.QueryRow(ctx, query, id).Scan(
		&att.ID, &att.UploadedBy, &att.Type, &att.ReferenceID, &att.PublicID, &att.URL, &att.SecureURL,
		&att.OriginalName, &att.Format, &att.ResourceType, &att.Size, &att.Width, &att.Height, &att.CreatedAt,
		&att.UploaderFirstName, &att.UploaderLastName, &att.UploaderEmail,
	)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM attachments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *repository) ListByReference(ctx context.Context, refType AttachmentType, refID uuid.UUID, limit, offset int) ([]AttachmentWithUser, int64, error) {
	countQuery := `SELECT COUNT(*) FROM attachments WHERE type = $1 AND reference_id = $2`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, refType, refID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			a.id, a.uploaded_by, a.type, a.reference_id, a.public_id, a.url, a.secure_url,
			a.original_name, a.format, a.resource_type, a.size, a.width, a.height, a.created_at,
			COALESCE(s.first_name, t.first_name, m.first_name) as uploader_first_name,
			COALESCE(s.last_name, t.last_name, m.last_name) as uploader_last_name,
			u.email as uploader_email
		FROM attachments a
		JOIN users u ON a.uploaded_by = u.id
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN teachers t ON u.id = t.user_id
		LEFT JOIN managers m ON u.id = m.user_id
		WHERE a.type = $1 AND a.reference_id = $2
		ORDER BY a.created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(ctx, query, refType, refID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanAttachments(rows, total)
}

func (r *repository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]AttachmentWithUser, int64, error) {
	countQuery := `SELECT COUNT(*) FROM attachments WHERE uploaded_by = $1`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			a.id, a.uploaded_by, a.type, a.reference_id, a.public_id, a.url, a.secure_url,
			a.original_name, a.format, a.resource_type, a.size, a.width, a.height, a.created_at,
			COALESCE(s.first_name, t.first_name, m.first_name) as uploader_first_name,
			COALESCE(s.last_name, t.last_name, m.last_name) as uploader_last_name,
			u.email as uploader_email
		FROM attachments a
		JOIN users u ON a.uploaded_by = u.id
		LEFT JOIN students s ON u.id = s.user_id
		LEFT JOIN teachers t ON u.id = t.user_id
		LEFT JOIN managers m ON u.id = m.user_id
		WHERE a.uploaded_by = $1
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanAttachments(rows, total)
}

func (r *repository) GetByPublicID(ctx context.Context, publicID string) (*Attachment, error) {
	query := `
		SELECT id, uploaded_by, type, reference_id, public_id, url, secure_url,
			original_name, format, resource_type, size, width, height, created_at
		FROM attachments
		WHERE public_id = $1`

	var att Attachment
	err := r.db.QueryRow(ctx, query, publicID).Scan(
		&att.ID, &att.UploadedBy, &att.Type, &att.ReferenceID, &att.PublicID, &att.URL, &att.SecureURL,
		&att.OriginalName, &att.Format, &att.ResourceType, &att.Size, &att.Width, &att.Height, &att.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *repository) scanAttachments(rows interface{ Next() bool; Scan(dest ...interface{}) error }, total int64) ([]AttachmentWithUser, int64, error) {
	var attachments []AttachmentWithUser
	for rows.Next() {
		var att AttachmentWithUser
		if err := rows.Scan(
			&att.ID, &att.UploadedBy, &att.Type, &att.ReferenceID, &att.PublicID, &att.URL, &att.SecureURL,
			&att.OriginalName, &att.Format, &att.ResourceType, &att.Size, &att.Width, &att.Height, &att.CreatedAt,
			&att.UploaderFirstName, &att.UploaderLastName, &att.UploaderEmail,
		); err != nil {
			return nil, 0, err
		}
		attachments = append(attachments, att)
	}
	return attachments, total, nil
}
