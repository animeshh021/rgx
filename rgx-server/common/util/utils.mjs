'use-strict'

import fs from 'node:fs'
import path from 'node:path'
import log4js from './log.mjs'
import tomlConfig from './toml-config.mjs'

const log = log4js.getLogger('utils')

export const appEnv = process.env.DCF_ENV || {
    dev: 'dev',
    'uat-1': 'uat',
    prd: 'prd'
}[process.env.EFX_ENV ?? 'dev']

export const configProfile = process.env.APP_CONFIG_PROFILE ?? 'kalypso'
export const configFile = process.env.APP_CONFIG_PROFILE || findInHierarchy(`service-config.${appEnv}.toml`)
if (!configFile) {
    log.fatal(`app config profile = '${configProfile}', env = '${appEnv}'`)
    log.fatal('Ensure config/$APP_CONFIG_PROFILE/service-config.$ENV.toml exists, or set APP_CONFIG_FILE to the config filename')
    log.fatal('could not determine the config filename, exiting.')
    process.exit(1)
}
const config = tomlConfig(configFile)

export function appDirs() {
    let dataDir = config.app.data_dir
    if (dataDir.includes('{HOME}')) {
        dataDir = path.join(homeDirectory(), dataDir.replace('{HOME}', ''))
    }
    return {
        cache: path.join(dataDir, 'cache')
    }
}

export function setupDirs () {
    const d = appDirs()

    if (!fs.existsSync(d.cache)) fs.mkdirSync(d.cache, { recursive: true })
    log.debug('cache dir:', d.cache)
}

function findInHierarchy (basename) {
    let curdir = '.'
    let absdir
    do {
        absdir = fs.realpathSync(curdir)
        const absfilename = path.join(absdir, 'config', configProfile, basename)
        if (fs.existsSync(absfilename)) {
            return absfilename
        }
        curdir += '/..'
    } while (absdir !== '/')
    return null
}

export function errorText (resp, code, message) {
    resp.setHeader('Content-type', 'text/plain')
    resp.status(code).send(message)
}

export function readJson (fname) {
    const d = fs.readFileSync(fname, 'utf8')
    return JSON.parse(d)
}

export const swaggerOptions = {
    definition: {
        openapi: '3.1.0',
        info: {
            title: 'rgx-server API',
            version: '0.1.0',
            description:
                'API for rgx-server',
            license: {
                name: 'Copyright {c)} AJ. All Rights Reserved.',
                url: 'https://xyz.com'
            },
            contact: {
                name: 'rgx dev team',
                url: 'https://todo',
                email: 'xxxxx@gmail.com'
            }
        },
        servers: [
            {
                url: 'https://rgx-api.uat.cloud.xyz'
            },
            {
                url: 'http://localhost:9020/rgx-server'
            }
        ]
    },
    apis: ['./routes/**/*.mjs']
}

/**
 * This is intended to help sort versions with small-ish parts only, like 1.2.3 or 1.2.104.5,
 * such as the version numbers use by NodeJS! It has *not* been tested for version strings
 * with timestamp-ish values, e.g. 7.0.15.20240503.200623.7 or version numbers with Build IDs 
 * like 1.17.3894701.14.
 */
export function normalizeVersion (v) {
    const parts = v.split('.').map(x => parseInt(x))
    if (parts.length > 5) {
        log.warn(`warning: the versopm string '${v}' has too many parts!`)
    }
    let c = ''
    for (const p of parts) {
        c += String(p).padStart(5, '0') + '.'
    }
    return c
}

export function strcmp (a, b) {
    if (a < b) return -1
    if (a > b) return 1
    return 0
}

export function homeDirectory () {
    if (process.platform.startsWith('win')) return process.env.USERPROFILE
    return process.env.HOME
}

export function visibleServer() {
    const serverName = config.app.visible_server
    const serverPort = config.app.visible_port
    const contextPath = config.app.context_path
    const proto = serverName === 'localhost' || serverPort === '80' ? 'http' : 'https'
    return `${proto}://${serverName}:${serverPort}/${contextPath}`
}

export const supportedPlatforms = "It currently supports the following. OSes: [windows, macos, linux], CPU Architectures: [x86-64, arm64]."

export function ipAddress (req) {
    if (req.headers['x-forwarded-for']) {
        const ips = req.headers['x-forwarded-for']
        .split(',')
        .map(ip => ip.trim())
        return ips[0] // per HTTP spec, the first address is always the client
    }

    const ip = req.connection.remoteAddress || req.socket.remoteAddress || (req.connection.socket ? req.connection.remoteAddress : null)
    return ip
}

export function compatibleTimestamp () {
    let ts = new Date().toISOString()
    ts = ts.replace('T', ' ')
    ts = ts.substring(0, 19) + '+00:00'
    return ts
}