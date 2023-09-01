package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID   primitive.ObjectID `bson:"_id,omitempty"`
	Password string             `bson:"password"`
	IDCard   string             `bson:"idcard"`
	NickName string             `bson:"nick_name"`
	Tel      string             `bson:"tel"`
	Email    string             `bson:"email"`
	Money    map[string]int32   `bson:"money"`
}
type UserTelIndex struct {
	Tel    string             `bson:"tel"`
	UserID primitive.ObjectID `bson:"user_id"`
}
type UserEmailIndex struct {
	Email  string             `bson:"email"`
	UserID primitive.ObjectID `bson:"user_id"`
}
type UserIDCardIndex struct {
	IDCard string             `bson:"idcard"`
	UserID primitive.ObjectID `bson:"user_id"`
}
type UserNickNameIndex struct {
	NickName string             `bson:"nick_name"`
	UserID   primitive.ObjectID `bson:"user_id"`
}
