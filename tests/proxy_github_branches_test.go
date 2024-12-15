package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"testing"

	"github.com/dxta-dev/app/mock/har"
)

var ts *httptest.Server
var serverShutdown func()

var repoBranchPattern = regexp.MustCompile(`^/repos/([^/]+)/([^/]+)/branches$`)

func createRepoBranchesHandler(page int) func(http.ResponseWriter, *http.Request) {
	harFilePaths := []string{
		"./github/branches/branches-200-page-1.json",
		"./github/branches/branches-200-page-2.json",
	}

	harMaps, err := har.Load(harFilePaths, 1)
	if err != nil {
		panic("Failed to load HAR files")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")

		pg := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil {
				pg = p
			} else {
				http.Error(w, "Invalid page number", http.StatusBadRequest)
				return
			}
		}

		if page == pg-1 {
			var h = harMaps[harFilePaths[page]]
			har.EntryToResponse(w, &h.Log.Entries[0], map[string]string{"Content-Encoding": "deflate"})
			return
		}

		http.NotFound(w, r)

	}
}

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/torvalds/linux/branches", createRepoBranchesHandler(0))

	ts = httptest.NewServer(mux)

	_, cancel := context.WithCancel(context.Background())
	serverShutdown = cancel

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		fmt.Println("Received shutdown signal")
		cancel()
	}()

	exitCode := m.Run()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ts.Close()
	}()

	wg.Wait()

	os.Exit(exitCode)
}

func TestBranches(t *testing.T) {
	resp, err := http.Get(ts.URL + "/repos/torvalds/linux/branches")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Status Code:", resp.Status)

	// Print the headers
	fmt.Println("Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestBranchesPage1(t *testing.T) {
	resp, err := http.Get(ts.URL + "/repos/torvalds/linux/branches?page=1")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Status Code:", resp.Status)

	// Print the headers
	fmt.Println("Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestBranchesPage2(t *testing.T) {
	resp, err := http.Get(ts.URL + "/repos/torvalds/linux/branches?page=2")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Status Code:", resp.Status)

	// Print the headers
	fmt.Println("Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
