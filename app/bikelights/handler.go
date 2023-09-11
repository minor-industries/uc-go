package bikelights

import (
	"github.com/pkg/errors"
	"os"
	"uc-go/pkg/protocol/rpc"
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
			//go func() {
			//	//a.Logs.Log("setting config")
			//	msg := &cfg.Config{}
			//	a.Logs.Log("no")
			//	return
			//	_, err := msg.UnmarshalMsg(body)
			//	if err != nil {
			//		a.Logs.Log("no")
			//		return
			//	}
			//	a.Logs.Log("unmarshalled")
			//	//a.Cfg.Edit(func(c *cfg.Config) {
			//	//	*c = *msg
			//	//})
			//	//ss := a.Cfg.SnapShot()
			//	//if err := storage.WriteMsgp(a.Logs, a.Lfs, &ss, configFile); err != nil {
			//	//	a.Logs.Error(errors.Wrap(err, "save config"))
			//	//}
			//	//a.Logs.Log("set config")
			//}()
			return nil
		}),
	}
}
