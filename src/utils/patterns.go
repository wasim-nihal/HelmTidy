package utils

import "regexp"

// matches tpl comments
var RgxTplComments = regexp.MustCompile(`{{/\*(.|\s)+?\*/}}`)

// matches tpl definations
var RgxTplDefinations = regexp.MustCompile(`{{-?\s?define\s?\"(.+)\"\s?\-}}`)

// matches tpl usages (either through `template` or `include`)
var RgxTplUsages = regexp.MustCompile(`(include|template)\s*\"(.*?)\".*?`)

// matches yaml comments
var RgxYamlComments = regexp.MustCompile(`^\s*#.*$`)
