package main

import "github.com/rlr524/snippetbox/pkg/models"

// Define a templateData type to act as the holding structure for any dynamic data that is passed to
// the HTML templates. Note the model is being referenced directly from the models package here, not
// via the Applications struct defined in main.go. Because of this, the model is referenced using the actual
// model name, Snippet, not the alias it's assigned in the Application struct, which is SnippetsModel.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
