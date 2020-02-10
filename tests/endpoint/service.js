const { expect } = require('chai');
const { client } = require('./test-helpers');

const { get } = client(process.env.API_ENDPOINT || 'http://localhost:5000', '../../schema/');

describe('Health Checks', () => {

    describe("Liveness", () => {
        it("returns 200 with no content", async () => {
            await get('/live').expect(200, '');
        });
    });

    describe("Readiness", () => {
        it("returns 200 with no content", async () => {
            await get('/ready').expect(200, '');
        })
        .slow(500).timeout(1000);
    });

    describe("Report", () => {

        var hr;

        before(() => {
            hr = get("/health");
        });

        it("should return 200 with json content", async () => {
            await hr
                .expect('Content-Type', /json/)
                .expect('x-schema','http://schemas.sentex.io/service/health.json')
                .expect(200);
        })
        .slow(500).timeout(1000);

        it("validates", async () => {
            await hr.validate();
        })
        .slow(500).timeout(1000);

        it("all health checks status equal OK", async () => {
            const body = (await hr).body;
            ["liveness", "readiness"].forEach(n => {
                expect(body[n].every(c => c.status === "OK"),"status should be OK").to.be.true;
            });
        })
        .slow(500).timeout(1000);
    });
});

describe('Metadata', () => {
    var hr;

    before(() => {
        hr = get("/metadata");
    });

    it("should return 200 with json content", async () => {
        await hr
            .expect('Content-Type', /json/)
            .expect('x-schema','http://schemas.sentex.io/service/metadata.json')
            .expect(200);
    });

    it("validates", async () => {
        await hr.validate();
    })
});

describe('Schemas', () => {
    var hr;

    before(() => {
        hr = get("/schemas");
    });

    it("should return 200 with json content", async () => {
        await hr
            .expect('Content-Type', /json/)
            .expect('x-schema','http://schemas.sentex.io/service/schemalist.json')
            .expect(200);
    });

    it("validates", async () => {
        await hr.validate();
    })
});