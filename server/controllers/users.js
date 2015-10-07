#!/usr/bin/env node

module.exports = {
  googleLogin: googleLogin,
  signin: signin
}

function googleLogin(req, res) {
  res.end('im loading the contacts bro');
}

function signin(req, res) {
  res.end('im loading the contacts bro');
}
