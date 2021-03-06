package system

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"gmanager/utils/base"
)

type SysDepartment struct {
	// columns START
	Id        int    `json:"id" gconv:"id,omitempty"`                // 主键
	ParentId  int    `json:"parentId" gconv:"parent_id,omitempty"`   // 上级机构
	Name      string `json:"name" gconv:"name,omitempty"`            // 部门/11111
	Code      string `json:"code" gconv:"code,omitempty"`            // 机构编码
	Sort      int    `json:"sort" gconv:"sort,omitempty"`            // 序号
	Linkman   string `json:"linkman" gconv:"linkman,omitempty"`      // 联系人
	LinkmanNo string `json:"linkmanNo" gconv:"linkman_no,omitempty"` // 联系人电话
	Remark    string `json:"remark" gconv:"remark,omitempty"`        // 机构描述
	// columns END
	ParentName string `json:"parentName" gconv:"parentName,omitempty"` // 父节点名称

	base.BaseModel
}

func (model SysDepartment) Get() SysDepartment {
	if model.Id <= 0 {
		glog.Error(model.TableName() + " get id error")
		return SysDepartment{}
	}

	var resData SysDepartment
	err := model.dbModel("t").Where(" id = ?", model.Id).Fields(model.columns()).Struct(&resData)
	if err != nil {
		glog.Error(model.TableName()+" get one error", err)
		return SysDepartment{}
	}

	return resData
}

func (model SysDepartment) GetOne(form *base.BaseForm) SysDepartment {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["id"] != "" {
		where += " and id = ? "
		params = append(params, gconv.Int(form.Params["id"]))
	}
	if form.Params != nil && form.Params["parentId"] != "" {
		where += " and parent_id = ? "
		params = append(params, gconv.Int(form.Params["parentId"]))
	}

	var resData SysDepartment
	err := model.dbModel("t").Where(where, params).Fields(model.columns()).Struct(&resData)
	if err != nil {
		glog.Error(model.TableName()+" get one error", err)
		return SysDepartment{}
	}

	return resData
}

func (model SysDepartment) List(form *base.BaseForm) []SysDepartment {
	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["name"] != "" {
		where += " and name like ? "
		params = append(params, "%"+form.Params["name"]+"%")
	}

	var resData []SysDepartment
	err := model.dbModel("t").Fields(
		model.columns()).Where(where, params).OrderBy(form.OrderBy).Structs(&resData)
	if err != nil {
		glog.Error(model.TableName()+" list error", err)
		return []SysDepartment{}
	}

	return resData
}

func (model SysDepartment) Page(form *base.BaseForm) []SysDepartment {
	if form.Page <= 0 || form.Rows <= 0 {
		glog.Error(model.TableName()+" page param error", form.Page, form.Rows)
		return []SysDepartment{}
	}

	where := " 1 = 1 "
	var params []interface{}
	if form.Params != nil && form.Params["name"] != "" {
		where += " and t.name like ? "
		params = append(params, "%"+form.Params["name"]+"%")
	}

	num, err := model.dbModel("t").Where(where, params).Count()
	form.TotalSize = num
	form.TotalPage = num / form.Rows

	// 没有数据直接返回
	if num == 0 {
		form.TotalPage = 0
		form.TotalSize = 0
		return []SysDepartment{}
	}

	var resData []SysDepartment
	pageNum, pageSize := (form.Page-1)*form.Rows, form.Rows
	dbModel := model.dbModel("t").Fields(model.columns() + ",sdp.name parentName,su1.real_name as updateName,su2.real_name as createName")
	dbModel = dbModel.LeftJoin("sys_user su1", " t.update_id = su1.id ")
	dbModel = dbModel.LeftJoin("sys_user su2", " t.update_id = su2.id ")
	dbModel = dbModel.LeftJoin("sys_department sdp", " sdp.id = t.parent_id ")
	err = dbModel.Where(where, params).Limit(pageNum, pageSize).OrderBy(form.OrderBy).Structs(&resData)
	if err != nil {
		glog.Error(model.TableName()+" page list error", err)
		return []SysDepartment{}
	}

	return resData
}

func (model SysDepartment) Delete() int64 {
	if model.Id <= 0 {
		glog.Error(model.TableName() + " delete id error")
		return 0
	}

	r, err := model.dbModel().Where(" id = ?", model.Id).Delete()
	if err != nil {
		glog.Error(model.TableName()+" delete error", err)
		return 0
	}

	res, err2 := r.RowsAffected()
	if err2 != nil {
		glog.Error(model.TableName()+" delete res error", err2)
		return 0
	}

	LogSave(model, DELETE)
	return res
}

func (model *SysDepartment) Update() int64 {
	r, err := model.dbModel().Data(model).Where(" id = ?", model.Id).Update()
	if err != nil {
		glog.Error(model.TableName()+" update error", err)
		return 0
	}

	res, err2 := r.RowsAffected()
	if err2 != nil {
		glog.Error(model.TableName()+" update res error", err2)
		return 0
	}

	LogSave(model, UPDATE)
	return res
}

func (model *SysDepartment) Insert() int64 {
	r, err := model.dbModel().Data(model).Insert()
	if err != nil {
		glog.Error(model.TableName()+" insert error", err)
		return 0
	}

	res, err2 := r.RowsAffected()
	if err2 != nil {
		glog.Error(model.TableName()+" insert res error", err2)
		return 0
	} else if res > 0 {
		lastId, err2 := r.LastInsertId()
		if err2 != nil {
			glog.Error(model.TableName()+" LastInsertId res error", err2)
			return 0
		} else {
			model.Id = gconv.Int(lastId)
		}
	}

	LogSave(model, INSERT)
	return res
}

func (model SysDepartment) dbModel(alias ...string) *gdb.Model {
	var tmpAlias string
	if len(alias) > 0 {
		tmpAlias = " " + alias[0]
	}
	tableModel := g.DB().Table(model.TableName() + tmpAlias).Safe()
	return tableModel
}

func (model SysDepartment) PkVal() int {
	return model.Id
}

func (model SysDepartment) TableName() string {
	return "sys_department"
}

func (model SysDepartment) columns() string {
	sqlColumns := "t.id,t.parent_id as parentId,t.name,t.code,t.sort,t.linkman,t.linkman_no as linkmanNo,t.remark,t.enable,t.update_time as updateTime,t.update_id as updateId,t.create_time as createTime,t.create_id as createId"
	return sqlColumns
}
