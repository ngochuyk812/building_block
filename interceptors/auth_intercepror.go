package interceptors

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ngochuyk812/building_block/infrastructure/helpers"
	auth_context "github.com/ngochuyk812/building_block/pkg/auth"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
)

func NewAuthInterceptor(secret string, policies *map[string][]string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (res connect.AnyResponse, err error) {
			tokenStr := req.Header().Get("Authorization")
			path := req.Spec().Procedure
			allowedRoles, _ := (*policies)[path]

			if len(allowedRoles) > 0 && tokenStr == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthorized: missing token"))

			}
			if tokenStr != "" {
				tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
				token, err := jwt.ParseWithClaims(tokenStr, &auth_context.AuthContext{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(secret), nil
				})

				if err != nil || !token.Valid {
					return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthorized: invalid token"))
				}
				claims, err := auth_context.VerifyJWT(tokenStr, secret)
				if err == nil {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf(err.Error()))
				}

				userRoles := claims.Roles
				valid := hasValidRole(userRoles, allowedRoles)

				if valid == false {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("unauthorized access"))
				}
				ctx = helpers.NewContext(ctx, helpers.AuthContextKey, &auth_context.AuthContext{
					IdSite:     claims.IdSite,
					IdAuthUser: claims.IdAuthUser,
					Roles:      claims.Roles,
					UserName:   claims.UserName,
					Email:      claims.Email,
					UserIP:     req.Header().Get("X-Forwarded-For"),
				})
			}

			response, errService := next(ctx, req)
			return response, errService
		}
	}
}

func hasValidRole(userRoles, allowedRoles []string) bool {
	for _, ur := range userRoles {
		for _, ar := range allowedRoles {
			if ur == ar {
				return true
			}
		}
	}
	return false
}
