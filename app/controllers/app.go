package controllers

import (
	"github.com/revel/revel"
	"obrolansubuh.com/models"
	osc "obrolansubuh.com/modules/gorm/app/controllers"
)

type App struct {
	osc.GormController
}

func (c App) Index() revel.Result {
	var posts []models.Post
	c.Trx.Limit(5).Order("created_at desc").
		Where("published = true").
		Find(&posts)

	return c.Render(posts)
}

func (c App) PostList(page int) revel.Result {
	var posts []models.Post

	perpage := 10
	c.Trx.Preload("Author").Preload("Categories").Limit(perpage).Offset((page - 1) * perpage).
		Order("created_at desc").
		Where("published = true").
		Find(&posts)

	var postCount, prevPage, nextPage int
	c.Trx.Model(models.Post{}).Count(&postCount)

	if page >= 2 {
		prevPage = page - 1
	}

	if page*perpage < postCount {
		nextPage = page + 1
	}

	return c.Render(posts, prevPage, nextPage)
}

func (c App) Post(id int64, slug string) revel.Result {
	var post models.Post
	c.Trx.Where("id = ?", id).Find(&post)

	if !post.Published {
		return c.NotFound("Tulisan tidak ditemukan :(")
	}

	return c.Render(post)
}

func (c App) Writers() revel.Result {
	var writers []models.Contributor
	c.Trx.Find(&writers)

	return c.Render(writers)
}

func (c App) About() revel.Result {
	var about models.SiteInfo
	c.Trx.First(&about)

	return c.Render(about)
}
