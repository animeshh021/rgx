import { Level } from 'level'

import * as utils from './utils.mjs'
import tomlConfig from './toml-config.mjs'
import log4js from './log.mjs'

const config = tomlConfig(utils.configFile)
const log = log4js.getLogger('cache')

const cache = new Level(utils.appDirs().cache, { keyEncoding: 'binary', valueEncoding: 'json' })
const cachePeriodHours = parseInt(config.app.cache_period_hours)
const maxage = (cachePeriodHours && cachePeriodHours <= 150 && cachePeriodHours > 0 ? cachePeriodHours : 8) * 60 * 60 * 1000

export async function cget (key) {
    try {
        const r = await cache.get(key)
        if (r?.ts) {
            const age = new Date().getTime() - r.ts
            return age < maxage ? r?.val : null
        } else {
            return null
        }
    } catch (e) {
        return null
    }
}

export async function cput (key, val) {
    const ts = new Date().getTime()
    const err = await cache.put(key, { ts, val })
    if (err) {
        log.error(`could not write cache for '${key}':`, err.message)
        // throw err
    }
}

export async function removeKeys (pattern) {
    for await (const key of cache.keys({ gte: `${pattern}!`, lte: `${pattern}~` })) {
        await cache.del(key)
    }
}