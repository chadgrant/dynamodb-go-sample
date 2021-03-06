import Product from "../models/Product"

export default interface ProductRepository {
    getPaged(category:string, last?: string, lastPrice? : number, size? : number) : Promise<Product[]>
    get(id:string) : Promise<Product | null>
    add(product:Product) : Promise<Product>
    edit(product:Product) : Promise<any>
    delete(id:string) : Promise<any>
}
