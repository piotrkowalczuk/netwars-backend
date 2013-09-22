express = require 'express'
config = require './config'

app = express()
app.set('title', 'Netwars')
app.use(express.bodyParser());


app.listen config.app.port
console.log 'Listening on port '+config.app.port