import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';
const html=fs.readFileSync(new URL('../pages/payment-result.html',import.meta.url),'utf8');
const script=html.match(/<script type="module">([\s\S]*?)<\/script>/)[1].replace(/^\s*import .*;\s*$/mg,'').replace('    renderResult();','    await renderResult();');
async function render(responses,query='?vnp_TxnRef=test&vnp_ResponseCode=00&vnp_Amount=999999999'){
 const card={innerHTML:'',setAttribute(){}};let calls=0;
 const scope={URLSearchParams,window:{location:{search:query}},document:{getElementById:id=>id==='resultCard'?card:null},setTimeout:fn=>{fn();return 1;},formatCurrency:n=>String(n),invoiceAPI:{verifyVNPayReturn:async()=>{const r=responses[Math.min(calls++,responses.length-1)];if(r instanceof Error)throw r;return r;}}};
 await new vm.Script('(async()=>{'+script+'})()').runInNewContext(scope);return {html:card.innerHTML,calls};
}
const result=(status,extra={})=>({verified:true,success:status==='success',response_code:'00',transaction_status:'00',data:{id:'test',status,amount:100000,...extra}});
test('unverified query cannot claim payment success',async()=>{
 for(const response of [new Error('invalid signature'),new Error('offline'),{success:true,data:{status:'success'}}]){const r=await render([response]);assert.match(r.html,/Chưa thể xác minh thanh toán/);assert.doesNotMatch(r.html,/Thanh toán thành công/);}
});
test('pending confirmation remains pending and offers status retry',async()=>{const r=await render([result('pending')]);assert.equal(r.calls,5);assert.match(r.html,/Đang chờ xác nhận thanh toán/);assert.match(r.html,/Kiểm tra lại/);});
test('polls until IPN confirms and displays backend amount',async()=>{const r=await render([result('pending'),result('success')]);assert.equal(r.calls,2);assert.match(r.html,/Thanh toán thành công/);assert.match(r.html,/>100000</);assert.doesNotMatch(r.html,/999999999/);});
test('signed failed status never renders success',async()=>{const r=await render([{...result('pending'),transaction_status:'02'}]);assert.match(r.html,/Thanh toán không thành công/);});
test('escapes transaction fields and marks overpayments for review',async()=>{const r=await render([result('success',{id:'<img onerror=alert(1)>',requires_review:true})]);assert.match(r.html,/Thanh toán cần đối soát/);assert.doesNotMatch(r.html,/<img/);assert.match(r.html,/&lt;img/);});
