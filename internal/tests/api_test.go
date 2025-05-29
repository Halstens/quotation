package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositTransaction(t *testing.T) {
	baseURL := os.Getenv("TEST_API_URL")
	url := baseURL + "/api/v1/wallet"
	walletID := "189c100b-9f38-4ef3-8e0e-af9f95e33b53"

	requestBody := map[string]interface{}{
		"walletId":      walletID,
		"operationType": "DEPOSIT",
		"amount":        30,
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
	assert.Equal(t, "success", response["status"])

}

func TestWithdrawTransaction(t *testing.T) {
	baseURL := os.Getenv("TEST_API_URL")
	url := baseURL + "/api/v1/wallet"
	walletID := "189c100b-9f38-4ef3-8e0e-af9f95e33b53"

	requestBody := map[string]interface{}{
		"walletId":      walletID,
		"operationType": "WITHDRAW",
		"amount":        30,
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
	assert.Equal(t, "success", response["status"])
}

func TestShowBalance(t *testing.T) {
	baseURL := os.Getenv("TEST_API_URL")
	existingID := "189c100b-9f38-4ef3-8e0e-af9f95e33b53"
	url := fmt.Sprintf("%s/api/v1/wallets/?id=%s", baseURL, existingID)

	resp, err := http.Get(url)
	require.NoError(t, err, "Ошибка при выполнении запроса")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Неверный статус-код")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Ошибка при чтении ответа")

	t.Logf("Raw response: %s", string(body))

	var result map[string]int64
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, "Ошибка парсинга JSON")

	balance, exists := result["balance"]
	require.True(t, exists, "Ключ 'balance' отсутствует в ответе")
	assert.Equal(t, int64(5750), balance, "Неверное значение баланса")
}
