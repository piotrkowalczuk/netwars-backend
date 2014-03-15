package forum

import (
	"github.com/modcloth/sqlutil"
)

type Forum struct {
	Id sqlutil.NullInt64 `db:"forum_id" json:"id"`
	Name sqlutil.NullString `db:"forum_name" json:"name"`
	Description sqlutil.NullString `db:"forum_desc" json:"description"`
	Order sqlutil.NullInt64 `db:"forum_order" json:"order"`
	Type sqlutil.NullInt64 `db:"forum_type" json:"type"`
	Topics sqlutil.NullInt64 `db:"forum_topics" json:"topics"`
	ShowTopics sqlutil.NullInt64 `db:"show_topics" json:"showTopics"`
}


