package graph

import "10/graph/model"

// It serves as dependency injection for your app, add any dependencies you require here.

//This are the types we use to output our data

type Resolver struct{
	users []*model.User
	newuser []*model.NewUser
	events []*model.Event
	eventusers []*model.EventUser
}

