const request = require('supertest');
const expect = require('chai').expect;

const {get,post} = request(process.env.API_ENDPOINT || 'http://localhost:5000');

describe('Health Checks', () => {

    it("/live returns 200 with no content", async () => {
        await get('/live').expect(200, '');
    });

    it("/ready returns 200 with no content", async () => {
        await get('/ready').expect(200, '');
    })
    .slow(500).timeout(1000);

    describe("Report", () => {

        var hr;

        before(() => {
            hr = get("/health");
        });

        it("should return 200 with json content", async () => {
            await hr
            .expect('Content-Type',/json/)
            .expect(200);
         })
         .slow(500).timeout(1000);
        

        it("returns valid report", async () => {
            const b = (await hr).body;
            expect(b.report_as_of_utc,"report as of").to.not.be.empty;
            expect(b.up_since,"up since").to.not.be.empty;
            expect(b.duration_ms,"duration ms").to.be.greaterThan(0);
            expect(b.readiness, "readiness").to.be.a('array');
            expect(b.liveness,"liveness").to.be.a('array');
        })
        .slow(500).timeout(1000);

        it("all health checks status equal OK", async () => {
            const b = (await hr).body;
            ["liveness","readiness"].forEach(n=>{
                b[n].forEach(c => {healthCheck(n,c)});
            });
        })
        .slow(500).timeout(1000);
    });
});

function healthCheck(context,hc) {
    expect(hc, `[${context}] health check not null`).to.not.be.undefined.and.not.be.null;
    expect(hc.name, `[${context}] health check name`).to.not.be.undefined.and.not.be.null;
    const name = hc.name;
    expect(hc.status, `[${context}] health check (status): ${name}`).to.equal("OK");
    expect(hc.duration_ms, `[${context}] health check (duration_ms): ${name}`).to.be.greaterThan(0);
    expect(hc.tested_at_utc, `[${context}] health check (tested_at_utc): ${name}`).to.not.be.empty;
}

describe('Metadata', () => {
    
    it("should return 200 with json content", async () => {
        await get('/metadata')
        .expect('Content-Type',/json/)
        .expect(200);
    });

    it("should have required properties", async () => {
        const b = (await get('/metadata')).body;

        ["vendor","group","service","friendly","description",
        "build_repo","build_number","built_by","git_hash","git_branch",
        "compiler_version"].forEach((prop) => {
            expect(b[prop],prop).to.not.be.empty;
        });
    });
});