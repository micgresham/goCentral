package goCentral 


import (
    "fmt"
    "os"
    "io/ioutil"
    "io"
//    "bytes"
//    "bufio"
//    "mime/multipart"
    "net/http"
    "time"
//    "log"
//    "strings"
    "github.com/buger/jsonparser"
    "github.com/akamensky/argparse"
    "sigs.k8s.io/yaml"
//    "gopkg.in/yaml.v2"
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "crypto/rand"
    "encoding/hex"
    "encoding/base64"
)

// Structure field names need to be capitalized to be exported (public)
type central struct {
    Base_url string `yaml:"base_url"`
    Customer_id string `yaml:"customer_id"`
    Client_id string `yaml:"client_id"`
    Client_secret string `yaml:"client_secret"`
    Token string `yaml:"token"`
    Refresh_token string `yaml:"refresh_token"`
}


func createHash(key string) string {
        hasher := md5.New()
        hasher.Write([]byte(key))
        return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
        block, _ := aes.NewCipher([]byte(createHash(passphrase)))
        gcm, err := cipher.NewGCM(block)
        if err != nil {
                panic(err.Error())
        }
        nonce := make([]byte, gcm.NonceSize())
        if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
                panic(err.Error())
        }
        ciphertext := gcm.Seal(nonce, nonce, data, nil)
        return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
        key := []byte(createHash(passphrase))
        block, err := aes.NewCipher(key)
        if err != nil {
                panic(err.Error())
        }
        gcm, err := cipher.NewGCM(block)
        if err != nil {
                panic(err.Error())
        }
        nonceSize := gcm.NonceSize()
        nonce, ciphertext := data[:nonceSize], data[nonceSize:]
        plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
        if err != nil {
                panic(err.Error())
        }
        return plaintext
}

//--------------------------------------------------------

func Test_central(central_info central) (int, string, string) {

  token := central_info.Token
  base_url := central_info.Base_url
  client_id := central_info.Client_id
  client_secret := central_info.Client_secret
  refresh_token := central_info.Refresh_token
  api_function_url := fmt.Sprintf("%s/configuration/v2/groups",base_url)
  oath2_url := fmt.Sprintf("%s/oauth2/token",base_url)

  c := http.Client{Timeout: time.Duration(10) * time.Second}
  req, err := http.NewRequest("GET", api_function_url, nil)
  if err != nil {
      fmt.Printf("error %s", err)
      return 500,"",""
  }
  q := req.URL.Query()
  q.Add("limit","1")
  q.Add("offset","0")
  req.URL.RawQuery = q.Encode()

  req.Header.Add("Content-Type", `application/json`)
  req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(token)))
  req.Header.Add("limit","1")
  resp, err := c.Do(req)
  if err != nil {
      fmt.Printf("error %s", err)
      return 500,"",""
  }
  

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  fmt.Printf("%s",body)
  fmt.Printf("**************\n")
  if resp.StatusCode == 401 {
    fmt.Println("ACCESS TOKEN is INVALID or EXPIRED.  Refreshing tokens...")

    c := http.Client{Timeout: time.Duration(10) * time.Second}
    req, err := http.NewRequest("POST", oath2_url, nil)
    if err != nil {
        fmt.Printf("error %s", err)
        return 500,"",""
    }
    q := req.URL.Query()
    q.Add("grant_type","refresh_token")
    q.Add("client_id",client_id)
    q.Add("client_secret",client_secret)
    q.Add("refresh_token",refresh_token)
    req.URL.RawQuery = q.Encode()

    req.Header.Add("Content-Type", `application/json`)
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s",fmt.Sprintf(token)))
    req.Header.Add("limit","1")
    resp2, err := c.Do(req)
    if err != nil {
        fmt.Printf("error %s", err)
        return 500,"",""
    }

  defer resp2.Body.Close()
  body, err := ioutil.ReadAll(resp2.Body)
  fmt.Printf("%s",body)
  fmt.Printf("**************\n")
  refresh_token, err := jsonparser.GetString(body, "refresh_token")
  if err != nil {
     fmt.Printf("error %s", err)
     return 500,"",""
   }
  token, err := jsonparser.GetString(body, "access_token")
  if err != nil {
     fmt.Printf("error %s", err)
     return 500,"",""
   }


  return resp2.StatusCode,token,refresh_token

  }

 return resp.StatusCode,token,refresh_token

}


