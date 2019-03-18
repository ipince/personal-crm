#!/usr/bin/env node

var usersController = require('../controllers/users');
var contactsController = require('../controllers/contacts');
var routes = require('./routes');

module.exports = function(app) {
  filter(app);
  route(app);
}

function filter(app) {
  app.use(function(req, res, next) {
    console.log('im in your middleware, eating yo cheese');
    next();
  });
}

function route(app) {
  app.get(routes.login, usersController.googleLogin);
  app.get(routes.oauthcallback, usersController.signin);
  app.get(routes.contacts, contactsController.getContacts);
}
