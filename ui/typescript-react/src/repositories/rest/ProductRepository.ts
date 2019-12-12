import ProductRepository from '../ProductRepository'
import Product,{ProductResponse} from "../../models/Product"

export default class RESTProductRepository implements ProductRepository {
    endpoint:string;

    constructor(endpoint:string) {
        this.endpoint = endpoint;
    }

    async getPaged(category:string, last?:string, lastPrice?: number, size:number = 25): Promise<Product[]> {
        const r = await fetch(`${this.endpoint}/product/${category.toLowerCase()}?last=${((last !== undefined) ? last : "")}&lastprice=${((lastPrice !== undefined ? lastPrice : ""))}`);
        if (!r.ok) {
            return new Array<Product>();
        }
        const resp = await (r.json() as Promise<ProductResponse>);
        return await resp.results;
    }    
    
    async get(id: string): Promise<Product | null> {
        const r = await fetch(`${this.endpoint}/product/${id}`);
        if (!r.ok) {
            return null;
        }
        return await (r.json() as Promise<Product>);
    }

    async add(product: Product): Promise<Product> {
        const r = await fetch(`${this.endpoint}/product/`,{
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product)
        });
        if (!r.ok) {
            throw new Error("couldn't add product");
        }
        const r2 = await fetch(`${this.endpoint}${r.headers.get('Location')}`);
        return  await (r2.json() as Promise<Product>);
    }

    async edit(product: Product): Promise<any> {
        const r = await fetch(`${this.endpoint}/product/${product.id}`,{
            method: 'PUT',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product)
        });
        if (!r.ok) {
            throw new Error("couldn't add product");
        }
    }

    async delete(id: string): Promise<any> {
        return await fetch(`${this.endpoint}/product/${id}`,{method: 'DELETE'});
    }
}