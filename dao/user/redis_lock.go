package user

import (
	"context"
	"time"

	"github.com/chenjie199234/account/ecode"
)

//send email or send tel rate limit

// 3 times per 1 min
func (d *Dao) RedisLockLoginTelDynamic(ctx context.Context, tel string) error {
	rate := map[string][2]uint64{"dynamic_tel_login_lock_{" + tel + "}": {3, 60}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 3 times per 1 min
func (d *Dao) RedisLockLoginEmailDynamic(ctx context.Context, email string) error {
	rate := map[string][2]uint64{"dynamic_email_login_lock_{" + email + "}": {3, 60}}
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
func (d *Dao) RedisLockNickNameOP(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"nickname_op_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

//other rate limit

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
