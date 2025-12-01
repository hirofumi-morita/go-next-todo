package handlers

import "todo-app/graph"

var GQLClient *graph.Client

func InitGraphQLClient(client *graph.Client) {
	GQLClient = client
}
