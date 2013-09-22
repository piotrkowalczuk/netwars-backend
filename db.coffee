knex = require 'knex'
config = require './config'
redis = require 'redis'

module.exports =
  postgres: knex.initialize(config.database)
  redis: redis.createClient()