package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templates embed.FS

// ShowMetrics is an HTTP handler that shows all the metrics
func ShowMetrics(s Storager, tempFile string) gin.HandlerFunc {
	// Parse the template file
	tmpl, errFunc := ParseTemplate(tempFile)
	if errFunc != nil {
		return func(c *gin.Context) {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	// Return the actual handler function
	return func(c *gin.Context) {
		
		data := s.GetAllMetrics()

		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Execute the template with the data
		err := tmpl.Execute(c.Writer, data)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// ParseTemplate parses the template file and returns the parsed template
func ParseTemplate(tempFile string) (*template.Template, func(*gin.Context)) {
	tmpl, err := template.ParseFS(templates, filepath.Join("templates", tempFile))
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		return nil, func(c *gin.Context) {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}
	return tmpl, nil
}
