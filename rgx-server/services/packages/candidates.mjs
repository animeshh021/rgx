import * as utils from '../../common/util/utils.mjs'

const allCandidates = utils.readJson('./common/app/candidates.json')

export function candidate (name) {
    return allCandidates[name]
}

export function all () {
    return Object.keys(allCandidates).toSorted().map(k => {
        return {
            name: k,
            description: allCandidates[k].description
        }
    })
}