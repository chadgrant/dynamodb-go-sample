import React, { Component } from 'react';

export default class Product extends Component {
    constructor(props){
      super(props);
      this.state = {isEdit:false}
      this.edit = this.edit.bind(this);
      this.editSubmit = this.editSubmit.bind(this);
      this.delete = this.delete.bind(this);
    }
    delete(){
        this.props.delete(this.props.product.id);
    }
    edit(){
      this.setState((prevState,props) => ({
        isEdit : !prevState.isEdit
      }));
    }
    editSubmit(){
      this.setState((prevState,props) => ({
        isEdit : !prevState.isEdit
      }));
       
      this.props.edit(
        this.props.product,
        this.nameInput.value,
        this.priceInput.value,
        this.descriptionInput.value
      );
    }
    render() {
      const product = this.props.product;
      return (
        this.state.isEdit === true ? (
          <tr className="bg-warning" key={product.id}>
            <td>
              <input ref={nameInput => this.nameInput = nameInput} defaultValue ={product.name} size="50"/>
            </td>
            <td>
              <input ref={priceInput => this.priceInput = priceInput} defaultValue={product.price} size="7"/>
            </td>
            <td>
              <input ref={descriptionInput => this.descriptionInput = descriptionInput} defaultValue={product.description} size="80"/>
            </td>
            <td>
              <i className="far fa-save" onClick={this.editSubmit}></i>
            </td>
            <td><i className="fas fa-trash"></i></td>
          </tr>
        ) : (
          <tr key={product.id}>
            <td>{product.name}</td>
            <td>${product.price}</td>
            <td>{product.description}</td>
            <td><i className="far fa-edit" onClick={this.edit}></i></td>
            <td><i className="fas fa-trash" onClick={this.delete}></i></td>
          </tr>
        )
      );
    }
  }