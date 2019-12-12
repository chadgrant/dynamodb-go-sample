import React, { FC } from 'react';
import ProductModel from '../models/Product';
import Product from './Product'

interface ProductProps {
    products: ProductModel[],
    remove: (id: string) => void
    edit: (product:ProductModel,name:string,price:number,description:string)=>void
}

const Products:FC<ProductProps> = (props) => {
    return (
        <tbody>
        {props.products.map((product) => (
            <Product
                key={product.id}
                product={product}
                edit={props.edit}
                remove={props.remove}
            />
        ))}
        </tbody>
    )
}

export default Products;