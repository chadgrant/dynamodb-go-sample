const { expect } = require('chai');
const { client } = require('./test-helpers');

const { get } = client(process.env.API_ENDPOINT || 'http://localhost:5000');

describe('CORS', () => {

    const corsHeaders = {
        'access-control-allow-origin': '*',
        'access-control-allow-credentials': 'true',
        'access-control-expose-headers': 'Location'
    };

    const paths = ['/live', '/health', '/ready', '/categories', '/products/hats'];

    paths.forEach(p => {
        it(`${p} returns cors headers`, async () => {
            const resp = await get(p);
            for (let prop in corsHeaders) {
                expect(resp.headers[prop], `${p} cors header: ${prop}`).to.equal(corsHeaders[prop]);
            }
        }).slow(500).timeout(1000);
    });
});