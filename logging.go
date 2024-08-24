package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	logger            = log.New(os.Stderr, boldGreen("[DeepSpace] "), log.LstdFlags)
	serverErrorLogger = log.New(getPalaceServerErrorLog(), "", log.LstdFlags)
)

var (
	boldWhite   = color.New(color.FgHiWhite, color.Bold).SprintFunc()
	boldGreen   = color.New(color.FgGreen, color.Bold).SprintFunc()
	boldGreenf  = color.New(color.FgGreen, color.Bold).SprintfFunc()
	boldYellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	boldYellowf = color.New(color.FgYellow, color.Bold).SprintfFunc()
	boldRed     = color.New(color.FgRed, color.Bold).SprintFunc()
	green       = color.New(color.FgHiGreen).SprintFunc()
	red         = color.New(color.FgRed).SprintFunc()
)

const asciiDeepSpace = `
____    _____   _____   ____    ____    ____       _       ____   _____ 
|  _ \  | ____| | ____| |  _ \  / ___|  |  _ \     / \     / ___| | ____|
| | | | |  _|   |  _|   | |_) | \___ \  | |_) |   / _ \   | |     |  _|  
| |_| | | |___  | |___  |  __/   ___) | |  __/   / ___ \  | |___  | |___ 
|____/  |_____| |_____| |_|     |____/  |_|     /_/   \_\  \____| |_____|
`

func logServerStarts(baseUrl string) {
	logger.Println(boldWhite("DeepSpace Starts => change base_url to "+strconv.Quote(baseUrl)) + "\n" + asciiDeepSpace)
}

func logRequest(
	method string,
	path string,
	query string,
	requestContentType string,
	requestID string,
	responseStatus string,
	responseContentType string,
	responseTTFT int,
	deepseekRequestID string,
	deepseekServerTiming int,
	deepseekContextCacheID string,
	deepseekUID string,
	deepseekGID string,
	deepseek *DeepSeek,
	latency time.Duration,
	tokenFinishLatency time.Duration,
	err error,
	warnings []error,
) {
	if query != "" {
		path += "?" + query
	}
	if strings.HasPrefix(responseStatus, "2") {
		responseStatus = green(responseStatus)
	} else {
		responseStatus = red(responseStatus)
	}
	logger.Printf("%s %s %s %.2fs\n",
		boldYellowf("%-6s", method),
		boldWhite(path),
		responseStatus,
		float64(latency)/float64(time.Second),
	)
	if requestContentType != "" {
		logger.Printf("  - Request Headers: \n")
		logger.Printf("    - Content-Type:   %s\n", requestContentType)
		if requestID != "" {
			logger.Printf("    - X-Request-Id:   %s\n", requestID)
		}
	}
	if deepseekRequestID != "" {
		logger.Printf("  - Response Headers: \n")
		logger.Printf("    - Content-Type:          %s\n", responseContentType)
		logger.Printf("    - Dse-Request-Id:        %s\n", deepseekRequestID)
		logger.Printf("    - Server-Timing:         %ss\n", boldYellowf("%.4f", float64(deepseekServerTiming)/1000.00))
		if deepseekContextCacheID != "" {
			logger.Printf("    - Dse-Context-Cache-Id:  %s\n", deepseekContextCacheID)
		}
		if deepseekUID != "" {
			logger.Printf("    - Dse-Uid:               %s\n", deepseekUID)
			logger.Printf("    - Dse-Gid:               %s\n", deepseekGID)
		}
	}
	if deepseek != nil && deepseek.ID != "" {
		logger.Printf("  - Response: \n")
		logger.Printf("    - id:                %s\n", deepseek.ID)
		if responseTTFT > 0 {
			logger.Printf("    - ttft:              %ss\n", boldYellowf("%.4f", float64(responseTTFT)/1000.00))
		}
		if usage := deepseek.Usage; usage != nil {
			if tokenFinishLatency > 0 {
				timePerOutputToken := ((float64(tokenFinishLatency) -
					float64(responseTTFT)*float64(time.Millisecond)) /
					float64(time.Second)) /
					float64(usage.CompletionTokens-_boolToInt(responseTTFT != 0))
				logger.Printf("    - tpot:              %ss/token\n", boldYellowf("%.4f", timePerOutputToken))
				logger.Printf("    - otps:              %stokens/s\n", boldYellowf("%.4f", 1/timePerOutputToken))
			}
			logger.Printf("    - prompt_tokens:     %d\n", usage.PromptTokens)
			logger.Printf("    - completion_tokens: %d\n", usage.CompletionTokens)
			logger.Printf("    - total_tokens:      %d\n", usage.TotalTokens)
			if usage.CachedTokens > 0 {
				logger.Printf("    - cached_tokens:     %d\n", usage.CachedTokens)
			}
		}
	}
	if err != nil {
		if errorMsg := err.Error(); errorMsg != "" {
			var (
				indent = "  "
				render = boldRed
			)
			if errors.As(err, new(*deepseekError)) {
				logger.Printf("  - %s: \n", boldYellow("DeepSeek Error"))
				indent += "  "
				render = boldYellow
			}
			for _, line := range strings.Split(errorMsg, "\n") {
				logger.Printf("%s%s\n", indent, render(line))
			}
		}
	}
	if len(warnings) > 0 {
		for _, warning := range warnings {
			logger.Printf("  %s %s\n", boldYellow("[WARNING]"), boldYellow(warning.Error()))
		}
	}
}

func _boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func logNewRow(id int64) {
	logger.Println(
		boldWhite("  New Row Inserted:"),
		boldGreenf("last_insert_id=%d", id),
	)
}

func logExport(file *os.File) {
	logger.Println("export to", boldGreen(file.Name()), "successfully")
}

func logFatal(err error) {
	if errorMsg := err.Error(); errorMsg != "" {
		for _, line := range strings.Split(errorMsg, "\n") {
			fmt.Fprintln(os.Stderr, boldRed(line))
		}
	}
	os.Exit(2)
}
