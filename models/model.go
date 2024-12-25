package models

import (
	"encoding/json"
	"fmt"
)

type ServerData struct {
	IpAddr []string `json:"IP_Addr"`
}

func (s *ServerData) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("error unmarshalling ServerData: %w", err)
	}
	return nil
}
