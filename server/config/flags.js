#!/usr/bin/env node

require('dotenv').load()

module.exports = {
  CONTACTS_URL: process.env.CONTACTS_URL,
  CLIENT_ID: process.env.CLIENT_ID,
  CLIENT_SECRET: process.env.CLIENT_SECRET,
  HOSTNAME: process.env.HOSTNAME,
  PORT: process.env.PORT,
}

