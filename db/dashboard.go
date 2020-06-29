// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import (
	"fmt"

	"github.com/cuttle-ai/octopus-service/config"
	"github.com/jinzhu/gorm"
)

/*
 * This file contains the database interactions for the dashboards
 */

//Dashboard is the dashboard model
type Dashboard struct {
	gorm.Model
	//Name of the dashbaord
	Name string
	//Description of the dashboard
	Description string
	//UserID of the user who created the dashboard
	UserID uint
	//IsPublic indicates whether the dashboard is public
	IsPublic bool
	//HasPublicWidgets indicates that the dashboard has widgets that do have public access
	HasPublicWidgets bool
	//ShowNavigation indicates whether the navigation for the dashboard has to be made visible
	ShowNavigation bool
	//DashboardPages has the list of pages in the dashboard
	DashboardPages []DashboardPage
}

//DashboardUserMappings has the mappings for the users and the permissions they have with the dashboard
type DashboardUserMappings struct {
	gorm.Model
	//DashboardID  is the id of the dashboard
	DashboardID uint
	//UserID is the id of the user
	UserID uint
	//Share indicates whether the user has the permissions for sharing the dashboard
	Share bool
	//Manage indicates whether the user has the permissions for managing the user permissions to the dashboard
	Manage bool
	//Edit indicates whether the user has the permissions for editing the dashboard
	Edit bool
}

//DashboardPage is the model for representing a dashboard page
type DashboardPage struct {
	gorm.Model
	//DashboardID  is the id of the dashboard
	DashboardID uint
	//Name is the name of the dashboard page
	Name string
	//Number is the page number in the dashboard for cronological order
	Number uint
	//GridSize is the size of the each grid uint. grid will be a square of the given size
	GridSize uint
	//Width is the width of the page grid.
	Width uint
	//Height is the height of the page grid.
	Height uint
	//HasWidgetAdded indicates that a widget has been added to the dashboard page
	HasWidgetAdded bool
	//PageGridItems has the list of grid items in the page
	PageGridItems []PageGridItem
}

const (
	//PageDefaultWidth is the default width of the page
	PageDefaultWidth = uint(100)
	//PageDefaultHeight is the default height of the page
	PageDefaultHeight = uint(100)
	//PageDefaultGridSize is the default page grid size
	PageDefaultGridSize = uint(10)
)

//PageGridItem is a grid item in the page layout. It will be linked to the underlying widget.
//Page Grid items together builds the page layout
type PageGridItem struct {
	gorm.Model
	//DashboardPageID is the id of the page in which the page grid item is present
	DashboardPageID uint
	//WidgetID is id of the underlying widget
	WidgetID uint
	//X is the x position of the grid item in the grid (in grid units)
	X uint
	//Y is the y position of the grid item in the grid (in grid units)
	Y uint
	//Width is the width of the grid item in the grid (in grid units)
	Width uint
	//Height is the height of the grid item in the grid (in grid units)
	Height uint
}

//AddWidget will add a widget to the dashboard
func (d *Dashboard) AddWidget(ctx *config.AppContext, w Widget, width, height uint) error {
	/*
	 * We will get the last page in the dashboard
	 * We will try to add the widget inside the page grid
	 * If not possible we will create a new page and add the widget to that
	 */
	//getting the last page in the dashboard
	page, err := d.GetLastPage(ctx)
	if err != nil {
		//error while getting the last page of the dashboard
		ctx.Log.Error("error while getting the last page of the dashboard", d.ID)
		return err
	}

	//will try to add widget inside the page grid
	ok, err := page.AddWidget(ctx, w, width, height)
	if err != nil {
		//error while addding a widget to the page
		ctx.Log.Error("error while adding the widget to the page", page.ID)
		return err
	}

	if ok {
		ctx.Log.Info("added the widget", w.ID, "to the dashboard", d.ID)
		return nil
	}

	ctx.Log.Info("couldn't add the widget to the page. So creating a page and adding the widget to that for dashboard", d.ID)
	//will create a page if not able to add widget to the page
	page, err = d.CreatePage(ctx, page.Number+1)
	if err != nil {
		//error while creating a page
		ctx.Log.Error("error while creating a page for the dashboard", d.ID)
		return err
	}
	//and add the widget to the page
	ok, err = page.AddWidget(ctx, w, width, height)
	if err != nil {
		//error while addding a widget to the page
		ctx.Log.Error("error while adding the widget to the newly created page", page.ID)
		return err
	}
	return nil
}

//GetLastPage returns the last page in the dashboard
func (d *Dashboard) GetLastPage(ctx *config.AppContext) (*DashboardPage, error) {
	dp := &DashboardPage{}
	err := ctx.Db.Where("dashboard_id = ?", d.ID).Order("number DESC").First(dp).Error
	if err != nil {
		return nil, err
	}
	return dp, nil
}

//CreatePage creates a new page for the given dashboard with the given page number
func (d *Dashboard) CreatePage(ctx *config.AppContext, pageNumber uint) (*DashboardPage, error) {
	newPage := &DashboardPage{
		DashboardID:    d.ID,
		Name:           fmt.Sprintf("%s - %d", d.Name, pageNumber),
		Number:         pageNumber,
		GridSize:       PageDefaultGridSize,
		Width:          PageDefaultWidth,
		Height:         PageDefaultHeight,
		HasWidgetAdded: true,
	}
	err := ctx.Db.Create(newPage).Error
	if err != nil {
		return nil, err
	}
	return newPage, nil
}

//AddWidget will try to add a widget to the dashboard page. If succeeds will return true. Else false.
func (dp *DashboardPage) AddWidget(ctx *config.AppContext, w Widget, width, height uint) (bool, error) {
	pageGrid := &PageGridItem{
		DashboardPageID: dp.ID,
		WidgetID:        w.ID,
		X:               0,
		Y:               0,
		Width:           width,
		Height:          height,
	}
	if len(dp.PageGridItems) == 0 {
		err := ctx.Db.Create(pageGrid).Error
		if err == nil {
			dp.PageGridItems = append(dp.PageGridItems, *pageGrid)
		}
		return err == nil, err
	}
	return false, nil
}

//GetPageLayout will return the page layout filled with the occupied positions as true
func (dp *DashboardPage) GetPageLayout() [][]bool {
	/*
	 * We will first initialize the grid
	 * Then we will set the cells in the grid as true where the grid item exists
	 */
	//initializing the grid
	grid := [][]bool{}
	if dp.Height == 0 || dp.Width == 0 {
		return grid
	}

	for i := uint(0); i < dp.Height; i++ {
		row := []bool{}
		for j := uint(0); j < dp.Width; j++ {
			row = append(row, false)
		}
		grid = append(grid, row)
	}

	//setting the cells in the grid as true where the grid items exist
	for _, v := range dp.PageGridItems {
		if v.Y < 0 || v.X < 0 || v.Y >= uint(len(grid)) || v.X >= uint(len(grid[0])) {
			continue
		}
		for i := v.Y; i < v.Y+v.Height; i++ {
			for j := v.X; j < v.X+v.Width; j++ {
				grid[i][j] = true
			}
		}
	}
	return grid
}
