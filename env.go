package djbot

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/errormsg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type EnvVar struct {
	Var  interface{}
	Type stypes.Type
}

type EnvOwner struct {
	Env map[string]EnvVar
}
type EnvManager struct {
	Owner map[string]*EnvOwner
}

func (base *EnvManager) GetOwner(owner string) *EnvOwner {
	if _, ok := base.Owner[owner]; !ok {
		base.copyDefaultEnv(owner)
	}
	return base.Owner[owner]
}

func marshelJson(i interface{}, filename string) error {
	saveJson, _ := json.Marshal(i)
	err := ioutil.WriteFile(filename, saveJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

//TODO? change this into io.writer
func (base EnvManager) Save(filename string) error {
	saveJson, _ := json.Marshal(base)
	err := ioutil.WriteFile(filename, saveJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

//TODO? change this into io.writer
func (base *EnvManager) Load(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, base)
	if err != nil {
		return err
	}
	return nil
}
func (base *EnvManager) updateEnv() {
	defaultow := base.Owner["default"]
	for _, owner := range base.Owner {
		for key := range owner.Env {
			//delete old env
			if _, ok := defaultow.Env[key]; !ok {
				delete(owner.Env, key)
			}
		}
		for key := range defaultow.Env {
			//update new env
			if _, ok := owner.Env[key]; !ok {
				owner.Env[key] = defaultow.Env[key]
			}
			//update new type
			if owner.Env[key].Type != defaultow.Env[key].Type {
				owner.Env[key] = defaultow.Env[key]
			}
		}
	}
}
func (base *EnvManager) copyDefaultEnv(owner string) {
	base.Owner[owner] = &EnvOwner{make(map[string]EnvVar)}
	for key, iter := range base.Owner["default"].Env {
		base.Owner[owner].Env[key] = iter
	}
}

func (base *EnvOwner) GetEnv(key string) (interface{}, error) {
	sar := interface{}(nil)
	if env, ok := base.Env[key]; ok {
		sar = env.Var
	}
	if sar == nil {
		return nil, errors.New(errormsg.AcessUndefinedEnv)
	}

	return sar, nil
}

func (base *EnvManager) MakeDefaultEnv(key string, i interface{}, t stypes.Type) error {
	defaultow := base.Owner["default"]
	defaultow.Env[key] = EnvVar{i, t}
	return nil
}

func (base *EnvOwner) SetEnvWithStr(key string, value string) error {
	if env, ok := base.Env[key]; ok {
		i, err := stypes.TypeConvertOne(value, env.Type)
		if err != nil {
			return err
		}
		base.Env[key] = EnvVar{i, env.Type}
		return nil
	}

	return errors.New(errormsg.AcessUndefinedEnv)
}

func (base *EnvOwner) SetEnvWithInterface(key string, in interface{}) error {
	if env, ok := base.Env[key]; ok {
		if stypes.GetType(in) == env.Type {
			base.Env[key] = EnvVar{in, env.Type}
		} else {
			return errors.New(errormsg.ConvertingError)
		}
		return nil
	}
	return errors.New(errormsg.AcessUndefinedEnv)
}
