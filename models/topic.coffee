db = require '../db'

module.exports =
  getById: (id, next) ->
    getById(id, next)
  getBy: (fieldName, value, next) ->
    getBy(fieldName, value, next)


getBy = (fieldName, value, next) ->
  db.postgres()
    .from('forum_topic')
    .select(
      'forum_id',
      'topic_id',
      'first_poster',
      'first_poster_name',
      'last_poster',
      'last_poster_name',
      'last_post_id',
      'last_post_date',
      'topic_posts',
      'topic_views',
      'topic_closed',
      'topic_pined',
      'topic_visible_from',
      'topic_visible_to',
      'topic_deleted',
      'change_date',
      'change_user_id',
      'change_ip'
    )
    .where(fieldName, value)
    .exec (error, reply) ->
      if not error
        next(reply[0])
      else
        next(false)

getById = (id, next) ->
  getBy('topic_id', id, next)