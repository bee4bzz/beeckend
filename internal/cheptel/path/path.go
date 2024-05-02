package path

const (
	CheptelsGroup = "/cheptels"
	CheptelParam  = ":cheptel-id"

	Create = CheptelsGroup
	Query  = CheptelsGroup
	Update = CheptelsGroup
)

var (
	Delete = MakeDeletePath(CheptelParam)
)

func MakeDeletePath(cheptelID string) string {
	return CheptelsGroup + "/" + cheptelID
}
