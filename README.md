# TCP Proxy Through HTTP Proxy (CONNECT)

This project is a simple TCP proxy written in Go that forwards TCP traffic through an HTTP CONNECT proxy. It supports any kind of TCP traffic, allowing you to connect to a remote host and port via an intermediate HTTP proxy server. It creates a local port for you to connect to, and any TCP traffic sent to this local port will be forwarded to the remote TCP destination.

## Features

- Listens on a local port and forwards incoming TCP connections to a remote host and port
- Supports connecting through an HTTP CONNECT proxy
- Configurable via command-line flags

## Usage

To run the TCP proxy, use the following command:

```
./tcp-proxy-through-http-proxy [flags]
```

The available flags are:

- `-local-port`: Local port to listen on (default: "8088")
- `-remote-host`: Remote host to connect to (default: "target-mysql-host")
- `-remote-port`: Remote port to connect to (default: "3306")
- `-proxy-url`: HTTP proxy URL (default: "http://squid-host:3128")

For example, to start the TCP proxy listening on port 9000, connecting to `target-mysql-host:3306`, and using the proxy URL `http://proxy.example.com:8080`, run:

```
./tcp-proxy-through-http-proxy -local-port=9000 -remote-host=target-mysql-host -remote-port=3306 -proxy-url=http://proxy.example.com:8080
```

## Configuration

The TCP proxy can be configured using command-line flags. Here's a description of each flag:

- `-local-port`: Specifies the local port on which the TCP proxy should listen for incoming connections. Default is "8088".
- `-remote-host`: Specifies the remote host to which the TCP proxy should forward the incoming connections. Default is "target-mysql-host".
- `-remote-port`: Specifies the remote port to which the TCP proxy should forward the incoming connections. Default is "3306".
- `-proxy-url`: Specifies the URL of the HTTP CONNECT proxy through which the TCP proxy should establish the connection. Default is "http://squid-host:3128".

## Alternative: socat

An alternative to using this TCP proxy is to use the `socat` utility. You can achieve similar functionality with the following command:

```
socat TCP-LISTEN:8088,reuseaddr,fork PROXY:proxy-host:target-host:target-port,proxyport=3128
```

However, this Go-based TCP proxy provides additional features and flexibility compared to socat:

- It supports HTTPS scheme for the HTTP proxy, which socat does not support.
- It is easier to extend and customize since it is written in Go.

## License

This project is licensed under the [MIT License](LICENSE).