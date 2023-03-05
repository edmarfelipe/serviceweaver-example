package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ServiceWeaver/weaver"
	"github.com/edmarfelipe/serviceweaver-example/commentservice"
	"github.com/edmarfelipe/serviceweaver-example/postservice"
)

func main() {
	ctx := context.Background()
	root := weaver.Init(ctx)

	commentService, err := weaver.Get[commentservice.Service](root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	postService, err := weaver.Get[postservice.Service](root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	r := http.NewServeMux()
	r.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, fmt.Sprintf("Method %q not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}

		slug := r.URL.Query().Get("slug")
		if len(slug) == 0 {
			http.Error(w, "Parameter slug is not defined in the url", http.StatusBadRequest)
			return
		}

		post, err := postService.GetPost(r.Context(), slug)
		if err != nil {
			fmt.Printf("GetByPost: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			fmt.Printf("Encode: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	r.HandleFunc("/comments/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			msg := fmt.Sprintf("Method %q not allowed", r.Method)
			http.Error(w, msg, http.StatusMethodNotAllowed)
		}

		postId, err := strconv.Atoi(r.URL.Query().Get("postId"))
		if err != nil {
			http.Error(w, "Parameter postId is not valid value", http.StatusBadRequest)
			return
		}

		comments, err := commentService.GetByPost(r.Context(), postId)
		if err != nil {
			fmt.Printf("GetByPost: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(comments)
		if err != nil {
			fmt.Printf("Encode: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	lis, err := root.Listener("main", weaver.ListenerOptions{LocalAddress: ":3030"})
	if err != nil {
		fmt.Printf("root.Listener: %s", err.Error())
		return
	}

	err = http.Serve(lis, r)
	if err != nil {
		fmt.Printf("http.Serve: %s", err.Error())
		return
	}
}
