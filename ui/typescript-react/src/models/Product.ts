export class ProductInfo {
  id: string
  price: number

  constructor(id:string, price:number) {
    this.id = id;
    this.price = price;
  }
}

export default class Product extends ProductInfo {
    category: string
    name: string
    description: string
  
    constructor(id:string, category:string, name:string, price:number, description:string) {
      super(id,price);
      this.category = category;
      this.name = name;
      this.description = description;
    }
  }

  export class ProductResponse {
    results : Product[];
    next :string;
  
    constructor(results:Product[], next:string) {
      this.results = results;
      this.next =next;
    }
  }
  