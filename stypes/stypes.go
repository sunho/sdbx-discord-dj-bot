package stypes

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
)

type Type int

const (
	TypeInt Type = iota
	TypeString
	TypeStrings
	TypeBool
	TypeRange
	TypeOther
)

type Range struct {
}

func GetType(in interface{}) Type {
	switch in.(type) {
	case int:
		return TypeInt
	case string:
		return TypeString
	case []string:
		return TypeStrings
	case bool:
		return TypeBool
	case Range:
		return TypeRange
	default:
		return TypeOther
	}
}

func e(str string) error {
	return errors.New(str)
}

func TypeConvertOne(str string, t Type) (interface{}, error) {
	s := strings.Split(str, ",")
	return typeConvertOne(s, t)
}

func typeConvertOne(nstr []string, t Type) (interface{}, error) {
	//exception for commandmanager
	if len(nstr) == 0 {
		if t == TypeStrings {
			return []string{""}, nil
		}
		return nil, e(msg.NotEnoughMinerals)
	}
	switch t {
	case TypeInt:
		i, err := strconv.Atoi(nstr[0])
		if err != nil {
			return nil, e(msg.TypesDontMatch)
		}
		return i, nil
	case TypeString:
		return nstr[0], nil
	case TypeStrings:
		return nstr, nil
	case TypeBool:
		if nstr[0] == "true" {
			return true, nil
		} else if nstr[0] == "false" {
			return false, nil
		} else {
			return nil, e(msg.TypesDontMatch)
		}
	default:
		return nil, e(msg.UndefinedType)
	}
}

func TypeConvertMany(strs []string, types []Type) ([]interface{}, error) {
	//void commands
	if len(types) == 0 {
		return nil, nil
	}
	nstrs := strs
	ia := []interface{}{}
	for _, t := range types {
		i, err := typeConvertOne(nstrs, t)
		if err != nil {
			return nil, err
		}
		ia = append(ia, i)
		//TODO do something better
		//strings
		if reflect.TypeOf(i) == reflect.TypeOf([]string{}) {
			return ia, nil
		}
		nstrs = nstrs[1:]
	}
	if len(types) == len(strs) {
		return ia, nil
	}
	return nil, e(msg.SoEnoughMinerals)
}
