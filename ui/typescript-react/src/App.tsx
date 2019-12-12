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
  pages: ProductInfo[]
}

const App:FC<{state:AppState, categoryRepo:CategoryRepository, productRepo:ProductRepository}> = (props) => {
  const [appState,setAppState] = useState<AppState>(props.state);

  useEffect(()=>{
    props.categoryRepo.getAll().then(cats=>{
      props.productRepo.getPaged(props.state.category, undefined, undefined, PAGE_SIZE).then(products=>{
        setAppState((p:AppState) => {
          return {
            categories: cats, 
            category: p.category, 
            products: products, 
            page: p.page, 
            pages: p.pages};
          });
      });
    });
  },[])

  const changeCategory = async (category:string) => {
    const products = await props.productRepo.getPaged(category, undefined, undefined, PAGE_SIZE);
    setAppState((p:AppState) => {
      return {
        categories: p.categories, 
        category: category, 
        products: products, 
        page: p.page, 
        pages: p.pages};
      });
  } 

  const next = async () => {
    const {id,price} = appState.products[appState.products.length - 1];
    const products = await props.productRepo.getPaged(appState.category, id, price, PAGE_SIZE);
    setAppState((p:AppState) => {
      const pages = [...p.pages,{id,price}];
      return {
        categories: p.categories, 
        category: p.category, 
        products: products,
        page: p.page+1, 
        pages: pages
      }
    });
  }

  const prev = async () => {
    let id= "", price = 0;
    if (appState.page > 2) {
        const t = appState.pages[appState.page-2];
        id = t.id;
        price = t.price;
    }
    const products = await props.productRepo.getPaged(appState.category, id, price, PAGE_SIZE);
    setAppState((p:AppState) => {
      return {
        categories: p.categories, 
        category: p.category, 
        products: products,
        page: p.page-1, 
        pages: p.pages
      }
    });
  }

  const add = () => {
    setAppState((p:AppState) => {
      var newp = {
        categories: p.categories, 
        category: p.category, 
        products: p.products,
        page: p.page, 
        pages: p.pages
      };
      newp.products.unshift({id : "", category: p.category, name: "", price: 1, description: "" })
      return newp;
    });
  }

  const remove = async (id:string) => {
    let r = window.confirm("Do you want to delete this item");
    if( r === true) {
      await props.productRepo.delete(id);

      setAppState((p:AppState) => {
        return {
          categories: p.categories, 
          category: p.category, 
          products: p.products.filter(p => p.id !== id),
          page: p.page-1, 
          pages: p.pages
        }
     });
    }   
  }

  const edit = async (product:Product,name:string,price:number,description:string) => {
    product.name = name;
    product.price = price;
    product.description = description;

    await props.productRepo.edit(product);
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
                  disabled={(appState.products.length < PAGE_SIZE)}>Next</button>
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