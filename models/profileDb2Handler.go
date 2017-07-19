package models

import (
	"github.com/go-contrib/uuid"
	"fmt"
	"flag"
	"database/sql"
	_"bitbucket.org/phiggins/db2cli"
	"DVP-DB2ProfileMigrator/sherardFunctions"
)

var (
	config = sherardFunctions.LoadConfiguration()
	connStr = flag.String("conn", "DATABASE=" + config.DBTWO.Database + "; HOSTNAME=" + config.DBTWO.Host + "; PORT=" + config.DBTWO.Port + "; PROTOCOL=" + config.DBTWO.Protocol + "; UID=" + config.DBTWO.User + "; PWD=" + config.DBTWO.Password + ";", "connection string to use")
)

func MigrateProfile(uid uuid.UUID, query string, tenant int, company int) {

	sherardFunctions.Block{
		Try: func() {
			s := fmt.Sprintf("Profile Migration. Get Profile From DB2. Process Start -> %s.  Tenant : %d , Company : %d ",uid,tenant,company)
			fmt.Println(s)

			db, err := sql.Open("db2-cli", *connStr)
			if err != nil {
				sherardFunctions.Throw(err)
			}
			defer db.Close()

			rows, err := db.Query(query)
			defer rows.Close()
			if err != nil {
				sherardFunctions.Throw(err)
			}

			var profiles []ExternalUsers
			for rows.Next() {
				var Title string
				var Firstname string
				var Lastname string
				var Gender string
				var Name string
				var Phone string
				var Email string
				var Locale string
				var Zipcode string
				var thirdpartyreference string
				err = rows.Scan(&Title, &Firstname, &Lastname, &Name, &Gender, &Phone, &Email, &Locale, &Zipcode,&thirdpartyreference)
				if err != nil {
					sherardFunctions.Throw(err)
				}

				profile := ExternalUsers{
					thirdpartyreference:thirdpartyreference,
					Tenant :tenant,
					Company :company,
					Title:Title,
					Firstname: Firstname,
					Lastname:Lastname,
					Gender:Gender,
					Name:Name,
					Phone:Phone,
					Email:Email,
					Locale:Locale,
					Address: Address{
						City:"",
						Country:"",
						Number:"",
						Province:"",
						Street:"",
						Zipcode:Zipcode,
					},
				}
				profiles = append(profiles, profile)
			}
			SaveProfilesToMongo(profiles,uid,tenant,company)

		},
		Catch: func(e sherardFunctions.Exception) {
			s := fmt.Sprintf("Profile Migration Process Fail -> %s.  Tenant : %d , Company : %d  %v\n",uid,tenant,company,e)
			fmt.Println(s)
		},
		Finally: func() {
			s := fmt.Sprintf("Profile Migration Process Complete -> %s.  Tenant : %d , Company : %d ",uid,tenant,company)
			fmt.Println(s)
		},
	}.Do()
}

/*
func execQuery(st *sql.Stmt) error {
	rows, err := st.Query()
	if err != nil {
		return err
	}
	defer rows.Close()
	var profiles []ExternalUsers


	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {

		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		for i, col := range columns {

			var v interface{}

			val := values[i]

			b, ok := val.([]byte)

			if (ok) {
				v = string(b)
			} else {
				v = val
			}

			fmt.Println(col, v)
		}
	}


	for rows.Next() {



		var ACTNO string
		var ACTKWD string
		var ACTDESC string

		err = rows.Scan(&ACTNO, &ACTKWD, &ACTDESC)
		if err != nil {
			return err
		}
		//fmt.Printf("ACT: ACTNO %v | ACTKWD %v | ACTDESC %v \n", ACTNO, ACTKWD, ACTDESC)

		profile := ExternalUsers{
			Title:"Mr",
			Firstname: ACTKWD,
			Lastname:ACTDESC,
			Gender:"Male",
			Name:"Test",
			Phone:ACTNO,
			Email:ACTNO + "@duo.lk",
			Locale:"en",
			Address: Address{
				City:"",
				Country:"",
				Number:"",
				Province:"",
				Street:"",
				Zipcode:"123",
			},
		}
		profiles = append(profiles, profile)
	}
	SaveProfilesToMongo(profiles)
	return rows.Err()
}

func dbOperations() error {
	db, err := sql.Open("db2-cli", *connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	// Attention: If you have to go through DB2-Connect you have to terminate SQL-statements with ';'
	st, err := db.Prepare("call profilelist();")
	if err != nil {
		return err
	}
	defer st.Close()

	err = execQuery(st)
	if err != nil {
		return err
	}
	return nil

	*//*for i := 0; i < int(*repeat); i++ {
		err = execQuery(st)
		if err != nil {
			return err
		}
	}
	return nil*//*
}*/
