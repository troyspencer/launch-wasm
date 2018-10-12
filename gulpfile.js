var gulp = require('gulp'),
  browserSync = require('browser-sync').create(),
  exec = require('child_process').exec

  
// wasm watch
gulp.task('wasm-watch', ['wasm-build'], function() {
  setTimeout(function() { browserSync.reload(); }, 300); 
})


// wasm build
gulp.task('wasm-build', function(cb) {
  exec('env GOOS=js GOARCH=wasm go build -o server/static/build/main.wasm go/main.go', function(err, stdout, stderr) {
    if (stdout) {
      console.log(stdout)
    }
    if (stderr) {
      console.log(stderr);
    }
    cb(err)
  })
})


// Static server
gulp.task('browsersync', function() {
  browserSync.init({
    "callbacks": {
      ready: function(err, bs) {
        bs.utils.serveStatic.mime.define({ 'application/wasm': ['wasm'] });
      }
    },
    server: {
        baseDir: "./server/static"
    },
    "browser": "google chrome",
    "open": false
  });

  gulp.watch("./go/*.go", ['wasm-watch'])
});

gulp.task('default', ['browsersync'])