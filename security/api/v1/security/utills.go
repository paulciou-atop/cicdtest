package security

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"

	localconfig "security/config"

	"github.com/pkg/errors"
	"github.com/smallstep/certificates/api"
	"github.com/smallstep/certificates/authority/config"
	"github.com/smallstep/certificates/cas/apiv1"
	"github.com/smallstep/certificates/pki"
	"go.step.sm/cli-utils/errs"
	"go.step.sm/cli-utils/fileutil"
)

func removeAll(path string) {
	e := os.RemoveAll(path)
	if e != nil {
		log.Fatal(e)
	}
}

func NewFile(name string) *File {
	return &File{
		name: name,
	}
}

type File struct {
	name string
}

func (f *File) Save(message string) error {
	if err := ioutil.WriteFile(f.name, []byte(message), 0644); err != nil {
		return err
	}

	return nil
}

func onboardPKI() (*config.Config, string, error) {
	pwd := []byte(localconfig.Getpassword())
	var opts = []pki.Option{
		pki.WithAddress(localconfig.GetAddress()),
		pki.WithDNSNames([]string{localconfig.GetDns()}),
		pki.WithProvisioner(localconfig.GetProvisioner()),
	}
	p, err := pki.New(apiv1.Options{
		Type:      apiv1.SoftCAS,
		IsCreator: true,
	}, opts...)
	if err != nil {
		return nil, "", err
	}
	name := localconfig.GetName()
	log.Println("Generating root certificate...")
	root, err := p.GenerateRootCertificate(name, name, name, pwd)
	if err != nil {
		return nil, "", err
	}
	log.Println("Generating intermediate certificate...")
	err = p.GenerateIntermediateCertificate(name, name, name, root, pwd)

	if err != nil {
		return nil, "", err
	}
	// Write files to disk
	if err := p.WriteFiles(); err != nil {
		return nil, "", err
	}

	log.Println("Generating admin provisioner...")
	if err := p.GenerateKeyPairs(pwd); err != nil {
		return nil, "", err
	}

	caConfig, err := p.GenerateConfig()
	if err != nil {
		return nil, "", err
	}

	b, err := json.MarshalIndent(caConfig, "", "   ")
	if err != nil {
		return nil, "", errors.Wrapf(err, "error marshaling %s", p.GetCAConfigPath())
	}
	if err := fileutil.WriteFile(p.GetCAConfigPath(), b, 0666); err != nil {
		return nil, "", errs.FileError(err, p.GetCAConfigPath())
	}

	return caConfig, p.GetRootFingerprint(), nil

}

func getPEM(i interface{}) ([]byte, error) {
	block := new(pem.Block)
	switch i := i.(type) {
	case api.Certificate:
		block.Type = "CERTIFICATE"
		block.Bytes = i.Raw
	case *x509.Certificate:
		block.Type = "CERTIFICATE"
		block.Bytes = i.Raw
	case *rsa.PrivateKey:
		block.Type = "RSA PRIVATE KEY"
		block.Bytes = x509.MarshalPKCS1PrivateKey(i)
	case *ecdsa.PrivateKey:
		var err error
		block.Type = "EC PRIVATE KEY"
		block.Bytes, err = x509.MarshalECPrivateKey(i)
		if err != nil {
			return nil, errors.Wrap(err, "error marshaling private key")
		}
	case ed25519.PrivateKey:
		var err error
		block.Type = "PRIVATE KEY"
		block.Bytes, err = x509.MarshalPKCS8PrivateKey(i)
		if err != nil {
			return nil, errors.Wrap(err, "error marshaling private key")
		}
	default:
		return nil, errors.Errorf("unsupported key type %T", i)
	}
	return pem.EncodeToMemory(block), nil
}
