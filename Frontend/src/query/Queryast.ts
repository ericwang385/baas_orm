export interface QueryAST{
    //Body: Body,
    SessionId: string,
    Where?: Expr,
    Args?: string[],
    Table: string,
    Preload?: string[],
    OrderBy?: string[],
    Offset?: number,
    Limit?: number,
    Desc?: boolean
}

// interface Body {
//     Query: QueryAST,
//     Select: Select,
//     Value?: string[],
// }

// interface Select {
//     Distinct: Boolean,
//     GroupBy: string,
//     Projection: string[],
//     From: Table[],
//     Where?: {
//         expr: Expr,
//         args: string[]
//     },
// }

// type Table = {
//     name: string,
//     columns: string[]
// }

type Expr = {
    T: string,
    V: string,
    C: Expr[]
}