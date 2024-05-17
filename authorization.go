package alor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

/*
	Механизм refresh token позволяет получать JWT с новым сроком жизни.
	Для этого отправьте POST запрос на адрес
	https://oauthdev.alor.ru/refresh?token={refreshToken} (тестовый контур)
	или
	https://oauth.alor.ru/refresh?token={refreshToken} (боевой контур).
	Если у refresh token не истек срок жизни и он не был отозван,
	то в теле ответа в поле AccessToken вернётся свежий JWT токен.

*/

const jwtTokenTtl = 60 // Время жизни токена JWT в секундах

type JSResp struct {
	AccessToken string
}

// GetJWT получим accessToken
// TODO определять время работы токена
func (c *Client) GetJWT() error {
	if c.refreshToken == "" {
		c.accessToken = ""
		return nil
	}
	// если не пустой токен и время окончания токена больше текущего время
	if c.accessToken != "" && c.cancelTimeToken.After(time.Now()) {
		//c.debug("GetJWT Не надо формировать новый токен")
		return nil

	}
	//c.debug("GetJWT Формируем новый токен")

	queryURL, _ := url.Parse(c.OauthURL)
	queryURL.Path = path.Join(queryURL.Path, "refresh")

	q := queryURL.Query()
	q.Set("token", c.refreshToken)
	// добавляем к URL параметры
	queryURL.RawQuery = q.Encode()

	//c.debug("full url: %s", queryURL.String())
	req, err := http.NewRequest(http.MethodPost, queryURL.String(), nil)
	if err != nil {
		return err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка получения JWT токена: статус %d", res.StatusCode)

	}
	//defer res.Body.Close()
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	data, err := io.ReadAll(res.Body)
	//data, err := c.callAPI(ctx, r)
	if err != nil {
		c.accessToken = ""
		return err
	}
	var result JSResp
	err = json.Unmarshal(data, &result)
	if err != nil {
		//c.debug("error  %s", err.Error())
		return err
	}
	c.cancelTimeToken = time.Now().Add(jwtTokenTtl * time.Second)
	c.accessToken = result.AccessToken

	return nil

}
