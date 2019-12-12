import React, { FC, useState } from 'react';

type CategoriesProps = {
    category: string,
    categories: string[],
    changeCategory: (id: string) => any
}

type CategoriesState = {
    category: string
}

const Categories:FC<CategoriesProps> = (props) => {
    const [state,setState] = useState<CategoriesState>({category:props.category});

    const onChange = (e: React.ChangeEvent<HTMLSelectElement>) =>{
        setState({category: e.target.value});
        props.changeCategory(e.target.value);
    }

    return (
        <>
            Products : <select className="sel-cat" onChange={onChange} value={state.category}>
            {props.categories.map(c => {
                return (<option key={c}>{c}</option>);
            })}
            </select>
        </>
    )
}

export default Categories;