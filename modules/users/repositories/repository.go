package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Sunwatcha303/OAuth-golang-demo/modules/entities"
	"github.com/Sunwatcha303/OAuth-golang-demo/pkg/databases"
	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	db *databases.Database
}

func NewUsersRepository(db *databases.Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetAccessTokenByUserId(id int) (accessToken *entities.AccessToken, err error) {

	key := fmt.Sprintf("user:%d:access_token", id)

	data, err := r.db.Redis.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("access token not found")
	} else if err != nil {
		return nil, fmt.Errorf("failed to get token from redis: %v", err)
	}

	err = json.Unmarshal([]byte(data), accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize token: %v", err)
	}

	return accessToken, err
}

func (r *UserRepository) GetRefreshTokenByUserId(id int) (string, error) {
	var refreshToken string

	err := r.db.PostgreSQL.QueryRow("SELECT refresh_token FROM users WHERE id = $1", id).Scan(&refreshToken)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return refreshToken, nil
}

func (r *UserRepository) SaveAccessToken(id int, token string, expiresAt time.Time) error {
	key := fmt.Sprintf("user:%d:access_token", id)

	accessToken := entities.AccessToken{
		Id:         id,
		Token:      token,
		Expires_at: expiresAt,
	}

	data, err := json.Marshal(accessToken)
	if err != nil {
		return fmt.Errorf("failed to serialize token: %v", err)
	}
	expiration := time.Until(expiresAt)

	err = r.db.Redis.Set(context.Background(), key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set token in redis: %v", err)
	}

	return nil
}

func (r *UserRepository) GetUserById(id int) (user *entities.User, err error) {
	user = &entities.User{}
	query := `SELECT id, sub, email, name, picture, refresh_token, expires_at, scope, created_at, updated_at 
              FROM users WHERE id = $1`

	err = r.db.PostgreSQL.Get(user, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Println("Error fetching user:", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *entities.User) (int, error) {
	query := `INSERT INTO users (sub, email, name, picture, refresh_token, expires_at, scope)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          RETURNING id`

	var newID int
	err := r.db.PostgreSQL.Get(&newID, query, user.Sub, user.Email, user.Name, user.Picture, user.Refresh_token, user.Expires_at, user.Scope)
	if err != nil {
		log.Println("Error inserting user:", err)
		return 0, err
	}

	return newID, nil
}

func (r *UserRepository) GetUserBySup(sub string) (user *entities.User, err error) {
	user = &entities.User{}
	query := `SELECT id, sub, email, name, picture, refresh_token, expires_at, scope, created_at, updated_at 
              FROM users WHERE sub = $1`

	err = r.db.PostgreSQL.Get(user, query, sub)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		log.Println("Error fetching user:", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *entities.User) error {
	user.Updated_at = time.Now()

	query := `
	UPDATE users 
	SET name = $1, email = $2, picture = $3, 
	    refresh_token = $4, expires_at = $5, 
	    scope = $6, updated_at = $7
	WHERE id = $8
`
	result, err := r.db.PostgreSQL.Exec(query, user.Name, user.Email, user.Picture,
		user.Refresh_token, user.Expires_at, user.Scope,
		user.Updated_at, user.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
