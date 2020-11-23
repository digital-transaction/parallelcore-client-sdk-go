//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseTokenAndExpireTimestamp(s string) (token string, expireTimestamp int64, err error) {
	list := strings.Split(s, " ")
	if len(list) != 2 {
		err = fmt.Errorf("Wrong format message from RenewToken().")
	}
	if err == nil {
		token = list[0]
		expireTimestamp, err = strconv.ParseInt(list[1], 10, 64)
		if err != nil {
			err = fmt.Errorf("Wrong format of expireTimestamp. %w", err)
		}
	}
	return token, expireTimestamp, err
}

func formatTokenTimestamp(token string, expireTimestamp int64) string {
	return fmt.Sprintf("Token{len:%d, expire{%d, %q}}", len(token), expireTimestamp, time.Unix(expireTimestamp, 0).String())
}
