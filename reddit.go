package main
import "fmt"
import "net/http"
import "strings"
import "encoding/base64"
import "io/ioutil"
import "encoding/json"
import "time"
var client_id string = ""
var client_secret string = ""
var user_agent string = "copypasta_webhook"
type error200 string
func (e error200) Error() string{
  return string(e)
}
var e200 error200 = "200"
func auth() (string, error){
  var json_data interface{}
  method := "POST"
  base_url := "https://www.reddit.com/api/v1/access_token"
  body := strings.NewReader("grant_type=client_credentials")
  req,err := http.NewRequest(method, base_url, body)
  if err != nil{
    return "", err
  }
  basic := client_id+":"+client_secret
  basic_enc := base64.StdEncoding.EncodeToString([]byte(basic))
  req.Header.Add("Authorization", "Basic "+basic_enc)
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("User-Agent", user_agent)
  resp,err := http.DefaultClient.Do(req)
  defer resp.Body.Close()
  if err != nil || resp.StatusCode != 200{
    if resp.StatusCode != 200{
      return "", e200
    }
    return "", err
  }
  content, err := ioutil.ReadAll(resp.Body)
  json.Unmarshal(content, &json_data)

  access_token := json_data.(map[string]interface{})["access_token"].(string)

  return access_token, nil

}
func getPosts()([]string, error){
  var json_data interface{}
  method := "GET"
  base_url := "https://oauth.reddit.com/r/copypasta/hot"
  req,err := http.NewRequest(method, base_url, nil)
  bearer, err := auth()
  if err != nil{
    return nil, err
  }
  req.Header.Add("Authorization", "Bearer "+ bearer)
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("User-Agent", user_agent)
  resp,err := http.DefaultClient.Do(req)
  defer resp.Body.Close()
  if err != nil || resp.StatusCode != 200{
    if resp.StatusCode != 200{
      return nil, e200
    }
    return nil, err
  }
  content, err := ioutil.ReadAll(resp.Body)
  json.Unmarshal(content, &json_data)
  data := json_data.(map[string]interface{})["data"].(map[string]interface{})["children"].([]interface{})
  i := 0 //counter
  ret := make([]string, 0)
  for _,v := range data{
    if i == 3{
      break;
    }
    self_text := v.(map[string]interface{})["data"].(map[string]interface{})["selftext"].(string)
    stickied := v.(map[string]interface{})["data"].(map[string]interface{})["stickied"].(bool)
    if !stickied{
      ret = append(ret, self_text)
      i++
    }

  }
  return ret, nil
}
func postToDisc(content string){
  method := "POST"
  base_url := ""
  body := strings.NewReader("content="+content)
  req,_ := http.NewRequest(method, base_url, body)
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  resp,_ := http.DefaultClient.Do(req)
  bytes,_ := ioutil.ReadAll(resp.Body)
  fmt.Printf("%s", bytes)
  resp.Body.Close()
}
func main(){
  posts, err := getPosts()
  if err != nil{
    fmt.Println("Catastrophic Failure")
  }
  for _,v := range posts {
    postToDisc(v)

    time.Sleep(500 * time.Millisecond)
  }
}
