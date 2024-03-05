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
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	nonce := make([]byte, 16)
	rand.Read(nonce)
	user = &model.User{
		Password: util.SignMake("", nonce),
		IDCard:   "",
		Tel:      tel,
		Email:    "",
		OAuths:   map[string]string{},
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
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	nonce := make([]byte, 16)
	rand.Read(nonce)
	user = &model.User{
		Password: util.SignMake("", nonce),
		IDCard:   "",
		Tel:      "",
		Email:    email,
		OAuths:   map[string]string{},
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
func (d *Dao) MongoCreateUserByOAuth(ctx context.Context, oauthservicename, oauthid string) (user *model.User, e error) {
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
				user, e = d.MongoGetUserByOAuth(ctx, oauthservicename, oauthid)
			}
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	nonce := make([]byte, 16)
	rand.Read(nonce)
	user = &model.User{
		Password: util.SignMake("", nonce),
		IDCard:   "",
		Tel:      "",
		Email:    "",
		OAuths:   map[string]string{oauthservicename: oauthid},
		Money:    map[string]int32{},
	}
	var r *mongo.InsertOneResult
	if r, e = d.mongo.Database("account").Collection("user").InsertOne(sctx, user); e != nil {
		return
	}
	user.UserID = r.InsertedID.(primitive.ObjectID)
	_, e = d.mongo.Database("account").Collection("user_oauth_index").InsertOne(sctx, &model.UserOAuthIndex{
		Service: oauthservicename + "|" + oauthid,
		UserID:  user.UserID,
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
			e = ecode.ErrDBDataConflict
		}
		return nil, e
	}
	if user.Tel != tel {
		return nil, ecode.ErrDBDataConflict
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
			e = ecode.ErrDBDataConflict
		}
		return nil, e
	}
	if user.Email != email {
		return nil, ecode.ErrDBDataConflict
	}
	return user, nil
}
func (d *Dao) MongoGetUserByIDCard(ctx context.Context, idcard string) (*model.User, error) {
	index, e := d.MongoGetUserIDCardIndex(ctx, idcard)
	if e != nil {
		return nil, e
	}
	//between find email index and find user
	//another thread may change the association between user and idcard
	user, e := d.MongoGetUser(ctx, index.UserID)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrDBDataConflict
		}
		return nil, e
	}
	if user.IDCard != idcard {
		return nil, ecode.ErrDBDataConflict
	}
	return user, nil
}
func (d *Dao) MongoGetUserByOAuth(ctx context.Context, oauthservicename, oauthid string) (*model.User, error) {
	index, e := d.MongoGetUserOAuthIndex(ctx, oauthservicename, oauthid)
	if e != nil {
		return nil, e
	}
	//between find email index and find user
	//another thread may change the association between user and oauth
	user, e := d.MongoGetUser(ctx, index.UserID)
	if e != nil {
		if e == ecode.ErrUserNotExist {
			e = ecode.ErrDBDataConflict
		}
		return nil, e
	}
	if user.OAuths[oauthservicename] != oauthid {
		return nil, ecode.ErrDBDataConflict
	}
	return user, nil
}
func (d *Dao) MongoUpdateUserOAuth(ctx context.Context, userid primitive.ObjectID, oauthservicename, newoauthid string) (olduser *model.User, e error) {
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
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	updater := bson.M{}
	if newoauthid == "" {
		updater = bson.M{"$unset": bson.M{"oauths." + oauthservicename: 1}}
	} else {
		updater = bson.M{"$set": bson.M{"oauths." + oauthservicename: newoauthid}}
	}
	olduser = &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, updater).Decode(olduser); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if olduser.OAuths[oauthservicename] == newoauthid {
		return
	}
	oldOauthID := olduser.OAuths[oauthservicename]
	if newoauthid != "" {
		if _, e = d.mongo.Database("account").Collection("user_oauth_index").InsertOne(sctx, bson.M{"service": oauthservicename + "|" + newoauthid, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrTelAlreadyUsed
			}
			return
		}
	}
	if oldOauthID != "" {
		if _, e = d.mongo.Database("account").Collection("user_oauth_index").DeleteOne(sctx, bson.M{"service": oauthservicename + "|" + oldOauthID, "user_id": userid}); e != nil {
			return
		}
	}
	if newoauthid == "" {
		delete(olduser.OAuths, oauthservicename)
	} else {
		olduser.OAuths[oauthservicename] = newoauthid
	}
	//now olduser is new,check if we need to delete it
	e = d._MongoDelUselessUser(sctx, olduser)
	//turn back to olduser
	if oldOauthID == "" {
		delete(olduser.OAuths, oauthservicename)
	} else {
		olduser.OAuths[oauthservicename] = oldOauthID
	}
	return
}
func (d *Dao) MongoGetUserOAuthIndex(ctx context.Context, oauthservicename, oauthid string) (*model.UserOAuthIndex, error) {
	index := &model.UserOAuthIndex{}
	e := d.mongo.Database("account").Collection("user_oauth_index").FindOne(ctx, bson.M{"service": oauthservicename + "|" + oauthid}).Decode(index)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return nil, e
	}
	return index, nil
}
func (d *Dao) MongoUpdateUserTel(ctx context.Context, userid primitive.ObjectID, newTel string) (olduser *model.User, e error) {
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
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	olduser = &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"tel": newTel}}).Decode(olduser); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if olduser.Tel == newTel {
		return
	}
	oldTel := olduser.Tel
	if newTel != "" {
		if _, e = d.mongo.Database("account").Collection("user_tel_index").InsertOne(sctx, bson.M{"tel": newTel, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrTelAlreadyUsed
			}
			return
		}
	}
	if oldTel != "" {
		if _, e = d.mongo.Database("account").Collection("user_tel_index").DeleteOne(sctx, bson.M{"tel": oldTel, "user_id": userid}); e != nil {
			return
		}
	}
	olduser.Tel = newTel
	//now olduser is new,check if we need to delete it
	e = d._MongoDelUselessUser(sctx, olduser)
	//turn back to olduser
	olduser.Tel = oldTel
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
func (d *Dao) MongoUpdateUserEmail(ctx context.Context, userid primitive.ObjectID, newEmail string) (olduser *model.User, e error) {
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
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	olduser = &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"email": newEmail}}).Decode(olduser); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if olduser.Email == newEmail {
		return
	}
	oldEmail := olduser.Email
	if newEmail != "" {
		if _, e = d.mongo.Database("account").Collection("user_email_index").InsertOne(sctx, bson.M{"email": newEmail, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrEmailAlreadyUsed
			}
			return
		}
	}
	if oldEmail != "" {
		if _, e = d.mongo.Database("account").Collection("user_email_index").DeleteOne(sctx, bson.M{"email": oldEmail, "user_id": userid}); e != nil {
			return
		}
	}
	olduser.Email = newEmail
	//now olduser is new,check if we need to delete it
	e = d._MongoDelUselessUser(sctx, olduser)
	//turn back to olduser
	olduser.Email = oldEmail
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
func (d *Dao) MongoUpdateUserIDCard(ctx context.Context, userid primitive.ObjectID, newIDCard string) (olduser *model.User, e error) {
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
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	olduser = &model.User{}
	if e = d.mongo.Database("account").Collection("user").FindOneAndUpdate(sctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"idcard": newIDCard}}).Decode(olduser); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	if olduser.IDCard == newIDCard {
		return
	}
	oldIDCard := olduser.IDCard
	if newIDCard != "" {
		if _, e = d.mongo.Database("account").Collection("user_idcard_index").InsertOne(sctx, bson.M{"idcard": newIDCard, "user_id": userid}); e != nil {
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrEmailAlreadyUsed
			}
			return
		}
	}
	if oldIDCard != "" {
		if _, e = d.mongo.Database("account").Collection("user_idcard_index").DeleteOne(sctx, bson.M{"idcard": oldIDCard, "user_id": userid}); e != nil {
			return
		}
	}
	olduser.IDCard = newIDCard
	//now olduser is new,check if we need to delete it
	e = d._MongoDelUselessUser(sctx, olduser)
	//turn back to olduser
	olduser.IDCard = oldIDCard
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
		} else if e = s.CommitTransaction(sctx); e != nil {
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
	if user.IDCard != "" || user.Email != "" || user.Tel != "" || len(user.OAuths) != 0 {
		return nil
	}
	_, e := d.mongo.Database("account").Collection("user").DeleteOne(ctx, bson.M{"_id": user.UserID})
	return e
}
