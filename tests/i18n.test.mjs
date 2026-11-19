import test from 'node:test';
import assert from 'node:assert/strict';
import vm from 'node:vm';
import fs from 'node:fs';
function translator(language) {
 const document={readyState:'loading',addEventListener(){},documentElement:{lang:language},querySelectorAll(){return []},body:{nodeType:1,tagName:'BODY',matches(){return false},hasAttribute(){return false},childNodes:[]},title:''};
 const scope={window:{dispatchEvent(){}},document,localStorage:{getItem:()=>language,setItem(){}},Node:{TEXT_NODE:3,ELEMENT_NODE:1},CustomEvent:class{}};
 vm.runInNewContext(fs.readFileSync(new URL('../js/i18n.js',import.meta.url),'utf8'),scope);
 return scope.window.SiteI18n;
}
test('legacy error aliases translate in the correct direction',()=>{
 const t=translator('vi');const vi='Email hoặc mật khẩu không đúng';
 assert.equal(t.translate(vi),vi);assert.equal(t.translate('Failed to create user'),'Không thể tạo người dùng');
 t.setLanguage('en');assert.equal(t.translate(vi),'Email or password is incorrect');assert.equal(t.translate('Không thể tạo người dùng'),'Failed to create user');
 t.setLanguage('vi');assert.equal(t.translate(vi),vi);
});
test('GIS and recovery labels have complete translations',()=>{
 const t=translator('en');for(const text of ['Cột','Dốc','Lõi','Ngôn ngữ','Nhập email đã đăng ký. Mã đặt lại mật khẩu sẽ được gửi tới email đã đăng ký.']) assert.doesNotMatch(t.translate(text),/[À-ỹĐđ]/);
});
test('all pages load localization and the 3D selector has a dedicated visible host',()=>{
 for(const path of ['index.html',...fs.readdirSync(new URL('../pages/',import.meta.url)).filter(p=>p.endsWith('.html')).map(p=>'pages/'+p)]) assert.match(fs.readFileSync(new URL('../'+path,import.meta.url),'utf8'),/js\/i18n\.js/);
 assert.match(fs.readFileSync(new URL('../index.html',import.meta.url),'utf8'),/gis-language-control/);
});
test('Vietnamese backend response messages translate completely into English',()=>{
 const t=translator('en');
 for(const dir of ['handlers','services']) for(const file of fs.readdirSync(new URL('../backend/internal/'+dir+'/',import.meta.url)).filter(f=>f.endsWith('.go')&&!f.endsWith('_test.go'))){
  const source=fs.readFileSync(new URL('../backend/internal/'+dir+'/'+file,import.meta.url),'utf8');
  for(const [,text] of source.matchAll(/"(?:error|message|title)"\s*:\s*"([^"\n]+)"/g)) if(/[À-ỹĐđ]/.test(text)) assert.doesNotMatch(t.translate(text),/[À-ỹĐđ]/,file+': '+text);
 }
});
