// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package db

import "github.com/jinzhu/gorm"

/*
 * This file contains the defition of the widget in the dashboard
 */

//Widget represents a widget which can be a visualization or something else in a dashboard page
type Widget struct {
	gorm.Model
}
