webpackJsonp([1],{"3fs2":function(e,t,s){var o=s("RY/4"),a=s("dSzd")("iterator"),r=s("/bQp");e.exports=s("FeBl").getIteratorMethod=function(e){if(void 0!=e)return e[a]||e["@@iterator"]||r[o(e)]}},BO1k:function(e,t,s){e.exports={default:s("fxRn"),__esModule:!0}},"RY/4":function(e,t,s){var o=s("R9M2"),a=s("dSzd")("toStringTag"),r="Arguments"==o(function(){return arguments}());e.exports=function(e){var t,s,n;return void 0===e?"Undefined":null===e?"Null":"string"==typeof(s=function(e,t){try{return e[t]}catch(e){}}(t=Object(e),a))?s:r?o(t):"Object"==(n=o(t))&&"function"==typeof t.callee?"Arguments":n}},Znwo:function(e,t){},fxRn:function(e,t,s){s("+tPU"),s("zQR9"),e.exports=s("g8Ux")},g8Ux:function(e,t,s){var o=s("77Pl"),a=s("3fs2");e.exports=s("FeBl").getIterator=function(e){var t=a(e);if("function"!=typeof t)throw TypeError(e+" is not iterable!");return o(t.call(e))}},iuFH:function(e,t,s){"use strict";Object.defineProperty(t,"__esModule",{value:!0});var o=s("BO1k"),a=s.n(o),r={data:function(){return{saved:null,customColors:[{color:"#32cd32",percentage:"50"},{color:"#ffa500",percentage:"85"},{color:"#ff0000",percentage:"100"}],expands:[],timer:""}},created:function(){this.fetchData(),this.timer=setInterval(this.fetchData,1e3)},watch:{},methods:{fetchData:function(){var e=this;this.$axios.get("/api/servers").then(function(t){e.saved=[];var s=!0,o=!1,r=void 0;try{for(var n,l=a()(t.data.tmpServers);!(s=(n=l.next()).done);s=!0){var i=n.value,d=(new Date).getTime();Math.floor(d/1e3)-i.lastTime>30?(i.status="Offline",i.IPv4="Off",i.IPv6="Off",i.system="-",i.uptime="-",i.downloadSpeed="-",i.uploadSpeed="-",i.CPUUsed=NaN,i.memoryRate=NaN,i.diskRate=NaN,i.memoryUsed=i.memoryTotal=i.swapUsed=i.swapTotal=i.diskUsed=i.diskTotal=0):(i.status="Online",i.IPv4=""===i.IPv4Addr?"Off":"On",i.IPv6=""===i.IPv6Addr?"Off":"On",i.memoryRate=parseInt(i.memoryUsed/i.memoryTotal*100),i.diskRate=parseInt(i.diskUsed/i.diskTotal*100),i.memoryUsed=i.memoryUsed.toFixed(2),i.memoryTotal=i.memoryTotal.toFixed(2),i.swapUsed=i.swapUsed.toFixed(2),i.swapTotal=i.swapTotal.toFixed(2),i.diskUsed=i.diskUsed.toFixed(2),i.diskTotal=i.diskTotal.toFixed(2)),e.saved.push(i)}}catch(e){o=!0,r=e}finally{try{!s&&l.return&&l.return()}finally{if(o)throw r}}})},Save:function(e){var t=this;this.$axios.post("/api/servers",{from:"temp",action:"save",UUID:e}).then(function(s){"success"===s.data?(t.fetchData(),t.$message({showClose:!0,message:"Save success!",type:"success"}),t.expands.splice(t.expands.indexOf(e),1)):"fail"===s.data&&t.$message({showClose:!0,message:"Save fail!",type:"error"})})},Delete:function(e){var t=this;this.$axios.post("/api/servers",{from:"temp",action:"delete",UUID:e}).then(function(s){"success"===s.data?(t.fetchData(),t.$message({showClose:!0,message:"Delete success!",type:"success"}),t.expands.splice(t.expands.indexOf(e),1)):"fail"===s.data&&t.$message({showClose:!0,message:"Delete fail!",type:"error"})})},rowClick:function(e,t,s){this.expands.indexOf(e.UUID)<0?this.expands.push(e.UUID):this.expands.splice(this.expands.indexOf(e.UUID),1)}}},n={render:function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("el-table",{staticStyle:{width:"100%"},attrs:{data:e.saved,"default-sort":{prop:"status",order:"descending"},"row-key":"UUID","expand-row-keys":e.expands},on:{"row-click":e.rowClick}},[s("el-table-column",{attrs:{type:"expand"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("el-button",{attrs:{type:"success",icon:"el-icon-check",round:""},on:{click:function(s){return e.Save(t.row.UUID)}}},[e._v("Save")]),e._v(" "),s("el-button",{attrs:{type:"danger",icon:"el-icon-delete",round:""},on:{click:function(s){return e.Delete(t.row.UUID)}}},[e._v("Delete")]),e._v(" "),s("el-form",{directives:[{name:"show",rawName:"v-show",value:"Online"===t.row.status,expression:"props.row.status === 'Online'? true : false"}],staticClass:"demo-table-expand",attrs:{"label-position":"left","label-width":"80px"}},[s("el-form-item",{attrs:{label:"Memory"}},[s("span",[e._v(e._s(t.row.memoryUsed)+"GiB / "+e._s(t.row.memoryTotal)+"GiB")])]),e._v(" "),s("el-form-item",{attrs:{label:"Swap"}},[s("span",[e._v(e._s(t.row.swapUsed)+"GiB / "+e._s(t.row.swapTotal)+"GiB")])]),e._v(" "),s("el-form-item",{attrs:{label:"Disk"}},[s("span",[e._v(e._s(t.row.diskUsed)+"GiB / "+e._s(t.row.diskTotal)+"GiB")])])],1)]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"Name",prop:"serverName"}}),e._v(" "),s("el-table-column",{attrs:{label:"Status",prop:"status",sortable:""},scopedSlots:e._u([{key:"default",fn:function(t){return[s("el-tag",{attrs:{type:"Online"===t.row.status?"success":"danger","disable-transitions":""}},[e._v(e._s(t.row.status))])]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"IPv4",prop:"IPv4"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("el-tag",{attrs:{type:"On"===t.row.IPv4?"success":"danger","disable-transitions":""}},[e._v(e._s(t.row.IPv4))])]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"IPv6",prop:"IPv6"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("el-tag",{attrs:{type:"On"===t.row.IPv6?"success":"danger","disable-transitions":""}},[e._v(e._s(t.row.IPv6))])]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"System",prop:"system"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("pre",[e._v(e._s(t.row.system))])]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"Location",prop:"location"}}),e._v(" "),s("el-table-column",{attrs:{label:"Uptime",prop:"uptime"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("pre",[e._v(e._s(t.row.uptime))])]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"Network",prop:"network"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("p",[e._v("⬆️"+e._s(t.row.uploadSpeed))]),e._v(" "),s("p",[e._v("⬇️"+e._s(t.row.downloadSpeed))])]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"CPU",prop:"CPUUsed"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("el-progress",{attrs:{"text-inside":!0,"stroke-width":26,percentage:t.row.CPUUsed,color:e.customColors}})]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"RAM",prop:"memoryRate"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("el-progress",{attrs:{"text-inside":!0,"stroke-width":26,percentage:t.row.memoryRate,color:e.customColors}})]}}])}),e._v(" "),s("el-table-column",{attrs:{label:"Disk",prop:"diskRate"},scopedSlots:e._u([{key:"default",fn:function(t){return[s("el-progress",{attrs:{"text-inside":!0,"stroke-width":26,percentage:t.row.diskRate,color:e.customColors}})]}}])})],1)},staticRenderFns:[]};var l=s("VU/8")(r,n,!1,function(e){s("Znwo")},null,null);t.default=l.exports}});
//# sourceMappingURL=1.361b18633f2dd2dbf562.js.map