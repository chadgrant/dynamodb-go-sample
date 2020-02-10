const request = require('supertest');
const Ajv = require('Ajv');
const { resolve } = require('path');
const { readdirSync } = require('fs');

const schemas = {};
const loaded = {};
const validator = new Ajv({allErrors: true});

function getFiles(dir,match) {
  const dirents = readdirSync(dir, { withFileTypes: true });
  const files = dirents.map((dirent) => {
    const res = resolve(dir, dirent.name);
    return dirent.isDirectory() ? getFiles(res) : res;
  });
  return Array.prototype.concat(...files).filter(f=>f.match(match));
}

function setHeaders(hdrs, req) {
    for (let p in hdrs) {
        req.set(p, hdrs[p]);
    }
    return req;
}

function addValidator(req) {
    req['validate'] = (schema) => {
        return req.expect(()=>{
            let key = schema || req.response.headers['x-schema'];
            if (!schemas[key].validate(key, req.response.Body)) {
                throw schemas[key].errors;
            }
        });
    };
    return req;
}

async function loadServiceSchemas(baseUrl) {
    if (loaded[baseUrl] != undefined) {
        return;
    }
    loaded[baseUrl] = true
    const all = (await request(baseUrl).get("/schemas")).body;
    for (let i in all) {
        let id = all[i].uri;
        let url = new URL(all[i].url);
        if (schemas[id] === undefined) {
            let s = (await request(url.protocol + "//" + url.host).get(url.pathname + url.search)).body;
            schemas[id] = validator.addSchema(s,id);
        }
    }
}

function client(apiBaseUrl, schemaDirectory) {
    const baseUrl = apiBaseUrl;
    const headers = { 'Origin': 'http://unittests.com' };

    if (schemaDirectory) {
        getFiles(schemaDirectory, "json$").forEach(f => { 
            let s = require(f);
            let id = s["$id"] || s["id"]; 
            if (schemas[id] === undefined) {
                schemas[id] = validator.addSchema(s,id);
            }
        });
    }

    loadServiceSchemas(apiBaseUrl);

    return {
        get: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).get(url)));
        },
        post: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).post(url)));
        },
        put: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).put(url)));
        },
        patch: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).patch(url)));
        },
        del: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).del(url)));
        }
    };
}

exports.client = client;