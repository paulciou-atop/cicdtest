package security

import (
	"io/ioutil"
	"log"
	"net/http"
	localca "security/ca"

	"security/config"

	"github.com/smallstep/certificates/ca"
	"github.com/urfave/cli"
	"go.step.sm/cli-utils/step"
)

const rootName = "root.crt"
const srvCsrName = "srv.crt"
const srvkeyName = "srv.key"

func NewAtopPki() *atopPki {
	removeAll(step.ProfilePath())
	return &atopPki{}
}

type atopPki struct {
	srv         *ca.CA
	fingerprint string
}

func (a *atopPki) LoadConfiguration() {
	f := NewFile(config.GetpwdFileName())
	err := f.Save(config.Getpassword())
	if err != nil {
		log.Fatal(err)
	}
	c, RootFinger, err := onboardPKI()
	if err != nil {
		log.Fatal(err)
	}
	c.Password = config.Getpassword()
	a.srv, err = ca.New(c)
	if err != nil {

		log.Fatal(err)
	}
	a.fingerprint = RootFinger
}

func (a *atopPki) CABootstrapServer(ch chan bool) {

	go func() {
		ch <- true
		if err := a.srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

}

func (a *atopPki) Askcertificate(ch chan bool) {
	<-ch
	args := []string{"ca", "token", config.GetDns(), "--password-file=" + config.GetpwdFileName(),
		"--ca-url=" + config.GetCaUrl()}

	app := cli.NewApp()
	app.Commands = []cli.Command{localca.TokenCommand()}
	err := app.Run(args)
	if err != nil {
		log.Fatal(err)
	}
	token := config.GetToken()

	client, err := ca.NewClient(config.GetCaUrl(), ca.WithRootSHA256(a.fingerprint))
	if err != nil {
		log.Fatal(err)
	}
	root, err := client.Root(a.fingerprint)
	if err != nil {
		log.Fatal(err)
	}

	req, pk, err := ca.CreateSignRequest(token)
	if err != nil {
		log.Fatal(err)
	}
	sign, err := client.Sign(req)
	if err != nil {
		log.Fatal(err)
	}

	savePEM(srvkeyName, pk)
	savePEM(rootName, root.RootPEM)
	certPEM, err := getPEM(sign.ServerPEM)
	if err != nil {
		log.Fatal(err)
	}
	caPEM, err := getPEM(sign.CaPEM)
	if err != nil {
		log.Fatal(err)
	}
	chain := append(certPEM, caPEM...)
	saveFile(srvCsrName, chain)

}

func (a *atopPki) Stop() {
	err := a.srv.Stop()
	if err != nil {
		log.Fatal(err)
	}
}

func savePEM(filename string, ca interface{}) error {
	b, err := getPEM(ca)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0644)

}

func saveFile(filename string, b []byte) error {

	return ioutil.WriteFile(filename, b, 0644)

}
