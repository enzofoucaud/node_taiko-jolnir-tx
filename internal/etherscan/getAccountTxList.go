package etherscan

import (
	"encoding/json"
	"net/http"
)

type AccountTxList struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		TimeStamp string `json:"timeStamp"`
		Hash      string `json:"hash"`
		To        string `json:"to"`
	} `json:"result"`
}

func GetAccountTxList(url, api, address string) (AccountTxList, error) {
	resp, err := http.Get(url + "api?module=account&action=txlist&page=1&offset=25&sort=desc&address=" + address + "&tag=latest&apikey=" + api)
	if err != nil {
		return AccountTxList{}, err
	}
	defer resp.Body.Close()

	var accountTxList AccountTxList

	err = json.NewDecoder(resp.Body).Decode(&accountTxList)
	if err != nil {
		return AccountTxList{}, err
	}

	return accountTxList, nil
}
