'use-strict'

import { parse } from 'node-html-parser'

import { normalizeVersion, strcmp, ipAddress, supportedPlatforms } from '../../common/util/utils.mjs'
import * as fetcher from '../../common/util/fetcher.mjs'
import log4js from '../../common/util/log.mjs'
import * as cache from '../../common/util/cache.mjs'


const log = log4js.getLogger('golang-service')
const goDownloadURL = 'https://go.dev/dl/'

const goDownloadRegex = /go([0-9.]+)/

async function goReleases () {
    const r = await fetcher.getText(goDownloadURL)
    if(!r.ok) return r

    const oslist = ['macos', 'linux', 'windows']
    const archlist = ['arm64', 'x86-64']
    const root = parse(r.data)
    const elems = root.querySelectorAll('tr')
    const releases = []

    for (const elem of elems) {
        const cells = elem.querySelectorAll('td')
        let a
        try {
            a = {
                filename: cells[0].textContent?.trim(),
                kind: cells[1].textContent?.trim()?.toLowerCase(),
                os: cells[2].textContent?.trim()?.toLowerCase(),
                arch: cells[3].textContent?.trim()?.toLowerCase(),
                sha256: cells[5].textContent?.trim()?.toLowerCase(),
            }
        } catch (ex) {}
        if (a && a.kind === 'archive' && oslist.includes(a.os) && archlist.includes(a.arch)) {
            const m = a.filename.match(goDownloadRegex)
            let minorVersion
            if (m) {
                a.version = m[1].substring(0, m[1].length - 1)
                const p = a.version.split('.')
                if (p.length >= 2) {
                    minorVersion = parseInt(p[1])
                }
            } if (minorVersion && minorVersion > 10) {
                releases.push(a)
            }
        }
    }
    return { ok: true, releases }
}

export async function majorVersions () {
    const ckey = 's:golang:majorversions'
    const data = await cache.cget(ckey)
    if (data) return { ok: true, data }

    const s = new Set()
    const r = await goReleases()
    if (!r.ok) return r

    r.releases.forEach(x => {
        const p = x.version.split('.')
        if (p.length >= 2) {
            const m = p[0] + '.' + p[1]
            if (m) s.add(m)
        }
    })
    const sortedVersions = [...s].toSorted((a,b) => strcmp(normalizeVersion(a), normalizeVersion(b)))
    await cache.cput(ckey, sortedVersions)
    return { ok: true, data: sortedVersions }
}

const scriptLinks = {
    windows: '/static/assets/install-scripts/golang/rgx-setup.cmd',
    macos: '/static/assets/install-scripts/golang/rgx-setup.sh',
    linux: '/static/assets/install-scripts/golang/rgx-setup.sh',
}

function createRecipe (version, pkg, opsys) {
    const recipe = {
        script: scriptLinks[opsys],
        script_dir: 'golang',
        package_version: version,
        artifacts: [
            {
                artifact_type: 'golang-sdk',
                action: 'extract',
                name: pkg.name,
                extract_dir: `golang/go-${version}`,
                extract_target: `golang/go-${version}`,
                version,
                link: pkg.link,
                checksum: pkg.checksum,
                checksum_type: 'sha256'
            }
        ]
    }
    return recipe
}

function mapOS (os) {
    switch (os) {
        case 'macos':
        case 'linux':
        case 'windows':
            return os;
        default:
            log.error('unknown os:', os)
            return ''
    }
}

function mapArch (arch) {
    switch (arch) {
        case 'x64':
            return 'x86-64'
        case 'x86-64':
            return arch
        case 'arm64':
            return arch
        default:
            log.error('unknown cpu arch:', arch)
            return ''
    }
}

export async function latestRelease (p) {
    const { majorVersion, os: providedOS, arch: providedArch } = p
    const os = mapOS(providedOS)
    const arch = mapArch(providedArch)
    if (!os || !arch) return { ok: false, error: `Unsupported os/arch combination: ${providedOS}/${providedArch}. ${supportedPlatforms}` }

    const ckey = `s:golang:latestrelease:${majorVersion}-${os}-${arch}`
    const data = await cache.cget(ckey)
    if (data) return { ok: true, data }

    try{
        const s = new Set()
        const r = await goReleases()
        if(!r.ok) return r

        r.releases.filter(x => x.version.startsWith(majorVersion) && x.os === os && x.arch === arch).forEach(x => s.add(x.version))
        const sortedVersions = [...s].toSorted((a,b) => strcmp(normalizeVersion(a), normalizeVersion(b)))
        const latest = sortedVersions[sortedVersions.length - 1]

        const candidates = r.releases.filter(x => x.version === latest && x.os === os && x.arch === arch)
        if (candidates.length === 0) return { ok: false, code: 404, error: 'nothing found' }
        if (candidates.length > 1) log.warn('found more than one candidate, returning the first!')
    
        const release = candidates[0]
        release.link = `${goDownloadURL}${release.filename}`
        delete release.kind

        const recipe = createRecipe(latest, {
            name: release.filename,
            os,
            arch,
            link: release.link,
            checksum: release.sha256
        }, providedOS)

        await cache.cput(ckey, recipe)

        return {ok: true, data: recipe}
    } catch (e) {
        return { ok: false, error: e.message }
    }
}

export async function clearCache () {
    await cache.removeKeys('s:golang:')
}