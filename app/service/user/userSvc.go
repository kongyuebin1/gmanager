package user

import (
	"errors"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"gmanager/app/model/user"
	"gmanager/utils/base"
)

// 请求参数
type Request struct {
	user.Entity
}

func GetById(id int64) (*user.Entity, error) {
	if id <= 0 {
		glog.Error(" get id error")
		return new(user.Entity), errors.New("参数不合法")
	}

	return user.Model.FindOne(" id = ?", id)
}

func GetOne(form *base.BaseForm) (*user.Entity, error) {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["id"] != "" {
		where += " and id = ? "
		params = append(params, gconv.Int(form.Params["id"]))
	}

	return user.Model.FindOne(where, params)
}

func List(form *base.BaseForm) ([]*user.Entity, error) {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["name"] != "" {
		where += " and name like ? "
		params = append(params, "%"+form.Params["name"]+"%")
	}

	return user.Model.Order(form.OrderBy).FindAll(where, params)
}

func Delete(id int64) (int64, error) {
	if id <= 0 {
		glog.Error("delete id error")
		return 0, errors.New("参数不合法")
	}

	r, err := user.Model.Delete(" id = ?", id)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func Update(request *Request) (int64, error) {
	entity := (*user.Entity)(nil)
	err := gconv.StructDeep(request.Entity, &entity)
	if err != nil {
		return 0, errors.New("数据错误")
	}

	if entity.Id <= 0 {
		glog.Error("update id error")
		return 0, errors.New("参数不合法")
	}

	r, err := user.Model.Where(" id = ?", entity.Id).Update(entity)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

func Insert(request *Request) (int64, error) {
	entity := (*user.Entity)(nil)
	err := gconv.StructDeep(request.Entity, &entity)
	if err != nil {
		return 0, errors.New("数据错误")
	}

	if entity.Id > 0 {
		glog.Error("insert id error")
		return 0, errors.New("参数不合法")
	}

	r, err := user.Model.Insert(entity)
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}

// 分页查询
func Page(form *base.BaseForm) ([]user.Entity, error) {
	if form.Page <= 0 || form.Rows <= 0 {
		glog.Error("page param error", form.Page, form.Rows)
		return []user.Entity{}, nil
	}

	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["operObject"] != "" {
		where += " and oper_object like ? "
		params = append(params, "%"+form.Params["operObject"]+"%")
	}
	if form.Params != nil && form.Params["operTable"] != "" {
		where += " and oper_table like ? "
		params = append(params, "%"+form.Params["operTable"]+"%")
	}
	if form.Params != nil && gconv.Int(form.Params["logType"]) > 0 {
		where += " and log_type = ? "
		params = append(params, gconv.Int(form.Params["logType"]))
	}
	if form.Params != nil && gconv.Int(form.Params["operType"]) > 0 {
		where += " and oper_type = ? "
		params = append(params, gconv.Int(form.Params["operType"]))
	}

	num, err := user.Model.As("t").FindCount(where, params)
	form.TotalSize = num
	form.TotalPage = num / form.Rows

	if err != nil {
		glog.Error("page count error", err)
		return []user.Entity{}, err
	}

	// 没有数据直接返回
	if num == 0 {
		form.TotalPage = 0
		form.TotalSize = 0
		return []user.Entity{}, err
	}

	var resData []user.Entity
	dbModel := user.Model.As("t").Fields(user.Model.Columns() + ",su1.real_name as updateName,su2.real_name as createName")
	dbModel = dbModel.LeftJoin("sys_user su1", " t.update_id = su1.id ")
	dbModel = dbModel.LeftJoin("sys_user su2", " t.update_id = su2.id ")
	err = dbModel.Where(where, params).Order(form.OrderBy).Page(form.Page, form.Rows).M.Structs(&resData)
	if err != nil {
		glog.Error("page list error", err)
		return []user.Entity{}, err
	}

	return resData, nil
}