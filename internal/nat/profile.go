package nat

import (
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type Queue []string

const (
	ReferralEventName = "referral-notification"
	ReferralTitleKey  = "referral-notification"
	ReferralMessage   = "referral message from user profile"
)

type Event struct {
	EventName string `json:"event_name,omitempty" bson:"event_name,omitempty"`
	UserId    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Data      Data   `json:"data,omitempty" bson:"data,omitempty"`
}

type Data struct {
	TitleKey string `json:"title_key,omitempty" bson:"title_key,omitempty"`
	Code     string `json:"code,omitempty" bson:"code,omitempty"`
	UserId   string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	Message  string `json:"message,omitempty" bson:"message,omitempty"`
	Symbol   string `json:"symbol,omitempty" bson:"symbol,omitempty"`
	Amount   string `json:"amount,omitempty" bson:"amount,omitempty"`
	Tx       string `json:"tx,omitempty" bson:"tx,omitempty"`
}

func NewEvent(refCode, userId string) Event {
	return Event{EventName: ReferralEventName,
		UserId: userId,
		Data:   Data{Code: refCode, Tx: uuid.NewString(), UserId: userId},
	}
}

func NewData(code, userId string) Data {
	return Data{Code: code, Tx: uuid.NewString(), UserId: userId}
}

func PublishNewProfile(refCode, userId string) {
	event := NewData(refCode, userId)
	data, err := json.Marshal(event)
	if err != nil {
		log.Errorf("Publish message has an error")
		return
	}

	// TODO: refactor nats connect code
	err = natClient.Publish(NatCfg.Topic.NewProfileTopic, data)
	if err != nil {
		log.Errorf("Publish message has an error")
		return
	}

	// Flush the connection to ensure that the message is sent
	natClient.Flush()

	if err := natClient.LastError(); err != nil {
		log.Errorf("error flushing connection: %v", err)
		return
	}

	log.Infof("Published message: %s", data)
}
