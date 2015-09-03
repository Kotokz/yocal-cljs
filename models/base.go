package models

// GetById is a ORM helper func to find object base on ID
func GetById(id int64, obj interface{}) error {
	has, err := orm.Id(id).Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

// GetByExample finds object base on the information in object, like name etc
func GetByExample(obj interface{}) error {
	has, err := orm.Get(obj)
	if err != nil {
		return err
	}
	if !has {
		return ErrNotExist
	}
	return nil
}

// Count counts object number base on partial object information
func Count(obj interface{}) (int64, error) {
	return orm.Count(obj)
}

// IsExist check if similar object already in db, to quickly find user name duplicate etc
func IsExist(obj interface{}) bool {
	has, _ := orm.Get(obj)
	return has
}

// Insert is a helper func to insert object into db
func Insert(obj interface{}) error {
	_, err := orm.Insert(obj)
	return err
}

// Find finds objects base on limit offset, to help pagination
func Find(limit, start int, objs interface{}) error {
	return orm.Limit(limit, start).Find(objs)
}

// DeleteById deletes object from DB base on id.
func DeleteById(id int64, obj interface{}) error {
	_, err := orm.Id(id).Delete(obj)
	return err
}

// DeleteByExample deletes object from DB base on object info
func DeleteByExample(obj interface{}) error {
	_, err := orm.Delete(obj)
	return err
}

func Obj2Table(objs []string) []string {
	var res = make([]string, len(objs))
	for i, c := range objs {
		res[i] = orm.ColumnMapper.Obj2Table(c)
	}
	return res
}

// UpdateById updates object from DB base on object ID
func UpdateById(id int64, object interface{}, cols ...string) error {
	_, err := orm.Cols(cols...).Id(id).Update(object)
	return err
}
