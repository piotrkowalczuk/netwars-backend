package forum

import (
	"github.com/modcloth/sqlutil"
	"time"
)

type Forum struct {
	Id int64 `db:"forum_id" json:"id"`
	Name sqlutil.NullString `db:"forum_name" json:"name"`
	Description sqlutil.NullString `db:"forum_desc" json:"description"`
	Order sqlutil.NullInt64 `db:"forum_order" json:"order"`
	Type sqlutil.NullInt64 `db:"forum_type" json:"type"`
	Topics sqlutil.NullInt64 `db:"forum_topics" json:"topics"`
	ShowTopics sqlutil.NullInt64 `db:"show_topics" json:"showTopics"`
}

type Topic struct {
	Id int64 `db:"topic_id" json:"id"`
	ForumId int64 `db:"forum_id" json:"forumId"`
	Name *string `db:"topic_name" json:"name"`
	AuthorId sqlutil.NullInt64 `db:"first_poster" json:"authorId"`
	AuthorName sqlutil.NullString `db:"first_poster_name" json:"authorName"`
	LastPostAuthorId sqlutil.NullInt64 `db:"last_poster" json:"lastPostAuthorId"`
	LastPostAuthorName sqlutil.NullString `db:"last_poster_name" json:"lastPostAuthorName"`
	LastPostId sqlutil.NullInt64 `db:"last_post_id" json:"lastPostId"`
	LastPostDate *time.Time `db:"last_post_date" json:"lastPostDate"`
	NbOfPosts sqlutil.NullInt64 `db:"topic_posts" json:"nbOfPosts"`
	NbOfViews sqlutil.NullInt64 `db:"topic_views" json:"nbOfViews"`
	IsClosed *int16 `db:"topic_closed" json:"isClosed"`
	IsPinned *int16 `db:"topic_pined" json:"isPinned"`
	IsDeleted *int16 `db:"topic_deleted" json:"isDeleted"`
	VisibleFrom *time.Time `db:"topic_visible_from" json:"visibleFrom"`
	VisibleTo *time.Time `db:"topic_visible_to" json:"visibleTo"`
	ChangeAt *time.Time `db:"change_date" json:"changeAt"`
	ChangerId sqlutil.NullInt64 `db:"change_user_id" json:"changerId"`
	ChangerIP sqlutil.NullString `db:"change_ip" json:"changerIP"`
}

type Post struct {
	Id int64 `db:"post_id" json:"id"`
	TopicId int64 `db:"topic_id" json:"topicId"`
	AuthorId sqlutil.NullInt64 `db:"user_id" json:"authorId"`
	AuthorName sqlutil.NullString `db:"user_name" json:"authorName"`
	AuthorIP sqlutil.NullString `db:"ip_address" json:"authorIP"`
	CreatedAt *time.Time `db:"post_date" json:"createdAt"`
	Content *string `db:"post_body" json:"content"`
	NbOfChanges sqlutil.NullInt64 `db:"mod_counter" json:"nbOfChanges"`
	ChangeAt *time.Time `db:"mod_date" json:"changeAt"`
	ChangerId sqlutil.NullInt64 `db:"mod_user_id" json:"changerId"`
	ChangerName sqlutil.NullString `db:"mod_user_name" json:"changerName"`
}

type CreateTopicRequest struct {
	Post Post `json:"post"`
	Topic Topic `json:"topic"`
}
