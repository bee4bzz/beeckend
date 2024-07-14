package token

import "fmt"

func AddConfirmQuery(path, tokenID, token string) string {
	return path + fmt.Sprintf("?token=%s&challenge=%s", token, tokenID)
}
