package router

import (
	"github.com/ktpswjz/httpserver/example/webserver/server/errors"
	"github.com/ktpswjz/httpserver/types"
	"github.com/ktpswjz/om/data/model"
)

func (s *innerRouter) checkOmwToken(ip, token string) (types.Error, error) {
	if len(token) < 1 {
		return errors.AuthNoToken, nil
	}

	entity, ok := s.omwToken.Get(token, true)
	if !ok {
		return errors.AuthTokenInvalid, nil
	}
	tokenEntity := entity.(*model.Token)
	if tokenEntity == nil {
		return errors.AuthTokenInvalid, nil
	}

	if ip != tokenEntity.LoginIP {
		return errors.AuthTokenIllegal, nil
	}

	return nil, nil
}
