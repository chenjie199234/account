package user

import (
	"context"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoGetUserByUserID(ctx context.Context, userid primitive.ObjectID) (*model.User, error) {
	user := &model.User{}
	if e := d.mongo.Database("account").Collection("user").FindOne(ctx, bson.M{"_id": userid}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return user, nil
}
func (d *Dao) MongoGetUserByTel(ctx context.Context, tel string) (*model.User, error) {
	telindex := &model.UserTelIndex{}
	if e := d.mongo.Database("account").Collection("user_tel_index").FindOne(ctx, bson.M{"tel": tel}).Decode(telindex); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return d.MongoGetUserByUserID(ctx, telindex.UserID)
}
func (d *Dao) MongoGetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	emailindex := &model.UserEmailIndex{}
	if e := d.mongo.Database("account").Collection("user_email_index").FindOne(ctx, bson.M{"email": email}).Decode(emailindex); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return d.MongoGetUserByUserID(ctx, emailindex.UserID)
}
func (d *Dao) MongoGetUserByIDCard(ctx context.Context, idcard string) (*model.User, error) {
	idcardindex := &model.UserIDCardIndex{}
	if e := d.mongo.Database("account").Collection("user_idcard_index").FindOne(ctx, bson.M{"idcard": idcard}).Decode(idcardindex); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return d.MongoGetUserByUserID(ctx, idcardindex.UserID)
}
func (d *Dao) MongoGetUserByNickName(ctx context.Context, nickname string) (*model.User, error) {
	nicknameindex := &model.UserNickNameIndex{}
	if e := d.mongo.Database("account").Collection("user_nick_name_index").FindOne(ctx, bson.M{"nick_name": nickname}).Decode(nicknameindex); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return d.MongoGetUserByUserID(ctx, nicknameindex.UserID)

}
func (d *Dao) MongoUpdateUserTel(ctx context.Context, userid primitive.ObjectID, newTel string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = sctx.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	user := &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOne(sctx, bson.M{"_id": userid}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.Tel == newTel {
		return
	}
	if _, e = d.mongo.Database("account").Collection("user_tel_index").InsertOne(sctx, bson.M{"tel": newTel, "user_id": userid}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrTelAlreadyUsed
		}
		return
	}
	if _, e = d.mongo.Database("account").Collection("user_tel_index").DeleteOne(sctx, bson.M{"tel": user.Tel, "user_id": userid}); e != nil {
		return
	}
	_, e = d.mongo.Database("account").Collection("user").UpdateOne(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"tel": newTel}})
	return
}
func (d *Dao) MongoUpdateUserEmail(ctx context.Context, userid primitive.ObjectID, newEmail string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = sctx.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	user := &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOne(sctx, bson.M{"_id": userid}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.Email == newEmail {
		return
	}
	if _, e = d.mongo.Database("account").Collection("user_email_index").InsertOne(sctx, bson.M{"email": newEmail, "user_id": userid}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrEmailAlreadyUsed
		}
		return
	}
	if _, e = d.mongo.Database("account").Collection("user_email_index").DeleteOne(sctx, bson.M{"email": user.Email, "user_id": userid}); e != nil {
		return
	}
	_, e = d.mongo.Database("account").Collection("user").UpdateOne(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"email": newEmail}})
	return
}
func (d *Dao) MongoUpdateUserIDCard(ctx context.Context, userid primitive.ObjectID, newIDCard string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = sctx.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	user := &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOne(sctx, bson.M{"_id": userid}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.IDCard == newIDCard {
		return
	}
	//实名认证后无法更改
	if user.IDCard != "" {
		e = ecode.ErrIDCardAlreadySetted
		return
	}
	if _, e = d.mongo.Database("account").Collection("user_idcard_index").InsertOne(sctx, bson.M{"idcard": newIDCard, "user_id": userid}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrEmailAlreadyUsed
		}
		return
	}
	_, e = d.mongo.Database("account").Collection("user").UpdateOne(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"idcard": newIDCard}})
	return
}
func (d *Dao) MongoUpdateUserNickName(ctx context.Context, userid primitive.ObjectID, newNickName string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = sctx.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	user := &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOne(sctx, bson.M{"_id": userid}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.NickName == newNickName {
		return
	}
	if _, e = d.mongo.Database("account").Collection("user_nick_name_index").InsertOne(sctx, bson.M{"nick_name": newNickName, "user_id": userid}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrNickNameAlreadyUsed
		}
		return
	}
	if _, e = d.mongo.Database("account").Collection("user_nick_name_index").DeleteOne(sctx, bson.M{"nick_name": user.NickName, "user_id": userid}); e != nil {
		return
	}
	_, e = d.mongo.Database("account").Collection("user").UpdateOne(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"nick_name": newNickName}})
	return
}
func (d *Dao) MongoUpdateUserPassword(ctx context.Context, userid primitive.ObjectID, oldpassword, newpassword string) (e error) {

}
