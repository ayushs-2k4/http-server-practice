package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
			if errors.Is(err, io.EOF) {
				fmt.Println("req closed")
				conn.Close()
				return
			}
			log.Fatal(err)
		}
		go func() {
			conn := conn
			var req Request
			req.HeadersCompleted = false
			req.RequestCompleted = false
			var totalBytes []byte

			for {
				for !req.RequestCompleted {
					if len(totalBytes) != 0 {
						parseRequestIncrementally(&req, &totalBytes)
						if req.RequestCompleted {
							break
						}
					}
					reqBytes := make([]byte, 1024)
					k, err := conn.Read(reqBytes)
					if err != nil {
						if errors.Is(err, io.EOF) {
							fmt.Println("req closed")
							conn.Close()
							return
						}
						log.Fatal(err)
					}
					reqBytes = reqBytes[:k]
					totalBytes = append(totalBytes, reqBytes...)
					parseRequestIncrementally(&req, &totalBytes)
				}

				// req got completed, there may be some more requests lined up
				fmt.Println(req)
				req = Request{}

				res := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nhello"

				conn.Write([]byte(res))
				//conn.Close()
			}

		}()
	}
}

func parseRequestIncrementally(req *Request, reqBytes *[]byte) {
	if req == nil {
		panic("req is nil")
	}
	if reqBytes == nil {
		panic("req bytes is nil")
	}

	if req.Method == nil {
		methodEndInd := bytes.Index(*reqBytes, []byte(" "))
		if methodEndInd == -1 {
			// even method name didn't came
			return
		}

		methodName := (*reqBytes)[:methodEndInd]
		*reqBytes = (*reqBytes)[methodEndInd+1:]

		req.Method = methodName
	}

	if req.Endpoint == nil {
		endpointEndInd := bytes.Index(*reqBytes, []byte(" "))
		if endpointEndInd == -1 {
			// endpoint has not came completely yet
			return
		}

		endpoint := (*reqBytes)[:endpointEndInd]
		*reqBytes = (*reqBytes)[endpointEndInd+1:]

		req.Endpoint = endpoint
	}

	if req.Protocol == nil {
		protocolEndInd := bytes.Index(*reqBytes, []byte("/"))
		if protocolEndInd == -1 {
			//protocol hasn't came yet completely
			return
		}

		protocol := (*reqBytes)[:protocolEndInd]
		*reqBytes = (*reqBytes)[protocolEndInd+1:]

		req.Protocol = protocol
	}

	if req.ProtocolVersion == nil {
		protocolVersionEndInd := bytes.Index(*reqBytes, []byte("\r\n"))
		if protocolVersionEndInd == -1 {
			//protocol version hasn't came yet completely
			return
		}

		protocolVersion := (*reqBytes)[:protocolVersionEndInd]
		*reqBytes = (*reqBytes)[protocolVersionEndInd+2:]

		req.ProtocolVersion = protocolVersion
	}

	// headers
	if !req.HeadersCompleted {
		for {
			headerEndInd := bytes.Index(*reqBytes, []byte("\r\n"))
			if headerEndInd == -1 {
				//this specific hasn't came yet completely
				return
			}
			if headerEndInd == 0 {
				// headers ended
				*reqBytes = (*reqBytes)[headerEndInd+2:]
				req.HeadersCompleted = true
				break
			}

			headerBytes := (*reqBytes)[:headerEndInd]
			*reqBytes = (*reqBytes)[headerEndInd+2:]

			key, value := parseHeader(headerBytes)
			if req.Headers == nil {
				req.Headers = make(map[string]string)
			}
			req.Headers[string(key)] = string(value)
		}
	}

	if contentLength, ok := req.Headers["Content-Length"]; ok {
		contentLengthInt, _ := strconv.Atoi(contentLength)

		if len(*reqBytes) >= contentLengthInt {
			req.Body = (*reqBytes)[:contentLengthInt]
			req.RequestCompleted = true
		}
	}

	req.RequestCompleted = true
}

func parseHeaders(headersString []byte) map[string]string {
	ans := map[string]string{}

	headers := bytes.Split(headersString, []byte("\n"))
	for _, header := range headers {
		key, value := parseHeader(header)

		ans[string(key)] = string(value)
	}

	return ans
}

func parseHeader(headerString []byte) ([]byte, []byte) {
	keyInd := bytes.Index(headerString, []byte(": "))
	if keyInd == -1 {
		panic("")
	}
	key := headerString[:keyInd]
	value := headerString[keyInd+2:]
	if value[len(value)-1] == '\r' {
		value = value[:len(value)-1]
	}

	return key, value
}

type Request struct {
	Endpoint         []byte
	Method           []byte
	Protocol         []byte
	ProtocolVersion  []byte
	Body             []byte
	Headers          map[string]string
	HeadersCompleted bool
	RequestCompleted bool
}
