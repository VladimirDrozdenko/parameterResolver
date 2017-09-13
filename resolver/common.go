package resolver

import "regexp"

const ssmNonSecurePrefix = "ssm:"
const ssmSecurePrefix = "ssm-secure:"

const secureStringType = "SecureString"
const stringType = "String"

//
// SSM Parameter placeholder - relaxed regular expression
var parameterPlaceholder = regexp.MustCompile("{{\\s*(" + ssmNonSecurePrefix + "[\\w-/]+)\\s*}}")
var secureParameterPlaceholder = regexp.MustCompile("{{\\s*(" + ssmSecurePrefix + "[\\w-/]+)\\s*}}")

type ResolveOptions struct {
	IgnoreSecureParameters bool
}

type SsmParameterInfo struct {
	Name  string
	Type  string
	Value string
}
