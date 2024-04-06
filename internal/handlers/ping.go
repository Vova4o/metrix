package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Connection to the database successful"})
	}
}
