// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_holder").hide();
	$("#success_create").hide();
	$("#error_holder").hide();
	$("#error_query").hide();
	
	$scope.queryAllFile = function(){

		appFactory.queryAllFile(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_file = array;
		});
	}

	$scope.queryFile = function(){

		var id = $scope.file_id;

		appFactory.queryFile(id, function(data){
			$scope.query_file = data;

			if ($scope.query_file == "Could not locate file"){
				console.log()
				$("#error_query").show();
			} else{
				$("#error_query").hide();
			}
		});
	}

	$scope.recordFile = function(){

		appFactory.recordFile($scope.file, function(data){
			$scope.create_file = data;
			$("#success_create").show();
		});
	}

	$scope.changeFilehash = function(){

		appFactory.changeFilehash($scope.filehash, function(data){
			$scope.change_filehash = data;
			if ($scope.change_filehash == "Error: no file found"){
				$("#error_holder").show();
				$("#success_holder").hide();
			} else{
				$("#success_holder").show();
				$("#error_holder").hide();
			}
		});
	}

});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.queryAllFile = function(callback){

    	$http.get('/get_all_file/').success(function(output){
			callback(output)
		});
	}

	factory.queryFile = function(id, callback){
    	$http.get('/get_file/'+id).success(function(output){
			callback(output)
		});
	}

	factory.recordFile = function(data, callback){

		var file = data.id + "-" + data.filehash + "-" + data.timestamp;

    	$http.get('/add_file/'+file).success(function(output){
			callback(output)
		});
	}

	factory.changeFilehash = function(data, callback){

		var filehash = data.id + "-" + data.filehash;

    	$http.get('/change_filehash/'+filehash).success(function(output){
			callback(output)
		});
	}

	return factory;
});


