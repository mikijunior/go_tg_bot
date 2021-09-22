package sqlstore

import (
	"bot/models"
	"database/sql"
	"errors"
	"log"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *models.User) error {
	return r.store.db.QueryRow(
		"INSERT INTO users (username, chat_id) VALUES ($1, $2) RETURNING id",
		u.Username,
		u.ChatId,
	).Scan(&u.ID)
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	u := &models.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, username, chat_id FROM users WHERE username = $1",
		username,
	).Scan(
		&u.ID,
		&u.Username,
		&u.ChatId,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}

		return nil, err
	}

	return u, nil
}


func (r *UserRepository) FindByChatId(chatId int) (*models.User, error) {
	u := &models.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, username, chat_id FROM users WHERE chat_id = $1",
		chatId,
	).Scan(
		&u.ID,
		&u.Username,
		&u.ChatId,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.store.db.Query(`SELEST * FROM users`)

	if err != nil {
		log.Fatal("Unable to load users")
	}

	var users []models.User

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(&user.ID, &user.Username, &user.ChatId)
		
		if err != nil {
			log.Fatal("Unable to scan user row", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}
