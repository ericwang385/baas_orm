import {parse} from "../src/common/parser/parser";

const a = parse("age>$1 and sex=$2");
console.log(a.children[0]);