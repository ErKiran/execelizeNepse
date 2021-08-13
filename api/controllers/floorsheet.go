package controllers

import (
	responses "nepse-backend/api/response"
	"nepse-backend/nepse"
	"nepse-backend/nepse/neweb"
	"nepse-backend/utils"
	"net/http"
)

type FloorsheetResult struct {
	Ticker                string             `json:"ticker"`
	BuyerQuantityMap      map[string]int     `json:"buyerQuantityMap"`
	BuyerTurnOverMap      map[string]float64 `json:"buyerTurnOverMap"`
	BuyerAveragePriceMap  map[string]float64 `json:"buyerAveragePriceMap"`
	SellerQuantityMap     map[string]int     `json:"sellerQuantityMap"`
	SellerTurnOverMap     map[string]float64 `json:"sellerTurnOverMap"`
	SellerAveragePriceMap map[string]float64 `json:"sellerAveragePriceMap"`
}

func (s *Server) GetFloorsheet(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	start := r.URL.Query().Get("start")

	randomId := r.URL.Query().Get("id")

	var nepseBeta nepse.NepseInterface

	nepseBeta, err := neweb.Neweb()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	nepseSectors, err := nepseBeta.GetStocks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var id string

	for _, stock := range nepseSectors {
		if ticker == stock.Ticker {
			id = stock.Id
		}
	}

	floorsheetInfo, err := nepseBeta.GetFloorsheet(id, start, randomId, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalElement := floorsheetInfo.Totaltrades

	floorsheetInfoAgg, err := nepseBeta.GetFloorsheet(id, start, randomId, totalElement)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result FloorsheetResult

	result.BuyerQuantityMap = make(map[string]int)
	result.SellerQuantityMap = make(map[string]int)
	result.BuyerTurnOverMap = make(map[string]float64)
	result.SellerTurnOverMap = make(map[string]float64)
	result.BuyerAveragePriceMap = make(map[string]float64)
	result.SellerAveragePriceMap = make(map[string]float64)

	for _, sheetData := range floorsheetInfoAgg.Floorsheets.Content {
		if sheetData.Buyermemberid != "" {
			result.BuyerQuantityMap[sheetData.Buyermemberid] += sheetData.Contractquantity
		}
		if sheetData.Sellermemberid != "" {
			result.SellerQuantityMap[sheetData.Sellermemberid] += sheetData.Contractquantity
		}
		if sheetData.Buyermemberid != "" {
			result.BuyerTurnOverMap[sheetData.Buyermemberid] += sheetData.Contractamount
		}
		if sheetData.Sellermemberid != "" {
			result.SellerTurnOverMap[sheetData.Sellermemberid] += sheetData.Contractamount
		}
	}

	for k, v := range result.BuyerQuantityMap {
		result.BuyerAveragePriceMap[k] = utils.ToFixed(float64(result.BuyerTurnOverMap[k])/float64(v), 2)
	}
	for k, v := range result.SellerQuantityMap {
		result.SellerAveragePriceMap[k] = utils.ToFixed(float64(result.SellerTurnOverMap[k])/float64(v), 2)
	}
	result.Ticker = ticker

	responses.JSON(w, http.StatusOK, result)
}
