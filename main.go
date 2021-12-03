package main

import (
	"context"
	"fmt"

	"main/db"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	// Create a post
	post, err := client.Post.CreateOne(
		db.Post.Title.Set("Prisma"),
	).Exec(ctx)
	if err != nil {
		return err
	}

	// Create a category
	category, err := client.Category.
		CreateOne(db.Category.Name.Set("Excellent Go ORMs")).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Link the category to the post
	_, err = client.Post.FindUnique(db.Post.ID.Equals(post.ID)).
		Update(db.Post.Categories.Link(
			db.Category.ID.InIfPresent([]int{category.ID}),
		),
		).
		Exec(ctx)
	if err != nil {
		return err
	} else {
		fmt.Println("Success!")
	}

	return nil
}
