package user

import (
	"context"
	"time"

	"github.com/chenjie199234/account/ecode"

	credis "github.com/chenjie199234/Corelib/redis"
	"github.com/chenjie199234/Corelib/util/ctime"
)

const CodeExpire = ctime.Duration(time.Minute * 5)

func (d *Dao) RedisSetCode(ctx context.Context, target, action, receiver string) (string, bool, error) {
	code, duplicatereciver, e := d.redis.MakeVerifyCode(ctx, target, action, receiver, uint(CodeExpire.StdDuration().Seconds()))
	if e != nil && e == credis.ErrVerifyCodeCheckTimesUsedup {
		e = ecode.ErrBan
	}
	return code, duplicatereciver, e
}
func (d *Dao) RedisCheckCode(ctx context.Context, target, action, code, mustreceiver string) error {
	e := d.redis.CheckVerifyCode(ctx, target, action, code, mustreceiver)
	if e != nil {
		if e == credis.ErrVerifyCodeMissing {
			e = ecode.ErrCodeNotExist
		} else if e == credis.ErrVerifyCodeReceiverMissing {
			e = ecode.ErrPasswordWrong
		} else if e == credis.ErrVerifyCodeWrong {
			e = ecode.ErrPasswordWrong
		} else if e == credis.ErrVerifyCodeCheckTimesUsedup {
			e = ecode.ErrBan
		}
	}
	return e
}
func (d *Dao) RedisDelCode(ctx context.Context, target, action string) error {
	return d.redis.DelVerifyCode(ctx, target, action)
}
func (d *Dao) RedisCodeCheckTimes(ctx context.Context, target, action, mustreceiver string) error {
	e := d.redis.HasCheckTimes(ctx, target, action, mustreceiver)
	if e != nil {
		if e == credis.ErrVerifyCodeMissing {
			e = ecode.ErrCodeNotExist
		} else if e == credis.ErrVerifyCodeReceiverMissing {
			e = ecode.ErrCodeNotExist
		} else if e == credis.ErrVerifyCodeCheckTimesUsedup {
			e = ecode.ErrBan
		}
	}
	return e
}
