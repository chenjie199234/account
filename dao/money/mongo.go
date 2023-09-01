package money

import (
	"context"

	"github.com/chenjie199234/account/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// action: spend/recharge/refund/all
// page: 0:means return all,>0:means return the required page,if page overflow,the last page will return
func (d *Dao) MongoGetMoneyLogs(ctx context.Context, userid primitive.ObjectID, action string, pagesize, page uint32) ([]*model.MoneyLog, uint32, uint32, error) {
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
	var cur *mongo.Cursor
	if page == 0 {
		cur, e = d.mongo.Database("account").Collection("money_log").Find(ctx, filter, opts)
	} else {
		skip := int64((page - 1) * pagesize)
		if skip >= totalsize {
			if totalsize%int64(pagesize) > 0 {
				page = uint32(totalsize/int64(pagesize) + 1)
			} else {
				page = uint32(totalsize / int64(pagesize))
			}
			skip = int64((page - 1) * pagesize)
		}
		opts = opts.SetSkip(skip).SetLimit(int64(pagesize))
		cur, e = d.mongo.Database("account").Collection("money_log").Find(ctx, filter, opts)
	}
	result := make([]*model.MoneyLog, 0, cur.RemainingBatchLength())
	e = cur.All(ctx, &result)
	return result, page, uint32(totalsize), e
}
