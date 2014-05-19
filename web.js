

var http = require('http');
var port = process.env.PORT || 8080;

var server = http.createServer(
  function(req, res) {
    res.writeHead(
      200,
      {'Content-Type': 'text/plain'}
    );
    res.end('Hello world from node\n');
  }
);

server.listen(port);

console.log('Up and running');
