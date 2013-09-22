loadData = required '../'

module.exports = (grunt) ->
  grunt.initConfig
    pkg: grunt.file.readJSON "package.json"

  grunt.registerTask "loadData", 'Task that insert specified fixtures into postgres database',