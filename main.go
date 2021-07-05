package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"strconv"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
)

type User struct {
	Id           uint `gorm:"primary_key"`
	Username     string
	Password     string
	PasswordShow string
	Quota        int64
	Download     int64
	Upload       int64
	UseDays      int
	ExpiryData   int64
}

type Server struct {
	Name string
	Addr string
	Port int
	Area string
}

//go:embed template.yaml
var template []byte

type clashConfig struct {
	Proxies     []proxy
	ProxyGroups []proxyGroup `yaml:"proxy-groups"`
}

type proxy struct {
	Name      string
	ProxyType string `yaml:"type"`
	Server    string
	Port      int
	Password  string
}

type proxyGroup struct {
	Name      string
	GroupType string `yaml:"type"`
	Proxies   []string
}

func main() {
	var connstr string
	flag.StringVar(&connstr, "db", "", "数据库链接")
	flag.Parse()
	db, err := gorm.Open("mysql", connstr)
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.GET("/clash", func(c *gin.Context) {
		passwd := c.Query("passwd")
		passwdBase64 := base64.StdEncoding.EncodeToString([]byte(passwd))
		usr := User{}
		db.First(&usr, "passwordShow = ?", passwdBase64)
		if usr.Id == 0 {
			c.AbortWithStatus(403)
			return
		}
		c.Writer.Write(template)
		var clashConf clashConfig
		var servers []Server
		db.Find(&servers)
		clashConf.ProxyGroups = []proxyGroup{{
			Name:      "Proxy",
			GroupType: "select",
		}}
		leftFlow := leftFlowFmt(usr.Quota - usr.Download - usr.Upload)
		servers = append(servers, Server{Name: fmt.Sprintf("流量剩余:%s", leftFlow)})
		for _, v := range servers {
			clashConf.Proxies = append(clashConf.Proxies, proxy{
				Name:      v.Name,
				ProxyType: "trojan",
				Server:    v.Addr,
				Port:      v.Port,
				Password:  passwd,
			})
			clashConf.ProxyGroups[0].Proxies = append(clashConf.ProxyGroups[0].Proxies, v.Name)
		}
		clashByte, err := yaml.Marshal(&clashConf)
		if err != nil {
			fmt.Println(err)
			return
		}

		c.Writer.WriteString("\n")
		c.Writer.Write(clashByte)
	})
	r.Run(":25001")
}

var levelName = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func leftFlowFmt(leftFlow int64) string {
	if leftFlow < 0 {
		return "+∞"
	}
	num := float64(leftFlow)
	var level int8
	for ; num > 1024 && level < 6; level++ {
		num = num / 1024
	}
	return strconv.FormatFloat(num, 'f', 2, 64) + levelName[level]
}
