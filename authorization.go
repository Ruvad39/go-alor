package alor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
//  Срок действия access токена составляет 30 минут. Ограничим 25 минут
const jwtTokenTtl = 25 // Время жизни токена JWT в минутах

type JSResp struct {
	AccessToken string
}

// GetJWT получим accessToken
func (c *Client) GetJWT() (string, error) {
	//log.Debug("зашли в НОВУЮ GetJWT")
	if c.refreshToken == "" {
		c.accessToken = ""
		return c.accessToken, nil
	}
	// если не пустой токен и время окончания токена больше текущего время
	if c.accessToken != "" && c.cancelTimeToken.After(time.Now()) {
		//c.debug("GetJWT Не надо формировать новый токен")
		return c.accessToken, nil

	}
	r := &request{
		method:           http.MethodPost,
		endpoint:         "/refresh",
		notAuthorization: true, // Проставить обязательно. Иначе будет ошибка
	}
	r.baseURL = getOauthEndPoint()
	r.setParam("token", c.refreshToken)

	var result JSResp
	data, err := c.callAPI(context.Background(), r)
	if err != nil {
		return "", fmt.Errorf("ошибка получения JWT токена: %w", err)
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("ошибка получения JWT токена: %w", err)
	}
	c.cancelTimeToken = time.Now().Add(jwtTokenTtl * time.Minute)
	c.accessToken = result.AccessToken

	return c.accessToken, nil

}

/*
// GetJWT получим accessToken
func (c *Client) GetJWT_OLD() (string, error) {
	if c.refreshToken == "" {
		c.accessToken = ""
		return c.accessToken, nil
	}
	// если не пустой токен и время окончания токена больше текущего время
	if c.accessToken != "" && c.cancelTimeToken.After(time.Now()) {
		//c.debug("GetJWT Не надо формировать новый токен")
		return c.accessToken, nil

	}
	queryURL, _ := url.Parse(getOauthEndPoint())
	queryURL.Path = path.Join(queryURL.Path, "refresh")

	q := queryURL.Query()
	q.Set("token", c.refreshToken)
	// добавляем к URL параметры
	queryURL.RawQuery = q.Encode()

	//c.debug("full url: %s", queryURL.String())
	req, err := http.NewRequest(http.MethodPost, queryURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("ошибка получения JWT токена: %w", err)
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка получения JWT токена: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка получения JWT токена: статус %d", res.StatusCode)
	}

	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		c.accessToken = ""
		return c.accessToken, fmt.Errorf("ошибка получения JWT токена: %w", err)
	}
	var result JSResp
	err = json.Unmarshal(data, &result)
	if err != nil {
		//c.debug("error  %s", err.Error())
		return c.accessToken, fmt.Errorf("ошибка получения JWT токена: %w", err)
	}
	c.cancelTimeToken = time.Now().Add(jwtTokenTtl * time.Minute)
	c.accessToken = result.AccessToken

	return c.accessToken, nil

}
*/
