package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool argument types
type Base64EncodeArgs struct {
	Text string `json:"text"`
}

type SHA256HashArgs struct {
	Text string `json:"text"`
}

type URLEncodeArgs struct {
	Text string `json:"text"`
}

type JSONValidateArgs struct {
	JSONString string `json:"json_string"`
}

type NoArgs struct{}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("========================================")
	log.Printf("MCP Streamable HTTP Server Starting")
	log.Printf("========================================")

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sample-mcp-server-streamable-http",
		Version: "1.0.0",
	}, nil)

	log.Printf("[MAIN] Registering tools...")

	// Create schemas for tool inputs
	base64EncodeSchema, _ := jsonschema.For[Base64EncodeArgs](nil)
	sha256HashSchema, _ := jsonschema.For[SHA256HashArgs](nil)
	getTimestampSchema, _ := jsonschema.For[NoArgs](nil)
	jsonValidateSchema, _ := jsonschema.For[JSONValidateArgs](nil)
	urlEncodeSchema, _ := jsonschema.For[URLEncodeArgs](nil)

	// Register base64_encode tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "base64_encode",
		Description: "Encodes text to base64 format",
		InputSchema: base64EncodeSchema,
	}, base64EncodeHandler)
	log.Printf("[MAIN] - Registered tool: base64_encode")

	// Register sha256_hash tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "sha256_hash",
		Description: "Generates SHA256 hash of the given text",
		InputSchema: sha256HashSchema,
	}, sha256HashHandler)
	log.Printf("[MAIN] - Registered tool: sha256_hash")

	// Register get_timestamp tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_timestamp",
		Description: "Returns the current Unix timestamp in seconds and milliseconds",
		InputSchema: getTimestampSchema,
	}, getTimestampHandler)
	log.Printf("[MAIN] - Registered tool: get_timestamp")

	// Register json_validate tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "json_validate",
		Description: "Validates if a string is valid JSON and returns parsed structure info",
		InputSchema: jsonValidateSchema,
	}, jsonValidateHandler)
	log.Printf("[MAIN] - Registered tool: json_validate")

	// Register url_encode tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "url_encode",
		Description: "URL encodes the given text",
		InputSchema: urlEncodeSchema,
	}, urlEncodeHandler)
	log.Printf("[MAIN] - Registered tool: url_encode")

	log.Printf("[MAIN] All 5 tools registered successfully")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Create streamable HTTP handler
	log.Printf("[MAIN] Setting up Streamable HTTP transport on :%s", port)

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		log.Printf("[HTTP] New connection from %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		log.Printf("[HTTP] Headers: User-Agent=%s, Origin=%s", r.Header.Get("User-Agent"), r.Header.Get("Origin"))
		return server
	}, nil)

	// Create HTTP server with streamable HTTP and health endpoints
	mux := http.NewServeMux()

	// Streamable HTTP endpoint
	mux.Handle("/mcp", handler)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HEALTH] Health check from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Catch-all for debugging
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/mcp" && r.URL.Path != "/health" {
			log.Printf("[WARN] Unknown endpoint requested: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.NotFound(w, r)
		}
	})

	log.Printf("========================================")
	log.Printf("[MAIN] Server ready and listening on :%s", port)
	log.Printf("[MAIN] Streamable HTTP endpoint: http://localhost:%s/mcp", port)
	log.Printf("[MAIN] Health endpoint: http://localhost:%s/health", port)
	log.Printf("========================================")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("[MAIN] FATAL: Server failed to start: %v", err)
	}
}

// Tool Handlers

func base64EncodeHandler(ctx context.Context, req *mcp.CallToolRequest, args Base64EncodeArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing base64_encode with args: %+v", args)

	encoded := base64.StdEncoding.EncodeToString([]byte(args.Text))
	responseText := fmt.Sprintf("Original: %s\nBase64 Encoded: %s", args.Text, encoded)

	log.Printf("[TOOL] Base64 encoding successful")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

func sha256HashHandler(ctx context.Context, req *mcp.CallToolRequest, args SHA256HashArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing sha256_hash with args: %+v", args)

	hash := sha256.Sum256([]byte(args.Text))
	hashHex := fmt.Sprintf("%x", hash)
	responseText := fmt.Sprintf("Original: %s\nSHA256 Hash: %s", args.Text, hashHex)

	log.Printf("[TOOL] SHA256 hash generated successfully")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

func getTimestampHandler(ctx context.Context, req *mcp.CallToolRequest, args NoArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing get_timestamp")

	now := time.Now()
	unixSec := now.Unix()
	unixMilli := now.UnixMilli()
	formattedTime := now.Format(time.RFC3339)

	responseText := fmt.Sprintf("Current Timestamp:\n- Unix (seconds): %d\n- Unix (milliseconds): %d\n- RFC3339: %s", unixSec, unixMilli, formattedTime)

	log.Printf("[TOOL] Timestamp generated: %d", unixSec)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

func jsonValidateHandler(ctx context.Context, req *mcp.CallToolRequest, args JSONValidateArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing json_validate with args: %+v", args)

	var parsed interface{}
	err := json.Unmarshal([]byte(args.JSONString), &parsed)

	if err != nil {
		responseText := fmt.Sprintf("Invalid JSON: %s\nError: %v", args.JSONString, err)
		log.Printf("[TOOL] JSON validation failed: %v", err)

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText},
			},
			IsError: true,
		}, nil, nil
	}

	// Determine JSON type
	var jsonType string
	switch parsed.(type) {
	case map[string]interface{}:
		jsonType = "Object"
	case []interface{}:
		jsonType = "Array"
	case string:
		jsonType = "String"
	case float64:
		jsonType = "Number"
	case bool:
		jsonType = "Boolean"
	case nil:
		jsonType = "Null"
	default:
		jsonType = "Unknown"
	}

	// Pretty print the JSON
	prettyJSON, _ := json.MarshalIndent(parsed, "", "  ")
	responseText := fmt.Sprintf("Valid JSON (%s)\n\nPretty printed:\n%s", jsonType, string(prettyJSON))

	log.Printf("[TOOL] JSON validation successful: type=%s", jsonType)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

func urlEncodeHandler(ctx context.Context, req *mcp.CallToolRequest, args URLEncodeArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing url_encode with args: %+v", args)

	encoded := url.QueryEscape(args.Text)
	responseText := fmt.Sprintf("Original: %s\nURL Encoded: %s", args.Text, encoded)

	log.Printf("[TOOL] URL encoding successful")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}
