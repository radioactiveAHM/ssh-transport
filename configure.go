package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type TransportWs struct {
	Listen   string `json:"listen"`
	Path     string `json:"path"`
	Tls      bool   `json:"tls"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

type TransportTls struct {
	Listen   string `json:"listen"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

type Config struct {
	SshAddr string       `json:"SshAddr"`
	Ws      TransportWs  `json:"ws"`
	Tls     TransportTls `json:"tls"`
}

func LoadConfig() Config {
	// TODO: Load config from Command-Line Flags
	ssh := flag.String("ssh", "127.0.0.1:22", "SSH address")
	// WS
	wsListen := flag.String("ws-listen", "", "Websocket listen address")
	wsPath := flag.String("ws-path", "/", "Websocket path")
	wsTls := flag.Bool("ws-tls", false, "Enable TLS for websocket")
	wsCert := flag.String("ws-cert", "cert.pem", "Certificate path for wss")
	wsKey := flag.String("ws-key", "key.pem", "Key path for wss")
	// TLS
	tlsListen := flag.String("tls-listen", "", "Websocket listen address")
	tlsCert := flag.String("tls-cert", "cert.pem", "Certificate path for tls")
	tlsKey := flag.String("tls-key", "key.pem", "Key path for tls")

	flag.Parse()

	if len(os.Args) == 1 {
		// Load config from file
		cfile, cfile_err := os.ReadFile("config.json")
		if cfile_err != nil {
			log.Fatalln(cfile_err.Error())
		}

		conf := Config{}
		conf_err := json.Unmarshal(cfile, &conf)
		if conf_err != nil {
			log.Fatalln(conf_err.Error())
		}
		return conf
	}

	return Config{
		SshAddr: *ssh,
		Ws: TransportWs{
			Listen:   *wsListen,
			Path:     *wsPath,
			Tls:      *wsTls,
			CertPath: *wsCert,
			KeyPath:  *wsKey,
		},
		Tls: TransportTls{
			Listen:   *tlsListen,
			CertPath: *tlsCert,
			KeyPath:  *tlsKey,
		},
	}
}
