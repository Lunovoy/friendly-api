import{r as s,n as P,p as b,j as t,o as M,C as k}from"./index-CvvDs7yX.js";import{u as N,F as I,E as F,a as T,b as h}from"./FiltersBlock-B0bDtalD.js";import{F as y}from"./Footer-DPWHb4zx.js";import{P as A}from"./Page-BGgwBIZw.js";import{u as L}from"./useTranslation-DdVfk-uP.js";import"./useTagStore-uBXJ3nh_.js";import"./Tabs-DUq1uYSp.js";const _="_footer_1kjwg_1",B={footer:_},H=s.memo(f=>{const{className:x}=f,{width:r}=P(),o=s.useMemo(()=>r??0,[r]),{t:S}=L("events-page"),[e,n]=s.useState([]),[v,l]=s.useState(!1),[i,d]=s.useState(""),[g,m]=s.useState(!1),[p,E]=s.useState("active"),u=N(a=>a.deleteEvent),[c]=b();s.useEffect(()=>{const a=c.get("q"),O=c.get("event");a&&d(a),O&&l(!0)},[c]);const j=s.useCallback(async()=>{e!=null&&e.length&&(await Promise.all(e.map(async a=>{await u(a.id)})),n([]))},[u,e]),C=s.useCallback(()=>{n([])},[]),w=s.useCallback(a=>{l(!0)},[]);return t.jsxs(A,{className:M(B.EventsPage,{},[x]),children:[t.jsxs("div",{className:"col-span-4 flex-grow p-0 flex flex-col gap-4",children:[t.jsx(I,{searchString:i,setSearchString:d,handleDeleteSelected:j,handleClearAll:C,handleEditSelected:w,selectedEvents:e,activeTab:p,setActiveTab:E}),t.jsx(F,{searchString:i,selectedEvents:e,setSelectedEvents:n,eventsType:p})]}),o>768&&(e!=null&&e.length)?t.jsx(T,{event:e[0],onModalClose:()=>n([])}):o>768?t.jsx(k,{className:"col-span-3",children:t.jsx("h1",{children:S("Выберите событие для просмотра")})}):null,e!=null&&e.length?t.jsx(h,{isOpen:v,setIsOpen:l,event:e[0]}):null,o<=768&&t.jsx(y,{handleAddButtonClick:()=>m(!0),className:"fixed bottom-4",children:t.jsx(h,{isOpen:g,setIsOpen:m})})]})});export{H as default};