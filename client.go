package alor

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"net/http"
	"os"
    "time"
    "net/url"
    "path"
    "strconv"
)

const (
    libraryName    = "ALOR API GO"
    libraryVersion = "0.0.1"
    devServer      = "https://apidev.alor.ru" //Тестовый контур
    prodServer     = "https://api.alor.ru"    //Боевой контур

)


type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

type Client struct {
    token       string 
    baseURL     string
    exchange    string  // с какой биржей работаем

    UserAgent   string    // если проставлен, пропишем User agent в http запросе
    httpClient  HTTPClient //*http.Client
    Logger      *slog.Logger

}


// создание клиента
func NewClient( opts ...ClientOption) (*Client, error) {
    c := &Client{
        baseURL:     prodServer,
        exchange:    "MOEX",            // по умолчанию работет с биржей MOEX
        httpClient:  http.DefaultClient,
        Logger :     slog.New(slog.NewTextHandler(os.Stdout, nil)), //io.Discard
    }
    // обрабратаем входящие параметры
    for _, opt := range opts {
        opt(c)
    }

    return c, nil
}


// выполним запрос  (вернем http.Response)
func (client *Client) RequestHttp(ctx context.Context, httpMethod string, url string, body interface{})(*http.Response, error){

    //client.Logger.Debug("RequestHttp", slog.Any("body", body))
    buf := new(bytes.Buffer)
    if body != nil {
        json.NewEncoder(buf).Encode(body)
        client.Logger.Debug("RequestHttp", slog.Any("buf", buf))
    } 
    
    req, err := http.NewRequestWithContext(ctx, httpMethod, url, buf)         
    if err != nil {
        client.Logger.Error("RequestHttp", "httpMethod", httpMethod, "url", url, "err", err.Error())
        return nil, err
    }     

    // если есть токен доступа = добавим его в заголовок
    if client.token != ""{
        bearer := "Bearer " + client.token 
        req.Header.Add("Authorization", bearer)
    }

    // добавляем заголовки
    if client.UserAgent != "" {
        req.Header.Set("User-Agent", client.UserAgent)
    }    
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    if client.token != ""{
        req.Header.Add("X-Api-Key", client.token)
    }
    
    resp, err := client.httpClient.Do(req)
    if err != nil {
        client.Logger.Error("RequestHttp", "httpMethod", httpMethod, "url", url, "err", err.Error())
        return nil, err
    }
    
    client.Logger.Debug("RequestHttp", "httpMethod", httpMethod, "url", url, "StatusCode", resp.StatusCode)

    return resp, err
}

// выполним запрос  (вернем []byte)
func (client *Client) GetHttp(ctx context.Context, httpMethod string, url string, body interface{})([]byte, error){
    resp, err := client.RequestHttp(ctx, httpMethod, url, body )

    if err != nil {
        client.Logger.Error("RequestHttp", "httpMethod", httpMethod, "url", url, "err", err.Error())
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        client.Logger.Error("RequestHttp", slog.Any("resp",resp))
    }

    defer resp.Body.Close()
    return ioutil.ReadAll(resp.Body)

}

func (client *Client) GetTime(ctx context.Context) (time.Time, error){
    endPoint := "/md/v2/time"
    url, err := url.Parse(client.baseURL)
    if err != nil {
        client.Logger.Error("ошибка разбора baseURL", "err", err.Error())
        return time.Now(), err
    }    

    url.Path = path.Join(url.Path, endPoint)

    resp, err := client.GetHttp(ctx,"GET", url.String(), nil)
    if err != nil {
        client.Logger.Error("GetSysTime GetHttp", "err", err.Error())
        return time.Now(), err
    }    
    tt, _:= strconv.ParseInt(string(resp), 10, 64)
    if err != nil {
        client.Logger.Error("GetSysTime ParseInt", "err", err.Error())
        return time.Now(), err
    }
    servTime := time.Unix(tt, 0) 

    return servTime, nil


}


// (debug) вернем текущую версию
func (c *Client) Version() string{
    return libraryVersion
}



// входящие параметры для создания клиента
type ClientOption func(c *Client)

// WithLogger задает логгер 
// По умолчанию логирование включено на ошибки
func WithLogger(logger *slog.Logger) ClientOption {
    return func(opts *Client) {
        opts.Logger = logger
    }
}

// установим свой HttpClient
// по умолчанию стоит http.DefaultClient
func WithGttpClient(client HTTPClient) ClientOption {
    return func(opts *Client) {
        opts.httpClient = client
    }
}

// url сервера
// по умолчанию стоит Боевой контур ("https://api.alor.ru")
func WithServer(params string) ClientOption {
    return func(opts *Client) {
        opts.baseURL = params
    }
}

// с какой биржей работаем
// MOEX - Московская биржа (стоит по умолчанию)
// SPBX - СПБ Биржа
func WithExchange(params string) ClientOption {
    return func(opts *Client) {
        opts.exchange = params
    }
}

