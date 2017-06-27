package main

import (
	"sync"
	"fmt"
	"time"
	"os"
	"bufio"
	"strings"
	"DVP-DB2ProfileMigrator/models"
	"github.com/go-contrib/uuid"
	"github.com/jasonlvhit/gocron"

	"DVP-DB2ProfileMigrator/sherardFunctions"
	"encoding/json"
	"strconv"
)

var wg sync.WaitGroup
var (
	config = sherardFunctions.LoadConfiguration()
)

func main() {

	// Create at least 1 goroutine
	wg.Add(1)

	go forever()

	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	s := gocron.NewScheduler()
	//s.Every(1).Minute().Do(migrateNewProfile)
	s.Every(1).Day().At(config.MigrationInfo.SchedulerTime).Do(migrateNewProfile,true) //'hour:min'
	<-s.Start()

	//wg.Wait()
}

func migrateNewProfile(newProfile bool) {

	fmt.Println("Start Schedule task.")

	var mInfo []string
	json.Unmarshal([]byte(config.MigrationInfo.CompanysName), &mInfo)

	query :="CALL PROFILELIST();"
	if newProfile {
		query ="CALL NEWPROFILELIST();"
	} else {
		query ="CALL PROFILELIST();"
	}
	for i := 0; i < len(mInfo); i ++ {

		sherardFunctions.Block{
			Try: func() {
				uid := uuid.NewV4()
				t, _ := strconv.Atoi(config.MigrationInfo.Tenant)
				c, _ := strconv.Atoi(mInfo[i])
				go models.MigrateProfile(uid, query, t, c)
			},
			Catch: func(e sherardFunctions.Exception) {
				s := fmt.Sprintf("Fail To Start Schedule Task. Tenant : %s , Company : %s ", config.MigrationInfo.Tenant, mInfo[i])
				fmt.Println(s)
			},
			Finally: func() {
				s := fmt.Sprintf("Start Schedule Task, Profile Migrattion Process. Tenant : %s , Company : %s ", config.MigrationInfo.Tenant, mInfo[i])
				fmt.Println(s)
			},
		}.Do()

	}

}

func forever() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Command To Excute : ")
		text, _ := reader.ReadString('\n')

		commandToExecute := strings.ToLower(strings.TrimSpace(text))

		switch  commandToExecute{
		case "exit":{
			fmt.Printf("graceful exit") // return the program name back to %s
			f()
			os.Exit(1) // graceful exit
		}
		case "help":{
			fmt.Println("Commands......")
			fmt.Println("exit - >  Exit From Application.")
			fmt.Println("upload - > Migrate Profile To Facetone.")
			fmt.Println("test\n")
		}
		case "upload":{
			go migrateNewProfile(false) //models.MigrateProfile(uid, "CALL PROFILELIST();", 1, 103)
		}
		case "save":{
			udata := []models.ExternalUsers{
				{
					Tenant :1,
					Company :103,
					Title:"Mr",
					Firstname:"Test",
					Lastname:"LastName",
					Gender:"Male",
					Name:"Test",
					Phone:"12346",
					Email:"test@fo.lk",
					Locale:"en",
					Address: models.Address{
						City:"",
						Country:"",
						Number:"",
						Province:"",
						Street:"",
						Zipcode:"123",
					},
				},
				{
					Tenant :1,
					Company :103,
					Title:"Mr",
					Firstname:"Test111",
					Lastname:"LastName111111",
					Gender:"Male",
					Name:"Tes1111111t",
					Phone:"123333333333346",
					Email:"test@fo.lk",
					Locale:"en",
					Address: models.Address{
						City:"",
						Country:"",
						Number:"",
						Province:"",
						Street:"",
						Zipcode:"125553",
					},
				},
			}
			uid := uuid.NewV4()
			go models.SaveProfiles(udata, uid, 1, 103);
		}
		case "test":{
			kvs := map[string]string{"a": "apple", "b": "banana"}
			for k, v := range kvs {
				fmt.Printf("%s -> %s\n", k, v)
			}
			fmt.Println("Application Up And Running.")
		}

			time.Sleep(time.Second)
		/*if false {

			wg.Add(1)
			go f()
		}*/
		}
	}
}

func f() {
	// When the termination condition for this goroutine is detected, do:
	wg.Done()
}