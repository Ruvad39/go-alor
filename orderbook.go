package alor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// /md/v2/orderbooks/{exchange}/{symbol} Получение информации о биржевом стакане
// https://apidev.alor.ru/md/v2/orderbooks/MOEX/LKOH?depth=20&format=Simple

// GetOrderBooks Получение информации о биржевом стакане
func (c *Client) GetOrderBooks(ctx context.Context, symbol string) (OrderBook, error) {
	queryURL, _ := url.Parse("/md/v2/orderbooks")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, symbol)
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	// выставим максимальное
	r.setParam("depth", 20)

	result := OrderBook{}
	data, err := c.callAPI(ctx, r)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil

}

type PriceVolume struct {
	Price  float64 `json:"price"`  // цена
	Volume int64   `json:"volume"` // объем
}

// PriceVolumeSlice Биды  Аски
type PriceVolumeSlice []PriceVolume

func (slice PriceVolumeSlice) Len() int { return len(slice) }

func (p PriceVolume) String() string {
	return fmt.Sprintf("PriceVolume{ Price: %v, Volume: %v }", p.Price, p.Volume)
}

// вернем второй элемент
func (slice PriceVolumeSlice) Second() (PriceVolume, bool) {
	if len(slice) > 1 {
		return slice[1], true
	}
	return PriceVolume{}, false
}

// вернем первый элемент
func (slice PriceVolumeSlice) First() (PriceVolume, bool) {
	if len(slice) > 0 {
		return slice[0], true
	}
	return PriceVolume{}, false
}

// вернем объем стакана
func (slice PriceVolumeSlice) SumDepth() int64 {
	var total int64
	for _, pv := range slice {
		total = total + pv.Volume
	}

	return total
}

func (slice PriceVolumeSlice) Copy() PriceVolumeSlice {
	var s = make(PriceVolumeSlice, len(slice))
	copy(s, slice)
	return s
}

// OrderBook биржевой стакан
type OrderBook struct {
	Bids        PriceVolumeSlice `json:"bids"`         // Биды
	Asks        PriceVolumeSlice `json:"asks"`         // Аски
	MsTimestamp int64            `json:"ms_timestamp"` // Время (UTC) в формате Unix Time Milliseconds
	Existing    bool             `json:"existing"`     // True - для данных из "снепшота", то есть из истории. False - для новых событий

}

func (b *OrderBook) LastTime() time.Time {
	return time.UnixMilli(b.MsTimestamp)
}

func (b *OrderBook) BestBid() (PriceVolume, bool) {
	if len(b.Bids) == 0 {
		return PriceVolume{}, false
	}

	return b.Bids[0], true
}

func (b *OrderBook) BestAsk() (PriceVolume, bool) {
	if len(b.Asks) == 0 {
		return PriceVolume{}, false
	}

	return b.Asks[0], true
}

func (b *OrderBook) String() string {
	sb := strings.Builder{}

	sb.WriteString("BOOK ")
	//sb.WriteString(b.Symbol)
	sb.WriteString("\n")
	sb.WriteString(b.LastTime().Format("2006-01-02T15:04:05-0700"))
	//sb.WriteString(b.LastTime().String())
	sb.WriteString("\n")

	if len(b.Asks) > 0 {
		sb.WriteString("ASKS:\n")
		for i := len(b.Asks) - 1; i >= 0; i-- {
			sb.WriteString("- ASK: ")
			sb.WriteString(b.Asks[i].String())
			sb.WriteString("\n")
		}
	}

	if len(b.Bids) > 0 {
		sb.WriteString("BIDS:\n")
		for _, bid := range b.Bids {
			sb.WriteString("- BID: ")
			sb.WriteString(bid.String())
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
