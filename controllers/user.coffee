config = require '../config'
user = require '../models/user'

module.exports = (app) ->
  app.get '/user/:id', getById

getById = (req, res) ->

  userId = req.params.id
  user.getById userId, (user) =>
    if user
      res.header "Content-Type", "application/json"
      res.send(user)
    else
      res.status(404).send()
