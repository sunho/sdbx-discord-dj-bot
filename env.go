package djbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

func e(str string) error {
	return errors.New(str)
}

type EnvVar struct {
	Var  interface{}
	Type stypes.Type
}

type EnvServer struct {
	sync.Mutex
	Env map[string]EnvVar
	ID  string
}

type EnvManager struct {
	sync.Mutex
	Servers map[string]*EnvServer
}

func (envm *EnvManager) GetServer(serverID string) *EnvServer {
	if _, ok := envm.Servers[serverID]; !ok {
		envm.copyDefaultEnv(serverID)
	}
	return envm.Servers[serverID]
}

func (envm *EnvManager) Save(filename string) {
	bytes, err := json.Marshal(envm)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filename, bytes, 0777)
	if err != nil {
		return
	}
}

func (envm *EnvManager) Load(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, envm)
	if err != nil {
		return
	}
}

func (envm *EnvManager) Update() {
	defaultserver := envm.Servers["default"]
	for _, server := range envm.Servers {
		for key := range server.Env {
			//delete old env
			if _, ok := defaultserver.Env[key]; !ok {
				delete(server.Env, key)
			}
		}
		for key := range defaultserver.Env {
			//add new env
			if _, ok := server.Env[key]; !ok {
				server.Env[key] = defaultserver.Env[key]
			}
			//update new type
			if server.Env[key].Type != defaultserver.Env[key].Type {
				server.Env[key] = defaultserver.Env[key]
			}
		}
	}
}

func (envm *EnvManager) copyDefaultEnv(serverID string) {
	envm.Lock()
	defer func() {
		envm.Unlock()
	}()
	envm.Servers[serverID] = &EnvServer{
		Env: make(map[string]EnvVar),
		ID:  serverID,
	}
	for key, env := range envm.Servers["default"].Env {
		envm.Servers[serverID].Env[key] = env
	}
}

func (envs *EnvServer) GetEnv(key string) interface{} {
	if env, ok := envs.Env[key]; ok {
		return env.Var
	}
	fmt.Println(key, "doesn't exist.")
	return nil
}

func (envm *EnvManager) MakeDefaultEnv(key string, value interface{}) {
	envm.Lock()
	defer func() {
		envm.Unlock()
	}()
	defaultserver := envm.Servers["default"]
	typ := stypes.GetType(value)
	if typ == stypes.TypeOther {
		return
	}
	defaultserver.Env[key] = EnvVar{value, typ}
}

func (envs *EnvServer) SetEnvWithStr(key string, value string) error {
	envs.Lock()
	defer func() {
		envs.Unlock()
	}()
	if env, ok := envs.Env[key]; ok {
		convvalue, err := stypes.TypeConvertOne(value, env.Type)
		if err != nil {
			return err
		}
		envs.Env[key] = EnvVar{convvalue, env.Type}
		return nil
	}

	return e(msg.AcessUndefinedEnv)
}

func (envs *EnvServer) SetEnvWithInterface(key string, value interface{}) error {
	envs.Lock()
	defer func() {
		envs.Unlock()
	}()
	if env, ok := envs.Env[key]; ok {
		if stypes.GetType(value) == env.Type {
			envs.Env[key] = EnvVar{value, env.Type}
		} else {
			return e(msg.TypesDontMatch)
		}
		return nil
	}
	return e(msg.AcessUndefinedEnv)
}

func NewEnvManager() EnvManager {
	return EnvManager{
		Servers: make(map[string]*EnvServer),
	}
}
