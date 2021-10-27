!function(e){"use strict";var t=1,n=e.plot.browser,o=n.getPixelRatio;function i(e,i){var l=e.filter(s);t=o(i.getContext("2d"));var u=l.map((function(e){var t=new Image,o=new Promise(function(e,t){return e.sourceDescription='<info className="'+t.className+'" tagName="'+t.tagName+'" id="'+t.id+'">',e.sourceComponent=t,function(o,i){e.onload=function(t){e.successfullyLoaded=!0,o(e)},e.onabort=function(t){e.successfullyLoaded=!1,console.log("Can't generate temp image from "+e.sourceDescription+". It is possible that it is missing some properties or its content is not supported by this browser. Source component:",e.sourceComponent),o(e)},e.onerror=function(t){e.successfullyLoaded=!1,console.log("Can't generate temp image from "+e.sourceDescription+". It is possible that it is missing some properties or its content is not supported by this browser. Source component:",e.sourceComponent),o(e)},function(e,t){"CANVAS"===e.tagName&&(o=e,t.src=o.toDataURL("image/png"));var o;"svg"===e.tagName&&function(e,t){n.isSafari()||n.isMobileSafari()?function(e,t){function n(e){var t="";const n=new Uint8Array(e),o=16384;for(var i=0;i<n.length;i+=o){t+=String.fromCharCode.apply(null,n.subarray(i,i+o))}return t}var o,i,s=g(r(document),e);s=a(s),i=n(new(TextEncoder||TextEncoderLite)("utf-8").encode(s)),o="data:image/svg+xml;base64,"+btoa(i),t.src=o}(e,t):function(e,t){var n=g(r(document),e);n=a(n);var o=new Blob([n],{type:"image/svg+xml;charset=utf-8"}),i=(self.URL||self.webkitURL||self).createObjectURL(o);t.src=i}(e,t)}(e,t);t.srcImgTagName=e.tagName,function(e,t){t.genLeft=e.getBoundingClientRect().left,t.genTop=e.getBoundingClientRect().top,"CANVAS"===e.tagName&&(t.genRight=t.genLeft+e.width,t.genBottom=t.genTop+e.height);"svg"===e.tagName&&(t.genRight=e.getBoundingClientRect().right,t.genBottom=e.getBoundingClientRect().bottom)}(e,t)}(t,e)}}(t,e));return o})),f=Promise.all(u).then(function(e){return function(n){var o=function(e,n){var o=function(e,n){var o=0;if(0===e.length)o=-1;else{var i=e[0].genLeft,s=e[0].genTop,r=e[0].genRight,g=e[0].genBottom,a=0;for(a=1;a<e.length;a++)i>e[a].genLeft&&(i=e[a].genLeft),s>e[a].genTop&&(s=e[a].genTop);for(a=1;a<e.length;a++)r<e[a].genRight&&(r=e[a].genRight),g<e[a].genBottom&&(g=e[a].genBottom);if(r-i<=0||g-s<=0)o=-2;else{for(n.width=Math.round(r-i),n.height=Math.round(g-s),a=0;a<e.length;a++)e[a].xCompOffset=e[a].genLeft-i,e[a].yCompOffset=e[a].genTop-s;!function(e,n){function o(e){return"svg"===e.srcImgTagName}void 0!==n.find(o)&&t<1&&(e.width=e.width*t,e.height=e.height*t)}(n,e)}}return o}(e,n);if(0===o)for(var i=n.getContext("2d"),s=0;s<e.length;s++)!0===e[s].successfullyLoaded&&i.drawImage(e[s],e[s].xCompOffset*t,e[s].yCompOffset*t);return o}(n,e);return o}}(i),c);return f}function s(e){var t=!0,n=!0;return null==e?n=!1:"CANVAS"===e.tagName&&(e.getBoundingClientRect().right!==e.getBoundingClientRect().left&&e.getBoundingClientRect().bottom!==e.getBoundingClientRect().top||(t=!1)),n&&t&&"visible"===window.getComputedStyle(e).visibility}function r(e){for(var t=e.styleSheets,n=[],o=0;o<t.length;o++)try{for(var i=t[o].cssRules||[],s=0;s<i.length;s++){var r=i[s];n.push(r.cssText)}}catch(e){console.log("Failed to get some css rules")}return n}function g(e,n){return['<svg class="snapshot '+n.classList+'" width="'+n.width.baseVal.value*t+'" height="'+n.height.baseVal.value*t+'" viewBox="0 0 '+n.width.baseVal.value+" "+n.height.baseVal.value+'" xmlns="http://www.w3.org/2000/svg">',"<style>","/* <![CDATA[ */",e.join("\n"),"/* ]]> */","</style>",n.innerHTML,"</svg>"].join("\n")}function a(e){var t="";return e.match(/^<svg[^>]+xmlns="http:\/\/www\.w3\.org\/2000\/svg"/)||(t=e.replace(/^<svg/,'<svg xmlns="http://www.w3.org/2000/svg"')),e.match(/^<svg[^>]+"http:\/\/www\.w3\.org\/1999\/xlink"/)||(t=e.replace(/^<svg/,'<svg xmlns:xlink="http://www.w3.org/1999/xlink"')),'<?xml version="1.0" standalone="no"?>\r\n'+t}function c(){return-100}e.plot.composeImages=i,e.plot.plugins.push({init:function(e){e.composeImages=i},name:"composeImages",version:"1.0"})}(jQuery);