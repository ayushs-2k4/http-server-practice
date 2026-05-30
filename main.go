package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
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
	firstLineInd := bytes.Index(req, []byte("\r\n"))
	if firstLineInd == -1 {
		panic("")
	}

	firstLine := req[:firstLineInd]

	methodEndInd := bytes.Index(firstLine, []byte(" "))
	if methodEndInd == -1 {
		panic("")
	}
	method := firstLine[:methodEndInd]

	endpointInd := bytes.Index(firstLine[methodEndInd+1:], []byte(" "))
	if endpointInd == -1 {
		panic("")
	}
	endpoint := firstLine[methodEndInd+1:][:endpointInd]

	protocolInd := bytes.Index(firstLine[methodEndInd+1:][endpointInd+1:], []byte("/"))
	if protocolInd == -1 {
		panic("")
	}
	protocol := firstLine[methodEndInd+1:][endpointInd+1:][:protocolInd]

	protocolVersion := firstLine[methodEndInd+1:][endpointInd+1:][protocolInd+1:]

	bodyInd := bytes.Index(req, []byte("\r\n\r\n"))
	if bodyInd == -1 {
		panic("")
	}

	headersString := req[firstLineInd+2 : bodyInd]
	headers := parseHeaders(headersString)
	body := req[bodyInd+4:]

	return &Request{
		Endpoint:        endpoint,
		Method:          method,
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
	Endpoint        []byte
	Method          []byte
	Protocol        []byte
	ProtocolVersion []byte
	Body            []byte
	Headers         map[string]string
}
