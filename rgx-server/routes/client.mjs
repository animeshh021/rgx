import express from 'express'

const router = express.Router()

router.get('/latest', (req, res) => {
    res.setHeader('Content-Type', 'text/plain; charset=utf-8')
    res.status(200).send('0.1.0\n')
})

export default router