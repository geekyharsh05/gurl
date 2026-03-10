# gurl - Modern HTTP CLI Client 

A fast, reliable HTTP client with JSON support and intuitive syntax. Inspired by curl but with additional modern features.

## Features

- Simple and intuitive command syntax
- Support for all HTTP methods
- Custom headers and request body
- JSON formatting for requests and responses
- Response output to file
- Follow redirects with configurable limits
- Timing information
- Insecure TLS connections for testing
- Color-coded output

## Installation

### Install on macOS or Linux (with Homebrew)

```bash
brew tap geekyharsh05/gurl
brew install gurl
```

### Install on Windows (with Scoop)

```bash
scoop bucket add geekyharsh05 https://github.com/geekyharsh05/scoop-bucket
scoop install gurl
```

### Install on Linux (Debian/Ubuntu)

```bash
wget https://github.com/geekyharsh05/gurl/releases/download/v1.0.1/gurl_1.0.1_linux_amd64.deb
sudo apt install ./gurl_1.0.1_linux_amd64.deb
```

### Using Go

```bash
go install github.com/geekyharsh05/gurl@latest
```

### Build from Source

```bash
git clone https://github.com/geekyharsh05/gurl.git
cd gurl
go build
```

## Usage Examples

### Basic GET request:

```bash
gurl request https://dummyjson.com/products
```

### Download files:

```bash
gurl request -o myfile.zip https://example.com/file.zip
```

### Using HTTP methods:

```bash
gurl request -X POST https://dummyjson.com/products
gurl request -X PUT https://dummyjson.com/products/1
```

### Send JSON data:

```bash
gurl request https://dummyjson.com/auth/login \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"username":"emilys","password":"emilyspass","expiresInMins":30}'
```

### Custom headers:

```bash
gurl request https://dummyjson.com/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN_HERE"
```

### Follow redirects:

```bash
gurl request -L https://dummyjson.com/http/200
```

### Format JSON response:

```bash
gurl request -j https://dummyjson.com/products/1
```

### Save response to file:

```bash
gurl request -o response.json https://dummyjson.com/users
```

### Verbose output with timing information:

```bash
gurl request -v https://dummyjson.com/test
```

### Ignore SSL certificate errors:

```bash
gurl request -k https://dummyjson.com/test
```

### Set wait-time:

```bash
gurl request --wait-time 5s https://dummyjson.com/test
```

### Verbose output:

```bash
gurl request -v --wait-time 3s https://dummyjson.com/test
```

## Download Command

The download command offers wget-like functionality for downloading files with progress tracking, resume capability, and more:

### Basic file download:

```bash
gurl download https://example.com/file.zip
```

### Specify output filename:

```bash
gurl download -o myfile.zip https://example.com/file.zip
```

### Save to specific directory:

```bash
gurl download -P ./downloads/ https://example.com/file.zip
```

### Resume partial download:

```bash
gurl download -c https://example.com/large-file.iso
```

### Quiet mode (no progress bar):

```bash
gurl download -q https://example.com/file.zip
```

### Retry on failure:

```bash
gurl download --tries 5 --retry-delay 5s https://example.com/file.zip
```

### Download with custom timeout:

```bash
gurl download -t 60s https://example.com/large-file.iso
```

## Options

### Global Flags:

- `-k, --insecure`: Allow insecure server connections
- `--timeout duration`: Request timeout (default 30s)
- `--max-redirects int`: Maximum number of redirects to follow (default 10)
- `--wait-time duration`: Wait for the specified duration before making the request

### Request Command Flags:

- `-X, --method string`: HTTP method (default "GET")
- `-d, --data string`: Request body
- `-H, --header strings`: Custom headers
- `-v, --verbose`: Show verbose output
- `-o, --output string`: Write response to file
- `-j, --json`: Format response as JSON
- `-L, --follow`: Follow redirects
- `--json-request`: Set Content-Type to application/json
- `--form`: Set Content-Type to application/x-www-form-urlencoded
- `--no-pretty`: Disable automatic JSON formatting

### Download Command Flags:

- `-o, --output string`: Save file with the specified name
- `-P, --directory string`: Save files to the specified directory
- `-c, --continue`: Resume getting a partially-downloaded file
- `-q, --quiet`: Quiet mode - don't show progress bar
- `-k, --insecure`: Allow insecure server connections
- `-t, --timeout duration`: Set timeout for download (default 30s)
- `-n, --no-redirect`: Don't follow redirects
- `--max-redirects int`: Maximum number of redirects to follow (default 10)
- `-r, --tries int`: Number of retry attempts (default 3, 0 for no retries)
- `--retry-delay duration`: Delay between retries (default 2s)
