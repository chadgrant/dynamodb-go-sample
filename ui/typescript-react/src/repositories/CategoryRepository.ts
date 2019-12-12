export default interface CategoryRepository {
    getAll() : Promise<string[]>
}