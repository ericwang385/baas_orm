import { Entity } from "./Entity";
import { Session } from "session/Session";

export type EntityTarget<T extends Entity> = new(data: {[key: string]: any}, sess: Session) => T
