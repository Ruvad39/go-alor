package alor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Endpoints
const (
	libraryName    = "ALOR API GO"
	libraryVersion = "0.0.3"
	apiProdURL     = "https://api.alor.ru"    // Боевой контур
	apiDevURL      = "https://apidev.alor.ru" // Тестовый контур
	oauthProdURL   = "https://oauth.alor.ru"
	oauthDevURL    = "https://oauthdev.alor.ru"
)

// UseDevelop использовать тестовый или боевой сервер
var UseDevelop = false

// getAPIEndpoint return the base endpoint of the Rest API according the UseDevelop flag
func getAPIEndpoint() (string, string) {
	if UseDevelop {
		return apiDevURL, oauthDevURL
	}
	return apiProdURL, oauthProdURL
}

// NewClient создание нового клиента
func NewClient(token string) *Client {
	apidURL, oauthURL := getAPIEndpoint()
	return &Client{
		refreshToken: token,
		//Portfolio:    portfolio,
		ApiURL:     apidURL,
		OauthURL:   oauthURL,
		Exchange:   "MOEX", // по умолчанию работаем с биржей MOEX
		UserAgent:  "Alor/golang",
		HTTPClient: http.DefaultClient,
		Logger:     log.New(os.Stderr, "go-alor ", log.LstdFlags),
	}
}

// Client define API client
type Client struct {
	Portfolio       string    // ID портфеля с которым работаем по умолчанию
	refreshToken    string    // Refresh токен пользователя
	accessToken     string    // JWT токен для дальнейшей авторизации
	cancelTimeToken time.Time // Время завершения действия JWT токена
	Exchange        string    // С какой биржей работаем по умолчанию
	ApiURL          string
	OauthURL        string
	UserAgent       string
	HTTPClient      *http.Client
	Debug           bool
	Logger          *log.Logger
	TimeOffset      int64
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}

	err = c.GetJWT()
	if err != nil {
		c.debug("error  %s", err.Error())
		return err
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

	if c.accessToken != "" {
		header.Set("Authorization", "Bearer "+c.accessToken)
	}

	fullURL := fmt.Sprintf("%s%s", c.ApiURL, r.endpoint)
	// только если ранее мы не заполнили полный путь
	if r.fullURL == "" {

		if queryString != "" {
			fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
		}
		r.fullURL = fullURL
	}
	//c.debug("full url: %s, body: %s", r.fullURL, bodyString)
	c.debug("full url: %s, body: %s", r.fullURL, r.body)

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
	c.debug("request: %#v", req)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	//data, err = ioutil.ReadAll(res.Body)
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
	c.debug("response status code: %d", res.StatusCode)
	//c.debug("debug: GET %s -> %d", r.fullURL, res.StatusCode)

	if res.StatusCode >= http.StatusBadRequest {
		apiErr := new(APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
			apiErr.Code = strconv.Itoa(res.StatusCode)
			apiErr.Message = http.StatusText(res.StatusCode)
		}
		return nil, apiErr
		//c.debug("Erorr response body: %s", string(data))
		//return nil, fmt.Errorf("error HTTP %d: %s", res.StatusCode, http.StatusText(res.StatusCode))
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
}

func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%s, msg=%s", e.Code, e.Message)
}
