package money

import (
	"context"

	"github.com/chenjie199234/account/ecode"
	"github.com/chenjie199234/account/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// action: spend/recharge/refund/all
func (d *Dao) MongoGetMoneyLogs(ctx context.Context, userid primitive.ObjectID, opaction string) ([]*model.MoneyLog, error) {
	filter := bson.M{"user_id": userid}
	if opaction == "spend" || opaction == "recharge" || opaction == "refund" {
		filter["action"] = opaction
	}
	opts := options.Find().SetSort(bson.M{"_id": -1})
	cur, e := d.mongo.Database("account").Collection("money_log").Find(ctx, filter, opts)
	if e != nil {
		return nil, e
	}
	result := make([]*model.MoneyLog, 0, cur.RemainingBatchLength())
	e = cur.All(ctx, &result)
	return result, e
}
func (d *Dao) MongoInsertMoneyLogs(ctx context.Context, log *model.MoneyLog) error {
	r, e := d.mongo.Database("account").Collection("money_log").InsertOne(ctx, log)
	if e == nil {
		log.LogID = r.InsertedID.(primitive.ObjectID)
		return nil
	}
	if !mongo.IsDuplicateKeyError(e) {
		return e
	}
	dblog := &model.MoneyLog{}
	if e := d.mongo.Database("account").Collection("money_log").FindOne(ctx, bson.M{"user_id": log.UserID, "action": log.Action, "unique_id": log.UniqueID}).Decode(dblog); e != nil {
		return e
	}
	if dblog.SrcDst != log.SrcDst || dblog.MoneyType != log.MoneyType || dblog.MoneyAmount != log.MoneyAmount {
		e = ecode.ErrDBDataConflict
	} else {
		log.LogID = dblog.LogID
	}
	return e
}
