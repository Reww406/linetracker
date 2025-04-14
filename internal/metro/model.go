package metro

type LineCode string

const (
	RedLine    LineCode = "RD"
	OrangeLine LineCode = "OR"
	SilverLine LineCode = "SV"
	BlueLine   LineCode = "BL"
	GreenLine  LineCode = "GR"
)

func ToLineCodes(codes []string) []LineCode {
	result := make([]LineCode, len(codes))
	for i, code := range codes {
		result[i] = LineCode(code)
	}

	return result	
}