func Read_DB() central {

     var central_info central
     var tmp_byte []byte

     filename:= "CTconfig.yml"

     yamlFile, err := ioutil.ReadFile(filename)
     if err != nil {
         fmt.Printf("yamlFile.Get err   #%v ", err)
     }
     err = yaml.Unmarshal(yamlFile, &central_info)
     if err != nil {
         fmt.Printf("Unmarshal: %v", err)
     }

     tmp_byte, err = base64.StdEncoding.DecodeString(central_info.Base_url)
     central_info.Base_url = string(decrypt(tmp_byte, passphrase))

     tmp_byte, err = base64.StdEncoding.DecodeString(central_info.Customer_id)
     central_info.Customer_id = string(decrypt(tmp_byte, passphrase))

     tmp_byte, err = base64.StdEncoding.DecodeString(central_info.Client_id)
     central_info.Client_id = string(decrypt(tmp_byte, passphrase))

     tmp_byte, err = base64.StdEncoding.DecodeString(central_info.Client_secret)
     central_info.Client_secret = string(decrypt(tmp_byte, passphrase))

     tmp_byte, err = base64.StdEncoding.DecodeString(central_info.Token)
     central_info.Token = string(decrypt(tmp_byte, passphrase))

     tmp_byte, err = base64.StdEncoding.DecodeString(central_info.Refresh_token)
     central_info.Refresh_token = string(decrypt(tmp_byte, passphrase))

     return(central_info)
}

func Write_DB(central_info_global central) int {

     var central_info central

     // Now encrypt it all into the structure
     central_info.Base_url = string(base64.StdEncoding.EncodeToString(encrypt([]byte(central_info_global.Base_url), passphrase)))
     central_info.Customer_id = string(base64.StdEncoding.EncodeToString(encrypt([]byte(central_info_global.Customer_id), passphrase)))
     central_info.Client_id = string(base64.StdEncoding.EncodeToString(encrypt([]byte(central_info_global.Client_id), passphrase)))
     central_info.Client_secret = string(base64.StdEncoding.EncodeToString(encrypt([]byte(central_info_global.Client_secret), passphrase)))
     central_info.Token = string(base64.StdEncoding.EncodeToString(encrypt([]byte(central_info_global.Token), passphrase)))
     central_info.Refresh_token = string(base64.StdEncoding.EncodeToString(encrypt([]byte(central_info_global.Refresh_token), passphrase)))

     yaml_vars, err := yaml.Marshal(&central_info)
     if err != nil {
       fmt.Printf("err: %v\n", err)
       return(1)
     }

     filename:= "CTconfig.yml"
      _, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
      err = ioutil.WriteFile(filename, yaml_vars, 0644)

     return(0)
}

func Init_DB() {

     var base_url string
     var customer_id string
     var client_id string
     var client_secret string
     var token string
     var refresh_token string
     var central_info central

     fmt.Println("Welcome to the database initialization")
     fmt.Println("")
     fmt.Print("Provide the Central API URL: ")
     fmt.Scanln(&base_url)
     fmt.Print("Provide the Central customer ID: ")
     fmt.Scanln(&customer_id)
     fmt.Print("Provide the Central client ID: ")
     fmt.Scanln(&client_id)
     fmt.Print("Provide the Central secret: ")
     fmt.Scanln(&client_secret)
     fmt.Print("Provide the Central token: ")
     fmt.Scanln(&token)
     fmt.Print("Provide the Central refresh token: ")
     fmt.Scanln(&refresh_token)

     central_info.Base_url = base_url
     central_info.Customer_id = customer_id 
     central_info.Client_id = client_id
     central_info.Client_secret = client_secret
     central_info.Token = token 
     central_info.Refresh_token = refresh_token 

     write_DB(central_info)

}
