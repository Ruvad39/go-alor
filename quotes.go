package alor

import (
    "context"
    "net/url"
    "path"
    "encoding/json"
   
)

// реализован только этот метод. так как только он позволяет делать запрос без авторизации
// пример: https://apidev.alor.ru/md/v2/Securities/MOEX:SiH4/quotes
func (client *Client) GetQuotes(ctx context.Context, symbol string) ([]Quote, error){
	endPoint := "md/v2/Securities"
	endPoint2 := "quotes"
	ticker := client.exchange +":"+symbol // "MOEX:"+ symbol

 	result := make([]Quote, 0)

	url, err := url.Parse(client.baseURL)
    if err != nil {
        client.Logger.Error("ошибка разбора baseURL", "err", err.Error())
        return result, err
    }
    url.Path = path.Join(url.Path, endPoint,ticker,endPoint2)
    //queryURL.Path = path.Join(queryURL.Path, "md/v2/Securities",ticker,"quotes")


 	resp, err := client.GetHttp(ctx,"GET", url.String(), nil)
    // res, err := client.GetHttp(ctx,"GET", queryURL.String())
    if err != nil {
        client.Logger.Error("GetQuotes ", "err", err.Error())
        return result, err
    }    

	err = json.Unmarshal(resp, &result)
	if err != nil {
		client.Logger.Error("GetQuotes ", "err", err.Error())
		return result, err
	}
    return result, nil


}
