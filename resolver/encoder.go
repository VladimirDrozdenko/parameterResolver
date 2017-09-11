package resolver

type formatEncoder interface {
	encode(value string) string
}

type textEncoder struct {
}

func (e textEncoder) encode(value string) string {
	return value
}

type xmlEncoder struct {
}

func (e xmlEncoder) encode(value string) string {
	return value
}

type jsonEncoder struct {
}

func (e jsonEncoder) encode(value string) string {
	return value
}

type ymlEncoder struct {
}

func (e ymlEncoder) encode(value string) string {
	return value
}

func NewFormatEncoder(encoding string) formatEncoder {
	switch encoding {
	case XmlOutputFormat:
		return new(xmlEncoder)
	case JsonOutputFormat:
		return new(jsonEncoder)
	case YmlOutputFormat:
		return new(ymlEncoder)
	default:
		// text encoding is default
		return new(textEncoder)
	}
}
