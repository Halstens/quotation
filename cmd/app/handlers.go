package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/quotes/internal/models"
)

func (app *application) showQuotes(w http.ResponseWriter, r *http.Request) {
	authorRequest := r.URL.Query().Get("author")
	if len([]rune(authorRequest)) == 0 {
		data, err := app.quotes.GetQuotes()
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
				fmt.Println("Не найдено")
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)

	} else {
		author, text, err := app.quotes.GetQuoteByAuthor(authorRequest)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.notFound(w)
			} else {
				app.serverError(w, err)
				fmt.Println("Не найдено")
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"author": author,
			"quote":  text,
		})

	}

}

func (app *application) showRandomQuotes(w http.ResponseWriter, r *http.Request) {
	author, text, err := app.quotes.GetRundomQuote()
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"author": author,
		"quote":  text,
	})

}

func (app *application) AddNewQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow: ", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	var request models.Quotes
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len([]rune(request.Author)) == 0 || len([]rune(request.Text)) == 0 {
		http.Error(w, "Error of valid author and text", http.StatusBadRequest)
		return
	}
	id, status, err := app.quotes.CreateQuote(request.Author, request.Text)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     strconv.Itoa(id),
		"status": status,
	})
}

func (app *application) deleteQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow: ", http.MethodDelete)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/quotes/")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	status, err := app.quotes.DeleteQuotesRec(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     strconv.Itoa(id),
		"status": status,
	})

}
