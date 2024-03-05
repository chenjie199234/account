package user

import (
	"context"
	"time"

	"github.com/chenjie199234/account/ecode"
)

// 3 times per 1 min
func (d *Dao) RedisLockLoginDynamic(ctx context.Context, src string) error {
	rate := map[string][2]uint64{"dynamic_login_lock_{" + src + "}": {3, 60}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockTelOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"tel_op_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockEmailOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"email_op_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockIDCardOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"idcard_op_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockOAuthOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"oauth_op_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockUpdatePassword(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"update_password_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 1 times per second
func (d *Dao) RedisLockDuplicateCheck(ctx context.Context, srctype, userid string) error {
	ok, e := d.redis.SetNX(ctx, srctype+"_duplicate_check_lock_{"+userid+"}", 1, time.Second).Result()
	if e == nil && !ok {
		e = ecode.ErrTooFast
	}
	return e
}
