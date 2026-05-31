package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func main() {
	baseURL := flag.String("base", "http://localhost:8080", "API base URL")
	tokenFile := flag.String("token-file", ".night-hawk-token", "path to saved JWT token")
	token := flag.String("token", "", "JWT token override")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	cmd := args[0]

	var err error
	switch cmd {
	case "health":
		err = request(client, http.MethodGet, *baseURL+"/api/v1/health", "", nil, false)
	case "products":
		err = request(client, http.MethodGet, *baseURL+"/api/v1/products", "", nil, true)
	case "register":
		username, password, parseErr := usernamePassword(args)
		if parseErr != nil {
			err = parseErr
			break
		}
		err = requestJSON(client, http.MethodPost, *baseURL+"/api/v1/users/register", "", map[string]string{
			"username": username,
			"password": password,
		}, false)
	case "login":
		username, password, parseErr := usernamePassword(args)
		if parseErr != nil {
			err = parseErr
			break
		}
		err = login(client, *baseURL, *tokenFile, username, password)
	case "me":
		authToken := strings.TrimSpace(*token)
		if authToken == "" {
			authToken, err = readToken(*tokenFile)
			if err != nil {
				break
			}
		}
		err = request(client, http.MethodGet, *baseURL+"/api/v1/users/me", authToken, nil, false)
	case "orders":
		authToken := strings.TrimSpace(*token)
		if authToken == "" {
			authToken, err = readToken(*tokenFile)
			if err != nil {
				break
			}
		}
		productID, quantity, parseErr := orderArgs(args)
		if parseErr != nil {
			err = parseErr
			break
		}
		err = requestJSON(client, http.MethodPost, *baseURL+"/api/v1/orders", authToken, map[string]int64{
			"product_id": productID,
			"quantity":   quantity,
		}, true)
	case "me-wrong":
		err = request(client, http.MethodGet, *baseURL+"/api/v1/users/me", "wrong-token", nil, false)
	case "token":
		var saved string
		saved, err = readToken(*tokenFile)
		if err == nil {
			fmt.Println(saved)
		}
	default:
		err = fmt.Errorf("unknown command: %s", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Night-Hawk API test tool

Usage:
  go run ./cmd/apitest health
  go run ./cmd/apitest products
  go run ./cmd/apitest register <username> <password>
  go run ./cmd/apitest login <username> <password>
  go run ./cmd/apitest me
  go run ./cmd/apitest me-wrong
  go run ./cmd/apitest orders [product_id quantity]
  go run ./cmd/apitest token

Options:
  -base       API base URL, default http://localhost:8080
  -token      JWT token override for protected commands
  -token-file saved token path, default .night-hawk-token

Example:
	go run ./cmd/apitest health
	go run ./cmd/apitest products
	go run ./cmd/apitest register TopBeauty 200506
	go run ./cmd/apitest login TopBeauty 200506
	go run ./cmd/apitest me
	go run ./cmd/apitest orders
	go run ./cmd/apitest orders 1 2`)
}

func usernamePassword(args []string) (string, string, error) {
	if len(args) != 3 {
		return "", "", fmt.Errorf("usage: %s <username> <password>", args[0])
	}

	return args[1], args[2], nil
}

func orderArgs(args []string) (int64, int64, error) {
	if len(args) == 1 {
		return 1, 2, nil
	}
	if len(args) != 3 {
		return 0, 0, fmt.Errorf("usage: %s [product_id quantity]", args[0])
	}

	productID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid product_id: %w", err)
	}

	quantity, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid quantity: %w", err)
	}

	return productID, quantity, nil
}

func requestJSON(client *http.Client, method, url, token string, body any, failOnHTTPError bool) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return request(client, method, url, token, bytes.NewReader(payload), failOnHTTPError)
}

func request(client *http.Client, method, url, token string, body io.Reader, failOnHTTPError bool) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("HTTP %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	printBody(responseBody)

	if failOnHTTPError && resp.StatusCode >= 400 {
		return fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	return nil
}

func login(client *http.Client, baseURL, tokenFile, username, password string) error {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/v1/users/login", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("HTTP %d %s\n", resp.StatusCode, http.StatusText(resp.StatusCode))
	printBody(responseBody)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	var parsed apiResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return err
	}

	var data struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(parsed.Data, &data); err != nil {
		return err
	}
	if data.Token == "" {
		return errors.New("login response does not contain token")
	}

	if err := os.WriteFile(tokenFile, []byte(data.Token), 0600); err != nil {
		return err
	}

	fmt.Printf("\nSaved token to %s\n", tokenFile)
	return nil
}

func readToken(tokenFile string) (string, error) {
	tokenBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("cannot read token from %s; run login first or pass -token: %w", tokenFile, err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return "", fmt.Errorf("token file %s is empty", tokenFile)
	}

	return token, nil
}

func printBody(body []byte) {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, body, "", "  "); err == nil {
		fmt.Println(pretty.String())
		return
	}

	fmt.Println(string(body))
}
