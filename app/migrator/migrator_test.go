package migrator

import (
	"os"
	"testing"

	"github.com/unknwon/com"
	"gopkg.in/yaml.v2"
)

func prepare() error {
	hexoData := `---
title: Go Channels
date: 2017-07-16 18:51:47
categories: [Tech]
tags: [Golang, 并发]
---
this is a hexo post
`
	return com.WriteFile("PugoTest/Hexo/index.md", []byte(hexoData))
}

func TestYaml(t *testing.T) {
	data := `
title: Go Channels
date: 2017-07-16 18:51:47
categories: [Tech]
tags: [Golang, 并发]
`
	var header PostHeader
	err := yaml.Unmarshal([]byte(data), &header)
	if err != nil {
		t.Fatal(err)
	}
}
func TestMigrator(t *testing.T) {
	if err := prepare(); err != nil {
		t.Fatal(err)
	}

	os.MkdirAll("./PugoTest/Pugo", os.ModePerm)
	migrator := NewMigrator("./PugoTest/Hexo", "./PugoTest/Pugo")
	if err := migrator.Migrate(); err != nil {
		t.Fatal(err)
	}
}
