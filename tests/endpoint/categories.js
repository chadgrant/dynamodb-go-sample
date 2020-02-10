const { expect } = require('chai');
const { client } = require('./test-helpers');

const { get } = client(process.env.API_ENDPOINT || 'http://localhost:5000', '../../schema/');

describe('Categories', () => {

    it("returns 200 with content type of json", async () => {
       await get("/categories")
            .expect(200)
            .expect('Content-Type', /json/)
            .expect('X-Schema','http://schemas.sentex.io/store/categories.json');
    });

    it("validates", async () => {
        await get("/categories").validate();
    });

    it("returns categories", async () => {
        const body = (await get("/categories")).body;
        body.forEach(c => {
            expect(c, "category").is.not.empty;
        });
    });
});