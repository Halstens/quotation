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

// func (wr *WalletRepository) UpdateBalance(id uuid.UUID, amount int, opType models.OperationType) error {
// 	var query string
// 	if amount <= 0 {
// 		return fmt.Errorf("amount must be positive")
// 	}

// 	tx, err := wr.DB.Begin()
// 	if err != nil {
// 		return fmt.Errorf("fail transaction: %w", err)
// 	}
// 	fmt.Println("trans")

// 	defer tx.Rollback()

// 	fmt.Println("check ok")
// 	switch opType {
// 	case "DEPOSIT":
// 		query = "UPDATE wallets SET balance = balance + $1 WHERE id = $2 RETURNING balance"
// 	case "WITHDRAW":
// 		query = "UPDATE wallets SET balance = balance - $1 WHERE id = $2 AND balance >= $1 RETURNING balance"
// 	default:
// 		return fmt.Errorf("invalid op type")
// 	}
// 	fmt.Println("case")

// 	var newBalance int
// 	err = tx.QueryRow(query, amount, id).Scan(&newBalance)
// 	if err != nil {
// 		if err == sql.ErrNoRows && opType == "WITHDRAW" {
// 			return fmt.Errorf("insufficient funds")
// 		}
// 		return fmt.Errorf("failed to update balance: %w", err)
// 	}
// 	fmt.Println("get balance", newBalance)

// 	if err := tx.Commit(); err != nil {
// 		return fmt.Errorf("failed to commit transaction: %w", err)
// 	}
// 	fmt.Println("commit")
// 	return err
// }

// func (wr *WalletRepository) GetBalance(id string) (int64, error) {
// 	var balance int64
// 	err := wr.DB.QueryRow("SELECT balance FROM wallets WHERE id = $1", id).Scan(&balance)
// 	return balance, err
// }

// func (wr *WalletRepository) UpdateBalanceWithRetry(id uuid.UUID, amount int, opType models.OperationType, maxRetries int) error {
// 	var lastErr error

// 	for i := 0; i < maxRetries; i++ {
// 		err := wr.UpdateBalance(id, amount, opType)
// 		if err == nil {
// 			return nil
// 		}

// 		// Если словили дедлок
// 		var pgErr *pq.Error
// 		if errors.As(err, &pgErr) && pgErr.Code == "40P01" {
// 			lastErr = err
// 			delay := time.Duration(math.Pow(2, float64(i))) * time.Millisecond
// 			delay += time.Duration(rand.Intn(100)) * time.Millisecond

// 			time.Sleep(delay)
// 			continue
// 		}

// 		return err
// 	}

// 	return fmt.Errorf("max retries (%d) reached, last error: %v", maxRetries, lastErr)
// }
