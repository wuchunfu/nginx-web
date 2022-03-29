package certx

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

func TestLego(t *testing.T) {
	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	myUser := MyUser{
		Email: "you@yours.com",
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	// This CA URL is configured for a local dev instance of Boulder running in Dockerfile in a VM.
	config.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// We specify an HTTP port of 5002 and an TLS port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 9180 and 9181.
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "9180"))
	if err != nil {
		log.Fatal(err)
	}
	//err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "9181"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}

	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{"mydomain.com"},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)

	err = os.WriteFile("fullchain.cer", certificates.Certificate, 0744)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("private.key", certificates.PrivateKey, 0744)
	if err != nil {
		log.Fatal(err)
	}
}
