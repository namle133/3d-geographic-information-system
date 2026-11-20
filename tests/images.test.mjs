import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';
function media(){const source=fs.readFileSync(new URL('../js/api/images.js',import.meta.url),'utf8').replace(/^import .*;\s*$/mg,'').replace(/export (async )?function/g,'$1function').replace('export default','const imagesAPI =');const scope={window:{addEventListener(){}},URL,Map,Promise,FormData,apiClient:{}};vm.runInNewContext(source+';globalThis.test={validateImage,imageGallery};',scope);return scope.test;}
test('limits supported image types and file sizes before upload',()=>{const m=media();assert.throws(()=>m.validateImage({type:'image/svg+xml',size:10}));assert.throws(()=>m.validateImage({type:'image/jpeg',size:2*1024*1024+1}));m.validateImage({type:'image/png',size:1024});});
test('galleries accept generated image references without embedding untrusted URLs',()=>{const m=media();const html=m.imageGallery(['<img src=x onerror=alert(1)>','https://example.com/photo.png','12345678-1234-4234-8234-123456789abc']);assert.doesNotMatch(html,/onerror|example.com|src=/);assert.match(html,/data-image-id="12345678/);});
test('multipart requests preserve FormData and do not send a JSON content type',async()=>{
 const form=new FormData();form.append('purpose','maintenance');let sent;
 const scope={API_CONFIG:{BASE_URL:'http://test',TIMEOUT:1000},STORAGE_KEYS:{},localStorage:{getItem:()=>null},window:{},FormData,AbortController,setTimeout,clearTimeout,fetch:async(url,config)=>{sent=config;return {ok:true,status:201,headers:{get:()=> 'application/json'},json:async()=>({data:{id:'saved'}})}}};
 let source=fs.readFileSync(new URL('../js/api/client.js',import.meta.url),'utf8').replace(/^import .*;\s*$/mg,'').replace(/^export .*;\s*$/mg,'');
 vm.runInNewContext(source+';globalThis.client=apiClient;',scope);
 await scope.client.fetch('/images',{method:'POST',body:form,headers:{Authorization:'Bearer test'}});
 assert.equal(sent.body,form);assert.equal(sent.headers['Content-Type'],undefined);
});
