import { Entity } from "../entity/Entity";
import { EntityTarget } from "../entity/EnitiyTarget";
import { Query } from "../query/Query";
import { QueryAST } from "../query/Queryast";
import { Context } from "../context/Context";
import { QueryResponse } from "../common/Type";


export class Session {
    entityMap: Entity[][] = [];
    closed: boolean = false;
    axios = require('axios').default;

    constructor(public id: string, public ctx: Context) {
        
    }

    close(): void {
        this.closed = true;
    }

    save(e?: Entity[]): void {
        // if e exist save this specific entity
        // otherwise save the entire session
        if (typeof e === 'undefined') {
            for (const e of this.entityMap) {
                this._saveEntity(e)
                return
            }
        } else if (this.entityMap.includes(e)) {
            //console.log("save start")
            this._saveEntity(e)
            return
        } else {
            throw new Error('entity not found in current session');
        }
    }

    delete(e?: Entity[]): void {
        //same logic as save
        if (typeof e === 'undefined') {
            for (const e of this.entityMap) {
                this._deleteEntity(e)
                return
            }
        } else if (this.entityMap.includes(e)) {
            //console.log("delete start")
            this._deleteEntity(e)
            return
        } else {
            throw new Error("entity not found in current session")
        }
        return
    }

    select<T extends Entity>(entity: EntityTarget<T>): Query<T> {
        if (this.closed) {
            throw new Error('Session Closed');
        } else {
            return new Query<T>(entity, this);
        }
    }

    async loadColmn(entity: Entity, s: string): Promise<any> {
        //load for lazy cols, get data by pkcol and colName
        //console.log(entity.tableName)
        let res = await this.axios.post('http://127.0.0.1:8889/query/lazy', {
            table: entity.tableName,
            colName: s,
            primarykey: entity.pkcolumnName
        },
            { headers: { 'cookie': 'uid=1' } })
            
        return res.data.data
    }

    async query(q: QueryAST): Promise<QueryResponse> {
        q.SessionId = this.id
        let res = await this.axios.post('http://127.0.0.1:8889/query/select', q, { headers: { 'cookie': 'uid=1' } }).catch(err =>{
            console.log(err)
        })
        return res.data
    }

    private _saveEntity(entities: Entity[]): void {
        for (let entity of entities) {
            let row = entity.export()
            if (row !== null) {
                this.axios.post('http://127.0.0.1:8889/update', {
                primarykey: row.pkcolumn.toString(),
                colname:    Object.keys(row.dirtyData),
                value:      Object.values(row.dirtyData),
                table:      entity.tableName
            },{
                headers: {'cookie': 'uid=1'}
            }).then(function (res) {
                console.log(res)
            }).catch(function (err) {
                console.log(err)
            })
            } else {
                continue
            }
        }
        return
    }

    private _deleteEntity(entities: Entity[]): void {
        for (let entity of entities) {
            let row = entity.export()
            if (row !== null) {
                console.log(row)
                console.log(row.dirtyData)
                this.axios.post('http://127.0.0.1:8889/delete', {
                primarykey: row.pkcolumn.toString(),
                colname:    Object.keys(row.dirtyData),
                value:      Object.values(row.dirtyData),
                table:      entity.tableName
            },{
                headers: {'cookie': 'uid=1'}
            }).then(function (res) {
                console.log(res)
            }).catch(function (err) {
                console.log(err)
            })
            } else {
                continue
            }
        }
    }
}