package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/quotes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.showQuotes(w, r) //get + filter
		case http.MethodPost:
			app.AddNewQuote(w, r) // post
		default:
			w.Header().Set("Allow", "GET or POST")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/quotes/random", app.showRandomQuotes) //get random
	mux.HandleFunc("/quotes/", app.deleteQuote)            //del from id

	return mux
}
