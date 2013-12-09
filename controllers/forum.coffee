config = require '../config'
forum = require '../models/forum'

module.exports = (app) ->
  app.get '/forum/:id', getById

getById = (req, res) ->

  forumId = req.params.id
  forum.getById forumId, (forum) =>
    console.log forum
    if forum
      res.header "Content-Type", "application/json"
      res.send(forum)
    else
      res.status(404).send()
