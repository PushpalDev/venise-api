package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/dernise/base-api/server"
	"github.com/dernise/base-api/services"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

func main() {
	api := server.API{Router: gin.Default(), Config: viper.New()}

	api.SetupViper()

	api.EmailSender = services.NewSendGridEmailSender(api.Config)

	session, err := mgo.Dial(api.Config.GetString("db_host"))
	if err != nil {
		panic(err)
	}
	defer session.Close()
	api.Database = session.DB(api.Config.GetString("db_name"))

	govalidator.SetFieldsRequiredByDefault(true)

	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
