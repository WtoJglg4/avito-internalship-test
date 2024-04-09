package repository

import (
	"github/avito/entities"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const (
	usersTable = "users"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user entities.User) (int, error) {
	var id int

	sql, args := sqlbuilder.PostgreSQL.NewInsertBuilder().
		InsertInto(usersTable).
		SQL("(login, password_hash, role)").
		Values(user.Login, user.Password, user.Role).
		SQL("returning id;").Build()

	row := r.db.QueryRow(sql, args...)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (entities.User, error) {
	var user entities.User
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sql, args := sb.
		Select("*").
		From(usersTable).
		Where(
			sb.Equal("login", login),
			sb.Equal("password_hash", password),
		).Build()

	logrus.Println(sql, args)
	err := r.db.Get(&user, sql, args...)
	return user, err
}
