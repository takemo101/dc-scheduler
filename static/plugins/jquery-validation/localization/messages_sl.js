!function(e){"function"==typeof define&&define.amd?define(["jquery","../jquery.validate"],e):"object"==typeof module&&module.exports?module.exports=e(require("jquery")):e(jQuery)}((function(e){return e.extend(e.validator.messages,{required:"To polje je obvezno.",remote:"Prosimo popravite to polje.",email:"Prosimo vnesite veljaven email naslov.",url:"Prosimo vnesite veljaven URL naslov.",date:"Prosimo vnesite veljaven datum.",dateISO:"Prosimo vnesite veljaven ISO datum.",number:"Prosimo vnesite veljavno število.",digits:"Prosimo vnesite samo števila.",creditcard:"Prosimo vnesite veljavno številko kreditne kartice.",equalTo:"Prosimo ponovno vnesite vrednost.",extension:"Prosimo vnesite vrednost z veljavno končnico.",maxlength:e.validator.format("Prosimo vnesite največ {0} znakov."),minlength:e.validator.format("Prosimo vnesite najmanj {0} znakov."),rangelength:e.validator.format("Prosimo vnesite najmanj {0} in največ {1} znakov."),range:e.validator.format("Prosimo vnesite vrednost med {0} in {1}."),max:e.validator.format("Prosimo vnesite vrednost manjše ali enako {0}."),min:e.validator.format("Prosimo vnesite vrednost večje ali enako {0}.")}),e}));