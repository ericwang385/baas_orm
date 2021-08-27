package table

import (
	"database/sql"
	"errors"
	"feorm/config"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func NewTable(schema string, name string) error {
	var f *os.File
	filename := name + ".ts"
	table := TableMap[strings.Join([]string{schema, name}, ".")]
	code, err := CodeGen(table)
	if err != nil {
		return err
	}
	f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	//err = f.Truncate(0)
	//if err != nil {
	//	return err
	//}
	_, err = f.Write([]byte(code))
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

func CodeGen(table *Table) (string, error) {
	imports := `import { Column } from "../../common/Type"
import { Relation, SaveRequest } from "../../common/Type"
import { users } from "./user"
import { Entity } from "../../entity/Entity"
import { Session } from "../../session/Session"`
	body, err := GenBody(table)
	if err != nil {
		return "", err
	}
	code := strings.Join([]string{imports, "export class " + table.Name + " extends Entity{", body, "}"}, "\n")
	return code, nil
}

func GenBody(table *Table) (string, error) {
	cons := GenConstructor()
	funs, err := GenMethods(table)
	if err != nil {
		return "", err
	}
	props, err := GenProps(table)
	if err != nil {
		return "", err
	}
	//relFuns, err := GenRelMethods(table)
	//if err != nil {
	//	return "", nil
	//}
	//relProps := GenRelProps(table.Relations)
	out := strings.Join([]string{"\t", cons, funs, props}, "\n\t")
	return out, nil
}

func GenConstructor() string {
	consStr := `constructor(data: {[key: string]: any}, sess: Session) {
		super()
        this._data = data;
        this.sess = sess;
		this.tableName = "public.items";
		this.pkcolumn = data["id"]
		this.pkcolumnName = "id"
    }`
	return consStr
}

func GenMethods(table *Table) (string, error) {
	// methods behavior may vary by auth
	static := []string{
		`export(): SaveRequest|null {
        if (this.isDirty) {
            let out = {dirtyData: this._dirtyData, pkcolumn: this.pkcolumn[0]}
			return out
		}
		return null
    }`,
	}
	getter := make([]string, 0)
	setter := make([]string, 0)
	lazyGetter := make([]string, 0)
	lazySetter := make([]string, 0)
	for _, lazyCol := range table.lazyColumnInfo {
		getterStr, err := GenLazyGetter(lazyCol)
		if err != nil {
			return "", err
		}
		setterStr, err := GenSetter(lazyCol)
		if err != nil {
			return "", err
		}
		lazyGetter = append(lazyGetter, getterStr)
		lazySetter = append(lazySetter, setterStr)
	}
	for _, col := range table.columnInfo {
		tmpGet, err := GenGetter(col)
		if err != nil {
			return "", err
		}
		tmpSet, err := GenSetter(col)
		if err != nil {
			return "", err
		}
		getter = append(getter, tmpGet)
		setter = append(setter, tmpSet)
	}
	tableGetter := GenRelMethods(table.Relations)

	outStr := append(static, getter...)
	outStr = append(outStr, setter...)
	outStr = append(outStr, lazyGetter...)
	outStr = append(outStr, lazySetter...)
	outStr = append(outStr, tableGetter...)

	out := strings.Join(outStr, "\n\t")
	return out, nil
}

func GenRelMethods(rels []config.Relation) []string {
	// currently no setter for relation
	getter := make([]string, 0)
	for i, rel := range rels {
		fun := fmt.Sprintf("public get %s(){\n\t return this.relations[%d]\n\t}", rel.ForeignName, i)
		getter = append(getter, fun)
	}
	return getter
}

func GenProps(table *Table) (string, error) {
	var out string
	staticPropsStr := GenStaticProps()
	colPropsStr, err := GenColumnProps(table.allColumnInfo)
	if err != nil {
		return "", err
	}
	if len(table.TableDefine.Relations) > 0 {
		relPropStr := GenRelProps(table.TableDefine.Relations)
		out = strings.Join([]string{staticPropsStr, colPropsStr, relPropStr}, "\n\t")
	} else {
		out = strings.Join([]string{staticPropsStr, colPropsStr}, "")
	}
	return out, nil
}

func GenStaticProps() string {
	dataStr := "private _data: any;"
	dirtyData := "private _dirtyData: { [key: string]: any } = {};"
	staticProps := strings.Join([]string{dataStr, dirtyData}, "\n\t")
	return staticProps
}

func GenColumnProps(info []sql.ColumnType) (string, error) {
	colprops := make([]string, 0)
	columns := make([]string, 0)

	for _, colT := range info {
		generatedType, err := TypeConvert(colT.DatabaseTypeName())
		if err != nil {
			return "", err
		}
		datatype := strings.Join([]string{"dataType: ", generatedType, " "}, "\"")
		name := strings.Join([]string{"name: ", colT.Name(), " "}, "\"")
		isNullable, _ := colT.Nullable()
		nullable := "nullable: " + strconv.FormatBool(isNullable)
		columns = append(columns, strings.Join([]string{datatype, name, nullable}, ",\n\t"))
		colprops = append(colprops, strings.Join([]string{colT.Name(), ": ", generatedType, "[]", ";"}, ""))
	}
	out := strings.Join([]string{
		"static columns: Column[] = [{",
		strings.Join(columns, "\n},{\n\t"),
		"}];",
	}, "\n\t")
	return out, nil
}

func GenRelProps(rels []config.Relation) string {
	relStr := make([]string, 0)
	for _, rel := range rels {
		relStr = append(relStr, rel.ForeignName)
	}
	out := strings.Join([]string{
		"relations: Relation[] = [",
		strings.Join(relStr, ","),
		"];"}, "\n\t")
	return out
}

func GenGetter(info sql.ColumnType) (string, error) {
	getterDecl := strings.Join([]string{"public get", info.Name(), "(){"}, " ")
	getterBody := strings.Join([]string{
		"return this._data.",
		info.Name(),
		";"}, "")
	out := strings.Join([]string{getterDecl, getterBody, "}"}, "\n\t")
	return out, nil
}

func GenSetter(info sql.ColumnType) (string, error) {
	generatedType, err := TypeConvert(info.DatabaseTypeName())
	if err != nil {
		return "", err
	}
	setterDecl := strings.Join([]string{"public set ", info.Name(), "(data: ", generatedType, "[]){"}, "")
	setterBody := strings.Join([]string{
		"this._data.",
		info.Name(),
		" = data"}, "")
	out := strings.Join([]string{setterDecl, setterBody, "}"}, "\n\t")
	return out, nil
}

func GenLazyGetter(info sql.ColumnType) (string, error) {
	generatedType, err := TypeConvert(info.DatabaseTypeName())
	if err != nil {
		return "", err
	}
	getterDecl := strings.Join([]string{"public get", info.Name(), "(){"}, " ")
	getterError := "new Error('lazy load column plz use get" + info.Name() + " instead')"
	getterBody := strings.Join([]string{
		"return this._data.",
		info.Name(),
		";"}, "")
	getterCond := strings.Join([]string{"if(this._data." + info.Name() + "==\"object\"){", getterError, "}", "else{", getterBody, "}"}, "\n\t")
	asyncDecl := strings.Join([]string{"async get", info.Name(), "():Promise<", generatedType + ">{"}, "")
	asyncReturn := strings.Join([]string{"return await this.sess.loadColmn(this,\"", info.Name(), "\")"}, "")
	asyncCond := strings.Join([]string{
		"if(this._data." + info.Name() + "==\"object\"){",
		asyncReturn,
		"}else{",
		getterBody,
	}, "\n\t")
	out := strings.Join([]string{getterDecl, getterCond, "}", asyncDecl, asyncCond, "}"}, "\n\t")
	return out, nil
}

func TypeConvert(dbType string) (string, error) {
	switch dbType {
	case "INT4":
		return "number", nil
	case "TEXT":
		return "string", nil
	default:
		return "", errors.New("CODEGEN ERROR: TypeConvert failed, dbType mismatch")
	}
}
