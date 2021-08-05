package setup

import (
	"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/util"
)

func Auth(prov provider.Provider) {
	session, success := util.GetCred(prov.Name())

	if !success {
		session = authFlow()
	}

	_ = session
}

func authFlow() string{

}