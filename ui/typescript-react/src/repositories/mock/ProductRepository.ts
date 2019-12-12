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
        const arr:Product[] = new Array<Product>();
        const prods = this.products.get(category)
        if (prods) {
            let consuming = (last === undefined || last === "" || last === null);
            prods.forEach((v,_) => {
                if (arr.length < size) {
                    if (!consuming) {
                        if(v.id === last) {
                            consuming = true;
                        }
                    } else {                
                        arr.push(v);  
                    }
                }
            });
        }
        return Promise.resolve(arr);
    }

    get(id: string): Promise<Product | null> {
        this.products.forEach((_,k)=>{
            const map = this.products.get(k);
            const p = map!.get(id);
            if (p != undefined){
                return Promise.resolve(p);
            }
        });
        return Promise.resolve(null);
    }
    
    add(product: Product): Promise<Product> {
        const i = Date.now();
        const p = new Product(`${product.category.toLowerCase()}-${i}`,product.category,`${product.category} Product ${i}`,1,`Description of cool ${product.category} Product ${i}`)
        const map = this.products.get(p.category)
        map!.set(p.id,p);
        return Promise.resolve(p);
    }

    edit(product: Product): Promise<any> {
        const map = this.products.get(product.category);
        const p = map!.get(product.id);
        if (p) {
            p.name = product.name;
            p.price = product.price;
            p.description = product.description;
            map!.set(p.id,p);
        }
        return Promise.resolve();
    }

    delete(id: string): Promise<any> {
        this.products.forEach((_,k)=>{
            const map = this.products.get(k);
            map!.delete(id);
        });
        return Promise.resolve();
    }
}