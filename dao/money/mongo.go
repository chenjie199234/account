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
// page: 0:means return all,>0:means return the required page,if page overflow,the last page will return
func (d *Dao) MongoGetMoneyLogs(ctx context.Context, userid primitive.ObjectID, action string, page int64) ([]*model.MoneyLog, int64, int64, error) {
	filter := bson.M{"user_id": userid}
	if action == "spend" || action == "recharge" || action == "refund" {
		filter["action"] = action
	}
	totalsize, e := d.mongo.Database("account").Collection("money_log").CountDocuments(ctx, filter)
	if e != nil {
		return nil, 0, 0, e
	}
	if totalsize == 0 {
		return make([]*model.MoneyLog, 0), 0, 0, nil
	}
	opts := options.Find().SetSort(bson.M{"_id": -1})
	if page != 0 {
		skip := (page - 1) * DefaultMoneyLogsPageSize
		if skip >= totalsize {
			if totalsize%DefaultMoneyLogsPageSize > 0 {
				page = totalsize/DefaultMoneyLogsPageSize + 1
			} else {
				page = totalsize / DefaultMoneyLogsPageSize
			}
			skip = (page - 1) * DefaultMoneyLogsPageSize
		}
		opts = opts.SetSkip(skip).SetLimit(DefaultMoneyLogsPageSize)
	}
	cur, e := d.mongo.Database("account").Collection("money_log").Find(ctx, filter, opts)
	if e != nil {
		return nil, 0, 0, e
	}
	result := make([]*model.MoneyLog, 0, cur.RemainingBatchLength())
	e = cur.All(ctx, &result)
	return result, page, totalsize, e
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
		e = ecode.ErrDBConflict
	} else {
		log.LogID = dblog.LogID
	}
	return e
}
