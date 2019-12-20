//SPDX-License-Identifier: Apache-2.0

var filehash = require('./controller.js');

module.exports = function(app){

  app.get('/get_file/:id', function(req, res){
    filehash.get_file(req, res);
  });

  app.get('/add_file/:filehash', function(req, res){
    filehash.add_file(req, res);
  });

  app.get('/get_all_file', function(req, res){
    filehash.get_all_file(req, res);
  });
  
  app.get('/change_filehash/:filehash', function(req, res){
    filehash.change_filehash(req, res);
  });

}