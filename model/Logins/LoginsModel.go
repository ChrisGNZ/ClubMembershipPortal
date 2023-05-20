package Logins

type LoginStruct struct {
	SessionID     string `json:"sessionID,omitempty"`
	Username      string `json:"username,omitempty"`
	Nickname      string `json:"nickname,omitempty"`
	Picture       string `json:"picture,omitempty"`
	UserId        string `json:"user_Id,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified string `json:"email_Verified,omitempty"`
	GivenName     string `json:"given_Name,omitempty"`
	FamilyName    string `json:"family_Name,omitempty"`
	ClientIP      string `json:"clientIP,omitempty"`
}
