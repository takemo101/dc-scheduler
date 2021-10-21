!function(e){"use strict";var t=e.plot.saturated.floorInBase,i=function(e,t){var i=new e(t),o=i.setTime.bind(i);i.update=function(e){o(e),e=Math.round(1e3*e)/1e3,this.microseconds=1e3*(e-Math.floor(e))};var n=i.getTime.bind(i);return i.getTime=function(){return n()+this.microseconds/1e3},i.setTime=function(e){this.update(e)},i.getMicroseconds=function(){return this.microseconds},i.setMicroseconds=function(e){var t=n()+e/1e3;this.update(t)},i.setUTCMicroseconds=function(e){this.setMicroseconds(e)},i.getUTCMicroseconds=function(){return this.getMicroseconds()},i.microseconds=null,i.microEpoch=null,i.update(t),i};function o(e,t,i,o){if("function"==typeof e.strftime)return e.strftime(t);var n,s=function(e,t){return t=""+(null==t?"0":t),1===(e=""+e).length?t+e:e},r=function(e,t,i){var o,n=1e3*e+t;if(i<6&&i>0){var s=parseFloat("1e"+(i-6));o=("00000"+(n=Math.round(Math.round(n*s)/s))).slice(-6,-(6-i))}else o=("00000"+(n=Math.round(n))).slice(-6);return o},a=[],c=!1,m=e.getHours(),u=m<12;i||(i=["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"]),o||(o=["Sun","Mon","Tue","Wed","Thu","Fri","Sat"]),n=m>12?m-12:0===m?12:m;for(var d=-1,l=0;l<t.length;++l){var h=t.charAt(l);if(!isNaN(Number(h))&&Number(h)>0)d=Number(h);else if(c){switch(h){case"a":h=""+o[e.getDay()];break;case"b":h=""+i[e.getMonth()];break;case"d":h=s(e.getDate());break;case"e":h=s(e.getDate()," ");break;case"h":case"H":h=s(m);break;case"I":h=s(n);break;case"l":h=s(n," ");break;case"m":h=s(e.getMonth()+1);break;case"M":h=s(e.getMinutes());break;case"q":h=""+(Math.floor(e.getMonth()/3)+1);break;case"S":h=s(e.getSeconds());break;case"s":h=""+r(e.getMilliseconds(),e.getMicroseconds(),d);break;case"y":h=s(e.getFullYear()%100);break;case"Y":h=""+e.getFullYear();break;case"p":h=u?"am":"pm";break;case"P":h=u?"AM":"PM";break;case"w":h=""+e.getDay()}a.push(h),c=!1}else"%"===h?c=!0:a.push(h)}return a.join("")}function n(e){function t(e,t,i,o){e[t]=function(){return i[o].apply(i,arguments)}}var i={date:e};void 0!==e.strftime&&t(i,"strftime",e,"strftime"),t(i,"getTime",e,"getTime"),t(i,"setTime",e,"setTime");for(var o=["Date","Day","FullYear","Hours","Minutes","Month","Seconds","Milliseconds","Microseconds"],n=0;n<o.length;n++)t(i,"get"+o[n],e,"getUTC"+o[n]),t(i,"set"+o[n],e,"setUTC"+o[n]);return i}function s(e,t){var o=864e13;if(t&&"seconds"===t.timeBase?e*=1e3:"microseconds"===t.timeBase&&(e/=1e3),e>o?e=o:e<-o&&(e=-o),"browser"===t.timezone)return i(Date,e);if(t.timezone&&"utc"!==t.timezone){if("undefined"!=typeof timezoneJS&&void 0!==timezoneJS.Date){var s=i(timezoneJS.Date,e);return s.setTimezone(t.timezone),s.setTime(e),s}return n(i(Date,e))}return n(i(Date,e))}var r={microsecond:1e-6,millisecond:.001,second:1,minute:60,hour:3600,day:86400,month:2592e3,quarter:7776e3,year:525949.2*60},a={microsecond:.001,millisecond:1,second:1e3,minute:6e4,hour:36e5,day:864e5,month:2592e6,quarter:7776e6,year:525949.2*60*1e3},c={microsecond:1,millisecond:1e3,second:1e6,minute:6e7,hour:36e8,day:864e8,month:2592e9,quarter:7776e9,year:31556951999999.996},m=[[1,"microsecond"],[2,"microsecond"],[5,"microsecond"],[10,"microsecond"],[25,"microsecond"],[50,"microsecond"],[100,"microsecond"],[250,"microsecond"],[500,"microsecond"],[1,"millisecond"],[2,"millisecond"],[5,"millisecond"],[10,"millisecond"],[25,"millisecond"],[50,"millisecond"],[100,"millisecond"],[250,"millisecond"],[500,"millisecond"],[1,"second"],[2,"second"],[5,"second"],[10,"second"],[30,"second"],[1,"minute"],[2,"minute"],[5,"minute"],[10,"minute"],[30,"minute"],[1,"hour"],[2,"hour"],[4,"hour"],[8,"hour"],[12,"hour"],[1,"day"],[2,"day"],[3,"day"],[.25,"month"],[.5,"month"],[1,"month"],[2,"month"]],u=m.concat([[3,"month"],[6,"month"],[1,"year"]]),d=m.concat([[1,"quarter"],[2,"quarter"],[1,"year"]]);function l(e){var i,o=e.options,n=[],m=s(e.min,o),l=0,h=o.tickSize&&"quarter"===o.tickSize[1]||o.minTickSize&&"quarter"===o.minTickSize[1]?d:u;i="seconds"===o.timeBase?r:"microseconds"===o.timeBase?c:a,null!==o.minTickSize&&void 0!==o.minTickSize&&(l="number"==typeof o.tickSize?o.tickSize:o.minTickSize[0]*i[o.minTickSize[1]]);for(var f=0;f<h.length-1&&!(e.delta<(h[f][0]*i[h[f][1]]+h[f+1][0]*i[h[f+1][1]])/2&&h[f][0]*i[h[f][1]]>=l);++f);var M=h[f][0],g=h[f][1];if("year"===g){if(null!==o.minTickSize&&void 0!==o.minTickSize&&"year"===o.minTickSize[1])M=Math.floor(o.minTickSize[0]);else{var k=parseFloat("1e"+Math.floor(Math.log(e.delta/i.year)/Math.LN10)),p=e.delta/i.year/k;M=p<1.5?1:p<3?2:p<7.5?5:10,M*=k}M<1&&(M=1)}e.tickSize=o.tickSize||[M,g];var y=e.tickSize[0],S=y*i[g=e.tickSize[1]];"microsecond"===g?m.setMicroseconds(t(m.getMicroseconds(),y)):"millisecond"===g?m.setMilliseconds(t(m.getMilliseconds(),y)):"second"===g?m.setSeconds(t(m.getSeconds(),y)):"minute"===g?m.setMinutes(t(m.getMinutes(),y)):"hour"===g?m.setHours(t(m.getHours(),y)):"month"===g?m.setMonth(t(m.getMonth(),y)):"quarter"===g?m.setMonth(3*t(m.getMonth()/3,y)):"year"===g&&m.setFullYear(t(m.getFullYear(),y)),S>=i.millisecond&&(S>=i.second?m.setMicroseconds(0):m.setMicroseconds(1e3*m.getMilliseconds())),S>=i.minute&&m.setSeconds(0),S>=i.hour&&m.setMinutes(0),S>=i.day&&m.setHours(0),S>=4*i.day&&m.setDate(1),S>=2*i.month&&m.setMonth(t(m.getMonth(),3)),S>=2*i.quarter&&m.setMonth(t(m.getMonth(),6)),S>=i.year&&m.setMonth(0);var v,T,z=0,b=Number.NaN;do{if(T=b,v=m.getTime(),b=o&&"seconds"===o.timeBase?v/1e3:o&&"microseconds"===o.timeBase?1e3*v:v,n.push(b),"month"===g||"quarter"===g)if(y<1){m.setDate(1);var q=m.getTime();m.setMonth(m.getMonth()+("quarter"===g?3:1));var B=m.getTime();m.setTime(b+z*i.hour+(B-q)*y),z=m.getHours(),m.setHours(0)}else m.setMonth(m.getMonth()+y*("quarter"===g?3:1));else"year"===g?m.setFullYear(m.getFullYear()+y):"seconds"===o.timeBase?m.setTime(1e3*(b+S)):"microseconds"===o.timeBase?m.setTime((b+S)/1e3):m.setTime(b+S)}while(b<e.max&&b!==T);return n}e.plot.plugins.push({init:function(t){t.hooks.processOptions.push((function(t){e.each(t.getAxes(),(function(e,t){var i=t.options;if("time"===i.mode){if(t.tickGenerator=l,"tickFormatter"in i&&"function"==typeof i.tickFormatter)return;t.tickFormatter=function(e,t){var n=s(e,t.options);if(null!=i.timeformat)return o(n,i.timeformat,i.monthNames,i.dayNames);var m,u=t.options.tickSize&&"quarter"===t.options.tickSize[1]||t.options.minTickSize&&"quarter"===t.options.minTickSize[1];m="seconds"===i.timeBase?r:"microseconds"===i.timeBase?c:a;var d,l,h=t.tickSize[0]*m[t.tickSize[1]],f=t.max-t.min,M=i.twelveHourClock?" %p":"",g=i.twelveHourClock?"%I":"%H";if(d="seconds"===i.timeBase?1:"microseconds"===i.timeBase?1e6:1e3,h<m.second){var k=-Math.floor(Math.log10(h/d));String(h).indexOf("25")>-1&&k++,l="%S.%"+k+"s"}else l=h<m.minute?g+":%M:%S"+M:h<m.day?f<2*m.day?g+":%M"+M:"%b %d "+g+":%M"+M:h<m.month?"%b %d":u&&h<m.quarter||!u&&h<m.year?f<m.year?"%b":"%b %Y":u&&h<m.year?f<m.year?"Q%q":"Q%q %Y":"%Y";return o(n,l,i.monthNames,i.dayNames)}}}))}))},options:{xaxis:{timezone:null,timeformat:null,twelveHourClock:!1,monthNames:null,timeBase:"seconds"},yaxis:{timeBase:"seconds"}},name:"time",version:"1.0"}),e.plot.formatDate=o,e.plot.dateGenerator=s,e.plot.dateTickGenerator=l,e.plot.makeUtcWrapper=n}(jQuery);