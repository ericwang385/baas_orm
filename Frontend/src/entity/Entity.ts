import { Column } from "common/Type";
import { Relation, SaveRequest } from "common/Type";
import { Session } from "session/Session";

//definetion of one entity, each entity is a single row
export abstract class Entity {
    public sess: Session;
    public relations: Relation[];
    public isDirty: Boolean = false;
    public pkcolumn: any;
    public pkcolumnName: string;
    public tableName: string;
    static columns: Column[];
    
    save(){this.sess.save([this])}
	delete(){this.sess.delete([this])}
    abstract export(): SaveRequest|null
}