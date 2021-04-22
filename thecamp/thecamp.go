// Package thecamp is a thin wrapper around the thecamp.or.kr website.
// Messages are automatically chunked to around 1500 characters.
package thecamp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DateFormat for formatting dates for Search.
const DateFormat = "20060102"

// Endpoint URLs.
const (
	endpointBase   = "https://www.thecamp.or.kr"
	endpointLogin  = endpointBase + "/login/loginA.do"
	endpointSearch = endpointBase + "/main/cafeCreateCheckA.do"
	endpointSend   = endpointBase + "/consolLetter/insertConsolLetterA.do"
)

// Internal HTTP client with a sane timeout.
var httpClient = &http.Client{Timeout: 15 * time.Second}

// Client holds authentication information to make requests to thecamp.
type Client struct {
	CookieIUID  string
	CookieToken string
}

// New creates a new client with the given credentials.
func New(ctx context.Context, user, pass string) (*Client, error) {
	client := Client{}

	form := make(url.Values, 4)
	form.Set("state", "email-login")
	form.Set("autoLoginYn", "N")
	form.Set("userId", user)
	form.Set("userPwd", pass)

	// Not logged in yet, but we just send empty cookies so it's okay
	resp, err := client.request(ctx, endpointLogin, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	for _, c := range resp.Cookies() {
		switch c.Name {
		case "iuid":
			client.CookieIUID = c.Value
		case "Token":
			client.CookieToken = c.Value
		}
	}

	if client.CookieIUID == "" || client.CookieToken == "" {
		return nil, fmt.Errorf("cookies not in response")
	}

	return &client, nil
}

// Utility function to send a request with cookies and Content-Type set correctly.
func (c *Client) request(ctx context.Context, url string, form url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost, url,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{
		Name:  "iuid",
		Value: c.CookieIUID,
	})
	req.AddCookie(&http.Cookie{
		Name:  "Token",
		Value: c.CookieToken,
	})

	return httpClient.Do(req)
}

// Search for the soldier code. birthday and entrance should be formatted like DateFormat.
func (c *Client) Search(ctx context.Context, name, birthday, entrance string) (int64, error) {
	form := make(url.Values, 6)
	form.Set("name", name)
	form.Set("birth", birthday)
	form.Set("enterDate", entrance)
	// TODO: allow modifything these values, refer to existing libraries
	form.Set("trainUnitCd", "20020191700") // 육군훈련소?
	form.Set("grpCd", "0000010001")        // 육군?

	resp, err := c.request(ctx, endpointSearch, form)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		ResultCode string `json:"resultCd"`
		ResultMsg  string `json:"resultMsg"`

		ListResult []struct {
			TraineeMgrSeq int64 `json:"traineeMgrSeq"`
		} `json:"listResult"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return -1, err
	}

	if body.ResultCode != "9999" {
		return -1, fmt.Errorf("invalid result code: %s (%s)", body.ResultCode, body.ResultMsg)
	}
	if len(body.ListResult) == 0 && body.ListResult[0].TraineeMgrSeq == 0 {
		return -1, fmt.Errorf("invalid response")
	}

	return body.ListResult[0].TraineeMgrSeq, nil
}

// Send a message. The code can be obtained through Search. The message will be
// automatically be chunked to around 1500 characters.
func (c *Client) Send(ctx context.Context, code int64, title, contents string) error {
	lf := "<br/>"
	lines := strings.SplitAfter(strings.ReplaceAll(contents, "\n", lf+"\n"), "\n")

	var (
		chunk      []string
		size, page int
	)
	for _, line := range lines {
		// Assumes no one line is more than 1400 bytes
		if len(line)+size < 1500 {
			chunk = append(chunk, line)
			size += len(line)
		} else {
			// Too large, flush
			err := c.sendChunk(ctx, code, title, strings.Join(chunk, ""), page)
			if err != nil {
				return err
			}

			chunk[0] = line
			chunk = chunk[:1]
			size = 0
			page++
		}
	}

	if len(chunk) > 0 {
		// Flush remaning data
		err := c.sendChunk(ctx, code, title, strings.Join(chunk, ""), page)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) sendChunk(ctx context.Context, code int64, title, chunk string, page int) error {
	form := make(url.Values, 5)
	form.Set("traineeMgrSeq", strconv.FormatInt(code, 10))
	form.Set("sympathyLetterSubject", fmt.Sprintf("%s: %d", title, page))
	form.Set("sympathyLetterContent", chunk)
	form.Set("boardDiv", "sympathyLetter")
	form.Set("tempSaveYn", "N")

	resp, err := c.request(ctx, endpointSend, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var body struct {
		ResultCode string `json:"resultCd"`
		ResultMsg  string `json:"resultMsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	if body.ResultCode != "0000" {
		return fmt.Errorf("unexpected result code: %s (%s)", body.ResultCode, body.ResultMsg)
	}

	return nil
}
