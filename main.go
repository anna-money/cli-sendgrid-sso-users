package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"

	"github.com/sendgrid/sendgrid-go"
	logger "github.com/sirupsen/logrus"
)

const (
	sendgridAPIHost = "https://api.sendgrid.com"
)

var (
	updateUsersFlag, createUsersFlag, allUsersFlag, usersWithoutSSOflag bool
	configPath, sendgridToken                                           string
)

type (
	YamlConfig struct {
		Groups []struct {
			Admin struct {
				Users []struct {
					FirstName string `yaml:"first_name"`
					LastName  string `yaml:"last_name"`
					Email     string `yaml:"email"`
				} `yaml:"users"`
				Scopes []string `yaml:"scopes"`
			} `yaml:"admin,omitempty"`
			Developer struct {
				Users []struct {
					FirstName string `yaml:"first_name"`
					LastName  string `yaml:"last_name"`
					Email     string `yaml:"email"`
				} `yaml:"users"`
				Scopes []string `yaml:"scopes"`
			} `yaml:"developer,omitempty"`
			Support struct {
				Users []struct {
					FirstName string `yaml:"first_name"`
					LastName  string `yaml:"last_name"`
					Email     string `yaml:"email"`
				} `yaml:"users"`
				Scopes []string `yaml:"scopes"`
			} `yaml:"support,omitempty"`
		} `yaml:"groups"`
	}
	ResponseCreateUserSSO struct {
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Email     string   `json:"email"`
		IsAdmin   bool     `json:"is_admin"`
		Scopes    []string `json:"scopes"`
	}
	ResponseGetAllTeammatesWithoutSSO struct {
		Result []struct {
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			IsAdmin   bool   `json:"is_admin"`
			IsSso     bool   `json:"is_sso"`
			LastName  string `json:"last_name"`
			UserType  string `json:"user_type"`
			Username  string `json:"username"`
		} `json:"result"`
	}
	ResponseUpdateUserSSO struct {
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		IsAdmin   bool     `json:"is_admin"`
		Scopes    []string `json:"scopes"`
	}
)

func (c *YamlConfig) getConf() *YamlConfig {
	yamlFile, err := ioutil.ReadFile("config/users.yaml")
	if err != nil {
		logger.Fatalf("get config/users.yaml error: #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		logger.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func CreateTeammatesSSO(firstName string, lastName string, email string, isAdmin bool, scopes []string, group string) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := sendgridAPIHost
	request := sendgrid.GetRequest(apiKey, "/v3/sso/teammates", host)
	request.Method = "POST"
	m := ResponseCreateUserSSO{firstName, lastName, email, isAdmin, scopes}
	b, _ := json.Marshal(m)
	request.Body = b
	response, err := sendgrid.API(request)
	if err != nil {
		logger.Fatalf("error send request to sendgrid api: %v", err)
	}
	if response.StatusCode == 201 {
		logger.Infof("username %s/%s created", group, email)
	} else if response.StatusCode == 400 {
		logger.Infof("username %s/%s already exist", group, email)
	} else {
		logger.Errorf("error statusCode from sendgrid, expected 200, status code: %d", response.StatusCode)
		logger.Errorf("response: %s", response.Body)
	}
}

func UpdateTeammatesSSO(firstName string, lastName string, email string, isAdmin bool, scopes []string, group string) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := sendgridAPIHost
	request := sendgrid.GetRequest(apiKey, "/v3/sso/teammates/"+email, host)
	request.Method = "PATCH"
	m := ResponseUpdateUserSSO{firstName, lastName, isAdmin, scopes}
	b, _ := json.Marshal(m)
	request.Body = b
	response, err := sendgrid.API(request)
	if err != nil {
		logger.Fatalf("error send request to sendgrid api: %v", err)
	}
	if response.StatusCode == 200 {
		logger.Infof("username %s/%s success update", group, email)
	} else {
		logger.Errorf("error statusCode from sendgrid, expected 200, status code: %d", response.StatusCode)
		logger.Errorf("response: %s", response.Body)
	}
}

func GetAllTeammatesWithoutSSO() {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := sendgridAPIHost
	request := sendgrid.GetRequest(apiKey, "/v3/teammates?limit=100", host)
	request.Method = "GET"
	response, err := sendgrid.API(request)
	if err != nil {
		logger.Fatalf("error send request to sendgrid api: %v", err)
	}
	m := []byte(response.Body)
	r := bytes.NewReader(m)
	decoder := json.NewDecoder(r)
	users := &ResponseGetAllTeammatesWithoutSSO{}
	err = decoder.Decode(users)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}
	for _, u := range users.Result {
		if u.IsSso == false {
			logger.Infof("username %s", u.Username)
		}
	}
}

