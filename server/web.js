

var httplib = require('http');
var httpslib = require('https');
var urllib = require('url');
var googleapis = require('googleapis');
require('dotenv').load();

OAuth2Client = googleapis.OAuth2Client;

var CLIENT_ID = 'REDACTED';
var CLIENT_SECRET = 'REDACTED';
var CALLBACK = 'http://127.0.0.1:8080/oauthcallback';
var CONTACTS_URL = 'https://www.google.com/m8/feeds/contacts/default/full?alt=json&max-results=10000&access_token=';

var port = process.env.PORT || 8080;


var client = new OAuth2Client(CLIENT_ID, CLIENT_SECRET, CALLBACK);

var server = httplib.createServer(
  function(req, res) {
    parsed = urllib.parse(req.url, true);
    if (parsed.pathname === '/login') {
      // redirect
      res.writeHead(
        302,
        {Location: client.generateAuthUrl(
          {access_type: 'offline',
           scope: 'https://www.google.com/m8/feeds'})});
      res.end();
      return;
    } else if (parsed.pathname === '/oauthcallback') {
      code = parsed.query.code;
      client.getToken(code, function(err, tokens) {
        client.setCredentials(tokens);
        console.log(tokens);
      });
      res.end('callback');
    } else if (parsed.pathname === '/contacts') {
      console.log('fetching contacts');
      var contacts;
      httpslib.get(CONTACTS_URL + client.credentials.access_token, function(res) {
        console.log('got data back');
        var responseContents = ''
        res.on('data', function(data) { responseContents += data; });
        res.on('end', function() {
          contacts = JSON.parse(responseContents);
          console.log('just after parse');
          console.log(contacts);
          for (var i = 0; i < contacts.feed.entry.length; i++) {
            console.log(contacts.feed.entry[i].title.$t);
          }
        });
      });
      res.end('fetched, maybe');
    }
    res.writeHead(
      200,
      {'Content-Type': 'text/plain'}
    );
    res.end('Hello from node!\n');
  }
);

server.listen(process.env.PORT);

console.log('Up and running');
