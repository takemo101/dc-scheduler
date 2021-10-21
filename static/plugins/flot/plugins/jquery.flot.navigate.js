!function(t){"use strict";var e={zoom:{interactive:!1,active:!1,amount:1.5},pan:{interactive:!1,active:!1,cursor:"move",frameRate:60,mode:"smart"},recenter:{interactive:!0},xaxis:{axisZoom:!0,plotZoom:!0,axisPan:!0,plotPan:!0,panRange:[void 0,void 0],zoomRange:[void 0,void 0]},yaxis:{axisZoom:!0,plotZoom:!0,axisPan:!0,plotPan:!0,panRange:[void 0,void 0],zoomRange:[void 0,void 0]}},a=t.plot.saturated,o=t.plot.browser,n=t.plot.uiConstants.SNAPPING_CONSTANT,i=t.plot.uiConstants.PANHINT_LENGTH_CONSTANT;function r(e,r){var s=null,f=!1,c="manual"===r.pan.mode,l="smartLock"===r.pan.mode,u=l||"smart"===r.pan.mode;var p,d="default",g=null,v=null,x={x:0,y:0},m=!1;function h(t,a){var n=Math.abs(t.originalEvent.deltaY)<=1?1+Math.abs(t.originalEvent.deltaY)/50:null;if(m&&M(t),e.getOptions().zoom.active)return t.preventDefault(),function(t,a,n){var i=o.getPageXY(t),r=e.offset();r.left=i.X-r.left,r.top=i.Y-r.top;var s=e.getPlaceholder().offset();s.left=i.X-s.left,s.top=i.Y-s.top;var f=e.getXAxes().concat(e.getYAxes()).filter((function(t){var e=t.box;if(void 0!==e)return s.left>e.left&&s.left<e.left+e.width&&s.top>e.top&&s.top<e.top+e.height}));0===f.length&&(f=void 0),a?e.zoomOut({center:r,axes:f,amount:n}):e.zoom({center:r,axes:f,amount:n})}(t,a<0,n),!1}function b(t){f=!0}function y(t){f=!1}function P(t){if(!f||!function(t){return 0===t.button}(t))return!1;m=!0;var a=o.getPageXY(t),n=e.getPlaceholder().offset();n.left=a.X-n.left,n.top=a.Y-n.top,0===(s=e.getXAxes().concat(e.getYAxes()).filter((function(t){var e=t.box;if(void 0!==e)return n.left>e.left&&n.left<e.left+e.width&&n.top>e.top&&n.top<e.top+e.height}))).length&&(s=void 0);var i=e.getPlaceholder().css("cursor");i&&(d=i),e.getPlaceholder().css("cursor",e.getOptions().pan.cursor),u?p=e.navigationState(a.X,a.Y):c&&(x.x=a.X,x.y=a.Y)}function w(t){if(m){var a=o.getPageXY(t),n=e.getOptions().pan.frameRate;-1!==n?!v&&n&&(v=setTimeout((function(){u?e.smartPan({x:p.startPageX-a.X,y:p.startPageY-a.Y},p,s,!1,l):c&&(e.pan({left:x.x-a.X,top:x.y-a.Y,axes:s}),x.x=a.X,x.y=a.Y),v=null}),1/n*1e3)):u?e.smartPan({x:p.startPageX-a.X,y:p.startPageY-a.Y},p,s,!1,l):c&&(e.pan({left:x.x-a.X,top:x.y-a.Y,axes:s}),x.x=a.X,x.y=a.Y)}}function M(t){if(m){v&&(clearTimeout(v),v=null),m=!1;var a=o.getPageXY(t);e.getPlaceholder().css("cursor",d),u?(e.smartPan({x:p.startPageX-a.X,y:p.startPageY-a.Y},p,s,!1,l),e.smartPan.end()):c&&(e.pan({left:x.x-a.X,top:x.y-a.Y,axes:s}),x.x=0,x.y=0)}}function Y(a){if(e.activate(),e.getOptions().recenter.interactive){var o,n=e.getTouchedAxis(a.clientX,a.clientY);e.recenter({axes:n[0]?n:null}),o=n[0]?new t.Event("re-center",{detail:{axisTouched:n[0]}}):new t.Event("re-center",{detail:a}),e.getPlaceholder().trigger(o)}}function O(t){return e.activate(),m&&M(t),!1}e.navigationState=function(t,e){var a=this.getAxes(),o={};return Object.keys(a).forEach((function(t){var e=a[t];o[t]={navigationOffset:{below:e.options.offset.below||0,above:e.options.offset.above||0},axisMin:e.min,axisMax:e.max,diagMode:!1}})),o.startPageX=t||0,o.startPageY=e||0,o},e.activate=function(){var t=e.getOptions();t.pan.active&&t.zoom.active||(t.pan.active=!0,t.zoom.active=!0,e.getPlaceholder().trigger("plotactivated",[e]))},e.zoomOut=function(t){t||(t={}),t.amount||(t.amount=e.getOptions().zoom.amount),t.amount=1/t.amount,e.zoom(t)},e.zoom=function(a){a||(a={});var o=a.center,n=a.amount||e.getOptions().zoom.amount,i=e.width(),r=e.height(),s=a.axes||e.getAxes();o||(o={left:i/2,top:r/2});var f=o.left/i,c=o.top/r,l={x:{min:o.left-f*i/n,max:o.left+(1-f)*i/n},y:{min:o.top-c*r/n,max:o.top+(1-c)*r/n}};for(var u in s)if(s.hasOwnProperty(u)){var p=s[u],d=p.options,g=l[p.direction].min,v=l[p.direction].max,x=p.options.offset;if((d.axisZoom||!a.axes)&&(a.axes||d.plotZoom)){if((g=t.plot.saturated.saturate(p.c2p(g)))>(v=t.plot.saturated.saturate(p.c2p(v)))){var m=g;g=v,v=m}if(d.zoomRange){if(v-g<d.zoomRange[0])continue;if(v-g>d.zoomRange[1])continue}var h=t.plot.saturated.saturate(x.below-(p.min-g)),b=t.plot.saturated.saturate(x.above-(p.max-v));d.offset={below:h,above:b}}}e.setupGrid(!0),e.draw(),a.preventEvent||e.getPlaceholder().trigger("plotzoom",[e,a])},e.pan=function(o){var n={x:+o.left,y:+o.top};isNaN(n.x)&&(n.x=0),isNaN(n.y)&&(n.y=0),t.each(o.axes||e.getAxes(),(function(t,e){var i=e.options,r=n[e.direction];if((i.axisPan||!o.axes)&&(i.plotPan||o.axes)){var s=e.p2c(i.panRange[0])-e.p2c(e.min),f=e.p2c(i.panRange[1])-e.p2c(e.max);if(void 0!==i.panRange[0]&&r>=f&&(r=f),void 0!==i.panRange[1]&&r<=s&&(r=s),0!==r){var c=a.saturate(e.c2p(e.p2c(e.min)+r)-e.c2p(e.p2c(e.min))),l=a.saturate(e.c2p(e.p2c(e.max)+r)-e.c2p(e.p2c(e.max)));isFinite(c)||(c=0),isFinite(l)||(l=0),i.offset={below:a.saturate(c+(i.offset.below||0)),above:a.saturate(l+(i.offset.above||0))}}}})),e.setupGrid(!0),e.draw(),o.preventEvent||e.getPlaceholder().trigger("plotpan",[e,o])},e.recenter=function(a){t.each(a.axes||e.getAxes(),(function(t,e){a.axes?"x"===this.direction?e.options.offset={below:0}:"y"===this.direction&&(e.options.offset={above:0}):e.options.offset={below:0,above:0}})),e.setupGrid(!0),e.draw()};var X=null,T={x:0,y:0};e.smartPan=function(t,o,i,r,s){var f,c,l,u,p,d,v=!!s||function(t){return Math.abs(t.y)<n&&Math.abs(t.x)>=n||Math.abs(t.x)<n&&Math.abs(t.y)>=n}(t),x=e.getAxes();(function(t){return Math.abs(t.x)>0&&Math.abs(t.y)>0})(t=s?function(t){switch(!X&&Math.max(Math.abs(t.x),Math.abs(t.y))>=n&&(X=Math.abs(t.x)<Math.abs(t.y)?"y":"x"),X){case"x":return{x:t.x,y:0};case"y":return{x:0,y:t.y};default:return{x:0,y:0}}}(t):function(t){return Math.abs(t.x)<n&&Math.abs(t.y)>=n?{x:0,y:t.y}:Math.abs(t.y)<n&&Math.abs(t.x)>=n?{x:t.x,y:0}:t}(t))&&(o.diagMode=!0),v&&!0===o.diagMode&&(o.diagMode=!1,function(t,e,a){var o;Object.keys(t).forEach((function(n){o=t[n],0===a[o.direction]&&(o.options.offset.below=e[n].navigationOffset.below,o.options.offset.above=e[n].navigationOffset.above)}))}(x,o,t)),g=v?{start:{x:o.startPageX-e.offset().left+e.getPlotOffset().left,y:o.startPageY-e.offset().top+e.getPlotOffset().top},end:{x:o.startPageX-t.x-e.offset().left+e.getPlotOffset().left,y:o.startPageY-t.y-e.offset().top+e.getPlotOffset().top}}:{start:{x:o.startPageX-e.offset().left+e.getPlotOffset().left,y:o.startPageY-e.offset().top+e.getPlotOffset().top},end:!1},isNaN(t.x)&&(t.x=0),isNaN(t.y)&&(t.y=0),i&&(x=i),Object.keys(x).forEach((function(e){if(c=x[e],l=c.min,u=c.max,f=c.options,d=t[c.direction],p=T[c.direction],(f.axisPan||!i)&&(i||f.plotPan)){var o=p+c.p2c(f.panRange[0])-c.p2c(l),n=p+c.p2c(f.panRange[1])-c.p2c(u);if(void 0!==f.panRange[0]&&d>=n&&(d=n),void 0!==f.panRange[1]&&d<=o&&(d=o),0!==d){var r=a.saturate(c.c2p(c.p2c(l)-(p-d))-c.c2p(c.p2c(l))),s=a.saturate(c.c2p(c.p2c(u)-(p-d))-c.c2p(c.p2c(u)));isFinite(r)||(r=0),isFinite(s)||(s=0),c.options.offset.below=a.saturate(r+(c.options.offset.below||0)),c.options.offset.above=a.saturate(s+(c.options.offset.above||0))}}})),T=t,e.setupGrid(!0),e.draw(),r||e.getPlaceholder().trigger("plotpan",[e,t,i,o])},e.smartPan.end=function(){g=null,X=null,T={x:0,y:0},e.triggerRedrawOverlay()},e.getTouchedAxis=function(t,a){var o=e.getPlaceholder().offset();return o.left=t-o.left,o.top=a-o.top,e.getXAxes().concat(e.getYAxes()).filter((function(t){var e=t.box;if(void 0!==e)return o.left>e.left&&o.left<e.left+e.width&&o.top>e.top&&o.top<e.top+e.height}))},e.hooks.drawOverlay.push((function(t,e){if(g){e.strokeStyle="rgba(96, 160, 208, 0.7)",e.lineWidth=2,e.lineJoin="round";var a,o,n=Math.round(g.start.x),r=Math.round(g.start.y);if(s?"x"===s[0].direction?(o=Math.round(g.start.y),a=Math.round(g.end.x)):"y"===s[0].direction&&(a=Math.round(g.start.x),o=Math.round(g.end.y)):(a=Math.round(g.end.x),o=Math.round(g.end.y)),e.beginPath(),!1===g.end)e.moveTo(n,r-i),e.lineTo(n,r+i),e.moveTo(n+i,r),e.lineTo(n-i,r);else{var f=r===o;e.moveTo(n-(f?0:i),r-(f?i:0)),e.lineTo(n+(f?0:i),r+(f?i:0)),e.moveTo(n,r),e.lineTo(a,o),e.moveTo(a-(f?0:i),o-(f?i:0)),e.lineTo(a+(f?0:i),o+(f?i:0))}e.stroke()}})),e.hooks.bindEvents.push((function(t,e){var a=t.getOptions();a.zoom.interactive&&e.mousewheel(h),a.pan.interactive&&(t.addEventHandler("dragstart",P,e,0),t.addEventHandler("drag",w,e,0),t.addEventHandler("dragend",M,e,0),e.bind("mousedown",b),e.bind("mouseup",y)),e.dblclick(Y),e.click(O)})),e.hooks.shutdown.push((function(t,e){e.unbind("mousewheel",h),e.unbind("mousedown",b),e.unbind("mouseup",y),e.unbind("dragstart",P),e.unbind("drag",w),e.unbind("dragend",M),e.unbind("dblclick",Y),e.unbind("click",O),v&&clearTimeout(v)}))}t.plot.plugins.push({init:function(t){t.hooks.processOptions.push(r)},options:e,name:"navigate",version:"1.3"})}(jQuery);