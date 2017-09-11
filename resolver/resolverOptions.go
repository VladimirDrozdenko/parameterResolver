package resolver

import (
	"errors"
)

//
// Supported value encodings
// default: txt
const TextOutputFormat = "txt"
const XmlOutputFormat = "xml"
const YmlOutputFormat = "yml"
const JsonOutputFormat = "json"


//
// Map of supported encodings
// The resolved values in the output document will be encoded in accordance with the selected format.
var supportedEncodings = map[string]bool {
	TextOutputFormat : true,
	XmlOutputFormat : true,
	YmlOutputFormat : true,
	JsonOutputFormat: true,
}


type ResolveOptions struct {
	IgnoreSecureParameters bool
	ValueEncoding string
}


func validateResolveOptions(options *ResolveOptions) error {
	_, ok := supportedEncodings[options.ValueEncoding]
	if !ok {
		return errors.New("ResolveOptions.ValueEncoding=" + options.ValueEncoding + " is invalid.")
	}
	return nil
}