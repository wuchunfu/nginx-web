package certx

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/wuchunfu/nginx-web/middleware/logx"
	"github.com/wuchunfu/nginx-web/utils/nginxx"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// MyUser You'll need a user or account type that implements acme.User
type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func GetCertInfo(domain string) (key *x509.Certificate) {
	ts := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: ts}
	response, err := client.Get("https://" + domain)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(response.Body)
	key = response.TLS.PeerCertificates[0]
	return
}

func IssueCert(domain string) error {
	// Create a user. New accounts need an email and private key to start.
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Println(err)
		return err
	}

	myUser := MyUser{
		Email: "you@yours.com",
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)

	if false {
		// This CA URL is configured for a local dev instance of Boulder running in Dockerfile in a VM.
		config.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}
	config.Certificate.KeyType = certcrypto.RSA2048

	// A client facilitates communication with the CA server.
	client, err := lego.NewClient(config)
	if err != nil {
		log.Println(err)
		return err
	}

	// We specify an HTTP port of 5002 and an TLS port of 5001 on all interfaces
	// because we aren't running as root and can't bind a listener to port 80 and 443
	// (used later when we attempt to pass challenges). Keep in mind that you still
	// need to proxy challenge traffic to port 5002 and 5001.
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "9180"))
	if err != nil {
		log.Println(err)
		return err
	}
	//err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "9181"))
	//if err != nil {
	//	log.Fatal(err)
	//}

	// New users will need to register
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Println(err)
		return err
	}

	myUser.Registration = reg

	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Println(err)
		return err
	}

	saveDir, pathErr := nginxx.GetConfPath("ssl/" + domain)
	if pathErr != nil {
		logx.GetLogger().Sugar().Errorf("Get nginx conf path failed: %s", pathErr.Error())
		return pathErr
	}

	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		err = os.Mkdir(saveDir, 0755)
		if err != nil {
			log.Println("fail to create", saveDir)
			return err
		}
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	err = os.WriteFile(filepath.Join(saveDir, "fullchain.cer"), certificates.Certificate, 0744)
	if err != nil {
		log.Println(err)
		return err
	}

	err = os.WriteFile(filepath.Join(saveDir, domain+".key"), certificates.PrivateKey, 0744)
	if err != nil {
		log.Println(err)
		return err
	}

	output, err := nginxx.Reload()
	if err != nil {
		logx.GetLogger().Sugar().Infof("Reload nginx failed: %s", output)
		logx.GetLogger().Sugar().Errorf("Reload nginx failed: %s", err.Error())
		return err
	}
	return nil
}
