config = require '../config'
topic = require '../models/topic'

module.exports = (app) ->
  app.get '/topic/:id', getById

getById = (req, res) ->

  topicId = req.params.id
  topic.getById topicId, (topic) =>
    if topic
      res.header "Content-Type", "application/json"
      res.send(topic)
    else
      res.status(404).send()
