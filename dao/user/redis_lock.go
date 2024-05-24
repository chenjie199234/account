package user

import (
	"context"
	"time"

	"github.com/chenjie199234/account/ecode"
)

// 3 times per 1 min
func (d *Dao) RedisLockLoginDynamic(ctx context.Context, src string) error {
	rate := map[string][2]uint64{"rate_dynamic_login_{" + src + "}": {3, 60}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockTelOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"rate_tel_op_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockEmailOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"rate_email_op_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockOAuthOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"rate_oauth_op_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockUpdatePassword(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"rate_update_password_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockResetPassword(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"rate_reset_password_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 1 time per second
func (d *Dao) RedisLockDuplicateCheck(ctx context.Context, srctype, userid string) error {
	ok, e := d.redis.SetNX(ctx, "rate_"+srctype+"_duplicate_check_{"+userid+"}", 1, time.Second).Result()
	if e == nil && !ok {
		e = ecode.ErrTooFast
	}
	return e
}
