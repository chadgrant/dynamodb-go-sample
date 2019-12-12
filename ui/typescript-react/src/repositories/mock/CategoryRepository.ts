import CategoryRepository from '../CategoryRepository';

export default class MockCategoryRepository implements CategoryRepository {
    getAll(): Promise<string[]> {
        return Promise.resolve(["Hats", "Shirts", "Pants", "Shoes", "Ties", "Belts", "Socks", "Accessory"]);
    }
}