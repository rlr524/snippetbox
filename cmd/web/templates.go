package main

import (
	"github.com/rlr524/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
)

// Define a templateData type to act as the holding structure for any dynamic data that is passed to
// the HTML templates. Note the model is being referenced directly from the models package here, not
// via the Applications struct defined in main.go. Because of this, the model is referenced using the actual
// model name, Snippet, not the alias it's assigned in the Application struct, which is SnippetsModel.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with the extension "page.gohtml". This
	// essentially provides a slice of all the "page" templates in the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.gohtml"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one by one.
	for _, page := range pages {
		// Extract the file name (e.g. "home.page.gohtml") from the full file path and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the page template file in to a template set.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any "layout" templates to the template set.
		_, err = ts.ParseGlob(filepath.Join(dir, "*.layout.gohtml"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any "partial templates to the template set.
		_, err = ts.ParseGlob(filepath.Join(dir, "*.partial.gohtml"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page as the key
		cache[name] = ts
	}
	// Return the map
	return cache, nil
}
