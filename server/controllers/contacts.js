#!/usr/bin/env node

module.exports = {
  getContacts: getContacts
}

function getContacts(req, res) {
  res.end('im loading the contacts bro');
}
