!function(t){"function"==typeof define&&define.amd?define(["jquery","datatables.net-bs4","datatables.net-searchbuilder"],(function(e){return t(e)})):"object"==typeof exports?module.exports=function(e,n){return e||(e=window),n&&n.fn.dataTable||(n=require("datatables.net-bs4")(e,n).$),n.fn.dataTable.searchBuilder||require("datatables.net-searchbuilder")(e,n),t(n)}:t(jQuery)}((function(t){"use strict";var e=t.fn.dataTable;return t.extend(!0,e.SearchBuilder.classes,{clearAll:"btn btn-light dtsb-clearAll"}),t.extend(!0,e.Group.classes,{add:"btn btn-light dtsb-add",clearGroup:"btn btn-light dtsb-clearGroup",logic:"btn btn-light dtsb-logic"}),t.extend(!0,e.Criteria.classes,{condition:"form-control dtsb-condition",data:"form-control dtsb-data",delete:"btn btn-light dtsb-delete",left:"btn btn-light dtsb-left",right:"btn btn-light dtsb-right",value:"form-control dtsb-value"}),e.searchPanes}));