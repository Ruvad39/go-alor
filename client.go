package alor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	libraryName    = "ALOR API GO"
	libraryVersion = "0.0.3"
)

// Endpoints
const (
	apiProdURL   = "https://api.alor.ru"      // Боевой контур
	apiDevURL    = "https://apidev.alor.ru"   // Тестовый контур
	oauthProdURL = "https://oauth.alor.ru"    // Боевой контур авторизации
	oauthDevURL  = "https://oauthdev.alor.ru" // Тестовый контур авторизации
	wssDevURL    = "wss://apidev.alor.ru/ws"  // Тестовый контур wss
	wssProdURL   = "wss://api.alor.ru/ws"     // Боевой контур wss
)

var ErrNotFound = errors.New("404 Not Found")

var logLevel = &slog.LevelVar{} // INFO
var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: logLevel,
})).With(slog.String("package", "go-alor"))

func SetLogger(logger *slog.Logger) {
	log = logger
}

//const (
//	ErrNotFound = "404 Not Found"
//)

// demoServer использовать тестовый или боевой сервер
var demoServer = false

// getAPIEndpoint return the base endpoint of the Rest API according the UseDevelop flag
func getAPIEndpoint() string {
	if demoServer {
		return apiDevURL
	}
	return apiProdURL
}

func getOauthEndPoint() string {
	if demoServer {
		return oauthDevURL
	}
	return oauthProdURL
}

func getWsEndpoint() string {
	if demoServer {
		return wssDevURL
	}
	return wssProdURL
}

// NewClient создание нового клиента
func NewClient(token string) *Client {
	//Log.Info("NewClient")
	return &Client{
		refreshToken: token,
		//Portfolio:    portfolio,
		Exchange:   "MOEX", // по умолчанию работаем с биржей MOEX
		HTTPClient: http.DefaultClient,
		//Logger:     log.New(os.Stderr, "go-alor ", log.LstdFlags),
	}
}

// Client define API client
type Client struct {
	portfolioID     string    // Номер счета по умолчанию (Portfolio)
	fortsPortfolio  string    // Номер счета для работы с ФОРТС
	stockPortfolio  string    // Номер счета для работы с фондовым рынком
	fxPortfolio     string    // Номер счета для работы с валютным рынком
	refreshToken    string    // Refresh токен пользователя
	accessToken     string    // JWT токен для дальнейшей авторизации
	cancelTimeToken time.Time // Время завершения действия JWT токена
	Exchange        string    // С какой биржей работаем по умолчанию
	HTTPClient      *http.Client
	Stream
	//Portfolio       string    // ID портфеля с которым работаем по умолчанию
}

// SetPortfolioID установим номер счета по умолчанию
func (c *Client) SetPortfolioID(portfolio string) {
	c.portfolioID = portfolio
}

// SetFortsPortfolio установим номер счета для работы с рынком фортс
func (c *Client) SetFortsPortfolio(portfolio string) {
	c.fortsPortfolio = portfolio
}

// SetFxAPortfolio установим номер счета для работы с валютным рынком
func (c *Client) SetFxAPortfolio(portfolio string) {
	c.fxPortfolio = portfolio
}

// SetStockPortfolio установим номер счета для работы с фондовым рынком
func (c *Client) SetStockPortfolio(portfolio string) {
	c.stockPortfolio = portfolio
}

// SetLogDebug установим уровень логгирования Debug
func (c *Client) SetLogDebug(debug bool) {
	if debug {
		logLevel.Set(slog.LevelDebug)
		//log.Debug("установлен уровень Debug")
	} else {
		logLevel.Set(slog.LevelInfo)
	}

}
func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	queryString := r.query.Encode()
	//body := &bytes.Buffer{}
	//bodyString := r.form.Encode()
	header := http.Header{}
	if r.header != nil {
		header = r.header.Clone()
	}

	// если запрос нужно делать с авторизацией (по умолчанию)
	if !r.notAuthorization {
		// получим токен авторизации
		_, err = c.GetJWT()
		if err != nil {
			log.Debug("parseRequest GetJWT", "error", err.Error())
			return err
		}
		if c.accessToken != "" {
			header.Set("Authorization", "Bearer "+c.accessToken)
		}
	}

	baseURL := r.baseURL
	// если ранее не заполнили базовый путь (может быть заполнен в GetJWT())
	if baseURL == "" {
		baseURL = getAPIEndpoint()
	}
	// TODO переделать url.Parse + path.Join
	fullURL := fmt.Sprintf("%s%s", baseURL, r.endpoint)
	// только если ранее мы не заполнили полный путь
	if r.fullURL == "" {
		if queryString != "" {
			fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
		}
		r.fullURL = fullURL
	}
	//c.debug("full url: %s, body: %s", r.fullURL, r.body)
	log.Debug("parseRequest", "url", r.fullURL, slog.Any("body", r.body))

	r.header = header
	//r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	//c.debug("request: %#v", req)
	log.Debug("callAPI", slog.Any("request", req))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	//c.debug("response: %#v", res)
	//c.debug("response body: %s", string(data))
	//c.debug("response status code: %d", res.StatusCode)
	//Log.Debug("callAPI", "status code", res.StatusCode, slog.Any("body", res.Body))
	log.Debug("callAPI", "status code", res.StatusCode)
	//c.debug("debug: GET %s -> %d", r.fullURL, res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(APIError)
		apiErr.Status = res.StatusCode
		//c.debug("Error response body: %s", string(data))
		log.Error("callAPI", "Error response body", res.StatusCode, slog.Any("data", data))
		// TODO обработать ошибку StatusNotFound
		if res.StatusCode == http.StatusNotFound {
			return []byte{}, ErrNotFound
		}

		err := json.Unmarshal(data, apiErr)
		if err != nil {
			//c.debug("failed to unmarshal json: %s", e)
			log.Error("callAPI json.Unmarshal", "err", err.Error())
			apiErr.Code = strconv.Itoa(res.StatusCode)
			apiErr.Message = http.StatusText(res.StatusCode)
		}
		return nil, apiErr
		//c.debug("Erorr response body: %s", string(data))

	}
	return data, nil
}

// (debug) вернем текущую версию
func (c *Client) Version() string {
	return libraryVersion
}

// getRequestID Получение уникального кода запроса
// Текущее время в наносекундах, прошедших с 01.01.1970 в UTC
func (c *Client) getRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

type APIError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	OrderNumber string `json:"orderNumber"`
	Status      int    `json:"-"` // статус HTTP ответа
}

func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%s, msg=%s", e.Code, e.Message)
}

func (e APIError) HTTPStatus() int {
	return e.Status
}
