import React, {FC, useState, useEffect} from 'react';
import CategoryRepository from './repositories/CategoryRepository';
import ProductRepository from './repositories/ProductRepository';
import Product, {ProductInfo} from './models/Product';
import Categories from './components/Categories';
import Products from './components/Products';
import './App.css';

const PAGE_SIZE = 25;

interface AppState {
  categories: string[]
  category : string
  products: Product[]
  page: number
  pages: ProductInfo[],
  eof: boolean
}

const App:FC<{state:AppState, categoryRepo:CategoryRepository, productRepo:ProductRepository}> = (props) => {
  const [appState,setAppState] = useState<AppState>(props.state);

  useEffect(()=>{
    props.categoryRepo.getAll().then(cats=>{
      props.productRepo.getPaged(props.state.category, undefined, undefined, PAGE_SIZE).then(products=>{
        setAppState((p:AppState) => {
          return mutate(p, n=>{
            n.categories = cats;
            n.products = products;
            return n;
          });
        });
      });
    });
  },[props.state.category, props.categoryRepo, props.productRepo])

  const changeCategory = async (category:string) => {
    const products = await props.productRepo.getPaged(category, undefined, undefined, PAGE_SIZE);
    setAppState((p:AppState) => {
      return mutate(p, n=>{
        n.category = category;
        n.products = products;
        return n;
      });
    });
  } 

  const next = async () => {
    const {id,price} = appState.products[appState.products.length - 1];
    const products = await props.productRepo.getPaged(appState.category, id, price, PAGE_SIZE);
    setAppState((p:AppState) => {
      const pages = [...p.pages,{id,price}];
      return mutate(p, n=>{
        n.page++;
        n.products = (products.length === 0) ? p.products : products;
        n.eof = products.length === 0;
        n.pages = pages;
        return n;
      });
    });
  }

  const prev = async () => {
    let id="", price=0;
    if (appState.page > 2) {
        const t = appState.pages[appState.page-2];
        id = t.id;
        price = t.price;
    }
    const products = await props.productRepo.getPaged(appState.category, id, price, PAGE_SIZE);
    setAppState((p:AppState) => {
      return mutate(p, n=>{
        n.page--;
        n.products = products;
        return n;
      });
    });
  }

  const add = () => {
    setAppState((p:AppState) => {
      return mutate(p, n=>{
        n.products.unshift({id : "", category: p.category, name: "", price: 1, description: "" })
        return n;
      });
    });
  }

  const remove = async (id:string) => {
    if(window.confirm("Do you want to delete this item") === true) {
      await props.productRepo.delete(id);

      setAppState((p:AppState) => {
        return mutate(p, n=>{
          n.products = p.products.filter(p => p.id !== id);
          return n;
        });        
     });
    }   
  }

  const edit = async (product:Product,name:string,price:number,description:string) => {
    product.name = name;
    product.price = price;
    product.description = description;

    await props.productRepo.edit(product);
  }

  const mutate = (prev:AppState, callback:(x:AppState) => AppState) => {
    return callback({
      categories: prev.categories, 
      category: prev.category, 
      products: prev.products,
      page: prev.page, 
      pages: prev.pages,
      eof: false
    });
  }

  return (
    <div className="container-fluid">
      <div className="row mt-3">
        <div className="col-lg-12">
          <div className="card">
            <div className="card-header">
              <Categories 
                  categories={appState.categories} 
                  category={appState.category}
                  changeCategory={changeCategory} />
              <button
                  className="btn btn-dark float-right btn-next"
                  onClick={next}
                  disabled={(appState.products.length < PAGE_SIZE || appState.eof)}>Next</button>
              <button
                  className="btn btn-dark float-right btn-prev"
                  onClick={prev}
                  disabled={(appState.page < 1)}>Previous</button>
              <button
                  className="btn btn-primary float-right btn-add"
                  onClick={add}>Add New</button>                
              </div>
            <div className="card-body">
              <table className="table table-hover">
                <thead className="thead-dark"><tr><th>Name</th><th>Price</th><th>Description</th><th>Edit/Save</th><th>Delete</th></tr></thead>
                  <Products 
                    products={appState.products} 
                    remove={remove}
                    edit={edit}                
                  />
              </table>
          </div>
        </div>
      </div>
    </div>
  </div>      
  )
}
export default App;