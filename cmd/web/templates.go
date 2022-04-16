package main

import "github.com/rlr524/snippetbox/pkg/models"

// Define a templateData type to act as the holding structure for any dynamic data that is passed to
// the HTML templates.
type templateData struct {
	Snippet *models.Snippet
}
