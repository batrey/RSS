package handlers

import (
	"strconv"
)

const AddForm = `
<form method="GET" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>`

// Checks to see if the string number is a valid number
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Uses the abbreviation send from user and matches it with the Category in the DB
func CategoryCheck(category string) string {
	switch category {
	case "bbc":
		return "BBC News - UK"
	case "bbc-tech":
		return "BC News - Technology"
	case "sky":
		return "K News - The latest headlines from the UK | Sky News"
	case "sky-tech":
		return "Tech News - Latest Technology and Gadget News | Sky News"
	default:
		return "No category found"
	}
}
