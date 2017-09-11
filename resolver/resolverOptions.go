package resolver

//
// Supported value encodings
// default: txt
const TextOutputFormat = "txt"
const XmlOutputFormat = "xml"
const YmlOutputFormat = "yml"
const JsonOutputFormat = "json"

type ResolveOptions struct {
	IgnoreSecureParameters bool
	ValueEncoding          string
}
