package infra_repository

import (
	"database/sql"
	"fmt"
	"time"

	"forum/domain/entity"
	custom_errors "forum/domain/errors"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type SQLiteCommentRepository struct {
	db *sql.DB
	userRepo         repository.UserRepository
	commentReaction	repository.CommentReactionRepository
}

func NewSQLiteCommentRepository(
	db *sql.DB,
	userRepo *repository.UserRepository,
	commentReaction	*repository.CommentReactionRepository,
	) repository.CommentRepository {
	return &SQLiteCommentRepository{
		db: db,
		userRepo:  *userRepo,
		commentReaction: *commentReaction,
	}
}

func (r *SQLiteCommentRepository) Create(comment *entity.Comment) error {
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()

	query := `INSERT INTO comments (id, content, user_id, post_id, createdat)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, comment.ID.String(), comment.Content,
		comment.UserID.String(), comment.PostID.String(), comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	return nil
}

func (r *SQLiteCommentRepository) GetByID(commentID uuid.UUID) (*entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat FROM comments WHERE id = ?`

	row := r.db.QueryRow(query, commentID.String())

	comment := &entity.Comment{}
	var idStr, userIDStr, postIDStr string

	err := row.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrCommentNotFound
		}
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	comment.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	comment.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	comment.PostID, err = uuid.Parse(postIDStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	return comment, nil
}

func (r *SQLiteCommentRepository) GetByPostID(postID uuid.UUID) ([]entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat 
			  FROM comments WHERE post_id = ? ORDER BY createdat ASC`

	rows, err := r.db.Query(query, postID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	var comments []entity.Comment

	for rows.Next() {
		comment := entity.Comment{}
		var idStr, userIDStr, postIDStr string

		err := rows.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *SQLiteCommentRepository) GetByUserID(userID uuid.UUID) ([]*entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat 
			  FROM comments WHERE user_id = ? ORDER BY createdat DESC`

	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	var comments []*entity.Comment

	for rows.Next() {
		comment := &entity.Comment{}
		var idStr, userIDStr, postIDStr string

		err := rows.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *SQLiteCommentRepository) GetByPostIDWithPagination(postID uuid.UUID, limit, offset int) ([]*entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat 
			  FROM comments WHERE post_id = ? ORDER BY createdat ASC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, postID.String(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	var comments []*entity.Comment

	for rows.Next() {
		comment := &entity.Comment{}
		var idStr, userIDStr, postIDStr string

		err := rows.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comment.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *SQLiteCommentRepository) Update(comment *entity.Comment) error {
	query := `UPDATE comments SET content = ? WHERE id = ?`

	result, err := r.db.Exec(query, comment.Content, comment.ID.String())
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return custom_errors.ErrCommentNotFound
	}

	return nil
}

func (r *SQLiteCommentRepository) Delete(commentID uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = ?`

	result, err := r.db.Exec(query, commentID.String())
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return custom_errors.ErrCommentNotFound
	}

	return nil
}

func (r *SQLiteCommentRepository) GetCountByPostID(postID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = ?`

	var count int
	err := r.db.QueryRow(query, postID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	return count, nil
}

func (r *SQLiteCommentRepository) GetCountByUserID(userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE user_id = ?`

	var count int
	err := r.db.QueryRow(query, userID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	return count, nil
}

func (r *SQLiteCommentRepository) GetWithDetails(commentID uuid.UUID) (*entity.CommentWithDetails, error) {
	comment, err := r.GetByID(commentID)
	if err != nil {
		return nil, err
	}

	return &entity.CommentWithDetails{
		Comment: *comment,
	}, nil
}

func (r *SQLiteCommentRepository) GetByPostIDWithDetails(postID uuid.UUID) ([]entity.CommentWithDetails, error) {
	comments, err := r.GetByPostID(postID)
	if err != nil {
		return nil, err
	}

	var commentsWithDetails []entity.CommentWithDetails
	for _, comment := range comments {
		user,err:=r.userRepo.GetByID(comment.UserID)
		if err != nil {
			fmt.Printf("comment user error ")
		}
		likes,dislikes,err:=r.commentReaction.GetReactionCountsByCommentID(comment.ID)
		if err != nil {
			fmt.Printf("comment user fell error ")
		}
		commentsWithDetails = append(commentsWithDetails, entity.CommentWithDetails{
			Comment: comment,
			Author: *user,
			LikeCount: likes,
			DislikeCount: dislikes,
		})
	}

	return commentsWithDetails, nil
}
