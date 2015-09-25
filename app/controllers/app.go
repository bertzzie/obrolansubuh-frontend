package controllers

import (
	"github.com/revel/revel"
	"obrolansubuh.com/models"
	osc "obrolansubuh.com/modules/gorm/app/controllers"
)

type App struct {
	osc.GormController
}

// Post per page
const perpage = 10

func (c App) Index() revel.Result {
	var posts []models.Post
	c.Trx.Limit(5).Order("created_at desc").
		Where("published = true").
		Find(&posts)

	return c.Render(posts)
}

func (c App) Posts(page int) revel.Result {
	var posts []models.Post

	// default value
	if page == 0 {
		page = 1
	}

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
	c.Trx.Preload("Author").Where("id = ?", id).Find(&post)

	if !post.Published {
		return c.NotFound("Tulisan tidak ditemukan :(")
	}

	return c.Render(post, slug)
}

func (c App) Writers() revel.Result {
	var writers []models.Contributor
	c.Trx.Find(&writers)

	return c.Render(writers)
}

func (c App) WritersPosts(handle string, page int) revel.Result {
	writer, err := c.GetContributorByHandle(handle)
	if err != nil {
		return c.NotFound("Penulis tidak ditemukan :(")
	}

	// default value
	if page == 0 {
		page = 1
	}

	var posts []models.Post
	c.Trx.Preload("Author").Preload("Categories").Limit(perpage).Offset((page-1)*perpage).
		Order("created_at desc").
		Where("author_id = ? AND published = true", writer.ID).
		Find(&posts)

	var postCount, prevPage, nextPage int
	c.Trx.Model(models.Post{}).Where("author_id = ?", writer.ID).Count(&postCount)

	if page >= 2 {
		prevPage = page - 1
	}

	if page*perpage < postCount {
		nextPage = page + 1
	}

	return c.Render(writer, posts, prevPage, nextPage)
}

func (c App) CategoriesPosts(id int, slug string, page int) revel.Result {
	if page == 0 {
		page = 1
	}

	var posts []models.Post
	var category models.Category
	c.Trx.Preload("Author").Preload("Categories").Limit(perpage).Offset((page-1)*perpage).
		Order("created_at desc").
		Joins("inner join post_categories on post_categories.post_id = posts.id").
		Where("posts.published = true AND post_categories.category_id = ?", id).
		Find(&posts)

	c.Trx.Where("id = ?", id).Find(&category)

	var postCount, prevPage, nextPage int
	postCount = c.Trx.Model(&category).Association("Posts").Count()

	if page >= 2 {
		prevPage = page - 1
	}

	if page*perpage < postCount {
		nextPage = page + 1
	}

	return c.Render(category, posts, prevPage, nextPage, slug)
}

func (c App) About() revel.Result {
	var about models.SiteInfo
	c.Trx.First(&about)

	return c.Render(about)
}
