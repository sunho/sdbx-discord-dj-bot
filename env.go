package djbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

func e(str string) error {
	return errors.New(str)
}

type EnvVar struct {
	Var     interface{}
	Type    stypes.Type
	ForUser bool
}

type UpdatableVar interface {
	Get(string) string
	Set(string, string) error
}

type Updatable struct {
	Var *UpdatableVar
}

type EnvOwner struct {
	Env map[string]EnvVar
}

type EnvManager struct {
	Owner           map[string]*EnvOwner
	UpdateFunctions []func()
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
	for _, f := range base.UpdateFunctions {
		f()
	}
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
	base.updateEnv()
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
	if env, ok := base.Env[key]; ok {
		return env.Var, nil
	}
	fmt.Println(key + msg.Nil)
	return nil, e(msg.AcessUndefinedEnv)
}

func (base *EnvManager) MakeDefaultEnv(key string, i interface{}, foruser bool) error {
	defaultow := base.Owner["default"]
	t := stypes.GetType(i)
	if t == stypes.TypeOther {
		return e(msg.ConvertingError)
	}
	defaultow.Env[key] = EnvVar{i, t, foruser}
	return nil
}

// only for users
func (base *EnvOwner) SetEnvWithStr(key string, value string) error {
	if env, ok := base.Env[key]; ok {
		i, err := stypes.TypeConvertOne(value, env.Type)
		if err != nil {
			return err
		}
		base.Env[key] = EnvVar{i, env.Type, true}
		return nil
	}
	return e(msg.AcessUndefinedEnv)
}

func (base *EnvOwner) SetEnvWithInterface(key string, in interface{}) error {
	if env, ok := base.Env[key]; ok {
		if stypes.GetType(in) == env.Type {
			base.Env[key] = EnvVar{in, env.Type, env.ForUser}
		} else {
			return e(msg.ConvertingError)
		}
		return nil
	}
	if t := stypes.GetType(in); t != stypes.TypeOther {
		base.Env[key] = EnvVar{in, t, false}
	}
	return e(msg.NoSupportOther)
}

func NewEnvManager() EnvManager {
	return EnvManager{make(map[string]*EnvOwner), []func(){}}
}
func (base *EnvOwner) DeleteEnv(key string) {
	delete(base.Env, key)
}
