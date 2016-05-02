package orm

import (
	conn "github.com/lnx1337/mongorm/config"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/lnx1337/go/api"
	util "github.com/lnx1337/go/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"golang.org/x/net/context"
	"github.com/guregu/db"
)

type Orm struct {
	Model         interface{}
	PkStructField string
	PkBson        string
}

var Collection *mgo.Collection
var Sess *mgo.Session

func NewOrmWithContext(ctx  context.Context, model interface{}) Orm {
	self := Orm{Model: model}
	val := reflect.ValueOf(model)

	d  := self.GetModel("DbName", val)
	c  := self.GetModel("CollectionName", val)

	Sess = db.MongoDB(ctx,"main")
	Collection = Sess.DB(d).C(c)

	return self
}

func NewOrm(model interface{}) Orm {
	var err error
	self := Orm{Model: model}
	val := reflect.ValueOf(model)

	conn.Db = self.GetModel("DbName", val)
	conn.Col = self.GetModel("CollectionName", val)
	
	conn.InitDb()
	Sess = conn.Sess().Copy()

	Collection = Sess.DB(conn.Db).C(conn.Col)

	if err != nil {
		fmt.Println(err.Error())
	}
	return self
}

func (self *Orm) SetPK(structField string, bsonTag string) {
	self.PkStructField = structField
	self.PkBson = bsonTag
}

func (self *Orm) getPkValue() interface{} {
	var pkValue interface{}
	if len(self.PkStructField) == 0 {
		self.PkBson = "_id"
		pkValue = reflect.ValueOf(self.Model).Elem().FieldByName("ID").Interface().(bson.ObjectId)
	} else {
		pkValue = reflect.ValueOf(self.Model).Elem().FieldByName(self.PkStructField).Interface().(string)
	}
	return pkValue
}

// TODO:
// FindAll
// map[string]interface{}
// Sort()

func (self *Orm) FindById(id string) (interface{}, *api.Err) {
	errors := api.NewError()
	err := Collection.Find(
		bson.M{
			self.PkBson: id,
		},
	).One(self.Model)
	if err != nil {
		errors.SetErr(
			err.Error(),
			api.ErrMongoSelectAll,
			17,
		)
	}
	return self.Model, errors
}

func (self *Orm) FindByCondition(model interface{}, cond interface{}) *api.Err {
	errors := api.NewError()
	err := Collection.Find(cond).All(model)
	if err != nil {
		errors.SetErr(
			err.Error(),
			api.ErrMongoSelectAll,
			17,
		)
	}
	return errors
}

func (self *Orm) FindByPk() (interface{}, *api.Err) {
	errors := api.NewError()
	id := self.getPkValue()
	err := Collection.Find(
		bson.M{
			self.PkBson: id,
		},
	).One(self.Model)
	if err != nil {
		errors.SetErr(
			err.Error(),
			api.ErrMongoSelectAll,
			17,
		)
	}
	return self.Model, errors
}

func (self *Orm) Save() *api.Err {
	errors := api.NewError()
	err := Collection.Insert(self.Model)
	if err != nil {
		errors.SetErr(
			err.Error(),
			api.ErrMongoInsert,
			18,
		)
	}
	return errors
}

func (self *Orm) Update() *api.Err {
	errors := api.NewError()

	data, _ := util.ToMap(self.Model, "json")

	pkValue := self.getPkValue()

	err := Collection.Update(
		bson.M{
			self.PkBson: pkValue,
		},
		bson.M{
			"$set": data,
		},
	)

	if err != nil {
		errors.SetErr(
			err.Error(),
			api.ErrMongoUpdate,
			19,
		)
	}
	return errors
}

func (self *Orm) UpdateAllByConditions(a, b bson.M) *api.Err {
	errors := api.NewError()
	_, err := Collection.UpdateAll(a, b)
	if err != nil {
		errors.SetErr(
			err.Error(),
			api.ErrMongoUpdate,
			19,
		)
	}
	return errors
}

func (self *Orm) Delete() *api.Err {
	errors := api.NewError()
	id := self.getPkValue()
	err := Collection.Remove(
		bson.M{
			self.PkBson: id,
		},
	)
	if err != nil {
		errors.SetErr(
			err.Error(),
			api.ErrMongoDelete,
			20,
		)
	}
	return errors
}

func (self *Orm) Close() {
	conn.Sess().Close()
}

func (self *Orm) Validate() *api.Err {
	errors := api.NewError()
	valid := validation.Validation{}
	b, _ := valid.Valid(self.Model)
	if !b {
		for _, err := range valid.Errors {
			errors.SetErr(err.Key, err.Message, 0)
		}
	}
	return errors
}

func (self *Orm) GetModel(nameMethod string, val reflect.Value) string {
	ind := reflect.Indirect(val)
	fun := val.MethodByName(nameMethod)
	if fun.IsValid() {
		vals := fun.Call([]reflect.Value{})
		if len(vals) > 0 {
			val := vals[0]
			if val.Kind() == reflect.String {
				return val.String()
			}
		}
	}
	return self.snakeString(ind.Type().Name())
}

func (self *Orm) snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	for i, d := range s {
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		j = (d != '_')
		data = append(data, byte(d))
	}
	return strings.ToLower(string(data[:len(data)]))
}
