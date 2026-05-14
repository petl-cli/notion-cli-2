#!/usr/bin/env node
'use strict'

const https = require('https')
const fs = require('fs')
const path = require('path')

const CLI_NAME = 'notion-api'
const OWNER = 'petl-cli'
const REPO = 'notion-cli-2'
const VERSION = '0.1.0'

const PLATFORM_MAP = { darwin: 'darwin', linux: 'linux', win32: 'windows' }
const ARCH_MAP = { x64: 'amd64', arm64: 'arm64' }

const platform = PLATFORM_MAP[process.platform]
if (!platform) { console.error('Unsupported platform: ' + process.platform); process.exit(1) }

const arch = ARCH_MAP[process.arch]
if (!arch) { console.error('Unsupported architecture: ' + process.arch); process.exit(1) }

const ext = platform === 'windows' ? '.exe' : ''
const binaryName = CLI_NAME + '-' + platform + '-' + arch + ext
const url = 'https://github.com/' + OWNER + '/' + REPO + '/releases/download/v' + VERSION + '/' + binaryName
const binDir = path.join(__dirname, '..', 'bin')
const destPath = path.join(binDir, CLI_NAME + ext)

if (!fs.existsSync(binDir)) fs.mkdirSync(binDir, { recursive: true })

console.log('Downloading ' + CLI_NAME + ' for ' + platform + '/' + arch + '...')

function download(url, dest, cb) {
  const file = fs.createWriteStream(dest)
  https.get(url, (res) => {
    if (res.statusCode === 301 || res.statusCode === 302) {
      file.close()
      fs.unlink(dest, () => {})
      return download(res.headers.location, dest, cb)
    }
    if (res.statusCode !== 200) {
      file.close()
      fs.unlink(dest, () => {})
      return cb(new Error('Download failed with status ' + res.statusCode))
    }
    res.pipe(file)
    file.on('finish', () => file.close(cb))
  }).on('error', (err) => { fs.unlink(dest, () => {}); cb(err) })
}

download(url, destPath, (err) => {
  if (err) { console.error('Failed to download ' + CLI_NAME + ': ' + err.message); process.exit(1) }
  if (platform !== 'windows') fs.chmodSync(destPath, 0o755)
  console.log(CLI_NAME + ' installed successfully.')
})
