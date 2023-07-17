package k8simageadmissioncontroller

import (
	"crypto/tls"
	"fmt"
	"net/http"

	logging "github.com/sirupsen/logrus"
)

func RunWebhookServer(hostname string, port int, certFile, keyFile string, compressedImageSizeName int64) {
	logging.Debug("Loading certFile and keyFile ...")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		logging.Panic(err)
	}

	logging.Info(fmt.Sprintf("Serving on https://%s:%d ...", hostname, port))

	router := http.NewServeMux()
	router.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) { ValidatePod(w, r, compressedImageSizeName) })

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%d", hostname, port),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		Handler: WithLogging(router),
	}

	if err := server.ListenAndServeTLS("", ""); err != nil {
		logging.Panic(err)
	}
}
