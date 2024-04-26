package util

import (
	"context"

	"github.com/chenjie199234/account/config"

	"github.com/chenjie199234/Corelib/util/common"
)

func SendEmailCode(ctx context.Context, email, code, action string) error {
	body := "Your dynamic code is:" + code
	return config.GetEmail("qq_email").SendTextEmail(ctx, []string{email}, "Login Dynamic Code", common.STB(body))
}

// TODO
func SendTelCode(ctx context.Context, tel, code, action string) error {
	return nil
}
