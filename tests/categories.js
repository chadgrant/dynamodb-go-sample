const request = require('supertest');
const {expect} = require('chai');
const {get} = request(process.env.API_ENDPOINT || 'http://localhost:5000');

describe('Categories', () => {

    var req;

    before(()=>{
        req = get("/categories");
        console.log("endpoint=" + process.env.API_ENDPOINT)
    });

    it("returns 200 with content type of json",async () => {
        await req
        .expect(200)
        .expect('Content-Type',/json/);
    });

    it("returns categories", async () => {
        const b = (await req).body;

        expect(b,"results").is.a('array');
        expect(b.length,"results").is.greaterThan(0);
        b.forEach(c => {
            expect(c,"category").is.not.empty;
        });
    });
});