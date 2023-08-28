package bikelights

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"uc-go/app/bikelights/cfg"
	"uc-go/pkg/protocol/rpc"
	"uc-go/pkg/storage"
)

func (a *App) Handlers() map[string]rpc.Handler {
	return map[string]rpc.Handler{
		"get-config": rpc.HandlerFunc(func(method string, body []byte) error {
			ss := a.Cfg.SnapShot()
			if err := rpc.Send(os.Stdout, "show-config", &ss); err != nil {
				a.Logs.Error(errors.Wrap(err, "send show-config"))
			}
			return nil
		}),

		"reset-config": rpc.HandlerFunc(func(method string, body []byte) error {
			if err := a.Lfs.Remove(a.ConfigFile()); err != nil {
				a.Logs.Error(errors.Wrap(err, "remove file"))
			} else {
				a.Logs.Log("config reset")
			}
			return nil

		}),

		"set-config": rpc.HandlerFunc(func(method string, body []byte) error {
			fmt.Println("config:")
			msg := &cfg.Config{}
			_, err := msg.UnmarshalMsg(body)
			if err != nil {
				return errors.Wrap(err, "unmarshal")
			}
			a.Cfg.Edit(func(c *cfg.Config) {
				*c = *msg
			})
			ss := a.Cfg.SnapShot()
			if err := storage.WriteMsgp(a.Logs, a.Lfs, &ss, configFile); err != nil {
				a.Logs.Error(errors.Wrap(err, "save config"))
			}
			a.Logs.Log("set config")
			return nil
		}),
	}
}
