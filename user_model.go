package main

import (
	"github.com/jamieomatthews/validation"
	"github.com/martini-contrib/binding"
	"github.com/modcloth/sqlutil"
	"github.com/nu7hatch/gouuid"
	"net/http"
	"time"
)

const (
	SESSION_LIFE_TIME int = 86400
)

type SecureUser struct {
	Id           int64              `db:"user_id" json:"id"`
	Name         string             `db:"user_name" json:"name, string"`
	Email        string             `db:"email" json:"email, string"`
	NTCNick      sqlutil.NullString `db:"ntcnick" json:"ntcNick, string"`
	NickHistory  sqlutil.NullString `db:"nickhistory" json:"nickHistory, string"`
	Status       *uint16            `db:"user_status" json:"status"`
	EmailUsed    sqlutil.NullInt64  `db:"email_used" json:"emailUsed, string"`
	ReferrerId   sqlutil.NullInt64  `db:"referrer" json:"referrerId"`
	GaduGadu     sqlutil.NullString `db:"gg" json:"gaduGadu, string"`
	ExtraInfo    sqlutil.NullString `db:"extrainfo" json:"extraInfo"`
	Trial        *uint16            `db:"trial" json:"trial"`
	ShowEmail    sqlutil.NullString `db:"showemail" json:"showEmail"`
	NbOfRefs     sqlutil.NullInt64  `db:"refer_count" json:"nbOfRefs"`
	Suspended    sqlutil.NullString `db:"suspended" json:"suspended"`
	ChangeIp     sqlutil.NullString `db:"change_ip" json:"changeIp"`
	ChangeUserId sqlutil.NullInt64  `db:"change_user_id" json:"changeUserId"`
	ChangeAt     *time.Time         `db:"change_date" json:"changeAt"`
	LoginAt      *time.Time         `db:"last_login" json:"loginAt"`
	CreatedAt    *time.Time         `db:"created" json:"createdAt"`
}

type User struct {
	SecureUser
	Password     string            `db:"user_pass" json:"password"`
	PasswordSalt string            `db:"pass_salt" json:"passwordSalt"`
	PasswordType uint16            `db:"pass_type" json:"passwordType"`
	BadLogins    sqlutil.NullInt64 `db:"bad_logins" json:"badLogins"`
}

func (u *User) GetSecureUser() (su *SecureUser) {
	su.Id = u.Id
	su.Name = u.Name
	su.Email = u.Email
	su.NTCNick = u.NTCNick
	su.NickHistory = u.NickHistory
	su.Status = u.Status
	su.EmailUsed = u.EmailUsed
	su.ReferrerId = u.ReferrerId
	su.GaduGadu = u.GaduGadu
	su.ExtraInfo = u.ExtraInfo
	su.Trial = u.Trial
	su.ShowEmail = u.ShowEmail
	su.NbOfRefs = u.NbOfRefs
	su.Suspended = u.Suspended
	su.ChangeIp = u.ChangeIp
	su.ChangeUserId = u.ChangeUserId
	su.ChangeAt = u.ChangeAt
	su.LoginAt = u.LoginAt
	su.CreatedAt = u.CreatedAt

	return
}

type UserRegistration struct {
	PlainPassword string `json:"plainPassword"`
	Name          string `json:"name, string"`
	Email         string `json:"email, string"`
}

func (ur UserRegistration) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	v := validation.NewValidation(&errors, ur)

	v.Validate(&ur.Name).
		Message("Nazwa użytkownika może mieć od 3 do 20 znaków.").
		Key("name").
		TrimSpace().
		Range(3, 20)
	v.Validate(&ur.Email).
		Message("Niepoprawny adres email.").
		Key("email").
		TrimSpace().
		Email()
	v.Validate(&ur.PlainPassword).
		Message("Hasło może mieć od 6 do 20 znaków.").
		Key("plainPassword").
		TrimSpace().
		Range(6, 20)

	return *v.Errors.(*binding.Errors)
}

func (self *UserRegistration) createUser() *User {
	user := new(User)

	user.Name = self.Name
	user.Email = self.Email
	user.Password = self.PlainPassword

	return user
}

type BasicUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type UserSession struct {
	Id        int64      `db:"user_id" json:"id"`
	Name      string     `db:"user_name" json:"name"`
	Email     *string    `db:"email" json:"email"`
	Trial     *uint16    `db:"trial" json:"trial"`
	Token     string     `db:"-" json:"token"`
	Suspended *string    `db:"suspended" json:"suspended"`
	ChangeAt  *time.Time `db:"change_date" json:"changeAt"`
	LoginAt   *time.Time `db:"last_login" json:"loginAt"`
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

func (us *UserSession) getSessionKey() string {
	return createUserSessionKey(us.Token)
}

func createUserSessionKey(key string) string {
	return "user:" + key
}

type APICredentials struct {
	Id    int64  `form:"id" json:"id"`
	Token string `form:"token" json:"token"`
}

func (apic *APICredentials) getSessionKey() string {
	return createUserSessionKey(apic.Token)
}

type LoginCredentials struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}
