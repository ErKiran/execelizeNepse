package neweb

import (
	"context"
	"fmt"
	"nepse-backend/nepse"
	"net/http"
)

type test struct {
	Id string `json:"id"`
}

func (n *NewebAPI) GetFloorsheet(stockId, businessDate, randomId string, size int) (*nepse.FloorsheetResponse, error) {
	url := n.buildFloorsheetSlug(stockId, businessDate, size)
	fmt.Println("url", url)

	ok := test{Id: randomId}

	req, err := n.client.NewRequest(http.MethodPost, url, ok)
	if err != nil {
		return nil, err
	}

	res := &nepse.FloorsheetResponse{}

	if _, err := n.client.Do(context.Background(), req, res); err != nil {
		return nil, err
	}
	return res, nil
}