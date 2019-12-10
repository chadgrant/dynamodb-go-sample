import React, { Component } from 'react';
import Product from './Product'

export default class Products extends Component {
    render() {
        return (
            <tbody>
            {this.props.products.map((product,index) => (
                <Product
                    key={product.id}
                    product={product}
                    index={index}
                    edit={this.props.edit}
                    delete={this.props.delete}
                />
            ))}
            </tbody>
        )
    }
}