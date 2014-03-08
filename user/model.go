package user

import (
	"time"
)

type User struct {
	Id uint32 `db:"user_id"`
	Name string `db:"user_name"`
	Password string `db:"user_pass"`
	PasswordSalt string `db:"pass_salt"`
	PasswordType uint16 `db:"pass_type"`
	BadLogins uint32 `db:"bad_logins"`
	Email string `db:"email"`
	NTCNick string `db:"ntcnick"`
	NickHistory string `db:"nickhistory"`
	Status uint16 `db:"user_status"`
	EmailUsed uint16 `db:"email_used"`
	ReferrerId uint32 `db:"referrer"`
	GaduGadu string `db:"gg"`
	ExtraInfo string `db:"extrainfo"`
	Trial uint16 `db:"trial"`
	ShowEmail string `db:"showemail"`
	NbOfRefs uint32 `db:"refer_count"`
	Suspended string `db:"suspended"`
	ChangeIp string `db:"change_ip"`
	ChangeUserId uint32 `db:"change_user_id"`
	ChangeAt time.Time `db:"change_date"`
	LoginAt time.Time `db:"last_login"`
	CreatedAt time.Time `db:"created"`
}
