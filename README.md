# Notion Graph View

[Notion](https://www.notion.so) is a knowledge management system with a very nice API. You can easily link pages. This application attempts to visualise those links nicely in a similar way that [Obsidian](https://obsidian.md), another knowledge base project, does.

![Notion Graph Example](/docs/graph_example.png)

## Try it for yourself

- Generate a [Notion API Key](https://www.notion.so/my-integrations)
- Give the API key access to the databases that you want to have included in the graph
- Navigate to `https://christopherhex.github.io/notion-graph-view?notionKey=<YOUR_API_KEY>`

Currently the application will only look for `mentions` of pages. So references in database properties are not included yet in the algorithm. The application is intended to run fully client-side in the browser. Due to CORS restrictions, it's not possible to call the Notion API directly from a client-side PWA application. Therefore a simple CORS proxy is used to make the Notion API available client-side.

## Run it yourself

- Make sure you have `node` and `go` installed
- Checkout the project
- Deploy your own proxy and update the url in `src/wasm/notiongraph.go`
- Run `npm run build`
- Run `npm run dev`

## To Do

- [ ] Improve UI to make it render graph better
- [ ] Add caching in localstorage to reduce the number of Notion Calls
- [ ] Add manual refresh button
- [ ] Allow to click on notes and open the corresponding notion page
- [ ] Find other way to store/supply Notion token
