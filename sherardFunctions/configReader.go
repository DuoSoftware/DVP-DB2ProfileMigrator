package sherardFunctions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var dirPath string

type DB struct {
	Type     string `json:"Type"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Port     string `json:"Port"`
	Host     string `json:"Host"`
	Database string `json:"Database"`
}

type Mongo struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Database string `json:"Database"`
	User     string `json:"User"`
	Password string `json:"Password"`
	CollectionName string `json:"CollectionName"`
}

type DBTWO struct {
	Protocol string `json:"Protocol"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Port     string `json:"Port"`
	Host     string `json:"Host"`
	Database string `json:"Database"`
}

type Host struct {
	Domain      string `json:"Domain"`
	Port        string `json:"Port"`
	Version     string `json:"Version"`
	Hostpath    string `json:"Hostpath"`
	Logfilepath string `json:"Logfilepath"`
}

type Security struct {
	Ip          string `json:"Ip"`
	Port        string `json:"Port"`
	AccessToken string `json:"AccessToken"`
}

type MigrationInfo struct {
	SchedulerTime  string `json:"SchedulerTime"`
	Tenant  string `json:"Tenant"`
	Companys  []string `json:"Companys"`
	CompanysName  string `json:"CompanysName"`
}

type Configuration struct {
	DB       DB       `json:"DB"`
	Mongo    Mongo       `json:"Mongo"`
	DBTWO    DBTWO       `json:"DBTWO"`
	Host     Host     `json:"Host"`
	Security Security `json:"Security"`
	MigrationInfo MigrationInfo `json:"MigrationInfo"`
}

func GetDirPath() string {
	envPath := os.Getenv("GO_CONFIG_DIR")
	if envPath == "" {
		envPath = "./"
	}
	fmt.Println(envPath)
	return envPath
}

func LoadDefaultConfig() Configuration {
	confPath := filepath.Join("E:\\DuoProject\\Service\\GO-Projects\\src\\DVP-DB2ProfileMigrator", "conf.json")
	fmt.Println("GetDefaultConfig config path: ", confPath)
	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)
	if deferr != nil {
		log.Panic(deferr)
	}
	return defconfiguration
}

func LoadConfiguration() Configuration {
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "custom-environment-variables.json")
	fmt.Println("InitiateRedis config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}
	envConfig := Configuration{}
	envconfiguration := Configuration{}
	enverr := json.Unmarshal(content, &envconfiguration)
	if enverr != nil {
		fmt.Println("Fail to Load Environment Settings and Load Default Settings :", enverr)
		envConfig = LoadDefaultConfig()


	} else {

		envConfig.DB.Database = os.Getenv(envconfiguration.DB.Database)
		envConfig.DB.Database = os.Getenv(envconfiguration.DB.Database)
		envConfig.DB.Host = os.Getenv(envconfiguration.DB.Host)
		envConfig.DB.Password = os.Getenv(envconfiguration.DB.Password)
		envConfig.DB.Port = os.Getenv(envconfiguration.DB.Port)
		envConfig.DB.Type = os.Getenv(envconfiguration.DB.Type)
		envConfig.DB.User = os.Getenv(envconfiguration.DB.User)

		envConfig.Mongo.Host = os.Getenv(envconfiguration.Mongo.Host)
		envConfig.Mongo.Port = os.Getenv(envconfiguration.Mongo.Port)
		envConfig.Mongo.Database = os.Getenv(envconfiguration.Mongo.Database)
		envConfig.Mongo.Password = os.Getenv(envconfiguration.Mongo.Password)
		envConfig.Mongo.User = os.Getenv(envconfiguration.Mongo.User)
		envConfig.Mongo.CollectionName = os.Getenv(envconfiguration.Mongo.CollectionName)

		envConfig.DBTWO.Database = os.Getenv(envconfiguration.DBTWO.Database)
		envConfig.DBTWO.Host = os.Getenv(envconfiguration.DBTWO.Host)
		envConfig.DBTWO.Password = os.Getenv(envconfiguration.DBTWO.Password)
		envConfig.DBTWO.Port = os.Getenv(envconfiguration.DBTWO.Port)
		envConfig.DBTWO.Protocol = os.Getenv(envconfiguration.DBTWO.Protocol)
		envConfig.DBTWO.User = os.Getenv(envconfiguration.DBTWO.User)

		envConfig.Host.Domain = os.Getenv(envconfiguration.Host.Domain)
		envConfig.Host.Hostpath = os.Getenv(envconfiguration.Host.Hostpath)
		envConfig.Host.Logfilepath = os.Getenv(envconfiguration.Host.Logfilepath)
		envConfig.Host.Port = os.Getenv(envconfiguration.Host.Port)
		envConfig.Host.Version = os.Getenv(envconfiguration.Host.Version)

		envConfig.Security.AccessToken = os.Getenv(envconfiguration.Security.AccessToken)
		envConfig.Security.Ip = os.Getenv(envconfiguration.Security.Ip)
		envConfig.Security.Port = os.Getenv(envconfiguration.Security.Port)

		envConfig.MigrationInfo.Tenant = os.Getenv(envconfiguration.MigrationInfo.Tenant)
		envConfig.MigrationInfo.SchedulerTime = os.Getenv(envconfiguration.MigrationInfo.SchedulerTime)
		tempMigrationInfo := os.Getenv(envconfiguration.MigrationInfo.CompanysName)
		var mInfo []string
		json.Unmarshal([]byte(tempMigrationInfo), &mInfo)
		envConfig.MigrationInfo.Companys = mInfo


	}
	fmt.Println("Configurations :", envConfig)
	return envConfig
}
