package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/Aniket52kr/GO-Assignment/internal"
	socials "github.com/Aniket52kr/GO-Assignment/internal/auth"
	"github.com/Aniket52kr/GO-Assignment/middleware"
	"github.com/Aniket52kr/GO-Assignment/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl.html", nil)
}

func notFound(c *gin.Context) {
	c.HTML(http.StatusNotFound, "error.tmpl.html", gin.H{
		"error":   "404 Not Found",
		"message": "The requested page was not found.",
	})
}

func main() {
	godotenv.Load(".env")
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.RedirectTrailingSlash = true
	app.HandleMethodNotAllowed = true
	app.NoRoute(notFound)

	app.Static("/static", "./static")
	app.SetFuncMap(template.FuncMap{
		"formatAsTitle": internal.FormatAsTitle,
		"formatAsDate":  internal.FormatAsDate,
	})
	app.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte(os.Getenv("SECRET_KEY")))
	app.Use(sessions.Sessions("SocialEcho", store))
	app.Use(middleware.RecoveryMiddleware())

	// public routes:-
	app.GET("/", index)
	app.GET("/signup", routes.SignUp)
	app.GET("/login", routes.Login)
	app.GET("/logout", routes.Logout)
	app.GET("/feed", middleware.AuthMiddleware(), routes.UserFeed)
	app.GET("/feed/more", middleware.AuthMiddleware(), routes.LoadMoreFeed)

	// Oauth and verification routes:-
	auth := app.Group("/auth")
	{

		auth.GET("/signup/github", socials.GitHubSignUp)
		auth.GET("/signup/google", socials.GoogleSignUp)
		auth.GET("/login/github", socials.GitHubLogin)
		auth.GET("/login/google", socials.GoogleLogin)

		auth.GET("/github", socials.GitHubAuth)
		auth.GET("/google", socials.GoogleAuth)
		auth.GET("/verify", middleware.AuthMiddleware(), routes.SendVerificationMail)
		auth.GET("/verify/:id", routes.Verify)

		auth.POST("/signup", routes.SignUp)
		auth.POST("/login", routes.Login)
	}

	// user group routes:-
	user := app.Group("/user")
	user.GET("/:username", routes.GetUserByName)
	user.GET("/:username/posts", routes.GetUserPosts)
	user.GET("/:username/posts/more", routes.LoadMorePosts)
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/", routes.GetUser)
		user.GET("/settings/avatar", routes.UpdateAvatar)
		user.GET("/settings/username", routes.UpdateUsername)
		user.GET("/settings/password", routes.UpdatePassword)
		user.GET("/settings/delete", routes.DeleteUser)

		user.POST("/:username/toggle-follow", routes.ToggleFollow)
		user.POST("/settings/avatar", routes.UpdateAvatar)
		user.POST("/settings/username", routes.UpdateUsername)
		user.POST("/settings/password", routes.UpdatePassword)
		user.POST("/settings/delete", routes.DeleteUser)
	}

	// search group routes:-
	search := app.Group("/search")
	{
		search.GET("/", routes.SearchUser)
		search.GET("/more", routes.LoadMoreUsers)

		search.POST("/", routes.SearchUser)
		search.POST("/:username/toggle-follow", middleware.AuthMiddleware(), routes.ToggleSearchFollow)
	}

	// post group routes:-
	post := app.Group("/post")
	post.GET("/:id", routes.GetPost)
	post.Use(middleware.AuthMiddleware())
	{
		post.GET("/", routes.NewPost)
		post.GET("/:id/toggle-vote", routes.ToggleVote)
		post.GET("/:id/delete", routes.DeletePost)
		post.GET("/:id/comments", routes.LoadMoreComments)
		post.GET("/:id/comment/delete", routes.DeleteComment)

		post.POST("/", routes.NewPost)
		post.POST("/:id/comment", routes.Comment)
	}

	// Load custom port from .env or fallback to 8081
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := app.Run(":" + port); err != nil {
		panic(err)
	}
}
