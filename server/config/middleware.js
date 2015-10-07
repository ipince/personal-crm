#!/usr/bin/env node

var usersController = require('../controllers/users');
var contactsController = require('../controllers/contacts');

module.exports = function(app) {
  filters(app);
  routes(app);
}

function filters(app) {
  app.use(function(req, res, next) {
    console.log('im in your middleware, eating yo cheese');
    next();
  });
}

function routes(app) {
  app.get('/login', usersController.googleLogin);
  app.get('/oauthcallback', usersController.signin);
  app.get('/contacts', contactsController.getContacts);
}
