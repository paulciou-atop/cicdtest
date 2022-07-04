package security

import (
	"io"
	"os"
)

func InitCa() {

	atoppki := NewAtopPki()
	atoppki.LoadConfiguration()

	ch := make(chan bool, 1)
	atoppki.CABootstrapServer(ch)
	atoppki.Askcertificate(ch)
	atoppki.Stop()

}

func NewPkiService() *PkiService {

	InitCa()
	return &PkiService{}
}

type PkiService struct {
	UnimplementedPkiServer
}

func (p *PkiService) GetRootCrt(n *EmptyParams, stream Pki_GetRootCrtServer) error {
	// Maximum 1KB size per stream.
	f, err := os.Open(rootName)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)

	for {
		num, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&ResponseData{Data: buf[:num]}); err != nil {
			return err
		}

	}

	return nil
}

func (p *PkiService) GetSrvCrt(n *EmptyParams, stream Pki_GetSrvCrtServer) error {
	// Maximum 1KB size per stream.
	f, err := os.Open(srvCsrName)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)

	for {
		num, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&ResponseData{Data: buf[:num]}); err != nil {
			return err
		}

	}

	return nil
}

func (p *PkiService) GetSrvKey(n *EmptyParams, stream Pki_GetSrvKeyServer) error {
	// Maximum 1KB size per stream.
	f, err := os.Open(srvkeyName)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)

	for {
		num, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&ResponseData{Data: buf[:num]}); err != nil {
			return err
		}

	}

	return nil
}
