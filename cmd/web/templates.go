package main

import (
	"github.com/rlr524/snippetbox/pkg/forms"
	"github.com/rlr524/snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// Define a templateData type to act as the holding structure for any dynamic data that is passed to
// the HTML templates. Note the model is being referenced directly from the models package here, not
// via the Applications struct defined in main.go. Because of this, the model is referenced using the actual
// model name, Snippet, not the alias it's assigned in the Application struct, which is SnippetsModel.
type templateData struct {
	CurrentYear     int
	Toast           string
	Form            *forms.Form
	IsAuthenticated bool
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
}

// The humanDate function converts ISO 8601 time format to a more human-readable form. See the documentation
// for the time package for acceptable string formats. A custom template function like this is variadic, it can
// take in as many parameters as needed, but it can only return one value (and an error as needed).
func humanDate(t time.Time) string {
	return t.Format("January 02, 2006 at 15:04")
}

// The function map object acts as a lookup between the names of custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
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

		// Parse the page template file in to a template set. In doing this, also register the function map by creating
		// an empty template set with template.New(), using the Funcs() method to register the function map
		// and then parsing the file. These methods can be chained together as here.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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
