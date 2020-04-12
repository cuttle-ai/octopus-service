// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package dict will have necessary configuration dictionary in the application
package dict

import (
	bDict "github.com/cuttle-ai/brain/dict"
	"github.com/cuttle-ai/octopus-service/log"
	"github.com/cuttle-ai/octopus/interpreter"
	defaultRules "github.com/cuttle-ai/octopus/rules"
	"github.com/jinzhu/gorm"
)

//InitDictionary inits the dictionary for user token in the platform
func InitDictionary(db *gorm.DB) {
	l := log.NewLogger(0)
	aggDataset := bDict.NewDAgg(db, l)
	aggDict := bDict.NewDAgg(db, l)
	bDict.SetDefaultDatasetAggregator(aggDataset)
	interpreter.SetDefaultDICTAggregator(aggDict)
	defaultRules.LoadDefaultRules()
}
