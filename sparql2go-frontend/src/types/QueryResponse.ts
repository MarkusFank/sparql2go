export type QueryResponse = {
    count: number;
    result: {[key: string]: string}[] // TODO value always a string??
    vars: string[]
}