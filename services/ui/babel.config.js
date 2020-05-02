const moduleResolver = require('babel-plugin-module-resolver')

module.exports = {
  presets: [
    [
      moduleResolver,
      {
        "root": [
          "./src"
        ],
        "alias": {
          "src": "./src"
        }
      }
    ]
  ]
}