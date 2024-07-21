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
