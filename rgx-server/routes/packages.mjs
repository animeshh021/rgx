'use strict'

import express from 'express'
import golangRouter from './packages/golang.mjs'
import gcloudRouter from './packages/gcloud.mjs'
import * as candidates from '../services/packages/candidates.mjs'

const router = express.Router()

const supportedPackages = {
    golang: golangRouter,
    gcloud: gcloudRouter
}

for (const pkgAlias of Object.keys(supportedPackages)) {
    router.use(`/${pkgAlias}`, supportedPackages[pkgAlias])
}

router.get('/', (req, res) => {
    /**
     * @openapi
     * tags:
     *   name: packages
     *   description: Endpoints for querying packages that this server supports
     * /packages:
     *   get:
     *      summary: Show all packages
     *      tags: [packages]
     *      responses:
     *        200:
     *          description: A list of packages supported by the server
     *        500:
     *          description: Server error
     */

    return res.json(candidates.all())
})

export default router