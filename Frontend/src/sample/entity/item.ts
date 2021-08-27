import { Column } from "../../common/Type"
import { Relation, SaveRequest } from "../../common/Type"
import { users } from "./user"
import { Entity } from "../../entity/Entity"
import { Session } from "../../session/Session"

export class items extends Entity {
	
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
	private _dirtyData: {[key: string]:any} = {};
	private _data: {[key: string]:any};
	public get id (){
	return this._data.id;
	}
	public get uid (){
	return this._data.uid;
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
	public set value(data: string){
		this.isDirty = true;
		this._data.value = data;
		this._dirtyData["value"] = data;
	}
	public set control(data: string[]){
	this._data.control = data
	}
	public get name (){
		throw new Error('lazy load column plz use getname instead');
	}
	async getname():Promise<string>{
		if(typeof this._data.name == "object") {
			return await this.sess.loadColmn(this,"name")
		} else {
			return this._data.name
		}
	}
	public set name(data: string[]){
	this._data.name = data
	}
	public get users(){
	return this.relations[0]
	}

	columns: Column[] = [{
	dataType: "number" ,
	name: "id" ,
	nullable: false
},{
	dataType: "number" ,
	name: "uid" ,
	nullable: false
},{
	dataType: "string" ,
	name: "value" ,
	nullable: false
},{
	dataType: "string" ,
	name: "control" ,
	nullable: false
},{
	dataType: "string" ,
	name: "name" ,
	nullable: false
	}];
	relations: Relation[] = [
	users
	];
}