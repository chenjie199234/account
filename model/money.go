package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MoneyLog struct {
	LogID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"user_id"`
	Action      string             `bson:"action"`
	UniqueID    string             `bson:"unique_id"`
	SrcDst      string             `bson:"src_dst"`
	MoneyType   string             `bson:"money_type"`
	MoneyAmount uint32             `bson:"money_amount"`
	ExtData     string             `bson:"ext_data"`
}
