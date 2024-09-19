package auth

type HasAuthDto struct {
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserRole  string `json:"user_role"`
	Location  string `json:"location"`
}

func (d *HasAuthDto) AssignAuthData(claims JwtClaims) {
	d.UserId = claims.UserId
	d.UserName = claims.UserName
	d.UserEmail = claims.UserEmail
	d.UserRole = claims.UserRole

	// d.AuthUserId = claims.OldUserId
	// d.AuthUserType = claims.UserType
	// d.AuthUserIsAMaster = claims.OldUserIsAMaster
	// d.AuthUserName = claims.UserName
	// d.AuthUserEmail = claims.UserEmail
	// d.AuthUserIsADoa = claims.UserIsADoa

	// @TODO: Wait FE implement new login
	//d.AuthUserId = claims.UserId
	//d.AuthUserIsAMaster = claims.UserIsAMaster
}
