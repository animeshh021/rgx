'use-strict'

import express from 'express'
import * as utils from '../../common/util/utils.mjs'
import * as gcloud from '../../services/packages/gcloud.mjs'
import * as candidates from '../../services/packages/candidates.mjs'

const router = express.Router()

router.get('/', (req, res) => {
    /**
     * @openapi
     * tags:
     *   name: gcloud
     *   description: APIs for working with this package
     * /packages/gcloud:
     *   get:
     *      summary: Display information about this pacakge
     *      tags: [gcloud]
     *      responses:
     *        200:
     *          description: Information about this package
     *        500:
     *          description: Server error
     */

    return res.json(candidates.candidate('gcloud'))
})

router.get('/versions', async (req, res) => {
    /**
     * @openapi
     * /packages/gcloud/versions:
     *   get:
     *     summary: Display the major versions available
     *     tags: [gcloud]
     *     responses:
     *        200:
     *          description: A list of major versions
     *        500:
     *          description: Server error
     */

    const { lts } = req.query
    const r = await gcloud.majorVersions({ lts })
    if (!r.ok) return utils.errorText(res, 500, r.error)
    return res.json(r.data)
})

router.get('/release/:majorVersion/:os/:arch', async (req, res) => {
    /**
     * @openapi
     * /packages/gcloud/release/{majorVersion}/{os}/{arch}:
     *   get:
     *     summary: Display version and download information about a specific release
     *     tags: [gcloud]
     *     parameters:
     *       - in: path
     *         name: majorVersion
     *         type: integer
     *         required: true
     *         description: The major version of gcloud required, e.g. 499.0.0, 387.0.0
     *       - in: path
     *         name: os
     *         type: string
     *         required: true
     *         description: The OS for which this release is -- windows, linux, or mac
     *       - in: path
     *         name: arch
     *         type: string
     *         required: true
     *         description: The CPU architecture fpr which this release is -- x86-64 (Intel), or aarch64 (Mac M1-series and some others)
     *     responses:
     *        200:
     *          description: A list of major versions
     *        500:
     *          description: Server error
     */

    const { majorVersion, os, arch } = req.params
    const r = await gcloud.latestRelease({ majorVersion, os, arch }, {
        installation: req.headers['x-rgx-installation']
    })
    if (!r.ok) return utils.errorText(res, r?.code ? r.code : 500, r.error)
        return res.json(r.data)
})

router.get('/clear-cache', async (req, res) => {
    await gcloud.clearCache()
    return res.json({ ok: true, message: 'cache cleared' })
})

export default router