package main

import (
	"encoding/json"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var jsonpath = "./content.json"

func LoadContent() {
	iot, err := ioutil.ReadFile(jsonpath) //读取文件
	if checkErr(err) {                    //做错误判断
		json.Unmarshal(iot, &MbText)
	}
}

func SaveContent() {
	if js, err := json.MarshalIndent(MbText, "", "	"); checkErr(err) {
		err := ioutil.WriteFile(jsonpath, js, 0664) //接收文件名，字节类型的数据，文件的权限
		checkErr(err)
	}
}

func LoadConfig() {
	//读取yaml文件到缓存中
	config, err := ioutil.ReadFile("spider.yaml")
	checkErr(err)
	//yaml文件内容影射到结构体中
	err = yaml.Unmarshal(config, &cfg)
	checkErr(err)
}

type Config struct {
	Account struct {
		Cookie string  `yaml:"cookie"`
		Topic  *string `yaml:"topic"`
		UID    int     `yaml:"uid"`
	} `yaml:"account"`
	Message struct {
		Comment string `yaml:"comment"`
		Online  string `yaml:"online"`
		Weibo   string `yaml:"weibo"`
	} `yaml:"message"`
	Ws struct {
		Host string `yaml:"host"`
		Post int    `yaml:"post"`
	} `yaml:"ws"`
}
