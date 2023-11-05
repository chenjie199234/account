package user

import (
	"context"
	"crypto/rand"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"
	"github.com/chenjie199234/account/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoCreateUserByTel(ctx context.Context, tel string) (user *model.User, e error) {
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
			if mongo.IsDuplicateKeyError(e) {
				user, e = d.MongoGetUserByTel(ctx, tel)
			}
		} else if e = sctx.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	nonce := make([]byte, 16)
	rand.Read(nonce)
	user = &model.User{
		Password: util.SignMake("", nonce),
		IDCard:   "",
		NickName: "",
		Tel:      tel,
		Email:    "",
		Money:    map[string]int32{},
	}
	var r *mongo.InsertOneResult
	if r, e = d.mongo.Database("account").Collection("user").InsertOne(sctx, user); e != nil {
		return
	}
	user.UserID = r.InsertedID.(primitive.ObjectID)
	_, e = d.mongo.Database("account").Collection("user_tel_index").InsertOne(sctx, &model.UserTelIndex{
		Tel:    tel,
		UserID: user.UserID,
	})
	return
}
func (d *Dao) MongoCreateUserByEmail(ctx context.Context, email string) (user *model.User, e error) {
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
			if mongo.IsDuplicateKeyError(e) {
				user, e = d.MongoGetUserByEmail(ctx, email)
			}
		} else if e = sctx.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	nonce := make([]byte, 16)
	rand.Read(nonce)
	user = &model.User{
		Password: util.SignMake("", nonce),
		IDCard:   "",
		NickName: "",
		Tel:      "",
		Email:    email,
		Money:    map[string]int32{},
	}
	var r *mongo.InsertOneResult
	if r, e = d.mongo.Database("account").Collection("user").InsertOne(sctx, user); e != nil {
		return
	}
	user.UserID = r.InsertedID.(primitive.ObjectID)
	_, e = d.mongo.Database("account").Collection("user_email_index").InsertOne(sctx, &model.UserEmailIndex{
		Email:  email,
		UserID: user.UserID,
	})
	return
}
func (d *Dao) MongoGetUser(ctx context.Context, userid primitive.ObjectID) (*model.User, error) {
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
	index, e := d.MongoGetUserTelIndex(ctx, tel)
	if e != nil {
		return nil, e
	}
	//between find tel index and find user
	//another thread may change the association between user and tel
	user, e := d.MongoGetUser(ctx, index.UserID)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrDBConflict
		}
		return nil, e
	}
	if user.Tel != tel {
		return nil, ecode.ErrDBConflict
	}
	return user, nil
}
func (d *Dao) MongoGetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	index, e := d.MongoGetUserEmailIndex(ctx, email)
	if e != nil {
		return nil, e
	}
	//between find email index and find user
	//another thread may change the association between user and email
	user, e := d.MongoGetUser(ctx, index.UserID)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrDBConflict
		}
		return nil, e
	}
	if user.Email != email {
		return nil, ecode.ErrDBConflict
	}
	return user, nil
}
func (d *Dao) MongoGetUserByIDCard(ctx context.Context, idcard string) (*model.User, error) {
	index, e := d.MongoGetUserIDCardIndex(ctx, idcard)
	if e != nil {
		return nil, e
	}
	user, e := d.MongoGetUser(ctx, index.UserID)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrDBConflict
		}
		return nil, e
	}
	if user.IDCard != idcard {
		return nil, ecode.ErrDBConflict
	}
	return user, nil
}
func (d *Dao) MongoGetUserByNickName(ctx context.Context, nickname string) (*model.User, error) {
	index, e := d.MongoGetUserNickNameIndex(ctx, nickname)
	if e != nil {
		return nil, e
	}
	user, e := d.MongoGetUser(ctx, index.UserID)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrDBConflict
		}
		return nil, e
	}
	if user.NickName != nickname {
		return nil, ecode.ErrDBConflict
	}
	return user, nil
}
func (d *Dao) MongoUpdateUserTel(ctx context.Context, userid primitive.ObjectID, newTel string) (oldTel string, e error) {
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
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"tel": newTel}}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.Tel == newTel {
		return
	}
	if newTel != "" {
		if _, e = d.mongo.Database("account").Collection("user_tel_index").InsertOne(sctx, bson.M{"tel": newTel, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrTelAlreadyUsed
			}
			return
		}
	}
	oldTel = user.Tel
	if user.Tel != "" {
		if _, e = d.mongo.Database("account").Collection("user_tel_index").DeleteOne(sctx, bson.M{"tel": user.Tel, "user_id": userid}); e != nil {
			return
		}
	}
	user.Tel = newTel
	e = d._MongoDelUselessUser(sctx, user)
	return
}
func (d *Dao) MongoGetUserTelIndex(ctx context.Context, tel string) (*model.UserTelIndex, error) {
	index := &model.UserTelIndex{}
	if e := d.mongo.Database("account").Collection("user_tel_index").FindOne(ctx, bson.M{"tel": tel}).Decode(index); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return index, nil
}
func (d *Dao) MongoUpdateUserEmail(ctx context.Context, userid primitive.ObjectID, newEmail string) (oldEmail string, e error) {
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
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"email": newEmail}}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.Email == newEmail {
		return
	}
	if newEmail != "" {
		if _, e = d.mongo.Database("account").Collection("user_email_index").InsertOne(sctx, bson.M{"email": newEmail, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrEmailAlreadyUsed
			}
			return
		}
	}
	oldEmail = user.Email
	if user.Email != "" {
		if _, e = d.mongo.Database("account").Collection("user_email_index").DeleteOne(sctx, bson.M{"email": user.Email, "user_id": userid}); e != nil {
			return
		}
	}
	user.Email = newEmail
	e = d._MongoDelUselessUser(sctx, user)
	return
}
func (d *Dao) MongoGetUserEmailIndex(ctx context.Context, email string) (*model.UserEmailIndex, error) {
	index := &model.UserEmailIndex{}
	if e := d.mongo.Database("account").Collection("user_email_index").FindOne(ctx, bson.M{"email": email}).Decode(index); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return index, nil
}
func (d *Dao) MongoUpdateUserIDCard(ctx context.Context, userid primitive.ObjectID, newIDCard string) (oldIDCard string, e error) {
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
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"idcard": newIDCard}}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.IDCard == newIDCard {
		return
	}
	if newIDCard != "" {
		if _, e = d.mongo.Database("account").Collection("user_idcard_index").InsertOne(sctx, bson.M{"idcard": newIDCard, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrEmailAlreadyUsed
			}
			return
		}
	}
	oldIDCard = user.IDCard
	if user.IDCard != "" {
		if _, e = d.mongo.Database("account").Collection("user_idcard_index").DeleteOne(sctx, bson.M{"idcard": user.IDCard, "user_id": userid}); e != nil {
			return
		}
	}
	user.IDCard = newIDCard
	e = d._MongoDelUselessUser(sctx, user)
	return
}
func (d *Dao) MongoGetUserIDCardIndex(ctx context.Context, idcard string) (*model.UserIDCardIndex, error) {
	index := &model.UserIDCardIndex{}
	if e := d.mongo.Database("account").Collection("user_idcard_index").FindOne(ctx, bson.M{"idcard": idcard}).Decode(index); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return index, nil
}
func (d *Dao) MongoUpdateUserNickName(ctx context.Context, userid primitive.ObjectID, newNickName string) (oldNickName string, e error) {
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
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"nick_name": newNickName}}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if user.NickName == newNickName {
		return
	}
	if newNickName != "" {
		if _, e = d.mongo.Database("account").Collection("user_nick_name_index").InsertOne(sctx, bson.M{"nick_name": newNickName, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrNickNameAlreadyUsed
			}
			return
		}
	}
	oldNickName = user.NickName
	if user.NickName != "" {
		if _, e = d.mongo.Database("account").Collection("user_nick_name_index").DeleteOne(sctx, bson.M{"nick_name": user.NickName, "user_id": userid}); e != nil {
			return
		}
	}
	user.NickName = newNickName
	e = d._MongoDelUselessUser(sctx, user)
	return
}
func (d *Dao) MongoGetUserNickNameIndex(ctx context.Context, nickname string) (*model.UserNickNameIndex, error) {
	index := &model.UserNickNameIndex{}
	if e := d.mongo.Database("account").Collection("user_nick_name_index").FindOne(ctx, bson.M{"nick_name": nickname}).Decode(index); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return index, nil
}
func (d *Dao) MongoUpdateUserPassword(ctx context.Context, userid primitive.ObjectID, oldpassword, newpassword string) (e error) {
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

	nonce := make([]byte, 16)
	rand.Read(nonce)
	user := &model.User{}
	filter := bson.M{"_id": userid}
	updater := bson.M{"password": util.SignMake(newpassword, nonce)}
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, filter, bson.M{"$set": updater}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if e = util.SignCheck(oldpassword, user.Password); e != nil && e == ecode.ErrSignCheckFailed {
		e = ecode.ErrPasswordWrong
	}
	return
}
func (d *Dao) _MongoDelUselessUser(ctx context.Context, user *model.User) error {
	if user.NickName != "" || user.IDCard != "" || user.Email != "" || user.Tel != "" || len(user.OAuths) != 0 {
		return nil
	}
	_, e := d.mongo.Database("account").Collection("user").DeleteOne(ctx, bson.M{"_id": user.UserID})
	return e
}
