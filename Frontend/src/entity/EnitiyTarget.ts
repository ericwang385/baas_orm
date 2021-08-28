import { Entity } from "./Entity";
import { Session } from "session/Session";
//Entity类的抽象
export type EntityTarget<T extends Entity> = new(data: {[key: string]: any}, sess: Session) => T
