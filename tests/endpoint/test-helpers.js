const request = require('supertest');
const Ajv = require('Ajv');
const { resolve } = require('path');
const { readdirSync } = require('fs');

function getFiles(dir,m) {
  const dirents = readdirSync(dir, { withFileTypes: true });
  const files = dirents.map((dirent) => {
    const res = resolve(dir, dirent.name);
    return dirent.isDirectory() ? getFiles(res) : res;
  });
  return Array.prototype.concat(...files).filter(f=>f.match(m));
}

function setHeaders(hdrs, req) {
    for (let p in hdrs) {
        req.set(p, hdrs[p]);
    }
    return req;
}

function addValidator(req, schemas) {
    req['validate'] = (schema) => {
        return req.expect(()=>{
            let key = schema || req.response.headers['x-schema'];
            let tester = schemas[key];
            if (!tester.validate(key, req.response.Body)) {
                throw tester.errors;
            }
        });
    };
    return req;
}

function client(apiBaseUrl, schemaDirectory) {
    const baseUrl = apiBaseUrl;
    const headers = { 'Origin': 'http://unittests.com' };
    const schemas = {};
    let val = new Ajv({allErrors: true});
    if (schemaDirectory) {
        getFiles(schemaDirectory, "json$").forEach(f => { 
            let s = require(f);
            let id = s["$id"] || s["id"]; 
            schemas[id] = val.addSchema(s,id);
        });
    }
    return {
        get: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).get(url)), schemas);
        },
        post: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).post(url)), schemas);
        },
        put: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).put(url)), schemas);
        },
        del: (url) => {
            return addValidator(setHeaders(headers, request(baseUrl).del(url)), schemas);
        }
    };
}

exports.client = client;