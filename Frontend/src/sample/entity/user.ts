import { Entity } from "../../entity/Entity";
import { Column, SaveRequest } from "../../common/Type";


export class users extends Entity {

    // constructor(data: {[key: string]: any}, sess: Session) {
	// 	super()
    //     this._data = data;
    //     this.sess = sess;
	// 	this.tableName = "public.items";
    // }
    export(): SaveRequest|null{
        return
    }
    pkcolumn = ["id"];
    columns: Column[] = [{
        dataType: "number",
        name: "age",
        nullable: false,
    },{
        dataType: "string",
        name: "firstName",
        nullable: false,
    },{
        dataType: "string",
        name: "lastName",
        nullable: false,
    }];
}