package models

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"time"
	"DVP-DB2ProfileMigrator/sherardFunctions"
	"fmt"
	"github.com/go-contrib/uuid"
	"gopkg.in/mgo.v2/bson"
)

type (
	Address struct {
		Zipcode  string
		Number   string
		Street   string
		City     string
		Province string
		Country  string
	}

	ExternalUsers struct {
		ThreadPartyReference string
		Tenant               int
		Company              int
		Title                string
		Firstname            string
		Lastname             string
		Gender               string
		Name                 string
		Phone                string
		Email                string
		Locale               string
		Address              Address
		Tags                 [] string
	}
)

func SaveProfilesToMongo(userProfiles []ExternalUsers, uid uuid.UUID, tenant int, company int) {

	s := fmt.Sprintf("Profile Migrattion. Save To MongoDb. Process Start -> %s.  Tenant : %d , Company : %d ", uid, tenant, company)
	fmt.Println(s)

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{config.Mongo.Host + ":" + config.Mongo.Port},
		Timeout:  60 * time.Second,
		Database: config.Mongo.Database,
		Username: config.Mongo.User,
		Password: config.Mongo.Password,
	}


	//session, err := mgo.Dial("104.236.231.11,server2.example.com")
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		sherardFunctions.Throw(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(config.Mongo.Database).C(config.Mongo.CollectionName)
	for i := 0; i < len(userProfiles); i++ {
		tempProfile := userProfiles[i];
		// check profile is existing
		result := ExternalUsers{}
		err = c.Find(bson.M{"threadpartyreference": tempProfile.ThreadPartyReference,"tenant":tenant,"company":company}).One(&result)

		if err != nil {
			switch err {
			case mgo.ErrNotFound:
				{
					err = c.Insert(userProfiles[i])
					if err != nil {
						fmt.Println(err)
					}
				}
			default:
				fmt.Println(err)
			}
		}else {
			colQuerier := bson.M{"threadpartyreference": tempProfile.ThreadPartyReference}
			err = c.Update(colQuerier, tempProfile)
			if err != nil {
				fmt.Println(err)
			}
		}

		/*if ( result.ThreadPartyReference == "") {
			//result.ThreadPartyReference == "" ||
			err = c.Insert(userProfiles[i])
			//err = c.Insert(&ExternalUsers{userProfiles.Title, userProfiles.Firstname, userProfiles.Lastname, userProfiles.Gender, userProfiles.Name, userProfiles.Phone, userProfiles.Email, userProfiles.Locale, userProfiles.Address, userProfiles.Tags})
			if err != nil {
				fmt.Println(err)
			}
		} else {
			colQuerier := bson.M{"threadpartyreference": tempProfile.ThreadPartyReference}
			//change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}
			err = c.Update(colQuerier, tempProfile)
			if err != nil {
				fmt.Println(err)
			}
		}
*/
	}
	ss := fmt.Sprintf("Profile Migrattion. Save To MongoDb. Process Complete -> %s.  Tenant : %d , Company : %d ", uid, tenant, company)
	fmt.Println(ss)
}

func SaveProfiles(userProfiles []ExternalUsers, uid uuid.UUID, tenant int, company int) {

	sherardFunctions.Block{
		Try: func() {

			SaveProfilesToMongo(userProfiles, uid, tenant, company)
		},
		Catch: func(e sherardFunctions.Exception) {
			fmt.Printf("Caught %v\n", e)
		},
		Finally: func() {
			s := fmt.Sprintf("Profile Migrattion. Save To MongoDb. Process Completed. -> %s.  Tenant : %d , Company : %d ", uid, tenant, company)
			fmt.Println(s)
		},
	}.Do()




	/*err = c.Insert(&userProfiles)
	//err = c.Insert(&ExternalUsers{userProfiles.Title, userProfiles.Firstname, userProfiles.Lastname, userProfiles.Gender, userProfiles.Name, userProfiles.Phone, userProfiles.Email, userProfiles.Locale, userProfiles.Address, userProfiles.Tags})
	if err != nil {
		log.Fatal(err)
	}*/

	/*result := ExternalUsers{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)*/
}