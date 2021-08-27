import { Column } from "../../common/Type"
import { Relation, SaveRequest } from "../../common/Type"
import { users } from "./user"
import { Entity } from "../../entity/Entity"
import { Session } from "../../session/Session"
export class items extends Entity{
	
	constructor(data: {[key: string]: any}, sess: Session) {
		super()
        this._data = data;
        this.sess = sess;
		this.tableName = "public.items";
		this.pkcolumn = data["id"]
		this.pkcolumnName = "id"
    }
	export(): SaveRequest|null {
        if (this.isDirty) {
            let out = {dirtyData: this._dirtyData, pkcolumn: this.pkcolumn[0]}
			return out
		}
		return null
    }
	public get id (){
	return this._data.id;
	}
	public get uid (){
	return this._data.uid;
	}
	public get name (){
	return this._data.name;
	}
	public get secure (){
	return this._data.secure;
	}
	public get value (){
	return this._data.value;
	}
	public get control (){
	return this._data.control;
	}
	public set id(data: number[]){
	this._data.id = data
	}
	public set uid(data: number[]){
	this._data.uid = data
	}
	public set name(data: string[]){
	this._data.name = data
	}
	public set secure(data: string[]){
	this._data.secure = data
	}
	public set value(data: string[]){
	this._data.value = data
	}
	public set control(data: string[]){
	this._data.control = data
	}
	public get users(){
	 return this.relations[0]
	}
	private _data: any;
	private _dirtyData: { [key: string]: any } = {};
	static columns: Column[] = [{
	dataType: "number" ,
	name: "id" ,
	nullable: false
},{
	dataType: "number" ,
	name: "uid" ,
	nullable: false
},{
	dataType: "string" ,
	name: "name" ,
	nullable: false
},{
	dataType: "string" ,
	name: "secure" ,
	nullable: false
},{
	dataType: "string" ,
	name: "value" ,
	nullable: false
},{
	dataType: "string" ,
	name: "control" ,
	nullable: false
	}];
	private _data: any;
	relations: Relation[] = [
	users
	];
}