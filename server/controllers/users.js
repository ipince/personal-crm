
var config = require('../config'),
    client = config.client,
    parsed = config.parsed


var googleLogin = function(req, res){
  res.writeHead(
    302,
    {Location: client.generateAuthUrl(
      {access_type: 'offline',
       scope: 'https://www.google.com/m8/feeds'})});
  res.end();
  return;
}

var signin = function(req, res){
  var token;
  code = parsed.query.code;
  client.getToken(code, function(err, key) {
    client.setCredentials(key);
    token = key;
  });
  res.status(200).send(token).end();
}


module.exports = {
  signin: signin,
  googleLogin: googleLogin
}
