var url = require('url');
var googleapis = require('googleapis')
var routes = require('../config/routes')
var flags = require('../config/flags');

var callbackUrlBase = flags.HOSTNAME + ':' + flags.PORT;
var client = new googleapis.OAuth2Client(
    flags.CLIENT_ID,
    flags.CLIENT_SECRET,
    callbackUrlBase + routes.oauthcallback);

var googleLogin = function(req, res) {
  res.writeHead(
    302,
    {Location: client.generateAuthUrl(
      {access_type: 'offline',
       scope: 'https://www.google.com/m8/feeds'})});
  res.end();
  return;
}

var signin = function(req, res) {
  var token;
  code = url.parse(req.url, true).query.code;
  console.log('the code: ' + code);
  client.getToken(code, function(err, key) {
    client.setCredentials(key);
    console.log(JSON.stringify(client.credentials));
    console.log();
    console.log(JSON.stringify(client));
    token = key;
  });
  res.status(200).send(token).end();
}


module.exports = {
  signin: signin,
  googleLogin: googleLogin,

  // TODO(rodrigo): move somewhere else.
  oauthClient: client,
}

