package main

import (
	"log"
	"net/http"
)

func main() {
	config := LoadConfig()

	if config.Tls.Listen != "" {
		go SshTls(config.Tls.Listen, config.SshAddr)
	}
	if config.Ws.Listen != "" {
		http.HandleFunc(config.Ws.Path, handler(config.SshAddr))

		if config.Ws.Tls {
			log.Println("SSH over WSS Listening on ", config.Ws.Listen)
			http_server_error := http.ListenAndServeTLS(config.Ws.Listen, config.Ws.CertPath, config.Ws.KeyPath, nil)

			if http_server_error != nil {
				log.Fatalln(http_server_error.Error())
			}
		} else {
			log.Println("SSH over websocket Listening on ", config.Ws.Listen)
			http_server_error := http.ListenAndServe(config.Ws.Listen, nil)

			if http_server_error != nil {
				log.Fatalln(http_server_error.Error())
			}
		}
	}
}

func handler(sshAddr string) func(w http.ResponseWriter, r *http.Request) {
	wsHandler := func(w http.ResponseWriter, r *http.Request) {
		// Two different types of websockets exist: the secure type, which includes the `Sec-WebSocket-Key` and `Sec-WebSocket-Version`.
		// And the unsecure type, which does not encrypt any data.
		if r.Header["Sec-Websocket-Key"] != nil {
			SecureWS(w, r, sshAddr)
		} else {
			UnsecureWS(w, sshAddr)
		}
	}

	return wsHandler
}
