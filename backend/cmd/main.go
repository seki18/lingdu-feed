package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	connStr := "host=localhost port=5432 user=admin password=password dbname=community sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var user User

			rows.Scan(&user.ID, &user.Name)

			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	})

	r.Run(":18080")
}
