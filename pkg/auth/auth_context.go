package auth_context

import "github.com/golang-jwt/jwt/v5"

type AuthContext struct {
	IdSite     string   `json:"id_site"`
	IdAuthUser string   `json:"id_auth_user"`
	Roles      []string `json:"roles"`
	UserAgent  string   `json:"user_agent"`
	UserIP     string   `json:"user_ip"`
	UserName   string   `json:"user_name"`
	Email      string   `json:"email"`
	jwt.RegisteredClaims
}
