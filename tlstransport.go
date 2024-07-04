package main

import (
	"crypto/tls"
	"log"
	"net"
)

func SshTls(tlsListenAddr string, sshAddr string) {
	log.Println("SSH over TLS Listening on ", tlsListenAddr)
	// Load certificate
	cert, cert_e := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if cert_e != nil {
		log.Println(cert_e.Error())
		return
	}

	// TLS Config
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	// TLS Listener
	tls_conn, tls_e := tls.Listen("tcp", tlsListenAddr, tlsConfig)
	if tls_e != nil {
		log.Println(tls_e.Error())
		return
	}

	for {
		tlsConn, tlsConn_e := tls_conn.Accept()
		if tlsConn_e != nil {
			log.Println(tlsConn_e.Error())
			continue
		}

		go tlsConnHandler(tlsConn, sshAddr)
	}
}

func tlsConnHandler(tlsConn net.Conn, sshAddr string) {
	// connect to ssh tcp socket
	ssh, ssh_e := net.Dial("tcp", sshAddr)
	if ssh_e != nil {
		log.Println(ssh_e.Error())
		return
	}
	defer ssh.Close()

	tlsConnChan := make(chan []byte)
	sshChan := make(chan []byte)

	go HandleTCPRecv(tlsConn, tlsConnChan)
	go HandleTCPRecv(ssh, sshChan)

	for {
		select {
		case tlsConnBuff, tlsConnClosed := <-tlsConnChan:
			if tlsConnClosed {
				_, sshWrite_e := ssh.Write(tlsConnBuff)
				if sshWrite_e != nil {
					log.Println(sshWrite_e.Error())
					return
				}
			} else {
				return
			}
		case sshBuff, sshClosed := <-sshChan:
			if sshClosed {
				_, tlsConnWrite_e := tlsConn.Write(sshBuff)
				if tlsConnWrite_e != nil {
					log.Println(tlsConnWrite_e.Error())
					return
				}
			} else {
				return
			}
		}
	}
}
