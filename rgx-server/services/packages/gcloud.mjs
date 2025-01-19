import log4js from "../../common/util/log.mjs";
import * as cache from "../../common/util/cache.mjs";
import { supportedPlatforms, compatibleTimestamp } from "../../common/util/utils.mjs";

const log = log4js.getLogger("gcloud-service");
const gcloudUrl = "https://storage.googleapis.com/cloud-sdk-release/";

export async function majorVersions() {
    const ckey = 's:gcloud:majorversions'
    const data = await cache.cget(ckey);
    if (data) return { ok: true, data };

    const r = ["As of now, all versions from 100.0.0 to 502.0.0 are available. Also, installation takes around 10 minutes."];
    await cache.cput(ckey, r);
    return { ok: true, data: r };
}

function packageDetails(packageUrl) {
    const packageInfo = {
        name: '',
        version: '',
        os: '',
        arch: '',
    };

    const urlParts = packageUrl.split('/');
    const filename = urlParts[urlParts.length - 1];

    const windowsRegex = /google-cloud-sdk-(\d+\.\d+\.\d+)-windows-(\w+)-bundled-python\.zip/;
    const otherRegex = /google-cloud-sdk-(\d+\.\d+\.\d+)-(\w+)-(\w+)\.tar\.gz/;

    const windowsMatch = filename.match(windowsRegex);
    const otherMatch = filename.match(otherRegex);

    if (windowsMatch) {
        packageInfo.name = filename;
        packageInfo.version = windowsMatch[1];
        packageInfo.os = 'windows';
        packageInfo.arch = windowsMatch[2];
    } else if (otherMatch) {
        packageInfo.name = filename;
        packageInfo.version = otherMatch[1];
        packageInfo.os = otherMatch[2];
        packageInfo.arch = otherMatch[3];
    } 

    return packageInfo; 
}

function mapOS (os) {
    switch (os) {
        case 'macos':
            return 'darwin';
        case 'linux':
            return 'linux';
        case 'windows':
            return 'windows';
        default:
            log.error('unknown os:', os)
            return ''
    }
}

function mapArch (arch) {
    switch (arch) {
        case 'x64':
            return 'x86_64';
        case 'aarch64':
            return 'arm';
        case 'x86-64':
            return 'x86_64';
        case 'arm64':
            return 'arm';
        case 'arm':
            return 'arm';
        default:
            log.error('unknown cpu arch:', arch)
            return ''
    }
}

const scriptLinks = {
    windows: '/static/assets/install-scripts/gcloud/rgx-setup.cmd',
    macos: '/static/assets/install-scripts/gcloud/rgx-setup.sh',
    linux: '/static/assets/install-scripts/gcloud/rgx-setup.sh'
}

function createRecipe (version, pkg, opsys, pkgUrl) {
    const recipe = {
        script: scriptLinks[opsys],
        script_dir: `google-cloud-sdk/gcloudsdk-${version}`,
        package_version: version,
        artifacts: [
            {
                artifact_type: 'google-cloud-sdk',
                action: 'extract',
                name: pkg.name,
                extract_dir: `google-cloud-sdk/gcloudsdk-${version}`,
                extract_target: `google-cloud-sdk/gcloudsdk-${version}`,
                version,
                link: pkgUrl,
            }
        ]
    }
    return recipe
}

export async function latestRelease(p) {

    const { majorVersion, os: providedOS, arch: providedArch } = p;
    const os = mapOS(providedOS);
    const arch = mapArch(providedArch);
    if (!os || !arch) return { ok: false, error: `Unsupported os/arch combination: ${providedOS}/${providedArch}. ${supportedPlatforms}` };

    let url = '';
    if (os === 'windows') {
        url = gcloudUrl + 'google-cloud-sdk-' + majorVersion + '-windows-' + arch + '-bundled-python.zip';
    } else {
        url = gcloudUrl + 'google-cloud-sdk-' + majorVersion + '-' + os + '-' + arch + '.tar.gz';
    }

    const r = packageDetails(url);

    const ckey = `s:gcloud:latestrelease:${majorVersion}-${os}-${arch}`;
    const data = await cache.cget(ckey);
    if (data) { return { ok: true, data }; }

    try {
        const recipe = createRecipe(majorVersion, {
            name: r.name,
            os,
            arch,
            link: url,
        }, providedOS, url);
        
        await cache.cput(ckey, recipe);

        return { ok: true, data: recipe };
    } catch (e) {
        return { ok: false, error: e.message };
    }
}

export async function clearCache () {
    await cache.removeKeys('s:gcloud:');
}