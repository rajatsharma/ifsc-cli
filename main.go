package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Bank struct {
	Bank    string `json:"BANK"`
	Branch  string `json:"BRANCH"`
	City    string `json:"CITY"`
	Address string `json:"ADDRESS"`
}

func request(ctx context.Context, url string, response *Bank) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable make request %w", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to get response %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			println("unable to close body")
		}
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read response %w", err)
	}
	err = json.Unmarshal(b, response)
	if err != nil {
		return fmt.Errorf("unable to parse response %w", err)
	}

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// trap SIGINT, wait to trigger shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	var response Bank
	url := fmt.Sprintf("https://ifsc.razorpay.com/%s", "ICIC0006955")
	if err := request(ctx, url, &response); err != nil {
		println("unable to fetch")
	}

	fmt.Printf("Bank: %s\nBranch: %s\nCity: %s\nAddress: %s",
		response.Bank, response.Branch, response.City, response.Address)

	// Trap signals and shutdown gracefully
	go func() {
		<-signals
		cancel()
	}()
}
