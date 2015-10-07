
var https = require('https');
var config = require('../config'),
    CONTACTS_URL = config.CONTACTS_URL,
    access_token = config.client.credentials.access_token

var getContacts = function(req, res){
  console.log('fetching contacts');
  var contacts;
  https.get(CONTACTS_URL + access_token, function(res) {
    console.log('got data back');
    var responseContents = ''
    res.on('data', function(data) { responseContents += data; });
    res.on('end', function() {
      contacts = JSON.parse(responseContents);
      console.log('just after parse');
      console.log(contacts);
      res.status(200).send(contacts).end();
    });
  });
}

module.exports = {
  getContacts: getContacts
}
