import React, { FC, useState } from 'react';
import ProductModel from '../models/Product';
import ProductEdit from './ProductEdit'

interface ProductProps {
    product : ProductModel,
    remove : (id:string) => void
    edit : (product:ProductModel,name:string,price:number,description:string) => void
}

interface ProductState {
    isEdit:boolean
}

const Product:FC<ProductProps> = (props) => {
    const product = props.product;
    const [state,setState] = useState<ProductState>({isEdit:props.product.id === undefined || props.product.id === ""});
   
    const remove = (id:string) => {
        props.remove(id);
    };

    const removeProduct = (_:any) => {
      remove(props.product.id);
    };

    const toggleEdit = () => {
      setState((prev) => ({
        isEdit : !prev.isEdit
      }));
    };

    const edit = () => {
      setState((prev) => ({
        isEdit : !prev.isEdit
      }));
    };
          
    if (state.isEdit)
      return <ProductEdit product={product} edit={edit} remove={remove}/>

    return (
      <tr key={product.id} data-id={product.id}>
        <td>{product.name}</td>
        <td>${product.price}</td>
        <td>{product.description}</td>
        <td><i className="far fa-edit" onClick={toggleEdit}></i></td>
        <td><i className="fas fa-trash" onClick={removeProduct}></i></td>
      </tr>
    );
  }

  export default Product;