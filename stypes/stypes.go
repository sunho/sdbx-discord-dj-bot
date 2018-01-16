package stypes

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/errormsg"
)

type Type int

const (
	TypeInt Type = iota
	TypeString
	TypeStrings
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
	case Range:
		return TypeRange
	default:
		return TypeOther
	}
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
		return nil, errors.New(errormsg.NotEnoughMinerals)
	}
	switch t {
	case TypeInt:
		i, err := strconv.Atoi(nstr[0])
		if err != nil {
			return nil, errors.New(errormsg.ConvertingError)
		}
		return i, nil
	case TypeString:
		return nstr[0], nil
	case TypeStrings:
		return nstr, nil
	default:
		return nil, errors.New(errormsg.UndefinedType)
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
	return nil, errors.New(errormsg.SoEnoughMinerals)
}
