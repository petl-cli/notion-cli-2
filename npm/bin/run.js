#!/usr/bin/env node
'use strict'

const { spawnSync } = require('child_process')
const path = require('path')
const fs = require('fs')

const CLI_NAME = 'notion-api'
const ext = process.platform === 'win32' ? '.exe' : ''
const binaryPath = path.join(__dirname, CLI_NAME + ext)

if (!fs.existsSync(binaryPath)) {
  console.error(CLI_NAME + ' binary not found. Try reinstalling: npm install -g @petl-cli/' + CLI_NAME)
  process.exit(1)
}

const result = spawnSync(binaryPath, process.argv.slice(2), { stdio: 'inherit' })
process.exit(result.status ?? 1)
