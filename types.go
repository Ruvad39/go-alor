package alor

/*

positions[https://alor.dev/rawdocs2/WarpOpenAPIv2.yml#/components/schemas/positionposition

*/

// Interval период свечей
type Interval string

// Длительность таймфрейма. В качестве значения можно указать точное количество секунд или код таймфрейма
const (
	Interval_S15 Interval = "15"   // 15 секунд
	Interval_M1  Interval = "60"   // 60 секунд или 1 минута
	Interval_H1  Interval = "3600" // 3600 секунд или 1 час
	Interval_D1  Interval = "D"    // D — сутки (соответствует значению 86400)
	Interval_W1  Interval = "W"    // W — неделя (соответствует значению 604800)
	Interval_MN1 Interval = "M"    // M — месяц (соответствует значению 2592000)
	Interval_Y1  Interval = "Y"    // Y — год (соответствует значению 31536000)

)

// направление сдели ( BUY SELL )
type SideType string

const (
	SideTypeBuy  SideType = "buy"
	SideTypeSell SideType = "sell"
)

type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// OrderStatus статус заявки ( New Filled Canceled Rejected )
type OrderStatus string

const (
	OrderStatusWorking  OrderStatus = "working"  // На исполнении
	OrderStatusFilled   OrderStatus = "filled"   // Полностъю исполнилась (выполнилась)
	OrderStatusCanceled OrderStatus = "canceled" // Отменена
	OrderStatusRejected OrderStatus = "rejected" // отклонена

)

// условие по времени действия заявки
type TimeInForce string

var (
	TimeInForceGTC    TimeInForce = "goodtillcancelled" // Активна до отмены
	TimeInForceDAY    TimeInForce = "oneday"            // До конца дня
	TimeInForceFOK    TimeInForce = "fillorkill"        // Исполнить целиком или отклонить
	TimeInForceCancel TimeInForce = "immediateorcancel" // Снять остаток
)
