package zendesk
import (
  "encoding/json"
  "errors"
  "net/http"
  bytepool "github.com/viki-org/bytepool"
  "strconv"
)

type ZendeskClient struct {
  zendeskUrl string
  zendeskUsername string
  zendeskToken string
  organizationCodes map[string]string
  HttpClient *http.Client
}

func (client *ZendeskClient) SetUsername(username string) {
  client.zendeskUsername = username + "/token"
}

func (client *ZendeskClient) GetUsername() string { 
  return client.zendeskUsername 
}

func (client *ZendeskClient) SetToken(token string) {
  client.zendeskToken = token
}

func (client * ZendeskClient) GetToken() string{
  return client.zendeskToken
}

func (client *ZendeskClient) UpdateAsQC(isQC bool, email string) (error){
  return client.updateOrganization("QCs", isQC, email)
}

func (client *ZendeskClient) UpdateAsSubcriber(isSubcriber bool, email string) (error){
  return client.updateOrganization("Subcribers", isSubcriber, email)
}

func (client *ZendeskClient) updateOrganization(organization string, belongTo bool, email string) (error){
  id, err := client.getIdByEmail(email)
  if err != nil {
    return err
  }

  organizationCode := ""
  if belongTo {
    organizationCode = client.organizationCodes[organization]
  }
  o_id, r := client.putOrganization(organizationCode, id)
  if r == nil {
    if o_id == organizationCode {
      return nil
    }
    return errors.New("organization_id update mismatch")
  }
  return r
}

func (client *ZendeskClient) putOrganization(organizationCode string, id string) (string, error){
  url := client.zendeskUrl + "/users/" + id + ".json"
  sentJson := `{"user": {"organization_id":` + organizationCode + `}}`
  sentBody :=  pool.Checkout()
  sentBody.WriteString(sentJson)
  req, _ := http.NewRequest("PUT", url, sentBody)
  req.Header.Set("Content-Type", "application/json")
  req.SetBasicAuth(client.zendeskUsername, client.zendeskToken)
  resp, err:= client.HttpClient.Do(req)
  if err != nil {
    return "",err
  }

  defer resp.Body.Close()
  body := pool.Checkout()
  body.ReadFrom(resp.Body)

  var result interface{}
  err = json.Unmarshal(body.Bytes(), &result)
  if err != nil {
    return "", err
  }
  m := result.(map[string]interface{})
  var user map[string]interface{}
  user = m["user"].(map[string]interface{})
  if user["organization_id"] == "" {
    return "", nil
  }
  return strconv.FormatInt(int64(user["organization_id"].(float64)), 10), nil
}

func (client *ZendeskClient) getIdByEmail(email string) (string, error){
  url := client.zendeskUrl + "/users/search.json";
  req, _ := http.NewRequest("GET", url, nil)
  q := req.URL.Query()
  q.Add("query",email)
  req.URL.RawQuery = q.Encode()
  req.SetBasicAuth(client.zendeskUsername, client.zendeskToken)

  resp, err:= client.HttpClient.Do(req)
  if err != nil {
    return "", err
  }

  defer resp.Body.Close()
  body := pool.Checkout()
  body.ReadFrom(resp.Body)

  var result interface{}
  err = json.Unmarshal(body.Bytes(), &result)
  if err != nil {
    return "", err
  }
  m := result.(map[string]interface{})
  var users []interface{}
  users = m["users"].([]interface{})
  var user map[string]interface{}
  user = users[0].(map[string]interface{})
  return strconv.FormatInt(int64(user["id"].(float64)), 10), nil
}

var (
  pool = bytepool.New(100, 102400)
  Client *ZendeskClient
)

func init() {
  Client = &ZendeskClient{
    zendeskUrl : "https://viki.zendesk.com/api/v2",
    HttpClient : &http.Client{},
  }
  Client.organizationCodes = make(map[string]string)
  Client.organizationCodes["QCs"] =  "4375387268"
  Client.organizationCodes["Subcribers"] =  "4322274827"
}


