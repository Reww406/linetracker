package metro

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"


type LineCode string

const (
	RedLine    LineCode = "RD"
	OrangeLine LineCode = "OR"
	SilverLine LineCode = "SV"
	BlueLine   LineCode = "BL"
	GreenLine  LineCode = "GR"
)

func ToLineCodes(codes []types.AttributeValue) []LineCode {
	result := make([]LineCode, len(codes))
	for i, code := range codes {
		// TODO What the hell is happening here.
		if sv, ok := code.(*types.AttributeValueMemberS); ok {
			result[i] = LineCode(sv.Value)
		}
	}

	return result
}
