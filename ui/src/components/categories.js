import React, { Component } from 'react';

export default class Categories extends Component {
    constructor(props){
        super(props);
        this.change = this.change.bind(this);
     }

    change(e) {
        this.props.changeCategory(e.target.value)
    }

    render() {
        return (
            <select className="sel-cat" onChange={this.change}>
            {this.props.categories.map((cat) => (
                <option key={cat}>{cat}</option>    
            ))}
            </select>
        )
    }
}