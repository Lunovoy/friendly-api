import{g as r,h as f,j as o}from"./import-CewOev11.js";import{I as K,r as y}from"./index-DP60Imk7.js";class c{*[Symbol.iterator](){yield*this.iterable}get size(){return this.keyMap.size}getKeys(){return this.keyMap.keys()}getKeyBefore(e){let t=this.keyMap.get(e);return t?t.prevKey:null}getKeyAfter(e){let t=this.keyMap.get(e);return t?t.nextKey:null}getFirstKey(){return this.firstKey}getLastKey(){return this.lastKey}getItem(e){return this.keyMap.get(e)}at(e){const t=[...this.getKeys()];return this.getItem(t[e])}constructor(e,{expandedKeys:t}={}){this.keyMap=new Map,this.iterable=e,t=t||new Set;let l=a=>{if(this.keyMap.set(a.key,a),a.childNodes&&(a.type==="section"||t.has(a.key)))for(let d of a.childNodes)l(d)};for(let a of e)l(a);let n,i=0;for(let[a,d]of this.keyMap)n?(n.nextKey=a,d.prevKey=n.key):(this.firstKey=a,d.prevKey=void 0),d.type==="item"&&(d.index=i++),n=d,n.nextKey=void 0;this.lastKey=n==null?void 0:n.key}}function x(s){let[e,t]=K(s.expandedKeys?new Set(s.expandedKeys):void 0,s.defaultExpandedKeys?new Set(s.defaultExpandedKeys):new Set,s.onExpandedChange),l=r(s),n=y.useMemo(()=>s.disabledKeys?new Set(s.disabledKeys):new Set,[s.disabledKeys]),i=f(s,y.useCallback(d=>new c(d,{expandedKeys:e}),[e]),null);return y.useEffect(()=>{l.focusedKey!=null&&!i.getItem(l.focusedKey)&&l.setFocusedKey(null)},[i,l.focusedKey]),{collection:i,expandedKeys:e,disabledKeys:n,toggleKey:d=>{t(u(e,d))},setExpandedKeys:t,selectionManager:new o(i,l)}}function u(s,e){let t=new Set(s);return t.has(e)?t.delete(e):t.add(e),t}export{x as $};