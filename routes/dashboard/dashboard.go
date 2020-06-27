// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package dashboard has the implementation of the dashboard api for the server
package dashboard

import "github.com/cuttle-ai/octopus-service/db"

//Dashboard data transilation object
type Dashboard struct {
	db.Dashboard
}
