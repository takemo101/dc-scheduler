Package.describe({name:"dangrossman:bootstrap-daterangepicker",version:"3.1.0",summary:"Date range picker component",git:"https://github.com/dangrossman/daterangepicker",documentation:"README.md"}),Package.onUse((function(e){e.versionsFrom("METEOR@0.9.0.1"),e.use("momentjs:moment@2.22.1",["client"]),e.use("jquery@3.3.1",["client"]),e.addFiles("daterangepicker.js",["client"]),e.addFiles("daterangepicker.css",["client"])}));