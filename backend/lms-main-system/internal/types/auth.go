package types

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
