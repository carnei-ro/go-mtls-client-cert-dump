package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var root_ca = flag.String("root_ca", "root-ca.crt", "Provide the file name for Client Root CA Certificate")
var server_cert = flag.String("server_cert", "server.crt", "Provide the file name for Server Certificate")
var server_key = flag.String("server_key", "server.key", "Provide the file name for Server Key")
var mtls = flag.Bool("mtls", true, "Enable Mutual authentication")

func main() {
	flag.Parse()

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(rw, ":8443/client-cert to dump client certificate\n")
	})

	http.HandleFunc("/client-cert", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		if cs := r.TLS; len(cs.PeerCertificates) != 0 {
			fmt.Fprintf(rw, "Length of PeerCertificates: %#v\n", len(cs.PeerCertificates))
			fmt.Fprint(rw, "----------------------------------------------------------------\n\n")
			for i := 0; i < len(cs.PeerCertificates); i++ {
				pem.Encode(rw, &pem.Block{Type: "CERTIFICATE", Bytes: cs.PeerCertificates[i].Raw})
			}
		} else {
			fmt.Fprintf(rw, "This is not a mTLS connection. Length of PeerCertificates: %#v\n", len(cs.PeerCertificates))
		}
	})

	server := &http.Server{Addr: ":8443"}
	if *mtls {
		server = createServerWithMTLS()
	}

	go http.ListenAndServe(":8081", nil)
	// Start the server loading the certificate and key
	err := server.ListenAndServeTLS(*server_cert, *server_key)
	if err != nil {
		log.Fatal("Unable to start server", err)
	}
}

func createServerWithMTLS() *http.Server {
	flag.Parse()
	// Add the cert chain as the intermediate signs both the servers and the clients certificates
	clientCACert, err := ioutil.ReadFile(*root_ca)
	if err != nil {
		log.Fatal(err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	return &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}
}
