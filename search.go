package database

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/zlt-com/go-common"
)

// Find Find
func (m *Database) Find() (replys interface{}, count int, err error) {
	if m.where == nil {
		// m.where = map[string]interface{}{}
		m.where = make([]Condition, 0)
	}

	//反射数据模型，x为查询结果集
	t := reflect.TypeOf(m.model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	x := reflect.New(slice.Type())
	x.Elem().Set(slice)

	//使用缓存
	if UseRedcache {
		rcache := new(RedCache).Where(m.whereConditions...).Limit(m.limit).Offset(m.offset).OrderBy(m.order).Model(m.model)
		if count, cacheCount, err := m.Count(); err != nil {
			return nil, count, err
		} else {
			if cacheCount < count {
				return m.syncData(rcache, x, slice, t)
			} else {
				if cacheValue, err := rcache.Select(); err == nil {
					if cacheValue == nil {
						return m.syncData(rcache, x, slice, t)
					} else {
						// replys := []interface{}{}
						replys := []reflect.Value{}
						for _, v := range cacheValue {
							if err = common.Byte2Object(v.([]byte), &m.model); err != nil {
								return nil, count, err
							} else {
								// n.Set(reflect.ValueOf(&model))
								clone := common.DeepClone(m.model)
								replys = append(replys, reflect.ValueOf(clone).Elem())
							}
						}
						r := reflect.New(slice.Type())
						r = reflect.Append(r.Elem(), replys...)
						return r.Interface(), count, nil
					}
				}
			}

		}

	} else {
		return m.fromDB(x, slice, t)
	}
	return nil, 0, nil
}

// Count Count
func (m *Database) Count() (count, cacheCount int, err error) {
	// return 1643372, 0, nil
	rcache := new(RedCache).Model(m.model).Where(m.whereConditions...)
	db := defaultDB
	tableName := common.ReflectMethod(m.model, "TableName")[0].Interface().(string)
	for _, where := range m.whereConditions {
		db = db.Where(where.ToString(tableName), where.Value)
	}
	if UseRedcache {
		if SyncCache, err = rcache.GetSyncStatus(); err == nil {
			if !SyncCache[common.ReflectInterfaceName(m.model)] {
				if err := db.Model(m.model).Count(&count); err.Error != nil {
					fmt.Println(err.Error)
					return 0, 0, err.Error
				}
				if cacheCount, err = rcache.Count(); err != nil {
					return count, cacheCount, err
				} else {
					if count == cacheCount {
						SyncCache[common.ReflectInterfaceName(m.model)] = true
					}
					rcache.SetSyncStatus(SyncCache)
					return
				}

			} else {
				cacheCount, err := rcache.Count()
				rcache.SetSyncStatus(SyncCache)
				return cacheCount, cacheCount, err
			}
		} else {
			return 0, 0, err
		}
	} else {
		if err := db.Model(m.model).Count(&count); err.Error != nil {
			fmt.Println(err.Error)
			return 0, 0, err.Error
		}
		return count, cacheCount, nil
	}

}

// Update Update
func (m *Database) Update(u interface{}) (i interface{}, err error) {
	result := defaultDB.Model(u).Update(u)
	if result.Error != nil {
		return nil, result.Error
	} else {
		if UseRedcache {
			rcache := new(RedCache).Model(result.Value).Model(u)
			if err = rcache.Create(); err != nil {
				return nil, err
			}
		}
	}
	return &result.Value, err
}

// Create Create
func (m *Database) Create(u interface{}) (i interface{}, err error) {
	result := defaultDB.Model(u).Create(u)
	if result.Error != nil {
		return nil, result.Error
	}
	if UseRedcache {
		rcache := new(RedCache).Model(result.Value).Model(u)
		if err = rcache.Create(); err != nil {
			fmt.Println(err)
		}
	}

	return result.Value, err
}

// Delete Delete
func (m *Database) Delete(u interface{}) (b bool, err error) {
	if UseRedcache {
		rcache := new(RedCache).Model(u)
		if err = rcache.Delete(); err != nil {
			fmt.Println(err)
			return false, err
		}
	}
	if err := defaultDB.Delete(m); err.Error != nil {
		fmt.Println(err.Error)
		return false, err.Error
	}
	return true, nil
}

//从数据库查询
func (m *Database) fromDB(x, slice reflect.Value, t reflect.Type) (interface{}, int, error) {
	defer timeMeasurement("fromDB", time.Now())
	db := defaultDB
	tableName := common.ReflectMethod(m.model, "TableName")[0].Interface().(string)
	// defer db.Close()
	// for _, where := range m.whereConditions {
	// 	db = db.Where(where.ToString(tableName), where.Value)
	// }
	count := 0
	if err := db.Model(m.model).Count(&count); err.Error != nil {
		fmt.Println(err.Error)
	}
	// ids := make([]interface{}, 0)

	// db.Table(tableName).Select("id").Scan(ids)
	// result := db.Order(m.order).Offset(m.offset).Limit(m.limit).Model(m.model).Select("id").Find(&ids)
	// if result.Error != nil {
	// 	fmt.Println(result.Error)
	// }
	// if result := db.Where(ids).Find(x.Interface()); result.Error != nil {
	// tableName := common.ReflectInterfaceName(m.model)
	joinSqlHeader := fmt.Sprintf("SELECT * FROM %s inner join (SELECT id FROM %s where 1=1 ", tableName, tableName)
	joinSqlFooter := fmt.Sprintf(") b  on %s.id=b.id", tableName)
	joinSqlContent := ""
	whereSlice := make([]interface{}, 0)
	for _, where := range m.whereConditions {
		joinSqlContent += " and " + where.ToString(tableName)
		whereSlice = append(whereSlice, where.Value)
	}
	joinSqlContent += " ORDER BY id desc "
	if m.limit != nil {
		joinSqlContent += fmt.Sprintf(" LIMIT %d ", m.limit)
	}
	if m.offset != nil {
		joinSqlContent += fmt.Sprintf(" OFFSET %d ", m.offset)
	}
	joinSql := joinSqlHeader + joinSqlContent + joinSqlFooter
	// ids := make([]int, 0)
	// db.Order(m.order).Offset(m.offset).Limit(m.limit).Model(m.model).Select("id").Find(x.Interface()).Pluck("id", &ids)
	// if result := db.Table(tableName).Select("*").Where("id in (?)", ids).Find(x.Interface()); result.Error != nil {
	// if result := db.Table(tableName).Select("*").Joins(joinSql).Find(x.Interface()); result.Error != nil {
	// fmt.Println(joinSql)
	if result := db.Raw(joinSql, whereSlice...).Find(x.Interface()); result.Error != nil {
		return nil, 0, result.Error

	} else {
		//获得数据后反射返回结果
		if result.Value != nil {
			getValue := reflect.ValueOf(result.Value)
			if getValue.Kind() == reflect.Ptr || getValue.Kind() == reflect.Slice {
				replys := []reflect.Value{}
				for i := 0; i < getValue.Elem().Len(); i++ {
					replys = append(replys, getValue.Elem().Index(i).Convert(t))
				}
				r := reflect.New(slice.Type())
				r = reflect.Append(r.Elem(), replys...)
				return r.Interface(), count, nil
			}

		}
		return nil, count, errors.New("nil")
	}
}

//同步数据库与缓存数据
func (m *Database) syncData(rcache *RedCache, x, slice reflect.Value, t reflect.Type) (interface{}, int, error) {
	defer timeMeasurement("syncData", time.Now())
	if dbValue, count, err := m.fromDB(x, slice, t); err == nil {
		getValue := reflect.ValueOf(dbValue)
		batchCacheCreateArray := []interface{}{}
		if getValue.Kind() == reflect.Slice {
			for i := 0; i < getValue.Len(); i++ {
				gv := getValue.Index(i).Interface()
				rcache = rcache.Model(gv)
				batchCacheCreateArray = append(batchCacheCreateArray, gv)
			}
			if err := rcache.BatchCreate(batchCacheCreateArray); err != nil {
				return nil, count, err
			}
		}
		return dbValue, count, err
	} else {
		return nil, 0, err
	}
}

func timeMeasurement(name string, start time.Time) {
	// elapsed := time.Since(start)
	// fmt.Printf(name+" Execution time: %s。\n", elapsed)
}
