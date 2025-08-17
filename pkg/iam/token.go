package iam

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// YandexCloudIAMURL URL для получения IAM токена
	YandexCloudIAMURL = "https://iam.api.cloud.yandex.net/iam/v1/tokens"
	// TokenTTL время жизни токена (по умолчанию 12 часов)
	TokenTTL = 12 * time.Hour
)

// IAMTokenRequest структура запроса для получения IAM токена
type IAMTokenRequest struct {
	JWT string `json:"jwt"`
}

// IAMTokenResponse структура ответа с IAM токеном
type IAMTokenResponse struct {
	IAMToken  string    `json:"iamToken"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// OAuthTokenRequest структура запроса для получения IAM токена через OAuth токен
type OAuthTokenRequest struct {
	YandexPassportOAuthToken string `json:"yandexPassportOauthToken"`
}

// ServiceAccountKey структура для ключа сервисного аккаунта
type ServiceAccountKey struct {
	ID               string `json:"id"`
	ServiceAccountID string `json:"service_account_id"`
	CreatedAt        string `json:"created_at"`
	KeyAlgorithm     string `json:"key_algorithm"`
	PublicKey        string `json:"public_key"`
	PrivateKey       string `json:"private_key"`
}

// TokenManager менеджер для работы с IAM токенами
type TokenManager struct {
	serviceAccountKey *ServiceAccountKey
	oauthToken        string
	currentToken      string
	expiresAt         time.Time
	client            *http.Client
}

// NewTokenManagerWithServiceAccount создает новый менеджер токенов с ключом сервисного аккаунта
func NewTokenManagerWithServiceAccount(key *ServiceAccountKey) *TokenManager {
	return &TokenManager{
		serviceAccountKey: key,
		client:            &http.Client{Timeout: 10 * time.Second},
	}
}

// NewTokenManagerWithOAuth создает новый менеджер токенов с OAuth токеном
func NewTokenManagerWithOAuth(oauthToken string) *TokenManager {
	return &TokenManager{
		oauthToken: oauthToken,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// GetToken возвращает действующий IAM токен
func (tm *TokenManager) GetToken() (string, error) {
	// Проверяем, нужно ли обновить токен
	if tm.currentToken == "" || time.Now().Add(5*time.Minute).After(tm.expiresAt) {
		if err := tm.refreshToken(); err != nil {
			return "", fmt.Errorf("failed to refresh IAM token: %w", err)
		}
	}

	return tm.currentToken, nil
}

// refreshToken обновляет IAM токен
func (tm *TokenManager) refreshToken() error {
	if tm.serviceAccountKey != nil {
		return tm.refreshTokenWithServiceAccount()
	} else if tm.oauthToken != "" {
		return tm.refreshTokenWithOAuth()
	}

	return fmt.Errorf("no authentication method available")
}

// refreshTokenWithServiceAccount обновляет токен используя ключ сервисного аккаунта
func (tm *TokenManager) refreshTokenWithServiceAccount() error {
	// Создаем JWT токен
	jwtToken, err := tm.createJWT()
	if err != nil {
		return fmt.Errorf("failed to create JWT: %w", err)
	}

	// Делаем запрос к API
	request := IAMTokenRequest{JWT: jwtToken}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := tm.client.Post(YandexCloudIAMURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var response IAMTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	tm.currentToken = response.IAMToken
	tm.expiresAt = response.ExpiresAt

	return nil
}

// refreshTokenWithOAuth обновляет токен используя OAuth токен
func (tm *TokenManager) refreshTokenWithOAuth() error {
	request := OAuthTokenRequest{YandexPassportOAuthToken: tm.oauthToken}
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := tm.client.Post(YandexCloudIAMURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	var response IAMTokenResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	tm.currentToken = response.IAMToken
	tm.expiresAt = response.ExpiresAt

	return nil
}

// createJWT создает JWT токен для сервисного аккаунта
func (tm *TokenManager) createJWT() (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"iss": tm.serviceAccountKey.ServiceAccountID,
		"aud": "https://iam.api.cloud.yandex.net/iam/v1/tokens",
		"iat": now.Unix(),
		"exp": now.Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	token.Header["kid"] = tm.serviceAccountKey.ID

	// Парсим приватный ключ
	privateKey, err := tm.parsePrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// Подписываем токен
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// parsePrivateKey парсит приватный ключ из PEM формата
func (tm *TokenManager) parsePrivateKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(tm.serviceAccountKey.PrivateKey))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA private key")
	}

	return rsaKey, nil
}

// LoadServiceAccountKeyFromFile загружает ключ сервисного аккаунта из файла
func LoadServiceAccountKeyFromFile(filename string) (*ServiceAccountKey, error) {
	data, err := readFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var key ServiceAccountKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key: %w", err)
	}

	return &key, nil
}

// LoadServiceAccountKeyFromJSON загружает ключ сервисного аккаунта из JSON строки
func LoadServiceAccountKeyFromJSON(jsonData string) (*ServiceAccountKey, error) {
	var key ServiceAccountKey
	if err := json.Unmarshal([]byte(jsonData), &key); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key: %w", err)
	}

	return &key, nil
}

// readFile читает файл с диска
func readFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
