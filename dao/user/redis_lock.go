package user

import (
	"context"

	"github.com/chenjie199234/account/ecode"
)

// 5 times per 10 min
func (d *Dao) RedisLockLoginTelDynamic(ctx context.Context, tel string) error {
	rate := map[string][2]uint64{"dynamic_tel_login_lock_{" + tel + "}": {5, 600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per 10 min
func (d *Dao) RedisLockLoginEmailDynamic(ctx context.Context, email string) error {
	rate := map[string][2]uint64{"dynamic_email_login_lock_{" + email + "}": {5, 600}}
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

// 5 times per hour
func (d *Dao) RedisLockUpdateNickName(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"update_nickname_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockUpdateTel(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"update_tel_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}

// 5 times per hour
func (d *Dao) RedisLockUpdateEmail(ctx context.Context, userid string) error {
	rate := map[string][2]uint64{"update_email_lock_{" + userid + "}": {5, 3600}}
	success, e := d.redis.RateLimit(ctx, rate)
	if e == nil && !success {
		e = ecode.ErrTooFast
	}
	return e
}
