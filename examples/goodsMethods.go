package examples

import (
	"errors"

	"../arca"
)

// Good whatever
type Good struct {
	ID          int
	Description string
	Price       int
}

// Goods whatever
type Goods []Good

// This is my Data-Base!
var goods = Goods{
	Good{1, "Computer", 1000},
	Good{2, "Smartphone", 2000},
	Good{3, "Wine", 1500},
}
var lastGoodsID = len(goods)

var goodsCRUD = arca.DIRUD{
	Read: func(requestParams *interface{}, context *interface{},
		response chan interface{}) error {
		go (func() { response <- goods })()
		return nil
	},
	Update: func(requestParams *interface{}, context *interface{},
		response chan interface{}) error {
		params := (*requestParams).(map[string]interface{})
		preid, ok := params["ID"]
		if !ok {
			return errors.New("params in request doesn't contain ID")
		}
		preid2, ok := preid.(float64)
		if !ok {
			return errors.New("ID in params isn't int")
		}

		id := int(preid2)
		for index, good := range goods {
			if good.ID == id {
				if description, ok := params["Description"]; ok {
					goods[index].Description = description.(string)
				}
				if price, ok := params["Price"]; ok && price != nil {
					preprice := price.(float64)
					goods[index].Price = int(preprice)
				}
				go (func() { response <- goods[index] })()
				return nil
			}
		}
		return errors.New("nothing")
	},
	Insert: func(requestParams *interface{}, context *interface{},
		response chan interface{}) error {
		params := (*requestParams).(map[string]interface{})
		lastGoodsID++
		newGood := Good{ID: lastGoodsID}
		if description, ok := params["Description"]; ok {
			newGood.Description = description.(string)
		}
		if price, ok := params["Price"]; ok && price != nil {
			preprice := price.(float64)
			newGood.Price = int(preprice)
		}
		goods = append(goods, newGood)
		go (func() { response <- newGood })()
		return nil
	},
	Delete: func(requestParams *interface{}, context *interface{},
		response chan interface{}) error {
		params := (*requestParams).(map[string]interface{})
		preid, ok := params["ID"]
		if !ok {
			return errors.New("params in request doesn't contain ID")
		}
		preid2, ok := preid.(float64)
		if !ok {
			return errors.New("ID in params isn't int")
		}

		id := int(preid2)
		deletedGood := Good{ID: id}
		for i, good := range goods {
			if good.ID == id {
				goods = append(goods[:i], goods[i+1:]...)
				go (func() { response <- deletedGood })()
				return nil
			}
		}
		return errors.New("nothing")
	},
}
