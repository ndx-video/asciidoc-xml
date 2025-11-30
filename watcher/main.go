package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ndx-video/asciidoc-xml/lib"
)

type Config struct {
	WatchDir string `json:"watchDir"`
	Port     int    `json:"watchPort"` // Default 8006
}

type Watcher struct {
	Config      Config
	Clients     map[chan string]bool
	ClientsLock sync.Mutex
	Running     bool
	StopChan    chan struct{}
	Status      string
	logger      *lib.Logger
}

var (
	watchDir string
	cliPath  string
	port     int
	watcher  *Watcher
)

func main() {
	showVersion := flag.Bool("version", false, "Show version information and exit")
	flag.StringVar(&watchDir, "watch", ".", "Directory to watch")
	flag.StringVar(&cliPath, "cli", "adc", "Path to adc CLI tool")
	flag.IntVar(&port, "port", 8006, "Port for watcher daemon")
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("adc-watcher version %s\n", version)
		os.Exit(0)
	}

	// Load config
	cfg := Config{
		WatchDir: watchDir,
		Port:     port,
	}

	if data, err := os.ReadFile("adc.json"); err == nil {
		var loadedCfg Config
		if err := json.Unmarshal(data, &loadedCfg); err == nil {
			if loadedCfg.WatchDir != "" {
				// If flag is default, overwrite with config
				isFlagSet := false
				flag.Visit(func(f *flag.Flag) {
					if f.Name == "watch" {
						isFlagSet = true
					}
				})
				if !isFlagSet {
					cfg.WatchDir = loadedCfg.WatchDir
				}
			}
			if loadedCfg.Port != 0 {
				// If flag is default, overwrite with config
				isFlagSet := false
				flag.Visit(func(f *flag.Flag) {
					if f.Name == "port" {
						isFlagSet = true
					}
				})
				if !isFlagSet {
					cfg.Port = loadedCfg.Port
				}
			}
		}
	}

	// Initialize logger with default config
	defaultLogConfig := lib.LogConfig{
		Level:  "info",
		Format: "text",
		Console: lib.ConsoleConfig{
			Enabled: true,
			Level:   "info",
		},
		File: lib.FileConfig{
			Enabled:  false,
			Path:     "./logs",
			Filename: "watcher.log",
			MaxSize:  10 * 1024 * 1024,
			MaxFiles: 5,
			Rotation: "size",
		},
	}
	
	logger, err := lib.NewLogger(defaultLogConfig)
	if err != nil {
		log.Printf("Failed to initialize logger: %v, using default logging", err)
	}

	watcher = &Watcher{
		Config:   cfg,
		Clients:  make(map[chan string]bool),
		Running:  false,
		Status:   "Stopped",
		logger:   logger,
	}

	// Start HTTP Server
	mux := http.NewServeMux()
	mux.HandleFunc("/events", watcher.handleEvents)
	mux.HandleFunc("/start", watcher.handleStart)
	mux.HandleFunc("/stop", watcher.handleStop)
	mux.HandleFunc("/status", watcher.handleStatus)
	
	// Enable CORS
	handler := corsMiddleware(mux)

	addr := fmt.Sprintf(":%d", cfg.Port)
	if watcher.logger != nil {
		watcher.logger.Info(nil, "Watcher daemon starting",
			"port", cfg.Port,
			"address", addr,
			"watch_dir", cfg.WatchDir,
		)
	} else {
		fmt.Printf("Watcher daemon listening on %s\n", addr)
	}
	if err := http.ListenAndServe(addr, handler); err != nil {
		if watcher.logger != nil {
			watcher.logger.Fatal(nil, "Watcher daemon failed to start", "error", err.Error())
		}
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func (w *Watcher) handleEvents(httpw http.ResponseWriter, r *http.Request) {
	// CORS headers must be set first
	httpw.Header().Set("Access-Control-Allow-Origin", "*")
	httpw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	httpw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// SSE headers
	httpw.Header().Set("Content-Type", "text/event-stream")
	httpw.Header().Set("Cache-Control", "no-cache")
	httpw.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string, 1)
	
	w.ClientsLock.Lock()
	w.Clients[messageChan] = true
	w.ClientsLock.Unlock()

	// Send initial status
	messageChan <- fmt.Sprintf("Watcher status: %s", w.Status)

	defer func() {
		w.ClientsLock.Lock()
		delete(w.Clients, messageChan)
		w.ClientsLock.Unlock()
		close(messageChan)
	}()

	flusher, ok := httpw.(http.Flusher)
	if !ok {
		http.Error(httpw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(httpw, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (w *Watcher) broadcast(msg string) {
	w.ClientsLock.Lock()
	defer w.ClientsLock.Unlock()
	
	if w.logger != nil {
		w.logger.Info(nil, "Broadcasting message", "message", msg)
	} else {
		fmt.Println(msg) // Also log to stdout
	}
	
	for client := range w.Clients {
		select {
		case client <- msg:
		default:
			// Skip if blocked
		}
	}
}

func (w *Watcher) handleStart(httpw http.ResponseWriter, r *http.Request) {
	// CORS headers
	httpw.Header().Set("Access-Control-Allow-Origin", "*")
	httpw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	httpw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		httpw.WriteHeader(http.StatusOK)
		return
	}
	
	if w.Running {
		w.broadcast("Watcher already running")
		httpw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(httpw).Encode(map[string]string{"status": "already running"})
		return
	}
	
	w.Running = true
	w.StopChan = make(chan struct{})
	w.Status = "Running"
	w.broadcast(fmt.Sprintf("Watcher started on %s", w.Config.WatchDir))
	
	go w.run()
	
	httpw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(httpw).Encode(map[string]string{"status": "started"})
}

func (w *Watcher) handleStop(httpw http.ResponseWriter, r *http.Request) {
	// CORS headers
	httpw.Header().Set("Access-Control-Allow-Origin", "*")
	httpw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	httpw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		httpw.WriteHeader(http.StatusOK)
		return
	}
	
	if !w.Running {
		if w.logger != nil {
			w.logger.Warn(nil, "Attempted to stop watcher that is not running")
		}
		w.broadcast("Watcher not running")
		httpw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(httpw).Encode(map[string]string{"status": "not running"})
		return
	}
	
	if w.logger != nil {
		w.logger.Info(nil, "Stopping watcher")
	}
	
	close(w.StopChan)
	w.Running = false
	w.Status = "Stopped"
	w.broadcast("Watcher stopped")
	
	httpw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(httpw).Encode(map[string]string{"status": "stopped"})
}

func (w *Watcher) handleStatus(httpw http.ResponseWriter, r *http.Request) {
	// CORS headers
	httpw.Header().Set("Access-Control-Allow-Origin", "*")
	httpw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	httpw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		httpw.WriteHeader(http.StatusOK)
		return
	}
	
	httpw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(httpw).Encode(map[string]interface{}{
		"status": w.Status,
		"config": w.Config,
	})
}

func (w *Watcher) run() {
	absDir, err := filepath.Abs(w.Config.WatchDir)
	if err != nil {
		errMsg := fmt.Sprintf("Error getting absolute path: %v", err)
		if w.logger != nil {
			w.logger.Error(nil, "Failed to get absolute path",
				"watch_dir", w.Config.WatchDir,
				"error", err.Error(),
			)
		}
		w.broadcast(errMsg)
		return
	}

	if w.logger != nil {
		w.logger.Info(nil, "Watcher started",
			"watch_dir", absDir,
		)
	}

	lastMod := make(map[string]time.Time)
	scan(absDir, lastMod)

	if w.logger != nil {
		w.logger.Info(nil, "Initial scan completed",
			"directory", absDir,
			"files_tracked", len(lastMod),
		)
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.StopChan:
			if w.logger != nil {
				w.logger.Info(nil, "Watcher stopped")
			}
			return
		case <-ticker.C:
			changes := checkChanges(absDir, lastMod)
			if len(changes) > 0 {
				if w.logger != nil {
					w.logger.Info(nil, "File changes detected",
						"change_count", len(changes),
						"files", changes,
					)
				}
				w.broadcast(fmt.Sprintf("Detected %d file change(s)", len(changes)))
				
				// Send each file path via SSE for the web app to process
				for _, filePath := range changes {
					// Send absolute file path
					absPath, err := filepath.Abs(filePath)
					if err != nil {
						if w.logger != nil {
							w.logger.Warn(nil, "Failed to get absolute path for changed file",
								"file", filePath,
								"error", err.Error(),
							)
						}
						w.broadcast(fmt.Sprintf("Error getting absolute path for %s: %v", filePath, err))
						continue
					}
					// Broadcast the file path - web app will add it to queue
					w.broadcast(absPath)
				}
			}
		}
	}
}

func scan(dir string, modMap map[string]time.Time) {
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".adoc") {
			info, err := d.Info()
			if err == nil {
				modMap[path] = info.ModTime()
			}
		}
		return nil
	})
}

func checkChanges(dir string, modMap map[string]time.Time) []string {
	var changed []string

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".adoc") {
			info, err := d.Info()
			if err == nil {
				last, ok := modMap[path]
				if !ok || info.ModTime().After(last) {
					changed = append(changed, path)
					modMap[path] = info.ModTime()
				}
			}
		}
		return nil
	})

	return changed
}
