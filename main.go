package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			conn := conn
			var totalBytes []byte
			for !isReqComplete(totalBytes) {
				bytes := make([]byte, 1024)
				k, err := conn.Read(bytes)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(k)
				totalBytes = append(totalBytes, bytes[:k]...)
				fmt.Println(string(totalBytes))
			}

			parsedReq, err := parseRequest(totalBytes)
			if err != nil {
				return
			}

			switch parsedReq.Endpoint {

			}

			res := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nhello"

			conn.Write([]byte(res))
			conn.Close()
		}()
	}
}

func parseRequest(req []byte) (*Request, error) {
	reqString := string(req)
	firstLineInd := strings.Index(reqString, "\r\n")
	if firstLineInd == -1 {
		panic("")
	}

	firstLine := reqString[:firstLineInd]

	methodEndInd := strings.Index(firstLine, " ")
	if methodEndInd == -1 {
		panic("")
	}
	method := firstLine[:methodEndInd]
	fmt.Println(method)

	endpointInd := strings.Index(firstLine[methodEndInd+1:], " ")
	if endpointInd == -1 {
		panic("")
	}
	endpoint := firstLine[methodEndInd+1:][:endpointInd]
	fmt.Println(endpoint)

	protocolInd := strings.Index(firstLine[methodEndInd+1:][endpointInd+1:], "/")
	if protocolInd == -1 {
		panic("")
	}
	protocol := firstLine[methodEndInd+1:][endpointInd+1:][:protocolInd]
	fmt.Println(protocol)

	protocolVersion := firstLine[methodEndInd+1:][endpointInd+1:][protocolInd+1:]
	fmt.Println(protocolVersion)

	bodyInd := strings.Index(reqString, "\r\n\r\n")
	if bodyInd == -1 {
		panic("")
	}

	headersString := req[firstLineInd+2 : bodyInd]
	fmt.Println(headersString)
	headers := parseHeaders(headersString)
	body := req[bodyInd+4:]

	fmt.Println(headers)
	fmt.Println(body)

	return &Request{
		Endpoint:        endpoint,
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
		Body:            body,
		Headers:         headers,
	}, nil
}

func isReqComplete(req []byte) bool {
	if len(req) < 4 {
		return false
	}

	if bytes.Contains(req, []byte("\r\n\r\n")) {
		// now body starts
		parsed, _ := parseRequest(req)
		if parsed == nil {
			return false
		}

		if val, ok := parsed.Headers["Content-Length"]; ok {
			if valInt, err := strconv.Atoi(val); err == nil {
				if len(parsed.Body) >= valInt {
					return true
				}
				// request partial
				return false
			}
		}

		// no body expected
		return true
	}

	return false
}

func parseHeaders(headersString []byte) map[string]string {
	ans := map[string]string{}

	headers := bytes.Split(headersString, []byte("\n"))
	for _, header := range headers {
		keyInd := bytes.Index(header, []byte(": "))
		if keyInd == -1 {
			panic("")
		}
		key := header[:keyInd]
		value := header[keyInd+2:]
		if value[len(value)-1] == '\r' {
			value = value[:len(value)-1]
		}

		ans[string(key)] = string(value)
	}

	return ans
}

type Request struct {
	Endpoint        string
	Protocol        string
	ProtocolVersion string
	Body            []byte
	Headers         map[string]string
}
