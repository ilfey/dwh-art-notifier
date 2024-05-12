package nekos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetArt(category string) (*Response, error) {
	resp, err := http.Get("https://nekos.best/api/v2/" + category)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(string(responseBody))
	}

	res := new(Response)

	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
