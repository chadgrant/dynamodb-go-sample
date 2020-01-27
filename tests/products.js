const request = require('supertest');
const expect = require('chai').expect;

const {get,post,put,del} = request(process.env.API_ENDPOINT || 'http://localhost:5000');

describe('Products', () => {

    describe("get paged", () => {
        it("should return products", async () => {
            let counter = 0;
            let ps = await page("/products/hats");
            while(ps.next !== undefined) {
                ps = await page(ps.next);
                counter++;
            }
            expect(counter,"page counts").to.be.greaterThan(3);
        })
        .slow(2000)
        .timeout(5000);
    });

    describe("get", () => {

        ["hats","belts"].forEach(cat => {
            it(`should return ${cat}`, async () => {
                const b = (
                    await get(`/products/${cat}`)
                    .expect(200)
                    .expect('Content-Type',/json/)
                ).body;

                expect(b.results,"results").is.a('array');
                expect(b.results.length,"results").is.greaterThan(0);
                b.results.forEach(test);
            });
        });

        it("should return product by id", async () => {
            const b = (
                await get("/products/hats")
                .expect(200)
                .expect('Content-Type',/json/)
            ).body;

            const p = await getProduct("/product/" + b.results[0].id);
            expect(p.id,"id").to.equal(b.results[0].id);
        });

        it("should return 404 for fake product id", async () => {
            await get("/product/thisisnotanid").expect(404);
        });
    });

    describe("add", () => {

        it("should add product", async () => {
            const p = {
                category: "hats",
                name: "Test product",
                description: "this is a description",
                price: 1.22
            };

            const res = (
                await post('/products/')
                .send(p)
                .expect(201,'')
            );

            expect(res.headers.location,"location header").to.not.be.empty;
            expect(res.headers.location,"location header").to.match(/product/);

            const added = await getProduct(res.headers.location);
            expect(added.id,"id").to.not.be.empty;
            expect(added.category,"category").to.equal(p.category);
            expect(added.name,"name").to.equal(p.name);
            expect(added.description,"description").to.equal(p.description);
            expect(added.price).to.equal(p.price);
        });

        ["category","price","name"].forEach((prop) =>{
            let p = {
                name: "Test Product",
                category: "hats",
                description: "this is a description",
                price: 1.1
            };
            it(`should require ${prop}`, async () => {
                delete p[prop];
                await post("/products/")
                .send(p)
                .expect(400);
            });
        });

        ["description"].forEach(prop => {
            let p = {
                name: "Test Product",
                category: "hats",
                description: "this is a description",
                price: 1.25
            };
            it(`should not require ${prop}`, async () => {
                delete p[prop];
                await post("/products/")
                .send(p)
                .expect(201);
            });
        });
    });

    describe("update", () => {

        let products;

        before(async () => {
           products = (await get("/products/hats").expect(200)).body.results;
        });

        it("should update product", async () => {

            let p = products[Math.floor(Math.random()*products.length)];

            await put("/product/" + p.id)
            .send(p)
            .expect(204);
        });

        ["name","description","category"].forEach(prop => {
            it(`should update ${prop}`, async () => {
                let p = products[Math.floor(Math.random()*products.length)];
                await update(p, prop, `updated ${prop}`);
            });
        });
        

        it("should update price", async () => {
            let p = products[Math.floor(Math.random()*products.length)];

            await update(p, "price", Number((p.price).toFixed(2)));
        });

        ["category","price","name"].forEach((prop) =>{
            it(`should require ${prop}`, async () => {
                let p = products[Math.floor(Math.random()*products.length)];

                delete p[prop];
                await put(`/product/${p.id}`)
                .send(p)
                .expect(400);
            });
        });
    });

    describe("delete", () => {
        it("should delete product", async () => {
            const p = {
                category: "hats",
                name: "Test product to delete",
                price: 1.22
            };

            const loc = (
                await post('/products/')
                .send(p)
                .expect(201,'')
            ).headers.location;

            const added = await getProduct(loc);
            expect(added.name,"added name").to.equal(p.name);
            await del(loc).expect(204);
            await get(loc).expect(404);
        });
    });
});

async function getProduct(url) {
    const b = (
        await get(url)
        .expect(200)
        .expect('Content-Type',/json/)
    ).body;
    return test(b);
}

function test(p) {
    expect(p.id,"product id").to.not.be.empty;
    expect(p.category,"product category").to.not.be.empty;
    expect(p.name,"product name").to.not.be.empty;
    return p;
}

async function update(p, prop, val) {
    p[prop] = val;
    
    await put(`/product/${p.id}`)
        .send(p)
        .expect(204);
        
    const updated = (
        await get(`/product/${p.id}`)
        .expect(200)
    ).body;
            
    expect(updated[prop],`updated ${prop}`).to.equal(val);
    return test(updated);
}

async function page(url) {
    const resp = (await get(url).expect(200)).body;

    if (resp.next !== undefined) {
        expect(resp.results).is.a('array');
        expect(resp.results.length).is.greaterThan(0);
        resp.results.forEach(test);
    }

    return resp;
}