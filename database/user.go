package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Aniket52kr/GO-Assignment/models"
)

func CreateUser(user *models.User) bool {
	userEmail := ""
	if user.Email != nil {
		userEmail = *user.Email
	}
	userAvatar := ""
	if user.Avatar != nil {
		userAvatar = *user.Avatar
	}

	log.Printf("Creating user: %+v\n", user)

	_, err := db.Exec(`
		INSERT INTO t_users (email, username, password, id, verified, avatar, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userEmail, user.Username, user.Password, user.Id, user.Verified, userAvatar, user.CreatedAt)

	if err != nil {
		log.Println("Insert failed:", err)
		return false
	}
	return true
}

func CreateOAuthUser(id string) bool {
	_, err := db.Exec(`INSERT INTO o_users(id) VALUES (?)`, id)
	if err != nil {
		log.Println("CreateOAuthUser failed:", err)
		return false
	}
	return true
}

func ReadUserByName(username string) *models.User {
	var user models.User
	var email, avatar sql.NullString

	err := db.QueryRow(`
		SELECT email, username, password, id, verified, avatar, created_at
		FROM t_users WHERE username = ?`, username).
		Scan(&email, &user.Username, &user.Password, &user.Id, &user.Verified, &avatar, &user.CreatedAt)

	if err != nil {
		log.Println("ReadUserByName error:", err)
		return nil
	}

	if email.Valid {
		user.Email = &email.String
	}
	if avatar.Valid {
		user.Avatar = &avatar.String
	}
	return &user
}

func ReadUserByEmail(email string) *models.User {
	var user models.User
	var emailField, avatar sql.NullString

	err := db.QueryRow(`
		SELECT email, username, password, id, verified, avatar, created_at
		FROM t_users WHERE email = ?`, email).
		Scan(&emailField, &user.Username, &user.Password, &user.Id, &user.Verified, &avatar, &user.CreatedAt)

	if err != nil {
		log.Println("ReadUserByEmail error:", err)
		return nil
	}

	if emailField.Valid {
		user.Email = &emailField.String
	}
	if avatar.Valid {
		user.Avatar = &avatar.String
	}
	return &user
}

func ReadUserById(id string) *models.User {
	var user models.User
	var email, avatar sql.NullString

	err := db.QueryRow(`
		SELECT email, username, password, id, verified, avatar, created_at
		FROM t_users WHERE id = ?`, id).
		Scan(&email, &user.Username, &user.Password, &user.Id, &user.Verified, &avatar, &user.CreatedAt)

	if err != nil {
		log.Println("ReadUserById error:", err)
		return nil
	}

	if email.Valid {
		user.Email = &email.String
	}
	if avatar.Valid {
		user.Avatar = &avatar.String
	}
	return &user
}

func IsOAuthUser(id string) bool {
	var count int
	_ = db.QueryRow(`SELECT COUNT(*) FROM o_users WHERE id = ?`, id).Scan(&count)
	return count > 0
}

func ReadUsers(username string, limit int, offset int) []models.User {
	var users []models.User

	rows, err := db.Query(`
		SELECT email, username, password, id, verified, avatar, created_at
		FROM t_users WHERE username LIKE ? ORDER BY username LIMIT ? OFFSET ?`,
		"%"+username+"%", limit, offset)
	if err != nil {
		log.Println("ReadUsers error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var email, avatar sql.NullString

		err := rows.Scan(&email, &user.Username, &user.Password, &user.Id, &user.Verified, &avatar, &user.CreatedAt)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		if email.Valid {
			user.Email = &email.String
		}
		if avatar.Valid {
			user.Avatar = &avatar.String
		}
		users = append(users, user)
	}
	return users
}

func UpdateUser(id string, updates map[string]any) bool {
	for column, value := range updates {
		query := fmt.Sprintf("UPDATE t_users SET %s = ? WHERE id = ?", strings.ReplaceAll(column, "`", ""))
		if _, err := db.Exec(query, value, id); err != nil {
			log.Println("UpdateUser error:", err)
			return false
		}
	}
	return true
}

func DeleteUser(id string) bool {
	if _, err := db.Exec(`DELETE FROM t_users WHERE id = ?`, id); err != nil {
		log.Println("DeleteUser error:", err)
		return false
	}
	return true
}

func Followed(userId, followId string) bool {
	var count int
	_ = db.QueryRow(`SELECT COUNT(*) FROM follows WHERE user_id = ? AND follow_id = ?`, userId, followId).Scan(&count)
	return count > 0
}

func ToggleFollow(userId, followId string) {
	var query string
	if Followed(userId, followId) {
		query = `DELETE FROM follows WHERE user_id = ? AND follow_id = ?`
	} else {
		query = `INSERT INTO follows(user_id, follow_id) VALUES (?, ?)`
	}
	if _, err := db.Exec(query, userId, followId); err != nil {
		log.Println("ToggleFollow error:", err)
	}
}

func ReadFollowers(userId string) []string {
	var followers []string
	rows, err := db.Query(`
		SELECT username FROM t_users WHERE id IN (SELECT user_id FROM follows WHERE follow_id = ?)`, userId)
	if err != nil {
		log.Println("ReadFollowers error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		rows.Scan(&username)
		followers = append(followers, username)
	}
	return followers
}

func ReadFollowersCount(userId string) int {
	var count int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM t_users WHERE id IN (SELECT user_id FROM follows WHERE follow_id = ?)`, userId).Scan(&count); err != nil {
		log.Println("ReadFollowersCount error:", err)
		return 0
	}
	return count
}

func ReadFollowing(userId string) []string {
	var following []string
	rows, err := db.Query(`
		SELECT username FROM t_users WHERE id IN (SELECT follow_id FROM follows WHERE user_id = ?)`, userId)
	if err != nil {
		log.Println("ReadFollowing error:", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		rows.Scan(&username)
		following = append(following, username)
	}
	return following
}

func ReadFollowingCount(userId string) int {
	var count int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM t_users WHERE id IN (SELECT follow_id FROM follows WHERE user_id = ?)`, userId).Scan(&count); err != nil {
		log.Println("ReadFollowingCount error:", err)
		return 0
	}
	return count
}

func CreateVerificationId(token string, id string) bool {
	_, err := db.Exec(`INSERT INTO shorturl(token, id) VALUES (?, ?)`, token, id)
	if err != nil {
		log.Println("CreateVerificationId error:", err)
		return false
	}
	return true
}

func ReadVerificationId(id string) string {
	var token string
	if err := db.QueryRow(`SELECT token FROM shorturl WHERE id = ?`, id).Scan(&token); err != nil {
		log.Println("ReadVerificationId error:", err)
		return ""
	}
	return token
}

func DeleteVerificationId(id string) bool {
	if _, err := db.Exec(`DELETE FROM shorturl WHERE id = ?`, id); err != nil {
		log.Println("DeleteVerificationId error:", err)
		return false
	}
	return true
}
