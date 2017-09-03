"use strict";
var app = angular.module("app", []);

app.config(function ($logProvider, $interpolateProvider) {
  
  $logProvider.debugEnabled(true);

  $interpolateProvider.startSymbol('[['); 
  $interpolateProvider.endSymbol(']]'); 
});

app.run(function ($rootScope, $interval, $window) {
  console.log("app run");
});