func GetAllTeammates() {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := sendgridAPIHost
	request := sendgrid.GetRequest(apiKey, "/v3/teammates?limit=200", host)
	request.Method = "GET"
	response, err := sendgrid.API(request)
	if err != nil {
		logger.Fatalf("error send request to sendgrid api: %v", err)
	}
	m := []byte(response.Body)
	r := bytes.NewReader(m)
	decoder := json.NewDecoder(r)
	users := &ResponseGetAllTeammatesWithoutSSO{}
	err = decoder.Decode(users)

	if err != nil {
		logger.Fatalf("error: %v", err)
	}
	for _, u := range users.Result {
		logger.Infof("username %s", u.Username)
	}
}

func validateArguments() error {
	if createUsersFlag == false && updateUsersFlag == false && usersWithoutSSOflag == false && allUsersFlag == false {
		return fmt.Errorf("need to choose one or both actions, use flags (--create/-c,--update/-u or --get-all/-a,--get-all-no-sso/-s)")
	}

	if configPath == "" {
		return fmt.Errorf("--config-path should not be empty")
	}

	return nil
}

func main() {
	logger.SetFormatter(&logger.TextFormatter{DisableColors: false})
	logger.SetLevel(logger.InfoLevel)

	pflag.StringVarP(&configPath, "config-path", "", "config/users.yaml", "Config file path")
	pflag.StringVarP(&sendgridToken, "sendgrid-token", "t", os.Getenv("SENDGRID_API_KEY"), "Config file path")
	pflag.BoolVarP(&createUsersFlag, "create", "c", false, "Create all users")
	pflag.BoolVarP(&updateUsersFlag, "update", "u", false, "Update all users")
	pflag.BoolVarP(&allUsersFlag, "get-all", "a", false, "Get all users")
	pflag.BoolVarP(&usersWithoutSSOflag, "get-all-no-sso", "n", false, "Get all no SSO users")
	pflag.Parse()

	if err := validateArguments(); err != nil {
		logger.Fatalf("Validation arguments error: %s", err.Error())
	}

	var c YamlConfig
	c.getConf()

	if usersWithoutSSOflag {
		GetAllTeammatesWithoutSSO()
	}

	if allUsersFlag {
		GetAllTeammates()
	}

	if createUsersFlag {
		for _, g := range c.Groups {
			for _, u := range g.Admin.Users {
				CreateTeammatesSSO(u.FirstName, u.LastName, u.Email, true, []string{}, "admin")
			}
			for _, u := range g.Developer.Users {
				CreateTeammatesSSO(u.FirstName, u.LastName, u.Email, false, g.Developer.Scopes, "developer")
			}
			for _, u := range g.Support.Users {
				CreateTeammatesSSO(u.FirstName, u.LastName, u.Email, false, g.Support.Scopes, "support")
			}
		}
	}

	if updateUsersFlag {
		for _, g := range c.Groups {
			for _, u := range g.Admin.Users {
				UpdateTeammatesSSO(u.FirstName, u.LastName, u.Email, true, []string{}, "admin")
			}
			for _, u := range g.Developer.Users {
				UpdateTeammatesSSO(u.FirstName, u.LastName, u.Email, false, g.Developer.Scopes, "developer")
			}
			for _, u := range g.Support.Users {
				UpdateTeammatesSSO(u.FirstName, u.LastName, u.Email, false, g.Support.Scopes, "support")
			}
		}
	}
}
