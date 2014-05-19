

var http = require('http');

var server = http.createServer(
  function(req, res) {
    res.writeHead(
      200,
      {'Content-Type': 'text/plain'}
    );
    res.end('Hello world from node\n');
  }
);

server.listen(8080);

console.log('Up and running');
