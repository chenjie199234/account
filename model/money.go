package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MoneyLog struct {
	LogID       bson.ObjectID `bson:"_id,omitempty"`
	UserID      bson.ObjectID `bson:"user_id"`
	Action      string        `bson:"action"` //spend,recharge,refund
	UniqueID    string        `bson:"unique_id"`
	SrcDst      string        `bson:"src_dst"`
	MoneyType   string        `bson:"money_type"`
	MoneyAmount uint32        `bson:"money_amount"`
}
