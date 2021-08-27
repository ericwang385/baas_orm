export type Column = {
    dataType: string
    name: string
    nullable: boolean
}

export type Relation = {}

export type QueryResponse = {
    tables: [{
        name: string,
        columns: string[],
        rows: [[]]
    }],
    duration: number
}

export type SaveRequest = {
    dirtyData: {[key: string]:any},
    pkcolumn: string[]
}