import { Context } from "../context/Context"
//import { users } from "./entity/user";
import { items } from "./entity/item";

(async function () {
    const ctx = new Context();
    ctx.registerEntity(items);
    const sess = ctx.getSession();
    let data = await sess.select(items).where("value=$1", "1").get_all()
    console.log(data[0].users.pkcolumn)
    console.log(await data[0].getname())
    //sess.save(data)
    //sess.delete(data)
})().catch(err => {
    console.log(err)
})
