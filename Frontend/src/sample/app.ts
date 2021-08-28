import { Context } from "../context/Context"
//import { users } from "./entity/user";
import { items } from "./entity/item";

(async function () {
    //生成context，注册Table对象
    const ctx = new Context();
    ctx.registerTable(items);
    const sess = ctx.getSession();
    //简单select
    let data = await sess.select(items).where("value=$1", "1").get_all()
    console.log(data[0].users.pkcolumn)
    console.log(await data[0].getname())
    //保存和删除删除
    //sess.save(data)
    //sess.delete(data)
})().catch(err => {
    console.log(err)
})
