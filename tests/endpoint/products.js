const { expect } = require('chai');
const { client } = require('./test-helpers');

const { get, post, put, del } = client(process.env.API_ENDPOINT || 'http://localhost:5000', '../../schema/');

describe('Products', () => {

    let products;

    before(async () => {
        products = (await getProducts()).body.results;
    });

    describe("Paged", () => {
        
        it("validates", async () => {
            await getProducts().validate();
        });

        it("should return products", async () => {
            let counter = 0;
            let ps = (await getProducts()).body;
            
            while (ps.next !== undefined) {
                ps = (await getProducts(ps.next)).body;
                counter++;
            }
            expect(counter, "page counts").to.be.greaterThan(3);
        }).slow(2000)
        .timeout(5000);

        ["hats", "belts"].forEach(cat => {
            it(`should return products for category: ${cat}`, async () => {
                const body = (await getProducts(`/products/${cat}`)).body;
                expect(body.results, "results").is.a('array');
                expect(body.results.length, "results").is.greaterThan(0);
            });
        });
    });

    describe("Get", () => {

        it("should return product by id", async () => {
            const body = (await getProduct(products[0].id)).body;
            expect(body.id, "id").to.equal(products[0].id);
        });

        it("validates", async () => {
            await getProduct(products[0].id);
        });

        it("should return 404 for fake product id", async () => {
            await get("/product/thisisnotanid").expect(404);
        });
    });

    describe("Add", () => {

        it("should add product", async () => {
            const product = createProduct();

            const res = (
                await post('/products/')
                    .send(product)
                    .expect(201, '')
            );

            expect(res.headers.location, "location header").to.not.be.empty;
            expect(res.headers.location, "location header").to.match(/\/product\/.+/);

            const added = (await get(res.headers.location).expect(200).validate()).body;
            expect(added.id, "id").to.not.be.empty;
            delete(added['id']);
            expect(added).to.deep.equal(product);
        });

        ["category", "price", "name"].forEach(prop => {
            const product = createProduct();

            it(`should require ${prop}`, async () => {
                delete product[prop];
                await post("/products/")
                    .send(product)
                    .expect(400);
            });
        });

        ["description"].forEach(prop => {
            const product = createProduct();

            it(`should not require ${prop}`, async () => {
                delete product[prop];
                await post("/products/")
                    .send(product)
                    .expect(201);
            });
        });
    });

    describe("Update", () => {

        it("should update product", async () => {

            let p = products[Math.floor(Math.random() * products.length)];

            await put("/product/" + p.id)
                .send(p)
                .expect(204);
        });

        ["name", "description", "category"].forEach(prop => {
            it(`should update ${prop}`, async () => {
                let p = products[Math.floor(Math.random() * products.length)];
                await update(p, prop, `updated ${prop}`);
            });
        });

        it("should update price", async () => {
            let p = products[Math.floor(Math.random() * products.length)];

            await update(p, "price", Number((p.price).toFixed(2)));
        });

        ["category", "price", "name"].forEach((prop) => {
            it(`should require ${prop}`, async () => {
                let p = products[Math.floor(Math.random() * products.length)];

                delete p[prop];
                await put(`/product/${p.id}`)
                    .send(p)
                    .expect(400);
            });
        });
    });

    describe("Delete", () => {
        it("should delete product", async () => {
            const product = createProduct();
            const loc = (
                await post('/products/')
                    .send(product)
                    .expect(201, '')
            ).headers.location;

            const added = (await get(loc)).body;
            expect(added.name, "added name").to.equal(product.name);
            await del(loc).expect(204);
            await get(loc).expect(404);
        });
    });
});

async function update(p, prop, val) {
    p[prop] = val;

    await put(`/product/${p.id}`)
        .send(p)
        .expect(204);

    const updated = (await getProduct(p.id)).body;
    expect(updated[prop], `updated ${prop}`).to.equal(val);
    return updated;
}

function createProduct(overrides) {
    const obj = {
        category: "hats",
        name: "Test product",
        description: "this is a description",
        price: 1.22
    };
    for (let prop in overrides) {
        obj[prop] = overrides[prop];
    }
    return obj;
}

getProducts = (path) => {
    const p = path || '/products/hats';
    return get(p)
           .expect(200)
           .expect('Content-Type', /json/)
           .expect('X-Schema','http://products.sentex.io/product.paged.json')
           .validate();
}

getProduct = (id) => {
    return get(`/product/${id}`)
            .expect(200)
            .expect('Content-Type', /json/)
            .expect('X-Schema','http://products.sentex.io/product.json')
            .validate();
}
