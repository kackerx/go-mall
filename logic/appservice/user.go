package appservice

import (
	"context"

	"github.com/kackerx/go-mall/api/reply"
	"github.com/kackerx/go-mall/common/log"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/logic/domainservice"
)

type UserAppSvc struct {
	userDomainSvc *domainservice.UserDomainSvc
}

func NewUserAppSvc(userDomainSvc *domainservice.UserDomainSvc) *UserAppSvc {
	return &UserAppSvc{userDomainSvc: userDomainSvc}
}
func (us *UserAppSvc) GenToken(ctx context.Context) (resp *reply.TokenResp, err error) {
	tokenInfo, err := us.userDomainSvc.GenAuthToken(ctx, 110, "h5", "")
	if err != nil {
		return nil, err
	}

	log.New(ctx).Info("gen token success", tokenInfo)

	resp = new(reply.TokenResp)
	if err = util.Copy(resp, tokenInfo); err != nil {
		return
	}

	return
}
