package handlers

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"

	"github.com/charmbracelet/log"
)

// generateHandler is a http.HandlerFunc that serves a form for generating a list of tag IDs
// This is accessible via GET and POST requests to /generate
// GET requests will return the form
// POST requests will generate and return a list of tag IDs
// Tag IDs are numbers of a specified length
func generateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		generateForm(w)
	case http.MethodPost:
		generateTags(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// generateForm serves the form for generating tag IDs
func generateForm(w http.ResponseWriter) {
	w.Write([]byte(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Generate tags</title>
	</head>
	<body>

	<h1>Generate tags</h1>

	<form method="post">
		<label for="length">Length:</label><br>
		<input type="number" id="length" name="length" min="1" step="1" required><br>
		<label for="count">Count:</label><br>
		<input type="number" id="count" name="count" min="1" step="1" required><br>
		<input type="submit" value="Generate">
	</form>

	</body>
	</html>
	`))
}

// generateTags generates and returns a list of tag IDs
func generateTags(w http.ResponseWriter, r *http.Request) {
	length := r.FormValue("length")
	if length == "" {
		http.Error(w, "Length not specified", http.StatusBadRequest)
		return
	}

	count := r.FormValue("count")
	if count == "" {
		http.Error(w, "Count not specified", http.StatusBadRequest)
		return
	}

	// Convert length and count to int
	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		http.Error(w, "Invalid length", http.StatusBadRequest)
		return
	}
	countInt, err := strconv.Atoi(count)
	if err != nil {
		http.Error(w, "Invalid count", http.StatusBadRequest)
		return
	}

	tags := numberGenerator(lengthInt, countInt)

	// Write tags as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// numberGenerator generates a list of tag IDs
// The length of the tag IDs is specified by the length parameter
// The count parameter specifies the number of tag IDs to generate
// The tag IDs are random numbers
// The tag IDs must be unique and returned in numerical order
func numberGenerator(length, count int) []int {
	log.Printf("Generating tags, length: %d, count: %d", length, count)
	if length < 1 || count < 1 || length > 9 {
		return nil
	}

	subInt := int(math.Pow10(length - 1))
	maxInt := int(math.Pow10(length)) - 1
	rangeSize := maxInt - subInt + 1

	if count > rangeSize {
		log.Printf("Error: Count exceeds the range of possible unique values for the specified length")
		return nil
	}

	tags := make(map[int]struct{})

	for len(tags) < count {
		tag := rand.Intn(rangeSize) + subInt
		tags[tag] = struct{}{}
	}

	tagList := make([]int, 0, count)
	for tag := range tags {
		tagList = append(tagList, tag)
	}

	sort.Ints(tagList)

	return tagList
}
