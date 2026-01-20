# Sample MCP Server with Streamable HTTP (Golang)

A Model Context Protocol (MCP) server implementation in Go using the official [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) with Streamable HTTP transport.

## Features

This MCP server provides 5 utility tools:

1. **base64_encode** - Encodes text to base64 format
2. **sha256_hash** - Generates SHA256 hash of the given text
3. **get_timestamp** - Returns the current Unix timestamp (seconds and milliseconds)
4. **json_validate** - Validates if a string is valid JSON and returns structure info
5. **url_encode** - URL encodes the given text

## Architecture

Built with the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk`), this server:
- Uses the SDK's built-in Streamable HTTP transport implementation
- Automatically handles JSON-RPC 2.0 protocol
- Provides comprehensive logging for debugging
- Includes CORS support for browser-based clients

### Endpoints

- `/mcp` - Streamable HTTP endpoint for JSON-RPC requests
- `/health` - Health check endpoint

## Requirements

- Go 1.23 or later
- Docker (for containerization)
- Kubernetes cluster (for deployment)

## Installation

```bash
cd sample-mcp-server-streamable-http
go mod download
go build -o mcp-server
```

## Running the Server

```bash
./mcp-server
```

The server will start on `http://localhost:8081`

## Kubernetes Deployment

### Build and Push Docker Image

```bash
cd sample-mcp-server-streamable-http

# Build for linux/amd64 architecture
docker build -t aliok/mcp-server-streamable-http:latest .

# Push to Docker Hub (or your registry)
docker push aliok/mcp-server-streamable-http:latest
```

**Note**: Replace `aliok` with your Docker Hub username or registry path. The Dockerfile is configured to build for linux/amd64 architecture.

### Deploy to Kubernetes

Deploy using kubectl:

```bash
# Apply manifests
kubectl apply -f k8s/
```

**Note**: The deployment uses `aliok/mcp-server-streamable-http:latest`. Update `k8s/deployment.yaml` if using a different registry.

### Verify Deployment

```bash
# Check deployment status
kubectl get deployments mcp-server-streamable-http

# Check pods
kubectl get pods -l app=mcp-server-streamable-http

# Check service
kubectl get svc mcp-server-streamable-http
```

### Access the Service

#### Port Forward (for testing)

```bash
kubectl port-forward svc/mcp-server-streamable-http 8081:80
```

Then access at `http://localhost:8081`

#### Ingress (for production)

Create an Ingress resource to expose the service:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: mcp-server-streamable-http
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: mcp-server-http.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: mcp-server-streamable-http
            port:
              number: 80
```

### Kubernetes Configuration

The deployment includes:
- **Replicas**: 2 pods for high availability
- **Health Checks**: Liveness and readiness probes on `/health`
- **Resources**: CPU and memory limits configured
- **Service**: ClusterIP service on port 80

To scale the deployment:

```bash
kubectl scale deployment mcp-server-streamable-http --replicas=3
```

### Clean Up

```bash
kubectl delete -k k8s/
```

## Testing with MCP Inspector

The [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is a visual interface for testing and debugging MCP servers.

### Prerequisites

Install the MCP Inspector:

```bash
npx @modelcontextprotocol/inspector
```

### Running the Inspector

1. Start your MCP server:

```bash
# Local development
./mcp-server

