package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

func handleConnection(clientConn net.Conn, proxyURL *url.URL, targetAddress string) {
	defer clientConn.Close()

	var proxyConn net.Conn
	var err error

	// Connect to the proxy server
	if proxyURL.Scheme == "https" {
		// For HTTPS proxies, use TLS
		proxyConn, err = tls.Dial("tcp", proxyURL.Host, &tls.Config{
			InsecureSkipVerify: true, // Note: This skips certificate verification. Use with caution.
		})
	} else {
		// For HTTP proxies, use regular TCP connection
		proxyConn, err = net.Dial("tcp", proxyURL.Host)
	}

	if err != nil {
		fmt.Printf("Failed to connect to proxy: %v\n", err)
		return
	}
	defer proxyConn.Close()

	connectReq := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Opaque: targetAddress},
		Host:   targetAddress,
		Header: make(http.Header),
	}

	// Add basic auth header if credentials are provided
	if proxyURL.User != nil {
		username := proxyURL.User.Username()
		password, _ := proxyURL.User.Password()
		auth := username + ":" + password
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		connectReq.Header.Add("Proxy-Authorization", basicAuth)
	}

	err = connectReq.Write(proxyConn)
	if err != nil {
		fmt.Printf("Failed to write CONNECT request: %v\n", err)
		return
	}
	resp, err := http.ReadResponse(bufio.NewReader(proxyConn), connectReq)
	if err != nil {
		fmt.Printf("Failed to read CONNECT response: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Proxy CONNECT response status: %v\n%s\n", resp.Status, body)
		return
	}
	go io.Copy(proxyConn, clientConn)
	io.Copy(clientConn, proxyConn)
}

func main() {
	localPort := flag.String("local-port", "8088", "Local port to listen on")
	remoteHost := flag.String("remote-host", "target-mysql-host", "Remote host to connect to")
	remotePort := flag.String("remote-port", "3306", "Remote port to connect to")
	proxyURL := flag.String("proxy-url", "http://squid-host:3128", "HTTP proxy URL")

	flag.Parse()

	listenAddr := fmt.Sprintf("0.0.0.0:%s", *localPort)
	parsedProxyURL, err := url.Parse(*proxyURL)
	if err != nil {
		fmt.Printf("Failed to parse proxy URL: %v\n", err)
		return
	}
	targetAddress := fmt.Sprintf("%s:%s", *remoteHost, *remotePort)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("Failed to start listener: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", listenAddr)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleConnection(clientConn, parsedProxyURL, targetAddress)
	}
}
