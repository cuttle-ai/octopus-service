// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import "github.com/jinzhu/gorm"

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
	Name uint
	//Number is the page number in the dashboard for cronological order
	Number uint
	//GridSize is the size of the each grid uint. grid will be a square of the given size
	GridSize uint
	//Width is the width of the page grid.
	Width uint
	//Height is the height of the page grid.
	Height uint
}

//PageGridItem is a grid item in the page layout. It will be linked to the underlying widget.
//Page Grid items together builds the page layout
type PageGridItem struct {
	gorm.Model
	//PageID is the id of the page in which the page grid item is present
	PageID uint
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
