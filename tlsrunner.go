package main

import (
	iris "github.com/kataras/iris/v12"

	"encoding/json"
	"io/ioutil"
)

type KeysTLS struct {
	Addr string `json:"addr"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func getTLSKeys() (*KeysTLS, error) {
	data, err := ioutil.ReadFile("tls.json")
	if err != nil {
		return nil, err
	}
	var res KeysTLS
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func getTLSRunner() (iris.Runner, error) {
	keys, err := getTLSKeys()
	if err != nil {
		return nil, err
	}
	return iris.TLS(keys.Addr, keys.Cert, keys.Key), nil
}
