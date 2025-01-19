'use-strict'

import express from 'express'
import * as utils from '../../common/util/utils.mjs'
import * as golang from '../../services/packages/golang.mjs'
import * as candidates from '../../services/packages/candidates.mjs'

const router = express.Router()

router.get('/', (req, res) => {
    /**
     * @openapi
     * tags:
     *   name: golang
     *   description: APIs for working with this package
     * /packages/golang:
     *   get:
     *      summary: Display information about this pacakge
     *      tags: [golang]
     *      responses:
     *        200:
     *          description: Information about this package
     *        500:
     *          description: Server error
     */

    return res.json(candidates.candidate('golang'))
})

router.get('/versions', async (req, res) => {
    /**
     * @openapi
     * /packages/golang/versions:
     *   get:
     *     summary: Display the major versions available
     *     tags: [golang]
     *     responses:
     *        200:
     *          description: A list of major versions
     *        500:
     *          description: Server error
     */

    const { lts } = req.query
    if (lts) return utils.errorText(res, 400, 'Go does not have LTS versions')

    const r = await golang.majorVersions()
    if (!r.ok) return utils.errorText(res, 500, r.error)
    return res.json(r.data)
})

router.get('/release/:majorVersion/:os/:arch', async (req, res) => {
    /**
     * @openapi
     * /packages/golang/release/{majorVersion}/{os}/{arch}:
     *   get:
     *     summary: Display version and download information about a specific release
     *     tags: [golang]
     *     parameters:
     *       - in: path
     *         name: majorVersion
     *         type: float
     *         required: true
     *         description: The major version of golang required, e.g. 1.21, 1.22
     *       - in: path
     *         name: os
     *         type: string
     *         required: true
     *         description: The OS for which this release is -- windows, linux, or mac
     *       - in: path
     *         name: arch
     *         type: string
     *         required: true
     *         description: The CPU architecture fpr which this release is -- x64 (Intel), or aarch64 (Mac M1-series and some others)
     *     responses:
     *        200:
     *          description: A list of major versions
     *        500:
     *          description: Server error
     */

    const { majorVersion, os, arch } = req.params
    const r = await golang.latestRelease({ majorVersion, os, arch }, {
        installation: req.headers['x-rgx-installation']
    })
    if (!r.ok) return utils.errorText(res, r?.code ? r.code : 500, r.error)
        return res.json(r.data)
})

router.get('/clear-cache', async (req, res) => {
    await golang.clearCache()
    return res.json({ ok: true, message: 'cache cleared' })
})

export default router