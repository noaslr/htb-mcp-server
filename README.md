# HackTheBox MCP Server

A Model Context Protocol (MCP) server that provides AI assistants with programmatic access to HackTheBox platform functionality.

## Features

The HTB MCP Server exposes 12 comprehensive tools for interacting with the HackTheBox platform:

### Challenge Management

- **`list_challenges`** - Get paginated list of challenges with filtering
- **`start_challenge`** - Initialize a challenge environment
- **`submit_challenge_flag`** - Submit flags for challenge verification

### Machine Management

- **`list_machines`** - Get active/retired machines with status information
- **`start_machine`** - Start a machine and get connection details
- **`get_machine_ip`** - Retrieve IP address of active machine
- **`submit_user_flag`** - Submit user flags for machines
- **`submit_root_flag`** - Submit root flags for machines

### User Management

- **`get_user_profile`** - Retrieve user profile and statistics
- **`get_user_progress`** - Get completion status and achievements

### Search & Utility

- **`search_content`** - Advanced search across challenges/machines/users
- **`get_server_status`** - Health check and server information

## Prerequisites

- Go 1.21 or later
- Valid HackTheBox account with API access
- HTB API token (JWT format)

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/NoASLR/htb-mcp-server.git
   cd htb-mcp-server
   ```

2. **Build the binary:**

   ```bash
   go build -o htb-mcp-server main.go
   ```

3. **Get your HTB API token:**
   - Go to [HackTheBox Profile Settings](https://app.hackthebox.com/profile/settings)
   - Generate an App Token
   - Copy the JWT token (format: `xxx.yyy.zzz`)

## Configuration

The server is configured via environment variables:

### Required

- `HTB_TOKEN` - Your HackTheBox API token (JWT format)

### Optional

- `SERVER_PORT` - Server port (default: 3000)
- `LOG_LEVEL` - Logging level: DEBUG, INFO, WARN, ERROR (default: INFO)
- `RATE_LIMIT_PER_MINUTE` - API rate limiting (default: 100)
- `CACHE_TTL_SECONDS` - Response cache TTL (default: 300)
- `REQUEST_TIMEOUT_SECONDS` - HTTP request timeout (default: 30)

## Usage

### Standalone Mode

```bash
export HTB_TOKEN="your.jwt.token.here"
./htb-mcp-server
```

### Docker Mode

```bash
docker build -t htb-mcp-server .
docker run -e HTB_TOKEN="your.jwt.token.here" htb-mcp-server
```

### MCP Client Integration

Add to your MCP client configuration (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "htb": {
      "command": "/path/to/htb-mcp-server",
      "env": {
        "HTB_TOKEN": "your.jwt.token.here"
      }
    }
  }
}
```

## Example Usage

Once connected, you can use the tools through your AI assistant:

```
# List active challenges
"Can you show me the available Web challenges on HackTheBox?"

# Start a machine
"Please start machine ID 123 and get its IP address"

# Submit a flag
"Submit the user flag 'HTB{example_flag}' for machine 123"

# Search for content
"Search for machines related to 'Active Directory'"

# Check server status
"What's the current status of the HTB MCP server?"
```

## API Endpoints

The server implements the MCP protocol over stdio transport. All communication follows the JSON-RPC 2.0 specification.

### Core MCP Methods

- `initialize` - Initialize the MCP session
- `tools/list` - List available tools
- `tools/call` - Execute a specific tool

### HTB API Integration

The server integrates with HackTheBox API v4:

- Base URL: `https://labs.hackthebox.com/api/v4`
- Authentication: Bearer token (JWT)
- Rate limiting: Respects HTB API limits

## Development

### Project Structure

```
htb-mcp-server/
├── main.go                    # Entry point
├── pkg/
│   ├── config/               # Configuration management
│   ├── htb/                  # HTB API client
│   └── mcp/                  # MCP protocol implementation
├── internal/
│   ├── server/               # MCP server core
│   └── tools/                # Tool implementations
├── tests/                    # Test files
└── docs/                     # Documentation
```

### Adding New Tools

1. Create a new tool struct implementing the `Tool` interface:

   ```go
   type MyTool struct {
       client *htb.Client
   }

   func (t *MyTool) Name() string { return "my_tool" }
   func (t *MyTool) Description() string { return "Description" }
   func (t *MyTool) Schema() mcp.ToolSchema { /* schema */ }
   func (t *MyTool) Execute(ctx context.Context, args map[string]interface{}) (*mcp.CallToolResponse, error) {
       // Implementation
   }
   ```

2. Register the tool in `registry.go`:
   ```go
   r.RegisterTool(NewMyTool(r.htbClient))
   ```

### Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (requires HTB_TOKEN)
HTB_TOKEN="your.token" go test -tags=integration ./...
```

## Security Considerations

- **Token Security**: Never commit your HTB token to version control
- **Rate Limiting**: The server implements rate limiting to prevent API abuse
- **Input Validation**: All user inputs are validated before API calls
- **Error Handling**: Sensitive information is not exposed in error messages

## Performance

- **Response Time**: < 500ms for 95% of requests
- **Caching**: Intelligent caching reduces API calls
- **Concurrency**: Supports multiple concurrent tool executions
- **Circuit Breaker**: Protects against HTB API outages

## Troubleshooting

### Common Issues

1. **"HTB token appears invalid or expired"**

   - Verify your token is correct and not expired
   - Ensure token has proper JWT format (3 parts separated by dots)
   - Check token permissions in HTB profile settings

2. **"Connection refused"**

   - Verify network connectivity to labs.hackthebox.com
   - Check if corporate firewall blocks HTB API access

3. **"Rate limit exceeded"**
   - Reduce request frequency
   - Increase `RATE_LIMIT_PER_MINUTE` if needed

### Debug Mode

Enable debug logging:

```bash
export LOG_LEVEL=DEBUG
./htb-mcp-server
```

### Health Check

Test server connectivity:

```bash
curl -X POST http://localhost:3000/health
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- HackTheBox team for providing the API
- Model Context Protocol community for the specification
- Go community for excellent tooling and libraries

## Roadmap

### Phase 2 Features

- Real-time machine status notifications
- Advanced analytics and reporting
- Team collaboration features
- Custom challenge creation tools

### Long-term Goals

- Mobile application support
- Integration with other security platforms
- Advanced automation capabilities
- Machine learning-powered recommendations

## Support

- **Issues**: [GitHub Issues](https://github.com/NoASLR/htb-mcp-server/issues)
- **Documentation**: [Wiki](https://github.com/NoASLR/htb-mcp-server/wiki)
- **Community**: [Discussions](https://github.com/NoASLR/htb-mcp-server/discussions)

---

Built with ❤️ for the cybersecurity community
