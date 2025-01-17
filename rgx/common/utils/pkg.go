package utils

const Version = "0.1.0"

const ApplicationName = "rgx"
const ApplicationShortDescription = "A CLI based package manager"
const ApplicationDescription = "rgx allows you to manage software packages through CLI"
const UserAgent = ApplicationName + "/" + Version

var Config RgxConfig
var CurrentRuntimeConfig = GetRuntimeConfig()
