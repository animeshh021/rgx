'use-strict'

import log4js from 'log4js'

// allowed values: https://stritti.github.io/log4js/docu/users-guide.html#configuration
const logLevel = process.env.APP_LOG_LEVEL ?? 'TRACE'

log4js.configure({
    appenders: {
        main: { type: 'stdout' },
        web: { type: 'stdout' }
    },
    categories: {
        default: { appenders: ['main'], level: logLevel },
        web: { appenders: ['web'], level: 'INFO' }
    }
})

export const weblogger = log4js.connectLogger(log4js.getLogger('web'), { level: log4js.levels.INFO })
export default log4js