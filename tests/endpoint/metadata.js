const { client } = require('./test-helpers');

const { get } = client(process.env.API_ENDPOINT || 'http://localhost:5000', '../../schema/');

describe('Metadata', () => {
    var hr;

    before(() => {
        hr = get("/metadata");
    });

    it("should return 200 with json content", async () => {
        await hr
            .expect('Content-Type', /json/)
            .expect(200);
    });

    it("validates", async () => {
        await hr.validate('types/metadata.json');
    })
});