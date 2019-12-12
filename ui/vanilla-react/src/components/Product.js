import React, { Component } from 'react';

export default class Product extends Component {

    state = {
      isEdit:this.props.product.id === undefined
    };
    
    delete = () => {
        this.props.delete(this.props.product.id);
    }

    edit = () => {
      this.setState((prev) => ({
        isEdit : !prev.isEdit
      }));
    }

    editSubmit = () => {
      this.setState((prev) => ({
        isEdit : !prev.isEdit
      }));
       
      this.props.edit(
        this.props.product,
        this.nameInput.value,
        this.priceInput.value,
        this.descriptionInput.value
      );
    }
    
    render() {
      return (this.state.isEdit === true ? 
        this.renderEdit(this.props.product) : 
        this.renderList(this.props.product)
      );
    }

    renderList = (product) => {
      return (
        <tr key={product.id} data-id={product.id}>
          <td>{product.name}</td>
          <td>${product.price}</td>
          <td>{product.description}</td>
          <td><i className="far fa-edit" onClick={this.edit}></i></td>
          <td><i className="fas fa-trash" onClick={this.delete}></i></td>
        </tr>
      );
    }

    renderEdit = (product) => {
      let delbttn = <i/>
      if (product.id !== undefined) {
        delbttn = <i className="fas fa-trash"/>
      }
      return (
          <tr className="bg-warning" key={product.id} data-id={product.id}>
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
            <td>
              {delbttn}
            </td>
          </tr>
      );
    }
  }