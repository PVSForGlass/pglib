package pglib

import (
	"encoding/gob"
	"io"
	"log"
	"net/rpc"
)

type FileData struct {
	Name     string
	Contents []byte
}

type Api interface {
	UploadFile(data FileData) (bool, error)
}

func (ac *apiClient) UploadFile(data FileData) (res bool, err error) {
	err = ac.c.Call("PGAPI.UploadFile", data, &res)
	return res, err
}

func ServeApi(rcvr Api, rwc io.ReadWriteCloser) {
	if err := rpc.RegisterName("PGAPI", &apiServer{rcvr}); err != nil {
		log.Panic(err)
	}
	rpc.ServeConn(rwc)
}

func ConnectApi(rwc io.ReadWriteCloser) Api {
	return &apiClient{rpc.NewClient(rwc)}
}

type apiServer struct {
	rcvr Api
}

func (as *apiServer) UploadFile(data FileData, res *bool) (err error) {
	*res, err = as.rcvr.UploadFile(data)
	return err
}

type apiClient struct {
	c *rpc.Client
}

func init() {
	gob.Register(FileData{})
}
