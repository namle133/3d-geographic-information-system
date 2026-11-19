import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs';
import vm from 'node:vm';
async function page(role='resident') {
 const elements=new Map();
 const element=id=>{if(!elements.has(id))elements.set(id,{id,style:{},dataset:{},value:'',innerHTML:'',textContent:'',classList:{add(){},remove(){}},addEventListener(){},setAttribute(){},querySelectorAll(){return []},focus(){},showModal(){this.open=true},close(){this.open=false}});return elements.get(id)};
 const items=[{id:'1',title:'<img src=x onerror=alert(1)>',message:'Full message\n\nSecond paragraph <script>bad</script>',type:'maintenance_update',is_read:false,created_at:'2026-11-19'},{id:'2',title:'Invoice ready',message:'Charges are available',type:'invoice_generated',is_read:true,created_at:'2026-11-18'}];
 let readCalls=0;
 const scope={document:{getElementById:element,querySelectorAll(){return []},querySelector(){return element('tab')},createElement(){return element('option')}},window:{dispatchEvent(){}},CustomEvent:class{},initLayout:async()=>true,authManager:{getUserRole:()=>role},notificationAPI:{getUserNotifications:async()=>({data:items,unread_count:items.filter(n=>!n.is_read).length}),markNotificationAsRead:async()=>{readCalls++}},usersAPI:{getUsers:async()=>({data:[]})},formatDate:x=>x,showToast(){}};
 let script=fs.readFileSync(new URL('../pages/notifications.html',import.meta.url),'utf8').match(/<script type="module">([\s\S]*?)<\/script>/)[1].replace(/^\s*import .*;\s*$/mg,'');
 script=script.replace('    loadNotifications();\n  ', '    await loadNotifications();\n  ');
 await new vm.Script('(async()=>{'+script+';window.test={openNotification,renderList,setTab:tab=>{currentTab=tab}};})()').runInNewContext(scope);
 return {element,items,api:scope.window.test,readCalls:()=>readCalls};
}
test('cards escape notification text and search combines with category filters',async()=>{
 const p=await page();assert.doesNotMatch(p.element('notifList').innerHTML,/<img src/);assert.match(p.element('notifList').innerHTML,/&lt;img/);
 p.element('notificationSearch').value='invoice';p.api.renderList();assert.match(p.element('notifList').innerHTML,/Invoice ready/);assert.doesNotMatch(p.element('notifList').innerHTML,/data-id="1"/);
 p.element('notificationCategory').value='maintenance';p.api.renderList();assert.match(p.element('notifList').innerHTML,/Không có thông báo phù hợp/);
});
test('opening a personal notification exposes full plain text and marks it read once',async()=>{
 const p=await page();await p.api.openNotification('1');assert.equal(p.element('notificationDetail').open,true);assert.equal(p.element('notificationDetailMessage').textContent,p.items[0].message);assert.equal(p.readCalls(),1);await p.api.openNotification('1');assert.equal(p.readCalls(),1);
});
test('system scope does not mark another recipient notification as read',async()=>{
 const p=await page('admin');p.api.setTab('system');await p.api.openNotification('1');assert.equal(p.readCalls(),0);assert.equal(p.items[0].is_read,false);
});
