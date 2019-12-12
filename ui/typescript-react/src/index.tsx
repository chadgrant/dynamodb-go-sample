import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import Product, { ProductInfo } from './models/Product';
import RESTCategoryRepository from './repositories/rest/CategoryRepository';
import RESTProductRepository from './repositories/rest/ProductRepository';
//import MockCategoryRepository from './repositories/mock/CategoryRepository';
//import MockProductRepository from './repositories/mock/ProductRepository';
//const categoryRepo = new MockCategoryRepository();
//const productRepo = new MockProductRepository(categoryRepo,100);

const API_ENDPOINT = window.API_ENDPOINT ? window.API_ENDPOINT : 'http://localhost:5000';
const categoryRepo = new RESTCategoryRepository(API_ENDPOINT);
const productRepo = new RESTProductRepository(API_ENDPOINT);

export interface CustomWindow extends Window { API_ENDPOINT: string;  }
declare let window: CustomWindow;

const defaultState = {
    categories: new Array<string>(),
    category: "Hats",
    products: new Array<Product>(),
    page: 0,
    pages : new Array<ProductInfo>(),
    eof: false
};
 
ReactDOM.render(<App state={defaultState} categoryRepo={categoryRepo} productRepo={productRepo}/>, document.getElementById('root'));
