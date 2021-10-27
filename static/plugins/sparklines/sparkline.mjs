export default class t{constructor(i,o={}){this.element=i,this.options={...t.options,...o},this.element.innerHTML="<canvas></canvas>",this.canvas=this.element.firstChild,this.context=this.canvas.getContext("2d"),this.ratio=window.devicePixelRatio||1,this.options.tooltip&&(this.canvas.style.position="relative",this.canvas.addEventListener("mousemove",(t=>{const i=t.offsetX||t.layerX||0,o=(this.options.width-2*this.options.dotRadius)/(this._points.length-1),n=(s=0,e=Math.round((i-this.options.dotRadius)/o),h=this._points.length-1,Math.max(s,Math.min(e,h)));var s,e,h;this.canvas.title=this.options.tooltip(this._points[n],n,this._points)}),!1))}set points(t){this.draw(t)}get points(){return this._points}draw(t=[]){this._points=t,this.canvas.width=this.options.width*this.ratio,this.canvas.style.width=`${this.options.width}px`;const i=this.options.height||this.element.offsetHeight;this.canvas.height=i*this.ratio,this.canvas.style.height=`${i}px`;const o=this.options.lineWidth*this.ratio,n=Math.max(this.options.dotRadius*this.ratio,o/2),s=Math.max(this.options.dotRadius*this.ratio,o/2),e=this.canvas.width-2*n,h=this.canvas.height-2*s,l=Math.min.apply(Math,t),a=Math.max.apply(Math,t),c=null!=this.options.minValue?this.options.minValue:Math.min(l,null!=this.options.maxMinValue?this.options.maxMinValue:l),p=null!=this.options.maxValue?this.options.maxValue:Math.max(a,null!=this.options.minMaxValue?this.options.minMaxValue:a);let r=n,x=n,d=n;const m=i=>p===c?s+h/2:s+h-(t[i]-c)/(p-c)*h,u=e/(t.length-1),g=(t,i,o)=>{t&&(this.context.save(),this.context.strokeStyle=t.color||"black",this.context.lineWidth=(t.width||1)*this.ratio,this.context.globalAlpha=t.alpha||1,this.context.beginPath(),this.context.moveTo("right"!=t.direction?n:i,o),this.context.lineTo("left"!=t.direction?e+n:i,o),this.context.stroke(),this.context.restore())},f=(t,i,o,n)=>{this.context.beginPath(),this.context.fillStyle=t,this.context.arc(o,n,this.options.dotRadius*this.ratio,0,2*Math.PI,!1),this.context.fill(),g(i,o,n)};if(this.context.save(),this.context.strokeStyle=this.options.lineColor,this.context.fillStyle=this.options.lineColor,this.context.lineWidth=o,this.context.lineCap="round",this.context.lineJoin="round",this.options.fillBelow&&t.length>1){this.context.save(),this.context.beginPath(),this.context.moveTo(d,m(0));for(let i=1;i<t.length;i++)d+=u,r=t[i]==l?d:r,x=t[i]==a?d:x,this.context.lineTo(d,m(i));this.context.lineTo(e+n,h+s+o/2),this.context.lineTo(n,h+s+o/2),this.context.fill(),this.options.fillLighten>0?(this.context.fillStyle="white",this.context.globalAlpha=this.options.fillLighten,this.context.fill(),this.context.globalAlpha=1):this.options.fillLighten<0&&(this.context.fillStyle="black",this.context.globalAlpha=-this.options.fillLighten,this.context.fill()),this.context.restore()}d=n,this.context.beginPath(),this.context.moveTo(d,m(0));for(let i=1;i<t.length;i++)d+=u,this.context.lineTo(d,m(i));this.context.stroke(),this.context.restore(),g(this.options.bottomLine,0,s),g(this.options.topLine,0,h+s+o/2),f(this.options.startColor,this.options.startLine,n+(1==t.length?e/2:0),m(0)),f(this.options.endColor,this.options.endLine,n+(1==t.length?e/2:e),m(t.length-1)),f(this.options.minColor,this.options.minLine,r+(1==t.length?e/2:0),m(t.indexOf(l))),f(this.options.maxColor,this.options.maxLine,x+(1==t.length?e/2:0),m(t.indexOf(a)))}static init(i,o){return new t(i,o)}static draw(i,o,n){const s=new t(i,n);return s.draw(o),s}}t.options={width:100,height:null,lineColor:"black",lineWidth:1.5,startColor:"transparent",endColor:"black",maxColor:"transparent",minColor:"transparent",minValue:null,maxValue:null,minMaxValue:null,maxMinValue:null,dotRadius:2.5,tooltip:null,fillBelow:!0,fillLighten:.5,startLine:!1,endLine:!1,minLine:!1,maxLine:!1,bottomLine:!1,topLine:!1,averageLine:!1};