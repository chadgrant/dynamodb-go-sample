import React, { Component } from 'react';
import './App.css';
import Categories from './components/categories';
import Products from './components/products';

const API_ENDPOINT = window.API_ENDPOINT ? window.API_ENDPOINT : 'http://localhost:5000';
const PAGE_SIZE = 25;

export default class App extends Component {
  constructor(props){
    super(props);
    this.state = { categories: [], category : "", products: [], page: 0, pages: [] }
    this.next = this.next.bind(this);
    this.prev = this.prev.bind(this);
    this.add = this.add.bind(this);
    this.edit = this.edit.bind(this);
    this.delete = this.delete.bind(this);
    this.changeCategory = this.changeCategory.bind(this);
  }

  componentDidMount() {
    fetch(`${API_ENDPOINT}/category`)
    .then(res => res.json())
    .then((data) => {
      this.setState((prev) => {
        prev.categories = data;
        prev.category = data[0];
        return prev;
      });
      this.changeCategory(this.state.category);
    })
    .catch(console.log)
  }

  changeCategory(category) {
    this.setState((prev) => {
      prev.pages = [];
      prev.page = 0;
      return prev;
    });
    this.loadPage(category, "","")
  }

  next() {
    const {id,price} = this.state.products[this.state.products.length - 1];
    this.setState((prev) => {
      prev.pages[prev.page] = {id,price};
      prev.page++;
      return prev;
    });
    this.loadPage(this.state.category, id, price);
  }

  prev() {
    let id= "", price = "";
    if (this.state.page > 2) {
        const t = this.state.pages[this.state.page-2];
        id = t.id;
        price = t.price;
    }
    this.setState((prev) => {
      prev.page--;
      return prev;
    });
    this.loadPage(this.state.category, id, price);
  }

  loadPage(category, id, price) {
    fetch(`${API_ENDPOINT}/product/${category.toLowerCase()}?last=${id}&lastprice=${price}`)
    .then(res => res.json())
    .then((data) => {
      this.setState((prev) => {
          prev.category = category;
          prev.products = data.results;
          return prev;
      });
    })
    .catch(console.log)    
  }

  add() {
    this.setState((prev) => {
      prev.products.unshift({
        category: prev.category,
        name: "",
        price: 1,
        description: ""
      });
      return prev;
    });
  }

  delete(id) {
    let r = window.confirm("Do you want to delete this item");
    if( r === true) {
      fetch(`${API_ENDPOINT}/product/${id}`,{method: 'DELETE'})
      .catch(console.log);

      this.setState((prev) => ({
        category: prev.category,
        products: prev.products.filter(p => p.id !== id)
     }));
    }    
  }

  edit(product, name, price, description) {
    product.name = name;
    product.price = Number(price);
    product.description = description;
    
    const method = (product.id === undefined) ? 'POST':'PUT';
    const url = (product.id === undefined) ? `${API_ENDPOINT}/product/`:`${API_ENDPOINT}/product/${product.id}`;

    fetch(url,{
      method: method,
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(product)
    })
    .then(r => {
      fetch(`${API_ENDPOINT}${r.headers.get('Location')}`)
      .then(res => res.json())
      .then((data) => {
        product = data;
      })
      .catch(console.log)
    })
    .catch(console.log);

    this.setState((prev) => ({
      products: prev.products.map(p =>{ 
        if (p.id === product.id) {
          p.name = name;
          p.price = price;
          p.description = description;
        }
        return p;
      })
   }));
  }

  render() {
    const bprev = (this.state.page < 1),
          bnext = (this.state.products.length < PAGE_SIZE);

    return (
      <div className="container-fluid">
        <div className="row mt-3">
          <div className="col-lg-12">
            <div className="card">
              <div className="card-header">
                Products : <Categories 
                categories={this.state.categories} 
                category={this.state.categories}
                changeCategory={this.changeCategory} />
                <button
                    className="btn btn-dark pull-right btn-next"
                    onClick={this.next}
                    disabled={bnext}>Next</button>
                <button
                    className="btn btn-dark pull-right btn-next"
                    onClick={this.prev}
                    disabled={bprev}>Previous</button>
                <button
                  className="btn btn-dark pull-right btn-add"
                  onClick={this.add}>
                  Add New
                </button>                
                </div>
              <div className="card-body">
                <table className="table table-hover">
                  <thead className="thead-dark"><tr><th>Name</th><th>Price</th><th>Description</th><th>Edit/Save</th><th>Delete</th></tr></thead>
                    <Products 
                      products={this.state.products} 
                      delete={this.delete}
                      edit={this.edit}                
                    />
                </table>
            </div>
          </div>
        </div>
      </div>
    </div>      
    )
  }
}
