import CategoryRepository from './CategoryRepository';

export default class RESTCategoryRepository implements CategoryRepository {
    baseurl:string;

    constructor(baseurl:string) {
        this.baseurl = baseurl;
    }

    async getAll(): Promise<string[]> {
        const resp = await fetch(`${this.baseurl}/category`);
        return await resp.json() as Promise<string[]>;
    }
}