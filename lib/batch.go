package lib

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// BatchConfig holds configuration for batch processing
type BatchConfig struct {
	MaxWorkers        int
	ParallelThreshold int
	EnableParallel    bool
	DryRun            bool
	ValidateOnly      bool
}

// ProcessingLimits defines safety limits for file processing
type ProcessingLimits struct {
	MaxFileSize    int64 // bytes
	MaxArchiveSize int64 // bytes
	MaxFileCount   int
}

// Default limits
const (
	DefaultMaxFileSize    = 10 * 1024 * 1024 // 10MB
	DefaultMaxArchiveSize = 10 * 1024 * 1024 // 10MB
	DefaultMaxFileCount   = 10000
	DefaultThreshold      = 2
)

// BatchResult holds the results of a batch operation
type BatchResult struct {
	TotalFiles   int
	SuccessCount int
	ErrorCount   int
	Errors       []FileError
	Status       string
	Duration     time.Duration
}

// FileError represents an error occurring during file processing
type FileError struct {
	File  string
	Error error
}

// ProgressCallback is a function type for reporting progress
type ProgressCallback func(current int, total int, file string, err error)

// ProcessFilesParallel processes a list of files, optionally in parallel
// logger is optional - if nil, no logging is performed
func ProcessFilesParallel(files []string, processFn func(string) error, config BatchConfig, limits ProcessingLimits, progressCb ProgressCallback, logger *Logger) BatchResult {
	ctx := context.Background()
	startTime := time.Now()
	result := BatchResult{
		TotalFiles: len(files),
	}

	if logger != nil {
		logger.Info(ctx, "Starting batch processing",
			"total_files", len(files),
			"max_workers", config.MaxWorkers,
			"parallel_threshold", config.ParallelThreshold,
			"enable_parallel", config.EnableParallel,
			"dry_run", config.DryRun,
			"validate_only", config.ValidateOnly,
		)
	}

	// Validate file count limit
	if len(files) > limits.MaxFileCount {
		err := fmt.Errorf("file count %d exceeds limit %d", len(files), limits.MaxFileCount)
		result.Status = "Failed"
		result.ErrorCount = 1
		result.Errors = append(result.Errors, FileError{
			File:  "batch",
			Error: err,
		})
		if logger != nil {
			logger.Error(ctx, "Batch processing failed: file count limit exceeded",
				"file_count", len(files),
				"max_file_count", limits.MaxFileCount,
			)
		}
		return result
	}

	// Dry Run / Validation Only setup
	if config.DryRun {
		if logger != nil {
			logger.Info(ctx, "Running in dry-run mode")
		}
		result.Status = "Dry run successful"
		for i, file := range files {
			if err := checkLimits(file, limits); err != nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, FileError{File: file, Error: err})
				if logger != nil {
					logger.Warn(ctx, "File validation failed in dry-run",
						"file", file,
						"error", err.Error(),
					)
				}
				if progressCb != nil {
					progressCb(i+1, len(files), file, err)
				}
			} else {
				result.SuccessCount++
				if logger != nil {
					logger.Debug(ctx, "File validated in dry-run", "file", file)
				}
				if progressCb != nil {
					progressCb(i+1, len(files), file, nil)
				}
			}
		}
		result.Duration = time.Since(startTime)
		if logger != nil {
			logger.Info(ctx, "Dry-run completed",
				"success_count", result.SuccessCount,
				"error_count", result.ErrorCount,
				"duration_ms", result.Duration.Milliseconds(),
			)
		}
		return result
	}

	// Determine effective workers
	workers := config.MaxWorkers
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	if !config.EnableParallel && len(files) < config.ParallelThreshold {
		workers = 1
	}
	if logger != nil {
		logger.Debug(ctx, "Worker configuration",
			"effective_workers", workers,
			"file_count", len(files),
			"parallel_enabled", config.EnableParallel && len(files) >= config.ParallelThreshold,
		)
	}

	// Worker pool setup
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)
	resultsChan := make(chan FileError, len(files))

	for i, file := range files {
		wg.Add(1)
		
		// If sequential (1 worker), just run directly to maintain order if desired, 
		// but here we use the same channel structure for simplicity.
		// For strict progress reporting order, channels might need adjustment, 
		// but basic callbacks work fine from routines.
		
		go func(f string, index int) {
			defer wg.Done()
			sem <- struct{}{} // Acquire token
			defer func() { <-sem }() // Release token

			var err error
			
			// Check file limits first
			if limitErr := checkLimits(f, limits); limitErr != nil {
				err = limitErr
			} else {
			// Process file
			err = processFn(f)
		}

		// Report result
		if err != nil {
			resultsChan <- FileError{File: f, Error: err}
			if logger != nil {
				logger.Error(ctx, "File processing failed",
					"file", f,
					"error", err.Error(),
				)
			}
		} else {
			resultsChan <- FileError{File: f, Error: nil}
			if logger != nil {
				logger.Debug(ctx, "File processed successfully", "file", f)
			}
		}
		
		// Optional: Report progress (concurrency safe callback required)
		if progressCb != nil {
			// Note: This simple callback from goroutine means order isn't guaranteed in UI
			progressCb(0, len(files), f, err) 
		}
		}(file, i)
	}

	wg.Wait()
	close(resultsChan)

	// Collect results
	for res := range resultsChan {
		if res.Error != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, res)
		} else {
			result.SuccessCount++
		}
	}

	result.Status = "Completed"
	if result.ErrorCount > 0 {
		result.Status = "Completed with errors"
	}
	result.Duration = time.Since(startTime)

	if logger != nil {
		logger.Info(ctx, "Batch processing completed",
			"total_files", result.TotalFiles,
			"success_count", result.SuccessCount,
			"error_count", result.ErrorCount,
			"duration_ms", result.Duration.Milliseconds(),
			"status", result.Status,
		)
	}

	return result
}

func checkLimits(file string, limits ProcessingLimits) error {
	info, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	if info.Size() > limits.MaxFileSize {
		return fmt.Errorf("file size %d exceeds limit %d", info.Size(), limits.MaxFileSize)
	}
	return nil
}

