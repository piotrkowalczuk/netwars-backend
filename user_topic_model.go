package main

type UserTopic struct {
	UserId   *int64 `json:"userId"`
	TopicId  *int64 `json:"topicId"`
	Observed *bool  `json:"observed"`
	PostSeen *int   `json:"postSeen"`
}

func NewUserTopic(userId int64, topicId int64, postSeen int) *UserTopic {
	userTopic := new(UserTopic)
	userTopic.UserId = &userId
	userTopic.TopicId = &topicId
	userTopic.PostSeen = &postSeen

	return userTopic
}
