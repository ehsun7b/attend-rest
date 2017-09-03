package controllers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/coopernurse/gorp"
	"github.com/ehsun7b/attend-rest/app/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(InitDb)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)

	revel.TimeFormats = append(revel.TimeFormats, "2016-06-01 12:12")
}

func getParamString(param string, defaultValue string) string {
	p, found := revel.Config.String(param)
	if !found {
		if defaultValue == "" {
			revel.ERROR.Fatal("Cound not find parameter: " + param)
		} else {
			return defaultValue
		}
	}
	return p
}

func getConnectionString() string {
	host := getParamString("db.host", "")
	port := getParamString("db.port", "3306")
	user := getParamString("db.user", "")
	pass := getParamString("db.password", "")
	dbname := getParamString("db.name", "attend")
	protocol := getParamString("db.protocol", "tcp")
	dbargs := getParamString("dbargs", " ")

	if strings.Trim(dbargs, " ") != "" {
		dbargs = "?" + dbargs
	} else {
		dbargs = ""
	}
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s",
		user, pass, protocol, host, port, dbname, dbargs)
}

//InitDb creates a connections
func InitDb() {
	connectionString := getConnectionString()
	if db, err := sql.Open("mysql", connectionString); err != nil {
		revel.ERROR.Fatal(err)
	} else {
		Dbm = &gorp.DbMap{
			Db: db,
			Dialect: gorp.MySQLDialect{
				Engine:   "InnoDB",
				Encoding: "UTF8",
			},
		}
	}

	defineEventTable(Dbm)
	// in the case of any DDL change we will uncomment this line of one run to drop old tables
	Dbm.DropTablesIfExists()
	if err := Dbm.CreateTablesIfNotExists(); err != nil {
		revel.ERROR.Fatal(err)
	}
}

func defineEventTable(dbm *gorp.DbMap) {
	// set "id" as primary key and autoincrement
	t := dbm.AddTable(models.Event{}).SetKeys(true, "id")
	t.ColMap("title").SetMaxSize(models.EventTitleMaxLength).SetNotNull(true)
	t.ColMap("tag").SetNotNull(true).SetUnique(true)
	t.ColMap("status").SetNotNull(true)
	t.ColMap("address").SetMaxSize(models.EventAddressMaxLength)
	t.ColMap("category").SetMaxSize(models.EventCategoryMaxLength).SetNotNull(true)
	t.ColMap("created_at").SetNotNull(true)
	t.ColMap("time").SetNotNull(true)
	t.ColMap("message").SetMaxSize(models.EventMessageMaxLength)
	t.ColMap("map_link").SetMaxSize(models.EventLinkleMaxLength)
	t.ColMap("tag").SetMaxSize(models.EventTagMaxLength)
}
