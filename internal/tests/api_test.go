package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"

	"github.com/quotes/internal/config"
	"github.com/quotes/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShowQuotesByOne(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	testAuthor := "Confucius"
	testText := "poel, popil, i snova poel"
	url := cfg.TestApiUrl + cfg.ServerPort + "/quotes" + "?author=" + testAuthor

	quote := models.Quote{
		Author: testAuthor,
		Text:   testText,
	}

	resp, err := http.Get(url)
	require.NoError(t, err, "Ошибка при выполнении запроса", url)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Неверный статус-код")

	var actualQuote models.Quote
	err = json.NewDecoder(resp.Body).Decode(&actualQuote)
	require.NoError(t, err, "Ошибка декодирования")

	assert.Equal(t, quote, actualQuote, "Неверный ответ в теле запроса")
}

func TestShowAllQuotes(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	respQuoteOne := models.Quotes{
		ID:     1,
		Author: "Confucius",
		Text:   "poel, popil, i snova poel",
	}
	respQuoteSecond := models.Quotes{
		ID:     2,
		Author: "Kto to tam",
		Text:   "pospal, vstal, poel i snova pospal",
	}

	var quotes []models.Quotes
	quotes = append(quotes, respQuoteOne, respQuoteSecond)

	url := cfg.TestApiUrl + cfg.ServerPort + "/quotes"
	resp, err := http.Get(url)
	require.NoError(t, err, "Ошибка при выполнении запроса", url)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Неверный статус-код")

	var actualQuotes []models.Quotes
	err = json.NewDecoder(resp.Body).Decode(&actualQuotes)
	require.NoError(t, err, "Ошибка декодирования")

	assert.Equal(t, quotes, actualQuotes, "Неверный ответ в теле запроса")

}

func TestShowRandomQuotes(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	url := cfg.TestApiUrl + cfg.ServerPort + "/quotes/random"
	resp, err := http.Get(url)
	require.NoError(t, err, "Ошибка при выполнении запроса", url)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Неверный статус-код")

	var actualQuote models.Quote
	err = json.NewDecoder(resp.Body).Decode(&actualQuote)
	require.NoError(t, err, "Ошибка декодирования")

	assert.IsType(t, "", actualQuote.Author)
	assert.IsType(t, "", actualQuote.Text)
}

func TestAddNewQuote(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	url := cfg.TestApiUrl + cfg.ServerPort + "/quotes"

	requestBody := map[string]interface{}{
		"author": "Trump",
		"quote":  "I like Mexico!",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "Success!", response["status"])
}

func TestDeleteQuote(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	testId := "3"
	url := cfg.TestApiUrl + cfg.ServerPort + "/quotes/" + testId
	req, err := http.NewRequest("DELETE", url, nil)
	require.NoError(t, err, url)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, url)
	defer resp.Body.Close()

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "success delete", response["status"])
}
