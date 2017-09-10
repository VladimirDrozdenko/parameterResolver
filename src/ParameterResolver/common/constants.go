package common

import "regexp"

//
// Maximum number of parameters that can be requested from SSM Parameter store in one GetParameters request
const MaxParametersRetrievedFromSsm = 10

//
// Maximum file size in bytes
const MaxFileSizeInBytes = 1024 * 1024 * 1024

//
// Supported output file formats
const TextOutputFormat = "txt"
const XmlOutputFormat = "xml"
const YmlOutputFormat = "yml"
const JsonOutputFormat = "json"

//
// Supported options
const FailOnParameterNotFoundOption = "failonparameternotfound"
const IgnoreParameterNotFoundOption = "ignoreparameternotfound"


const SSMNonSecurePrefix = "ssm:"
const SSMSecurePrefix = "ssm-secure:"

//
// SSM Parameter placeholder - relaxed regular expression
var ParameterPlaceholder = regexp.MustCompile("{{\\s*" + SSMNonSecurePrefix + "([\\w-/]+)\\s*}}")
var SecureParameterPlaceholder = regexp.MustCompile("{{\\s*" + SSMSecurePrefix + "([\\w-/]+)\\s*}}")

//
// Map of supported document formats.
// The resolved values in the output document will be encoded in accordance with the selected format.
var formatMap = map[string]bool {
	TextOutputFormat : true,
	XmlOutputFormat : true,
	YmlOutputFormat : true,
	JsonOutputFormat: true,
}

//
// Map of resolver options.
var optionMap = map[string]bool {
	FailOnParameterNotFoundOption : true,
	IgnoreParameterNotFoundOption : true,
}

//
// Testing related stuff
const CheckMark = "\u2713"
const BallotX = "\u2717"
