# SSH Transport: SSH over WebSocket/Tls in Go

The SSH Transport project aims to provide SSH functionality over WebSocket/TLS connections using Go.

## Features

* Supports both command-line flags and file-based configuration.
* Works with secure and unsecure WebSocket and WSS connections.
* Designed for high performance with a low memory footprint.
* Supports SSH over TLS protocol.

## How to build

1. Clone the repository
2. Build the project using one of the following methods:
    * Run `goreleaser release`.
    * Use `go build`.
    * Use `go build -ldflags "-w -s"` for smaller binary size.

## Usage

### Using CLI Arguments as Configuration

1. Execute `ssh-transport`.
2. Sample: `ssh-transport -ssh=127.0.0.1:22 -ws-listen=0.0.0.0:80`.
3. Use the `-help` argument to get detailed usage information.

### Using a File for Configuration

1. Create a config.json file in the same directory as the ssh-transport binary.
2. Fill in the configuration details [config.json](./config.json).
3. Run `ssh-transport` without any cli arguments.
4. To disable WS or TLS leave the listen address empty like this `"listen": ""`

## NGINX reverse proxy websocket

```NGINX
location /sshws {
    proxy_pass http://127.0.0.1:5153/;
    proxy_redirect off;
    proxy_http_version 1.1;
    proxy_set_header Upgrade websocket;
    proxy_set_header Connection Upgrade;
    proxy_set_header Host $http_host;
    proxy_set_header Sec-WebSocket-Version $http_sec_websocket_version;
    proxy_set_header Sec-WebSocket-Key $http_sec_websocket_key;
    proxy_read_timeout 52w;
}
```

## License

This project is licensed under the Apache-2.0 License. See the LICENSE file for details.
