db = require '../db'

module.exports =
  getById: (id, next) ->
    getById(id, next)
  getBy: (fieldName, value, next) ->
    getBy(fieldName, value, next)


getBy = (fieldName, value, next) ->
  db.postgres()
    .from('forum')
    .select(
      'forum_id',
      'forum_name',
      'forum_desc',
      'forum_order',
      'forum_type',
      'forum_topics',
      'show_topics'
    )
    .where(fieldName, value)
    .exec (error, reply) ->
      if not error
        next(reply[0])
      else
        next(false)

getById = (id, next) ->
  getBy('forum_id', id, next)
