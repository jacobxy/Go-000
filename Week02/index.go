package main

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"
)

func AskSomeThing(name string) ([]interface{}, error) {
	return nil, sql.ErrNoRows
}

func Insert(name string) error {}
func Update(name string) error {}
func Delete(name string) error {}

func DAO(option string, name string) (interface{}, error) {
	r, err := AskSomeThing(name)
	switch option {
	case "query":
		return r, err
	case "insert":
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, Insert(name)
			}
			return nil, errors.Wrapf(err, "option:%s param:%s", option, name)
		} else {
			return r, errors.New("something aleardy exist")
		}
	case "update":
		_, err := AskSomeThing(name)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("Update param:%s NoRows", name)
				return nil, nil
			}
			return nil, errors.Wrapf(err, "option:%s param:%s", option, name)
		} else {
			return nil, Update(name)
		}
	case "delete":
		return nil, Delete(name)
	}
	return nil, nil
}
