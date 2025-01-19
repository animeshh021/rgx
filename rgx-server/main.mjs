import * as path from 'node:path'
import { fileURLToPath } from 'node:url'

import express from 'express'
import swaggerUi from 'swagger-ui-express'
import swaggerJSDoc from 'swagger-jsdoc'

import log4js, { weblogger } from './common/util/log.mjs'
import clientRouter from './routes/client.mjs'
import pkgRouter from './routes/packages.mjs'
import * as utils from './common/util/utils.mjs'
import tomlConfig from './common/util/toml-config.mjs'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const log = log4js.getLogger('main')
const config = tomlConfig(utils.configFile)
const appPort = process.env.PORT || config.app.port || 9020
const port = parseInt(appPort)
const contextPath = config.app.route || 'rgx-server'
const routePrefix = process.env.DCF_ENV ? '' : `/${contextPath}`
const app = express()

utils.setupDirs()
app.use(express.json({ limit: '256kb' }))
app.set('json-spaces', 2)
app.use(weblogger)

app.use(`${routePrefix}/static`, express.static(path.join(__dirname, 'static')))

const swaggerSpecs = swaggerJSDoc(utils.swaggerOptions)
app.use(`${routePrefix}/api-docs`, swaggerUi.serve, swaggerUi.setup(swaggerSpecs))

app.get('/', (req, res) => {
    res.redirect(`/${routePrefix}`)
})

app.get('/favicon.ico', (req, res) => {
    res.sendFile(path.join(__dirname, './static/assets/icons/favicon.png'), {
        headers: { 'Content-Type': 'image/png' }
    })
})

app.get(`${routePrefix}`, (req, res) => {
    res.setHeader('Content-Type', 'text/plain')
    res.send(`rgx allows you to manage FOSS software packages.
[${new Date().toISOString()}] /rgx-server is running. Go to /rgx-server/api-docs to see the API docs.`)
})

app.use(`${routePrefix}/client`, clientRouter)
app.use(`${routePrefix}/packages`, pkgRouter)

app.listen(port, () => {
    log.info(`server started on port ${port} with route prefix = '${routePrefix}'`)
})