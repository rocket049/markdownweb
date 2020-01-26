package main

import (
	"github.com/kataras/iris"

	"encoding/json"
	"io/ioutil"
)

type KeysTLS struct {
	Addr   string `json:"addr"`
	Domain string `json:"domain"`
	Email  string `json:"email"`
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
	return iris.AutoTLS(keys.Addr, keys.Domain, keys.Email), nil
}
