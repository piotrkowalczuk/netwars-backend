package user

import (
	"github.com/modcloth/sqlutil"
	"github.com/nu7hatch/gouuid"
	"database/sql"
	"time"
)

type SecureUser struct {
	Id  int64 `db:"user_id" json:"id"`
	Name string `db:"user_name" json:"name, string"`
	Email string `db:"email" json:"email, string"`
	NTCNick sqlutil.NullString `db:"ntcnick" json:"ntcNick, string"`
	NickHistory sqlutil.NullString `db:"nickhistory" json:"nickHistory, string"`
	Status *uint16 `db:"user_status" json:"status"`
	EmailUsed sqlutil.NullInt64 `db:"email_used" json:"emailUsed, string"`
	ReferrerId sqlutil.NullInt64 `db:"referrer" json:"referrerId"`
	GaduGadu sqlutil.NullString `db:"gg" json:"gaduGadu, string"`
	ExtraInfo sqlutil.NullString `db:"extrainfo" json:"extraInfo"`
	Trial *uint16 `db:"trial" json:"trial"`
	ShowEmail sqlutil.NullString `db:"showemail" json:"showEmail"`
	NbOfRefs sqlutil.NullInt64 `db:"refer_count" json:"nbOfRefs"`
	Suspended sqlutil.NullString `db:"suspended" json:"suspended"`
	ChangeIp sqlutil.NullString `db:"change_ip" json:"changeIp"`
	ChangeUserId sqlutil.NullInt64 `db:"change_user_id" json:"changeUserId"`
	ChangeAt *time.Time `db:"change_date" json:"changeAt"`
	LoginAt *time.Time `db:"last_login" json:"loginAt"`
	CreatedAt *time.Time `db:"created" json:"createdAt"`
}

type User struct {
	SecureUser
	Password string `db:"user_pass" json:"password"`
	PasswordSalt string `db:"pass_salt" json:"passwordSalt"`
	PasswordType uint16 `db:"pass_type" json:"passwordType"`
	BadLogins sqlutil.NullInt64 `db:"bad_logins" json:"badLogins"`
}

type UserRegistration struct {
	PlainPassword string `json:"plainPassword"`
	Name string `json:"name, string"`
	Email string `json:"email, string"`
	GaduGadu string `json:"gaduGadu, string"`
}

func (self *UserRegistration) isValid() (isValid bool) {
	isValid = true

	return
}

func (self *UserRegistration) createUser() *User {
	user := new(User)

	user.Name = self.Name
	user.Email = self.Email
	user.GaduGadu = sqlutil.NullString{sql.NullString{self.GaduGadu, true}}
	user.Password = self.PlainPassword

	return user
}

type UserSession struct {
	Id int64 `db:"user_id" json:"id"`
	Name string `db:"user_name" json:"name"`
	Email *string `db:"email" json:"email"`
	Trial *uint16 `db:"trial" json:"trial"`
	Token string `db:"-" json:"token"`
	Suspended *string `db:"suspended" json:"suspended"`
	ChangeAt *time.Time `db:"change_date" json:"changeAt"`
	LoginAt *time.Time `db:"last_login" json:"loginAt"`
	CreatedAt *time.Time `db:"created" json:"createdAt"`
}

func NewUserSession(user *User) (userSession *UserSession) {
	token, _ := uuid.NewV4()

	userSession = &UserSession{
		*&user.Id,
		user.Name,
		&user.Email,
		user.Trial,
		token.String(),
		&user.Suspended.String,
		user.ChangeAt,
		user.LoginAt,
		user.CreatedAt,
	}

	return
}

type APICredentials struct {
	Id int64 `form:"id" json:"id"`
	Token string `form:"token" json:"token"`
}

type LoginCredentials struct {
	Email string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}
