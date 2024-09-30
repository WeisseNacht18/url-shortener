package configValidator

import (
	"errors"
	"strconv"
	"strings"
)

func IsValidServerHost(host string) error {
	splitedHost := strings.Split(host, ":")

	if len(splitedHost) > 2 {
		return errors.New("invalid hostname")
	}

	num, err := strconv.Atoi(splitedHost[1])
	if err != nil {
		return err
	}

	if num < 0 || num > 65536 {
		return errors.New("invalid port")
	}

	return nil
}
