import {Column} from "./types/Column";
import {Relation} from "./types/Relation";
import {users} from "./users"

export class items {
    constructor(data: {[key: string]: any}, sess: Session) {
        this._data = data
        this.sess = sess
    }

    save() {
        if (this.isDirty === true) {
            //返回data
        }
        return
    }

    delete() {

    }

    public get id() {
        if (this._data['id'] && this._data['id'].invalid === "lazy") {
            throw new Error(`prelod use async getId()`)
        } else {
            return this._id;
        }

    }
    async getId():Promise<number>{
    	return this._data['id']=this.session.loadColmn(this,"id")
	}

    public get uid() {
        return this._uid;
    }

    public get name() {
        return this._name;
    }

    public get value() {
        return this._value;
    }

    public get control() {
        return this._control;
    }

    public set id(data: number []) {
        this._id = data
    }

    public set uid(data: number []) {
        this._uid = data
    }

    public set name(data: string []) {
        this._name = data
    }

    public set value(data: string []) {
        this._value = data
    }

    public set control(data: string []) {
        this._control = data
    }

    pkcolumn = [""];

    columns: Column[] = [{
        datatype: "number",
        colname: "id",
        nullable: false
    }, {
        datatype: "number",
        colname: "uid",
        nullable: false
    }, {
        datatype: "string",
        colname: "name",
        nullable: false
    }, {
        datatype: "string",
        colname: "value",
        nullable: false
    }, {
        datatype: "string",
        colname: "control",
        nullable: false
    }];
    public isDirty: Boolean
    private dirtyColumn: string[]
    private dirtyData: { [key: string]: any }
    private _data: { [key: string]: any }
    private _id: number[];
    private _uid: number[];
    private _name: string[];
    private _value: string[];
    private _control: string[];
    relations: Relation[] = [
        users
    ];
}