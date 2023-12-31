package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Query[T any](url string, target T) error {
	fmt.Printf("GET %q\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%q returned %q", url, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
}
