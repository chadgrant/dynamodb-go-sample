import React, { FC } from 'react';
import Product from '../models/Product';

interface ProductEditProps {
    product: Product,
    remove: (id: string) => void
    edit: (product:Product,name:string,price:number,description:string)=>void
}

const ProductEdit:FC<ProductEditProps> = (props) => {
    const product = props.product;
    const delbttn = (product.id !== undefined) ? <i/> : <i className="fas fa-trash"/>;
    const name = React.createRef<HTMLInputElement>();
    const price = React.createRef<HTMLInputElement>();
    const description = React.createRef<HTMLInputElement>();

    const edit = (_:any) => {
      props.edit(product, name.current!.value, Number(price.current!.value), description.current!.value);
    }

    return (
      <tr className="bg-warning" key={product.id} data-id={product.id}>
        <td>
          <input ref={name} defaultValue ={product.name} size={50}/>
        </td>
        <td>
          <input ref={price} defaultValue={product.price} size={7}/>
        </td>
        <td>
          <input ref={description} defaultValue={product.description} size={80}/>
        </td>
        <td>
          <i className="far fa-save" onClick={edit}></i>
        </td>
        <td>
          {delbttn}
        </td>
      </tr>
    );
}

export default ProductEdit;