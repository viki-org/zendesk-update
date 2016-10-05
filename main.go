package main 
import (
  "encoding/json"
  "fmt"
  //"io/ioutil"
  "net/http"
  "net/url"
  "github.com/viki-org/bytepool"
  "strconv"
)

type CookieStorage struct {
  cookieMap map[string] []*http.Cookie
}

func (c *CookieStorage) SetCookies(u *url.URL, cookies []*http.Cookie) {
  c.cookieMap [u.Host] = cookies
}

func (c *CookieStorage) Cookies(u *url.URL) []*http.Cookie {
  return c.cookieMap[u.Host]
}

// type ZendeskUser struct {
//   Id    string 
//   Url   string
//   Name  string
//   ExternalId string `json:"external_id"`
//   Alias      string
//   CreatedAt string `json:"created_at"`
//   UpdatedAt string `json:"updated_at"`
//   TimeZone  string `json:"time_zone"`
//   LastLoginAt string `json:"last_login_at"`
//   Email     string 
//   Phone     string
//   Signature string
//   Details   string
//   Notes     string
//   Locale    string
//   LocaleId  int `json:"locale_id"`
//   OrganizationId string `json:"organization_id"`
//   Role           string
//   CustomeRoleId  string `json:"custom_role_id"`
//   Moderator      bool
//   TicketRestriction string `json:"ticket_restriction"`
//   OnlyPrivateComments bool `json:"false"`
//   Active         bool
//   Shared         bool
//   SharedAgent    bool
//   Verified       bool
//   Photo          interface{}
// }

// type SearchResult struct {
//   Users    []ZendeskUser `json:"users"`
//   NextPage string `json:"next_page"`
//   PrevPage string `json:"previous_page"`
//   Count    int `json:"count"`
// }
var (
  pool = bytepool.New(100, 102400)
)

type ZendeskClient struct {
  zendeskUrl string
  zendeskUsername string
  zendeskToken string
  Cookie map[string] []*http.Cookie
  organization = map[string]string{"QCs": "4375387268", "Subcribers":"4322274827"}
  HttpClient &http.Client
}

func (client *ZendeskClient) SetUrl(url string) {
  client.zendeskUrl = url
}

func (client *ZendeskClient) SetUsername(username string) {
  client.zendeskUsername = username + "/token"
}

func (client *ZendeskClient) SetToken(token string) {
  client.zendeskToken = token
}

func (client *ZendeskClient) UpdateAsQC(isQC bool, email string) {
  client.updateOrganization("QCs", isQc)
}

func (client *ZendeskClient) UpdateAsSubcriber(isSubcriber bool, email string) {
  client.updateOrganization("Subcribers", isSubcriber)
}

func (client *ZendeskClient) updateOrganization(organization string, belongTo bool, email string) {
  client.HttpClient := &http.Client{}
  cookieStorage := &CookieStorage{}
  cookieStorage.cookieMap = make(map[string] []*http.Cookie)
  client.HttpClient.Jar = cookieStorage

}

func (client *ZendeskClient) getIdByEmail(email string) (<-chan string, error){
  c := make(chan string)
  go func(){
    req, _ := http.NewRequest("GET", zendeskUrl + "/users/search.json", nil)
    q := req.URL.Query()
    q.Add("query",email)
    req.URL.RawQuery = q.Encode()
    req.SetBasicAuth(zendeskUsername + "/token", zendeskToken)

    resp, err := client.Do(req)
    if err != nil {
      fmt.Println(err)
      c <- nil
      return c, err
    }

    defer resp.Body.Close()
    body := pool.Checkout()
    body.ReadFrom(resp.Body)
    var result interface{}
    err = json.Unmarshal(body.Bytes(), &result)
    if err != nil {
      fmt.Println(err)
      c <- nil
      return c, err
    }

    m := result.(map[string]interface{})
    var users []interface{}
    users = m["users"].([]interface{})
    var user map[string]interface{}
    user = users[0].(map[string]interface{})
    c <- strconv.FormatInt(int64(user["id"].(float64)), 10)
  }()
  return c, nil
}

func main() {
  zenddeskClient := &ZendeskClient{}
}