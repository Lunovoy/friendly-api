import{r as s,n as P,p as b,j as t,o as k,C as M}from"./index-BiqYdbNS.js";import{u as N,F as I,E as F,a as T,b as h}from"./FiltersBlock-Iz-HQHXj.js";import{F as y}from"./Footer-CSK3jV5d.js";import{P as A}from"./Page-Dv3Yxcn8.js";import{u as L}from"./useTranslation-z379OQBV.js";import"./useTagStore-CfbAclN9.js";import"./Tabs-C9hVvUkL.js";const _="_footer_1kjwg_1",B={footer:_},H=s.memo(f=>{const{className:x}=f,{width:i}=P(),n=s.useMemo(()=>i??0,[i]),{t:S}=L("events-page"),[e,o]=s.useState([]),[v,l]=s.useState(!1),[d,m]=s.useState(""),[g,r]=s.useState(!1),[p,E]=s.useState("active"),u=N(a=>a.deleteEvent),[c]=b();s.useEffect(()=>{const a=c.get("q"),O=c.get("event");a&&m(a),O&&l(!0)},[c]);const j=s.useCallback(async()=>{e!=null&&e.length&&(await Promise.all(e.map(async a=>{await u(a.id)})),o([]))},[u,e]),C=s.useCallback(()=>{o([])},[]),w=s.useCallback(a=>{l(!0)},[]);return t.jsxs(A,{className:k(B.EventsPage,{},[x]),children:[t.jsxs("div",{className:"col-span-4 flex-grow p-0 flex flex-col gap-4",children:[t.jsx(I,{onCreateClick:n>768?()=>r(!0):void 0,searchString:d,setSearchString:m,handleDeleteSelected:j,handleClearAll:C,handleEditSelected:w,selectedEvents:e,activeTab:p,setActiveTab:E}),t.jsx(F,{searchString:d,selectedEvents:e,setSelectedEvents:o,eventsType:p})]}),n>768&&(e!=null&&e.length)?t.jsx(T,{event:e[0],onModalClose:()=>o([])}):n>768?t.jsx(M,{className:"col-span-3",children:t.jsx("h1",{children:S("Выберите событие для просмотра")})}):null,e!=null&&e.length?t.jsx(h,{isOpen:v,setIsOpen:l,event:e[0]}):null,n<=768&&t.jsx(y,{handleAddButtonClick:()=>r(!0),className:"fixed bottom-4"}),t.jsx(h,{isOpen:g,setIsOpen:r})]})});export{H as default};
