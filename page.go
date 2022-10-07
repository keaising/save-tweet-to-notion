package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dstotijn/go-notion"
)

func createPageOnDate(ti time.Time) error {
	ctx := context.TODO()
	var err error

	yearName := fmt.Sprintf("%d", ti.Year())
	monthName := fmt.Sprintf("%d.%d", ti.Year(), ti.Month())

	yearNode := pageTree.GetByTitle(yearName)
	if yearNode == nil {
		if err = clt.RateLimiter.Wait(ctx); err != nil {
			return err
		}
		_, err = clt.CreatePage(ctx, notion.CreatePageParams{
			ParentType: notion.ParentTypePage,
			ParentID:   config.RootPageID,
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: yearName,
					},
				},
			},
		})
		if err != nil {
			log.Println("create year page failed:", err)
			return err
		}
		pageTree, _ = getRootTree()
	}
	yearNode = pageTree.GetByTitle(yearName)
	if yearNode == nil {
		log.Println("create year node failed")
		return fmt.Errorf("create year node failed")
	}

	monthNode := pageTree.GetByTitle(monthName)
	if monthNode != nil {
		return nil
	}
	if err = clt.RateLimiter.Wait(ctx); err != nil {
		return err
	}
	_, err = clt.CreatePage(ctx, notion.CreatePageParams{
		ParentType: notion.ParentTypePage,
		ParentID:   yearNode.ID,
		Title: []notion.RichText{
			{
				Text: &notion.Text{
					Content: monthName,
				},
			},
		},
	})
	if err != nil {
		log.Println("create month page failed:", err)
		return err
	}
	pageTree, err = getRootTree()

	return nil
}

func getRootTree() (*TreeNode, error) {
	ctx := context.TODO()
	t := TreeNode{ID: config.RootPageID}

	if err := clt.RateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	page, err := clt.FindPageByID(ctx, config.RootPageID)
	if err != nil {
		log.Println("get root tree: find page by id failed", err)
		return nil, err
	}
	propeties, ok := page.Properties.(notion.PageProperties)
	if !ok {
		log.Println("root is not a page")
		return nil, fmt.Errorf("root is not a page")
	}
	for _, p := range propeties.Title.Title {
		if p.Type == notion.RichTextTypeText {
			t.Title += p.Text.Content
		}
	}

	_ = iterTreeNode(&t)

	return &t, nil
}

func iterTreeNode(t *TreeNode) error {
	ctx := context.TODO()
	cursor := ""

	for {
		blocks, err := clt.FindBlockChildrenByID(ctx, t.ID, &notion.PaginationQuery{
			StartCursor: cursor,
			PageSize:    100,
		})
		if err != nil {
			log.Println("iter tree: find blocks failed", t.ID, err)
			return err
		}
		for _, block := range blocks.Results {
			if block.Type != notion.BlockTypeChildPage {
				continue
			}
			n := TreeNode{
				ID:    block.ID,
				Title: block.ChildPage.Title,
			}
			err = iterTreeNode(&n)
			if err != nil {
				log.Println("iter tree: iter children failed", n.ID, err)
				return err
			}
			t.Children = append(t.Children, &n)
		}
		if !blocks.HasMore {
			break
		}
	}
	return nil
}
