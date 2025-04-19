package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/barelydead/rested/storage"
	"github.com/manifoldco/promptui"
)

func SetRequestMethod(req *storage.RestedRequest) {
	prompt := promptui.Select{
		Label: "HTTP Method",
		Items: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v", err)
	}

	req.Method = result
}

func SetRequestUrl(req *storage.RestedRequest) {
	prompt := promptui.Prompt{
		Label: "Request URL",
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Invalid request URL: %v", err)
	}

	req.URL = result
}

func SetRequestName(req *storage.RestedRequest) {
	prompt := promptui.Prompt{
		Label: "Request name",
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Invalid request URL: %v", err)
	}

	req.RequestName = result
}

func SetBody(req *storage.RestedRequest) {
	// Use system temp dir
	tmpfile, err := os.CreateTemp("", "rested-body-*.txt")
	if err != nil {
		log.Printf("Could not create temp file: %v", err)
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up

	// Preload existing body (if any)
	if req.Body != "" {
		tmpfile.WriteString(req.Body)
		tmpfile.Sync()
	}
	tmpfile.Close()

	// Use user's preferred editor or fallback
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // or "vi" or "notepad" for Windows
	}

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Editor failed: %v", err)
		return
	}

	// Read file content back
	editedBody, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		log.Printf("Failed to read edited body: %v", err)
		return
	}

	req.Body = strings.TrimSpace(string(editedBody))
}

func SetHeader(req *storage.RestedRequest) {
	keyPrompt := promptui.Prompt{
		Label: "Header Key (e.g. Content-Type)",
	}

	key, err := keyPrompt.Run()
	if err != nil || strings.TrimSpace(key) == "" {
		log.Println("Header key input cancelled or invalid")
		return
	}

	valuePrompt := promptui.Prompt{
		Label: fmt.Sprintf("Value for '%s'", key),
	}

	value, err := valuePrompt.Run()
	if err != nil {
		log.Println("Header value input cancelled or invalid")
		return
	}

	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	req.Headers[key] = value
}

func SendRequest(req *storage.RestedRequest) {
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	var body io.Reader
	if req.Body != "" {
		body = bytes.NewBufferString(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, parsedURL.String(), body)
	if err != nil {
		log.Fatalf("Failed to build request: %v", err)
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	PrintRequest(resp)
}

func saveRequestToCollection(db *storage.DB, req *storage.RestedRequest) {
	if len(db.Collections) == 0 {
		fmt.Println("⚠️ No collections found. Please create one first.")
		return
	}

	if req.RequestName == "" {
		fmt.Println("⚠️ Request must have a request name to be saved. Please create one first.")
		prompt := promptui.Prompt{
			Label: "Press enter to continue.",
		}

		prompt.Run()
		return
	}

	titles := []string{}
	for _, col := range db.Collections {
		titles = append(titles, col.Title)
	}

	prompt := promptui.Select{
		Label: "Save request to collection",
		Items: titles,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	db.Collections[index].Requests = append(db.Collections[index].Requests, *req)

	fmt.Printf("✅ Request '%s' added to collection '%s'\n", req.RequestName, db.Collections[index].Title)
}

func NewRequest(db *storage.DB, request storage.RestedRequest) {
	for {
		prompt := promptui.Select{
			Label: request.RequestName,
			Items: []string{
				"Send request",
				"Set request name",
				"Set method",
				"Set request URL",
				"Set headers",
				"Set body",
				"Save to collection",
				"Back",
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		switch result {
		case "Set request name":
			SetRequestName(&request)
		case "Set method":
			SetRequestMethod(&request)
		case "Set request URL":
			SetRequestUrl(&request)
		case "Set headers":
			SetHeader(&request)
		case "Set body":
			SetBody(&request)
		case "Send request":
			SendRequest(&request)
		case "Save to collection":
			saveRequestToCollection(db, &request)
		case "Back":
			Root(db)
		}
	}
}

func PrintRequest(resp *http.Response) {
	fmt.Println("\n------ Response ------")
	fmt.Printf("Status: %s\n", resp.Status)

	fmt.Println("Headers:")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	fmt.Println("Body:")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("  Error reading body: %v\n", err)
	} else {
		fmt.Printf("  %s\n", string(body))
	}
	fmt.Println("----------------------")
}
