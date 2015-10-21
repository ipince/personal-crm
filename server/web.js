#!/usr/bin/env node

var express = require('express');
var middleware = require('./config/middleware');

var app = express();

app.set('port', process.env.PORT || 8080);

middleware(app);

var server = app.listen(app.get('port'));
console.log('Up and running');

module.exports = app;

