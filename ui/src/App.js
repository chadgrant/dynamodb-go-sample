import React, { Component } from 'react';
import './App.css';
import Products from './components/products';

const API_ENDPOINT = window.API_ENDPOINT ? window.API_ENDPOINT : 'http://localhost:5000';

export default class App extends Component {
  constructor(props){
    super(props);
    this.state = { category : "Hats", products: [] }
    this.add = this.add.bind(this);
    this.edit = this.edit.bind(this);
    this.delete = this.delete.bind(this);
  }

  componentDidMount() {
    fetch(`${API_ENDPOINT}/product/${this.state.category.toLowerCase()}`)
    .then(res => res.json())
    .then((data) => {
      this.setState((prev) => ({
          category: prev.category,
          products: data.results
      }));
    })
    .catch(console.log)
  }

  add(name, price, description) {
    // fetch(`${API_ENDPOINT}/product/`, {
    //   method: 'PUT',
    //   headers: {
    //     'Accept': 'application/json',
    //     'Content-Type': 'application/json',
    //   },
    //   body: JSON.stringify({
    //     category: this.state.category,
    //     name: name,
    //     price: parseFloat(price),
    //     description: description
    //   })
    // })
    // .then(console.log)
    // .catch(console.log);

  //   this.setState((prev) => ({
  //     category: prev.category,
  //     products: prev.products.unshift()
  //  }));
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
    product.price = parseFloat(price);
    product.description = description;
    
    fetch(`${API_ENDPOINT}/product/${product.id}`,{
      method: 'PUT',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(product)
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
    return (
      <div className="container-fluid">
        <div className="row mt-3">
          <div className="col-lg-12">
            <div className="card">
              <div className="card-header">{this.state.category} Registry</div>
              <div className="card-body">
                <table className="table table-hover">
                  <thead className="thead-dark"><tr><th>Name</th><th>Price</th><th>Description</th><th>Edit/Save</th><th>Delete</th></tr></thead>
                    <Products 
                      products={this.state.products} 
                      delete={this.delete}
                      edit={this.edit}                
                    />
                </table>
                <button
                  className="btn btn-dark pull-left"
                  onClick={this.add}>
                  Add New
                </button>
            </div>
          </div>
        </div>
      </div>
    </div>      
    )
  }
}
