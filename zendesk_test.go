package zendesk
import (
  "fmt"
  //"io/ioutil"
  "net"
  "github.com/viki-org/gspec"
  "net/http/httptest"
  "net/http"
  "testing"
)

var (
  testClient *ZendeskClient
)

func init() {
  testClient = &ZendeskClient{
    zendeskUrl : "https://viki.zendesk.com/api/v2",
    zendeskUsername: "mariliam@viki.com",
    zendeskToken: "abcd",
    HttpClient : &http.Client{},
  }
  testClient.organizationCodes = make(map[string]string)
  testClient.organizationCodes["QCs"] =  "4375387268"
  testClient.organizationCodes["Subcribers"] =  "4322274827"
}

func TestValidateUserAndToken(t *testing.T) {
  spec := gspec.New(t)

  isValid := testClient.validateUserAndToken()
  spec.Expect(isValid).ToEqual(true)
}

func TestUpdateAsQC(t *testing.T) {
  spec := gspec.New(t)

  //create test server
  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Println(err)
  }
  handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      fmt.Fprintln(w, `{"users":[{"id":1234}]}`)
    }
    if r.Method == "PUT" {
      fmt.Fprintln(w, `{"user":{"organization_id":`+ testClient.organizationCodes["QCs"]+ `}}`)
    }
  })
  ts := &httptest.Server{
    Listener: l,
    Config: &http.Server{Handler: handler},
  }
  ts.Start()
  defer ts.Close()

  testClient.zendeskUrl = ts.URL
  err = testClient.UpdateAsQC(true, "george@viki.com")
  spec.Expect(err).ToBeNil()
}

func TestUpdateAsSubcriber(t *testing.T) {
  spec := gspec.New(t)

  //create test server
  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Println(err)
  }
  handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      fmt.Fprintln(w, `{"users":[{"id":1234}]}`)
    }
    if r.Method == "PUT" {
      fmt.Fprintln(w, `{"user":{"organization_id":`+ testClient.organizationCodes["Subcribers"]+ `}}`)
    }
  })
  ts := &httptest.Server{
    Listener: l,
    Config: &http.Server{Handler: handler},
  }
  ts.Start()
  defer ts.Close()

  testClient.zendeskUrl = ts.URL
  err = testClient.UpdateAsSubcriber(true, "george@viki.com")
  spec.Expect(err).ToBeNil()
}

func TestUpdateOrganization(t *testing.T) {
  spec := gspec.New(t)

  //create test server
  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Println(err)
  }
  handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
      fmt.Fprintln(w, `{"users":[{"id":1234}]}`)
    }
    if r.Method == "PUT" {
      fmt.Fprintln(w, `{"user":{"organization_id":""}}`)
    }
  })
  ts := &httptest.Server{
    Listener: l,
    Config: &http.Server{Handler: handler},
  }
  ts.Start()
  defer ts.Close()

  testClient.zendeskUrl = ts.URL
  err = testClient.updateOrganization("QCs", false, "george@viki.com")
  spec.Expect(err).ToBeNil()
}

func TestGetIdByEmail(t *testing.T) {
  spec := gspec.New(t)

  //create test server
  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Println(err)
  }
  handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"users":[{"id":1234}]}`)
  })
  ts := &httptest.Server{
    Listener: l,
    Config: &http.Server{Handler: handler},
  }
  ts.Start()
  defer ts.Close()

  testClient.zendeskUrl = ts.URL
  id, err := testClient.getIdByEmail("george@viki.com")
  spec.Expect(id).ToEqual("1234")
  spec.Expect(err).ToBeNil()
}

func TestPutOrganization(t *testing.T) {
  spec := gspec.New(t)

  //create test server
  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Println(err)
  }
  handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"user":{"organization_id":`+ testClient.organizationCodes["QCs"]+ `}}`)
  })
  ts := &httptest.Server{
    Listener: l,
    Config: &http.Server{Handler: handler},
  }
  ts.Start()
  defer ts.Close()
  testClient.zendeskUrl = ts.URL
  organization_id, err := testClient.putOrganization(testClient.organizationCodes["QCs"], "1234")
  spec.Expect(organization_id).ToEqual(testClient.organizationCodes["QCs"])
  spec.Expect(err).ToBeNil()
}

func TestRemoveOrganization(t *testing.T) {
  spec := gspec.New(t)

  //create test server
  l, err := net.Listen("tcp", "127.0.0.1:0")
  if err != nil {
    fmt.Println(err)
  }
  handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"user":{"organization_id":""}}`)
  })
  ts := &httptest.Server{
    Listener: l,
    Config: &http.Server{Handler: handler},
  }
  ts.Start()
  defer ts.Close()
  testClient.zendeskUrl = ts.URL
  organization_id, err := testClient.putOrganization("", "1234")
  spec.Expect(organization_id).ToEqual("")
  spec.Expect(err).ToBeNil()
}