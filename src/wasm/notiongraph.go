package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"syscall/js"

	"github.com/tidwall/gjson"
)

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
	Databases []NotionDatabase
}

var mu sync.Mutex
var notionToken string

func getData() NotionGraph {

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
	return NotionGraph{
		Pages: pagesToCheck,
		Links: pageLinks,
		Databases: availableDatabases,
	};


}



func genericNotionRequest(method string, url string, data []byte) []byte {
	full_url := "https://cors-proxy.creinto.workers.dev/v1" + url

	// Data := []byte(``)

	req, _ := http.NewRequest(method, full_url, bytes.NewBuffer(data))

	req.Header.Add("x-forwarded-to", "api.notion.com")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Notion-Version", "2022-02-22")
	req.Header.Add("Authorization", "Bearer " + notionToken)


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

func main(){
	c := make(chan int)

	getDataFunction := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		notionToken = args[0].String()


		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			// reject := args[1]

			go func(){
				res := getData()

				fmt.Println(res);

				objectConstructor := js.Global().Get("String")

				// resData := map[string]interface{}{
				// 	"Pages": res.Pages,
				// 	"Databases": res.Databases,
				// 	"Links": res.Links,
				// }

				jsonMapAsStringFormat, _ := json.Marshal(res)

				resolve.Invoke(objectConstructor.New(string(jsonMapAsStringFormat)))
			}()
	
			return nil
		});

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})

	js.Global().Set("getData", getDataFunction)

	<-c
}
