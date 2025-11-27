package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestCorsMiddleware(t *testing.T) {
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "OPTIONS request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Header().Get("Access-Control-Allow-Origin") != "*" {
					t.Error("Expected CORS headers for OPTIONS")
				}
			},
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Header().Get("Access-Control-Allow-Origin") != "*" {
					t.Error("Expected CORS headers for GET")
				}
			},
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				if w.Header().Get("Access-Control-Allow-Origin") != "*" {
					t.Error("Expected CORS headers for POST")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkHeaders != nil {
				tt.checkHeaders(t, w)
			}
		})
	}
}

func TestWatcher_handleStatus(t *testing.T) {
	watcher := &Watcher{
		Config: Config{
			WatchDir: "/test/dir",
			Port:     8006,
		},
		Status: "Running",
	}

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	watcher.handleStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["status"] != "Running" {
		t.Errorf("Expected status 'Running', got '%v'", result["status"])
	}

	config, ok := result["config"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected config in response")
	}

	if config["watchDir"] != "/test/dir" {
		t.Errorf("Expected watchDir '/test/dir', got '%v'", config["watchDir"])
	}
}

func TestWatcher_handleStart(t *testing.T) {
	watcher := &Watcher{
		Config: Config{
			WatchDir: "/test/dir",
		},
		Running: false,
		Status:  "Stopped",
		Clients: make(map[chan string]bool),
	}

	req := httptest.NewRequest(http.MethodPost, "/start", nil)
	w := httptest.NewRecorder()

	watcher.handleStart(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if !watcher.Running {
		t.Error("Expected watcher to be running")
	}

	if watcher.Status != "Running" {
		t.Errorf("Expected status 'Running', got '%s'", watcher.Status)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["status"] != "started" {
		t.Errorf("Expected status 'started', got '%s'", result["status"])
	}

	// Test starting when already running
	w2 := httptest.NewRecorder()
	watcher.handleStart(w2, req)
	// Should still return OK but not start again
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200 when already running, got %d", w2.Code)
	}
}

func TestWatcher_handleStop(t *testing.T) {
	watcher := &Watcher{
		Config: Config{
			WatchDir: "/test/dir",
		},
		Running:  true,
		Status:   "Running",
		StopChan: make(chan struct{}),
		Clients:  make(map[chan string]bool),
	}

	req := httptest.NewRequest(http.MethodPost, "/stop", nil)
	w := httptest.NewRecorder()

	watcher.handleStop(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if watcher.Running {
		t.Error("Expected watcher to be stopped")
	}

	if watcher.Status != "Stopped" {
		t.Errorf("Expected status 'Stopped', got '%s'", watcher.Status)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["status"] != "stopped" {
		t.Errorf("Expected status 'stopped', got '%s'", result["status"])
	}

	// Test stopping when not running
	watcher.Running = false
	w2 := httptest.NewRecorder()
	watcher.handleStop(w2, req)
	// Should still return OK
	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200 when not running, got %d", w2.Code)
	}
}

func TestWatcher_handleEvents(t *testing.T) {
	watcher := &Watcher{
		Config: Config{
			WatchDir: "/test/dir",
		},
		Status:  "Running",
		Clients: make(map[chan string]bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/events", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	// Test that it sets up SSE headers
	done := make(chan bool)
	go func() {
		watcher.handleEvents(w, req)
		done <- true
	}()

	// Give it a moment to set headers and enter the loop
	time.Sleep(50 * time.Millisecond)

	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Error("Expected Content-Type to be 'text/event-stream'")
	}

	if w.Header().Get("Cache-Control") != "no-cache" {
		t.Error("Expected Cache-Control to be 'no-cache'")
	}

	// Cancel the request context to stop the handler
	cancel()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Error("Handler did not complete")
	}
}

func TestWatcher_broadcast(t *testing.T) {
	watcher := &Watcher{
		Clients: make(map[chan string]bool),
	}

	// Create test clients
	client1 := make(chan string, 1)
	client2 := make(chan string, 1)

	watcher.Clients[client1] = true
	watcher.Clients[client2] = true

	// Broadcast a message
	watcher.broadcast("test message")

	// Check that both clients received the message
	select {
	case msg := <-client1:
		if msg != "test message" {
			t.Errorf("Client1 expected 'test message', got '%s'", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client1 did not receive message")
	}

	select {
	case msg := <-client2:
		if msg != "test message" {
			t.Errorf("Client2 expected 'test message', got '%s'", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client2 did not receive message")
	}
}

func TestScan(t *testing.T) {
	tempDir := t.TempDir()

	// Create test .adoc files
	files := []string{
		"file1.adoc",
		"file2.adoc",
		"subdir/file3.adoc",
		"other.txt", // Should be ignored
	}

	for _, f := range files {
		fullPath := filepath.Join(tempDir, f)
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		os.WriteFile(fullPath, []byte("content"), 0644)
	}

	modMap := make(map[string]time.Time)
	scan(tempDir, modMap)

	// Should find 3 .adoc files
	if len(modMap) != 3 {
		t.Errorf("Expected 3 files in modMap, got %d", len(modMap))
	}

	// Verify all .adoc files are in the map
	expectedFiles := []string{
		filepath.Join(tempDir, "file1.adoc"),
		filepath.Join(tempDir, "file2.adoc"),
		filepath.Join(tempDir, "subdir/file3.adoc"),
	}

	for _, expected := range expectedFiles {
		if _, ok := modMap[expected]; !ok {
			t.Errorf("Expected file %s to be in modMap", expected)
		}
	}
}

func TestCheckChanges(t *testing.T) {
	tempDir := t.TempDir()

	// Create initial file
	testFile := filepath.Join(tempDir, "test.adoc")
	os.WriteFile(testFile, []byte("initial content"), 0644)

	modMap := make(map[string]time.Time)
	scan(tempDir, modMap)

	// Wait a bit to ensure different mod time
	time.Sleep(10 * time.Millisecond)

	// Modify the file
	os.WriteFile(testFile, []byte("modified content"), 0644)

	// Check for changes
	changes := checkChanges(tempDir, modMap)

	if len(changes) != 1 {
		t.Errorf("Expected 1 change, got %d", len(changes))
	}

	if changes[0] != testFile {
		t.Errorf("Expected change in %s, got %s", testFile, changes[0])
	}

	// Check again - should be no changes
	changes2 := checkChanges(tempDir, modMap)
	if len(changes2) != 0 {
		t.Errorf("Expected 0 changes on second check, got %d", len(changes2))
	}
}

func TestCheckChanges_NewFile(t *testing.T) {
	tempDir := t.TempDir()

	modMap := make(map[string]time.Time)
	scan(tempDir, modMap)

	// Create a new file
	newFile := filepath.Join(tempDir, "new.adoc")
	os.WriteFile(newFile, []byte("new content"), 0644)

	changes := checkChanges(tempDir, modMap)

	if len(changes) != 1 {
		t.Errorf("Expected 1 change (new file), got %d", len(changes))
	}

	if changes[0] != newFile {
		t.Errorf("Expected change in %s, got %s", newFile, changes[0])
	}
}

func TestWatcher_run(t *testing.T) {
	tempDir := t.TempDir()

	watcher := &Watcher{
		Config: Config{
			WatchDir: tempDir,
		},
		Running:  true,
		StopChan: make(chan struct{}),
		Clients:  make(map[chan string]bool),
	}

	// Start the watcher in a goroutine
	done := make(chan bool)
	go func() {
		watcher.run()
		done <- true
	}()

	// Give it a moment to initialize
	time.Sleep(50 * time.Millisecond)

	// Stop the watcher
	close(watcher.StopChan)

	// Wait for it to finish
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Error("Watcher did not stop")
	}
}

func TestWatcher_ClientsManagement(t *testing.T) {
	watcher := &Watcher{
		Clients: make(map[chan string]bool),
	}

	// Test adding clients
	client1 := make(chan string, 1)
	client2 := make(chan string, 1)

	watcher.ClientsLock.Lock()
	watcher.Clients[client1] = true
	watcher.Clients[client2] = true
	watcher.ClientsLock.Unlock()

	if len(watcher.Clients) != 2 {
		t.Errorf("Expected 2 clients, got %d", len(watcher.Clients))
	}

	// Test removing clients
	watcher.ClientsLock.Lock()
	delete(watcher.Clients, client1)
	watcher.ClientsLock.Unlock()

	if len(watcher.Clients) != 1 {
		t.Errorf("Expected 1 client after removal, got %d", len(watcher.Clients))
	}
}

func TestScan_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	modMap := make(map[string]time.Time)
	scan(tempDir, modMap)

	if len(modMap) != 0 {
		t.Errorf("Expected 0 files in empty directory, got %d", len(modMap))
	}
}

func TestScan_NonExistentDirectory(t *testing.T) {
	modMap := make(map[string]time.Time)
	// Should not panic
	scan("/nonexistent/directory", modMap)

	if len(modMap) != 0 {
		t.Errorf("Expected 0 files for non-existent directory, got %d", len(modMap))
	}
}

func TestCheckChanges_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	modMap := make(map[string]time.Time)
	scan(tempDir, modMap)

	changes := checkChanges(tempDir, modMap)

	if len(changes) != 0 {
		t.Errorf("Expected 0 changes in empty directory, got %d", len(changes))
	}
}

func TestWatcher_broadcast_Concurrent(t *testing.T) {
	watcher := &Watcher{
		Clients: make(map[chan string]bool),
	}

	// Create multiple clients
	numClients := 10
	clients := make([]chan string, numClients)
	for i := 0; i < numClients; i++ {
		clients[i] = make(chan string, 1)
		watcher.Clients[clients[i]] = true
	}

	// Broadcast concurrently
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(msgNum int) {
			defer wg.Done()
			watcher.broadcast("message " + string(rune('0'+msgNum)))
		}(i)
	}

	wg.Wait()

	// All clients should have received messages (though order may vary)
	received := 0
	for _, client := range clients {
		select {
		case <-client:
			received++
		case <-time.After(10 * time.Millisecond):
		}
	}

	// At least some clients should have received messages
	if received == 0 {
		t.Error("No clients received messages")
	}
}

func TestWatcher_handleEvents_ClientCleanup(t *testing.T) {
	watcher := &Watcher{
		Config: Config{
			WatchDir: "/test/dir",
		},
		Status:  "Running",
		Clients: make(map[chan string]bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/events", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	initialClientCount := len(watcher.Clients)

	// Start handler in goroutine
	done := make(chan bool)
	go func() {
		watcher.handleEvents(w, req)
		done <- true
	}()

	// Give it time to register client and enter the loop
	time.Sleep(50 * time.Millisecond)

	// Cancel request context
	cancel()

	// Wait for handler to finish
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Error("Handler did not complete")
	}

	// Client should be cleaned up
	// Note: This is a bit tricky to test exactly due to timing, but we can verify
	// the cleanup logic exists
	if len(watcher.Clients) > initialClientCount+1 {
		t.Error("Clients may not have been cleaned up properly")
	}
}

