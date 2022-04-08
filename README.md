# Notion Graph View

[Notion](https://www.notion.so) is a knowledge management system with a very nice API. You can easily link pages. This project attempts to visualise those links nicely in a similar way that [Obsidian](https://obsidian.md), another knowledge base project, does.

![Notion Graph Example](/assets/graph_example.png)

## How to create the knowledge graph.

- [Generate a Notion API Key](https://www.notion.so/my-integrations) and add it to a `.env.local` file in the root of the repository as `NOTION_API_KEY=<secret>`
- Give the integration access to all databases that you want to include in the graph.
- Run the tool: `go run notiongraph.go`. This will generate a JSON file `test.json` that contains all the graph data
- To visualise the graph, use a simple http server (e.g. [http-server](https://www.npmjs.com/package/http-server)) to host the `index.html` file
- Open your browser and the graph should appear

## How to use
