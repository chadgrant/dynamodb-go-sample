import React, { Component } from 'react';

export default class Categories extends Component {

    state = {
        category : this.props.category
    }

    onChange = (e) =>{
        this.setState({category: e.target.value});
        this.props.changeCategory(e.target.value);
    }

    render() {
        return (
            <React.Fragment>
                Products : <select className="sel-cat" onChange={this.onChange}>
                {this.props.categories.map(this.renderItem)}
                </select>
            </React.Fragment>
        )
    }

    renderItem = (cat) => {
        return cat === this.state.category ? 
            (<option key={cat} selected="true">{cat}</option>) 
            : 
            (<option key={cat}>{cat}</option>)
    }
}