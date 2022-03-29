package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

// This follows the RAW notion api structure
type RawNotionDataItem struct {
	Id string
	Object string
	Url string
	Name string 
}

type NotionSearchResult struct {
	Object string
	Results []RawNotionDataItem
}

type NotionDatabase struct {
	Id string
	Url string
	Name string
}

type NotionPage struct {
	Id string
	ParentDatabaseId string
	Name string
	Url string
	PageMentions []string
}

type NotionPageLink struct {
	FromPage string
	ToPage string
}

type NotionGraph struct {
	Pages []NotionPage
	Links []NotionPageLink
}

var mu sync.Mutex

func main(){
	godotenv.Load(".env.local")
  
	var wg sync.WaitGroup

	var pagesToCheck []NotionPage
	var pageLinks []NotionPageLink

	availableDatabases := NotionGetAvailableDatabases()

	for _, db := range availableDatabases {
		wg.Add(1)
		go func(id string){
				mu.Lock()
				defer mu.Unlock()
				defer wg.Done()
				pagesToCheck = append(pagesToCheck, NotionGetDatabasePages(id)...)
		}(db.Id)
	}

	// Wait for pages list to be populated
	fmt.Println("Getting pages")
	wg.Wait()
	fmt.Println("Got all pages")
	
	for _, page := range pagesToCheck {
		wg.Add(1)
		go func(page NotionPage){
			defer wg.Done()


			mentions := GetNotionPageMentions(page.Id)
			page.PageMentions = mentions

			if len(mentions) > 0 {
				mu.Lock()
				defer mu.Unlock()
				for _, mention := range mentions {
					pageLinks = append(pageLinks, NotionPageLink{
						FromPage: page.Id,
						ToPage: mention,
					})
				}
			}
		}(page)
	}
	fmt.Println("Getting mentions")
	wg.Wait()
	fmt.Println("Got all mentions")

	// Write To Json Output
	file, _ := json.MarshalIndent(NotionGraph{
		Pages: pagesToCheck,
		Links: pageLinks,
	}, "", " ")
 
	_ = ioutil.WriteFile("test.json", file, 0644)

}


func genericNotionRequest(method string, url string, data []byte) []byte {
	full_url := "https://api.notion.com/v1" + url

	// Data := []byte(``)

	req, _ := http.NewRequest(method, full_url, bytes.NewBuffer(data))

	
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	// @NOCOMMIT
	req.Header.Add("Authorization", "Bearer " + os.Getenv("NOTION_API_KEY"))


	res, _ := http.DefaultClient.Do(req)


	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func GetNotionPageMentions(id string) []string {

	var parsedValues []string

	data_res := genericNotionRequest("GET","/blocks/"+ id + "/children",[]byte(``))

	values := gjson.GetBytes(data_res,`results.#(type=="paragraph")#.paragraph.rich_text.#(type=="mention")#.mention.page.id|@flatten`)

	for _, item := range values.Array() {
		parsedValues = append(parsedValues, item.String())
	}

	return parsedValues
}

// Get a listing of all databases available to the integration
func NotionGetAvailableDatabases() []NotionDatabase {
	
	var foundDatabases []NotionDatabase
	data_res := genericNotionRequest("POST","/search",[]byte(`{"filter":{"property": "object", "value":"database"}}`))
	data_parsed := gjson.GetBytes(data_res,`results`)

	for _, item := range data_parsed.Array() {

		databaseId := gjson.Get(item.Raw,`id`).String()
		databaseUrl := gjson.Get(item.Raw,`url`).String()
		databaseName := gjson.Get(item.Raw,`title.0.plain_text`).String()


		foundDatabases = append(foundDatabases, NotionDatabase{
			Id: databaseId,
			Name: databaseName,
			Url: databaseUrl,
		})
	}

	return foundDatabases
}

func NotionGetDatabasePages(database_id string) []NotionPage {
	var foundPages []NotionPage

	data_res := genericNotionRequest("POST","/databases/"+database_id+"/query",[]byte(``))

	// @TODO introduce recursive checking if check is not complete
	data_parsed := gjson.GetBytes(data_res,`results`)

	for _, item := range data_parsed.Array() {
		pageId := gjson.Get(item.Raw,`id`).String()
		pageParentId := gjson.Get(item.Raw,`parent.database_id`).String()
		pageUrl := gjson.Get(item.Raw,`url`).String()
		pageName := gjson.Get(item.Raw,`title.0.plain_text`).String()

		foundPages = append(foundPages, NotionPage{
			Id: pageId,
			Name: pageName,
			Url: pageUrl,
			ParentDatabaseId: pageParentId,
		})
	}

	return foundPages
}