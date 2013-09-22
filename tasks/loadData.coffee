db = require '../db'

module.exports = (grunt)->

  grunt.registerTask 'loadData', 'Task that load data to postgres database', () ->

    fixtures = require '../'+grunt.option('fixtures')

    data = fixtures.data
    table = fixtures.table

    if table and data

      done = @async()
      grunt.log.writeln 'Inserting data...'

      db.postgres(table)
        .insert(data)
        .exec (error, reply) ->
          if not error
            grunt.log.ok "#{data.length} object inerted"
          else
            grunt.log.error error
          done()

