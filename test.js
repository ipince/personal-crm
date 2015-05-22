console.log('test');

var google = require('googleapis');

var fs = require('fs');
var contents = fs.readFile(
  'hello.txt',
  'utf8',
  function(err, contents) {
    console.log(contents);
  }
);

