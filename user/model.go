package user

import (
	"database/sql"
	"time"
)

type User struct {
	Id sql.NullInt64 `db:"user_id", json:"id"`
	Name string `db:"user_name", json:"name"`
	Password string `db:"user_pass", json:"password"`
	PasswordSalt string `db:"pass_salt", json:"passwordSalt"`
	PasswordType uint16 `db:"pass_type", json:"passwordType"`
	BadLogins sql.NullInt64 `db:"bad_logins", json:"badLogins"`
	Email sql.NullString `db:"email", json:"email"`
	NTCNick sql.NullString `db:"ntcnick", json:"ntcNick"`
	NickHistory sql.NullString `db:"nickhistory", json:"nickHistory"`
	Status *uint16 `db:"user_status", json:"status"`
	EmailUsed sql.NullInt64 `db:"email_used", json:"emailUsed"`
	ReferrerId sql.NullInt64 `db:"referrer", json:"referrerId"`
	GaduGadu sql.NullString `db:"gg", json:"gaduGadu"`
	ExtraInfo sql.NullString `db:"extrainfo", json:"extraInfo"`
	Trial *uint16 `db:"trial", json:"trial"`
	ShowEmail sql.NullString `db:"showemail", json:"showEmail"`
	NbOfRefs sql.NullInt64 `db:"refer_count", json:"nbOfRefs"`
	Suspended sql.NullString `db:"suspended", json:"suspended"`
	ChangeIp sql.NullString `db:"change_ip", json:"changeIp"`
	ChangeUserId sql.NullInt64 `db:"change_user_id", json:"changeUserId"`
	ChangeAt *time.Time `db:"change_date", json:"changeAt"`
	LoginAt *time.Time `db:"last_login", json:"loginAt"`
	CreatedAt *time.Time `db:"created", json:"createdAt"`
}

type Credentials struct {
	Email string `form:"email"`
	Password string `form:"password"`
}
