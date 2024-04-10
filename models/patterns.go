package models

import "regexp"

var regex_comments = regexp.MustCompile(`{{/\*(.|\s)+?\*/}}`)
var regex_defination = regexp.MustCompile(`{{-?\s?define\s?\"(.+)\"\s?\-}}`)
var regex_tpl_usage = regexp.MustCompile(`(include|template)\s*\"(.*?)\".*?`)
var regex_yaml_comments = regexp.MustCompile(`^\s*#.*$`)
