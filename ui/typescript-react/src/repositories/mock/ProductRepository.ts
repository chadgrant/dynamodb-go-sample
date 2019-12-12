import CategoryRepository from '../CategoryRepository';
import ProductRepository from '../ProductRepository';
import Product from '../../models/Product';

export default class MockProductRepository implements ProductRepository {
    products: Map<string,Map<string,Product>>;
    constructor(catRepo:CategoryRepository, max:number) {
        this.products = new Map<string,Map<string,Product>>();
        catRepo.getAll().then(cats => {
            cats.forEach(c => {
                let map = new Map<string,Product>();
                for (let i=0; i < max; i++) {
                    let p = new Product(`${c.toLowerCase()}-${i}`,c,`${c} Product ${i}`,1,`Description of cool ${c} Product ${i}`);
                    map.set(p.id,p);
                }
                this.products.set(c, map);
            });
        });
    }

    getPaged(category: string, last?: string, lastPrice?: number, size:number = 25): Promise<Product[]> {
        const arr:Product[] = new Array(size);
        const prods = this.products.get(category)
        if (prods) {
            let counter = 0;
            prods.forEach((v,_) => {
                if (counter >= arr.length) { return false }
                arr[counter++] = v;
            });
        }
        return Promise.resolve(arr);
    }
    
    add(product: Product): Promise<Product> {
        throw new Error("Method not implemented.");
    }
    edit(product: Product): Promise<any> {
        throw new Error("Method not implemented.");
    }
    delete(id: string): Promise<any> {
        throw new Error("Method not implemented.");
    }
}