package adddatas

import (
	"testRouterAPI/go_recipe/recipe" 

	//"encoding/json"
	"fmt"
	"log"
	"testRouterAPI/go_recipe/router"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
)

//var qq recipe.Recipe
var sum int 


func CreateProcess(id string) *recipe.Process {
	var p2 []*recipe.Process_Config
	//sql
	//---> result
	// for  Len(result) {}
	p2 = append(p2, CreateProcessConfig(""))
	return nil
}

func CreateProcessConfig(id string) *recipe.Process_Config {
	return nil
}

func CreateProcessConfigStep(id string) *recipe.Process_Config_Step {
	return nil
}

func CreateProcessConfigSetting(id string) *recipe.Process_Config_Step_Setting {
	return nil
}

func Add() {
	var aa recipe.Recipe
	var p []*recipe.Process
	// var p2 []recipe.Process_Config
	// var p3 []recipe.Process_Config_Step
	// var p4 []recipe.Process_Config_Step_Setting


	ids := []string{}
	for _, id := range ids {
		p = append(p, CreateProcess(id))
	}

	aa.SpecID = ""
	aa.Factory = 2

	p[0].Name = "aa"

	aa.Processes = p

	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	s := router.ExampleSpec()

	encoded, err := proto.Marshal(s)
	if err != nil {
		return
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("recipe"))
		if err != nil {
			return err
		}
		return b.Put([]byte("T00001"), []byte(encoded))
	})
}

func View() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("recipe"))
		//v := b.Get([]byte("T00001"))
		//fmt.Printf("%s", v)

		if err := b.ForEach(func(k, v []byte) error {
			var ouput recipe.Recipe
			if err := proto.Unmarshal(v, &ouput); err != nil {
				return err
			}
			fmt.Printf("%s\n%s\n%s\n", ouput.Id, ouput.ProductID, ouput.Processes)

			return nil
		}); err != nil {
			return err
		}

		return nil
	})
}
