db = require '../db'

module.exports =
  getById: (id, next) ->
    getById(id, next)
  getBy: (fieldName, value, next) ->
    getBy(fieldName, value, next)


getBy = (fieldName, value, next) ->
  db.postgres()
    .from('users')
    .select(
      'user_id',
      'user_name',
      'last_login',
      'bad_logins',
      'email',
      'ntcnick',
      'nickhistory',
      'user_status',
      'change_date',
      'change_user_id',
      'change_ip',
      'email_used',
      'referrer',
      'gg',
      'extrainfo',
      'created',
      'trial',
      'showemail',
      'refer_count',
      'suspended'
    )
    .where(fieldName, value)
    .exec (error, reply) ->
      if not error
        next(reply[0])
      else
        next(false)

getById = (id, next) ->
  getBy('user_id', id, next)
