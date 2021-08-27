import { QueryAST } from "./Queryast";
import { parse } from "../common/parser/parser";
import { Entity } from "../entity/Entity";
import { EntityTarget } from "../entity/EnitiyTarget";
import { Session } from "../session/Session";

export class Query<T extends Entity> {
    content: QueryAST;

    constructor(public entity: EntityTarget<T>, public sess: Session){
        this.sess = sess;
        this.entity = entity;
        this.content = {
            SessionId: "",
            Table: "public.items",
            Desc: false
        }
    }
    limit(n: number): Query<T> {
        this.content.Limit = n;
        return this;
    }
    offset(n: number): Query<T> {
        this.content.Offset = n;
        return this
    }
    where(expr: string, ...args: string[]): Query<T> {
        const ast = parse(expr);
        this.content.Where = ast;
        this.content.Args = args;
        return this;
    }
    orderBy(name: string[]): Query<T> {
        this.content.OrderBy = name;
        return;
    }

    get_by_pid(name: string) {
        return this.where("$1=$2", this.sess.ctx.EntityMap.get(this.entity.name) ,name).get_one()
    }

    async get_all(): Promise<any[]> {
        return await this.get_n(null);
    }
    async get_one(): Promise<any> {
        let res = await this.get_n(1)
        return res;
    }

    async get_n(n: number): Promise<any[]> {
        var entities = [];
        this.content.Limit = n;
        let res = await this.sess.query(this.content)
        console.log(res)
        for (let table of res.tables) {
            for (let row of table.rows) {
                var data = {}
                for (let i=0; i<row.length; i++) {
                    data[table.columns[i]] = row[i]
                }
                entities.push(new this.entity(data, this.sess))
            }
        }
        // const out = new this.entity(res.tables[0].data, this.sess);
        this.sess.entityMap.push(entities)
        return entities;
    }
    
}