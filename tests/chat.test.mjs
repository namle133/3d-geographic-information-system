import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';

async function page(role = 'resident', overrides = {}, media = {}) {
  const elements = new Map();
  const element = id => {
    if (!elements.has(id)) elements.set(id, {value:'', innerHTML:'', textContent:'', hidden:false, disabled:false, scrollTop:0, scrollHeight:400, clientHeight:400,
      events:{}, classList:{add(){}}, addEventListener(type,fn){this.events[type]=fn}, querySelectorAll(){return []}, focus(){}});
    return elements.get(id);
  };
  const calls = [], items = [{id:'one', body:'<img src=x onerror=alert(1)>', from_admin:true, created_at:'2026-11-19T01:00:00Z'}];
  const api = {start:async()=>({data:{id:'resident-thread',status:'open'}}),people:async()=>({data:[]}),messages:async()=>({data:items,conversation:{id:'resident-thread',status:'open'},has_more:false}),read:async()=>{},send:async(id,body,nonce,imageIDs)=>{calls.push({id,body,nonce,imageIDs});return {data:{id:'sent',body,from_admin:false,created_at:'2026-11-19T02:00:00Z'}}},...overrides};
  const scope={imageGallery:()=>"",hydrateImages:()=>{},imagesAPI:{...media},document:{getElementById:element,querySelectorAll:()=>[],hidden:false},window:{addEventListener(){}},initLayout:async()=>true,authManager:{getUserRole:()=>role},chatAPI:api,formatDate:x=>x,crypto:{randomUUID:()=>String(Math.random())},setInterval(){}};
  const script=fs.readFileSync(new URL('../js/modules/chat.js',import.meta.url),'utf8').replace(/^import .*;\s*$/mg,'');
  await new vm.Script('(async()=>{'+script+';return {selectPerson,refreshMessages,drawPeople};})()').runInNewContext(scope).then(result=>scope.controls=result);
  return {element, calls, controls:scope.controls};
}
test('message text is escaped and residents load only their own conversation',async()=>{
  let opened;
  const p=await page('resident',{start:async user=>{opened=user;return {data:{id:'resident-thread'}}}});
  assert.equal(opened,undefined);
  assert.match(p.element('chatMessages').innerHTML,/&lt;img/);
  assert.doesNotMatch(p.element('chatMessages').innerHTML,/<img/);
});
test('a failed send retains the draft and reuses its idempotency reference',async()=>{
  const sent=[];
  const p=await page('resident',{send:async(id,body,nonce)=>{sent.push(nonce);if(sent.length===1)throw new Error('Offline');return {data:{id:'saved',body,created_at:'2026-11-19',from_admin:false}}}});
  p.element('chatDraft').value='Please help';
  const submit=()=>p.element('chatForm').events.submit({preventDefault(){}});
  await submit();assert.equal(p.element('chatDraft').value,'Please help');assert.equal(p.element('chatError').hidden,false);
  await submit();assert.equal(sent[0],sent[1]);assert.equal(p.element('chatDraft').value,'');
});
test('switching conversations discards late responses from the previous user',async()=>{
  let release;
  const p=await page('admin',{start:async id=>id==='slow'?await new Promise(resolve=>release=resolve):{data:{id}},messages:async id=>({data:[{id,body:id,created_at:'2026-11-19',from_admin:false}],conversation:{id,status:'open'}})});
  const slow=p.controls.selectPerson({id:'slow',username:'Slow',account_status:'active'});
  await p.controls.selectPerson({id:'fast',username:'Fast',account_status:'active'});
  release({data:{id:'slow'}});await slow;
  assert.equal(p.element('chatName').textContent,'Fast');
  assert.match(p.element('chatMessages').innerHTML,/fast/);assert.doesNotMatch(p.element('chatMessages').innerHTML,/slow/);
});

test('an uploaded image can be sent without text and stays attached to its conversation',async()=>{
 const id='12345678-1234-4234-8234-123456789abc';
 const p=await page('resident',{}, {upload:async()=>({data:{id}})});
 await p.element('chatImageInput').events.change({target:{files:[{type:'image/png'}],value:'sample'}});
 assert.equal(p.element('chatDraft').required,false);
 await p.element('chatForm').events.submit({preventDefault(){}});
 assert.deepEqual([...p.calls[0].imageIDs],[id]);assert.equal(p.calls[0].body,'');
 assert.equal(p.element('chatDraft').required,true);
});
