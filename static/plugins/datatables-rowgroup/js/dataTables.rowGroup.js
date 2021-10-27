/*! For license information please see dataTables.rowGroup.js.LICENSE.txt */
!function(t){"function"==typeof define&&define.amd?define(["jquery","datatables.net"],(function(r){return t(r,window,document)})):"object"==typeof exports?module.exports=function(r,e){return r||(r=window),e&&e.fn.dataTable||(e=require("datatables.net")(r,e).$),t(e,r,r.document)}:t(jQuery,window,document)}((function(t,r,e,n){"use strict";var o=t.fn.dataTable,a=function(r,e){if(!o.versionCheck||!o.versionCheck("1.10.8"))throw"RowGroup requires DataTables 1.10.8 or newer";this.c=t.extend(!0,{},o.defaults.rowGroup,a.defaults,e),this.s={dt:new o.Api(r)},this.dom={};var n=this.s.dt.settings()[0],s=n.rowGroup;if(s)return s;n.rowGroup=this,this._constructor()};return t.extend(a.prototype,{dataSrc:function(r){if(r===n)return this.c.dataSrc;var e=this.s.dt;return this.c.dataSrc=r,t(e.table().node()).triggerHandler("rowgroup-datasrc.dt",[e,r]),this},disable:function(){return this.c.enable=!1,this},enable:function(t){return!1===t?this.disable():(this.c.enable=!0,this)},enabled:function(){return this.c.enable},_constructor:function(){var t=this,r=this.s.dt,e=r.settings()[0];r.on("draw.dtrg",(function(r,n){t.c.enable&&e===n&&t._draw()})),r.on("column-visibility.dt.dtrg responsive-resize.dt.dtrg",(function(){t._adjustColspan()})),r.on("destroy",(function(){r.off(".dtrg")}))},_adjustColspan:function(){t("tr."+this.c.className,this.s.dt.table().body()).find("td:visible").attr("colspan",this._colspan())},_colspan:function(){return this.s.dt.columns().visible().reduce((function(t,r){return t+r}),0)},_draw:function(){var t=this.s.dt,r=this._group(0,t.rows({page:"current"}).indexes());this._groupDisplay(0,r)},_group:function(t,r){for(var e,a=Array.isArray(this.c.dataSrc)?this.c.dataSrc:[this.c.dataSrc],s=o.ext.oApi._fnGetObjectDataFn(a[t]),i=this.s.dt,u=[],d=0,c=r.length;d<c;d++){var l,p=r[d];null!==(l=s(i.row(p).data()))&&l!==n||(l=this.c.emptyDataGroup),e!==n&&l===e||(u.push({dataPoint:l,rows:[]}),e=l),u[u.length-1].rows.push(p)}if(a[t+1]!==n)for(d=0,c=u.length;d<c;d++)u[d].children=this._group(t+1,u[d].rows);return u},_groupDisplay:function(t,r){for(var e,n=this.s.dt,o=0,a=r.length;o<a;o++){var s,i=r[o],u=i.dataPoint,d=i.rows;this.c.startRender&&(e=this.c.startRender.call(this,n.rows(d),u,t),(s=this._rowWrap(e,this.c.startClassName,t))&&s.insertBefore(n.row(d[0]).node())),this.c.endRender&&(e=this.c.endRender.call(this,n.rows(d),u,t),(s=this._rowWrap(e,this.c.endClassName,t))&&s.insertAfter(n.row(d[d.length-1]).node())),i.children&&this._groupDisplay(t+1,i.children)}},_rowWrap:function(r,e,o){return null!==r&&""!==r||(r=this.c.emptyDataGroup),r===n||null===r?null:("object"==typeof r&&r.nodeName&&"tr"===r.nodeName.toLowerCase()?t(r):r instanceof t&&r.length&&"tr"===r[0].nodeName.toLowerCase()?r:t("<tr/>").append(t("<td/>").attr("colspan",this._colspan()).append(r))).addClass(this.c.className).addClass(e).addClass("dtrg-level-"+o)}}),a.defaults={className:"dtrg-group",dataSrc:0,emptyDataGroup:"No group",enable:!0,endClassName:"dtrg-end",endRender:null,startClassName:"dtrg-start",startRender:function(t,r){return r}},a.version="1.1.3",t.fn.dataTable.RowGroup=a,t.fn.DataTable.RowGroup=a,o.Api.register("rowGroup()",(function(){return this})),o.Api.register("rowGroup().disable()",(function(){return this.iterator("table",(function(t){t.rowGroup&&t.rowGroup.enable(!1)}))})),o.Api.register("rowGroup().enable()",(function(t){return this.iterator("table",(function(r){r.rowGroup&&r.rowGroup.enable(t===n||t)}))})),o.Api.register("rowGroup().enabled()",(function(){var t=this.context;return!(!t.length||!t[0].rowGroup)&&t[0].rowGroup.enabled()})),o.Api.register("rowGroup().dataSrc()",(function(t){return t===n?this.context[0].rowGroup.dataSrc():this.iterator("table",(function(r){r.rowGroup&&r.rowGroup.dataSrc(t)}))})),t(e).on("preInit.dt.dtrg",(function(r,e,n){if("dt"===r.namespace){var s=e.oInit.rowGroup,i=o.defaults.rowGroup;if(s||i){var u=t.extend({},i,s);!1!==s&&new a(e,u)}}})),a}));