package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go_hw_2/models"
	"io"
	"net/http"
	"os"
)

type Command struct {
	Host    string
	Port    int
	Action  string
	Name    string
	Amount  int
	NewName string
}

func main() {
	host := flag.String("host", "0.0.0.0", "server host")
	port := flag.Int("port", 8080, "server port")
	action := flag.String("cmd", "", "command to execute: 'create', 'delete', 'update-amount', 'update-name', 'get'")
	name := flag.String("name", "", "name of account")
	amount := flag.Int("amount", 0, "amount of account; required for 'update-amount'")
	newName := flag.String("new-name", "", "new name of account for 'update-name'")

	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage:")
		flag.PrintDefaults()
	}

	flag.Parse()

	cmd := Command{
		Host:    *host,
		Port:    *port,
		Action:  *action,
		Name:    *name,
		Amount:  *amount,
		NewName: *newName,
	}

	if err := executeCommand(&cmd); err != nil {
		panic(err)
	}
}

func executeCommand(cmd *Command) error {
	switch cmd.Action {
	case "create":
		return createAccount(cmd)
	case "delete":
		return deleteAccount(cmd)
	case "update-amount":
		return updateAccountAmount(cmd)
	case "update-name":
		return updateAccountName(cmd)
	case "get":
		return getAccount(cmd)
	default:
		flag.Usage()
		return fmt.Errorf("unknown command: %s", cmd.Action)
	}
}

func createAccount(cmd *Command) error {
	reqBody, err := json.Marshal(models.CreateAccountRequest{Name: cmd.Name})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/account/create", cmd.Host, cmd.Port),
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return checkResponse(resp, http.StatusCreated)
}

func deleteAccount(cmd *Command) error {
	req, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("http://%s:%d/account/delete?name=%s", cmd.Host, cmd.Port, cmd.Name),
		nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return checkResponse(resp, http.StatusOK)
}

func updateAccountAmount(cmd *Command) error {
	reqBody, err := json.Marshal(models.UpdateAmountRequest{Name: cmd.Name, Amount: cmd.Amount})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("http://%s:%d/account/update/amount", cmd.Host, cmd.Port),
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return fmt.Errorf("failed to update account amount: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return checkResponse(resp, http.StatusOK)
}

func updateAccountName(cmd *Command) error {
	reqBody, err := json.Marshal(models.UpdateNameRequest{Name: cmd.Name, NewName: cmd.NewName})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("http://%s:%d/account/update/name", cmd.Host, cmd.Port),
		bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update account name: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return checkResponse(resp, http.StatusOK)
}

func getAccount(cmd *Command) error {
	resp, err := http.Get(
		fmt.Sprintf("http://%s:%d/account?name=%s", cmd.Host, cmd.Port, cmd.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err := checkResponse(resp, http.StatusOK); err != nil {
		return err
	}

	var account models.GetAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	fmt.Printf("Name: %s   Amount: %d\n", account.Name, account.Amount)
	return nil
}

func checkResponse(resp *http.Response, expectedStatus int) error {
	if resp.StatusCode == expectedStatus {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	return fmt.Errorf("unexpected response status: %s", string(body))
}
