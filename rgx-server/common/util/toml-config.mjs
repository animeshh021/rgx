import fs from 'node:fs'
import toml from 'toml'
import log4js from './log.mjs'

const log = log4js.getLogger('toml-config')

let appConfig

/** Read  TOML file into a dict.
 * If the TOML file haas an [env] section, those config
 * entries cn be over-ridden with environment variables.
 */
function read (configFile) {
    const config = toml.parse(fs.readFileSync(configFile, 'utf-8'))

    if ('env' in config) {
        for(const [k] of Object.entries(config.env)) {
            //log.debug(`toml-config: read: env config key found, ${k}=${v}`)
            let readEnv = process.env[k.toUpperCase()]
            if (readEnv) {
                log.debug(`toml-config: read: overriding ${k} with the value from ${k.toUpperCase()}:`, readEnv)
                if (Array.isArray(config.env[k])) {
                    config.env [k] = readEnv.split(',')
                } else {
                    readEnv = parseInt(readEnv) ? parseInt(readEnv) : readEnv
                    if (readEnv === '0') readEnv = 0
                    config.env[k] = readEnv
                }
            }
        }
    }
    appConfig = config
    return config
}

export default function (fname) {
    return appConfig || read(fname)
}