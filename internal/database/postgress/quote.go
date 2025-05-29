package postgress

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/quotes/internal/models"
)

type QuotesRepository struct {
	DB *sqlx.DB
}

func (qr *QuotesRepository) GetRundomQuote() (string, string, error) {
	var author, text string
	err := qr.DB.QueryRow("SELECT author, text FROM quotes ORDER BY random() LIMIT 1").Scan(&author, &text)
	return author, text, err
}

func (qr *QuotesRepository) GetQuoteByAuthor(aut string) (string, string, error) {
	var author, text string
	err := qr.DB.QueryRow("SELECT author, text FROM quotes WHERE author = $1", aut).Scan(&author, &text)
	return author, text, err
}

func (qr *QuotesRepository) GetQuotes() ([]models.Quotes, error) {
	rows, err := qr.DB.Query("SELECT * FROM quotes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quotes []models.Quotes
	for rows.Next() {
		var q models.Quotes
		if err := rows.Scan(&q.ID, &q.Author, &q.Text); err != nil {
			return nil, err
		}
		quotes = append(quotes, q)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return quotes, err
}

func (qr *QuotesRepository) CreateQuote(author, text string) (int, string, error) {
	var id int
	err := qr.DB.QueryRow("INSERT INTO quotes (author, text) VALUES ($1, $2) RETURNING id", author, text).Scan(&id)
	if err != nil {
		return 0, "denied", err
	}

	return id, "Success!", err
}

func (qr *QuotesRepository) DeleteQuotesRec(id int) (string, error) {
	result, err := qr.DB.Exec("DELETE FROM quotes WHERE id = $1", id)
	if err != nil {

		return "denied", fmt.Errorf("deleted in database error %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {

		return "denied", fmt.Errorf("delete not success %w", err)
	}

	return "success delete", err

}
