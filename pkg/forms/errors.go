package forms

// The Errors type is used to hold the validation error messages for forms. The name of the
// form field is used as the key.
type Errors map[string][]string

// The Add method adds error messages for a given field to the map. It has a function receiver defined on the
// errors type and names it e. It takes in two string params named field and message and appends a message to
// the Errors map using the field name as the key. It takes a value receiver to Errors because it's only adding new
// elements to the map, it is not modifying any elements.
func (e Errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// The Get method retrieves the first error message for a given field from the map.
func (e Errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
