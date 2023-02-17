package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetFilter(params map[string]interface{}) (interface{}, *options.FindOptions, error) {
	var setElements bson.A
	choice := options.FindOptions{}

	for i, v := range params {
		if i == "sortorder" || i == "sortfield" {
			continue
		}
		if i == "dateStart" {
			setElements = append(setElements, bson.D{{"dateStart", bson.D{{"$gte", v}}}})
			continue
		}
		if i == "dateEnd" {
			setElements = append(setElements, bson.D{{"dateEnd", bson.D{{"$lte", v}}}})
			continue
		}
		if i == "pausePeriod.dateStart" {
			setElements = append(setElements, bson.D{{"pausePeriod.dateStart", bson.D{{"$gte", v}}}})
			continue
		}
		if i == "pausePeriod.dateEnd" {
			setElements = append(setElements, bson.D{{"pausePeriod.dateEnd", bson.D{{"$lte", v}}}})
			continue
		}
		setElements = append(setElements, bson.M{i: v})
	}

	v, k := params["sortfield"]
	if !k {
		v = "meta.created"
	}
	if val, ok := params["sortorder"]; ok {
		if val == "ascending" {
			choice.Sort = bson.D{{v.(string), 1}}
		}
	} else {
		choice.Sort = bson.D{{v.(string), -1}}
	}

	filter := bson.D{{}}
	if len(setElements) > 0 {
		filter = bson.D{{"$and", setElements}}
	}

	return filter, &choice, nil
}

