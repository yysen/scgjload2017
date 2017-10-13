package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var vconfig config

type detailTab struct {
	Name  string
	Sizes []int
	PreID string
}
type mainTable struct {
	Name        string
	FieldSize   []int
	ID          string
	Result      string
	Status      string
	Shadowtable string
	Maxrows     int
	Details     []detailTab
}
type config struct {
	AreaCode    string
	URL         string
	FinishOut   string
	UserName    string
	Password    string
	Driver      string
	DBURL       string
	Table       []mainTable
}

func readConfig(fileName string) error {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(bs, &vconfig); err != nil {
		return err
	}
	return nil
}
