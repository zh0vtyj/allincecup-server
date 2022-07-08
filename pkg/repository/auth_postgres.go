package repository

import (
	server "allincecup-server"
	"allincecup-server/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user server.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (role_id, name, email, password_hash, phone_number) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)
	row := a.db.QueryRow(query, user.RoleId, user.Name, user.Email, user.Password, user.PhoneNumber)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a *AuthPostgres) GetUser(email, password string) (server.User, error) {
	var user server.User
	query := fmt.Sprintf("SELECT id, role_id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := a.db.Get(&user, query, email, password)

	return user, err
}

func (a *AuthPostgres) NewSession(session domain.Session) (*domain.Session, error) {
	var newSession domain.Session

	query := fmt.Sprintf(
		"INSERT INTO %s (user_id, role_id, refresh_token, client_ip, user_agent, expires_at) values ($1, $2, $3, $4, $5, $6) RETURNING *",
		sessionsTable)
	row := a.db.QueryRow(query, session.UserId, session.RoleId, session.RefreshToken, session.ClientIp, session.UserAgent, session.ExpiresAt)

	if err := row.Scan(
		&newSession.Id,
		&newSession.UserId,
		&newSession.RoleId,
		&newSession.RefreshToken,
		&newSession.IsBlocked,
		&newSession.ClientIp,
		&newSession.UserAgent,
		&newSession.ExpiresAt,
		&newSession.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &newSession, nil
}

func (a *AuthPostgres) DeleteSessionByRefresh(refresh string) error {
	var id int
	query := fmt.Sprintf("DELETE from %s WHERE refresh_token=$1 RETURNING id", sessionsTable)
	row := a.db.QueryRow(query, refresh)

	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthPostgres) GetSessionByRefresh(refresh string) (*domain.Session, error) {
	var session domain.Session
	query := fmt.Sprintf("SELECT * from %s WHERE refresh_token=$1 LIMIT 1", sessionsTable)
	row := a.db.QueryRow(query, refresh)

	if err := row.Scan(
		&session.Id,
		&session.UserId,
		&session.RoleId,
		&session.RefreshToken,
		&session.IsBlocked,
		&session.ClientIp,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &session, nil
}
