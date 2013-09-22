knex = require 'knex'
config = require './config'
redis = require 'redis'

module.exports =
  postgres: knex.Initialize(config.database)
  redis: redis.createClient()