package nat

import (
	"fmt"

	"github.com/nartvt/profile-service/internal/conf"
	"github.com/nats-io/nats.go"
)

var (
	natClient *nats.Conn
	NatCfg    NatServer
)

type NatServer struct {
	Url   string `json:"url,omitempty" bson:"url,omitempty"`
	Topic Topic  `json:"topic,omitempty" bson:"topic,omitempty"`
}

type Topic struct {
	NewProfileTopic string
}

func InitNats(config *conf.NatServer) {
	fmt.Println("NAT CONFIG: ", config.Topic)
	NatCfg = NatServer{
		Url: config.Url,
		Topic: Topic{
			NewProfileTopic: config.Topic.NewProfileTopic,
		},
	}
	fmt.Println("[NATS SERVER]: ", NatCfg.Url)
	nc, err := nats.Connect(NatCfg.Url)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		panic(err)
	}
	natClient = nc
}

func CloseNats() {
	if natClient != nil {
		natClient.Close()
	}
}
