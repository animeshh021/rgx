'use-strict'

import fetch, { Response } from 'node-fetch'
import https from 'node:https'
import { error, time } from 'node:console'
import exp from 'node:constants'


async function fetchWithTimeout (resource, options = {}) {
    const { timeout = 60 } = options
    if (options.timeout) delete options.timeout
    const controller = new AbortController()
    const id = setTimeout(() => controller.abort(), timeout * 1000)
    const response = await fetch(resource, {
        ...options,
        signal: controller.signal
    })
    clearTimeout(id)
    return response
}

function errorHandler (err) {
    return new Response(
        JSON.stringify({
            msgtype: 'error',
            message: err.message
        })
    )
}

const textErrorPrefix = '//__[error]__//:'

function textEErrorHandler (err) {
    return new Response(`${textErrorPrefix}${err.message}`)
}

function ensureSuccess (resp) {
    if(!resp.ok) {
        throw Error(`error fetching ${resp.url}: HTTP status ${resp.status}, ${resp.statusText}`)
    }
    return resp
}

export async function getJson (url, options = {}) {
    if (url.startsWith('https:')) options.agent = agent
    const resp = await fetchWithTimeout(url, options).then(ensureSuccess).catch(errorHandler)
    let data
    try {
        data = await resp.json()
    } catch (ex) {
        return { ok: false, error: ex.message }
    }
    return data?.msgtype === 'error' ? { ok: false, error: data.message } : { ok: true, data }
}

export async function getText (url) {
    // if (url.startsWith('https:')) options.agent = agent
    const resp = await fetchWithTimeout(url).then(ensureSuccess).catch(errorHandler)
    const data = await resp.text()
    return data.startsWith(textErrorPrefix) ? { ok: false, error: data.substring(textErrorPrefix.length) } : { ok: true, data }
}

export async function postJson (url, jsonObject, options = {}) {
    if (url.startsWith('https:')) options.agent = agent

    if (!options?.headers) {
        options = {
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json'
            }
        }
    } else if (options?.headers) {
        if (!options.headers.Accept) options.headers.Accept = 'application/json'
        if (!options.headers['Content-Type']) options.headers['Content-Type'] = 'application/json'
    }

    const resp = await fetchWithTimeout(url, {
        ...options,
        method: 'POST',
        body: JSON.stringify(jsonObject)
    }).then(ensureSuccess).catch(errorHandler)

    let data
    try {
        data = await resp.json()
    } catch (ex) {
        return { ok: false, error: ex.message }
    }
    return data?.msgtype === 'error' ? { ok: false, error: data.message } : { ok: true, data }
}

export async function collectStatusAfterPost (url, jsonObject, options = {}) {
    if (url.startsWith('https:')) options.agent = agent

    if (!options?.headers) {
        options = {
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json'
            }
        }
    } else if (options?.headers) {
        if (!options.headers.Accept) options.headers.Accept = 'application/json'
        if (!options.headers['Content-Type']) options.headers['Content-Type'] = 'application/json'
    }

    const resp = await fetchWithTimeout(url, {
        ...options,
        method: 'POST',
        body: JSON.stringify(jsonObject)
    }).then(ensureSuccess).catch(errorHandler)

    let status = 0
    let data 
    try {
        status = resp.status
        data = await resp.text()
    } catch (ex) {
        return { ok: false, error: ex.message }
    }
    return data.startsWith(textErrorPrefix) ? { ok: false, error: data.substring(textErrorPrefix.length) } : { ok: true, status, data }
}