# Or with Kubernetes port-forward
kubectl port-forward svc/mcp-server-streamable-http 8081:80
```

2. In a new terminal, launch the inspector with your server's HTTP endpoint:

```bash
npx @modelcontextprotocol/inspector http://localhost:8081/mcp
```

3. The inspector will open in your browser, showing:
   - **Server Info**: Connection status and server capabilities
   - **Tools Tab**: List of all 5 available tools
   - **Interactive Testing**: Forms to call each tool with parameters

### Using the Inspector

**Test the Base64 Encode Tool:**
1. Click on the "Tools" tab
2. Select "base64_encode" from the list
3. Fill in the parameter:
   - text: "Hello, World!"
4. Click "Call Tool" to see the base64 encoded result

**Test the SHA256 Hash Tool:**
1. Select "sha256_hash"
2. Enter text: "password123"
3. Click "Call Tool" to see the hash

**Test the Get Timestamp Tool:**
1. Select "get_timestamp"
2. Click "Call Tool" (no parameters needed)
3. View the current Unix timestamp

**Test the JSON Validate Tool:**
1. Select "json_validate"
2. Enter json_string: `{"name": "John", "age": 30}`
3. Click "Call Tool" to see validation result and pretty-printed JSON

**Test the URL Encode Tool:**
1. Select "url_encode"
2. Enter text: "hello world & special chars!"
3. Click "Call Tool" to see URL-encoded result

## Logging and Debugging

The server includes comprehensive logging to help diagnose connection and protocol issues.

### Log Categories

- `[MAIN]` - Server startup and initialization
- `[HTTP]` - HTTP connection lifecycle
- `[TOOL]` - Tool execution
- `[HEALTH]` - Health check requests
- `[WARN]` - Unknown endpoints or warnings

### Example Log Output

```
2026-01-20 10:30:15.123456 main.go:45: ========================================
2026-01-20 10:30:15.123789 main.go:46: MCP Streamable HTTP Server Starting
2026-01-20 10:30:15.123890 main.go:47: ========================================
2026-01-20 10:30:15.124001 main.go:53: [MAIN] Registering tools...
2026-01-20 10:30:15.124112 main.go:69: [MAIN] - Registered tool: base64_encode
2026-01-20 10:30:15.124223 main.go:76: [MAIN] - Registered tool: sha256_hash
2026-01-20 10:30:15.124334 main.go:83: [MAIN] - Registered tool: get_timestamp
2026-01-20 10:30:15.124445 main.go:90: [MAIN] - Registered tool: json_validate
2026-01-20 10:30:15.124556 main.go:97: [MAIN] - Registered tool: url_encode
2026-01-20 10:30:15.124667 main.go:100: [MAIN] All 5 tools registered successfully
```

### Viewing Logs in Kubernetes

```bash
# View logs from all pods
kubectl logs -l app=mcp-server-streamable-http --all-containers=true

# Follow logs in real-time
kubectl logs -f deployment/mcp-server-streamable-http

# View logs from specific pod
kubectl logs <pod-name>

# View logs with timestamps
kubectl logs <pod-name> --timestamps
```

## Usage

### Using the MCP Inspector (Recommended)

```bash
# Start the server
./mcp-server

# In another terminal, run the inspector
npx @modelcontextprotocol/inspector http://localhost:8081/mcp
```

### Using MCP Client Libraries

**TypeScript/JavaScript:**
```typescript
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StreamableHTTPClientTransport } from "@modelcontextprotocol/sdk/client/streamable-http.js";

const transport = new StreamableHTTPClientTransport(
  new URL("http://localhost:8081/mcp")
);
const client = new Client({ name: "my-client", version: "1.0.0" }, {});
await client.connect(transport);

// List tools
const tools = await client.listTools();

// Call a tool
const result = await client.callTool({
  name: "base64_encode",
  arguments: { text: "Hello, World!" }
});
```

**Python:**
```python
from mcp import ClientSession
from mcp.client.streamable_http import streamable_http_client

async with streamable_http_client("http://localhost:8081/mcp") as (read, write):
    async with ClientSession(read, write) as session:
        # Initialize
        await session.initialize()

        # List tools
        tools = await session.list_tools()

        # Call tool
        result = await session.call_tool("sha256_hash", {
            "text": "password123"
        })
```

## Tool Specifications

### 1. base64_encode

Encodes text to base64 format.

**Parameters:**
- `text` (string, required): The text to encode

**Example:**
```json
{
  "text": "Hello, World!"
}
```

### 2. sha256_hash

Generates SHA256 hash of the given text.

**Parameters:**
- `text` (string, required): The text to hash

**Example:**
```json
{
  "text": "password123"
}
```

### 3. get_timestamp

Returns the current Unix timestamp.

**Parameters:** None

### 4. json_validate

Validates if a string is valid JSON and returns structure info.

**Parameters:**
- `json_string` (string, required): The JSON string to validate

**Example:**
```json
{
  "json_string": "{\"name\": \"John\", \"age\": 30}"
}
```

### 5. url_encode

URL encodes the given text.

**Parameters:**
- `text` (string, required): The text to URL encode

**Example:**
```json
{
  "text": "hello world & special chars!"
}
```

## Protocol

This server uses the official MCP Go SDK which implements the MCP protocol version `2024-11-05` with Streamable HTTP transport. The SDK handles all JSON-RPC 2.0 communication and transport concerns.

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - Official MCP Go SDK
- `github.com/google/jsonschema-go` - JSON Schema generation

## Resources

- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [MCP Documentation](https://modelcontextprotocol.io/)
- [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector)

## License

MIT
