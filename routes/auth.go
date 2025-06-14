package routes

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Aniket52kr/GO-Assignment/database"
	"github.com/Aniket52kr/GO-Assignment/middleware"
	"github.com/Aniket52kr/GO-Assignment/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	issuer    string
	secretKey []byte
)

func init() {
	godotenv.Load(".env")
	issuer = os.Getenv("ISSUER")
	secretKey = []byte(os.Getenv("SECRET_KEY"))
}

func SignUp(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
			"type": "signup",
		})
	case "POST":
		var user models.User

		if err := c.Request.ParseForm(); err != nil {
			log.Println("ParseForm error:", err)
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{
				"error":   "400 Bad Request",
				"message": "Unable to parse form.",
			})
			return
		}

		if err := c.ShouldBindWith(&user, binding.Form); err != nil {
			log.Println("Form bind error:", err)
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{
				"error":   "400 Bad Request",
				"message": err.Error(),
			})
			return
		}

		// Sanitize email
		if user.Email != nil && *user.Email == "" {
			user.Email = nil
		}

		// Check for existing username
		if database.ReadUserByName(user.Username) != nil {
			c.HTML(http.StatusForbidden, "error.tmpl.html", gin.H{
				"error":   "403 Forbidden",
				"message": "Username already taken.",
			})
			return
		}

		// Check for existing email
		if user.Email != nil && database.ReadUserByEmail(*user.Email) != nil {
			c.HTML(http.StatusForbidden, "error.tmpl.html", gin.H{
				"error":   "403 Forbidden",
				"message": "An account already exists with this email.",
			})
			return
		}

		user.Id = uuid.NewString()
		user.Verified = false
		user.CreatedAt = time.Now()

		if err := user.HashPassword(); err != nil {
			log.Println("Password hashing error:", err)
			c.HTML(http.StatusInternalServerError, "error.tmpl.html", gin.H{
				"error":   "500 Internal Server Error",
				"message": "Failed to hash password.",
			})
			return
		}

		if ok := database.CreateUser(&user); !ok {
			log.Println("Database create user error")
			c.HTML(http.StatusInternalServerError, "error.tmpl.html", gin.H{
				"error":   "500 Internal Server Error",
				"message": "Failed to create user account.",
			})
			return
		}

		token, _ := middleware.CreateToken(user.Id)
		session := sessions.Default(c)
		session.Set("Authorization", token)
		session.Set("userId", user.Id)
		session.Save()

		c.Redirect(http.StatusFound, "/auth/verify?signup=true")
	}
}

func Login(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		c.HTML(http.StatusOK, "auth.tmpl.html", gin.H{
			"type": "login",
		})
	case "POST":
		var login models.Login

		if err := c.Request.ParseForm(); err != nil {
			log.Println("ParseForm error:", err)
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{
				"error":   "400 Bad Request",
				"message": "Unable to parse form.",
			})
			return
		}

		if err := c.ShouldBindWith(&login, binding.Form); err != nil {
			log.Println("Form bind error:", err)
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{
				"error":   "403 Forbidden",
				"message": err.Error(),
			})
			return
		}

		user := database.ReadUserByName(login.Username)
		if user == nil {
			c.HTML(http.StatusUnauthorized, "error.tmpl.html", gin.H{
				"error":   "401 Unauthorized",
				"message": "User does not exist.",
			})
			return
		}

		if !user.CheckPassword(login.Password) {
			c.HTML(http.StatusUnauthorized, "error.tmpl.html", gin.H{
				"error":   "401 Unauthorized",
				"message": "Incorrect password.",
			})
			return
		}

		token, _ := middleware.CreateToken(user.Id)
		session := sessions.Default(c)
		session.Set("Authorization", token)
		session.Set("userId", user.Id)
		session.Save()

		c.Redirect(http.StatusFound, "/feed")
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("userId")

	if userId == nil {
		c.HTML(http.StatusUnauthorized, "error.tmpl.html", gin.H{
			"error":   "401 Unauthorized",
			"message": "User not logged in.",
		})
		return
	}

	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	session.Save()

	c.HTML(http.StatusOK, "response.tmpl.html", gin.H{
		"message": "Logged out successfully.",
	})
}
