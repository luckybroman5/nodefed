package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-fed/activity/pub"
	"github.com/luckybroman5/fediverse/db"
	"github.com/luckybroman5/fediverse/service"
)

type Actor interface {
	PostInbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error)
	GetInbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error)
	PostOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error)
	GetOutbox(c context.Context, w http.ResponseWriter, r *http.Request) (bool, error)
}

func main() {
	s := &service.Service{}
	db := &db.InMem{
		content:  &sync.Map{},
		locks:    &sync.Map{},
		hostname: "localhost",
	}

	actor := pub.NewFederatingActor(s, s, db, s)
	mux := http.NewServeMux()
	// Map the `me` actor's inbox to the path `/actors/me/inbox`
	mux.HandleFunc("/actors/me/inbox", func(w http.ResponseWriter, r *http.Request) {
		if isActivityPubRequest, err := actor.GetInbox(w, r); err != nil {
			// Do something with `err`
			return
		} else if isActivityPubRequest {
			// Go-Fed handled the ActivityPub GET request to the inbox
			return
		} else if isActivityPubRequest, err := actor.PostInbox(w, r); err != nil {
			// Do something with `err`
			return
		} else if isActivityPubRequest {
			// Go-Fed handled the ActivityPub POST request to the inbox
			return
		}
		// Here we return an error, but you may just as well decide
		// to render a webpage instead. But be sure to apply appropriate
		// authorizations. There's no guarantees about authorization at
		// this point.
		http.Error(w, "Non-ActivityPub request", http.StatusBadRequest)
		return
	})
	// Map the `me` actor's inbox to the path `/arbitrary/me/outbox`
	mux.HandleFunc("/arbitrary/me/outbox", func(w http.ResponseWriter, r *http.Request) {
		if isActivityPubRequest, err := actor.GetOutbox(w, r); err != nil {
			// Do something with `err`
			return
		} else if isActivityPubRequest {
			// Go-Fed handled the ActivityPub GET request to the outbox
			return
		} else if isActivityPubRequest, err := actor.PostOutbox(w, r); err != nil {
			// Do something with `err`
			return
		} else if isActivityPubRequest {
			// Go-Fed handled the ActivityPub POST request to the outbox
			return
		}
		// Here we return an error, but you may just as well decide
		// to render a webpage instead. But be sure to apply appropriate
		// authorizations. There's no guarantees about authorization at
		// this point.
		http.Error(w, "Non-ActivityPub request", http.StatusBadRequest)
		return
	})
}
