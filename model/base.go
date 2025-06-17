package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	UserID   bson.ObjectID     `bson:"_id,omitempty"`
	Password string            `bson:"password"`
	IDCard   string            `bson:"idcard"`
	Tel      string            `bson:"tel"`
	Email    string            `bson:"email"`
	OAuths   map[string]string `bson:"oauths"` //key service name,value unique id in this service
	Money    map[string]int32  `bson:"money"`
	BTime    uint64            `bson:"btime"`   //unit nanoseconds,if this is 0 means not banned,>0 means banned
	BReason  string            `bson:"breason"` //ban reason
}
type UserTelIndex struct {
	Tel    string        `bson:"tel"`
	UserID bson.ObjectID `bson:"user_id"`
}
type UserEmailIndex struct {
	Email  string        `bson:"email"`
	UserID bson.ObjectID `bson:"user_id"`
}
type UserIDCardIndex struct {
	IDCard string        `bson:"idcard"`
	UserID bson.ObjectID `bson:"user_id"`
}
type UserOAuthIndex struct {
	Service string        `bson:"service"` //service_name+'|'+service_id
	UserID  bson.ObjectID `bson:"user_id"`
}
