import { Session } from "../session/Session";
import { Entity } from "entity/Entity";
import { EntityTarget } from "../entity/EnitiyTarget";

export class Context {
    TableMap: Map<string, any> = new Map();
    sessionMap: Map<string, any> = new Map();
    name: string = "admin";

    constructor(name?: string) {
        this.name = name;
    }

    registerTable<T extends Entity>(entityClass: EntityTarget<T>){
        //导入table对象
        this.TableMap.set(entityClass.name, entityClass);
    }

    getSession(): Session {
        //初始化Session
        const { nanoid } = require('nanoid');
        const id = nanoid(12);
        const sess = new Session(id, this)
        this.sessionMap.set(id, sess);
        return sess;
    }

    getSessionbyId(id: string): Session {
        //通过session id获取已有session
        try{
            const sess = this.sessionMap.get(id);
            return sess;
        } catch(e) {
            throw new Error('session not found');
        }
    }
}