package stypes

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/errormsg"
)

type Type int

const (
	TypeInt Type = iota
	TypeString
	TypeRange
	TypeStrings
)

type Range struct {
}

func InterfacesToStrings(slice []interface{}) ([]string, error) {
	sa := []string{}
	for _, i := range slice {
		if s, ok := i.(string); ok {
			sa = append(sa, s)
		} else {
			return nil, errors.New(errormsg.NoString)
		}
	}
	return sa, nil
}

func StringsToInterfaces(slice []string) ([]interface{}, error) {
	ia := []interface{}{}
	for _, s := range slice {
		ia = append(ia, s)
	}
	return ia, nil
}

func TypeConvertOne(str string, t Type) (interface{}, error) {
	s := []string{str}
	return typeConvertOne(s, t)
}

func typeConvertOne(nstr []string, t Type) (interface{}, error) {
	switch t {
	case TypeString:
		return nstr[0], nil
	case TypeInt:
		i, err := strconv.Atoi(nstr[0])
		if err != nil {
			return nil, errors.New(errormsg.ConvertingError)
		}
		return i, nil
	case TypeStrings:
		ia2, _ := StringsToInterfaces(nstr)
		return ia2, errors.New("strings") //TODO: do something better
	default:
		return nil, errors.New(errormsg.UndefinedType)
	}
}

func TypeConvertMany(strs []string, types []Type) ([]interface{}, error) {
	if len(types) == 0 {
		return []interface{}{}, nil
	}
	fmt.Println(strs)
	nstrs := strs
	ia := []interface{}{}
	for _, t := range types {
		if len(nstrs) == 0 {
			return nil, errors.New(errormsg.NotEnoughMinerals)
		}
		i, err := typeConvertOne(nstrs, t)
		if err != nil {
			if err.Error() == "strings" {
				ia = append(ia, i.([]interface{})...)
				return ia, nil
			} else {
				return nil, err
			}
		}
		ia = append(ia, i)
		nstrs = nstrs[1:]
	}
	if len(types) == len(strs) {
		return ia, nil
	}
	return nil, errors.New(errormsg.SoEnoughMinerals)
}
