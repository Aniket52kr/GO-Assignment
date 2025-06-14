package database

import (
	"log"

	"github.com/Aniket52kr/GO-Assignment/models"
)

func CreatePost(userId string, post *models.Post) bool {
	if _, err := db.Exec(
		`INSERT INTO posts(user_id, id, body, created_at)
		VALUES (?, ?, ?, ?)`,
		userId, post.Id, post.Body, post.CreatedAt,
	); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func ReadPost(id string) *models.Post {
	var post models.Post
	if err := db.QueryRow(`SELECT * FROM posts WHERE id = ?`, id).Scan(
		&post.UserId, &post.Id, &post.Body, &post.CreatedAt,
	); err != nil {
		log.Println(err)
		return nil
	}
	return &post
}

func ReadPostsCount(userId string) int {
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM posts WHERE user_id = ?`, userId).Scan(&count); err != nil {
		log.Println(err)
		return 0
	}
	return count
}

func ReadPosts(userId string, limit int, offset int) []models.Post {
	var posts []models.Post
	rows, err := db.Query(
		`SELECT * FROM posts WHERE user_id = ? ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		userId, limit, offset,
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.UserId, &post.Id, &post.Body, &post.CreatedAt)
		posts = append(posts, post)
	}
	return posts
}

func ReadFeedPosts(userId string, limit int, offset int) []models.Post {
	var posts []models.Post
	rows, err := db.Query(
		`SELECT * FROM posts WHERE user_id IN
		(SELECT follow_id FROM follows WHERE user_id = ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		userId, limit, offset,
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		rows.Scan(&post.UserId, &post.Id, &post.Body, &post.CreatedAt)
		posts = append(posts, post)
	}
	return posts
}

func DeletePost(id string) bool {
	if _, err := db.Exec(`DELETE FROM posts WHERE id = ?`, id); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func Voted(userId string, id string) bool {
	var count int
	db.QueryRow(
		`SELECT COUNT(*) FROM votes WHERE user_id = ? AND id = ?`,
		userId, id,
	).Scan(&count)

	return count > 0
}

func ToggleVote(userId string, id string) {
	var query string
	voted := Voted(userId, id)

	if voted {
		query = `DELETE FROM votes WHERE user_id = ? AND id = ?`
	} else {
		query = `INSERT INTO votes (user_id, id) VALUES (?, ?)`
	}
	if _, err := db.Exec(query, userId, id); err != nil {
		log.Println(err)
	}
}

func ReadVotes(id string) []string {
	var voters []string
	rows, err := db.Query(
		`SELECT username FROM t_users WHERE id IN
		(SELECT user_id FROM votes WHERE id = ?)`,
		id,
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		rows.Scan(&username)
		voters = append(voters, username)
	}
	return voters
}

func CreateComment(userId string, postId string, comment *models.Comment) bool {
	if _, err := db.Exec(
		`INSERT INTO comments (user_id, post_id, id, body, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		userId, postId, comment.Id, comment.Body, comment.CreatedAt,
	); err != nil {
		log.Println(err)
		return false
	}
	return true
}

func ReadComment(id string) *models.Comment {
	var comment models.Comment
	if err := db.QueryRow(`SELECT * FROM comments WHERE id = ?`, id).Scan(
		&comment.UserId,
		&comment.PostId,
		&comment.Id,
		&comment.Body,
		&comment.CreatedAt,
	); err != nil {
		log.Println(err)
		return nil
	}
	return &comment
}

func ReadComments(postId string, limit int, offset int) []models.Comment {
	var comments []models.Comment
	rows, err := db.Query(
		`SELECT * FROM comments WHERE post_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`,
		postId, limit, offset,
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		rows.Scan(
			&comment.UserId,
			&comment.PostId,
			&comment.Id,
			&comment.Body,
			&comment.CreatedAt,
		)
		comments = append(comments, comment)
	}
	return comments
}

func DeleteComment(id string) bool {
	if _, err := db.Exec(`DELETE FROM comments WHERE id = ?`, id); err != nil {
		log.Println(err)
		return false
	}
	return true
}
