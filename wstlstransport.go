package main

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func UnsecureWS(w http.ResponseWriter, sshAddr string) {
	h, ok := w.(http.Hijacker)
	if !ok {
		log.Println("response does not implement http.Hijacker")
		return
	}
	ws, _, err := h.Hijack()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Send Ok as websocket connection established
	_, okWrite_e := ws.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n"))
	if okWrite_e != nil {
		log.Println(okWrite_e.Error())
		return
	}

	// connect to ssh tcp socket
	ssh, ssh_e := net.Dial("tcp", sshAddr)
	if ssh_e != nil {
		log.Println(ssh_e.Error())
		return
	}
	defer ssh.Close()

	wsChan := make(chan []byte)
	sshChan := make(chan []byte)

	go HandleTCPRecv(ws, wsChan)
	go HandleTCPRecv(ssh, sshChan)

	for {
		select {
		case wsBuff, wsClosed := <-wsChan:
			if wsClosed {
				_, sshWrite_e := ssh.Write(wsBuff)
				if sshWrite_e != nil {
					log.Println(sshWrite_e.Error())
					return
				}
			} else {
				return
			}
		case sshBuff, sshClosed := <-sshChan:
			if sshClosed {
				_, wsWrite_e := ws.Write(sshBuff)
				if wsWrite_e != nil {
					log.Println(wsWrite_e.Error())
					return
				}
			} else {
				return
			}
		}
	}
}

func SecureWS(w http.ResponseWriter, r *http.Request, sshAddr string) {
	upgrader := websocket.Upgrader{}
	ws, ws_e := upgrader.Upgrade(w, r, nil)

	// if upgrade failed close
	if ws_e != nil {
		log.Println(ws_e.Error())
		return
	}
	defer ws.Close()

	// connect to ssh tcp socket
	ssh, ssh_e := net.Dial("tcp", sshAddr)
	if ssh_e != nil {
		log.Println(ssh_e.Error())
		return
	}
	defer ssh.Close()

	wsChan := make(chan []byte)
	sshChan := make(chan []byte)

	go HandleWSRecv(ws, wsChan)
	go HandleTCPRecv(ssh, sshChan)

	for {
		select {
		case wsBuff, wsClosed := <-wsChan:
			if wsClosed {
				_, sshWrite_e := ssh.Write(wsBuff)
				if sshWrite_e != nil {
					log.Println(sshWrite_e.Error())
					return
				}
			} else {
				return
			}
		case sshBuff, sshClosed := <-sshChan:
			if sshClosed {
				wsWrite_e := ws.WriteMessage(1, sshBuff)
				if wsWrite_e != nil {
					log.Println(wsWrite_e.Error())
					return
				}
			} else {
				return
			}
		}
	}
}

func HandleWSRecv(ws *websocket.Conn, wsChan chan []byte) {
	defer close(wsChan)

	for {
		_, buff, read_err := ws.ReadMessage()
		if read_err != nil {
			log.Println(read_err.Error())
			return
		}

		wsChan <- buff
	}
}

func HandleTCPRecv(tcp net.Conn, tcpChan chan []byte) {
	defer close(tcpChan)

	for {
		buff := make([]byte, 1024*16)
		buffLen, read_err := tcp.Read(buff)
		if read_err != nil {
			log.Println(read_err.Error())
			return
		}

		tcpChan <- buff[:buffLen]
	}
}
