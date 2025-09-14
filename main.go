package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	return db, err
}

type PingResponse struct {
	Message string `json:"message"`
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func PingHandler(c *gin.Context) {
	response := PingResponse{Message: "pong"}
	c.JSON(http.StatusOK, response)
}

func main() {
	r := gin.Default()

	r.GET("/ping", PingHandler)
	r.POST("/users", Create)

	r.Run("test_cicd:50051")
}

func Create(c *gin.Context) {
	db, err := initDB()
	if err != nil {
		return
	}

	defer db.Close()

	var userRes User
	err = c.ShouldBindJSON(&userRes)
	if err != nil {
		println(err)
		return
	}

	query := `INSERT INTO user (id, name, email) VALUES ($1, $2, $3)`

	_, err = db.Exec(query, userRes.ID, userRes.Name, userRes.Email)
	if err != nil {
		return
	}

	println("User name: ", userRes.Name, "\n Email: ", userRes.Email)

}
