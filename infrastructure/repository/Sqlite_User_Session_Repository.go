package infra_repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type SQLiteUserSessionRepository struct {
	db *sql.DB
}

func NewSQLiteUserSessionRepository(db *sql.DB) repository.UserSessionRepository {
	return &SQLiteUserSessionRepository{db: db}
}

func (r *SQLiteUserSessionRepository) Create(session *entity.UserSession) error {
	session.ID = uuid.New()
	session.CreatedAt = time.Now()
	query := `INSERT INTO user_sessions (id, user_id, session_token, expires_at, created_at)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, session.ID.String(), session.UserID.String(),
		session.SessionToken, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *SQLiteUserSessionRepository) GetByToken(token string) (*entity.UserSession, error) {
	query := `SELECT id, user_id, session_token, expires_at, created_at 
			  FROM user_sessions WHERE session_token = ?`

	row := r.db.QueryRow(query, token)

	session := &entity.UserSession{}
	var idStr, userIDStr string

	err := row.Scan(&idStr, &userIDStr, &session.SessionToken, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("you need to login")
		}
		return nil, fmt.Errorf("failed to scan session: %v", err)
	}

	session.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	session.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SQLiteUserSessionRepository) GetByUserID(userID uuid.UUID) (*entity.UserSession, error) {
	query := `SELECT id, user_id, session_token, expires_at, created_at 
			  FROM user_sessions WHERE user_id = ? `

	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions *entity.UserSession

	for rows.Next() {
		session := &entity.UserSession{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &userIDStr, &session.SessionToken, &session.ExpiresAt, &session.CreatedAt)
		if err != nil {
			return nil, err
		}

		session.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		session.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
	}

	return sessions, nil
}

func (r *SQLiteUserSessionRepository) Update(session *entity.UserSession) error {
	query := `UPDATE user_sessions SET expires_at = ? WHERE id = ?`

	_, err := r.db.Exec(query, session.ExpiresAt, session.ID.String())
	return err
}

func (r *SQLiteUserSessionRepository) Delete(sessionID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE id = ?`

	_, err := r.db.Exec(query, sessionID.String())
	return err
}

func (r *SQLiteUserSessionRepository) DeleteAllUserSessions(userID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE user_id = ?`

	_, err := r.db.Exec(query, userID.String())
	return err
}
