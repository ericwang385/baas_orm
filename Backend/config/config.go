package config

import "encoding/json"

type Config struct {
	Variables []VariableDefine
	Tables    []TableDefine
}

// VariableDefine 变量定义，可以使用上下文中预制的一些变量通过sql查询得到，可以配置有效时间
type VariableDefine struct {
	Name string
	//应该得到一列，最多1行
	Sql string
	//默认值
	Default interface{}
	//过期时间秒，0不会缓存
	Expire int
	//todo 是否分布式一致，不缓存的情况下无效
	Consist bool
}

type Relation struct {
	//都应该为一对多或一对一关系，反向引用会自动生成，目前支持等值条件
	ForeignName   string
	ForeignSchema string
	LocalColumn   string
	ForeignColumn string
	// 可选，防止和列名字重复
	Name string
	// 自动会在子表上生成反向引用，可选，防止和列名字重复
	ReverseName string
	// 默认预读取
	Preload bool
	// 反向引用默认预读取
	ReversePreload bool
}

type TableDefine struct {
	Schema string
	Name   string
	//行级基本控制，条件会编译为sql，查询时生效
	// "isAdmin=@isAdmin and dept=@dept" 会嵌入isAdmin和dept变量，然后查询时作为filter条件嵌入
	SelectGuard string
	UpdateGuard string
	DeleteGuard string
	// insert 之前会进行检查
	InsertGuard   string
	HiddenColumns []string
	//列权限控制
	ColumnRules      []ColumnRule
	GeneratedColumns map[string]string
	//关联实体
	Relations []Relation
	//延迟加载的字段，建议不常用的长字段使用，以提高性能
	LazyColumns []string
}

type ColumnRule struct {
	Columns    []string
	Operations []string
	// 匹配表达式，无表达式代表所有
	// 如果表达式引用了数据库中的列，那么会先执行查询，然后计算权限
	// 更新或者插入会同样验证更改后的数据
	// "@isAdmin=true or userid=@userid"
	Match string
	// deny
	Action string
}

func LoadConfigJson(data string) (Config, error) {
	var config Config
	err := json.Unmarshal([]byte(data), &config)
	return config, err
}
